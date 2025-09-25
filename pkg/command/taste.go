package command

import (
	"fmt"
	"strings"

	"github.com/wltechblog/DikuGo/pkg/types"
)

// TasteCommand represents the taste command
type TasteCommand struct{}

// Execute executes the taste command
func (c *TasteCommand) Execute(character *types.Character, args string) error {
	// Check if there are arguments
	if args == "" {
		return fmt.Errorf("taste what?")
	}

	// Check if the character is in a room
	if character.InRoom == nil {
		return fmt.Errorf("you are not in a room")
	}

	// Find the object in inventory
	var targetObject *types.ObjectInstance
	args = strings.ToLower(strings.TrimSpace(args))

	// Look for the object in inventory
	for _, obj := range character.Inventory {
		if strings.Contains(strings.ToLower(obj.Prototype.Name), args) ||
			strings.Contains(strings.ToLower(obj.Prototype.ShortDesc), args) {
			targetObject = obj
			break
		}
	}

	// Check if we found the object
	if targetObject == nil {
		return fmt.Errorf("you can't find it!")
	}

	// If it's a drink container, use sip command instead
	if targetObject.Prototype.Type == types.ITEM_DRINKCON {
		sipCmd := &SipCommand{}
		return sipCmd.Execute(character, args)
	}

	// Check if it's food
	if targetObject.Prototype.Type != types.ITEM_FOOD {
		return fmt.Errorf("taste that?!? Your stomach refuses!")
	}

	// Send messages to room and character
	for _, ch := range character.InRoom.Characters {
		if ch != character {
			ch.SendMessage(fmt.Sprintf("%s tastes %s.\r\n", character.Name, targetObject.Prototype.ShortDesc))
		}
	}
	character.SendMessage(fmt.Sprintf("You taste %s.\r\n", targetObject.Prototype.ShortDesc))

	// Apply small food effect - gain 1 fullness
	if len(character.Conditions) > types.COND_FULL {
		character.Conditions[types.COND_FULL] += 1

		// Cap fullness at 24
		if character.Conditions[types.COND_FULL] > 24 {
			character.Conditions[types.COND_FULL] = 24
		}

		// Send fullness message if now full
		if character.Conditions[types.COND_FULL] > 20 {
			character.SendMessage("You are full.\r\n")
		}
	}

	// Check for poison (Value[3] indicates poisoned food)
	if targetObject.Prototype.Value[3] != 0 && !isAffectedBy(character, types.AFF_POISON) {
		character.SendMessage("Ooups, it did not taste good at all!\r\n")

		// Apply poison effect (shorter duration for taste)
		poisonDuration := 2 // Fixed duration for taste

		// Create poison affect
		affect := &types.Affect{
			Type:      33, // SPELL_POISON from spells.go
			Duration:  poisonDuration,
			Modifier:  0,
			Location:  types.APPLY_NONE,
			Bitvector: types.AFF_POISON,
		}

		// Add poison affect to character
		addAffectToCharacter(character, affect)
	}

	// Reduce food value by 1 (tasting consumes a small amount)
	foodValue := targetObject.Prototype.Value[0]
	if targetObject.Value[0] > 0 {
		foodValue = targetObject.Value[0] // Use instance value if set
	}

	targetObject.Value[0] = foodValue - 1
	if targetObject.Value[0] < 0 {
		targetObject.Value[0] = 0
	}

	// If nothing left, destroy the food
	if targetObject.Value[0] <= 0 {
		character.SendMessage("There is nothing left now.\r\n")
		removeObjectFromInventory(character, targetObject)
	}

	return nil
}

// Name returns the command name
func (c *TasteCommand) Name() string {
	return "taste"
}

// Aliases returns command aliases
func (c *TasteCommand) Aliases() []string {
	return []string{}
}

// MinPosition returns the minimum position required
func (c *TasteCommand) MinPosition() int {
	return types.POS_RESTING
}

// Level returns the minimum level required
func (c *TasteCommand) Level() int {
	return 0
}

// LogCommand returns whether this command should be logged
func (c *TasteCommand) LogCommand() bool {
	return false
}

// isAffectedBy checks if a character is affected by a specific bitvector
func isAffectedBy(character *types.Character, bitvector int64) bool {
	return (character.AffectedBy & bitvector) != 0
}
