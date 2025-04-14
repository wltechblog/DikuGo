package command

import (
	"fmt"
	"strings"

	"github.com/wltechblog/DikuGo/pkg/types"
)

// RemoveCommand represents the remove command
type RemoveCommand struct{}

// Execute executes the remove command
func (c *RemoveCommand) Execute(character *types.Character, args string) error {
	// Check if there are arguments
	if args == "" {
		return fmt.Errorf("remove what?")
	}

	// Check if the target is "all"
	if strings.ToLower(args) == "all" {
		return c.removeAll(character)
	}

	// Find the target object
	var obj *types.ObjectInstance
	var position int
	for i, o := range character.Equipment {
		if o != nil && strings.Contains(strings.ToLower(o.Prototype.Name), strings.ToLower(args)) {
			obj = o
			position = i
			break
		}
	}

	if obj == nil {
		return fmt.Errorf("you're not wearing %s", args)
	}

	// Remove the object from the character's equipment
	character.Equipment[position] = nil

	// Add the object to the character's inventory
	obj.WornBy = nil
	obj.WornOn = -1
	obj.CarriedBy = character
	character.Inventory = append(character.Inventory, obj)

	// Send a message to the character
	return fmt.Errorf("you remove %s.\r\n", obj.Prototype.ShortDesc)
}

// removeAll removes all equipment
func (c *RemoveCommand) removeAll(character *types.Character) error {
	// Check if the character is wearing anything
	var hasEquipment bool
	for _, obj := range character.Equipment {
		if obj != nil {
			hasEquipment = true
			break
		}
	}

	if !hasEquipment {
		return fmt.Errorf("you're not wearing anything.\r\n")
	}

	// Remove each piece of equipment
	var sb strings.Builder
	for i, obj := range character.Equipment {
		if obj == nil {
			continue
		}

		// Remove the object from the character's equipment
		character.Equipment[i] = nil

		// Add the object to the character's inventory
		obj.WornBy = nil
		obj.WornOn = -1
		obj.CarriedBy = character
		character.Inventory = append(character.Inventory, obj)

		// Add a message to the buffer
		sb.WriteString(fmt.Sprintf("You remove %s.\r\n", obj.Prototype.ShortDesc))
	}

	// Send the messages to the character
	return fmt.Errorf("%s", sb.String())
}

// Name returns the name of the command
func (c *RemoveCommand) Name() string {
	return "remove"
}

// Aliases returns the aliases of the command
func (c *RemoveCommand) Aliases() []string {
	return []string{"rem"}
}

// MinPosition returns the minimum position required to execute the command
func (c *RemoveCommand) MinPosition() int {
	return types.POS_RESTING
}

// Level returns the minimum level required to execute the command
func (c *RemoveCommand) Level() int {
	return 0
}

// LogCommand returns whether the command should be logged
func (c *RemoveCommand) LogCommand() bool {
	return false
}
