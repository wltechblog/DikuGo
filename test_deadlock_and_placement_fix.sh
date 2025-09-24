#!/bin/bash

# Test script to verify that both the deadlock and character placement fixes work
echo "Testing deadlock and character placement fixes..."

# Start the server
echo "Starting DikuGo server..."
./dikugo &
SERVER_PID=$!

# Give server time to start
sleep 3

echo "Testing character placement and movement..."

# Test 1: Create a new character and verify it starts in room 3001
echo "=== Test 1: New character placement ==="
{
    echo "NewPlayer"    # Character name
    echo "test123"      # Password
    echo "test123"      # Confirm password
    echo "1"            # Male
    echo "1"            # Warrior class
    echo "y"            # Confirm stats
    sleep 1
    echo "look"         # Should show Temple of Midgaard (room 3001)
    sleep 1
    echo "quit"         # Quit
    sleep 1
} | nc localhost 4000 > test1_output.log 2>&1 &

# Wait for first test to complete
sleep 8

# Test 2: Test movement without deadlocks
echo "=== Test 2: Movement without deadlocks ==="
{
    echo "TestPlayer"   # Character name
    echo "test123"      # Password
    echo "test123"      # Confirm password
    echo "1"            # Male
    echo "1"            # Warrior class
    echo "y"            # Confirm stats
    sleep 1
    echo "look"         # Look around
    sleep 1
    echo "north"        # Try to move north
    sleep 1
    echo "south"        # Try to move south
    sleep 1
    echo "east"         # Try to move east
    sleep 1
    echo "west"         # Try to move west
    sleep 1
    echo "who"          # Check who's online
    sleep 1
    echo "quit"         # Quit
    sleep 1
} | nc localhost 4000 > test2_output.log 2>&1 &

# Wait for second test to complete
sleep 10

# Kill the server
echo "Stopping server..."
kill $SERVER_PID 2>/dev/null
wait $SERVER_PID 2>/dev/null

echo "=== Test Results ==="

# Check Test 1 results
echo "Test 1 - Character Placement:"
if grep -q "Temple of Midgaard" test1_output.log; then
    echo "✓ PASS: New character correctly placed in Temple of Midgaard (room 3001)"
else
    echo "✗ FAIL: New character not placed in Temple of Midgaard"
fi

# Check Test 2 results
echo "Test 2 - Movement without deadlocks:"
if grep -q "You cannot go that way" test2_output.log || grep -q "Exits:" test2_output.log; then
    echo "✓ PASS: Movement commands completed without hanging"
else
    echo "✗ FAIL: Movement commands may have hung or failed"
fi

# Check for any hanging or timeout issues
if ps aux | grep -q "[d]ikugo"; then
    echo "⚠ WARNING: Server process may still be running"
    pkill -f dikugo
else
    echo "✓ PASS: Server shut down cleanly"
fi

echo ""
echo "Detailed logs:"
echo "Test 1 output (first 10 lines):"
head -10 test1_output.log 2>/dev/null || echo "No test1 output found"
echo ""
echo "Test 2 output (first 10 lines):"
head -10 test2_output.log 2>/dev/null || echo "No test2 output found"

# Cleanup
rm -f test1_output.log test2_output.log

echo ""
echo "Test completed. If both tests passed, the deadlock and character placement fixes are working!"
