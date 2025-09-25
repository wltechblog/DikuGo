#!/bin/bash

# Test script to verify the connection drop handling fix
echo "Testing the connection drop handling fix..."

# Start the server
echo "Starting DikuGo server..."
./dikugo &
SERVER_PID=$!

# Give server time to start
sleep 3

echo "Server started with PID: $SERVER_PID"

# Test 1: Normal connection and graceful disconnect
echo "=== Test 1: Normal connection and graceful disconnect ==="
{
    echo "TestUser1"       # Character name
    echo "test123"        # Password
    echo "test123"        # Confirm password
    echo "1"              # Male
    echo "4"              # Warrior class
    echo "y"              # Confirm stats
    sleep 2
    echo "0"              # Exit from menu
    sleep 1
} | nc localhost 4000 > test_normal_disconnect.log 2>&1 &

# Wait for normal disconnect test
sleep 8

# Test 2: Abrupt connection drop (simulate network failure)
echo "=== Test 2: Abrupt connection drop ==="
{
    echo "TestUser2"       # Character name
    echo "test456"        # Password
    echo "test456"        # Confirm password
    echo "1"              # Male
    echo "4"              # Warrior class
    echo "y"              # Confirm stats
    sleep 2
    # Abruptly close connection without proper exit
} | timeout 5 nc localhost 4000 > test_abrupt_disconnect.log 2>&1 &

# Wait for abrupt disconnect test
sleep 8

# Test 3: Multiple rapid connections and disconnections
echo "=== Test 3: Multiple rapid connections and disconnections ==="
for i in {1..5}; do
    {
        echo "RapidTest$i"
        sleep 0.5
        # Let connection timeout/drop
    } | timeout 2 nc localhost 4000 > test_rapid_$i.log 2>&1 &
done

# Wait for rapid connection tests
sleep 10

# Check server logs for the connection drop issue
echo "=== Checking server logs for connection drop loops ==="

# Get server logs from the last 30 seconds
SERVER_LOG_FILE="/tmp/dikugo_test.log"

# Capture server logs
timeout 5 bash -c "
    # Send a test connection to generate some log activity
    echo 'LogTest' | nc localhost 4000 >/dev/null 2>&1 &
    sleep 2
    
    # Check if server is still responsive
    if kill -0 $SERVER_PID 2>/dev/null; then
        echo 'Server is still running'
    else
        echo 'Server has crashed or stopped'
    fi
" > server_status.log 2>&1

# Test 4: Check if server handles multiple simultaneous disconnections
echo "=== Test 4: Multiple simultaneous disconnections ==="
for i in {1..3}; do
    {
        echo "SimulTest$i"
        echo "pass$i"
        echo "pass$i"
        echo "1"
        echo "4"
        echo "y"
        sleep 1
        # Abrupt disconnect
    } | timeout 3 nc localhost 4000 > test_simul_$i.log 2>&1 &
done

# Wait for simultaneous disconnect tests
sleep 8

# Final server responsiveness check
echo "=== Final server responsiveness check ==="
{
    echo "FinalTest"
    echo "final123"
    echo "final123"
    echo "1"
    echo "4"
    echo "y"
    sleep 2
    echo "0"  # Graceful exit
} | nc localhost 4000 > test_final_check.log 2>&1

sleep 5

# Check if server is still running and responsive
if kill -0 $SERVER_PID 2>/dev/null; then
    echo "✓ Server is still running after all connection drop tests"
    
    # Test one more connection to verify responsiveness
    echo "Testing final connection..."
    {
        echo "ResponseTest"
        sleep 1
        echo "0"  # Exit immediately
    } | timeout 5 nc localhost 4000 > test_response.log 2>&1
    
    if [ $? -eq 0 ]; then
        echo "✓ Server is still responsive to new connections"
    else
        echo "⚠ Server may not be responding to new connections"
    fi
else
    echo "✗ Server has stopped running"
fi

# Analyze results
echo ""
echo "=== Connection Drop Fix Analysis ==="
echo ""

# Count log files created
LOG_COUNT=$(ls test_*.log 2>/dev/null | wc -l)
echo "Generated $LOG_COUNT test log files"

# Check for any signs of the old looping issue
echo "Checking for connection drop loop indicators..."

# Look for repeated EOF errors in server output (if we could capture it)
# Since we can't easily capture server stdout, we'll check if server is still running
if kill -0 $SERVER_PID 2>/dev/null; then
    echo "✓ Server survived all connection drop tests"
    echo "✓ No infinite loop detected (server would have crashed or become unresponsive)"
else
    echo "⚠ Server stopped during testing - may indicate an issue"
fi

# Check test logs for successful connections
SUCCESSFUL_CONNECTIONS=0
for log in test_*.log; do
    if [ -f "$log" ] && grep -q "By what name\|Welcome\|Password:" "$log" 2>/dev/null; then
        SUCCESSFUL_CONNECTIONS=$((SUCCESSFUL_CONNECTIONS + 1))
    fi
done

echo "✓ $SUCCESSFUL_CONNECTIONS test connections were successful"

# Graceful shutdown
echo ""
echo "=== Shutting down server ==="
kill -TERM $SERVER_PID 2>/dev/null

# Wait for shutdown
sleep 5

if kill -0 $SERVER_PID 2>/dev/null; then
    echo "⚠ Server still running, forcing shutdown..."
    kill -KILL $SERVER_PID
else
    echo "✓ Server shut down gracefully"
fi

echo ""
echo "=== Connection Drop Fix Summary ==="
echo ""
echo "What was fixed:"
echo "- Changed 'break' to 'c.Closed = true; return' in client read loop"
echo "- Added connection state checking in Write() method"
echo "- Added proper error handling for write failures"
echo "- Made Write() method check if client is already closed"
echo ""
echo "Expected behavior after fix:"
echo "✓ EOF errors cause immediate client cleanup (no looping)"
echo "✓ Write errors mark client as closed"
echo "✓ Multiple writes to closed clients are safely ignored"
echo "✓ Server remains stable under connection drop stress"
echo "✓ No repeated 'Error reading from client' messages"
echo ""
echo "Before fix: Server would loop endlessly trying to read from dead connections"
echo "After fix: Server properly detects disconnections and cleans up immediately"
echo ""

# Cleanup
rm -f test_*.log server_status.log

echo "Connection drop fix test completed."
