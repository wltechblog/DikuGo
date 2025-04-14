package world

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/wltechblog/DikuGo/pkg/ai"
	"github.com/wltechblog/DikuGo/pkg/config"
	"github.com/wltechblog/DikuGo/pkg/types"
)

// Storage interface for world data
type Storage interface {
	LoadRooms() ([]*types.Room, error)
	LoadObjects() ([]*types.Object, error)
	LoadMobiles() ([]*types.Mobile, error)
	LoadZones() ([]*types.Zone, error)
	LoadShops() ([]*types.Shop, error)
	SaveCharacter(character *types.Character) error
	LoadCharacter(name string) (*types.Character, error)
	CharacterExists(name string) bool
}

// TimeWeather represents the time and weather in the game world
type TimeWeather struct {
	Hour    int
	Day     int
	Month   int
	Year    int
	Weather int
}

// NewTimeWeather creates a new time and weather instance
func NewTimeWeather() *TimeWeather {
	return &TimeWeather{
		Hour:    0,
		Day:     1,
		Month:   1,
		Year:    1,
		Weather: 0,
	}
}

// World represents the game world
type World struct {
	config  *config.Config
	storage Storage
	mutex   sync.RWMutex

	// Game world data
	rooms      map[int]*types.Room         // Rooms indexed by VNUM
	objects    map[int]*types.Object       // Object prototypes indexed by VNUM
	mobiles    map[int]*types.Mobile       // Mobile prototypes indexed by VNUM
	zones      map[int]*types.Zone         // Zones indexed by VNUM
	shops      map[int]*types.Shop         // Shops indexed by VNUM
	characters map[string]*types.Character // Active characters by name

	// Game state
	time        *TimeWeather
	running     bool
	mobRespawns []*types.MobRespawn // Mobs scheduled for respawning
	aiManager   *ai.Manager         // AI manager
}

// NewWorld creates a new world instance
func NewWorld(cfg *config.Config, store Storage) (*World, error) {
	w := &World{
		config:     cfg,
		storage:    store,
		rooms:      make(map[int]*types.Room),
		objects:    make(map[int]*types.Object),
		mobiles:    make(map[int]*types.Mobile),
		zones:      make(map[int]*types.Zone),
		shops:      make(map[int]*types.Shop),
		characters: make(map[string]*types.Character),
		time:       NewTimeWeather(),
	}

	// Load world data
	if err := w.loadWorld(); err != nil {
		return nil, err
	}

	// Initialize the AI system
	w.InitAI()

	return w, nil
}

// loadWorld loads all world data from storage
func (w *World) loadWorld() error {
	log.Println("Loading world data...")

	// Load rooms
	rooms, err := w.storage.LoadRooms()
	if err != nil {
		return fmt.Errorf("failed to load rooms: %w", err)
	}
	for _, room := range rooms {
		w.rooms[room.VNUM] = room
	}
	log.Printf("Loaded %d rooms", len(w.rooms))

	// Load object prototypes
	objects, err := w.storage.LoadObjects()
	if err != nil {
		return fmt.Errorf("failed to load objects: %w", err)
	}
	for _, obj := range objects {
		w.objects[obj.VNUM] = obj
	}
	log.Printf("Loaded %d object prototypes", len(w.objects))

	// Load mobile prototypes
	mobiles, err := w.storage.LoadMobiles()
	if err != nil {
		return fmt.Errorf("failed to load mobiles: %w", err)
	}
	for _, mob := range mobiles {
		w.mobiles[mob.VNUM] = mob
	}
	log.Printf("Loaded %d mobile prototypes", len(w.mobiles))

	// Load zones
	zones, err := w.storage.LoadZones()
	if err != nil {
		return fmt.Errorf("failed to load zones: %w", err)
	}
	for _, zone := range zones {
		w.zones[zone.VNUM] = zone
	}
	log.Printf("Loaded %d zones", len(w.zones))

	// Load shops
	shops, err := w.storage.LoadShops()
	if err != nil {
		return fmt.Errorf("failed to load shops: %w", err)
	}
	for _, shop := range shops {
		w.shops[shop.VNUM] = shop

		// Link shop to room
		if room, ok := w.rooms[shop.RoomVNUM]; ok {
			room.Shop = shop
			log.Printf("Linked shop %d to room %d", shop.VNUM, shop.RoomVNUM)
		} else {
			log.Printf("Warning: Shop %d has invalid room VNUM %d", shop.VNUM, shop.RoomVNUM)
		}
	}
	log.Printf("Loaded %d shops", len(w.shops))

	return nil
}

