# Multicast protocol

The message composing, sending and listening for interaction through multicast
is implemented in /common/multicast, the interface defines messages and the
network abstraction

## ID Manager Announcements

ID Managers should announce themselves over the network

- Periodically, every few seconds
- Immediately after receiving a REGISTER message from a client

## ID Manager Updates

ID managers should send an UPDATED message through multicast when a user is registered
with them, to trigger other id managers to initiate gossip
