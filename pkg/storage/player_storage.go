package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/wltechblog/DikuGo/pkg/types"
)

// PlayerStorage is an interface for storing and retrieving player data
type PlayerStorage interface {
	// SavePlayer saves a player to storage
	SavePlayer(player *types.Character) error

	// LoadPlayer loads a player from storage
	LoadPlayer(name string) (*types.Character, error)

	// PlayerExists checks if a player exists in storage
	PlayerExists(name string) bool

	// DeletePlayer deletes a player from storage
	DeletePlayer(name string) error

	// ListPlayers returns a list of all player names
	ListPlayers() ([]string, error)
}

// FilePlayerStorage is a file-based implementation of PlayerStorage
type FilePlayerStorage struct {
	// Directory where player files are stored
	playerDir string
}

// NewFilePlayerStorage creates a new FilePlayerStorage
func NewFilePlayerStorage(playerDir string) (*FilePlayerStorage, error) {
	// Create the player directory if it doesn't exist
	err := os.MkdirAll(playerDir, 0755)
	if err != nil {
		return nil, fmt.Errorf("failed to create player directory: %w", err)
	}

	return &FilePlayerStorage{
		playerDir: playerDir,
	}, nil
}

// SavePlayer saves a player to a file
func (s *FilePlayerStorage) SavePlayer(player *types.Character) error {
	// Create the player file path
	filePath := s.getPlayerFilePath(player.Name)

	// Store the room VNUM before saving
	if player.InRoom != nil {
		player.RoomVNUM = player.InRoom.VNUM
	}

	// Create a copy of the player without the InRoom field
	playerCopy := *player
	playerCopy.InRoom = nil // Don't save the entire room structure
	playerCopy.World = nil  // Don't save the world reference

	// Marshal the player to JSON
	data, err := json.MarshalIndent(playerCopy, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal player: %w", err)
	}

	// Write the player data to the file
	err = os.WriteFile(filePath, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write player file: %w", err)
	}

	return nil
}

// LoadPlayer loads a player from a file
func (s *FilePlayerStorage) LoadPlayer(name string) (*types.Character, error) {
	// Create the player file path
	filePath := s.getPlayerFilePath(name)

	// Check if the player file exists
	if !s.PlayerExists(name) {
		return nil, fmt.Errorf("player %s does not exist", name)
	}

	// Read the player data from the file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read player file: %w", err)
	}

	// Unmarshal the player data
	var player types.Character
	err = json.Unmarshal(data, &player)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal player: %w", err)
	}

	// Note: The InRoom field will be set by the World.AddCharacter method
	// based on the RoomVNUM value

	return &player, nil
}

// PlayerExists checks if a player exists in storage
func (s *FilePlayerStorage) PlayerExists(name string) bool {
	// Create the player file path
	filePath := s.getPlayerFilePath(name)

	// Check if the player file exists
	_, err := os.Stat(filePath)
	return err == nil
}

// DeletePlayer deletes a player from storage
func (s *FilePlayerStorage) DeletePlayer(name string) error {
	// Create the player file path
	filePath := s.getPlayerFilePath(name)

	// Check if the player file exists
	if !s.PlayerExists(name) {
		return fmt.Errorf("player %s does not exist", name)
	}

	// Delete the player file
	err := os.Remove(filePath)
	if err != nil {
		return fmt.Errorf("failed to delete player file: %w", err)
	}

	return nil
}

// ListPlayers returns a list of all player names
func (s *FilePlayerStorage) ListPlayers() ([]string, error) {
	// Read the player directory
	files, err := os.ReadDir(s.playerDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read player directory: %w", err)
	}

	// Extract player names from file names
	var players []string
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		// Get the file name without extension
		name := strings.TrimSuffix(file.Name(), ".json")
		players = append(players, name)
	}

	return players, nil
}

// getPlayerFilePath returns the file path for a player
func (s *FilePlayerStorage) getPlayerFilePath(name string) string {
	// Convert the name to lowercase
	name = strings.ToLower(name)

	// Create the player file path
	return filepath.Join(s.playerDir, name+".json")
}
