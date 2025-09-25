package command

import (
	"fmt"
	"log"
	"strings"

	"github.com/wltechblog/DikuGo/pkg/types"
)

// LookCommand represents the look command
type LookCommand struct{}

// Execute executes the look command
func (c *LookCommand) Execute(character *types.Character, args string) error {
	// If no arguments, look at the room
	if args == "" {
		return c.lookAtRoom(character)
	}

	// Look at something in the room
	return c.lookAtTarget(character, args)
}

// lookAtRoom looks at the room the character is in
func (c *LookCommand) lookAtRoom(character *types.Character) error {
	// Check if the character is in a room
	if character.InRoom == nil {
		return fmt.Errorf("you are not in a room")
	}

	// Get the room
	room := character.InRoom

	// Lock the room to safely read its contents
	room.RLock()
	defer room.RUnlock()

	// Build the room description
	var sb strings.Builder

	// Room name
	sb.WriteString(fmt.Sprintf("\r\n%s\r\n", room.Name))

	// Room description
	sb.WriteString(fmt.Sprintf("%s\r\n", room.Description))

	// Exits
	sb.WriteString("\r\nExits: ")
	var exits []string
	log.Printf("Look: Room %d has %d exits", room.VNUM, len(room.Exits))
	for dir := 0; dir < 6; dir++ {
		exit := room.Exits[dir]
		if exit != nil {
			log.Printf("Look: Room %d, Exit direction %d, DestVnum: %d, Flags: %d", room.VNUM, dir, exit.DestVnum, exit.Flags)
			// Only show exits that are not closed doors (following original DikuMUD behavior)
			if exit.DestVnum != -1 && !exit.IsClosed() {
				exits = append(exits, directionName(dir))
			}
		} else {
			log.Printf("Look: Room %d, Exit direction %d is nil", room.VNUM, dir)
		}
	}
	if len(exits) > 0 {
		sb.WriteString(strings.Join(exits, ", "))
	} else {
		sb.WriteString("none")
	}
	sb.WriteString("\r\n")

	// Characters in the room (safely read while holding read lock)
	for _, ch := range room.Characters {
		if ch != character {
			if ch.IsNPC {
				sb.WriteString(fmt.Sprintf("%s is here.\r\n", ch.ShortDesc))
			} else {
				sb.WriteString(fmt.Sprintf("%s %s is here.\r\n", ch.Name, ch.Title))
			}
		}
	}

	// Objects in the room (safely read while holding read lock)
	for _, obj := range room.Objects {
		sb.WriteString(fmt.Sprintf("%s is here.\r\n", obj.Prototype.ShortDesc))
	}

	// Send the description to the character
	// TODO: Implement a way to send messages to characters
	// For now, just return the description as an error
	return fmt.Errorf("%s", sb.String())
}

// lookAtTarget looks at a target in the room
func (c *LookCommand) lookAtTarget(character *types.Character, target string) error {
	// Check if the character is in a room
	if character.InRoom == nil {
		return fmt.Errorf("you are not in a room")
	}

	// Get the room
	room := character.InRoom

	// Check if the target is a direction
	dir := getDirectionByName(target)
	if dir >= 0 && dir < 6 {
		return c.lookInDirection(character, dir)
	}

	// Check if the target is "in" something (looking inside a container)
	parts := strings.SplitN(target, " in ", 2)
	if len(parts) == 2 && parts[0] == "" {
		// Looking inside a container
		return c.lookInContainer(character, parts[1])
	}

	// Lock the room to safely read its contents
	room.RLock()
	defer room.RUnlock()

	// Check if the target is a character in the room
	for _, ch := range room.Characters {
		if ch != character && strings.Contains(strings.ToLower(ch.Name), strings.ToLower(target)) {
			return c.lookAtCharacter(character, ch)
		}
	}

	// Check if the target is an object in the room
	for _, obj := range room.Objects {
		if strings.Contains(strings.ToLower(obj.Prototype.Name), strings.ToLower(target)) {
			return c.lookAtObject(character, obj)
		}
	}

	// Check if the target is an object in the character's inventory
	for _, obj := range character.Inventory {
		if strings.Contains(strings.ToLower(obj.Prototype.Name), strings.ToLower(target)) {
			return c.lookAtObject(character, obj)
		}
	}

	// Check if the target is an extra description keyword in the room
	for _, extraDesc := range room.ExtraDescs {
		if strings.Contains(strings.ToLower(extraDesc.Keywords), strings.ToLower(target)) {
			return fmt.Errorf("%s\r\n", extraDesc.Description)
		}
	}

	// Target not found
	return fmt.Errorf("you don't see that here")
}

// Name returns the name of the command
func (c *LookCommand) Name() string {
	return "look"
}

