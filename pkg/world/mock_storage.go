package world

import (
	"fmt"

	"github.com/wltechblog/DikuGo/pkg/types"
)

// MockStorage is a mock implementation of the Storage interface for testing
type MockStorage struct {
	rooms      []*types.Room
	objects    []*types.Object
	mobiles    []*types.Mobile
	zones      []*types.Zone
	shops      []*types.Shop
	characters map[string]*types.Character
	chars      map[string]*types.Character // Alias for characters for backward compatibility
	charObjs   map[string][]*types.ObjectInstance
}

// NewMockStorage creates a new mock storage
func NewMockStorage() *MockStorage {
	characters := make(map[string]*types.Character)
	return &MockStorage{
		rooms:      make([]*types.Room, 0),
		objects:    make([]*types.Object, 0),
		mobiles:    make([]*types.Mobile, 0),
		zones:      make([]*types.Zone, 0),
		shops:      make([]*types.Shop, 0),
		characters: characters,
		chars:      characters, // Alias for characters
		charObjs:   make(map[string][]*types.ObjectInstance),
	}
}

// LoadRooms returns the mock rooms
func (s *MockStorage) LoadRooms() ([]*types.Room, error) {
	return s.rooms, nil
}

// LoadObjects returns the mock objects
func (s *MockStorage) LoadObjects() ([]*types.Object, error) {
	return s.objects, nil
}

// LoadMobiles returns the mock mobiles
func (s *MockStorage) LoadMobiles() ([]*types.Mobile, error) {
	return s.mobiles, nil
}

// LoadZones returns the mock zones
func (s *MockStorage) LoadZones() ([]*types.Zone, error) {
	return s.zones, nil
}

// LoadShops returns the mock shops
func (s *MockStorage) LoadShops() ([]*types.Shop, error) {
	return s.shops, nil
}

// SaveCharacter saves a character to the mock storage
func (s *MockStorage) SaveCharacter(character *types.Character) error {
	s.characters[character.Name] = character
	return nil
}

// LoadCharacter loads a character from the mock storage
func (s *MockStorage) LoadCharacter(name string) (*types.Character, error) {
	if character, ok := s.characters[name]; ok {
		return character, nil
	}
	return nil, fmt.Errorf("character not found: %s", name)
}

// CharacterExists checks if a character exists in the mock storage
func (s *MockStorage) CharacterExists(name string) bool {
	_, ok := s.characters[name]
	return ok
}

// DeleteCharacter deletes a mock character
func (s *MockStorage) DeleteCharacter(name string) error {
	delete(s.characters, name)
	delete(s.chars, name)
	return nil
}

// ListCharacters lists all mock characters
func (s *MockStorage) ListCharacters() ([]string, error) {
	names := make([]string, 0, len(s.characters))
	for name := range s.characters {
		names = append(names, name)
	}
	return names, nil
}

// LoadPlayerObjects returns mock player objects
func (s *MockStorage) LoadPlayerObjects(name string) ([]*types.ObjectInstance, error) {
	if objs, ok := s.charObjs[name]; ok {
		return objs, nil
	}
	return []*types.ObjectInstance{}, nil
}

// SavePlayerObjects saves mock player objects
func (s *MockStorage) SavePlayerObjects(name string, objects []*types.ObjectInstance) error {
	s.charObjs[name] = objects
	return nil
}
