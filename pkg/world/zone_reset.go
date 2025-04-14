package world

import (
	"log"

	"github.com/wltechblog/DikuGo/pkg/types"
)

// ResetZones resets all zones that are due for a reset
func (w *World) ResetZones() {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	// Update zone ages
	for _, zone := range w.zones {
		zone.Age++
		if zone.ShouldReset() {
			log.Printf("Resetting zone %d: %s", zone.VNUM, zone.Name)
			w.resetZone(zone)
			zone.Age = 0
		}
	}
}

// resetZone resets a single zone
func (w *World) resetZone(zone *types.Zone) {
	// Process all zone commands
	for _, cmd := range zone.Commands {
		// Skip if the command is not enabled
		if cmd.IfFlag != 0 {
			// TODO: Implement if-flag logic
		}

		// Process the command based on its type
		switch cmd.Command {
		case 'M': // Load mobile
			// Arg1 = Mobile VNUM, Arg2 = Max number, Arg3 = Room VNUM
			w.resetMobile(cmd.Arg1, cmd.Arg3, cmd.Arg2)
		case 'O': // Load object
			// Arg1 = Object VNUM, Arg2 = Max number, Arg3 = Room VNUM
			w.resetObject(cmd.Arg1, cmd.Arg3, cmd.Arg2)
		case 'G': // Give object to mobile
			w.resetGiveObject(cmd.Arg1, cmd.Arg2)
		case 'E': // Equip mobile with object
			w.resetEquipObject(cmd.Arg1, cmd.Arg2, cmd.Arg3)
		case 'P': // Put object in container
			w.resetPutObject(cmd.Arg1, cmd.Arg2, cmd.Arg3)
		case 'D': // Set door state
			w.resetDoor(cmd.Arg1, cmd.Arg2, cmd.Arg3)
		case 'R': // Remove object
			w.resetRemoveObject(cmd.Arg1, cmd.Arg2)
		}
	}
}

// resetMobile loads a mobile into a room
func (w *World) resetMobile(mobVnum, roomVnum, maxExisting int) {

	// Get the mobile prototype directly (we already have a lock)
	mobProto := w.mobiles[mobVnum]
	if mobProto == nil {
		log.Printf("Warning: Mobile %d not found", mobVnum)
		return
	}

	// Get the room directly (we already have a lock)
	room := w.rooms[roomVnum]
	if room == nil {
		log.Printf("Warning: Room %d not found", roomVnum)
		return
	}

	// Count existing mobs of this type in the world
	count := 0
	for _, ch := range w.characters {
		if ch.IsNPC && ch.Prototype != nil && ch.Prototype.VNUM == mobVnum {
			count++
		}
	}

	// Check if we've reached the maximum number of this mob
	if count >= maxExisting {
		return
	}

	// Create a new mobile instance
	mob := &types.Character{
		Name:        mobProto.Name,
		ShortDesc:   mobProto.ShortDesc,
		LongDesc:    mobProto.LongDesc,
		Description: mobProto.Description,
		Level:       mobProto.Level,
		Gold:        mobProto.Gold,
		Position:    types.POS_STANDING,
		IsNPC:       true,
		Prototype:   mobProto,
		World:       w,
		Inventory:   make([]*types.ObjectInstance, 0),
		Equipment:   make([]*types.ObjectInstance, types.NUM_WEARS),
	}

	// Add the mobile to the world
	w.characters[mob.Name] = mob

	// Lock the room before modifying its character list
	log.Printf("resetMobile: Acquiring lock for room %d", roomVnum)
	room.Lock()
	log.Printf("resetMobile: Acquired lock for room %d", roomVnum)

	// Re-check count and maxExisting within the room lock for safety, although world lock should suffice
	count = 0
	for _, ch := range room.Characters { // Iterate room characters safely under lock
		if ch.IsNPC && ch.Prototype != nil && ch.Prototype.VNUM == mobVnum {
			count++
		}
	}
	if count >= maxExisting {
		room.Unlock() // Unlock before returning
		log.Printf("resetMobile: Releasing lock for room %d (max reached after re-check)", roomVnum)
		return
	}

	// Add the mobile to the world characters map (already under world lock)
	w.characters[mob.Name] = mob

	// Add the mobile to the room (under room lock)
	mob.InRoom = room
	room.Characters = append(room.Characters, mob)

	room.Unlock() // Unlock the room
	log.Printf("resetMobile: Releasing lock for room %d", roomVnum)

	log.Printf("Loaded mobile %s (%d) into room %d", mob.Name, mobVnum, roomVnum)
}

