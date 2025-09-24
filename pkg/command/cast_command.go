package command

import (
	"github.com/wltechblog/DikuGo/pkg/types"
	"github.com/wltechblog/DikuGo/pkg/world"
)

// CastCommand represents the cast command
type CastCommand struct {
	World *world.World
}

// Execute executes the cast command
func (c *CastCommand) Execute(character *types.Character, args string) error {
	return DoCast(character, args, c.World)
}

// Name returns the name of the command
func (c *CastCommand) Name() string {
	return "cast"
}

// Aliases returns the aliases of the command
func (c *CastCommand) Aliases() []string {
	return []string{"c"}
}

// MinPosition returns the minimum position required to execute the command
func (c *CastCommand) MinPosition() int {
	return types.POS_STANDING
}

// Level returns the minimum level required to execute the command
func (c *CastCommand) Level() int {
	return 1
}

// LogCommand returns whether the command should be logged
func (c *CastCommand) LogCommand() bool {
	return true
}
