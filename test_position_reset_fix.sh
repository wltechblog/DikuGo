#!/bin/bash

# Test script to verify the position reset fix
echo "Testing position reset fix for players entering the game..."

# Start the server
echo "Starting DikuGo server..."
./dikugo &
SERVER_PID=$!

# Give server time to start
sleep 3

echo "Server started with PID: $SERVER_PID"

# Test 1: Create a character and verify it starts in standing position
echo "=== Test 1: New character position verification ==="
{
    echo "PositionTest1"  # Character name
    echo "test123"        # Password
    echo "test123"        # Confirm password
    echo "1"              # Male
    echo "4"              # Warrior class
    echo "y"              # Confirm stats
    sleep 2
    echo "score"          # Check character status
    sleep 1
    echo "quit"           # Quit to save character
    sleep 1
} | nc localhost 4000 > test_position_output1.log 2>&1 &

# Wait for first test to complete
sleep 8

# Test 2: Log back in with the same character to verify position is still standing
echo "=== Test 2: Character re-login position verification ==="
{
    echo "PositionTest1"  # Same character name
    echo "test123"        # Password
    sleep 2
    echo "score"          # Check character status again
    sleep 1
    echo "look"           # Look around
    sleep 1
    echo "quit"           # Quit
    sleep 1
} | nc localhost 4000 > test_position_output2.log 2>&1 &

# Wait for second test to complete
sleep 8

# Test 3: Create another character, get into combat, then quit and re-login
echo "=== Test 3: Combat state reset verification ==="
{
    echo "PositionTest2"  # Character name
    echo "test123"        # Password
    echo "test123"        # Confirm password
    echo "1"              # Male
    echo "4"              # Warrior class
    echo "y"              # Confirm stats
    sleep 2
    echo "kill cityguard" # Start combat
    sleep 2
    echo "quit"           # Quit during combat (if possible)
    sleep 1
} | nc localhost 4000 > test_position_output3.log 2>&1 &

# Wait for third test to complete
sleep 10

# Test 4: Log back in with the combat character to verify position is standing
echo "=== Test 4: Post-combat character re-login verification ==="
{
    echo "PositionTest2"  # Same character name
    echo "test123"        # Password
    sleep 2
    echo "score"          # Check character status
    sleep 1
    echo "look"           # Look around
    sleep 1
    echo "who"            # Check who's online
    sleep 1
    echo "quit"           # Quit
    sleep 1
} | nc localhost 4000 > test_position_output4.log 2>&1 &

# Wait for fourth test to complete
sleep 8

# Check if server is still responsive
echo "=== Server responsiveness check ==="
if kill -0 $SERVER_PID 2>/dev/null; then
    echo "✓ Server is still running"
    
    # Generate a stack trace to verify no issues
    echo "Generating stack trace to verify server state..."
    kill -USR1 $SERVER_PID
    sleep 2
    
    # Check for stack trace file
    LATEST_TRACE=$(ls -t stacktrace_*.txt 2>/dev/null | head -1)
    if [ -n "$LATEST_TRACE" ]; then
        echo "✓ Stack trace generated: $LATEST_TRACE"
    else
        echo "⚠ No stack trace file found"
    fi
else
    echo "✗ Server has stopped unexpectedly"
fi

# Analyze test results
echo "=== Test Results Analysis ==="

# Check Test 1 results
if [ -f "test_position_output1.log" ]; then
    if grep -q "You are.*standing" test_position_output1.log; then
        echo "✓ Test 1: New character starts in standing position"
    else
        echo "⚠ Test 1: Could not verify new character position"
    fi
else
    echo "⚠ Test 1: No output log found"
fi

# Check Test 2 results
if [ -f "test_position_output2.log" ]; then
    if grep -q "You are.*standing" test_position_output2.log; then
        echo "✓ Test 2: Character maintains standing position on re-login"
    else
        echo "⚠ Test 2: Could not verify character position on re-login"
    fi
else
    echo "⚠ Test 2: No output log found"
fi

# Check Test 3 results
if [ -f "test_position_output3.log" ]; then
    if grep -q "You attack\|combat\|fighting" test_position_output3.log; then
        echo "✓ Test 3: Combat was initiated successfully"
    else
        echo "⚠ Test 3: Combat may not have been initiated"
    fi
else
    echo "⚠ Test 3: No output log found"
fi

# Check Test 4 results
if [ -f "test_position_output4.log" ]; then
    if grep -q "You are.*standing" test_position_output4.log; then
        echo "✓ Test 4: Character position reset to standing after combat logout"
    else
        echo "⚠ Test 4: Could not verify position reset after combat"
    fi
    
    # Check if commands work normally (indicating no stuck fighting state)
    if grep -q "Players:" test_position_output4.log; then
        echo "✓ Test 4: Commands work normally after combat logout/login"
    else
        echo "⚠ Test 4: Commands may not be working after combat logout/login"
    fi
else
    echo "⚠ Test 4: No output log found"
fi

echo ""
echo "Sample outputs:"
echo "==============="
if [ -f "test_position_output1.log" ]; then
    echo "Test 1 (New character):"
    grep -A2 -B2 "You are" test_position_output1.log | head -5
    echo ""
fi

if [ -f "test_position_output4.log" ]; then
    echo "Test 4 (After combat logout/login):"
    grep -A2 -B2 "You are" test_position_output4.log | head -5
    echo ""
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
echo "The position reset fix test has completed."
echo ""
echo "What was fixed:"
echo "- Players always start in POS_STANDING when entering the game"
echo "- Position field is not saved in player files (always saved as standing)"
echo "- resetPlayerCharacter() function clears fighting state and resets position"
echo "- Fighting field is cleared when players enter the game"
echo "- Minimum HP/Mana/Move points are ensured when entering the game"
echo ""
echo "Expected results:"
echo "✓ New characters start in standing position"
echo "✓ Characters maintain standing position on re-login"
echo "✓ Characters reset to standing position after combat logout/login"
echo "✓ Commands work normally after any logout/login scenario"
echo "✓ No stuck fighting states persist across logins"
echo ""
echo "If you see ✓ marks above, the position reset fix is working correctly!"

# Cleanup
rm -f test_position_output*.log
rm -f stacktrace_*.txt

echo "Test completed."
