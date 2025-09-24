#!/bin/bash

# Test script to verify that movement deadlock is fixed
# This script starts the server, connects, creates a character, and tests movement

echo "Starting DikuGo server..."
./dikugo &
SERVER_PID=$!

# Give server time to start
sleep 3

echo "Testing movement without deadlock..."

# Create a simple test using netcat
{
    echo "Squash"      # Character name
    echo "test123"     # Password
    echo "test123"     # Confirm password
    echo "1"           # Male
    echo "1"           # Warrior class
    echo "y"           # Confirm stats
    sleep 1
    echo "look"        # Look around
    sleep 1
    echo "north"       # Try to move north
    sleep 1
    echo "south"       # Try to move south
    sleep 1
    echo "east"        # Try to move east
    sleep 1
    echo "west"        # Try to move west
    sleep 1
    echo "who"         # Check who's online
    sleep 1
    echo "quit"        # Quit
    sleep 1
} | nc localhost 4000 &

# Wait for the test to complete
sleep 10

# Kill the server
echo "Stopping server..."
kill $SERVER_PID 2>/dev/null
wait $SERVER_PID 2>/dev/null

echo "Movement test completed. Check the output above for any hanging or deadlock issues."
echo "If the test completed without hanging, the deadlock fix is working!"