// GetCharacter gets a character by name
func (w *World) GetCharacter(name string) (*types.Character, error) {
	log.Printf("GetCharacter: Acquiring read mutex for character %s", name)
	w.mutex.RLock()
	log.Printf("GetCharacter: Acquired read mutex for character %s", name)
	defer func() {
		log.Printf("GetCharacter: Releasing read mutex for character %s", name)
		w.mutex.RUnlock()
	}()

	// Check if character is already in memory
	if char, ok := w.characters[name]; ok {
		return char, nil
	}

	// Try to load character from storage
	return w.storage.LoadCharacter(name)
}

// GetCharacters returns a map of all characters in the game
func (w *World) GetCharacters() map[string]*types.Character {
	w.mutex.RLock()
	defer w.mutex.RUnlock()

	// Create a copy to avoid returning internal map directly
	charsCopy := make(map[string]*types.Character, len(w.characters))
	for k, v := range w.characters {
		charsCopy[k] = v
	}
	return charsCopy
}

// CharacterExists checks if a character exists
func (w *World) CharacterExists(name string) bool {
	log.Printf("CharacterExists: Acquiring read mutex for character %s", name)
	w.mutex.RLock()
	log.Printf("CharacterExists: Acquired read mutex for character %s", name)
	defer func() {
		log.Printf("CharacterExists: Releasing read mutex for character %s", name)
		w.mutex.RUnlock()
	}()

	// Check if character is already in memory
	if _, ok := w.characters[name]; ok {
		return true
	}

	// Try to load character from storage
	_, err := w.storage.LoadCharacter(name)
	return err == nil
}

// AddCharacter adds a character to the world
func (w *World) AddCharacter(character *types.Character) {
	// --- Step 1: Determine target room OUTSIDE world lock ---
	var targetRoom *types.Room
	targetRoomVNUM := 0 // Default to 0 if no specific room found

	// If the character has a RoomVNUM, try to use that
	if character.RoomVNUM > 0 {
		// Use GetRoom which acquires world RLock internally
		foundRoom := w.GetRoom(character.RoomVNUM)
		if foundRoom != nil {
			targetRoom = foundRoom
			targetRoomVNUM = targetRoom.VNUM
			log.Printf("AddCharacter: Target room determined (from save): %d", targetRoomVNUM)
		} else {
			log.Printf("AddCharacter: Saved room %d not found, will use default", character.RoomVNUM)
		}
	}

	// If we still don't have a room, use default rooms
	if targetRoom == nil {
		// Try room 0 (The Void) first
		foundRoom := w.GetRoom(0)
		if foundRoom != nil {
			targetRoom = foundRoom
			targetRoomVNUM = targetRoom.VNUM
			log.Printf("AddCharacter: Target room determined (default 0): %d", targetRoomVNUM)
		} else {
			log.Printf("AddCharacter: Starting room 0 not found, trying room 3001")
			// Try room 3001 (Temple of Midgaard)
			foundRoom = w.GetRoom(3001)
			if foundRoom != nil {
				targetRoom = foundRoom
				targetRoomVNUM = targetRoom.VNUM
				log.Printf("AddCharacter: Target room determined (default 3001): %d", targetRoomVNUM)
			} else {
				log.Printf("AddCharacter: No starting room found for character %s. Character will not be placed in a room initially.", character.Name)
				// targetRoom remains nil, targetRoomVNUM remains 0
			}
		}
	}

	// --- Step 2: Acquire world lock and modify world state ---
	log.Printf("AddCharacter: Acquiring world mutex for character %s", character.Name)
	w.mutex.Lock()
	log.Printf("AddCharacter: Acquired world mutex for character %s", character.Name)

	// Add character to the world map
	w.characters[character.Name] = character
	// Set the character's World field
	character.World = w

	// --- Step 3: If a target room was identified, lock it and add character ---
	if targetRoom != nil {
		// Re-fetch the room pointer *while holding the world lock* to ensure it wasn't deleted
		actualRoom := w.rooms[targetRoomVNUM]
		if actualRoom == targetRoom { // Pointer check for safety
			log.Printf("AddCharacter: Acquiring room mutex for room %d", actualRoom.VNUM)
			actualRoom.Lock() // Lock the specific room
			log.Printf("AddCharacter: Acquired room mutex for room %d", actualRoom.VNUM)

			// Add character to the room's list
			character.InRoom = actualRoom
			actualRoom.Characters = append(actualRoom.Characters, character)
			// Update character's RoomVNUM just in case it wasn't set or was default
			character.RoomVNUM = actualRoom.VNUM

			log.Printf("AddCharacter: Placed character %s in room %d", character.Name, actualRoom.VNUM)
			actualRoom.Unlock() // Unlock the room
			log.Printf("AddCharacter: Released room mutex for room %d", actualRoom.VNUM)
		} else {
			log.Printf("AddCharacter: Warning: Room %d changed or was removed before character %s could be placed.", targetRoomVNUM, character.Name)
			// Character remains in the world but not in a specific room initially
			character.InRoom = nil
			character.RoomVNUM = 0 // Ensure consistent state
		}
	} else {
		log.Printf("AddCharacter: Character %s added to world but not placed in a room.", character.Name)
		character.InRoom = nil
		character.RoomVNUM = 0 // Ensure consistent state
	}

	// --- Step 4: Release world lock ---
	log.Printf("AddCharacter: Releasing world mutex for character %s", character.Name)
	w.mutex.Unlock()
}

