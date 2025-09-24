package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/wltechblog/DikuGo/pkg/config"
	"github.com/wltechblog/DikuGo/pkg/game"
)

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

	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

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

	// Wait for a signal or context cancellation
	select {
	case sig := <-sigChan:
		fmt.Printf("Received signal: %v\n", sig)
	case <-ctx.Done():
		fmt.Println("Context canceled, shutting down")
	}

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
