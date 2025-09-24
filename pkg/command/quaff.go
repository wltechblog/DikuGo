package command

import (
	"fmt"
	"strings"

	"github.com/wltechblog/DikuGo/pkg/types"
	"github.com/wltechblog/DikuGo/pkg/world"
)

// DoQuaff handles the quaff command
func DoQuaff(character *types.Character, arguments string, world *world.World) error {
	// Check if arguments are empty
	if arguments == "" {
		return fmt.Errorf("quaff what?")
	}

	// Find the potion in inventory or held
	var potion *types.ObjectInstance
	arguments = strings.ToLower(arguments)

	// Check if holding a potion
	if character.Equipment[types.WEAR_HOLD] != nil {
		obj := character.Equipment[types.WEAR_HOLD]
		if strings.Contains(strings.ToLower(obj.Prototype.Name), arguments) {
			if obj.Prototype.Type == types.ITEM_POTION {
				potion = obj
			}
		}
	}

	// If not holding a potion, check inventory
	if potion == nil {
		for _, obj := range character.Inventory {
			if strings.Contains(strings.ToLower(obj.Prototype.Name), arguments) {
				if obj.Prototype.Type == types.ITEM_POTION {
					potion = obj
					break
				}
			}
		}
	}

	// Check if potion was found
	if potion == nil {
		return fmt.Errorf("you don't have that potion")
	}

	// Send messages
	world.Act("$n quaffs $p.", true, character, potion, nil, types.TO_ROOM)
	world.Act("You quaff $p which dissolves.", false, character, potion, nil, types.TO_CHAR)

	// Cast the spells in the potion
	for i := 1; i <= 3; i++ {
		spellID := potion.Prototype.Value[i]
		if spellID > 0 && spellID < types.MAX_SPELLS {
			level := potion.Prototype.Value[0]
			castSpell(world, spellID, level, character, "", types.SPELL_TYPE_POTION, character, nil)
		}
	}

	// Remove the potion
	if potion == character.Equipment[types.WEAR_HOLD] {
		world.UnequipChar(character, types.WEAR_HOLD)
	}
	world.ExtractObj(potion)

	return nil
}
