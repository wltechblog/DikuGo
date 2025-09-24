#!/bin/bash

# Test script to demonstrate USR1 signal stack trace functionality
echo "Testing USR1 signal stack trace functionality..."

# Start the server in background
echo "Starting DikuGo server..."
./dikugo &
SERVER_PID=$!

# Give server time to start
sleep 3

echo "Server started with PID: $SERVER_PID"

# Send USR1 signal to trigger stack trace
echo "Sending USR1 signal to generate stack trace..."
kill -USR1 $SERVER_PID

# Wait a moment for stack trace to be written
sleep 2

# Check if stack trace file was created
STACKTRACE_FILE=$(ls stacktrace_*.txt 2>/dev/null | head -1)
if [ -n "$STACKTRACE_FILE" ]; then
    echo "✓ Stack trace file created: $STACKTRACE_FILE"
    echo ""
    echo "Stack trace file contents (first 20 lines):"
    head -20 "$STACKTRACE_FILE"
    echo ""
    echo "... (file continues with full stack traces)"
    echo ""
    echo "File size: $(wc -l < "$STACKTRACE_FILE") lines"
else
    echo "✗ No stack trace file found"
fi

# Send another USR1 signal to test multiple stack traces
echo "Sending another USR1 signal..."
sleep 1
kill -USR1 $SERVER_PID
sleep 2

# Count stack trace files
STACKTRACE_COUNT=$(ls stacktrace_*.txt 2>/dev/null | wc -l)
echo "Total stack trace files created: $STACKTRACE_COUNT"

# Test that server is still running after USR1 signals
if kill -0 $SERVER_PID 2>/dev/null; then
    echo "✓ Server is still running after USR1 signals"
else
    echo "✗ Server stopped unexpectedly"
fi

# Now test graceful shutdown
echo "Testing graceful shutdown with SIGTERM..."
kill -TERM $SERVER_PID

# Wait for server to shut down
sleep 5

if kill -0 $SERVER_PID 2>/dev/null; then
    echo "⚠ Server still running, forcing shutdown..."
    kill -KILL $SERVER_PID
else
    echo "✓ Server shut down gracefully"
fi

echo ""
echo "Test completed!"
echo ""
echo "Usage instructions:"
echo "1. Start the server: ./dikugo"
echo "2. In another terminal, find the process ID: ps aux | grep dikugo"
echo "3. Send USR1 signal: kill -USR1 <PID>"
echo "4. Check for stack trace file: ls stacktrace_*.txt"
echo "5. View stack trace: cat stacktrace_<timestamp>.txt"

# Cleanup old stack trace files (optional)
echo ""
echo "Cleaning up test stack trace files..."
rm -f stacktrace_*.txt
echo "Cleanup completed."
