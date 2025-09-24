package world

import (
	"log"
	"math/rand"

	"github.com/wltechblog/DikuGo/pkg/ai"
	"github.com/wltechblog/DikuGo/pkg/types"
)

// InitAI initializes the AI system
func (w *World) InitAI() {
	// Create a new AI manager
	w.aiManager = ai.NewManager(w)
	log.Println("AI system initialized")
}

// GetMobiles returns all mobile characters in the world
func (w *World) GetMobiles() []*types.Character {
	w.mutex.RLock()
	defer w.mutex.RUnlock()

	// Convert the map to a slice, including only NPCs
	mobiles := make([]*types.Character, 0, len(w.characters))
	for _, char := range w.characters {
		if char.IsNPC {
			mobiles = append(mobiles, char)
		}
	}

	return mobiles
}

// FindObjectsInRoom returns all objects in a room
func (w *World) FindObjectsInRoom(room *types.Room) []*types.ObjectInstance {
	if room == nil {
		return nil
	}

	// Acquire read lock for the room before accessing its objects
	room.RLock()
	defer func() {
		room.RUnlock()
	}()

	// Copy the slice to avoid returning a reference that can be modified unsafely
	objectsCopy := make([]*types.ObjectInstance, len(room.Objects))
	copy(objectsCopy, room.Objects)
	return objectsCopy
}

// MoveCharacter moves a character from one room to another
func (w *World) MoveCharacter(character *types.Character, room *types.Room) error {
	if character == nil || room == nil {
		return nil
	}

	// Use the existing CharacterMove method
	w.CharacterMove(character, room)
	return nil
}

// GetRandomExitRoom returns a random exit room from the given room
func (w *World) GetRandomExitRoom(room *types.Room) (*types.Room, int) {
	if room == nil {
		return nil, -1
	}

	// Find all valid exits
	var validExits []int
	for dir, exit := range room.Exits {
		if exit != nil && exit.DestVnum != -1 {
			// Check if the destination room exists
			if destRoom := w.GetRoom(exit.DestVnum); destRoom != nil {
				validExits = append(validExits, dir)
			}
		}
	}

	// If there are no valid exits, return nil
	if len(validExits) == 0 {
		return nil, -1
	}

	// Pick a random exit
	dir := validExits[rand.Intn(len(validExits))]
	exit := room.Exits[dir]
	destRoom := w.GetRoom(exit.DestVnum)

	return destRoom, dir
}

// TickAI updates the AI system
func (w *World) TickAI() {
	// Skip if the AI manager is not initialized
	if w.aiManager == nil {
		return
	}

	// Update the AI
	w.aiManager.Tick()
}
