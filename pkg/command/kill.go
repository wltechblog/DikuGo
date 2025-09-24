package command

import (
	"fmt"
	"strings"

	"github.com/wltechblog/DikuGo/pkg/types"
)

// CombatManagerInterface defines the interface for a combat manager
type CombatManagerInterface interface {
	StartCombat(attacker, defender *types.Character) error
	StopCombat(character *types.Character)
	Update()
}

// KillCommand represents the kill command
type KillCommand struct {
	// CombatManager is the combat manager
	CombatManager CombatManagerInterface
}

// Execute executes the kill command
func (c *KillCommand) Execute(character *types.Character, args string) error {
	// Check if there are arguments
	if args == "" {
		return fmt.Errorf("kill whom?")
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
		return fmt.Errorf("you can't kill yourself")
	}

	// Check if the target is a player
	if !target.IsNPC {
		return fmt.Errorf("you can't kill other players")
	}

	// Start combat
	err := c.CombatManager.StartCombat(character, target)
	if err != nil {
		return err
	}

	// Send a message to the character
	return fmt.Errorf("you attack %s!\r\n", target.ShortDesc)
}

// Name returns the name of the command
func (c *KillCommand) Name() string {
	return "kill"
}

// Aliases returns the aliases of the command
func (c *KillCommand) Aliases() []string {
	return []string{"k"}
}

// MinPosition returns the minimum position required to execute the command
func (c *KillCommand) MinPosition() int {
	return types.POS_STANDING
}

// Level returns the minimum level required to execute the command
func (c *KillCommand) Level() int {
	return 0
}

// LogCommand returns whether the command should be logged
func (c *KillCommand) LogCommand() bool {
	return true
}

// findCharacterInRoom finds a character in a room by name
func findCharacterInRoom(room *types.Room, name string) *types.Character {
	// Convert the name to lowercase
	name = strings.ToLower(name)

	// Check each character in the room
	for _, ch := range room.Characters {
		// Check if the character's name contains the search name
		if strings.Contains(strings.ToLower(ch.Name), name) {
			return ch
		}

		// Check if the character's short description contains the search name
		if strings.Contains(strings.ToLower(ch.ShortDesc), name) {
			return ch
		}
	}

	return nil
}
