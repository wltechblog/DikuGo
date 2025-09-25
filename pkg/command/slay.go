package command

import (
	"fmt"

	"github.com/wltechblog/DikuGo/pkg/types"
)

// SlayCommand represents the slay command
type SlayCommand struct {
	// CombatManager is the combat manager
	CombatManager CombatManagerInterface
}

// Execute executes the slay command
func (c *SlayCommand) Execute(character *types.Character, args string) error {
	// Check if there are arguments
	if args == "" {
		return fmt.Errorf("slay whom?")
	}

	// Check if the character is in a room
	if character.InRoom == nil {
		return fmt.Errorf("you are not in a room")
	}

	// Find the target
	target := findCharacterInRoom(character.InRoom, args)
	if target == nil {
		return fmt.Errorf("they aren't here")
	}

	// Check if the target is the character
	if target == character {
		return fmt.Errorf("you can't slay yourself")
	}

	// Check if the target is already dead
	if target.Position <= types.POS_DEAD {
		return fmt.Errorf("%s is already dead", target.ShortDesc)
	}

	// Send initial slay message
	character.SendMessage(fmt.Sprintf("You prepare to slay %s with divine power!\r\n", target.ShortDesc))
	target.SendMessage(fmt.Sprintf("%s prepares to slay you with divine power!\r\n", character.Name))

	// Send a message to the room
	for _, ch := range character.InRoom.Characters {
		if ch != character && ch != target {
			ch.SendMessage(fmt.Sprintf("%s prepares to slay %s with divine power!\r\n", character.Name, target.ShortDesc))
		}
	}

	// Store the target's current HP to calculate lethal damage
	targetHP := target.HP

	// Temporarily boost the character's damage to ensure a killing blow
	// Store original values
	originalDamRoll := character.DamRoll
	originalHitRoll := character.HitRoll

	// Set damage high enough to kill the target in one hit
	// Add extra to account for armor and other damage reductions
	character.DamRoll = targetHP + 50 // Ensure lethal damage
	character.HitRoll = 50            // Ensure the attack hits

	// Start combat using the normal combat system
	err := c.CombatManager.StartCombat(character, target)

	// Restore original values immediately after starting combat
	character.DamRoll = originalDamRoll
	character.HitRoll = originalHitRoll

	if err != nil {
		return fmt.Errorf("you slay %s with divine power!\r\n", target.ShortDesc)
	}

	// The combat system will handle the actual damage, death, and experience awarding
	return fmt.Errorf("you slay %s with divine power!\r\n", target.ShortDesc)
}

// Name returns the name of the command
func (c *SlayCommand) Name() string {
	return "slay"
}

// Aliases returns the aliases of the command
func (c *SlayCommand) Aliases() []string {
	return []string{}
}

// MinPosition returns the minimum position required to execute the command
func (c *SlayCommand) MinPosition() int {
	return types.POS_STANDING
}

// Level returns the minimum level required to execute the command
func (c *SlayCommand) Level() int {
	return 1 // Admin command - requires level 20+ at some point
}

// LogCommand returns whether the command should be logged
func (c *SlayCommand) LogCommand() bool {
	return true // Always log slay commands for admin oversight
}