// Aliases returns the aliases of the command
func (c *LookCommand) Aliases() []string {
	return []string{"l"}
}

// MinPosition returns the minimum position required to execute the command
func (c *LookCommand) MinPosition() int {
	return types.POS_RESTING
}

// Level returns the minimum level required to execute the command
func (c *LookCommand) Level() int {
	return 0
}

// LogCommand returns whether the command should be logged
func (c *LookCommand) LogCommand() bool {
	return false
}

// lookInDirection looks in a specific direction
func (c *LookCommand) lookInDirection(character *types.Character, dir int) error {
	// Check if the character is in a room
	if character.InRoom == nil {
		return fmt.Errorf("you are not in a room")
	}

	// Get the room
	room := character.InRoom

	// Check if there is an exit in the specified direction
	exit := room.Exits[dir]
	if exit == nil {
		return fmt.Errorf("you see nothing special in that direction")
	}

	// If the exit has a description, show it
	if exit.Description != "" {
		return fmt.Errorf("%s\r\n", exit.Description)
	}

	// Otherwise, show a default message
	return fmt.Errorf("you see nothing special in that direction\r\n")
}

// lookInContainer looks inside a container
func (c *LookCommand) lookInContainer(character *types.Character, containerName string) error {
	// Check if the character is in a room
	if character.InRoom == nil {
		return fmt.Errorf("you are not in a room")
	}

	// Find the container
	var containerObj *types.ObjectInstance

	// First check inventory
	containerObj = findObjectInInventory(character, containerName)
	if containerObj == nil {
		// Then check room
		containerObj = findObjectInRoom(character.InRoom, containerName)
	}
	if containerObj == nil {
		// Then check equipment
		for _, item := range character.Equipment {
			if item != nil && strings.Contains(strings.ToLower(item.Prototype.Name), strings.ToLower(containerName)) {
				containerObj = item
				break
			}
		}
	}

	if containerObj == nil {
		return fmt.Errorf("you don't see %s here", containerName)
	}

	// Check if it's a container
	if containerObj.Prototype.Type != types.ITEM_CONTAINER {
		return fmt.Errorf("%s is not a container", containerObj.Prototype.ShortDesc)
	}

	// Check if the container is closed
	if containerObj.Prototype.Value[1]&types.CONT_CLOSED != 0 {
		return fmt.Errorf("%s is closed", containerObj.Prototype.ShortDesc)
	}

	// Show the contents
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("You look inside %s:\r\n", containerObj.Prototype.ShortDesc))

	if len(containerObj.Contains) == 0 {
		sb.WriteString("It is empty.\r\n")
	} else {
		for _, item := range containerObj.Contains {
			sb.WriteString(fmt.Sprintf("  %s\r\n", item.Prototype.ShortDesc))
		}
	}

	return fmt.Errorf("%s", sb.String())
}

// lookAtCharacter looks at a character
func (c *LookCommand) lookAtCharacter(character *types.Character, target *types.Character) error {
	// Build the character description
	var sb strings.Builder

	// Character name
	if target.IsNPC {
		sb.WriteString(fmt.Sprintf("%s\r\n", target.ShortDesc))
	} else {
		sb.WriteString(fmt.Sprintf("%s %s\r\n", target.Name, target.Title))
	}

	// Character description
	if target.Description != "" {
		sb.WriteString(fmt.Sprintf("%s\r\n", target.Description))
	}

	// Send the description to the character
	return fmt.Errorf("%s", sb.String())
}

// lookAtObject looks at an object
func (c *LookCommand) lookAtObject(character *types.Character, obj *types.ObjectInstance) error {
	// Build the object description
	var sb strings.Builder

	// Object name
	sb.WriteString(fmt.Sprintf("%s\r\n", obj.Prototype.ShortDesc))

	// Object description
	if obj.Prototype.Description != "" {
		sb.WriteString(fmt.Sprintf("%s\r\n", obj.Prototype.Description))
	}

	// Check for extra descriptions
	for _, extraDesc := range obj.Prototype.ExtraDescs {
		// We don't need to check keywords here since we're looking at the object itself
		sb.WriteString(fmt.Sprintf("%s\r\n", extraDesc.Description))
	}

	// Send the description to the character
	return fmt.Errorf("%s", sb.String())
}

// getDirectionByName returns the direction number for a direction name
func getDirectionByName(name string) int {
	name = strings.ToLower(name)
	switch name {
	case "north", "n":
		return types.DIR_NORTH
	case "east", "e":
		return types.DIR_EAST
	case "south", "s":
		return types.DIR_SOUTH
	case "west", "w":
		return types.DIR_WEST
	case "up", "u":
		return types.DIR_UP
	case "down", "d":
		return types.DIR_DOWN
	default:
		return -1
	}
}
