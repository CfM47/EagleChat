# EagleChat

REEEEE 🦅🦅🦅💥💥🦅🦅💥💥🦅

                           .a@@@@@#########@@@@a.
                       .a@@######@@@mm@@mm######@@@a.
                  .a####@@@@@@@@@@@@@@@@@@@mm@@##@@v;%%,.
               .a###v@@@@@@@@vvvvvvvvvvvvvv@@@@#@v;%%%vv%%,
            .a##vv@@@@@@@@vv%%%%;S,  .S;%%vv@@#v;%%'/%vvvv%;
          .a##@v@@@@@vv%%vvvvvv%%;SssS;%%vvvv@v;%%./%vvvvvv%;
        ,a##vv@@@vv%%%@@@@@@@@@@@@mmmmmmmmmvv;%%%%vvvvvvvvv%;
        .a##@@@@@@@@@@@@@@@@@@@@@@@mmmmmvv;%%%%%vvvvvvvvvvv%;
       ###vv@@@v##@v@@@@@@@@@@mmv;%;%;%;%;%;%;%;%;%;%;%,%vv%'
      a#vv@@@@v##v@@@@@@@@###@@@@@%v%v%v%v%v%v%v%      ;%%;'
     ',a@@@@@@@v@@@@@@@@v###v@@@nvnvnvnvnvnvnvnv'     .%;'
     a###@@@@@@@###v@@@v##v@@@mnmnmnmnmnmnmnmnmn.     ;'
    ,###vv@@@@v##v@@@@@@v@@@@v##v@@@@@v###v@@@##@.
    ###vv@@@@@@v@@###v@@@@@@@@v@@@@@@v##v@@@v###v@@.

A chat messaging application with a decentralization and freedom premise 

## Overview

EagleChat is an academic project developed for the *Distributed Systems* course.  
The objective of the project is to design, deploy, and manage a distributed application running on a Docker Swarm cluster, with special emphasis on service orchestration, and secure communication between components.

The system is designed to be reproducible, automated, and easy to manage through a set of CLI scripts that encapsulate the underlying Docker and networking logic.


## Architecture

The system is deployed on a Docker Swarm environment and assumes execution on manager nodes, which are responsible for:

- Service orchestration
- Container lifecycle management
- Overlay network management
- Secure distribution of credentials

### Core Components

- **ID Managers**  
  One or more manager services responsible for identity management and coordination.
  
- **Clients**  
  Interactive client containers that connect to the system through the Swarm network.

- **Overlay Network**  
  All services communicate over a dedicated Docker overlay network.

- **PKI / Credentials**  
  A custom Certificate Authority (CA) is generated automatically, along with per-manager credentials.


## Usage

First run `go work vendor` it is needed for containers to be built efficiently.

To simplify interaction with the system and avoid direct Docker commands, the project provides several helper scripts.

### `mgr.sh` – Swarm Management CLI

`mgr.sh` is the main entry point for managing the EagleChat environment on Docker Swarm.  
It acts as a thin CLI wrapper over shared logic implemented in `scripts/common.sh`.

#### Usage

```bash
./mgr.sh <command> [arguments]
````

#### Available Commands

* **initialize [--build]**
  Initializes the Swarm environment.
  Images are built only if `--build` is specified.

* **create-manager `<name>`**
  Creates and starts a new ID Manager container.

* **start-client `<name>`**
  Starts an interactive client container connected to the Swarm network.

* **delete `<name>`**
  Stops and removes a specific container.

* **ls**
  Lists all managed containers (short names).

* **cleanup**
  Completely tears down the environment:

  * Stops and removes containers
  * Removes the overlay network
  * Deletes generated credentials

This script is intended to be run **from a Swarm manager node**.


## Local Testing with Docker Compose

For development and testing purposes, the project also supports a **Docker Compose–based setup**, which runs locally without requiring a Swarm cluster.

### `run-compose.sh`

This script:

1. Removes any existing credentials.
2. Generates a new CA.
3. Generates credentials for one or more ID Managers.
4. Starts the Docker Compose environment.

#### Usage

```bash
./run-compose.sh
```

Single-manager scenario (default).

```bash
./run-compose.sh multi-manager
```

Multi-manager scenario using Compose profiles.


### `stop-compose.sh`

Stops the Docker Compose environment and performs cleanup.

#### Usage

```bash
./stop-compose.sh
```

or

```bash
./stop-compose.sh multi-manager
```

This script also:

* Removes generated credentials
* Verifies that the Compose network has been properly removed


## Security and Credentials

* All credentials are **generated automatically**.
* A temporary certificate directory is used during execution.
* Cleanup scripts ensure that no sensitive material remains after shutdown.

This approach ensures reproducibility and minimizes configuration errors.


## Requirements

* Docker
* Docker Swarm (initialized on manager nodes)
* Docker Compose (for local testing)
* Bash (for management scripts)
