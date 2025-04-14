package storage

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/wltechblog/DikuGo/pkg/config"
	"github.com/wltechblog/DikuGo/pkg/types"
)

// FileStorage implements the Storage interface using the original DikuMUD file formats
type FileStorage struct {
	config    *config.Config
	dataPath  string
	playerDir string
}

// NewFileStorage creates a new file storage instance
func NewFileStorage(cfg *config.Config) (*FileStorage, error) {
	// Create player directory if it doesn't exist
	if err := os.MkdirAll(cfg.Storage.PlayerDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create player directory: %w", err)
	}

	return &FileStorage{
		config:    cfg,
		dataPath:  cfg.Game.DataPath,
		playerDir: cfg.Storage.PlayerDir,
	}, nil
}

// LoadRooms loads all rooms from the world file
func (fs *FileStorage) LoadRooms() ([]*types.Room, error) {
	log.Println("Loading rooms from", filepath.Join(fs.dataPath, "tinyworld.wld"))

	// Parse the room file
	rooms, err := ParseRooms(filepath.Join(fs.dataPath, "tinyworld.wld"))
	if err != nil {
		return nil, fmt.Errorf("failed to parse rooms: %w", err)
	}

	// Create a map of rooms by VNUM for easy lookup
	roomMap := make(map[int]*types.Room)
	for _, room := range rooms {
		roomMap[room.VNUM] = room
	}

	// Validate room exits
	log.Printf("Validating room exits for %d rooms", len(rooms))
	if err := validateRoomExits(rooms, roomMap); err != nil {
		log.Printf("Warning: %v", err)
	}

	return rooms, nil
}

// validateRoomExits validates that all room exits have valid destination rooms
func validateRoomExits(rooms []*types.Room, roomMap map[int]*types.Room) error {
	validationErrors := 0

	for _, room := range rooms {
		log.Printf("Validating exits for room %d", room.VNUM)
		for dir := 0; dir < 6; dir++ {
			exit := room.Exits[dir]
			if exit == nil {
				log.Printf("Room %d has no exit in direction %d", room.VNUM, dir)
				continue
			}
			log.Printf("Room %d has exit in direction %d with DestVnum %d", room.VNUM, dir, exit.DestVnum)

			// Skip if the exit is not supposed to lead anywhere
			if exit.DestVnum == -1 {
				continue
			}

			// Check if the destination room exists
			if _, ok := roomMap[exit.DestVnum]; !ok {
				log.Printf("ERROR: Room %d has exit in direction %d with DestVnum %d but that room does not exist",
					room.VNUM, dir, exit.DestVnum)
				validationErrors++

				// Create a placeholder room
				placeholderRoom := &types.Room{
					VNUM:        exit.DestVnum,
					Name:        fmt.Sprintf("Room %d", exit.DestVnum),
					Description: fmt.Sprintf("This is a placeholder for room %d.", exit.DestVnum),
					SectorType:  0,
					Flags:       0,
					// Exits array is initialized with nil values by default
				}
				roomMap[exit.DestVnum] = placeholderRoom
				rooms = append(rooms, placeholderRoom)
				log.Printf("Created placeholder room %d for exit from room %d (direction: %d)", exit.DestVnum, room.VNUM, dir)
			}
		}
	}

	if validationErrors > 0 {
		return fmt.Errorf("found %d validation errors in room connections", validationErrors)
	}

	return nil
}

// LoadObjects loads all object prototypes from the object file
func (fs *FileStorage) LoadObjects() ([]*types.Object, error) {
	log.Println("Loading objects from", filepath.Join(fs.dataPath, "tinyworld.obj"))

	// Parse the object file
	objects, err := ParseObjects(filepath.Join(fs.dataPath, "tinyworld.obj"))
	if err != nil {
		return nil, fmt.Errorf("failed to parse objects: %w", err)
	}

	return objects, nil
}

