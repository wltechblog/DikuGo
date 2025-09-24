package world

import (
	"fmt"
	"log"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/wltechblog/DikuGo/pkg/ai"
	"github.com/wltechblog/DikuGo/pkg/config"
	"github.com/wltechblog/DikuGo/pkg/types"
	"github.com/wltechblog/DikuGo/pkg/utils"
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
	DeleteCharacter(name string) error
	CharacterExists(name string) bool
}

// TimeWeather represents the time and weather in the game world
type TimeWeather struct {
	Hours    int
	Day      int
	Month    int
	Year     int
	Sunlight int
	Weather  int
	Change   int
}

// NewTimeWeather creates a new time and weather instance
func NewTimeWeather() *TimeWeather {
	return &TimeWeather{
		Hours:    0,
		Day:      0,
		Month:    0,
		Year:     1,
		Sunlight: types.SUN_DARK,
		Weather:  types.SKY_CLOUDLESS,
		Change:   0,
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
	rand        *rand.Rand          // Random number generator

	// Message handler
	messageHandler func(*types.Character, string) // Function to handle messages to characters
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
		rand:       rand.New(rand.NewSource(time.Now().UnixNano())),
		messageHandler: func(ch *types.Character, message string) {
			// Default message handler - just log the message
			// log.Printf("Message to %s: %s", ch.Name, message)
		},
	}

	// Load world data
	if err := w.loadWorld(); err != nil {
		return nil, err
	}

	// Perform initial zone reset to spawn mobs
	log.Println("Performing initial zone reset...")
	// Force all zones to reset by setting their age to their lifespan
	for _, zone := range w.zones {
		zone.Age = zone.Lifespan
	}
	w.ResetZones()

	// Initialize the time system
	w.InitializeTime()

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

	// Process zone commands to set up default equipment for mob prototypes
	w.ProcessZoneCommands()

	// Load shops
	shops, err := w.storage.LoadShops()
	if err != nil {
		return fmt.Errorf("failed to load shops: %w", err)
	}

	// Debug: Print all loaded shops
	log.Printf("Loaded %d shops from storage", len(shops))

	// First, clear any existing shops
	w.shops = make(map[int]*types.Shop)

	// Then, for each room, clear any shop references
	for _, room := range w.rooms {
		room.Shop = nil
	}

	// Now process each shop
	for _, shop := range shops {
		log.Printf("Processing shop #%d: Room VNUM = %d, Keeper VNUM = %d, Items: %v",
			shop.VNUM, shop.RoomVNUM, shop.MobileVNUM, shop.Producing)

		// Create a deep copy of the shop
		shopCopy := &types.Shop{
			VNUM:       shop.VNUM,
			RoomVNUM:   shop.RoomVNUM,
			MobileVNUM: shop.MobileVNUM,
			ProfitBuy:  shop.ProfitBuy,
			ProfitSell: shop.ProfitSell,
			OpenHour:   shop.OpenHour,
			CloseHour:  shop.CloseHour,
		}

		// Deep copy the slices
		shopCopy.Producing = make([]int, len(shop.Producing))
		copy(shopCopy.Producing, shop.Producing)

		shopCopy.BuyTypes = make([]int, len(shop.BuyTypes))
		copy(shopCopy.BuyTypes, shop.BuyTypes)

		shopCopy.Messages = make([]string, len(shop.Messages))
		copy(shopCopy.Messages, shop.Messages)

		// Set the World field
		shopCopy.World = w

		// Store the deep copy in the world
		w.shops[shopCopy.VNUM] = shopCopy
		log.Printf("Added shop %d to world", shopCopy.VNUM)

		// Link shop to room
		if room, ok := w.rooms[shopCopy.RoomVNUM]; ok {
			room.Shop = shopCopy
			log.Printf("Linked shop %d to room %d", shopCopy.VNUM, shopCopy.RoomVNUM)

			// Check if the shopkeeper mob is defined in the mobile prototypes
			if _, ok := w.mobiles[shopCopy.MobileVNUM]; ok {
				log.Printf("Shop %d has shopkeeper mobile VNUM %d", shopCopy.VNUM, shopCopy.MobileVNUM)

				// Check if the shopkeeper is already in the room
				shopkeeperExists := false
				for _, mob := range room.Characters {
					if mob.IsNPC && mob.Prototype != nil && mob.Prototype.VNUM == shopCopy.MobileVNUM {
						shopkeeperExists = true
						log.Printf("Shopkeeper %d already exists in room %d", shopCopy.MobileVNUM, room.VNUM)
						break
					}
				}

				// If the shopkeeper doesn't exist, create it
				if !shopkeeperExists {
					// Create the shopkeeper
					mob := w.CreateMobFromPrototype(shopCopy.MobileVNUM, room)
					if mob != nil {
						log.Printf("Created shopkeeper %s (VNUM %d) in room %d",
							mob.Name, shopCopy.MobileVNUM, room.VNUM)
					} else {
						log.Printf("Failed to create shopkeeper %d in room %d",
							shopCopy.MobileVNUM, room.VNUM)
					}
				}
			} else {
				log.Printf("Warning: Shop %d has invalid shopkeeper mobile VNUM %d", shopCopy.VNUM, shopCopy.MobileVNUM)
			}
		} else {
			log.Printf("Warning: Shop %d has invalid room VNUM %d", shopCopy.VNUM, shopCopy.RoomVNUM)
		}
	}
	log.Printf("Loaded %d shops into world", len(w.shops))

	return nil
}

