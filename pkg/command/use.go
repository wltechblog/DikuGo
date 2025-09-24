package command

import (
	"fmt"
	"strings"

	"github.com/wltechblog/DikuGo/pkg/types"
	"github.com/wltechblog/DikuGo/pkg/world"
)

// DoUse handles the use command for wands and staves
func DoUse(character *types.Character, arguments string, world *world.World) error {
	// Check if arguments are empty
	if arguments == "" {
		return fmt.Errorf("use what?")
	}

	// Split arguments into item and target
	args := strings.SplitN(arguments, " ", 2)
	itemArg := strings.ToLower(args[0])
	targetArg := ""
	if len(args) > 1 {
		targetArg = strings.ToLower(args[1])
	}

	// Check if holding an item
	if character.Equipment[types.WEAR_HOLD] == nil {
		return fmt.Errorf("you are not holding anything")
	}

	// Check if the held item matches the argument
	item := character.Equipment[types.WEAR_HOLD]
	if !strings.Contains(strings.ToLower(item.Prototype.Name), itemArg) {
		return fmt.Errorf("you are not holding that")
	}

	// Check if the item is a staff or wand
	if item.Prototype.Type == types.ITEM_STAFF {
		// Handle staff
		return useStaff(character, item, world)
	} else if item.Prototype.Type == types.ITEM_WAND {
		// Handle wand
		return useWand(character, item, targetArg, world)
	} else {
		return fmt.Errorf("use is normally only for wands and staves")
	}
}

// useStaff handles using a staff
func useStaff(ch *types.Character, staff *types.ObjectInstance, world *world.World) error {
	// Send messages
	world.Act("$n taps $p three times on the ground.", true, ch, staff, nil, types.TO_ROOM)
	world.Act("You tap $p three times on the ground.", false, ch, staff, nil, types.TO_CHAR)

	// Check if there are charges left
	if staff.Prototype.Value[2] <= 0 {
		ch.SendMessage("The staff seems powerless.\r\n")
		return nil
	}

	// Decrement charges
	staff.Prototype.Value[2]--

	// Get spell ID
	spellID := staff.Prototype.Value[3]

	// Check if spell ID is valid
	if spellID <= 0 || spellID >= types.MAX_SPELLS {
		ch.SendMessage("The staff seems to have a magical malfunction.\r\n")
		return nil
	}

	// Cast the spell on everyone in the room except the caster
	for _, victim := range ch.InRoom.Characters {
		if victim != ch {
			DoCast(ch, fmt.Sprintf("'%s' %s", types.GetSpellName(spellID), victim.Name), world)
		}
	}

	return nil
}

// useWand handles using a wand
func useWand(ch *types.Character, wand *types.ObjectInstance, targetArg string, world *world.World) error {
	// Check if target argument is provided
	if targetArg == "" {
		return fmt.Errorf("what should the wand be pointed at?")
	}

	// Find the target (character or object)
	var victim *types.Character
	var obj *types.ObjectInstance

	// Try to find a character target
	victim = world.GetCharacterInRoom(ch.InRoom, targetArg)

	// If no character found, try to find an object target
	if victim == nil {
		// Check inventory
		for _, o := range ch.Inventory {
			if strings.Contains(strings.ToLower(o.Prototype.Name), targetArg) {
				obj = o
				break
			}
		}

		// Check room
		if obj == nil {
			for _, o := range ch.InRoom.Objects {
				if strings.Contains(strings.ToLower(o.Prototype.Name), targetArg) {
					obj = o
					break
				}
			}
		}

		// Check equipment
		if obj == nil {
			for _, o := range ch.Equipment {
				if o != nil && strings.Contains(strings.ToLower(o.Prototype.Name), targetArg) {
					obj = o
					break
				}
			}
		}
	}

	// Check if a target was found
	if victim == nil && obj == nil {
		return fmt.Errorf("nothing by that name here")
	}

	// Send messages
	if victim != nil {
		world.Act("$n points $p at $N.", true, ch, wand, victim, types.TO_ROOM)
		world.Act("You point $p at $N.", false, ch, wand, victim, types.TO_CHAR)
	} else {
		world.Act("$n points $p at something.", true, ch, wand, nil, types.TO_ROOM)
		world.Act("You point $p at it.", false, ch, wand, nil, types.TO_CHAR)
	}

	// Check if there are charges left
	if wand.Prototype.Value[2] <= 0 {
		ch.SendMessage("The wand seems powerless.\r\n")
		return nil
	}

	// Decrement charges
	wand.Prototype.Value[2]--

	// Get spell ID
	spellID := wand.Prototype.Value[3]

	// Check if spell ID is valid
	if spellID <= 0 || spellID >= types.MAX_SPELLS {
		ch.SendMessage("The wand seems to have a magical malfunction.\r\n")
		return nil
	}

	// Cast the spell using the cast command
	DoCast(ch, fmt.Sprintf("'%s' %s", types.GetSpellName(spellID), targetArg), world)

	return nil
}
