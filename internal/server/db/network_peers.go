//go:build linux && cgo && !agent

package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	dbCluster "github.com/lxc/incus/v6/internal/server/db/cluster"
	"github.com/lxc/incus/v6/shared/api"
)

const (
	networkPeerTypeLocal = iota
	networkPeerTypeRemote
)

var networkPeerTypeNames = map[int]string{
	networkPeerTypeLocal:  "local",
	networkPeerTypeRemote: "remote",
}

// CreateNetworkPeer creates a new Network Peer and returns its ID.
// If there is a mutual peering on the target network side the both peer entries are updated to link to each other's
// repspective network ID.
// Returns the local peer ID and true if a mutual peering has been created.
func (c *ClusterTx) CreateNetworkPeer(ctx context.Context, networkID int64, info *api.NetworkPeersPost) (int64, bool, error) {
	var err error
	var localPeerID int64
	var targetPeerNetworkID int64 = -1 // -1 means no mutual peering exists.

	if info.Type == "" || info.Type == networkPeerTypeNames[networkPeerTypeLocal] {
		// Insert a new pending local peer record.
		result, err := c.tx.ExecContext(ctx, `
		INSERT INTO networks_peers
		(network_id, name, description, type, target_network_project, target_network_name)
		VALUES (?, ?, ?, ?, ?, ?)
		`, networkID, info.Name, info.Description, networkPeerTypeLocal, info.TargetProject, info.TargetNetwork)
		if err != nil {
			return -1, false, err
		}

		localPeerID, err = result.LastInsertId()
		if err != nil {
			return -1, false, err
		}
	} else if info.Type == networkPeerTypeNames[networkPeerTypeRemote] {
		// Get the target integration.
		if info.TargetIntegration == "" {
			return -1, false, fmt.Errorf("Missing network integration name")
		}

		networkIntegration, err := dbCluster.GetNetworkIntegration(ctx, c.tx, info.TargetIntegration)
		if err != nil {
			return -1, false, err
		}

		// Insert a new remote peer record.
		result, err := c.tx.ExecContext(ctx, `
		INSERT INTO networks_peers
		(network_id, name, description, type, target_network_integration_id)
		VALUES (?, ?, ?, ?, ?)
		`, networkID, info.Name, info.Description, networkPeerTypeRemote, networkIntegration.ID)
		if err != nil {
			return -1, false, err
		}

		localPeerID, err = result.LastInsertId()
		if err != nil {
			return -1, false, err
		}
	} else {
		return -1, false, fmt.Errorf("Invalid network peer type %q", info.Type)
	}

	// Save config.
	err = dbCluster.CreateNetworkPeerConfig(ctx, c.tx, localPeerID, info.Config)
	if err != nil {
		return -1, false, err
	}

	if info.Type == "" || info.Type == networkPeerTypeNames[networkPeerTypeLocal] {
		// Check if we are creating a mutual peering of an existing peer and if so then update both sides
		// with the respective network IDs. This query looks up our network peer's network name and project
		// name and then checks if there are any unlinked (target_network_id IS NULL) peers that have
		// matching target network and project names for the network this peer belongs to. If so then it
		// returns the target peer's ID and network ID. This can then be used to update both our local peer
		// and the target peer itself with the respective network IDs of each side.
		q := `
		SELECT
			target_peer.id,
			target_peer.network_id
		FROM networks_peers AS local_peer
		JOIN networks AS local_network
			ON local_network.id = local_peer.network_id
		JOIN projects AS local_project
			ON local_project.id = local_network.project_id
		JOIN networks_peers AS target_peer
			ON target_peer.target_network_name = local_network.name
			AND target_peer.target_network_project = local_project.name
		JOIN networks AS target_peer_network
			ON target_peer.network_id = target_peer_network.id
			AND target_peer_project.name = ?
		JOIN projects AS target_peer_project
			ON target_peer_network.project_id = target_peer_project.id
			AND target_peer_network.name = ?
		WHERE
			local_peer.network_id = ?
			AND local_peer.id = ?
			AND target_peer.target_network_id IS NULL
		LIMIT 1
		`

		var targetPeerID int64 = -1

		err = c.tx.QueryRowContext(ctx, q, info.TargetProject, info.TargetNetwork, networkID, localPeerID).Scan(&targetPeerID, &targetPeerNetworkID)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return -1, false, fmt.Errorf("Failed looking up mutual peering: %w", err)
		} else if err == nil {
			peerNetworkMap := map[int64]struct {
				localNetworkID      int64
				targetPeerNetworkID int64
			}{
				localPeerID: {
					localNetworkID:      networkID,
					targetPeerNetworkID: targetPeerNetworkID,
				},
				targetPeerID: {
					localNetworkID:      targetPeerNetworkID,
					targetPeerNetworkID: networkID,
				},
			}

			// A mutual peering has been found, update both sides with their respective network IDs
			// and clear the joining target project and network names.
			for peerID, peerMap := range peerNetworkMap {
				_, err := c.tx.Exec(`
				UPDATE networks_peers SET
					target_network_id = ?,
					target_network_project = NULL,
					target_network_name = NULL
				WHERE networks_peers.network_id = ? AND networks_peers.id = ?
				`, peerMap.targetPeerNetworkID, peerMap.localNetworkID, peerID)
				if err != nil {
					return -1, false, fmt.Errorf("Failed updating mutual peering: %w", err)
				}
			}
		}
	}

	return localPeerID, targetPeerNetworkID > -1, nil
}

// NetworkPeer represents a peer connection.
type NetworkPeer struct {
	NetworkName string
	PeerName    string
}
