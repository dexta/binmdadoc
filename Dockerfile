# Build stage
FROM debian:trixie-slim AS builder

# Install curl to download Go
RUN apt-get update && apt-get install -y curl ca-certificates && rm -rf /var/lib/apt/lists/*

# Download and install Go directly from source (tarball)
ENV GO_VERSION=1.21.6
RUN curl -L -o go${GO_VERSION}.linux-amd64.tar.gz https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz \
    && tar -C /usr/local -xzf go${GO_VERSION}.linux-amd64.tar.gz \
    && rm go${GO_VERSION}.linux-amd64.tar.gz

ENV PATH="/usr/local/go/bin:${PATH}"

WORKDIR /app

# Copy the go mod files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o pastemark main.go

# Final stage
FROM debian:trixie-slim

WORKDIR /app

# Copy the compiled binary from the builder stage
COPY --from=builder /app/pastemark .

# Copy static assets
COPY frontend/ ./frontend/
COPY serverTemplates/ ./serverTemplates/

# Expose the Go server port
EXPOSE 8000

# Set executable permissions (just in case)
RUN chmod +x ./pastemark

CMD ["./pastemark"]
