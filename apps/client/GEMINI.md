# EagleChat Client

## Project Overview

This project is the client-side application for EagleChat, a peer-to-peer (P2P) chat system. It is written in Go and appears to handle direct communication with other clients for sending and receiving messages.

The application is structured as a Go module (`eaglechat/apps/client`) and follows common Go project conventions, with a `cmd` directory for the main executable and an `internal` directory for the core application logic.

Key components include:
*   **P2P Connections:** The `internal/middleware/p2pconn` package suggests that the client establishes direct TCP connections with other peers.
*   **Domain Model:** The `internal/domain` directory defines the core entities and logic of the application, such as `User` and `Message`.
*   **Cryptography:** The presence of `internal/utils/simplecrypto` indicates that the application likely implements end-to-end encryption for messages.

## Building and Running

The following are inferred commands for building and running the project.

**To build the client:**
```bash
go build -o eaglechat_client ./cmd
```

**To run the client:**
The `main.go` file in the `cmd` directory suggests the client takes an IP address and port as arguments to connect to a server for status polling. However, the core P2P logic in the `internal` directory might be part of a different executable or used as a library.

To run the status poller:
```bash
./eaglechat_client <server_ip> <server_port>
```

**To run tests:**
```bash
go test ./...
```

## Development Conventions

*   **Structure:** The project follows the standard Go project layout, with public-facing code in `cmd` and internal logic in `internal`.
*   **Dependencies:** Dependencies are managed using Go modules (`go.mod` and `go.sum`).
*   **Testing:** The project includes tests, as indicated by the `_test.go` files and the use of the `testify` library. All new functionality should be accompanied by tests.
