# Errors

Guidelines on how to handle errors

## Error Wrapping

To ensure that error information is not lost as it propagates up the call stack, we must use error wrapping. When returning an error from a function, wrap it with a descriptive message using `fmt.Errorf` and the `%w` verb.

This preserves the original error, allowing callers to inspect the entire error chain.

### Example

```go
import (
    "fmt"
    "os"
)

func readFile() error {
    _, err := os.Open("a_file_that_does_not_exist")
    if err != nil {
        return fmt.Errorf("failed to open file: %w", err)
    }
    return nil
}
```

## Logging Errors

All errors should be logged using the centralized `ezlog` package. This ensures consistent, traceable logging across the application.

When an error is handled (i.e., not returned to the caller), it should be logged with an appropriate severity level (e.g., `Error`, `Warn`).

Refer to the `ezlog` documentation in `common/ezlog/README.md` for detailed setup and usage instructions.

### Example

```go
import (
    "context"
    "eaglechat/common/ezlog"
    "errors"
)

func doSomething() error {
    return errors.New("something went wrong")
}

func processRequest(ctx context.Context) {
    if err := doSomething(); err != nil {
        ezlog.Log(ctx).Error("Failed to do something: %v", err)
        // Handle the error, but don't propagate it further
    }
}
```