// resetObject loads an object into a room
func (w *World) resetObject(objVnum, roomVnum, maxExisting int) {
	// Get the object prototype directly (we already have a lock)
	objProto := w.objects[objVnum]
	if objProto == nil {
		log.Printf("Warning: Object %d not found", objVnum)
		return
	}

	// Get the room directly (we already have a lock)
	room := w.rooms[roomVnum]
	if room == nil {
		log.Printf("Warning: Room %d not found", roomVnum)
		return
	}

	// Count existing objects of this type in the room
	count := 0
	for _, obj := range room.Objects {
		if obj.Prototype != nil && obj.Prototype.VNUM == objVnum {
			count++
		}
	}

	// Check if we've reached the maximum number of this object
	if count >= maxExisting {
		return
	}

	// Create a new object instance
	obj := w.CreateObjectFromPrototype(objVnum)
	if obj == nil {
		log.Printf("Warning: Failed to create object %d", objVnum)
		return
	}

	// Lock the room before modifying its object list
	log.Printf("resetObject: Acquiring lock for room %d", roomVnum)
	room.Lock()
	log.Printf("resetObject: Acquired lock for room %d", roomVnum)

	// Re-check count and maxExisting within the room lock for safety
	count = 0
	for _, objInstance := range room.Objects { // Iterate room objects safely under lock
		if objInstance.Prototype != nil && objInstance.Prototype.VNUM == objVnum {
			count++
		}
	}
	if count >= maxExisting {
		room.Unlock() // Unlock before returning
		log.Printf("resetObject: Releasing lock for room %d (max reached after re-check)", roomVnum)
		return
	}

	// Create object instance (already under world lock)
	// Use '=' instead of ':=' because obj is already declared in this scope
	obj = w.CreateObjectFromPrototype(objVnum)
	if obj == nil {
		room.Unlock() // Unlock before returning
		log.Printf("resetObject: Releasing lock for room %d (obj creation failed)", roomVnum)
		log.Printf("Warning: Failed to create object %d", objVnum) // Log moved here
		return
	}

	// Add object to room (under room lock)
	obj.InRoom = room
	room.Objects = append(room.Objects, obj)

	room.Unlock() // Unlock the room
	log.Printf("resetObject: Releasing lock for room %d", roomVnum)

	log.Printf("Loaded object %s (%d) into room %d", obj.Prototype.Name, objVnum, roomVnum)
}

// resetGiveObject gives an object to a mobile
func (w *World) resetGiveObject(objVnum, mobVnum int) {
	// Get the object prototype directly (we already have a lock)
	objProto := w.objects[objVnum]
	if objProto == nil {
		log.Printf("Warning: Object %d not found", objVnum)
		return
	}

	// Find the last loaded mobile of this type
	var mob *types.Character
	for _, ch := range w.characters {
		if ch.IsNPC && ch.Prototype != nil && ch.Prototype.VNUM == mobVnum {
			mob = ch
		}
	}

	if mob == nil {
		log.Printf("Warning: Mobile %d not found for giving object %d", mobVnum, objVnum)
		return
	}

	// Create a new object instance
	obj := w.CreateObjectFromPrototype(objVnum)
	if obj == nil {
		log.Printf("Warning: Failed to create object %d", objVnum)
		return
	}

	// Add the object to the mobile's inventory
	obj.CarriedBy = mob
	mob.Inventory = append(mob.Inventory, obj)

	log.Printf("Gave object %s (%d) to mobile %s (%d)", obj.Prototype.Name, objVnum, mob.Name, mobVnum)
}

// resetEquipObject equips a mobile with an object
func (w *World) resetEquipObject(objVnum, mobVnum, position int) {
	// Get the object prototype directly (we already have a lock)
	objProto := w.objects[objVnum]
	if objProto == nil {
		log.Printf("Warning: Object %d not found", objVnum)
		return
	}

	// Find the last loaded mobile of this type
	var mob *types.Character
	for _, ch := range w.characters {
		if ch.IsNPC && ch.Prototype != nil && ch.Prototype.VNUM == mobVnum {
			mob = ch
		}
	}

	if mob == nil {
		log.Printf("Warning: Mobile %d not found for equipping object %d", mobVnum, objVnum)
		return
	}

	// Create a new object instance
	obj := w.CreateObjectFromPrototype(objVnum)
	if obj == nil {
		log.Printf("Warning: Failed to create object %d", objVnum)
		return
	}

	// Equip the mobile with the object
	obj.CarriedBy = mob
	mob.Equipment[position] = obj

	log.Printf("Equipped mobile %s (%d) with object %s (%d) in position %d", mob.Name, mobVnum, obj.Prototype.Name, objVnum, position)
}

// resetPutObject puts an object in a container
func (w *World) resetPutObject(objVnum, containerVnum, maxExisting int) {
	// TODO: Implement
}

// resetDoor sets the state of a door
func (w *World) resetDoor(roomVnum, direction, state int) {
	// TODO: Implement
}

// resetRemoveObject removes an object from a room
func (w *World) resetRemoveObject(objVnum, roomVnum int) {
	// TODO: Implement
}

// CreateObjectFromPrototype creates a new object instance from a prototype
func (w *World) CreateObjectFromPrototype(vnum int) *types.ObjectInstance {
	// Get the object prototype directly (we already have a lock)
	objProto := w.objects[vnum]
	if objProto == nil {
		log.Printf("Warning: Object prototype %d not found", vnum)
		return nil
	}

	// Create a new object instance
	obj := &types.ObjectInstance{
		Prototype: objProto,
	}

	return obj
}
