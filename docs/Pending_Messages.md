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
