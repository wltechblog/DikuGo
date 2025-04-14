package command

import (
	"fmt"

	"github.com/wltechblog/DikuGo/pkg/types"
)

// QuitCommand represents the quit command
type QuitCommand struct{}

// Execute executes the quit command
func (c *QuitCommand) Execute(character *types.Character, args string) error {
	// Check if the character is fighting
	if character.Fighting != nil {
		return fmt.Errorf("no way! you're fighting for your life!")
	}

	// Send a message to the character
	// TODO: Implement a way to send messages to characters
	// For now, just return a message as an error
	message := fmt.Sprintf("Goodbye, %s!\r\n", character.Name)

	// Set a flag to indicate that the character is quitting
	// TODO: Implement a way to indicate that the character is quitting
	// For now, just return a special error
	return fmt.Errorf("%sQUIT", message)
}

// Name returns the name of the command
func (c *QuitCommand) Name() string {
	return "quit"
}

// Aliases returns the aliases of the command
func (c *QuitCommand) Aliases() []string {
	return []string{}
}

// MinPosition returns the minimum position required to execute the command
func (c *QuitCommand) MinPosition() int {
	return types.POS_SLEEPING
}

// Level returns the minimum level required to execute the command
func (c *QuitCommand) Level() int {
	return 0
}

// LogCommand returns whether the command should be logged
func (c *QuitCommand) LogCommand() bool {
	return false
}
