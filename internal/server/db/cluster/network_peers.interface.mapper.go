//go:build linux && cgo && !agent

package cluster

import "context"

// NetworkPeerGenerated is an interface of generated methods for NetworkPeer.
type NetworkPeerGenerated interface {
	// CreateNetworkPeerConfig adds new network_peer Config to the database.
	// generator: network_peer Create
	CreateNetworkPeerConfig(ctx context.Context, db dbtx, networkPeerID int64, config map[string]string) error

	// CreateNetworkPeer adds a new network_peer to the database.
	// generator: network_peer Create
	CreateNetworkPeer(ctx context.Context, db dbtx, object NetworkPeer) (int64, error)
}
