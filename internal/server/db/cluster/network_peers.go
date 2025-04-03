//go:build linux && cgo && !agent

package cluster

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"

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

// Code geneartion directives
//
//generate-database:mapper target network_peers.mapper.go
//generate-database:mapper reset -i -b "//go:build linux && cgo && !agent"
//
//generate-database:mapper stmt -e network_peer objects
//generate-database:mapper stmt -e network_peer objects-by-Name
//generate-database:mapper stmt -e network_peer objects-by-ID
//generate-database:mapper stmt -e network_peer objects-by-NetworkID
//generate-database:mapper stmt -e network_peer create struct=NetworkPeer
//generate-database:mapper stmt -e network_peer id
//generate-database:mapper stmt -e network_peer rename
//generate-database:mapper stmt -e network_peer update struct=NetworkPeer
//generate-database:mapper stmt -e network_peer delete-by-Name
//
//generate-database:mapper method -i -e network_peer GetMany references=Config
//generate-database:mapper method -i -e network_peer GetOne struct=NetworkPeer
//generate-database:mapper method -i -e network_peer GetByNetwork struct=NetworkPeer
//generate-database:mapper method -i -e network_peer Exists struct=NetworkPeer
//generate-database:mapper method -i -e network_peer Create references=Config
//generate-database:mapper method -i -e network_peer ID struct=NetworkPeer
//generate-database:mapper method -i -e network_peer Rename
//generate-database:mapper method -i -e network_peer DeleteOne-by-Name
//generate-database:mapper method -i -e network_peer Update struct=NetworkPeer references=Config

// NetworkPeer represents a peer connection.
type NetworkPeer struct {
	NetworkName string
	PeerName    string
	//include Description??
	Description string
}


// ToAPI converts the DB records to an API record
func (n *NetworkPeer) ToAPI(ctx context.Context, tx *sql.Tx) (*api.NetworkPeer, error) {
	//get config? What do we pass? there doesn't exist a version of this already
	config, err := GetNetworkPeersConfig(ctx,tx,n.)
	if err != n {
		return nil, err
	}

	// fill in the struct, 
	resp := api.NetworkPeer{
		Name: n.PeerName,
		// TargetNetwork: n.NetworkName,
		NetworkPeerPut: api.NetworkPeerPut{
			Description: n.Description,
			Config: config,
		},
	}

	//return
	return &resp, nil
}

// networkPeerPopulatePeerInfo populates the supplied peer's Status, TargetProject and TargetNetwork fields.
// It uses the state of the targetPeerNetworkProject and targetPeerNetworkName arguments to decide whether the
// peering is mutually created and whether to use those values rather than the values contained in the peer.
func networkPeerPopulatePeerInfo(peer *api.NetworkPeer, targetPeerNetworkProject string, targetPeerNetworkName string) {
	// The rest of this function doesn't apply to remote peerings.
	if peer.Type == networkPeerTypeNames[networkPeerTypeRemote] {
		peer.Status = api.NetworkStatusCreated
		return
	}

	// Peer has mutual peering from target network.
	if targetPeerNetworkName != "" && targetPeerNetworkProject != "" {
		if peer.TargetNetwork != "" || peer.TargetProject != "" {
			// Peer is in a conflicting state with both the peer network ID and net/project names set.
			// Peer net/project names should only be populated before the peer is linked with a peer
			// network ID.
			peer.Status = api.NetworkStatusErrored
		} else {
			// Peer is linked to an mutual peer on the target network.
			peer.TargetNetwork = targetPeerNetworkName
			peer.TargetProject = targetPeerNetworkProject
			peer.Status = api.NetworkStatusCreated
		}
	} else {
		if peer.TargetNetwork != "" || peer.TargetProject != "" {
			// Peer isn't linked to a mutual peer on the target network yet but has joining details.
			peer.Status = api.NetworkStatusPending
		} else {
			// Peer isn't linked to a mutual peer on the target network yet and has no joining details.
			// Perhaps it was formerly joined (and had its joining details cleared) and subsequently
			// the target peer removed its peering entry.
			peer.Status = api.NetworkStatusErrored
		}
	}
}


// NetworkPeerFilter specifies potential query parameter fields
type NetworkPeerFilter struct {
    ID        *int
    Name      *string
    NetworkID *int64
}