// RemoveCharacter removes a character from the world
func (w *World) RemoveCharacter(character *types.Character) {
	log.Printf("RemoveCharacter: Acquiring world mutex for character %s", character.Name)
	w.mutex.Lock()
	log.Printf("RemoveCharacter: Acquired world mutex for character %s", character.Name)

	// Remove character from the world map
	delete(w.characters, character.Name)

	// Remove character from the room (if they are in one)
	sourceRoom := character.InRoom // Get room pointer before potentially clearing it
	if sourceRoom != nil {
		log.Printf("RemoveCharacter: Acquiring room mutex for room %d", sourceRoom.VNUM)
		sourceRoom.Lock() // Lock the specific room
		log.Printf("RemoveCharacter: Acquired room mutex for room %d", sourceRoom.VNUM)

		newChars := make([]*types.Character, 0, len(sourceRoom.Characters)-1)
		for _, ch := range sourceRoom.Characters {
			if ch != character {
				newChars = append(newChars, ch)
			}
		}
		sourceRoom.Characters = newChars
		character.InRoom = nil // Clear character's room reference

		sourceRoom.Unlock() // Unlock the room
		log.Printf("RemoveCharacter: Released room mutex for room %d", sourceRoom.VNUM)
	}

	log.Printf("RemoveCharacter: Releasing world mutex for character %s", character.Name)
	w.mutex.Unlock()
}

// SaveCharacter saves a character to storage
func (w *World) SaveCharacter(character *types.Character) error {
	// Note: Saving might involve reading character state (like RoomVNUM).
	// Consider if a character-specific lock is needed, or if world RLock is sufficient.
	// For now, assume world write lock is okay, but this could be refined.
	log.Printf("SaveCharacter: Acquiring world mutex for character %s", character.Name)
	w.mutex.Lock()
	log.Printf("SaveCharacter: Acquired world mutex for character %s", character.Name)
	defer func() {
		log.Printf("SaveCharacter: Releasing world mutex for character %s", character.Name)
		w.mutex.Unlock()
	}()

	// Update RoomVNUM before saving, ensure InRoom is consistent
	if character.InRoom != nil {
		character.RoomVNUM = character.InRoom.VNUM
	} else {
		character.RoomVNUM = 0 // Or appropriate value for "not in a room"
	}

	// Save the character to storage
	return w.storage.SaveCharacter(character)
}

