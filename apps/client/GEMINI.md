# EagleChat Client

## Project Overview

This project is the client-side application for EagleChat, a peer-to-peer (P2P) chat system. It is written in Go and functions as a sophisticated peer capable of direct communication, service discovery, and resilient, guaranteed messaging.

The application is structured as a Go module (`eaglechat/apps/client`) and follows common Go project conventions. The core logic is implemented as a set of composable services within the `internal/middleware` directory.

## Development Workflow

To ensure accuracy and efficiency, all development tasks will follow a structured, four-step workflow:

1. **Investigate:** Before making any changes, thoroughly analyze the relevant parts of the codebase. This includes reading interface definitions, struct definitions, and the logic of existing functions to build a complete understanding of the current state.
2. **Plan:** Based on the investigation, create a detailed, step-by-step plan. This plan must outline the proposed changes, new components, and any modifications to existing structures. The plan should also include a section for assumptions, clarifying questions, and potential issues.
3. **Approve:** Submit the plan for review and approval. Do not proceed until the plan is confirmed.
4. **Implement:** Execute the approved plan. Read files immediately before modifying them to prevent conflicts with any other changes.

## Core Middleware Components

The `internal/middleware` directory contains several key services that compose the client's core functionality:

* **`p2pconn`**: Provides a reliable, ACK-based, guaranteed-consumption P2P connection over TCP. It manages the low-level wire protocol for sending data and receiving application-level acknowledgements.

* **`p2pconnpool`**: Manages a pool of `p2pconn` instances. It caches outgoing connections to avoid redundant dialing and aggregates all incoming messages from all connections into a single channel for consumption.

* **`idmanagerconn`**: Handles HTTP communication with a single ID Manager service. It includes a health check on construction to ensure the service is available.

* **`idmanagerpool`**: Manages the discovery of ID Manager services on the local network via multicast announcements. It maintains a pool of healthy, available ID Manager connections.

* **`messageCache`**: A repository for storing pending messages that could not be delivered immediately, ensuring no data is lost due to network failures or offline peers.

* **`Middleware`**: The central orchestrator that ties all other components together. It will implement the primary application logic, including message routing (for self vs. others) and the store-and-forward mechanism for offline messages.

## Building and Running

**To build the client:**

```bash
go build -o eaglechat_client ./cmd
```

**To run tests:**

```bash
go test ./...
```

## Development Conventions

* **Structure:** The project follows the standard Go project layout.
* **Dependencies:** Dependencies are managed using Go modules (`go.mod` and `go.sum`).
* **Testing:** All new functionality should be accompanied by tests where feasible.

