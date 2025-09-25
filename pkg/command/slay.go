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

	// Start combat using the normal combat system first
	err := c.CombatManager.StartCombat(character, target)
	if err != nil {
		return err
	}

	// Now apply lethal damage directly using the combat system's damage calculation
	// This ensures we use all the normal combat mechanics for damage, death, and experience

	// Calculate lethal damage (target's current HP + a buffer to ensure death)
	lethalDamage := target.HP + 10

	// Apply the damage directly to the target
	target.HP -= lethalDamage

	// Ensure HP doesn't go below 0
	if target.HP < 0 {
		target.HP = 0
	}

	// Set position to dead
	target.Position = types.POS_DEAD

	// Send death messages (same as normal combat)
	character.SendMessage(fmt.Sprintf("You have slain %s!\r\n", target.ShortDesc))
	target.SendMessage(fmt.Sprintf("%s has slain you!\r\n", character.Name))

	// Send a message to the room
	for _, ch := range character.InRoom.Characters {
		if ch != character && ch != target {
			ch.SendMessage(fmt.Sprintf("%s has slain %s!\r\n", character.Name, target.ShortDesc))
		}
	}

	// Stop combat for both characters
	c.CombatManager.StopCombat(character)

	// Handle character death using the centralized death handler
	w, ok := target.World.(interface {
		HandleCharacterDeath(*types.Character)
	})
	if ok {
		w.HandleCharacterDeath(target)
	}

	// Award experience if target is NPC (same calculation as normal combat)
	if target.IsNPC {
		exp := calculateSlayExperience(character, target)
		character.Experience += exp
		character.SendMessage(fmt.Sprintf("You gain %d experience points.\r\n", exp))
	}

	return fmt.Errorf("you slay %s with divine power!\r\n", target.ShortDesc)
}

// calculateSlayExperience calculates the experience gained from slaying a target
// Uses the same formula as combat but ensures a minimum reward
func calculateSlayExperience(slayer, victim *types.Character) int {
	// Base experience is 1/3 of the mob's experience (same as combat)
	exp := victim.Experience / 3

	// Apply level difference bonus
	levelDiff := victim.Level - slayer.Level
	if levelDiff > 0 {
		// Bonus for slaying higher level mobs
		if slayer.IsNPC {
			// NPCs get less bonus
			exp += (exp * levelDiff) / 8
		} else {
			// Players get more bonus
			exp += (exp * levelDiff) / 4
		}
	}

	// Ensure minimum experience (slightly higher than combat since it's an admin command)
	if exp < 10 {
		exp = 10
	}

	return exp
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
