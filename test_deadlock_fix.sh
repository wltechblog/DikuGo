#!/bin/bash

# Test script to verify the weather/time deadlock fix
echo "Testing weather/time deadlock fix..."

# Start the server
echo "Starting DikuGo server..."
./dikugo &
SERVER_PID=$!

# Give server time to start
sleep 3

echo "Server started with PID: $SERVER_PID"

# Test 1: Connect a client and try to move while weather system is running
echo "=== Test 1: Movement during weather updates ==="
{
    echo "TestPlayer"   # Character name
    echo "test123"      # Password
    echo "test123"      # Confirm password
    echo "1"            # Male
    echo "1"            # Warrior class
    echo "y"            # Confirm stats
    sleep 2
    echo "look"         # Look around
    sleep 1
    echo "north"        # Try to move
    sleep 1
    echo "south"        # Try to move
    sleep 1
    echo "east"         # Try to move
    sleep 1
    echo "west"         # Try to move
    sleep 1
    echo "who"          # Check who's online
    sleep 1
    echo "quit"         # Quit
    sleep 1
} | nc localhost 4000 > test_movement_output.log 2>&1 &

# Wait for test to complete
sleep 10

# Check if server is still responsive
echo "=== Test 2: Server responsiveness check ==="
if kill -0 $SERVER_PID 2>/dev/null; then
    echo "✓ Server is still running"
    
    # Try to generate a stack trace to see current state
    echo "Generating stack trace to verify no deadlock..."
    kill -USR1 $SERVER_PID
    sleep 2
    
    # Check for stack trace file
    LATEST_TRACE=$(ls -t stacktrace_*.txt 2>/dev/null | head -1)
    if [ -n "$LATEST_TRACE" ]; then
        echo "✓ Stack trace generated: $LATEST_TRACE"
        
        # Check if there are any goroutines stuck on RWMutex.RLock
        DEADLOCK_COUNT=$(grep -c "sync.RWMutex.RLock" "$LATEST_TRACE" 2>/dev/null || echo "0")
        if [ "$DEADLOCK_COUNT" -eq "0" ]; then
            echo "✓ No deadlocked goroutines found in stack trace"
        else
            echo "⚠ Found $DEADLOCK_COUNT potentially deadlocked goroutines"
            echo "First few lines of stack trace:"
            head -20 "$LATEST_TRACE"
        fi
    else
        echo "⚠ No stack trace file found"
    fi
else
    echo "✗ Server has stopped unexpectedly"
fi

# Test 3: Check movement output
echo "=== Test 3: Movement command results ==="
if [ -f "test_movement_output.log" ]; then
    if grep -q "You cannot go that way\|Exits:" test_movement_output.log; then
        echo "✓ Movement commands completed successfully"
    else
        echo "⚠ Movement commands may have issues"
        echo "Movement output (first 10 lines):"
        head -10 test_movement_output.log
    fi
else
    echo "⚠ No movement output log found"
fi

# Graceful shutdown
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
echo "=== Test Summary ==="
echo "The deadlock fix test has completed."
echo ""
echo "What was fixed:"
echo "- Removed redundant locking in SendToOutdoor() method"
echo "- PulseWeather() and PulseTime() already hold world mutex"
echo "- SendToOutdoor() no longer tries to acquire the same lock"
echo "- This eliminates the AB-BA deadlock between weather and movement systems"
echo ""
echo "If the tests above show:"
echo "✓ Server remained responsive"
echo "✓ No deadlocked goroutines in stack trace"  
echo "✓ Movement commands completed"
echo "✓ Graceful shutdown"
echo ""
echo "Then the deadlock fix is working correctly!"

# Cleanup
rm -f test_movement_output.log
rm -f stacktrace_*.txt

echo "Test completed."
