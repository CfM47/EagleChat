# Stage 1: Build the Go binary
FROM golang:latest AS builder

# Set the working directory for the build
WORKDIR /app

# Copy all go.mod and go.sum files to download dependencies efficiently
# This leverages Docker's layer caching.
COPY apps/id_manager/go.mod apps/id_manager/go.sum ./apps/id_manager/
COPY apps/client/go.mod apps/client/go.sum ./apps/client/
COPY packages/crypto/go.mod ./packages/crypto/

# Set the workdir to the specific app we are building to download its specific dependencies
WORKDIR /app/apps/id_manager
RUN go mod download

# Change workdir back to the root to copy all source code
WORKDIR /app
COPY . .

# Build the application from within its directory
WORKDIR /app/apps/id_manager
RUN CGO_ENABLED=0 GOOS=linux go build -o /id_manager ./cmd/main.go

# Stage 2: Create the final, minimal image
FROM alpine:latest

COPY --from=builder /id_manager /id_manager

EXPOSE 8080
CMD ["/id_manager"]
