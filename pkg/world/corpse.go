package world

import (
	"fmt"
	"log"

	"github.com/wltechblog/DikuGo/pkg/types"
)

// Constants for corpse timers (in minutes)
const (
	MAX_NPC_CORPSE_TIME = 15 // 15 minutes for NPC corpses
	MAX_PC_CORPSE_TIME  = 60 // 60 minutes for PC corpses
)

// MakeCorpse creates a corpse from a character and places it in the room
func (w *World) MakeCorpse(ch *types.Character) *types.ObjectInstance {
	if ch == nil || ch.InRoom == nil {
		log.Printf("MakeCorpse: Invalid character or room")
		return nil
	}

	// Create a new object instance for the corpse
	corpse := &types.ObjectInstance{
		Prototype: &types.Object{
			VNUM:        -1, // Special VNUM for corpses
			Name:        "corpse",
			ShortDesc:   fmt.Sprintf("corpse of %s", ch.ShortDesc),
			Description: fmt.Sprintf("Corpse of %s is lying here.", ch.ShortDesc),
			Type:        types.ITEM_CONTAINER,
			WearFlags:   types.ITEM_WEAR_TAKE,
			Value:       [4]int{0, 0, 0, 1}, // Value[3] = 1 identifies it as a corpse
		},
		Contains: make([]*types.ObjectInstance, 0),
		Timer:    MAX_NPC_CORPSE_TIME,
	}

	// If it's a player corpse, set a longer timer
	if !ch.IsNPC {
		corpse.Timer = MAX_PC_CORPSE_TIME
	}

	// Transfer character's inventory to the corpse
	for _, item := range ch.Inventory {
		// Remove the item from the character
		item.CarriedBy = nil

		// Add the item to the corpse
		item.InObj = corpse
		corpse.Contains = append(corpse.Contains, item)
	}

	// Clear the character's inventory
	ch.Inventory = make([]*types.ObjectInstance, 0)

	// Transfer character's equipment to the corpse
	for i, item := range ch.Equipment {
		if item != nil {
			// Unequip the item
			item.WornBy = nil
			item.WornOn = -1
			ch.Equipment[i] = nil

			// Add the item to the corpse
			item.InObj = corpse
			corpse.Contains = append(corpse.Contains, item)
		}
	}

	// If the character has gold, create money in the corpse
	if ch.Gold > 0 {
		money := w.CreateMoney(ch.Gold)
		if money != nil {
			money.InObj = corpse
			corpse.Contains = append(corpse.Contains, money)
		}
		ch.Gold = 0
	}

	// Add the corpse to the room
	corpse.InRoom = ch.InRoom
	ch.InRoom.Objects = append(ch.InRoom.Objects, corpse)

	log.Printf("Created corpse of %s in room %d", ch.Name, ch.InRoom.VNUM)
	return corpse
}

// CreateMoney creates a money object with the specified amount
func (w *World) CreateMoney(amount int) *types.ObjectInstance {
	if amount <= 0 {
		return nil
	}

	var name, shortDesc, desc string
	if amount == 1 {
		name = "coin gold"
		shortDesc = "a gold coin"
		desc = "One miserable gold coin."
	} else {
		name = "coins gold"
		shortDesc = fmt.Sprintf("%d gold coins", amount)
		desc = fmt.Sprintf("A pile of %d gold coins.", amount)
	}

	money := &types.ObjectInstance{
		Prototype: &types.Object{
			VNUM:        -2, // Special VNUM for money
			Name:        name,
			ShortDesc:   shortDesc,
			Description: desc,
			Type:        types.ITEM_MONEY,
			Value:       [4]int{amount, 0, 0, 0}, // Value[0] = amount
		},
	}

	return money
}

// PulseCorpses updates all corpses in the world
func (w *World) PulseCorpses() {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	// No need to get current time for now

	// Check all rooms for corpses
	for _, room := range w.rooms {
		room.Lock()

		// Create a list of corpses to remove
		var corpsesToRemove []*types.ObjectInstance

		// Check all objects in the room
		for _, obj := range room.Objects {
			// Check if the object is a corpse
			if obj.Prototype != nil && obj.Prototype.Type == types.ITEM_CONTAINER && obj.Prototype.Value[3] == 1 {
				// Decrement the timer
				obj.Timer--

				// If the timer has expired, mark the corpse for removal
				if obj.Timer <= 0 {
					corpsesToRemove = append(corpsesToRemove, obj)
				}
			}
		}

		// Remove expired corpses
		for _, corpse := range corpsesToRemove {
			w.removeCorpse(corpse, room)
		}

		room.Unlock()
	}
}

// removeCorpse removes a corpse and its contents from the room
func (w *World) removeCorpse(corpse *types.ObjectInstance, room *types.Room) {
	if corpse == nil || room == nil {
		return
	}

	// Move the contents of the corpse to the room
	for _, item := range corpse.Contains {
		// Remove the item from the corpse
		item.InObj = nil

		// Add the item to the room
		item.InRoom = room
		room.Objects = append(room.Objects, item)
	}

	// Clear the corpse's contents
	corpse.Contains = nil

	// Remove the corpse from the room
	for i, obj := range room.Objects {
		if obj == corpse {
			// Remove the corpse from the room's objects
			room.Objects = append(room.Objects[:i], room.Objects[i+1:]...)
			break
		}
	}

	log.Printf("Corpse decayed in room %d", room.VNUM)
}
