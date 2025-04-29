//go:build linux && cgo && !agent

package cluster

import "context"

// NetworkPeerGenerated is an interface of generated methods for NetworkPeer.
type NetworkPeerGenerated interface {
	// GetNetworkPeers returns all available network_peers.
	// generator: network_peer GetMany
	GetNetworkPeers(ctx context.Context, db dbtx, filters ...NetworkPeerFilter) ([]NetworkPeer, error)

	// GetNetworkPeer returns the network_peer with the given key.
	// generator: network_peer GetOne
	GetNetworkPeer(ctx context.Context, db dbtx, name string) (*NetworkPeer, error)

	// NetworkPeerExists checks if a network_peer with the given key exists.
	// generator: network_peer Exists
	NetworkPeerExists(ctx context.Context, db dbtx, name string) (bool, error)

	// CreateNetworkPeer adds a new network_peer to the database.
	// generator: network_peer Create
	CreateNetworkPeer(ctx context.Context, db dbtx, object NetworkPeer) (int64, error)

	// GetNetworkPeerID return the ID of the network_peer with the given key.
	// generator: network_peer ID
	GetNetworkPeerID(ctx context.Context, db tx, name string) (int64, error)

	// DeleteNetworkPeer deletes the network_peer matching the given key parameters.
	// generator: network_peer DeleteOne-by-Name
	DeleteNetworkPeer(ctx context.Context, db dbtx, name string) error

	// UpdateNetworkPeer updates the network_peer matching the given key parameters.
	// generator: network_peer Update
	UpdateNetworkPeer(ctx context.Context, db tx, name string, object NetworkPeer) error
}
