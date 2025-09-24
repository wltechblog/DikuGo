package world

import (
	"log"

	"github.com/wltechblog/DikuGo/pkg/types"
	"github.com/wltechblog/DikuGo/pkg/utils"
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
	// Keep track of the last loaded mob (regardless of VNUM)
	var lastLoadedMob *types.Character

	log.Printf("Starting zone reset for zone %d: %s", zone.VNUM, zone.Name)

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
			log.Printf("Zone reset: Loading mobile VNUM %d into room %d (max: %d)", cmd.Arg1, cmd.Arg3, cmd.Arg2)
			mob := w.resetMobile(cmd.Arg1, cmd.Arg3, cmd.Arg2)
			if mob != nil {
				// Store the loaded mob for later equipment commands
				lastLoadedMob = mob
				log.Printf("Zone reset: Stored mob %s (VNUM %d) as last loaded mob", mob.Name, cmd.Arg1)
			}
		case 'O': // Load object
			// Arg1 = Object VNUM, Arg2 = Max number, Arg3 = Room VNUM
			w.resetObject(cmd.Arg1, cmd.Arg3, cmd.Arg2)
		case 'G': // Give object to mobile
			// Use the last loaded mob
			log.Printf("Zone reset: Giving object VNUM %d to last loaded mob", cmd.Arg1)
			if lastLoadedMob != nil {
				log.Printf("Zone reset: Using last loaded mob %s (VNUM %d)", lastLoadedMob.Name, lastLoadedMob.Prototype.VNUM)
				w.resetGiveObjectToMob(cmd.Arg1, lastLoadedMob)
			} else {
				log.Printf("Zone reset: No last loaded mob found, falling back to old method")
				w.resetGiveObject(cmd.Arg1, cmd.Arg2) // Fallback to old method
			}
		case 'E': // Equip mobile with object
			// Use the last loaded mob
			log.Printf("Zone reset: Equipping last loaded mob with object VNUM %d in position %d", cmd.Arg1, cmd.Arg3)
			if lastLoadedMob != nil {
				log.Printf("Zone reset: Using last loaded mob %s (VNUM %d)", lastLoadedMob.Name, lastLoadedMob.Prototype.VNUM)
				w.resetEquipObjectToMob(cmd.Arg1, lastLoadedMob, cmd.Arg3)
			} else {
				log.Printf("Zone reset: No last loaded mob found, falling back to old method")
				w.resetEquipObject(cmd.Arg1, cmd.Arg2, cmd.Arg3) // Fallback to old method
			}
		case 'P': // Put object in container
			w.resetPutObject(cmd.Arg1, cmd.Arg2, cmd.Arg3)
		case 'D': // Set door state
			w.resetDoor(cmd.Arg1, cmd.Arg2, cmd.Arg3)
		case 'R': // Remove object
			w.resetRemoveObject(cmd.Arg1, cmd.Arg2)
		}
	}
}