// LoadMobiles loads all mobile prototypes from the mobile file
func (fs *FileStorage) LoadMobiles() ([]*types.Mobile, error) {
	log.Println("Loading mobiles from", filepath.Join(fs.dataPath, "tinyworld.mob"))

	// Parse the mobile file
	mobiles, err := ParseMobiles(filepath.Join(fs.dataPath, "tinyworld.mob"))
	if err != nil {
		return nil, fmt.Errorf("failed to parse mobiles: %w", err)
	}

	return mobiles, nil
}

// LoadZones loads all zones from the zone file
func (fs *FileStorage) LoadZones() ([]*types.Zone, error) {
	log.Println("Loading zones from", filepath.Join(fs.dataPath, "tinyworld.zon"))

	// Parse the zone file
	zones, err := ParseZones(filepath.Join(fs.dataPath, "tinyworld.zon"))
	if err != nil {
		return nil, fmt.Errorf("failed to parse zones: %w", err)
	}

	return zones, nil
}

// LoadShops loads all shops from the shop file
func (fs *FileStorage) LoadShops() ([]*types.Shop, error) {
	log.Println("Loading shops from", filepath.Join(fs.dataPath, "tinyworld.shp"))

	// Parse the shop file
	shops, err := ParseShops(filepath.Join(fs.dataPath, "tinyworld.shp"))
	if err != nil {
		return nil, fmt.Errorf("failed to parse shops: %w", err)
	}

	return shops, nil
}

// LoadCharacter loads a character from the player file
func (fs *FileStorage) LoadCharacter(name string) (*types.Character, error) {
	log.Println("Loading character", name)

	// Create player storage
	playerStorage, err := NewFilePlayerStorage(fs.playerDir)
	if err != nil {
		return nil, fmt.Errorf("failed to create player storage: %w", err)
	}

	// Load the player
	return playerStorage.LoadPlayer(name)
}

// SaveCharacter saves a character to the player file
func (fs *FileStorage) SaveCharacter(character *types.Character) error {
	log.Println("Saving character", character.Name)

	// Create player storage
	playerStorage, err := NewFilePlayerStorage(fs.playerDir)
	if err != nil {
		return fmt.Errorf("failed to create player storage: %w", err)
	}

	// Save the player
	return playerStorage.SavePlayer(character)
}

// DeleteCharacter deletes a character from the player file
func (fs *FileStorage) DeleteCharacter(name string) error {
	log.Println("Deleting character", name)

	// Create player storage
	playerStorage, err := NewFilePlayerStorage(fs.playerDir)
	if err != nil {
		return fmt.Errorf("failed to create player storage: %w", err)
	}

	// Delete the player
	return playerStorage.DeletePlayer(name)
}

// ListCharacters lists all characters in the player file
func (fs *FileStorage) ListCharacters() ([]string, error) {
	log.Println("Listing characters")

	// Create player storage
	playerStorage, err := NewFilePlayerStorage(fs.playerDir)
	if err != nil {
		return nil, fmt.Errorf("failed to create player storage: %w", err)
	}

	// List the players
	return playerStorage.ListPlayers()
}

// CharacterExists checks if a character exists in storage
func (fs *FileStorage) CharacterExists(name string) bool {
	log.Println("Checking if character exists", name)

	// Create player storage
	playerStorage, err := NewFilePlayerStorage(fs.playerDir)
	if err != nil {
		return false
	}

	// Check if the player exists
	return playerStorage.PlayerExists(name)
}

// LoadPlayerObjects loads a player's objects from the rent file
func (fs *FileStorage) LoadPlayerObjects(name string) ([]*types.ObjectInstance, error) {
	log.Println("Loading objects for player", name)
	// TODO: Implement player object loading from rent file
	return []*types.ObjectInstance{}, nil
}

// SavePlayerObjects saves a player's objects to the rent file
func (fs *FileStorage) SavePlayerObjects(name string, objects []*types.ObjectInstance) error {
	log.Println("Saving objects for player", name)
	// TODO: Implement player object saving to rent file
	return nil
}