// GetRoom returns a room by VNUM
func (w *World) GetRoom(vnum int) *types.Room {
	w.mutex.RLock()
	defer w.mutex.RUnlock()
	return w.rooms[vnum]
}

// GetRooms returns all rooms in the world
func (w *World) GetRooms() []*types.Room {
	w.mutex.RLock()
	defer w.mutex.RUnlock()

	// Convert the map to a slice
	rooms := make([]*types.Room, 0, len(w.rooms))
	for _, room := range w.rooms {
		rooms = append(rooms, room)
	}

	return rooms
}

// GetObjectPrototype returns an object prototype by VNUM
func (w *World) GetObjectPrototype(vnum int) *types.Object {
	w.mutex.RLock()
	defer w.mutex.RUnlock()
	return w.objects[vnum]
}

// GetZone returns a zone by VNUM
func (w *World) GetZone(vnum int) *types.Zone {
	w.mutex.RLock()
	defer w.mutex.RUnlock()
	return w.zones[vnum]
}

// GetShop returns a shop by VNUM
func (w *World) GetShop(vnum int) *types.Shop {
	w.mutex.RLock()
	defer w.mutex.RUnlock()
	return w.shops[vnum]
}

// CharacterMove moves a character from one room to another
func (w *World) CharacterMove(character *types.Character, destRoom *types.Room) {
	// Acquire world lock first
	log.Printf("CharacterMove: Acquiring world mutex for %s to room %v", character.Name, destRoom) // Log destRoom pointer
	w.mutex.Lock()
	log.Printf("CharacterMove: Acquired world mutex for %s", character.Name)
	defer func() {
		log.Printf("CharacterMove: Releasing world mutex for %s", character.Name)
		w.mutex.Unlock()
	}()

	sourceRoom := character.InRoom

	// --- Early exit if no actual move needed ---
	if sourceRoom == destRoom {
		log.Printf("CharacterMove: Source and destination room are the same (%v). No move needed.", sourceRoom)
		return
	}

	// --- Define lock order based on VNUM to prevent AB-BA deadlocks between rooms ---
	var lockRoom1, lockRoom2 *types.Room
	// Handle nil rooms gracefully in lock ordering
	if sourceRoom == nil {
		lockRoom1 = destRoom // Only need to lock destination
		lockRoom2 = nil
	} else if destRoom == nil {
		lockRoom1 = sourceRoom // Only need to lock source
		lockRoom2 = nil
	} else if sourceRoom.VNUM < destRoom.VNUM {
		lockRoom1 = sourceRoom
		lockRoom2 = destRoom
	} else { // sourceRoom.VNUM >= destRoom.VNUM
		lockRoom1 = destRoom
		lockRoom2 = sourceRoom
	}

	// --- Lock rooms in consistent order ---
	if lockRoom1 != nil {
		log.Printf("CharacterMove: Acquiring room mutex 1 for room %d", lockRoom1.VNUM)
		lockRoom1.Lock()
		log.Printf("CharacterMove: Acquired room mutex 1 for room %d", lockRoom1.VNUM)
		defer func() {
			log.Printf("CharacterMove: Releasing room mutex 1 for room %d", lockRoom1.VNUM)
			lockRoom1.Unlock()
		}()
	}
	if lockRoom2 != nil { // lockRoom2 is always different from lockRoom1 if not nil
		log.Printf("CharacterMove: Acquiring room mutex 2 for room %d", lockRoom2.VNUM)
		lockRoom2.Lock()
		log.Printf("CharacterMove: Acquired room mutex 2 for room %d", lockRoom2.VNUM)
		defer func() {
			log.Printf("CharacterMove: Releasing room mutex 2 for room %d", lockRoom2.VNUM)
			lockRoom2.Unlock()
		}()
	}

	// --- Perform the move now that all necessary locks are held ---

	// Remove character from source room (if exists)
	if sourceRoom != nil {
		// Verify character is actually in sourceRoom's list before modifying
		found := false
		newChars := make([]*types.Character, 0, len(sourceRoom.Characters)-1)
		for _, ch := range sourceRoom.Characters {
			if ch == character {
				found = true
			} else {
				newChars = append(newChars, ch)
			}
		}
		if found {
			sourceRoom.Characters = newChars
			log.Printf("CharacterMove: Removed %s from room %d", character.Name, sourceRoom.VNUM)
		} else {
			// This case should ideally not happen if character.InRoom is consistent
			log.Printf("CharacterMove: Warning - Character %s not found in expected source room %d list.", character.Name, sourceRoom.VNUM)
		}
	}

	// Add character to destination room (if exists)
	character.InRoom = destRoom // Update character's reference first
	if destRoom != nil {
		destRoom.Characters = append(destRoom.Characters, character)
		character.RoomVNUM = destRoom.VNUM // Update character's VNUM
		log.Printf("CharacterMove: Added %s to room %d", character.Name, destRoom.VNUM)
	} else {
		character.RoomVNUM = 0 // Or some indicator of being out-of-world/in-limbo
		log.Printf("CharacterMove: %s moved to nil room", character.Name)
	}
}

