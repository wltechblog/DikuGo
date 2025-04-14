package network

import (
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/wltechblog/DikuGo/pkg/combat"
	"github.com/wltechblog/DikuGo/pkg/command"
	"github.com/wltechblog/DikuGo/pkg/config"
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
	combatManager   *combat.Manager
}

// NewServer creates a new server instance
func NewServer(cfg *config.Config, w *world.World) (*Server, error) {
	// Create combat manager
	combatManager := combat.NewManager()

	// Initialize command registry
	cmdRegistry := command.InitRegistry(w, combatManager)

	return &Server{
		config:          cfg,
		world:           w,
		clients:         make(map[string]*Client),
		shutdownCh:      make(chan struct{}),
		commandRegistry: cmdRegistry,
		combatManager:   combatManager,
	}, nil
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
	for {
		select {
		case <-s.shutdownCh:
			return
		default:
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
	// Close the listener
	if s.listener != nil {
		s.listener.Close()
	}

	// Signal all goroutines to stop
	close(s.shutdownCh)

	// Close all client connections
	s.mutex.Lock()
	for _, client := range s.clients {
		client.Close()
	}
	s.mutex.Unlock()

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
