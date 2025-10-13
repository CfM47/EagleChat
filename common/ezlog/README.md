# ezlog Logging Library

This guide explains how to set up and use the `ezlog` library.

## 1. Setup

First, add the required dependencies to your project:

```bash
go get github.com/rs/zerolog
go get github.com/google/uuid
```

Next, initialize the logger factory in your `main` function. This must be done once at application startup.

```go
// In cmd/my_app/main.go
package main

import (
 "eaglechat/common/ezlog"
 "eaglechat/common/ezlog/zerolog" // The implementation package
)

func main() {
 // Create the concrete factory and inject it into the ezlog API.
 factory := logging.NewFactory()
 ezlog.SetLoggerFactory(factory)

 // The application is now ready to log.
 // ...
}
```

## 2. Usage

### Creating a Logging Context

To start a traceable operation, create a new context. This is typically done at the beginning of a request or a new task.

```go
// Create a root context for the "Main" component.
ctx := ezlog.NewLoggerContext("Main")

// Retrieve the logger from the context.
logger := ezlog.Log(ctx)

// Use the logger.
logger.Info("Application starting...")

// or, more idiomatically:
ezlog.Log(ctx).Info("Application starting...")
```

### Logging in Sub-components

To log within a sub-component while keeping the same trace ID, use `WithComponentPrefix` to create a new context and pass it down.

```go
func doWork(ctx context.Context) {
    // Create a new context for the "Worker" component.
    // It inherits the trace ID from the parent context.
    workerCtx := ezlog.WithComponentPrefix(ctx, "Worker")

    ezlog.Log(workerCtx).Info("Performing some work...")
}

// In main, call the function with the parent context.
doWork(ctx)
```

## 3. Example Output

The code from the sections above will produce the following output in the log file. Note that the trace ID is the same for all related entries.

```
[a1b2c3d4-e5f6-g7h8-i9j0-k1l2m3n4o5p6] [Main] Application starting...
[a1b2c3d4-e5f6-g7h8-i9j0-k1l2m3n4o5p6] [Worker] Performing some work...
```

## 4. Configuration

The log file path defaults to `/data/log`. This can be overridden by setting the `LOGGER_PATH` environment variable.
