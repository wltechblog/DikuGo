package command

import (
	"github.com/wltechblog/DikuGo/pkg/types"
	"github.com/wltechblog/DikuGo/pkg/world"
)

// ReciteCommand represents the recite command
type ReciteCommand struct {
	World *world.World
}

// Execute executes the recite command
func (c *ReciteCommand) Execute(character *types.Character, args string) error {
	return DoRecite(character, args, c.World)
}

// Name returns the name of the command
func (c *ReciteCommand) Name() string {
	return "recite"
}

// Aliases returns the aliases of the command
func (c *ReciteCommand) Aliases() []string {
	return []string{"read"}
}

// MinPosition returns the minimum position required to execute the command
func (c *ReciteCommand) MinPosition() int {
	return types.POS_RESTING
}

// Level returns the minimum level required to execute the command
func (c *ReciteCommand) Level() int {
	return 1
}

// LogCommand returns whether the command should be logged
func (c *ReciteCommand) LogCommand() bool {
	return true
}