// resetMobile loads a mobile into a room and returns the created mob
func (w *World) resetMobile(mobVnum, roomVnum, maxExisting int) *types.Character {

	// Get the mobile prototype directly (we already have a lock)
	mobProto := w.mobiles[mobVnum]
	if mobProto == nil {
		log.Printf("Warning: Mobile %d not found", mobVnum)
		return nil
	}

	// Debug: Print the mobile prototype stats
	log.Printf("DEBUG: Mobile prototype #%d stats: Level=%d, HitRoll=%d, DamRoll=%d, AC=%v, Gold=%d, Exp=%d",
		mobVnum, mobProto.Level, mobProto.HitRoll, mobProto.DamRoll, mobProto.AC, mobProto.Gold, mobProto.Experience)

	// Get the room directly (we already have a lock)
	room := w.rooms[roomVnum]
	if room == nil {
		log.Printf("Warning: Room %d not found", roomVnum)
		return nil
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

	// Use the createMobFromPrototypeInternal function to create a character with all the correct stats
	mob := w.createMobFromPrototypeInternal(mobVnum, mobProto, baseHP, baseMana, baseMove)

	// Log the mob stats for debugging
	log.Printf("Created mob %s (VNUM %d) with stats: Level=%d, HitRoll=%d, DamRoll=%d, AC=%v, Gold=%d, Exp=%d",
		mob.Name, mobVnum, mob.Level, mob.HitRoll, mob.DamRoll, mob.ArmorClass, mob.Gold, mob.Experience)

	// Add the mobile to the world
	w.characters[mob.Name] = mob

	// Lock the room before modifying its character list
	room.Lock()

	// Re-check count and maxExisting within the room lock for safety, although world lock should suffice
	count = 0
	for _, ch := range room.Characters { // Iterate room characters safely under lock
		if ch.IsNPC && ch.Prototype != nil && ch.Prototype.VNUM == mobVnum {
			count++
		}
	}
	if count >= maxExisting {
		room.Unlock() // Unlock before returning
		return nil
	}

	// Add the mobile to the world characters map (already under world lock)
	w.characters[mob.Name] = mob

	// Add the mobile to the room (under room lock)
	mob.InRoom = room
	room.Characters = append(room.Characters, mob)

	room.Unlock() // Unlock the room
	log.Printf("resetMobile: Releasing lock for room %d", roomVnum)

	log.Printf("Loaded mobile %s (%d) into room %d", mob.Name, mobVnum, roomVnum)

	// Equip the mob with its default equipment
	w.equipMobFromPrototype(mob, mobProto)

	// Return the created mob
	return mob
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
	room.Lock()

	// Re-check count and maxExisting within the room lock for safety
	count = 0
	for _, objInstance := range room.Objects { // Iterate room objects safely under lock
		if objInstance.Prototype != nil && objInstance.Prototype.VNUM == objVnum {
			count++
		}
	}
	if count >= maxExisting {
		room.Unlock() // Unlock before returning
		return
	}

	// Create object instance (already under world lock)
	// Use '=' instead of ':=' because obj is already declared in this scope
	obj = w.CreateObjectFromPrototype(objVnum)
	if obj == nil {
		room.Unlock()                                              // Unlock before returning
		log.Printf("Warning: Failed to create object %d", objVnum) // Log moved here
		return
	}

	// Add object to room (under room lock)
	obj.InRoom = room
	room.Objects = append(room.Objects, obj)

	room.Unlock() // Unlock the room

	log.Printf("Loaded object %s (%d) into room %d", obj.Prototype.Name, objVnum, roomVnum)
}

// resetGiveObjectToMob gives an object to a specific mobile instance
func (w *World) resetGiveObjectToMob(objVnum int, mob *types.Character) {
	// Get the object prototype directly (we already have a lock)
	objProto := w.objects[objVnum]
	if objProto == nil {
		log.Printf("Warning: Object %d not found", objVnum)
		return
	}

	if mob == nil {
		log.Printf("Warning: Mobile is nil for giving object %d", objVnum)
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

	log.Printf("Gave object %s (%d) to mobile %s (VNUM %d)", obj.Prototype.Name, objVnum, mob.Name, mob.Prototype.VNUM)
}

// resetGiveObject gives an object to a mobile (fallback method)
func (w *World) resetGiveObject(objVnum, mobVnum int) {
	// Get the object prototype directly (we already have a lock)
	objProto := w.objects[objVnum]
	if objProto == nil {
		log.Printf("Warning: Object %d not found", objVnum)
		return
	}

	// Find the last loaded mobile of this type
	var mob *types.Character
	var lastMob *types.Character

	// First, try to find the mob by iterating through all characters
	for _, ch := range w.characters {
		if ch.IsNPC && ch.Prototype != nil && ch.Prototype.VNUM == mobVnum {
			lastMob = ch // Keep track of the last mob found with this VNUM
		}
	}

	// Use the last mob found with this VNUM
	mob = lastMob

	if mob == nil {
		log.Printf("Warning: Mobile %d not found for giving object %d", mobVnum, objVnum)
		return
	}

	log.Printf("Found mob %s (VNUM %d) for giving object %d", mob.Name, mobVnum, objVnum)

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

// resetEquipObjectToMob equips a specific mobile instance with an object
func (w *World) resetEquipObjectToMob(objVnum int, mob *types.Character, position int) {
	// Get the object prototype directly (we already have a lock)
	objProto := w.objects[objVnum]
	if objProto == nil {
		log.Printf("Warning: Object %d not found", objVnum)
		return
	}

	if mob == nil {
		log.Printf("Warning: Mobile is nil for equipping object %d", objVnum)
		return
	}

	// Create a new object instance
	obj := w.CreateObjectFromPrototype(objVnum)
	if obj == nil {
		log.Printf("Warning: Failed to create object %d", objVnum)
		return
	}

	// Equip the mobile with the object
	obj.WornBy = mob
	obj.WornOn = position
	mob.Equipment[position] = obj

	log.Printf("Equipped mobile %s (VNUM %d) with object %s (%d) in position %d",
		mob.Name, mob.Prototype.VNUM, obj.Prototype.Name, objVnum, position)
}

// resetEquipObject equips a mobile with an object (fallback method)
func (w *World) resetEquipObject(objVnum, mobVnum, position int) {
	// Get the object prototype directly (we already have a lock)
	objProto := w.objects[objVnum]
	if objProto == nil {
		log.Printf("Warning: Object %d not found", objVnum)
		return
	}

	// Find the last loaded mobile of this type
	var mob *types.Character
	var lastMob *types.Character

	// First, try to find the mob by iterating through all characters
	for _, ch := range w.characters {
		if ch.IsNPC && ch.Prototype != nil && ch.Prototype.VNUM == mobVnum {
			lastMob = ch // Keep track of the last mob found with this VNUM
		}
	}

	// Use the last mob found with this VNUM
	mob = lastMob

	if mob == nil {
		log.Printf("Warning: Mobile %d not found for equipping object %d", mobVnum, objVnum)
		return
	}

	log.Printf("Found mob %s (VNUM %d) for equipping object %d", mob.Name, mobVnum, objVnum)

	// Create a new object instance
	obj := w.CreateObjectFromPrototype(objVnum)
	if obj == nil {
		log.Printf("Warning: Failed to create object %d", objVnum)
		return
	}

	// Equip the mobile with the object
	obj.WornBy = mob
	obj.WornOn = position
	mob.Equipment[position] = obj

	log.Printf("Equipped mobile %s (%d) with object %s (%d) in position %d",
		mob.Name, mobVnum, obj.Prototype.Name, objVnum, position)
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
