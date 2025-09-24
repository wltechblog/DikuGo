package network

import (
	"net"
	"testing"
	"time"

	"github.com/wltechblog/DikuGo/pkg/command"
	"github.com/wltechblog/DikuGo/pkg/config"
	"github.com/wltechblog/DikuGo/pkg/types"
	"github.com/wltechblog/DikuGo/pkg/world"
)

// MockStorage for testing
type MockStorage struct{}

func (m *MockStorage) LoadRooms() ([]*types.Room, error)                   { return nil, nil }
func (m *MockStorage) LoadObjects() ([]*types.Object, error)               { return nil, nil }
func (m *MockStorage) LoadMobiles() ([]*types.Mobile, error)               { return nil, nil }
func (m *MockStorage) LoadZones() ([]*types.Zone, error)                   { return nil, nil }
func (m *MockStorage) LoadShops() ([]*types.Shop, error)                   { return nil, nil }
func (m *MockStorage) SaveCharacter(character *types.Character) error      { return nil }
func (m *MockStorage) LoadCharacter(name string) (*types.Character, error) { return nil, nil }
func (m *MockStorage) DeleteCharacter(name string) error                   { return nil }
func (m *MockStorage) CharacterExists(name string) bool                    { return false }

func TestClientShutdown(t *testing.T) {
	// Create a mock world
	cfg := &config.Config{}
	mockStorage := &MockStorage{}
	w, err := world.NewWorld(cfg, mockStorage)
	if err != nil {
		t.Fatalf("Failed to create world: %v", err)
	}

	// Create a mock command registry
	cmdRegistry := command.NewRegistry()

	// Create a pipe to simulate a network connection
	serverConn, clientConn := net.Pipe()
	defer serverConn.Close()
	defer clientConn.Close()

	// Create a client
	client := NewClient(serverConn, w, cmdRegistry)

	// Start the client handler in a goroutine
	clientDone := make(chan struct{})
	go func() {
		client.Handle()
		close(clientDone)
	}()

	// Give the client a moment to start
	time.Sleep(100 * time.Millisecond)

	// Close the client
	client.Close()

	// Wait for the client handler to finish with a timeout
	select {
	case <-clientDone:
		// Client shut down successfully
		t.Log("Client shut down successfully")
	case <-time.After(2 * time.Second):
		t.Error("Client shutdown timed out")
	}
}

func TestServerShutdown(t *testing.T) {
	// Create a mock world
	cfg := &config.Config{}
	cfg.Server.Host = "localhost"
	cfg.Server.Port = 0 // Use port 0 to get a random available port
	mockStorage := &MockStorage{}
	w, err := world.NewWorld(cfg, mockStorage)
	if err != nil {
		t.Fatalf("Failed to create world: %v", err)
	}

	// Create a server
	server, err := NewServer(cfg, w)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	// Start the server
	go server.Start()

	// Give the server a moment to start
	time.Sleep(100 * time.Millisecond)

	// Shutdown the server
	shutdownDone := make(chan struct{})
	go func() {
		server.Shutdown()
		close(shutdownDone)
	}()

	// Wait for the server to shut down with a timeout
	select {
	case <-shutdownDone:
		// Server shut down successfully
		t.Log("Server shut down successfully")
	case <-time.After(5 * time.Second):
		t.Error("Server shutdown timed out")
	}
}
