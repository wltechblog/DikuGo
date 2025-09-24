package command

import (
	"fmt"
	"strings"

	"github.com/wltechblog/DikuGo/pkg/types"
	"github.com/wltechblog/DikuGo/pkg/world"
)

// DoRecite handles the recite command
func DoRecite(character *types.Character, arguments string, world *world.World) error {
	// Check if arguments are empty
	if arguments == "" {
		return fmt.Errorf("recite what?")
	}

	// Split arguments into scroll and target
	args := strings.SplitN(arguments, " ", 2)
	scrollArg := strings.ToLower(args[0])
	targetArg := ""
	if len(args) > 1 {
		targetArg = strings.ToLower(args[1])
	}

	// Find the scroll in inventory or held
	var scroll *types.ObjectInstance
	
	// Check if holding a scroll
	if character.Equipment[types.WEAR_HOLD] != nil {
		obj := character.Equipment[types.WEAR_HOLD]
		if strings.Contains(strings.ToLower(obj.Prototype.Name), scrollArg) {
			if obj.Prototype.Type == types.ITEM_SCROLL {
				scroll = obj
			}
		}
	}

	// If not holding a scroll, check inventory
	if scroll == nil {
		for _, obj := range character.Inventory {
			if strings.Contains(strings.ToLower(obj.Prototype.Name), scrollArg) {
				if obj.Prototype.Type == types.ITEM_SCROLL {
					scroll = obj
					break
				}
			}
		}
	}

	// Check if scroll was found
	if scroll == nil {
		return fmt.Errorf("you don't have that scroll")
	}

	// Find the target
	var victim *types.Character
	var obj *types.ObjectInstance

	if targetArg != "" {
		// Try to find a character target
		victim = world.GetCharacterInRoom(character.InRoom, targetArg)
		
		// If no character found, try to find an object target
		if victim == nil {
			// Check inventory
			for _, o := range character.Inventory {
				if strings.Contains(strings.ToLower(o.Prototype.Name), targetArg) {
					obj = o
					break
				}
			}
			
			// Check room
			if obj == nil {
				for _, o := range character.InRoom.Objects {
					if strings.Contains(strings.ToLower(o.Prototype.Name), targetArg) {
						obj = o
						break
					}
				}
			}
			
			// Check equipment
			if obj == nil {
				for _, o := range character.Equipment {
					if o != nil && strings.Contains(strings.ToLower(o.Prototype.Name), targetArg) {
						obj = o
						break
					}
				}
			}
		}
	} else {
		// If no target specified, target self
		victim = character
	}

	// Send messages
	world.Act("$n recites $p.", true, character, scroll, nil, types.TO_ROOM)
	world.Act("You recite $p which dissolves.", false, character, scroll, nil, types.TO_CHAR)

	// Cast the spells in the scroll
	for i := 1; i <= 3; i++ {
		spellID := scroll.Prototype.Value[i]
		if spellID > 0 && spellID < types.MAX_SPELLS {
			level := scroll.Prototype.Value[0]
			castSpell(world, spellID, level, character, "", types.SPELL_TYPE_SCROLL, victim, obj)
		}
	}

	// Remove the scroll
	if scroll == character.Equipment[types.WEAR_HOLD] {
		world.UnequipChar(character, types.WEAR_HOLD)
	}
	world.ExtractObj(scroll)

	return nil
}
