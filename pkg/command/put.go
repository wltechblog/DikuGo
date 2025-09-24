package command

import (
	"fmt"
	"strings"

	"github.com/wltechblog/DikuGo/pkg/types"
)

// PutCommand represents the put command
type PutCommand struct{}

// Name returns the name of the command
func (c *PutCommand) Name() string {
	return "put"
}

// Aliases returns the aliases of the command
func (c *PutCommand) Aliases() []string {
	return []string{}
}

// MinPosition returns the minimum position required to execute the command
func (c *PutCommand) MinPosition() int {
	return types.POS_STANDING
}

// Level returns the minimum level required to execute the command
func (c *PutCommand) Level() int {
	return 0
}

// LogCommand returns whether the command should be logged
func (c *PutCommand) LogCommand() bool {
	return false
}

// Execute executes the put command
func (c *PutCommand) Execute(character *types.Character, args string) error {
	// Check if there are arguments
	if args == "" {
		return fmt.Errorf("put what in what?")
	}

	// Check if the character is in a room
	if character.InRoom == nil {
		return fmt.Errorf("you are not in a room")
	}

	// Parse the arguments
	var target, container string
	parts := strings.SplitN(args, " in ", 2)
	if len(parts) < 2 {
		return fmt.Errorf("put what in what?")
	}
	target = parts[0]
	container = parts[1]

	// Check if the target is "all"
	if strings.ToLower(target) == "all" {
		return c.putAll(character, container)
	}

	// Find the target object in inventory
	obj := findObjectInInventory(character, target)
	if obj == nil {
		return fmt.Errorf("you don't have %s", target)
	}

	// Find the container
	containerObj := findObjectInInventory(character, container)
	if containerObj == nil {
		containerObj = findObjectInRoom(character.InRoom, container)
	}
	if containerObj == nil {
		return fmt.Errorf("you don't see %s here", container)
	}

	// Check if the container is a container
	if containerObj.Prototype.Type != types.ITEM_CONTAINER {
		return fmt.Errorf("%s is not a container", containerObj.Prototype.ShortDesc)
	}

	// Check if the container is closed
	if containerObj.Prototype.Value[1]&types.CONT_CLOSED != 0 {
		return fmt.Errorf("%s is closed", containerObj.Prototype.ShortDesc)
	}

	// Check if the object can fit in the container
	// TODO: Implement weight checks
	maxWeight := containerObj.Prototype.Value[0]
	if maxWeight > 0 {
		// Calculate current weight in container
		currentWeight := 0
		for _, item := range containerObj.Contains {
			currentWeight += item.Prototype.Weight
		}

		// Check if adding this item would exceed the weight limit
		if currentWeight+obj.Prototype.Weight > maxWeight {
			return fmt.Errorf("%s is too full to hold %s", containerObj.Prototype.ShortDesc, obj.Prototype.ShortDesc)
		}
	}

	// Check if the object is a container and contains the container (prevent recursive containers)
	if obj.Prototype.Type == types.ITEM_CONTAINER {
		// Check if the container is inside the object
		for _, item := range obj.Contains {
			if item == containerObj {
				return fmt.Errorf("you can't put something inside itself")
			}
		}
	}

	// Remove the object from the character's inventory
	for i, o := range character.Inventory {
		if o == obj {
			character.Inventory = append(character.Inventory[:i], character.Inventory[i+1:]...)
			break
		}
	}

	// Add the object to the container
	obj.CarriedBy = nil
	obj.InObj = containerObj
	containerObj.Contains = append(containerObj.Contains, obj)

	// Send a message to the character
	return fmt.Errorf("you put %s in %s.\r\n", obj.Prototype.ShortDesc, containerObj.Prototype.ShortDesc)
}

// putAll puts all objects from inventory into a container
func (c *PutCommand) putAll(character *types.Character, container string) error {
	// Check if the character is in a room
	if character.InRoom == nil {
		return fmt.Errorf("you are not in a room")
	}

	// Check if the character has any objects
	if len(character.Inventory) == 0 {
		return fmt.Errorf("you are not carrying anything")
	}

	// Find the container
	containerObj := findObjectInInventory(character, container)
	if containerObj == nil {
		containerObj = findObjectInRoom(character.InRoom, container)
	}
	if containerObj == nil {
		return fmt.Errorf("you don't see %s here", container)
	}

	// Check if the container is a container
	if containerObj.Prototype.Type != types.ITEM_CONTAINER {
		return fmt.Errorf("%s is not a container", containerObj.Prototype.ShortDesc)
	}

	// Check if the container is closed
	if containerObj.Prototype.Value[1]&types.CONT_CLOSED != 0 {
		return fmt.Errorf("%s is closed", containerObj.Prototype.ShortDesc)
	}

	// Get the maximum weight the container can hold
	maxWeight := containerObj.Prototype.Value[0]

	// Calculate current weight in container
	currentWeight := 0
	for _, item := range containerObj.Contains {
		currentWeight += item.Prototype.Weight
	}

	// Put each object in the container
	var sb strings.Builder
	var itemsToRemove []*types.ObjectInstance

	for _, obj := range character.Inventory {
		// Skip the container itself
		if obj == containerObj {
			continue
		}

		// Check if the object can fit in the container
		if maxWeight > 0 && currentWeight+obj.Prototype.Weight > maxWeight {
			sb.WriteString(fmt.Sprintf("%s is too full to hold %s.\r\n", containerObj.Prototype.ShortDesc, obj.Prototype.ShortDesc))
			continue
		}

		// Add the object to the list of items to remove from inventory
		itemsToRemove = append(itemsToRemove, obj)

		// Update the current weight
		currentWeight += obj.Prototype.Weight

		// Add a message to the buffer
		sb.WriteString(fmt.Sprintf("You put %s in %s.\r\n", obj.Prototype.ShortDesc, containerObj.Prototype.ShortDesc))
	}

	// Remove items from inventory and add to container
	for _, obj := range itemsToRemove {
		// Remove from inventory
		for i, o := range character.Inventory {
			if o == obj {
				character.Inventory = append(character.Inventory[:i], character.Inventory[i+1:]...)
				break
			}
		}

		// Add to container
		obj.CarriedBy = nil
		obj.InObj = containerObj
		containerObj.Contains = append(containerObj.Contains, obj)
	}

	// If no items were moved, return a message
	if len(itemsToRemove) == 0 {
		return fmt.Errorf("you couldn't put anything in %s", containerObj.Prototype.ShortDesc)
	}

	// Send the messages to the character
	return fmt.Errorf("%s", sb.String())
}
