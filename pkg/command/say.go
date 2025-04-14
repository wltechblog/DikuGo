package command

import (
	"fmt"

	"github.com/wltechblog/DikuGo/pkg/types"
)

// SayCommand represents the say command
type SayCommand struct{}

// Execute executes the say command
func (c *SayCommand) Execute(character *types.Character, args string) error {
	// Check if there are arguments
	if args == "" {
		return fmt.Errorf("say what?")
	}

	// Check if the character is in a room
	if character.InRoom == nil {
		return fmt.Errorf("you are not in a room")
	}

	// Get the room
	room := character.InRoom

	// Send the message to the character
	// TODO: Implement a way to send messages to characters
	// For now, just return the message as an error
	message := fmt.Sprintf("You say, '%s'\r\n", args)

	// Send the message to everyone else in the room
	for _, ch := range room.Characters {
		if ch != character {
			// TODO: Implement a way to send messages to characters
			// For now, just ignore this
			_ = fmt.Sprintf("%s says, '%s'\r\n", character.Name, args)
		}
	}

	return fmt.Errorf("%s", message)
}

// Name returns the name of the command
func (c *SayCommand) Name() string {
	return "say"
}

// Aliases returns the aliases of the command
func (c *SayCommand) Aliases() []string {
	return []string{"'"}
}

// MinPosition returns the minimum position required to execute the command
func (c *SayCommand) MinPosition() int {
	return types.POS_RESTING
}

// Level returns the minimum level required to execute the command
func (c *SayCommand) Level() int {
	return 0
}

// LogCommand returns whether the command should be logged
func (c *SayCommand) LogCommand() bool {
	return false
}
