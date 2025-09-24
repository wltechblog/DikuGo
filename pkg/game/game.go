package game

import (
	"fmt"
	"log"
	"time"

	"github.com/wltechblog/DikuGo/pkg/config"
	"github.com/wltechblog/DikuGo/pkg/network"
	"github.com/wltechblog/DikuGo/pkg/storage"
	"github.com/wltechblog/DikuGo/pkg/world"
)

// Game represents the main game instance
type Game struct {
	config     *config.Config
	world      *world.World
	storage    storage.Storage
	server     *network.Server
	running    bool
	shutdownCh chan struct{}
}

// NewGame creates a new game instance
func NewGame(cfg *config.Config) (*Game, error) {
	// Initialize storage
	store, err := storage.NewStorage(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize storage: %w", err)
	}

	// Convert storage.Storage to world.Storage
	worldStore := store.(world.Storage)

	// Initialize world
	w, err := world.NewWorld(cfg, worldStore)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize world: %w", err)
	}

	// Register special procedures for mobiles
	registerSpecialProcs(w)

	// Initialize network server
	server, err := network.NewServer(cfg, w)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize server: %w", err)
	}

	return &Game{
		config:     cfg,
		world:      w,
		storage:    store,
		server:     server,
		shutdownCh: make(chan struct{}),
	}, nil
}

// Start starts the game
func (g *Game) Start() error {
	if g.running {
		return fmt.Errorf("game is already running")
	}

	g.running = true
	log.Printf("Starting DikuGo on port %d", g.config.Server.Port)

	// Start the network server
	go g.server.Start()

	// Start the game loop
	go g.gameLoop()

	return nil
}

// Shutdown gracefully shuts down the game
func (g *Game) Shutdown() error {
	if !g.running {
		return nil
	}

	log.Println("Shutting down game...")
	g.running = false

	// Signal all goroutines to stop
	log.Println("Signaling all goroutines to stop...")
	close(g.shutdownCh)

	// Give goroutines a moment to react to shutdown signal
	time.Sleep(500 * time.Millisecond)

	// Shutdown the server
	log.Println("Shutting down network server...")
	if err := g.server.Shutdown(); err != nil {
		log.Printf("Error shutting down server: %v", err)
		// Continue with shutdown even if there's an error
	}

	// Save world state
	log.Println("Saving world state...")
	if err := g.world.Save(); err != nil {
		log.Printf("Error saving world state: %v", err)
		// Continue with shutdown even if there's an error
	}

	// Final cleanup
	log.Println("Shutdown complete")
	return nil
}

// ValidateRooms validates all room connections
func (g *Game) ValidateRooms() error {
	return g.world.ValidateRooms()
}

// gameLoop is the main game loop
func (g *Game) gameLoop() {
	// Define pulse intervals (in milliseconds)
	const (
		pulseViolence    = 2000  // Combat
		pulseMobile      = 10000 // Mobile movement and actions
		pulseZone        = 60000 // Zone resets
		pulseWeather     = 30000 // Weather changes
		pulseTime        = 60000 // Game time
		pulseCorpses     = 15000 // Corpse decay
		pulsePointUpdate = 60000 // Character point updates (HP, mana, move)
		pulseAffect      = 60000 // Affect updates (spell durations)
	)

	// Create tickers for different pulse types
	violenceTicker := time.NewTicker(time.Duration(pulseViolence) * time.Millisecond)
	mobileTicker := time.NewTicker(time.Duration(pulseMobile) * time.Millisecond)
	zoneTicker := time.NewTicker(time.Duration(pulseZone) * time.Millisecond)
	weatherTicker := time.NewTicker(time.Duration(pulseWeather) * time.Millisecond)
	timeTicker := time.NewTicker(time.Duration(pulseTime) * time.Millisecond)
	corpsesTicker := time.NewTicker(time.Duration(pulseCorpses) * time.Millisecond)
	pointUpdateTicker := time.NewTicker(time.Duration(pulsePointUpdate) * time.Millisecond)
	affectTicker := time.NewTicker(time.Duration(pulseAffect) * time.Millisecond)

	defer func() {
		violenceTicker.Stop()
		mobileTicker.Stop()
		zoneTicker.Stop()
		weatherTicker.Stop()
		timeTicker.Stop()
		corpsesTicker.Stop()
		pointUpdateTicker.Stop()
		affectTicker.Stop()
	}()

	for {
		select {
		case <-g.shutdownCh:
			return
		case <-violenceTicker.C:
			g.world.PulseViolence()
		case <-mobileTicker.C:
			g.world.PulseMobile()
		case <-zoneTicker.C:
			g.world.PulseZone()
		case <-weatherTicker.C:
			g.world.PulseWeather()
		case <-timeTicker.C:
			g.world.PulseTime()
		case <-corpsesTicker.C:
			g.world.PulseCorpses()
		case <-pointUpdateTicker.C:
			g.world.PulsePointUpdate()
		case <-affectTicker.C:
			g.world.PulseAffectUpdate()
		}
	}
}
