#!/bin/bash

# Test script to verify the fighting state fix
echo "Testing fighting state fix after killing mobs..."

# Start the server
echo "Starting DikuGo server..."
./dikugo &
SERVER_PID=$!

# Give server time to start
sleep 3

echo "Server started with PID: $SERVER_PID"

# Test 1: Kill a mob and verify commands work afterward
echo "=== Test 1: Kill mob and test commands afterward ==="
{
    echo "TestFighter"   # Character name
    echo "test123"       # Password
    echo "test123"       # Confirm password
    echo "1"             # Male
    echo "4"             # Warrior class
    echo "y"             # Confirm stats
    sleep 2
    echo "look"          # Look around
    sleep 1
    echo "kill cityguard" # Kill a mob
    sleep 3              # Wait for combat to finish
    echo "look"          # This should work if fighting state is cleared
    sleep 1
    echo "who"           # This should work if fighting state is cleared
    sleep 1
    echo "inventory"     # This should work if fighting state is cleared
    sleep 1
    echo "score"         # This should work if fighting state is cleared
    sleep 1
    echo "quit"          # Quit
    sleep 1
} | nc localhost 4000 > test_combat_output.log 2>&1 &

# Wait for test to complete
sleep 15

# Check if server is still responsive
echo "=== Test 2: Server responsiveness check ==="
if kill -0 $SERVER_PID 2>/dev/null; then
    echo "✓ Server is still running"
    
    # Generate a stack trace to verify no deadlock
    echo "Generating stack trace to verify no deadlock..."
    kill -USR1 $SERVER_PID
    sleep 2
    
    # Check for stack trace file
    LATEST_TRACE=$(ls -t stacktrace_*.txt 2>/dev/null | head -1)
    if [ -n "$LATEST_TRACE" ]; then
        echo "✓ Stack trace generated: $LATEST_TRACE"
        
        # Check if there are any goroutines stuck on fighting-related operations
        FIGHTING_ISSUES=$(grep -c "Fighting\|combat\|StopCombat" "$LATEST_TRACE" 2>/dev/null || echo "0")
        if [ "$FIGHTING_ISSUES" -eq "0" ]; then
            echo "✓ No fighting-related issues found in stack trace"
        else
            echo "⚠ Found $FIGHTING_ISSUES potential fighting-related issues in stack trace"
        fi
    else
        echo "⚠ No stack trace file found"
    fi
else
    echo "✗ Server has stopped unexpectedly"
fi

# Test 3: Check command execution results
echo "=== Test 3: Command execution results ==="
if [ -f "test_combat_output.log" ]; then
    # Check if commands after combat were executed successfully
    if grep -q "You are carrying:" test_combat_output.log; then
        echo "✓ Inventory command worked after combat"
    else
        echo "⚠ Inventory command may have failed after combat"
    fi
    
    if grep -q "You are" test_combat_output.log && grep -q "years old" test_combat_output.log; then
        echo "✓ Score command worked after combat"
    else
        echo "⚠ Score command may have failed after combat"
    fi
    
    if grep -q "Players:" test_combat_output.log; then
        echo "✓ Who command worked after combat"
    else
        echo "⚠ Who command may have failed after combat"
    fi
    
    # Check for combat messages
    if grep -q "You have slain\|has been killed" test_combat_output.log; then
        echo "✓ Combat completed successfully"
    else
        echo "⚠ Combat may not have completed"
    fi
    
    echo ""
    echo "Sample output from test:"
    echo "========================"
    tail -20 test_combat_output.log
    echo "========================"
else
    echo "⚠ No combat output log found"
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
echo "The fighting state fix test has completed."
echo ""
echo "What was fixed:"
echo "- EnhancedDikuCombatManager.StartCombat() now sets the Fighting field"
echo "- EnhancedDikuCombatManager.StopCombat() now clears Fighting field for both characters"
echo "- Death handling now calls StopCombat() to clear fighting state"
echo "- ProcessCombat() now properly handles dead targets and separated characters"
echo "- Fixed map iteration issues when modifying combat states"
echo ""
echo "Expected results:"
echo "✓ Combat commands work normally"
echo "✓ After killing a mob, all other commands should work"
echo "✓ No hanging or freezing when issuing commands after combat"
echo "✓ Server remains responsive throughout"
echo ""
echo "If you see ✓ marks above, the fighting state fix is working correctly!"

# Cleanup
rm -f test_combat_output.log
rm -f stacktrace_*.txt

echo "Test completed."
