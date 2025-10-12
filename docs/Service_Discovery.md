# Service discovery


## ID Manager Discovery

ID Managers discover each other and get discovered by clients by using ANNOUNCE messages sent over multicast, as defined in the common/multicast interface.

A default id manager can be set in the environment variables of a client, this will be used as a last resort, if no managers are discovered over multicast.

## Client Discovery

Clients must register themselves in an ID Manager to be discovered by other clients.

When a client registers to an ID Manager, it sends its username, public key and ip address. The ID Manager stores this information and makes it available to other clients that query for it.

A Client must periodically announce its presence to as many ID Managers as it can, in this request, it will also notify the manager of the pending messages the client has cached. If it does not announce itself for a long time, it will be considered disconnected by ID Managers, and therefore, other clients

