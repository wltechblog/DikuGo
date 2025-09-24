package command

import (
	"github.com/wltechblog/DikuGo/pkg/types"
	"github.com/wltechblog/DikuGo/pkg/world"
)

// QuaffCommand represents the quaff command
type QuaffCommand struct {
	World *world.World
}

// Execute executes the quaff command
func (c *QuaffCommand) Execute(character *types.Character, args string) error {
	return DoQuaff(character, args, c.World)
}

// Name returns the name of the command
func (c *QuaffCommand) Name() string {
	return "quaff"
}

// Aliases returns the aliases of the command
func (c *QuaffCommand) Aliases() []string {
	return []string{"drink"}
}

// MinPosition returns the minimum position required to execute the command
func (c *QuaffCommand) MinPosition() int {
	return types.POS_RESTING
}

// Level returns the minimum level required to execute the command
func (c *QuaffCommand) Level() int {
	return 1
}

// LogCommand returns whether the command should be logged
func (c *QuaffCommand) LogCommand() bool {
	return true
}