// GetCharacter gets a character by name
func (w *World) GetCharacter(name string) (*types.Character, error) {
	w.mutex.RLock()
	defer func() {
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
	w.mutex.RLock()
	defer func() {
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

	// If the character has a RoomVNUM set (including room 0), try to use that
	// Note: In DikuMUD, room 0 is a valid room (The Void), so we check for >= 0
	// We use -1 or some other negative value to indicate "no saved room"
	if character.RoomVNUM >= 0 {
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
		// Try room 3001 (Temple of Midgaard) first, as per original DikuMUD
		foundRoom := w.GetRoom(3001)
		if foundRoom != nil {
			targetRoom = foundRoom
			targetRoomVNUM = targetRoom.VNUM
			log.Printf("AddCharacter: Target room determined (default 3001): %d", targetRoomVNUM)
		} else {
			log.Printf("AddCharacter: Starting room 3001 not found, trying room 0")
			// Try room 0 (The Void) as fallback
			foundRoom = w.GetRoom(0)
			if foundRoom != nil {
				targetRoom = foundRoom
				targetRoomVNUM = targetRoom.VNUM
				log.Printf("AddCharacter: Target room determined (fallback 0): %d", targetRoomVNUM)
			} else {
				log.Printf("AddCharacter: No starting room found for character %s. Character will not be placed in a room initially.", character.Name)
				// targetRoom remains nil, targetRoomVNUM remains 0
			}
		}
	}

	// --- Step 2: Acquire world lock and modify world state ---
	w.mutex.Lock()

	// Add character to the world map
	w.characters[character.Name] = character
	// Set the character's World field
	character.World = w

	// Initialize skills for player characters
	if !character.IsNPC {
		w.InitializeCharacterSkills(character)
	}

	// --- Step 3: If a target room was identified, lock it and add character ---
	if targetRoom != nil {
		// Re-fetch the room pointer *while holding the world lock* to ensure it wasn't deleted
		actualRoom := w.rooms[targetRoomVNUM]
		if actualRoom == targetRoom { // Pointer check for safety
			actualRoom.Lock() // Lock the specific room

			// Add character to the room's list
			character.InRoom = actualRoom
			actualRoom.Characters = append(actualRoom.Characters, character)
			// Update character's RoomVNUM just in case it wasn't set or was default
			character.RoomVNUM = actualRoom.VNUM

			log.Printf("AddCharacter: Placed character %s in room %d", character.Name, actualRoom.VNUM)
			actualRoom.Unlock() // Unlock the room
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
	w.mutex.Unlock()
}

// RemoveCharacter removes a character from the world
func (w *World) RemoveCharacter(character *types.Character) {
	w.mutex.Lock()

	// Remove character from the world map
	delete(w.characters, character.Name)

	// Remove character from the room (if they are in one)
	sourceRoom := character.InRoom // Get room pointer before potentially clearing it
	if sourceRoom != nil {
		sourceRoom.Lock() // Lock the specific room

		newChars := make([]*types.Character, 0, len(sourceRoom.Characters)-1)
		for _, ch := range sourceRoom.Characters {
			if ch != character {
				newChars = append(newChars, ch)
			}
		}
		sourceRoom.Characters = newChars
		character.InRoom = nil // Clear character's room reference

		sourceRoom.Unlock() // Unlock the room
	}

	w.mutex.Unlock()
}

// DeleteCharacter deletes a character from storage
func (w *World) DeleteCharacter(name string) error {
	// Acquire world lock
	w.mutex.Lock()
	defer w.mutex.Unlock()

	// Check if character is in the world
	character, ok := w.characters[name]
	if ok {
		// Remove character from the world
		delete(w.characters, name)

		// Remove character from the room (if they are in one)
		sourceRoom := character.InRoom // Get room pointer before potentially clearing it
		if sourceRoom != nil {
			sourceRoom.Lock() // Lock the specific room

			newChars := make([]*types.Character, 0, len(sourceRoom.Characters)-1)
			for _, ch := range sourceRoom.Characters {
				if ch != character {
					newChars = append(newChars, ch)
				}
			}
			sourceRoom.Characters = newChars
			character.InRoom = nil // Clear character's room reference

			sourceRoom.Unlock() // Unlock the room
		}
	}

	// Delete character from storage
	// Use the storage interface to delete the character
	return w.storage.DeleteCharacter(name)
}

// SaveCharacter saves a character to storage
func (w *World) SaveCharacter(character *types.Character) error {
	// Note: Saving might involve reading character state (like RoomVNUM).
	// Consider if a character-specific lock is needed, or if world RLock is sufficient.
	// For now, assume world write lock is okay, but this could be refined.
	w.mutex.Lock()
	defer func() {
		w.mutex.Unlock()
	}()

	// Update RoomVNUM before saving, ensure InRoom is consistent
	if character.InRoom != nil {
		character.RoomVNUM = character.InRoom.VNUM
	} else {
		// Only set RoomVNUM to 0 if it's not already -1 (new character)
		// This preserves the -1 value for new characters who haven't been placed yet
		if character.RoomVNUM != -1 {
			character.RoomVNUM = 0 // Character was in a room but is now out of world
		}
		// If RoomVNUM is -1, leave it as -1 (new character, will be placed in default room)
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

// AddDelay adds a delay to a character's actions
func (w *World) AddDelay(ch *types.Character, delay int) {
	// For now, just log the delay
	// In the future, this could be used to implement a proper combat round system
	// log.Printf("Adding delay of %d to %s", delay, ch.Name)
}

// GetCharacterInRoom finds a character in a room by name
func (w *World) GetCharacterInRoom(room *types.Room, name string) *types.Character {
	if room == nil || name == "" {
		return nil
	}

	name = strings.ToLower(name)
	for _, ch := range room.Characters {
		if strings.Contains(strings.ToLower(ch.Name), name) {
			return ch
		}
	}
	return nil
}

// Damage applies damage to a character
func (w *World) Damage(ch *types.Character, victim *types.Character, damage int, spellID int) {
	if ch == nil || victim == nil {
		return
	}

	// Apply damage
	victim.HP -= damage
	if victim.HP < 0 {
		victim.HP = 0
	}

	// Check if victim is dead
	if victim.HP <= 0 {
		// Handle death
		victim.Position = types.POS_DEAD
		victim.SendMessage("You are DEAD!\r\n")
		w.Act("$n is DEAD!", true, victim, nil, nil, types.TO_ROOM)
	}
}

// CharFromRoom removes a character from a room (DEPRECATED - use CharacterMove instead)
// This method is kept for compatibility with spell system but should be avoided
func (w *World) CharFromRoom(ch *types.Character) {
	if ch == nil || ch.InRoom == nil {
		return
	}

	room := ch.InRoom

	// Lock the room before modifying
	room.Lock()
	defer room.Unlock()

	// Remove character from room
	for i, rch := range room.Characters {
		if rch == ch {
			// Remove character from slice
			room.Characters = append(room.Characters[:i], room.Characters[i+1:]...)
			break
		}
	}

	ch.InRoom = nil
}

// CharToRoom adds a character to a room (DEPRECATED - use CharacterMove instead)
// This method is kept for compatibility with spell system but should be avoided
func (w *World) CharToRoom(ch *types.Character, room *types.Room) {
	if ch == nil || room == nil {
		return
	}

	// Lock the room before modifying
	room.Lock()
	defer room.Unlock()

	// Add character to room
	room.Characters = append(room.Characters, ch)
	ch.InRoom = room
	ch.RoomVNUM = room.VNUM
}

// ObjectToChar adds an object to a character's inventory
func (w *World) ObjectToChar(obj *types.ObjectInstance, ch *types.Character) {
	if obj == nil || ch == nil {
		return
	}

	// Add object to inventory
	ch.Inventory = append(ch.Inventory, obj)
	obj.CarriedBy = ch
}

// UnequipChar removes an object from a character's equipment
func (w *World) UnequipChar(ch *types.Character, position int) {
	if ch == nil || position < 0 || position >= len(ch.Equipment) {
		return
	}

	obj := ch.Equipment[position]
	if obj == nil {
		return
	}

	// Remove object from equipment
	ch.Equipment[position] = nil
	obj.WornBy = nil
	obj.WornOn = -1

	// Add object to inventory
	w.ObjectToChar(obj, ch)
}

// ExtractObj removes an object from the game
func (w *World) ExtractObj(obj *types.ObjectInstance) {
	if obj == nil {
		return
	}

	// Remove from character's inventory
	if obj.CarriedBy != nil {
		for i, o := range obj.CarriedBy.Inventory {
			if o == obj {
				obj.CarriedBy.Inventory = append(obj.CarriedBy.Inventory[:i], obj.CarriedBy.Inventory[i+1:]...)
				break
			}
		}
		obj.CarriedBy = nil
	}

	// Remove from character's equipment
	if obj.WornBy != nil {
		for i, o := range obj.WornBy.Equipment {
			if o == obj {
				obj.WornBy.Equipment[i] = nil
				break
			}
		}
		obj.WornBy = nil
	}

	// Remove from container
	if obj.InObj != nil {
		for i, o := range obj.InObj.Contains {
			if o == obj {
				obj.InObj.Contains = append(obj.InObj.Contains[:i], obj.InObj.Contains[i+1:]...)
				break
			}
		}
		obj.InObj = nil
	}

	// Remove from room
	if obj.InRoom != nil {
		for i, o := range obj.InRoom.Objects {
			if o == obj {
				obj.InRoom.Objects = append(obj.InRoom.Objects[:i], obj.InRoom.Objects[i+1:]...)
				break
			}
		}
		obj.InRoom = nil
	}

	// Remove contents
	for _, o := range obj.Contains {
		w.ExtractObj(o)
	}
}

// GetShop returns a shop by VNUM
func (w *World) GetShop(vnum int) *types.Shop {
	w.mutex.RLock()
	defer w.mutex.RUnlock()
	return w.shops[vnum]
}

// SendMessageToCharacter sends a message to a character
func (w *World) SendMessageToCharacter(character *types.Character, message string) {
	if character == nil {
		return
	}

	// Use the message handler if set
	if w.messageHandler != nil {
		w.messageHandler(character, message)
	}
}

// SetMessageHandler sets the message handler function
func (w *World) SetMessageHandler(handler func(*types.Character, string)) {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	w.messageHandler = handler
}

// CharacterMove moves a character from one room to another
func (w *World) CharacterMove(character *types.Character, destRoom *types.Room) {
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
		lockRoom1.Lock()
		defer func() {
			lockRoom1.Unlock()
		}()
	}
	if lockRoom2 != nil { // lockRoom2 is always different from lockRoom1 if not nil
		lockRoom2.Lock()
		defer func() {
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

// createMobFromPrototypeInternal creates a new mobile instance from a prototype
// This is an internal helper function used by both CreateMobFromPrototype and ProcessMobRespawns
// It does not add the mob to the world or equip it - that's done by the calling functions
func (w *World) createMobFromPrototypeInternal(_ int, mobProto *types.Mobile, baseHP, baseMana, baseMove int) *types.Character {
	// Create a new mobile instance
	mob := &types.Character{
		Name:          mobProto.Name,
		ShortDesc:     mobProto.ShortDesc,
		LongDesc:      mobProto.LongDesc,
		Description:   mobProto.Description,
		Level:         mobProto.Level,
		Sex:           mobProto.Sex,
		Class:         mobProto.Class,
		Race:          mobProto.Race,
		Gold:          mobProto.Gold,
		Experience:    mobProto.Experience,
		Alignment:     mobProto.Alignment,
		Position:      mobProto.Position,
		HP:            baseHP,
		MaxHitPoints:  baseHP,
		ManaPoints:    baseMana,
		MaxManaPoints: baseMana,
		MovePoints:    baseMove,
		MaxMovePoints: baseMove,
		ArmorClass:    mobProto.AC,
		HitRoll:       mobProto.HitRoll,
		DamRoll:       mobProto.DamRoll,
		ActFlags:      mobProto.ActFlags,
		IsNPC:         true,
		Prototype:     mobProto,
		World:         w,
		Inventory:     make([]*types.ObjectInstance, 0),
		Equipment:     make([]*types.ObjectInstance, types.NUM_WEARS),
		// RoomVNUM and InRoom will be set when placed
	}

	// Create a deep copy of the abilities array
	mob.Abilities = [6]int{
		mobProto.Abilities[0],
		mobProto.Abilities[1],
		mobProto.Abilities[2],
		mobProto.Abilities[3],
		mobProto.Abilities[4],
		mobProto.Abilities[5],
	}

	return mob
}

// CreateMobFromPrototype creates a new mobile instance from a prototype
// Assumes world lock is NOT held by caller.
func (w *World) CreateMobFromPrototype(vnum int, room *types.Room) *types.Character {
	w.mutex.Lock()

	// Get the mobile prototype (under world lock)
	mobProto := w.mobiles[vnum]
	if mobProto == nil {
		w.mutex.Unlock() // Unlock before returning
		log.Printf("Warning: Mobile prototype %d not found", vnum)
		return nil
	}

	// Calculate hit points based on dice values
	baseHP := 0
	if mobProto.Dice[0] > 0 && mobProto.Dice[1] > 0 {
		// Calculate HP using the dice formula: XdY+Z
		// Use actual dice rolls like the original DikuMUD
		baseHP = utils.Dice(mobProto.Dice[0], mobProto.Dice[1]) + mobProto.Dice[2]
		// Ensure minimum HP based on level
		minHP := mobProto.Level * 8
		if baseHP < minHP {
			baseHP = minHP
		}
	} else {
		// Default HP based on level if no dice are specified
		baseHP = mobProto.Level * 8
	}

	// Calculate mana and move points based on level
	baseMana := mobProto.Level * 10
	baseMove := mobProto.Level * 10

	// Create a new mobile instance using the helper function
	mob := w.createMobFromPrototypeInternal(vnum, mobProto, baseHP, baseMana, baseMove)

	// Log the mob stats for debugging
	log.Printf("CreateMobFromPrototype: Created mob %s (VNUM %d) with stats: Level=%d, HitRoll=%d, DamRoll=%d, AC=%v, Gold=%d, Exp=%d",
		mob.Name, vnum, mob.Level, mob.HitRoll, mob.DamRoll, mob.ArmorClass, mob.Gold, mob.Experience)

	// Add the mobile to the world map (under world lock)
	w.characters[mob.Name] = mob // Consider potential name collisions?

	// Add the mobile to the room (if provided)
	if room != nil {
		// Verify room pointer is still valid under world lock
		actualRoom := w.rooms[room.VNUM]
		if actualRoom == room {
			actualRoom.Lock() // Lock the room before modifying

			mob.InRoom = actualRoom
			actualRoom.Characters = append(actualRoom.Characters, mob)
			mob.RoomVNUM = actualRoom.VNUM

			actualRoom.Unlock() // Unlock the room
		} else {
			log.Printf("CreateMobFromPrototype: Warning - Room %d changed or removed before mob %d could be placed.", room.VNUM, vnum)
			// Mob exists in world map but not placed in room
		}
	}

	// Equip the mob with its default equipment
	w.equipMobFromPrototype(mob, mobProto)

	w.mutex.Unlock() // Release world lock

	return mob
}

// AddMobile adds a mobile prototype to the world
func (w *World) AddMobile(mobile *types.Mobile) {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	// Add the mobile to the world's mobile prototypes
	w.mobiles[mobile.VNUM] = mobile
	log.Printf("Added mobile prototype #%d (%s) to the world", mobile.VNUM, mobile.Name)
}

// ScheduleMobRespawn schedules a mob for respawning
// Assumes world lock is NOT held by caller.
func (w *World) ScheduleMobRespawn(mob *types.Character) {
	if mob == nil || !mob.IsNPC || mob.Prototype == nil {
		return
	}

	roomVNUMToRespawn := mob.RoomVNUM // Capture VNUM before locking

	w.mutex.Lock()

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
			actualRoom.Lock()

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
		} else {
			log.Printf("ScheduleMobRespawn: Warning - Room %d changed or removed before mob %s could be removed from it.", sourceRoom.VNUM, mob.Name)
			mob.InRoom = nil // Still clear mob's reference
		}
	} else {
		mob.InRoom = nil // Ensure reference is cleared if mob wasn't in a room
	}

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

			// Calculate hit points based on dice values
			baseHP := 0
			if mobProto.Dice[0] > 0 && mobProto.Dice[1] > 0 {
				// Calculate HP using the dice formula: XdY+Z
				// Use actual dice rolls like the original DikuMUD
				baseHP = utils.Dice(mobProto.Dice[0], mobProto.Dice[1]) + mobProto.Dice[2]
				// Ensure minimum HP based on level
				minHP := mobProto.Level * 8
				if baseHP < minHP {
					baseHP = minHP
				}
			} else {
				// Default HP based on level if no dice are specified
				baseHP = mobProto.Level * 8
			}

			// Calculate mana and move points based on level
			baseMana := mobProto.Level * 10
			baseMove := mobProto.Level * 10

			// Create a new mobile instance using the helper function
			mob := w.createMobFromPrototypeInternal(respawn.MobVNUM, mobProto, baseHP, baseMana, baseMove)

			// Log the mob stats for debugging
			log.Printf("ProcessMobRespawns: Created mob %s (VNUM %d) with stats: Level=%d, HitRoll=%d, DamRoll=%d, AC=%v, Gold=%d, Exp=%d",
				mob.Name, respawn.MobVNUM, mob.Level, mob.HitRoll, mob.DamRoll, mob.ArmorClass, mob.Gold, mob.Experience)

			// Add the mobile to the world map (under world lock)
			w.characters[mob.Name] = mob // Again, potential name collision?

			// Equip the mob with its default equipment
			w.equipMobFromPrototype(mob, mobProto)

			// Add the mobile to the room (needs room lock)
			room.Lock()

			mob.InRoom = room
			room.Characters = append(room.Characters, mob)
			mob.RoomVNUM = room.VNUM

			room.Unlock()

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

// equipMobFromPrototype equips a mob with its default equipment from the prototype
// Assumes world lock is already held by caller
func (w *World) equipMobFromPrototype(mob *types.Character, mobProto *types.Mobile) {
	// Check if the mob has any equipment defined in the prototype
	if len(mobProto.Equipment) == 0 {
		return
	}

	// Equip the mob with each item in the prototype's equipment list
	for _, eq := range mobProto.Equipment {
		// Check if the equipment should be equipped based on chance
		if eq.Chance > 0 && (eq.Chance == 100 || (rand.Intn(100) < eq.Chance)) {
			// Get the object prototype
			objProto := w.objects[eq.ObjectVNUM]
			if objProto == nil {
				log.Printf("Warning: Object prototype %d not found for mob equipment", eq.ObjectVNUM)
				continue
			}

			// Create a new object instance
			obj := w.CreateObjectFromPrototype(eq.ObjectVNUM)
			if obj == nil {
				log.Printf("Warning: Failed to create object %d for mob equipment", eq.ObjectVNUM)
				continue
			}

			// Equip the mob with the object
			obj.WornBy = mob
			obj.WornOn = eq.Position
			mob.Equipment[eq.Position] = obj
		}
	}
}