// PulseViolence is implemented in combat.go

// PulseMobile handles mobile updates
func (w *World) PulseMobile() {
	// Create a snapshot of mobiles to iterate over, avoiding lock contention during iteration.
	// A read lock is needed briefly to safely copy the character map.
	w.mutex.RLock()
	mobilesSnapshot := make([]*types.Character, 0, len(w.characters))
	for _, char := range w.characters {
		if char.IsNPC {
			mobilesSnapshot = append(mobilesSnapshot, char)
		}
	}
	w.mutex.RUnlock() // Release the lock immediately after copying

	// Process the collected mobiles' AI behavior directly
	if w.aiManager != nil {
		// Update the AI with the snapshot of mobiles collected under the read lock
		w.aiManager.TickWithMobiles(mobilesSnapshot) // Use the snapshot
	}
}

// PulseZone handles zone resets
func (w *World) PulseZone() {
	// Reset zones that are due for a reset
	w.ResetZones() // This acquires world lock

	// Process mob respawns
	w.ProcessMobRespawns() // This acquires world lock
}

// CreateMobFromPrototype creates a new mobile instance from a prototype
// Assumes world lock is NOT held by caller.
func (w *World) CreateMobFromPrototype(vnum int, room *types.Room) *types.Character {
	log.Printf("CreateMobFromPrototype: Acquiring world lock for mob %d", vnum)
	w.mutex.Lock()
	log.Printf("CreateMobFromPrototype: Acquired world lock for mob %d", vnum)

	// Get the mobile prototype (under world lock)
	mobProto := w.mobiles[vnum]
	if mobProto == nil {
		w.mutex.Unlock() // Unlock before returning
		log.Printf("CreateMobFromPrototype: Releasing world lock (prototype %d not found)", vnum)
		log.Printf("Warning: Mobile prototype %d not found", vnum)
		return nil
	}

	// Create a new mobile instance
	mob := &types.Character{
		Name:        mobProto.Name,
		ShortDesc:   mobProto.ShortDesc,
		LongDesc:    mobProto.LongDesc,
		Description: mobProto.Description,
		Level:       mobProto.Level,
		Gold:        mobProto.Gold,
		Position:    types.POS_STANDING, // Or mobProto.DefaultPos?
		IsNPC:       true,
		Prototype:   mobProto,
		World:       w,
		Inventory:   make([]*types.ObjectInstance, 0),
		Equipment:   make([]*types.ObjectInstance, types.NUM_WEARS),
		// RoomVNUM and InRoom will be set when placed
	}

	// Add the mobile to the world map (under world lock)
	w.characters[mob.Name] = mob // Consider potential name collisions?

	// Add the mobile to the room (if provided)
	if room != nil {
		// Verify room pointer is still valid under world lock
		actualRoom := w.rooms[room.VNUM]
		if actualRoom == room {
			log.Printf("CreateMobFromPrototype: Acquiring room lock for room %d", actualRoom.VNUM)
			actualRoom.Lock() // Lock the room before modifying
			log.Printf("CreateMobFromPrototype: Acquired room lock for room %d", actualRoom.VNUM)

			mob.InRoom = actualRoom
			actualRoom.Characters = append(actualRoom.Characters, mob)
			mob.RoomVNUM = actualRoom.VNUM

			actualRoom.Unlock() // Unlock the room
			log.Printf("CreateMobFromPrototype: Released room lock for room %d", actualRoom.VNUM)
		} else {
			log.Printf("CreateMobFromPrototype: Warning - Room %d changed or removed before mob %d could be placed.", room.VNUM, vnum)
			// Mob exists in world map but not placed in room
		}
	}

	log.Printf("CreateMobFromPrototype: Releasing world lock for mob %d", vnum)
	w.mutex.Unlock() // Release world lock

	return mob
}

