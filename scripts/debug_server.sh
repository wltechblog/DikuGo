#!/bin/bash

# Production debugging script for DikuGo server
# This script provides convenient functions for debugging a running server

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

print_debug() {
    echo -e "${BLUE}[DEBUG]${NC} $1"
}

# Function to find DikuGo process
find_dikugo_process() {
    local pid=$(pgrep -f "dikugo")
    if [ -z "$pid" ]; then
        return 1
    fi
    echo "$pid"
    return 0
}

# Function to generate stack trace
generate_stack_trace() {
    local pid=$(find_dikugo_process)
    if [ $? -ne 0 ]; then
        print_error "DikuGo server is not running"
        return 1
    fi
    
    print_status "Found DikuGo server with PID: $pid"
    print_status "Sending USR1 signal to generate stack trace..."
    
    if kill -USR1 "$pid" 2>/dev/null; then
        print_status "USR1 signal sent successfully"
        sleep 2  # Give time for stack trace to be written
        
        # Find the most recent stack trace file
        local latest_trace=$(ls -t stacktrace_*.txt 2>/dev/null | head -1)
        if [ -n "$latest_trace" ]; then
            print_status "Stack trace written to: $latest_trace"
            echo ""
            echo "Stack trace summary:"
            echo "===================="
            head -5 "$latest_trace"
            echo ""
            echo "Total goroutines: $(grep -c "^goroutine" "$latest_trace")"
            echo "File size: $(wc -l < "$latest_trace") lines"
            echo ""
            echo "To view full stack trace: cat $latest_trace"
            return 0
        else
            print_error "Stack trace file not found"
            return 1
        fi
    else
        print_error "Failed to send USR1 signal (permission denied?)"
        return 1
    fi
}

# Function to show server status
show_server_status() {
    local pid=$(find_dikugo_process)
    if [ $? -ne 0 ]; then
        print_error "DikuGo server is not running"
        return 1
    fi
    
    print_status "DikuGo server status:"
    echo "====================="
    echo "PID: $pid"
    echo "Command: $(ps -p $pid -o cmd --no-headers)"
    echo "Started: $(ps -p $pid -o lstart --no-headers)"
    echo "CPU usage: $(ps -p $pid -o %cpu --no-headers)%"
    echo "Memory usage: $(ps -p $pid -o %mem --no-headers)%"
    echo "Virtual memory: $(ps -p $pid -o vsz --no-headers) KB"
    echo "Resident memory: $(ps -p $pid -o rss --no-headers) KB"
    
    # Check if server is responding on port 4000
    if nc -z localhost 4000 2>/dev/null; then
        print_status "Server is accepting connections on port 4000"
    else
        print_warning "Server may not be accepting connections on port 4000"
    fi
}

# Function to monitor server continuously
monitor_server() {
    local interval=${1:-5}
    print_status "Monitoring DikuGo server (interval: ${interval}s, press Ctrl+C to stop)"
    echo ""
    
    while true; do
        local pid=$(find_dikugo_process)
        if [ $? -ne 0 ]; then
            print_error "$(date): DikuGo server is not running"
            sleep "$interval"
            continue
        fi
        
        local cpu=$(ps -p $pid -o %cpu --no-headers | tr -d ' ')
        local mem=$(ps -p $pid -o %mem --no-headers | tr -d ' ')
        local rss=$(ps -p $pid -o rss --no-headers | tr -d ' ')
        
        echo "$(date): PID=$pid CPU=${cpu}% MEM=${mem}% RSS=${rss}KB"
        sleep "$interval"
    done
}

# Function to list stack trace files
list_stack_traces() {
    local traces=(stacktrace_*.txt)
    if [ ! -e "${traces[0]}" ]; then
        print_warning "No stack trace files found"
        return 1
    fi
    
    print_status "Available stack trace files:"
    echo "============================="
    for trace in "${traces[@]}"; do
        local timestamp=$(echo "$trace" | sed 's/stacktrace_\([0-9]*\)\.txt/\1/')
        local readable_date=$(date -d "@$timestamp" 2>/dev/null || echo "Unknown date")
        local size=$(wc -l < "$trace" 2>/dev/null || echo "?")
        echo "$trace - $readable_date ($size lines)"
    done
}

# Function to clean up old stack traces
cleanup_stack_traces() {
    local max_age_hours=${1:-24}
    local count=0
    
    print_status "Cleaning up stack trace files older than $max_age_hours hours..."
    
    for trace in stacktrace_*.txt; do
        if [ ! -e "$trace" ]; then
            continue
        fi
        
        local timestamp=$(echo "$trace" | sed 's/stacktrace_\([0-9]*\)\.txt/\1/')
        local current_time=$(date +%s)
        local age_hours=$(( (current_time - timestamp) / 3600 ))
        
        if [ "$age_hours" -gt "$max_age_hours" ]; then
            rm -f "$trace"
            count=$((count + 1))
            print_debug "Removed $trace (${age_hours}h old)"
        fi
    done
    
    print_status "Cleaned up $count old stack trace files"
}

# Function to show help
show_help() {
    echo "DikuGo Server Debug Script"
    echo "=========================="
    echo ""
    echo "Usage: $0 <command> [options]"
    echo ""
    echo "Commands:"
    echo "  status                    Show server status and resource usage"
    echo "  trace                     Generate a stack trace"
    echo "  monitor [interval]        Monitor server continuously (default: 5s)"
    echo "  list-traces              List all stack trace files"
    echo "  cleanup-traces [hours]   Clean up stack traces older than N hours (default: 24)"
    echo "  help                     Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0 status                # Show server status"
    echo "  $0 trace                 # Generate stack trace"
    echo "  $0 monitor 10            # Monitor every 10 seconds"
    echo "  $0 cleanup-traces 48     # Clean up traces older than 48 hours"
}

# Main script logic
case "${1:-help}" in
    "status")
        show_server_status
        ;;
    "trace")
        generate_stack_trace
        ;;
    "monitor")
        monitor_server "$2"
        ;;
    "list-traces")
        list_stack_traces
        ;;
    "cleanup-traces")
        cleanup_stack_traces "$2"
        ;;
    "help"|*)
        show_help
        ;;
esac
