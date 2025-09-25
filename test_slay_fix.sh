#!/bin/bash

# Test script to verify the slay command fix
echo "Testing the fixed slay command..."

# Start the server
echo "Starting DikuGo server..."
./dikugo &
SERVER_PID=$!

# Give server time to start
sleep 3

echo "Server started with PID: $SERVER_PID"

# Test: Create a character and test slay command
echo "=== Testing slay command fix ==="
{
    echo "SlayTest"        # Character name
    echo "slay123"        # Password
    echo "slay123"        # Confirm password
    echo "1"              # Male
    echo "4"              # Warrior class
    echo "y"              # Confirm stats
    sleep 2
    echo "score"          # Check initial status
    sleep 1
    echo "look"           # Look around
    sleep 1
    echo "slay cityguard" # Slay the cityguard
    sleep 3               # Give time for death processing
    echo "score"          # Check experience gain
    sleep 1
    echo "look"           # Look for corpse
    sleep 1
    echo "quit"           # Quit
    sleep 1
} | nc localhost 4000 > test_slay_fix.log 2>&1 &

# Wait for test to complete
sleep 12

# Check if server is still responsive
echo "=== Server responsiveness check ==="
if kill -0 $SERVER_PID 2>/dev/null; then
    echo "✓ Server is still running"
else
    echo "✗ Server has stopped unexpectedly"
fi

# Analyze test results
echo "=== Test Results Analysis ==="

if [ -f "test_slay_fix.log" ]; then
    echo "Slay command test results:"
    
    # Check if slay command executed
    if grep -q "You prepare to slay\|divine power" test_slay_fix.log; then
        echo "✓ Slay command executed"
    else
        echo "⚠ Slay command may not have executed"
    fi
    
    # Check if target died
    if grep -q "You have slain\|has slain" test_slay_fix.log; then
        echo "✓ Target was successfully slain"
    else
        echo "⚠ Target may not have died"
    fi
    
    # Check if experience was gained
    if grep -q "You gain.*experience" test_slay_fix.log; then
        echo "✓ Experience points were awarded"
    else
        echo "⚠ No experience gain detected"
    fi
    
    # Check if combat continued (this should NOT happen with the fix)
    if grep -q "You miss\|hits you\|slashs you" test_slay_fix.log; then
        echo "✗ Combat continued after slay (BUG - slay didn't work properly)"
    else
        echo "✓ Combat did not continue after slay (FIXED)"
    fi
    
    # Check if player died (this should NOT happen with the fix)
    if grep -q "You have been KILLED\|back at the temple" test_slay_fix.log; then
        echo "✗ Player died after slay command (BUG - target wasn't killed properly)"
    else
        echo "✓ Player did not die (FIXED)"
    fi
    
    echo ""
    echo "Sample output from test:"
    echo "========================"
    # Show the slay command and its immediate aftermath
    grep -A10 -B2 "slay cityguard" test_slay_fix.log
    echo ""
    
    echo "Experience and death messages:"
    echo "=============================="
    grep -E "experience|slain|divine power" test_slay_fix.log
    echo ""
    
else
    echo "⚠ No test output log found"
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
echo "=== Slay Command Fix Summary ==="
echo ""
echo "What was fixed:"
echo "- Slay command now applies lethal damage immediately after starting combat"
echo "- Uses normal combat system for experience calculation and death handling"
echo "- Prevents continued combat after the target is slain"
echo "- Ensures the slayer doesn't die from the target's counterattack"
echo ""
echo "Expected behavior:"
echo "✓ Target dies immediately from slay command"
echo "✓ Experience is awarded to the slayer"
echo "✓ Combat stops after the slay"
echo "✓ Slayer remains alive and can continue playing"
echo "✓ Death messages are sent to all participants"
echo "✓ Corpse is created and respawn is scheduled (for NPCs)"
echo ""

# Cleanup
rm -f test_slay_fix.log

echo "Slay command fix test completed."
