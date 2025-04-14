package command

import (
	"fmt"
	"strings"

	"github.com/wltechblog/DikuGo/pkg/types"
)

// DropCommand represents the drop command
type DropCommand struct{}

// Execute executes the drop command
func (c *DropCommand) Execute(character *types.Character, args string) error {
	// Check if there are arguments
	if args == "" {
		return fmt.Errorf("drop what?")
	}

	// Check if the character is in a room
	if character.InRoom == nil {
		return fmt.Errorf("you are not in a room")
	}

	// Check if the target is "all"
	if strings.ToLower(args) == "all" {
		return c.dropAll(character)
	}

	// Find the target object
	obj := findObjectInInventory(character, args)
	if obj == nil {
		return fmt.Errorf("you don't have %s", args)
	}

	// Remove the object from the character's inventory
	for i, o := range character.Inventory {
		if o == obj {
			character.Inventory = append(character.Inventory[:i], character.Inventory[i+1:]...)
			break
		}
	}

	// Add the object to the room
	obj.CarriedBy = nil
	obj.InRoom = character.InRoom
	character.InRoom.Objects = append(character.InRoom.Objects, obj)

	// Send a message to the character
	return fmt.Errorf("you drop %s.\r\n", obj.Prototype.ShortDesc)
}

// dropAll drops all objects in the character's inventory
func (c *DropCommand) dropAll(character *types.Character) error {
	// Check if the character is in a room
	if character.InRoom == nil {
		return fmt.Errorf("you are not in a room")
	}

	// Check if the character has any objects
	if len(character.Inventory) == 0 {
		return fmt.Errorf("you are not carrying anything")
	}

	// Drop each object
	var sb strings.Builder
	for _, obj := range character.Inventory {
		// Remove the object from the character's inventory
		for i, o := range character.Inventory {
			if o == obj {
				character.Inventory = append(character.Inventory[:i], character.Inventory[i+1:]...)
				break
			}
		}

		// Add the object to the room
		obj.CarriedBy = nil
		obj.InRoom = character.InRoom
		character.InRoom.Objects = append(character.InRoom.Objects, obj)

		// Add a message to the buffer
		sb.WriteString(fmt.Sprintf("You drop %s.\r\n", obj.Prototype.ShortDesc))
	}

	// Send the messages to the character
	return fmt.Errorf("%s", sb.String())
}

// Name returns the name of the command
func (c *DropCommand) Name() string {
	return "drop"
}

// Aliases returns the aliases of the command
func (c *DropCommand) Aliases() []string {
	return []string{}
}

// MinPosition returns the minimum position required to execute the command
func (c *DropCommand) MinPosition() int {
	return types.POS_STANDING
}

// Level returns the minimum level required to execute the command
func (c *DropCommand) Level() int {
	return 0
}

// LogCommand returns whether the command should be logged
func (c *DropCommand) LogCommand() bool {
	return false
}
