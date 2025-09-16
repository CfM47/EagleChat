# EagleChat

REEEEE 🦅🦅🦅💥💥🦅🦅💥💥🦅

## ID manager

### Persistent messages (not connected)
- [client] hey! i'm connected, this is my ip
- [server] hey! ok, these people want to send you messages
- [client] ok, give me the ip of x, y, z...
- [server] only z is connected, this is the ip
- [client]... (connected)
- [client] i received message sklfjakhjfg
- [server] ok, everyone, you can delete ^

### Connected
- [client] i'm connected, this is my ip, i want to chat with x
- [server] ok, x is connected, this is the ip
- direct connection

### Interface

- i'm registering, this is my pk
- i'm connected (staying connected => keep alive)
    this is my ip
    these are the messages i'm storing
    | **response**: you have these messages waiting for you, and the ones who have them, you can delete these messages you were storing **PROBLEM**
- i'm disconnecting
- i want to speak with x
    - x is connected, here is the ip and their pk
    - x is not connected, send messages to these ips for storage
- here's a message for y (if y is you, decrypt it)
- you can delete this message

