# Pending message handling

The id manager has a set of (recipient_id, message_id), for the pending messages.

## User connection

When a user connects, he gives the id manager the pending messages it has
stored, the id manager tells the user:

- which messages they can delete
- which other users that are recipients to their messages are connected,
sends their data

## Message sending failure

When a user attempts to send a message to a user that is not connected,
or sending a message over a p2p connection fails, and the client assumes the
user is disconnected

- the user stores the message
- tells the id manager it has a message pending, and asks for n other users' ip,
then sends the message to those users for storage

## Message Caching

If user a connects with id manager A and tells it there's a message from a to
some other user b then loses connection with that id manager and connects to
another one B that does not yet have this information, a cannot assume the
message has been received by b because it did not find it in B's pending messages

### Three levels of caching

- **Up for deletion**: These messages may be deleted when not found in an id manager's
pending message list, or after a long period of time
  - Pending messages from other clients will eventually become up for deletion

- **Temporarily immune**: These messages are guaranteed to exist in cache until
a given time, after that they become Up for deletion
  - Pending messages from other clients start being Temporarily Immune

- **Immune**: These messages are guaranteed to always exist in cache until
explicitly removed
  - Pending messages from oneself are Immune

### P2P message delivery validation

Immune pending messages are only deleted when p2p confirmation is established:
when a p2p connection is established, pending messages from oneself to the target
are announced, and respective actions (deletion if delivered, delivery and then
validation again if not) are taken
