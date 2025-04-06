//go:build linux && cgo && !agent

package cluster

import (
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
//generate-database:mapper stmt -e network_peer objects-by-NetworkID
//generate-database:mapper stmt -e network_peer create struct=NetworkPeer
//generate-database:mapper stmt -e network_peer update struct=NetworkPeerUpdate
//generate-database:mapper stmt -e network_peer update struct=NetworkPeerTarget
//generate-database:mapper stmt -e network_peer delete-by-NetworkID-and-ID
//
//generate-database:mapper method -i -e network_peer Create struct=NetworkPeer references=Config

// NOTE: config.go already has defined Config including Objects, Create, and Delete

// NetworkPeer is our main entity we are quering for
type NetworkPeer struct {
	ID          int `db:"primary=yes"`
	NetworkID   int `db:"primary=yes"`
	NetworkName string
	PeerName    string
	Description string
	// not exhaustive at all, my actually link other tables
}

// this struct defines what we will update for UpdateNeworkPeer()
type NetworkPeerUpdate struct {
	Description string `db:"sql=networks_peers.description"`
}

// this struct defines what will be updated in CreateNetworkPeer()
// need to be ptrs to set to NULL, when accessing struct field they get
// automatically dereferenced
type NetworkPeerTarget struct {
	TargetNetworkID      *int    `db:"sql=networks_peers.target_network_id"`
	TargetNetworkProject *string `db:"sql=networks_peers.target_network_project"`
	TargetNetworkName    *string `db:"sql=networks_peers.target_network_name"`
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

// all fields used in WHERE clauses for networks_peers

// NetworkPeerFilter specifies potential query parameter fields
type NetworkPeerFilter struct {
	ID          *int
	NetworkName *string
	NetworkID   *int64
}
