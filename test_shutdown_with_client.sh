#!/bin/bash

echo "Testing server shutdown with connected client..."

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

# Connect a client in the background
echo "Connecting test client..."
(
    sleep 1
    echo "TestUser"
    sleep 1
    echo "y"
    sleep 1
    echo "testpass"
    sleep 1
    echo "testpass"
    sleep 1
    echo "1"
    sleep 10  # Stay connected for a while
) | telnet localhost 4000 &
CLIENT_PID=$!

# Give the client time to connect
sleep 2

# Send SIGINT to trigger graceful shutdown
echo "Sending SIGINT to server with connected client..."
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
        kill -9 $CLIENT_PID 2>/dev/null
        exit 1
    fi
done

# Clean up client process
kill -9 $CLIENT_PID 2>/dev/null

echo "SUCCESS: Server with connected client shut down gracefully in ${COUNT} seconds"
echo "Shutdown test with client completed successfully!"
