package command

import (
	"fmt"
	"strings"

	"github.com/wltechblog/DikuGo/pkg/types"
)

// ExamineCommand represents the examine command
type ExamineCommand struct{}

// Name returns the name of the command
func (c *ExamineCommand) Name() string {
	return "examine"
}

// Aliases returns the aliases of the command
func (c *ExamineCommand) Aliases() []string {
	return []string{"look in", "exa", "exam"}
}

// MinPosition returns the minimum position required to execute the command
func (c *ExamineCommand) MinPosition() int {
	return types.POS_RESTING
}

// Level returns the minimum level required to execute the command
func (c *ExamineCommand) Level() int {
	return 0
}

// LogCommand returns whether the command should be logged
func (c *ExamineCommand) LogCommand() bool {
	return false
}

// Execute executes the examine command
func (c *ExamineCommand) Execute(character *types.Character, args string) error {
	// Check if there are arguments
	if args == "" {
		return fmt.Errorf("examine what?")
	}

	// Check if the character is in a room
	if character.InRoom == nil {
		return fmt.Errorf("you are not in a room")
	}

	// Find the target object
	var obj *types.ObjectInstance

	// First check inventory
	obj = findObjectInInventory(character, args)
	if obj == nil {
		// Then check room
		obj = findObjectInRoom(character.InRoom, args)
	}
	if obj == nil {
		// Then check equipment
		for _, item := range character.Equipment {
			if item != nil && strings.Contains(strings.ToLower(item.Prototype.Name), strings.ToLower(args)) {
				obj = item
				break
			}
		}
	}

	if obj == nil {
		return fmt.Errorf("you don't see %s here", args)
	}

	// If it's a container, show its contents
	if obj.Prototype.Type == types.ITEM_CONTAINER {
		// Check if the container is closed
		if obj.Prototype.Value[1]&types.CONT_CLOSED != 0 {
			return fmt.Errorf("%s is closed.", obj.Prototype.ShortDesc)
		}

		// Show the contents
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("You look inside %s:\r\n", obj.Prototype.ShortDesc))

		if len(obj.Contains) == 0 {
			sb.WriteString("It is empty.\r\n")
		} else {
			for _, item := range obj.Contains {
				sb.WriteString(fmt.Sprintf("  %s\r\n", item.Prototype.ShortDesc))
			}
		}

		return fmt.Errorf("%s", sb.String())
	}

	// If it's not a container, just examine it
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("You examine %s:\r\n", obj.Prototype.ShortDesc))
	sb.WriteString(fmt.Sprintf("%s\r\n", obj.Prototype.Description))

	// Check for extra descriptions
	for _, extraDesc := range obj.Prototype.ExtraDescs {
		if strings.Contains(strings.ToLower(extraDesc.Keywords), strings.ToLower(args)) {
			sb.WriteString(fmt.Sprintf("%s\r\n", extraDesc.Description))
		}
	}

	return fmt.Errorf("%s", sb.String())
}
