#!/bin/bash

# Start the server
./dikugo &
SERVER_PID=$!

# Wait for the server to start
sleep 2

# Use netcat to test the server
echo "TestUser" | nc localhost 4000

# Kill the server
kill $SERVER_PID
