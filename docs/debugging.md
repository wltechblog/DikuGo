# Debugging DikuGo

This document describes debugging features available in DikuGo.

## Stack Trace Generation

DikuGo supports generating stack traces on demand using the USR1 signal. This is useful for debugging deadlocks, performance issues, or understanding what the server is doing at any given moment.

### How to Generate Stack Traces

1. **Start the server**:
   ```bash
   ./dikugo
   ```

2. **Find the process ID**:
   ```bash
   ps aux | grep dikugo
   # or
   pgrep dikugo
   ```

3. **Send USR1 signal**:
   ```bash
   kill -USR1 <PID>
   ```

4. **Check for stack trace file**:
   ```bash
   ls stacktrace_*.txt
   ```

5. **View the stack trace**:
   ```bash
   cat stacktrace_<timestamp>.txt
   ```

### Stack Trace File Format

The stack trace file contains:

- **Header Information**:
  - Timestamp when the stack trace was generated
  - Process ID
  - Number of active goroutines

- **Full Stack Trace**:
  - Complete stack traces for all goroutines
  - Function call chains
  - File names and line numbers
  - Goroutine states (running, waiting, etc.)

### Example Stack Trace File

```
Stack trace generated at: 2025-09-24T16:30:45Z
Process ID: 12345
Number of goroutines: 15

=== FULL STACK TRACE ===
goroutine 1 [running]:
main.writeStackTrace()
    /path/to/cmd/dikugo/main.go:25 +0x123
main.main.func2()
    /path/to/cmd/dikugo/main.go:105 +0x45
...

goroutine 6 [chan receive]:
github.com/wltechblog/DikuGo/pkg/network.(*Server).handleConnections(0xc000...)
    /path/to/pkg/network/server.go:89 +0x234
...
```

### Use Cases

**Debugging Deadlocks**:
- If the server appears to hang, generate a stack trace to see where goroutines are blocked
- Look for goroutines waiting on mutexes or channels

**Performance Analysis**:
- Generate multiple stack traces over time to identify hot code paths
- Look for goroutines that are consistently running or consuming CPU

**Understanding Server State**:
- See what the server is doing during normal operation
- Identify which goroutines are handling different aspects (network, AI, game logic)

### Important Notes

- **Non-Disruptive**: Sending USR1 does not stop or restart the server
- **Multiple Traces**: You can generate multiple stack traces; each gets a unique timestamp
- **File Location**: Stack trace files are created in the current working directory
- **File Naming**: Files are named `stacktrace_<unix_timestamp>.txt`
- **Automatic Cleanup**: Consider cleaning up old stack trace files periodically

### Troubleshooting

**No stack trace file created**:
- Check that the process is still running: `kill -0 <PID>`
- Verify you have write permissions in the current directory
- Check the server logs for any error messages

**Empty or truncated stack trace**:
- The default buffer is 1MB, which should be sufficient for most cases
- If you have an extremely large number of goroutines, the trace might be truncated

**Permission denied**:
- Make sure you have permission to send signals to the process
- If running as a different user, you may need sudo: `sudo kill -USR1 <PID>`

### Integration with Monitoring

You can integrate stack trace generation with monitoring systems:

```bash
# Generate stack trace and send to monitoring system
kill -USR1 $(pgrep dikugo)
sleep 2
LATEST_TRACE=$(ls -t stacktrace_*.txt | head -1)
# Send $LATEST_TRACE to your monitoring/logging system
```

### Testing

Use the provided test script to verify the functionality:

```bash
./test_usr1_signal.sh
```

This script will:
1. Start the server
2. Send USR1 signals
3. Verify stack trace files are created
4. Test that the server continues running
5. Clean up test files
