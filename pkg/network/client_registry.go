package network

import (
	"sync"

	"github.com/wltechblog/DikuGo/pkg/types"
)

// ClientRegistry is a global registry that maps characters to clients
type ClientRegistry struct {
	clients map[string]*Client // Map of character names to clients
	mutex   sync.RWMutex
}

// Global instance of the client registry
var clientRegistry = &ClientRegistry{
	clients: make(map[string]*Client),
}

// RegisterClient registers a client for a character
func RegisterClient(character *types.Character, client *Client) {
	if character == nil || client == nil {
		return
	}

	clientRegistry.mutex.Lock()
	defer clientRegistry.mutex.Unlock()

	clientRegistry.clients[character.Name] = client

	// Set the Client field in the character
	character.Client = client
}

// UnregisterClient unregisters a client for a character
func UnregisterClient(character *types.Character) {
	if character == nil {
		return
	}

	clientRegistry.mutex.Lock()
	defer clientRegistry.mutex.Unlock()

	delete(clientRegistry.clients, character.Name)

	// Clear the Client field in the character
	character.Client = nil
}

// GetClient gets the client for a character
func GetClient(character *types.Character) *Client {
	if character == nil {
		return nil
	}

	clientRegistry.mutex.RLock()
	defer clientRegistry.mutex.RUnlock()

	return clientRegistry.clients[character.Name]
}

// SendMessageToCharacter sends a message to a character
func SendMessageToCharacter(character *types.Character, message string) {
	if character == nil || character.IsNPC {
		return
	}

	client := GetClient(character)
	if client != nil {
		client.Write(message)
	}
}
