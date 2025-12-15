# Stage 1: Build the Go binary
FROM golang:1.25.4-alpine3.22 AS builder

# Set the working directory for the build
WORKDIR /app

# Copy all go.mod and go.sum files to download dependencies efficiently
COPY go.work go.work.sum ./
COPY apps/id_manager/go.mod apps/id_manager/go.sum ./apps/id_manager/
COPY apps/client/go.mod apps/client/go.sum ./apps/client/
COPY common/go.mod ./common/
COPY vendor ./vendor


# Change workdir back to the root to copy all source code
WORKDIR /app
COPY ./apps/client ./apps/client
COPY ./common ./common

# Build the application
# WORKDIR /app/apps/client
# RUN CGO_ENABLED=0 GOOS=linux go build -o /client ./cmd/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -mod=vendor -o /client ./apps/client/cmd/main.go

# Stage 2: Create the final, minimal image
FROM alpine:3.22.2

RUN apk add --no-cache tzdata

COPY --from=builder /client /client

CMD ["/client"]
