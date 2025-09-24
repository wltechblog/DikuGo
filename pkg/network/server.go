package network

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/wltechblog/DikuGo/pkg/combat"
	"github.com/wltechblog/DikuGo/pkg/command"
	"github.com/wltechblog/DikuGo/pkg/config"
	"github.com/wltechblog/DikuGo/pkg/types"
	"github.com/wltechblog/DikuGo/pkg/world"
)

// Server handles network connections
type Server struct {
	config          *config.Config
	world           *world.World
	listener        net.Listener
	clients         map[string]*Client
	mutex           sync.RWMutex
	shutdownCh      chan struct{}
	commandRegistry *command.Registry
	combatManager   command.CombatManagerInterface
}

// NewServer creates a new server instance
func NewServer(cfg *config.Config, w *world.World) (*Server, error) {
	// Create combat manager
	// Use the EnhancedDikuCombatManager for authentic DikuMUD combat with enhancements
	combatManager := combat.NewEnhancedDikuCombatManager()

	// Initialize command registry
	cmdRegistry := command.InitRegistry(w, combatManager)

	// Create server
	server := &Server{
		config:          cfg,
		world:           w,
		clients:         make(map[string]*Client),
		shutdownCh:      make(chan struct{}),
		commandRegistry: cmdRegistry,
		combatManager:   combatManager,
	}

	// Set the message handler in the world
	w.SetMessageHandler(func(ch *types.Character, message string) {
		// Find the client for this character
		client := GetClient(ch)
		if client != nil {
			client.Write(message)
		}
	})

	return server, nil
}

// Start starts the server
func (s *Server) Start() error {
	addr := fmt.Sprintf("%s:%d", s.config.Server.Host, s.config.Server.Port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", addr, err)
	}
	s.listener = listener

	log.Printf("Server listening on %s", addr)

	// Accept connections in a goroutine
	go s.acceptConnections()

	// Start the combat update loop
	go s.updateCombat()

	return nil
}

// updateCombat updates all combats
func (s *Server) updateCombat() {
	// Create a ticker that fires every 2 seconds (same as pulseViolence)
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-s.shutdownCh:
			return
		case <-ticker.C:
			s.combatManager.Update()
		}
	}
}

// acceptConnections accepts incoming connections
func (s *Server) acceptConnections() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			select {
			case <-s.shutdownCh:
				return
			default:
				log.Printf("Error accepting connection: %v", err)
				continue
			}
		}

		// Create a new client
		client := NewClient(conn, s.world, s.commandRegistry)

		// Add the client to the map
		s.mutex.Lock()
		s.clients[client.ID] = client
		s.mutex.Unlock()

		// Handle the client in a goroutine
		go func() {
			client.Handle()

			// Remove the client from the map
			s.mutex.Lock()
			delete(s.clients, client.ID)
			s.mutex.Unlock()
		}()
	}
}

// Shutdown shuts down the server
func (s *Server) Shutdown() error {
	log.Println("Shutting down network server...")

	// Signal all goroutines to stop
	if s.shutdownCh != nil {
		select {
		case <-s.shutdownCh:
			// Channel already closed
		default:
			close(s.shutdownCh)
		}
	}

	// Close the listener
	if s.listener != nil {
		log.Println("Closing network listener...")
		if err := s.listener.Close(); err != nil {
			log.Printf("Error closing listener: %v", err)
		}
	}

	// Close all client connections
	log.Printf("Closing %d client connections...", len(s.clients))
	s.mutex.Lock()
	clientCount := len(s.clients)
	for _, client := range s.clients {
		client.Close()
	}
	s.clients = make(map[string]*Client) // Clear the map
	s.mutex.Unlock()

	// Give clients time to finish their cleanup
	if clientCount > 0 {
		log.Printf("Waiting for %d client connections to finish cleanup...", clientCount)
		time.Sleep(1 * time.Second)
	}

	log.Println("Network server shutdown complete")
	return nil
}

// GetClients returns a list of all clients
func (s *Server) GetClients() []*Client {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	clients := make([]*Client, 0, len(s.clients))
	for _, client := range s.clients {
		clients = append(clients, client)
	}

	return clients
}