// ScheduleMobRespawn schedules a mob for respawning
// Assumes world lock is NOT held by caller.
func (w *World) ScheduleMobRespawn(mob *types.Character) {
	if mob == nil || !mob.IsNPC || mob.Prototype == nil {
		return
	}

	roomVNUMToRespawn := mob.RoomVNUM // Capture VNUM before locking

	log.Printf("ScheduleMobRespawn: Acquiring world lock for mob %s (VNUM %d)", mob.Name, mob.Prototype.VNUM)
	w.mutex.Lock()
	log.Printf("ScheduleMobRespawn: Acquired world lock for mob %s", mob.Name)

	// --- Determine Zone (under world lock) ---
	var zone *types.Zone
	sourceRoom := mob.InRoom // Get current room pointer under lock
	if sourceRoom != nil && sourceRoom.Zone != nil {
		zone = sourceRoom.Zone
	} else {
		// Fallback: find zone by VNUM range
		for _, z := range w.zones {
			// Use captured VNUM as mob.InRoom might be nil now
			if roomVNUMToRespawn >= z.MinVNUM && roomVNUMToRespawn <= z.MaxVNUM {
				zone = z
				break
			}
		}
	}

	if zone == nil {
		w.mutex.Unlock()
		log.Printf("ScheduleMobRespawn: Releasing world lock (failed to find zone for mob %s)", mob.Name)
		log.Printf("Failed to find zone for mob %s (VNUM %d) in room %d",
			mob.Name, mob.Prototype.VNUM, roomVNUMToRespawn)
		return
	}

	// --- Create Respawn Entry (under world lock) ---
	respawnInterval := time.Duration(zone.Lifespan) * time.Minute
	respawnTime := time.Now().Add(respawnInterval)
	respawn := &types.MobRespawn{
		MobVNUM:     mob.Prototype.VNUM,
		RoomVNUM:    roomVNUMToRespawn, // Use captured VNUM
		RespawnTime: respawnTime,
	}
	w.mobRespawns = append(w.mobRespawns, respawn)

	// --- Remove Mob from World Map (under world lock) ---
	delete(w.characters, mob.Name)

	// --- Remove Mob from Room (needs room lock) ---
	if sourceRoom != nil {
		// Verify room pointer consistency
		actualRoom := w.rooms[sourceRoom.VNUM]
		if actualRoom == sourceRoom {
			log.Printf("ScheduleMobRespawn: Acquiring room lock for room %d", actualRoom.VNUM)
			actualRoom.Lock()
			log.Printf("ScheduleMobRespawn: Acquired room lock for room %d", actualRoom.VNUM)

			newChars := make([]*types.Character, 0, len(actualRoom.Characters)-1)
			found := false
			for _, ch := range actualRoom.Characters {
				if ch == mob {
					found = true
				} else {
					newChars = append(newChars, ch)
				}
			}
			if found {
				actualRoom.Characters = newChars
			} else {
				log.Printf("ScheduleMobRespawn: Warning - Mob %s not found in expected room %d list.", mob.Name, actualRoom.VNUM)
			}
			mob.InRoom = nil // Clear mob's room reference

			actualRoom.Unlock()
			log.Printf("ScheduleMobRespawn: Released room lock for room %d", actualRoom.VNUM)
		} else {
			log.Printf("ScheduleMobRespawn: Warning - Room %d changed or removed before mob %s could be removed from it.", sourceRoom.VNUM, mob.Name)
			mob.InRoom = nil // Still clear mob's reference
		}
	} else {
		mob.InRoom = nil // Ensure reference is cleared if mob wasn't in a room
	}

	log.Printf("ScheduleMobRespawn: Releasing world lock for mob %s", mob.Name)
	w.mutex.Unlock() // Release world lock

	log.Printf("Scheduled mob %s (VNUM %d) for respawn in room %d at %s",
		mob.Name, mob.Prototype.VNUM, roomVNUMToRespawn, respawnTime.Format(time.RFC3339))
}

