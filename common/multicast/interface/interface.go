package multicast

type MulticastNetwork interface {
	Broadcast(BroadcastMessage) error
	Announcements() <-chan BroadcastMessage
	Close() error
	Done() <-chan struct{}
}
