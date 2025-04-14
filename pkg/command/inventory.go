package command

import (
	"fmt"
	"strings"

	"github.com/wltechblog/DikuGo/pkg/types"
)

// InventoryCommand represents the inventory command
type InventoryCommand struct{}

// Execute executes the inventory command
func (c *InventoryCommand) Execute(character *types.Character, args string) error {
	// Check if the character has any objects
	if len(character.Inventory) == 0 {
		return fmt.Errorf("you are not carrying anything.\r\n")
	}

	// Build the inventory list
	var sb strings.Builder
	sb.WriteString("You are carrying:\r\n")
	for _, obj := range character.Inventory {
		sb.WriteString(fmt.Sprintf("  %s\r\n", obj.Prototype.ShortDesc))
	}

	// Send the inventory list to the character
	return fmt.Errorf("%s", sb.String())
}

// Name returns the name of the command
func (c *InventoryCommand) Name() string {
	return "inventory"
}

// Aliases returns the aliases of the command
func (c *InventoryCommand) Aliases() []string {
	return []string{"i"}
}

// MinPosition returns the minimum position required to execute the command
func (c *InventoryCommand) MinPosition() int {
	return types.POS_RESTING
}

// Level returns the minimum level required to execute the command
func (c *InventoryCommand) Level() int {
	return 0
}

// LogCommand returns whether the command should be logged
func (c *InventoryCommand) LogCommand() bool {
	return false
}
