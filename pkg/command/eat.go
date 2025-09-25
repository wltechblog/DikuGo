package command

import (
	"fmt"
	"strings"

	"github.com/wltechblog/DikuGo/pkg/types"
)

// EatCommand represents the eat command
type EatCommand struct{}

// Execute executes the eat command
func (c *EatCommand) Execute(character *types.Character, args string) error {
	// Check if there are arguments
	if args == "" {
		return fmt.Errorf("eat what?")
	}

	// Check if the character is in a room
	if character.InRoom == nil {
		return fmt.Errorf("you are not in a room")
	}

	// Find the food object in inventory
	var targetFood *types.ObjectInstance
	args = strings.ToLower(strings.TrimSpace(args))

	// Look for the food in inventory
	for _, obj := range character.Inventory {
		if strings.Contains(strings.ToLower(obj.Prototype.Name), args) ||
			strings.Contains(strings.ToLower(obj.Prototype.ShortDesc), args) {
			targetFood = obj
			break
		}
	}

	// Check if we found the object
	if targetFood == nil {
		return fmt.Errorf("you can't find it!")
	}

	// Check if it's actually food (admins can eat anything)
	if targetFood.Prototype.Type != types.ITEM_FOOD && character.Level < 22 {
		return fmt.Errorf("your stomach refuses to eat that!?!")
	}

	// Check if character is too full
	if len(character.Conditions) > types.COND_FULL && character.Conditions[types.COND_FULL] > 20 {
		return fmt.Errorf("you are too full to eat more!")
	}

	// Send messages to room and character
	for _, ch := range character.InRoom.Characters {
		if ch != character {
			ch.SendMessage(fmt.Sprintf("%s eats %s.\r\n", character.Name, targetFood.Prototype.ShortDesc))
		}
	}
	character.SendMessage(fmt.Sprintf("You eat %s.\r\n", targetFood.Prototype.ShortDesc))

	// Apply food effects - gain fullness
	if len(character.Conditions) > types.COND_FULL {
		foodValue := targetFood.Prototype.Value[0]
		if targetFood.Value[0] > 0 {
			foodValue = targetFood.Value[0] // Use instance value if set
		}

		character.Conditions[types.COND_FULL] += foodValue

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
	if targetFood.Prototype.Value[3] != 0 && character.Level < 21 {
		character.SendMessage("Ooups, it tasted rather strange ?!!?\r\n")

		// Send message to room
		for _, ch := range character.InRoom.Characters {
			if ch != character {
				ch.SendMessage(fmt.Sprintf("%s coughs and utters some strange sounds.\r\n", character.Name))
			}
		}

		// Apply poison effect
		poisonDuration := targetFood.Prototype.Value[0] * 2
		if poisonDuration > 0 {
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
	}

	// Remove the food from inventory and destroy it
	removeObjectFromInventory(character, targetFood)

	return nil
}

// Name returns the command name
func (c *EatCommand) Name() string {
	return "eat"
}

// Aliases returns command aliases
func (c *EatCommand) Aliases() []string {
	return []string{}
}

// MinPosition returns the minimum position required
func (c *EatCommand) MinPosition() int {
	return types.POS_RESTING
}

// Level returns the minimum level required
func (c *EatCommand) Level() int {
	return 0
}

// LogCommand returns whether this command should be logged
func (c *EatCommand) LogCommand() bool {
	return false
}

// addAffectToCharacter adds an affect to a character
func addAffectToCharacter(character *types.Character, newAffect *types.Affect) {
	// Add to the front of the affect list
	newAffect.Next = character.Affected
	character.Affected = newAffect

	// Set the bitvector flag
	character.AffectedBy |= newAffect.Bitvector
}

// removeObjectFromInventory removes an object from character's inventory
func removeObjectFromInventory(character *types.Character, obj *types.ObjectInstance) {
	for i, item := range character.Inventory {
		if item == obj {
			// Remove from slice
			character.Inventory = append(character.Inventory[:i], character.Inventory[i+1:]...)

			// Clear object references
			obj.CarriedBy = nil
			break
		}
	}
}
