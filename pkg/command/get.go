package command

import (
	"fmt"
	"strings"

	"github.com/wltechblog/DikuGo/pkg/types"
)

// GetCommand represents the get command
type GetCommand struct{}

// Execute executes the get command
func (c *GetCommand) Execute(character *types.Character, args string) error {
	// Check if there are arguments
	if args == "" {
		return fmt.Errorf("get what?")
	}

	// Check if the character is in a room
	if character.InRoom == nil {
		return fmt.Errorf("you are not in a room")
	}

	// Parse the arguments
	var target, container string
	parts := strings.SplitN(args, " from ", 2)
	if len(parts) > 1 {
		target = parts[0]
		container = parts[1]
	} else {
		target = args
	}

	// Check if the target is "all"
	if strings.ToLower(target) == "all" {
		return c.getAll(character, container)
	}

	// Find the target object
	var obj *types.ObjectInstance
	if container != "" {
		// Get from container
		containerObj := findObjectInRoom(character.InRoom, container)
		if containerObj == nil {
			containerObj = findObjectInInventory(character, container)
		}
		if containerObj == nil {
			return fmt.Errorf("you don't see %s here", container)
		}
		obj = findObjectInContainer(containerObj, target)
		if obj == nil {
			return fmt.Errorf("you don't see %s in %s", target, containerObj.Prototype.ShortDesc)
		}
	} else {
		// Get from room
		obj = findObjectInRoom(character.InRoom, target)
		if obj == nil {
			return fmt.Errorf("you don't see %s here", target)
		}
	}

	// Check if the object can be picked up
	if obj.Prototype.ExtraFlags&types.ITEM_NOPICK != 0 {
		return fmt.Errorf("you can't pick up %s", obj.Prototype.ShortDesc)
	}

	// Check if the character can carry more weight
	// TODO: Implement weight checks

	// Remove the object from its current location
	if container != "" {
		// Remove from container
		containerObj := findObjectInRoom(character.InRoom, container)
		if containerObj == nil {
			containerObj = findObjectInInventory(character, container)
		}
		removeObjectFromContainer(containerObj, obj)
	} else {
		// Remove from room
		removeObjectFromRoom(character.InRoom, obj)
	}

	// Add the object to the character's inventory
	obj.InRoom = nil
	obj.CarriedBy = character
	character.Inventory = append(character.Inventory, obj)

	// Send a message to the character
	return fmt.Errorf("you get %s.\r\n", obj.Prototype.ShortDesc)
}

// getAll gets all objects from a container or room
func (c *GetCommand) getAll(character *types.Character, container string) error {
	// Check if the character is in a room
	if character.InRoom == nil {
		return fmt.Errorf("you are not in a room")
	}

	// Get objects from container or room
	var objects []*types.ObjectInstance
	if container != "" {
		// Get from container
		containerObj := findObjectInRoom(character.InRoom, container)
		if containerObj == nil {
			containerObj = findObjectInInventory(character, container)
		}
		if containerObj == nil {
			return fmt.Errorf("you don't see %s here", container)
		}
		objects = containerObj.Contains
	} else {
		// Get from room
		objects = character.InRoom.Objects
	}

	// Check if there are any objects
	if len(objects) == 0 {
		if container != "" {
			containerObj := findObjectInRoom(character.InRoom, container)
			if containerObj == nil {
				containerObj = findObjectInInventory(character, container)
			}
			return fmt.Errorf("%s is empty", containerObj.Prototype.ShortDesc)
		}
		return fmt.Errorf("there is nothing here")
	}

	// Get each object
	var sb strings.Builder
	for _, obj := range objects {
		// Check if the object can be picked up
		if obj.Prototype.ExtraFlags&types.ITEM_NOPICK != 0 {
			sb.WriteString(fmt.Sprintf("You can't pick up %s.\r\n", obj.Prototype.ShortDesc))
			continue
		}

		// Check if the character can carry more weight
		// TODO: Implement weight checks

		// Remove the object from its current location
		if container != "" {
			// Remove from container
			containerObj := findObjectInRoom(character.InRoom, container)
			if containerObj == nil {
				containerObj = findObjectInInventory(character, container)
			}
			removeObjectFromContainer(containerObj, obj)
		} else {
			// Remove from room
			removeObjectFromRoom(character.InRoom, obj)
		}

		// Add the object to the character's inventory
		obj.InRoom = nil
		obj.CarriedBy = character
		character.Inventory = append(character.Inventory, obj)

		// Add a message to the buffer
		sb.WriteString(fmt.Sprintf("You get %s.\r\n", obj.Prototype.ShortDesc))
	}

	// Send the messages to the character
	return fmt.Errorf("%s", sb.String())
}

// Name returns the name of the command
func (c *GetCommand) Name() string {
	return "get"
}

// Aliases returns the aliases of the command
func (c *GetCommand) Aliases() []string {
	return []string{"take"}
}

// MinPosition returns the minimum position required to execute the command
func (c *GetCommand) MinPosition() int {
	return types.POS_STANDING
}

// Level returns the minimum level required to execute the command
func (c *GetCommand) Level() int {
	return 0
}

// LogCommand returns whether the command should be logged
func (c *GetCommand) LogCommand() bool {
	return false
}

// findObjectInRoom finds an object in a room by name
func findObjectInRoom(room *types.Room, name string) *types.ObjectInstance {
	// Convert the name to lowercase
	name = strings.ToLower(name)

	// Check each object in the room
	for _, obj := range room.Objects {
		// Check if the object's name contains the search name
		if strings.Contains(strings.ToLower(obj.Prototype.Name), name) {
			return obj
		}

		// Check if the object's short description contains the search name
		if strings.Contains(strings.ToLower(obj.Prototype.ShortDesc), name) {
			return obj
		}
	}

	return nil
}

// findObjectInInventory finds an object in a character's inventory by name
func findObjectInInventory(character *types.Character, name string) *types.ObjectInstance {
	// Convert the name to lowercase
	name = strings.ToLower(name)

	// Check each object in the inventory
	for _, obj := range character.Inventory {
		// Check if the object's name contains the search name
		if strings.Contains(strings.ToLower(obj.Prototype.Name), name) {
			return obj
		}

		// Check if the object's short description contains the search name
		if strings.Contains(strings.ToLower(obj.Prototype.ShortDesc), name) {
			return obj
		}
	}

	return nil
}

// findObjectInContainer finds an object in a container by name
func findObjectInContainer(container *types.ObjectInstance, name string) *types.ObjectInstance {
	// Convert the name to lowercase
	name = strings.ToLower(name)

	// Check each object in the container
	for _, obj := range container.Contains {
		// Check if the object's name contains the search name
		if strings.Contains(strings.ToLower(obj.Prototype.Name), name) {
			return obj
		}

		// Check if the object's short description contains the search name
		if strings.Contains(strings.ToLower(obj.Prototype.ShortDesc), name) {
			return obj
		}
	}

	return nil
}

// removeObjectFromRoom removes an object from a room
func removeObjectFromRoom(room *types.Room, obj *types.ObjectInstance) {
	for i, o := range room.Objects {
		if o == obj {
			room.Objects = append(room.Objects[:i], room.Objects[i+1:]...)
			break
		}
	}
}

// removeObjectFromContainer removes an object from a container
func removeObjectFromContainer(container *types.ObjectInstance, obj *types.ObjectInstance) {
	for i, o := range container.Contains {
		if o == obj {
			container.Contains = append(container.Contains[:i], container.Contains[i+1:]...)
			break
		}
	}
}
