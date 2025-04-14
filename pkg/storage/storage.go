package storage

import (
	"fmt"

	"github.com/wltechblog/DikuGo/pkg/config"
	"github.com/wltechblog/DikuGo/pkg/types"
)

// Storage defines the interface for game data storage
type Storage interface {
	// World data loading
	LoadRooms() ([]*types.Room, error)
	LoadObjects() ([]*types.Object, error)
	LoadMobiles() ([]*types.Mobile, error)
	LoadZones() ([]*types.Zone, error)
	LoadShops() ([]*types.Shop, error)

	// Character data
	LoadCharacter(name string) (*types.Character, error)
	SaveCharacter(character *types.Character) error
	DeleteCharacter(name string) error
	ListCharacters() ([]string, error)
	CharacterExists(name string) bool

	// Player object data (rent system)
	LoadPlayerObjects(name string) ([]*types.ObjectInstance, error)
	SavePlayerObjects(name string, objects []*types.ObjectInstance) error
}

// NewStorage creates a new storage instance based on configuration
func NewStorage(cfg *config.Config) (Storage, error) {
	switch cfg.Storage.Type {
	case "file":
		return NewFileStorage(cfg)
	default:
		return nil, fmt.Errorf("unsupported storage type: %s", cfg.Storage.Type)
	}
}
