#!/bin/bash

echo "Testing server shutdown behavior..."

# Start the server in the background
echo "Starting DikuGo server..."
./dikugo &
SERVER_PID=$!

# Wait for the server to start
echo "Waiting for server to start..."
sleep 3

# Check if server is running
if ! kill -0 $SERVER_PID 2>/dev/null; then
    echo "ERROR: Server failed to start"
    exit 1
fi

echo "Server started with PID $SERVER_PID"

# Send SIGINT to trigger graceful shutdown
echo "Sending SIGINT to server..."
kill -INT $SERVER_PID

# Wait for graceful shutdown with timeout
echo "Waiting for graceful shutdown..."
TIMEOUT=15
COUNT=0

while kill -0 $SERVER_PID 2>/dev/null; do
    sleep 1
    COUNT=$((COUNT + 1))
    if [ $COUNT -ge $TIMEOUT ]; then
        echo "ERROR: Server shutdown timed out after ${TIMEOUT} seconds"
        echo "Forcing shutdown with SIGKILL..."
        kill -9 $SERVER_PID
        exit 1
    fi
done

echo "SUCCESS: Server shut down gracefully in ${COUNT} seconds"
echo "Shutdown test completed successfully!"
