package ports

import "context"

// Notifier defines the interface for notifying peers of updates.
// This abstraction helps break circular dependencies.
type Notifier interface {
	NotifyPeersOfUpdate(ctx context.Context) error
	// TriggerSyncFromPeer is used to initiate a data pull from a specific peer.
	// This is typically called when a notification is received from that peer.
	TriggerSyncFromPeer(ctx context.Context, peerAddress string) error
}
