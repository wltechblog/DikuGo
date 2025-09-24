package command

import (
	"fmt"
	"strings"

	"github.com/wltechblog/DikuGo/pkg/types"
)

// WhoCommand represents the who command
type WhoCommand struct{}

// Execute executes the who command
func (c *WhoCommand) Execute(character *types.Character, args string) error {
	// Get the current list of characters from the world
	world, ok := character.World.(interface {
		GetCharacters() map[string]*types.Character
	})
	if !ok {
		return fmt.Errorf("world interface not available")
	}

	characters := world.GetCharacters()

	// Build the who list
	var sb strings.Builder

	sb.WriteString("\r\nPlayers online:\r\n")
	sb.WriteString("---------------\r\n")

	// Add each character to the list
	for _, ch := range characters {
		if !ch.IsNPCFlag() {
			sb.WriteString(fmt.Sprintf("[%2d] %s%s\r\n", ch.Level, ch.Name, ch.Title))
		}
	}

	sb.WriteString("\r\n")

	// Send the who list to the character
	// TODO: Implement a way to send messages to characters
	// For now, just return the who list as an error
	return fmt.Errorf("%s", sb.String())
}

// Name returns the name of the command
func (c *WhoCommand) Name() string {
	return "who"
}

// Aliases returns the aliases of the command
func (c *WhoCommand) Aliases() []string {
	return []string{}
}

// MinPosition returns the minimum position required to execute the command
func (c *WhoCommand) MinPosition() int {
	return types.POS_SLEEPING
}

// Level returns the minimum level required to execute the command
func (c *WhoCommand) Level() int {
	return 0
}

// LogCommand returns whether the command should be logged
func (c *WhoCommand) LogCommand() bool {
	return false
}
