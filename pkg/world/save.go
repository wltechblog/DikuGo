package world

import (
	"fmt"
	"log"
)

// Save saves the world state
func (w *World) Save() error {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	log.Println("Saving world state...")

	// Save all characters
	for _, character := range w.characters {
		if !character.IsNPC {
			if err := w.storage.SaveCharacter(character); err != nil {
				return fmt.Errorf("failed to save character %s: %w", character.Name, err)
			}
		}
	}

	log.Println("World state saved")
	return nil
}

// ValidateRooms validates all room connections
func (w *World) ValidateRooms() error {
	w.mutex.RLock()
	defer w.mutex.RUnlock()

	log.Println("Validating room connections...")

	// Check all rooms
	var invalidExits int
	for _, room := range w.rooms {
		// Check all exits
		for dir, exit := range room.Exits {
			if exit != nil && exit.DestVnum != -1 {
				// Check if the destination room exists
				destRoom := w.GetRoom(exit.DestVnum)
				if destRoom == nil {
					log.Printf("Room %d has invalid exit %d to non-existent room %d", room.VNUM, dir, exit.DestVnum)
					invalidExits++
				}
			}
		}
	}

	if invalidExits > 0 {
		return fmt.Errorf("found %d invalid exits", invalidExits)
	}

	log.Println("All room connections are valid")
	return nil
}
