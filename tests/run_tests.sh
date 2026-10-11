#!/bin/bash
set -e

# Ensure we run from the project root
cd "$(dirname "$0")/.."

# Ensure cleanup on exit
cleanup() {
    echo "Cleaning up..."
    if [ -n "$NODE_PID" ]; then
        kill $NODE_PID 2>/dev/null || true
    fi
    if [ -n "$GO_PID" ]; then
        kill -INT $GO_PID 2>/dev/null || true
    fi
    if [ -n "$PASTEMARK_EXEC_PATH" ]; then
        rm -rf "$PASTEMARK_EXEC_PATH"
    fi
}
trap cleanup EXIT

echo "Starting Node server..."
cd jsServer
pnpm install
node app.js &
NODE_PID=$!
cd ..

echo "Starting Go server..."
# Build the server
go build -o server_bin main.go

# Create and use a temporary directory for document stores
export PASTEMARK_EXEC_PATH=$(mktemp -d)
echo "Using temp store path: $PASTEMARK_EXEC_PATH"

# Copy templates and frontend to the temp dir so the server can find them when ExecPath is set to the temp dir
cp -r serverTemplates $PASTEMARK_EXEC_PATH/
cp -r frontend $PASTEMARK_EXEC_PATH/

# Run Go server in the background
./server_bin &
GO_PID=$!

echo "Waiting for Node server..."
npx -y wait-on tcp:8080

echo "Waiting for Go server..."
npx -y wait-on tcp:8000

echo "Running Go tests..."
# Run the tests from the root directory so the go mod is respected
go test -v ./tests/...

echo "Tests finished."
