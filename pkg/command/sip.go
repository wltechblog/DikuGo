package command

import (
	"fmt"
	"strings"

	"github.com/wltechblog/DikuGo/pkg/types"
)

// SipCommand represents the sip command
type SipCommand struct{}

// Execute executes the sip command
func (c *SipCommand) Execute(character *types.Character, args string) error {
	// Check if there are arguments
	if args == "" {
		return fmt.Errorf("sip what?")
	}

	// Check if the character is in a room
	if character.InRoom == nil {
		return fmt.Errorf("you are not in a room")
	}

	// Find the drink container in inventory
	var targetDrink *types.ObjectInstance
	args = strings.ToLower(strings.TrimSpace(args))

	// Look for the drink container in inventory
	for _, obj := range character.Inventory {
		if strings.Contains(strings.ToLower(obj.Prototype.Name), args) ||
			strings.Contains(strings.ToLower(obj.Prototype.ShortDesc), args) {
			targetDrink = obj
			break
		}
	}

	// Check if we found the object
	if targetDrink == nil {
		return fmt.Errorf("you can't find it!")
	}

	// Check if it's actually a drink container
	if targetDrink.Prototype.Type != types.ITEM_DRINKCON {
		return fmt.Errorf("you can't sip from that!")
	}

	// Check if character is too drunk to sip
	if len(character.Conditions) > types.COND_DRUNK && character.Conditions[types.COND_DRUNK] > 10 {
		character.SendMessage("You simply fail to reach your mouth!\r\n")

		// Send message to room
		for _, ch := range character.InRoom.Characters {
			if ch != character {
				ch.SendMessage(fmt.Sprintf("%s tries to sip, but fails!\r\n", character.Name))
			}
		}
		return nil
	}

	// Check if container has liquid
	liquidAmount := targetDrink.Prototype.Value[1]
	if targetDrink.Value[1] > 0 {
		liquidAmount = targetDrink.Value[1] // Use instance value if set
	}

	if liquidAmount <= 0 {
		return fmt.Errorf("it is empty")
	}

	// Get liquid type
	liquidType := targetDrink.Prototype.Value[2]
	if targetDrink.Value[2] > 0 {
		liquidType = targetDrink.Value[2] // Use instance value if set
	}

	// Sip amount is always 1 unit
	sipAmount := 1

	// Apply drink effects (same as drink command but smaller amount)
	if len(character.Conditions) > types.COND_THIRST {
		// Apply thirst, hunger, and drunk effects based on liquid type
		drunkEffect := getDrinkEffect(liquidType, DRINK_EFFECT_DRUNK)
		fullEffect := getDrinkEffect(liquidType, DRINK_EFFECT_FULL)
		thirstEffect := getDrinkEffect(liquidType, DRINK_EFFECT_THIRST)

		character.Conditions[types.COND_DRUNK] += (drunkEffect * sipAmount) / 4
		character.Conditions[types.COND_FULL] += (fullEffect * sipAmount) / 4
		character.Conditions[types.COND_THIRST] += (thirstEffect * sipAmount) / 4

		// Cap conditions at reasonable limits
		if character.Conditions[types.COND_DRUNK] > 24 {
			character.Conditions[types.COND_DRUNK] = 24
		}
		if character.Conditions[types.COND_FULL] > 24 {
			character.Conditions[types.COND_FULL] = 24
		}
		if character.Conditions[types.COND_THIRST] > 24 {
			character.Conditions[types.COND_THIRST] = 24
		}
	}

	// Reduce liquid in container by sip amount
	targetDrink.Value[1] = liquidAmount - sipAmount
	if targetDrink.Value[1] < 0 {
		targetDrink.Value[1] = 0
	}

	// Send messages
	drinkName := getDrinkName(liquidType)
	character.SendMessage(fmt.Sprintf("You sip the %s.\r\n", drinkName))

	// Send message to room
	for _, ch := range character.InRoom.Characters {
		if ch != character {
			ch.SendMessage(fmt.Sprintf("%s sips %s from %s.\r\n",
				character.Name, drinkName, targetDrink.Prototype.ShortDesc))
		}
	}

	// Check for poison in the liquid (Value[3] indicates poisoned liquid)
	if targetDrink.Prototype.Value[3] != 0 && character.Level < 21 {
		character.SendMessage("Ooups, it tasted rather strange!\r\n")

		// Send message to room
		for _, ch := range character.InRoom.Characters {
			if ch != character {
				ch.SendMessage(fmt.Sprintf("%s coughs and utters some strange sounds.\r\n", character.Name))
			}
		}

		// Apply poison effect (shorter duration for sip)
		poisonDuration := 2 // Fixed duration for sip

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

	return nil
}

// Name returns the command name
func (c *SipCommand) Name() string {
	return "sip"
}

// Aliases returns command aliases
func (c *SipCommand) Aliases() []string {
	return []string{}
}

// MinPosition returns the minimum position required
func (c *SipCommand) MinPosition() int {
	return types.POS_RESTING
}

// Level returns the minimum level required
func (c *SipCommand) Level() int {
	return 0
}

// LogCommand returns whether this command should be logged
func (c *SipCommand) LogCommand() bool {
	return false
}