// GetZoneForRoom returns the zone that contains the given room VNUM
func (w *World) GetZoneForRoom(roomVNUM int) *types.Zone {
	// Get the room (acquires world RLock)
	room := w.GetRoom(roomVNUM)
	if room == nil {
		return nil
	}

	// If the room already has a zone reference, return it (no extra lock needed)
	// Note: Assumes Zone pointer on Room is set atomically during load/creation.
	if room.Zone != nil {
		return room.Zone
	}

	// Otherwise, find the zone that contains this room (needs world RLock)
	w.mutex.RLock()
	defer w.mutex.RUnlock()
	for _, zone := range w.zones {
		if roomVNUM >= zone.MinVNUM && roomVNUM <= zone.MaxVNUM {
			return zone
		}
	}

	return nil
}

// ProcessMobRespawns processes mob respawns
// Assumes world lock is held by caller (e.g., PulseZone -> ResetZones).
func (w *World) ProcessMobRespawns() {
	// Check if there are any mobs to respawn
	if len(w.mobRespawns) == 0 {
		return
	}

	now := time.Now()
	var remainingRespawns []*types.MobRespawn
	mobsToRespawn := w.mobRespawns // Process current list

	w.mobRespawns = nil // Clear original slice, will rebuild with remaining

	for _, respawn := range mobsToRespawn {
		if now.After(respawn.RespawnTime) {
			// Get room and prototype (safe under world lock held by caller)
			room := w.rooms[respawn.RoomVNUM]
			if room == nil {
				log.Printf("ProcessMobRespawns: Failed to find room %d for mob respawn (VNUM %d)",
					respawn.RoomVNUM, respawn.MobVNUM)
				continue // Skip this respawn
			}
			mobProto := w.mobiles[respawn.MobVNUM]
			if mobProto == nil {
				log.Printf("ProcessMobRespawns: Failed to find mobile prototype %d for respawn", respawn.MobVNUM)
				continue // Skip this respawn
			}

			// Create a new mobile instance
			mob := &types.Character{
				Name:        mobProto.Name,
				ShortDesc:   mobProto.ShortDesc,
				LongDesc:    mobProto.LongDesc,
				Description: mobProto.Description,
				Level:       mobProto.Level,
				Gold:        mobProto.Gold,
				Position:    types.POS_STANDING, // Or mobProto.DefaultPos?
				IsNPC:       true,
				Prototype:   mobProto,
				World:       w,
				Inventory:   make([]*types.ObjectInstance, 0),
				Equipment:   make([]*types.ObjectInstance, types.NUM_WEARS),
			}

			// Add the mobile to the world map (under world lock)
			w.characters[mob.Name] = mob // Again, potential name collision?

			// Add the mobile to the room (needs room lock)
			log.Printf("ProcessMobRespawns: Acquiring room lock for room %d", room.VNUM)
			room.Lock()
			log.Printf("ProcessMobRespawns: Acquired room lock for room %d", room.VNUM)

			mob.InRoom = room
			room.Characters = append(room.Characters, mob)
			mob.RoomVNUM = room.VNUM

			room.Unlock()
			log.Printf("ProcessMobRespawns: Released room lock for room %d", room.VNUM)

			log.Printf("Respawned mob %s (VNUM %d) in room %d",
				mob.Name, respawn.MobVNUM, respawn.RoomVNUM)
		} else {
			// Keep this respawn entry for later
			remainingRespawns = append(remainingRespawns, respawn)
		}
	}

	// Update the respawn list
	w.mobRespawns = remainingRespawns
}
