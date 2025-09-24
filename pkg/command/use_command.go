package command

import (
	"github.com/wltechblog/DikuGo/pkg/types"
	"github.com/wltechblog/DikuGo/pkg/world"
)

// UseCommand represents the use command
type UseCommand struct {
	World *world.World
}

// Execute executes the use command
func (c *UseCommand) Execute(character *types.Character, args string) error {
	return DoUse(character, args, c.World)
}

// Name returns the name of the command
func (c *UseCommand) Name() string {
	return "use"
}

// Aliases returns the aliases of the command
func (c *UseCommand) Aliases() []string {
	return []string{}
}

// MinPosition returns the minimum position required to execute the command
func (c *UseCommand) MinPosition() int {
	return types.POS_STANDING
}

// Level returns the minimum level required to execute the command
func (c *UseCommand) Level() int {
	return 1
}

// LogCommand returns whether the command should be logged
func (c *UseCommand) LogCommand() bool {
	return true
}
