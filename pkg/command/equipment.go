package command

import (
	"fmt"
	"strings"

	"github.com/wltechblog/DikuGo/pkg/types"
)

// EquipmentCommand represents the equipment command
type EquipmentCommand struct{}

// Execute executes the equipment command
func (c *EquipmentCommand) Execute(character *types.Character, args string) error {
	// Check if the character has any equipment
	var hasEquipment bool
	for _, obj := range character.Equipment {
		if obj != nil {
			hasEquipment = true
			break
		}
	}

	if !hasEquipment {
		return fmt.Errorf("you are not using any equipment.\r\n")
	}

	// Build the equipment list
	var sb strings.Builder
	sb.WriteString("You are using:\r\n")
	for i, obj := range character.Equipment {
		if obj != nil {
			sb.WriteString(fmt.Sprintf("  %-20s %s\r\n", wearPositionName(i)+":", obj.Prototype.ShortDesc))
		}
	}

	// Send the equipment list to the character
	return fmt.Errorf("%s", sb.String())
}

// Name returns the name of the command
func (c *EquipmentCommand) Name() string {
	return "equipment"
}

// Aliases returns the aliases of the command
func (c *EquipmentCommand) Aliases() []string {
	return []string{"eq"}
}

// MinPosition returns the minimum position required to execute the command
func (c *EquipmentCommand) MinPosition() int {
	return types.POS_SLEEPING
}

// Level returns the minimum level required to execute the command
func (c *EquipmentCommand) Level() int {
	return 0
}

// LogCommand returns whether the command should be logged
func (c *EquipmentCommand) LogCommand() bool {
	return false
}
