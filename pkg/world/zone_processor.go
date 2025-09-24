package world

import (
	"log"

	"github.com/wltechblog/DikuGo/pkg/types"
)

// ProcessZoneCommands processes zone commands to set up default equipment for mob prototypes
func (w *World) ProcessZoneCommands() {
	log.Println("Processing zone commands to set up default equipment for mob prototypes...")

	// Process all zones
	for _, zone := range w.zones {
		// Process all zone commands
		var lastMobVnum int
		for _, cmd := range zone.Commands {
			switch cmd.Command {
			case 'M': // Load mobile
				// Remember the last loaded mob VNUM
				lastMobVnum = cmd.Arg1
			case 'E': // Equip mobile with object
				// Arg1 = Object VNUM, Arg2 = unused in this context, Arg3 = Position
				objVnum := cmd.Arg1
				position := cmd.Arg3
				mobVnum := lastMobVnum

				// Get the mobile prototype
				mobProto := w.mobiles[mobVnum]
				if mobProto == nil {
					log.Printf("Warning: Mobile %d not found for equipping object %d", mobVnum, objVnum)
					continue
				}

				// Get the object prototype
				objProto := w.objects[objVnum]
				if objProto == nil {
					log.Printf("Warning: Object %d not found for mob equipment", objVnum)
					continue
				}

				// Add the equipment to the mobile prototype
				// Check if this equipment is already defined
				alreadyDefined := false
				for _, eq := range mobProto.Equipment {
					if eq.ObjectVNUM == objVnum && eq.Position == position {
						alreadyDefined = true
						break
					}
				}

				if !alreadyDefined {
					mobProto.Equipment = append(mobProto.Equipment, types.MobEquipment{
						ObjectVNUM: objVnum,
						Position:   position,
						Chance:     100, // Default to 100% chance
					})
					log.Printf("Added equipment %d (position %d) to mobile %d", objVnum, position, mobVnum)
				}
			}
		}
	}

	log.Println("Finished processing zone commands.")
}
