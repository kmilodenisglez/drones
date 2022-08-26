# docker build --no-cache --force-rm --tag drones_restapi .
FROM golang:1.16-alpine AS builder

RUN apk add --no-cache git

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Set the Current Working Directory inside the container
WORKDIR /tmp/go-drones-app

# Retrieve application dependencies.
# This allows the container build to reuse cached dependencies.
# Expecting to copy go.mod and if present go.sum.
COPY go.* ./
RUN go mod download
RUN go install github.com/swaggo/swag/cmd/swag@v1.7.0

# Copy local code to the container image.
COPY . .

RUN swag init --parseDependency --parseInternal --parseDepth 1 --md docs/md_endpoints
RUN go mod vendor

# testing
RUN go test -v
# Build the binary
RUN go build -v -o ./out/drones-server .

# Start fresh from a smaller image
FROM alpine:3.16.2
# RUN apk add ca-certificates

COPY --from=builder /tmp/go-drones-app/out/drones-server /app/drones-server

EXPOSE 7001
ENTRYPOINT ["/app/drones-server"]