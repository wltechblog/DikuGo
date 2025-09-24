package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
	"time"

	"github.com/wltechblog/DikuGo/pkg/config"
	"github.com/wltechblog/DikuGo/pkg/game"
)

// writeStackTrace writes a stack trace to a file with timestamp
func writeStackTrace() {
	filename := fmt.Sprintf("stacktrace_%d.txt", time.Now().Unix())
	file, err := os.Create(filename)
	if err != nil {
		log.Printf("Failed to create stack trace file %s: %v", filename, err)
		return
	}
	defer file.Close()

	// Write header with timestamp
	fmt.Fprintf(file, "Stack trace generated at: %s\n", time.Now().Format(time.RFC3339))
	fmt.Fprintf(file, "Process ID: %d\n", os.Getpid())
	fmt.Fprintf(file, "Number of goroutines: %d\n\n", runtime.NumGoroutine())

	// Get stack trace for all goroutines
	buf := make([]byte, 1024*1024) // 1MB buffer
	stackSize := runtime.Stack(buf, true)

	fmt.Fprintf(file, "=== FULL STACK TRACE ===\n")
	file.Write(buf[:stackSize])

	log.Printf("Stack trace written to %s", filename)
}

func main() {
	// No need to initialize random number generator in Go 1.20+
	// It's automatically seeded with a random value

	// Parse command line flags
	configPath := flag.String("config", "config.yaml", "Path to configuration file")
	testMode := flag.Bool("test", false, "Run in test mode to validate room connections")
	flag.Parse()

	// Load configuration
	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize the game
	gameInstance, err := game.NewGame(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize game: %v", err)
	}

	// If in test mode, validate room connections and exit
	if *testMode {
		log.Println("Running in test mode to validate room connections")
		if err := gameInstance.ValidateRooms(); err != nil {
			log.Fatalf("Validation failed: %v", err)
		}
		log.Println("Validation successful, exiting")
		return
	}

	// Set up signal handling for graceful shutdown and stack traces
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR1)

	// Use a WaitGroup to track when the game has fully started
	var wg sync.WaitGroup
	wg.Add(1)

	// Create a context that can be canceled
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start the game in a goroutine
	go func() {
		defer wg.Done()
		if err := gameInstance.Start(); err != nil {
			log.Printf("Game error: %v", err)
			cancel() // Cancel the context to signal shutdown
			os.Exit(1)
		}
	}()

	// Wait for signals or context cancellation
	for {
		select {
		case sig := <-sigChan:
			switch sig {
			case syscall.SIGUSR1:
				fmt.Println("Received USR1 signal, writing stack trace...")
				writeStackTrace()
				// Continue running after writing stack trace
				continue
			case syscall.SIGINT, syscall.SIGTERM:
				fmt.Printf("Received shutdown signal: %v\n", sig)
				// Break out of loop to start shutdown
				goto shutdown
			default:
				fmt.Printf("Received unexpected signal: %v\n", sig)
				continue
			}
		case <-ctx.Done():
			fmt.Println("Context canceled, shutting down")
			goto shutdown
		}
	}

shutdown:

	// Graceful shutdown with timeout
	fmt.Println("Shutting down...")

	// Create a channel to signal shutdown completion
	shutdownComplete := make(chan struct{})

	// Perform shutdown in a goroutine
	go func() {
		if err := gameInstance.Shutdown(); err != nil {
			log.Printf("Error during shutdown: %v", err)
		}
		close(shutdownComplete)
	}()

	// Wait for shutdown to complete with a timeout
	select {
	case <-shutdownComplete:
		fmt.Println("Shutdown completed gracefully")
	case <-time.After(10 * time.Second):
		fmt.Println("Shutdown timed out, forcing exit")
	}

	// Final cleanup
	signal.Stop(sigChan)
	fmt.Println("Exiting")
}
