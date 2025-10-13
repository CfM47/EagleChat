# Technical Debt

## P2P Connection Acknowledgement Model

**Date:** 2025-10-10

**Issue:**
The current ACK (acknowledgement) protocol within the `p2pconn` layer confirms message delivery as soon as the message is **consumed** by the application's receiving channel. It does not wait for the application to confirm that the message has been fully **handled** (e.g., processed, saved to a database, etc.).

**Risk:**
If the application crashes or encounters a fatal error *after* consuming the message from the channel but *before* it finishes processing it, the sender will have received an ACK for a message that was effectively lost. This can lead to data inconsistency, where the sender believes a message was delivered and processed when it was not.

**Future Improvement:**
A more robust solution would be to implement a "Commit-Ack" model. This would involve changing the `Receive()` channel to pass a struct containing both the message payload and an `Ack()` callback function. The application would be responsible for calling `Ack()` only after it has successfully handled the message, providing a true end-to-end guarantee.
