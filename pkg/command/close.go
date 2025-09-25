package command

import (
	"fmt"
	"strings"

	"github.com/wltechblog/DikuGo/pkg/types"
)

// CloseCommand represents the close command
type CloseCommand struct{}

// Execute executes the close command
func (c *CloseCommand) Execute(character *types.Character, args string) error {
	// Check if there are arguments
	if args == "" {
		return fmt.Errorf("close what?")
	}

	// Check if the character is in a room
	if character.InRoom == nil {
		return fmt.Errorf("you are not in a room")
	}

	// Parse arguments - could be "door", "chest", "door north", etc.
	parts := strings.Fields(strings.ToLower(strings.TrimSpace(args)))
	if len(parts) == 0 {
		return fmt.Errorf("close what?")
	}

	objectName := parts[0]
	var direction string
	if len(parts) > 1 {
		direction = parts[1]
	}

	// First try to find an object (container) in inventory or room
	targetObject := c.findObject(character, objectName)
	if targetObject != nil {
		return c.closeContainer(character, targetObject)
	}

	// If no object found, try to find a door/exit
	exitDir := c.findDoor(character, objectName, direction)
	if exitDir >= 0 {
		return c.closeDoor(character, exitDir)
	}

	// Nothing found
	return fmt.Errorf("I see no %s here.", objectName)
}

// findObject looks for an object in inventory and room
func (c *CloseCommand) findObject(character *types.Character, name string) *types.ObjectInstance {
	// Check inventory first
	for _, obj := range character.Inventory {
		if c.matchesName(obj, name) {
			return obj
		}
	}

	// Check room objects
	for _, obj := range character.InRoom.Objects {
		if c.matchesName(obj, name) {
			return obj
		}
	}

	return nil
}

// matchesName checks if an object matches the given name
func (c *CloseCommand) matchesName(obj *types.ObjectInstance, name string) bool {
	return strings.Contains(strings.ToLower(obj.Prototype.Name), name) ||
		strings.Contains(strings.ToLower(obj.Prototype.ShortDesc), name)
}

// findDoor looks for a door/exit by name and optional direction
func (c *CloseCommand) findDoor(character *types.Character, doorName string, direction string) int {
	// If a direction was specified, check that specific exit
	if direction != "" {
		exitDir := c.parseDirection(direction)
		if exitDir >= 0 && exitDir < 6 {
			exit := character.InRoom.Exits[exitDir]
			if exit != nil && c.matchesDoorName(exit, doorName) {
				return exitDir
			}
		}
		return -1
	}

	// No direction specified, search all exits for matching keyword
	for dir := 0; dir < 6; dir++ {
		exit := character.InRoom.Exits[dir]
		if exit != nil && c.matchesDoorName(exit, doorName) {
			return dir
		}
	}

	return -1
}

// matchesDoorName checks if an exit matches the given door name
func (c *CloseCommand) matchesDoorName(exit *types.Exit, name string) bool {
	if exit.Keywords == "" {
		// If no keywords, match common door names
		return name == "door"
	}
	return strings.Contains(strings.ToLower(exit.Keywords), name)
}

// parseDirection converts direction string to direction constant
func (c *CloseCommand) parseDirection(dir string) int {
	switch strings.ToLower(dir) {
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

// closeContainer handles closing a container object
func (c *CloseCommand) closeContainer(character *types.Character, container *types.ObjectInstance) error {
	// Check if it's actually a container
	if container.Prototype.Type != types.ITEM_CONTAINER {
		return fmt.Errorf("that's not a container.")
	}

	// Check if it's already closed
	if (container.Prototype.Value[1] & types.CONT_CLOSED) != 0 {
		return fmt.Errorf("but it's already closed!")
	}

	// Check if it can be closed
	if (container.Prototype.Value[1] & types.CONT_CLOSEABLE) == 0 {
		return fmt.Errorf("that's impossible.")
	}

	// Close the container by setting the CLOSED flag
	container.Value[1] = container.Prototype.Value[1] | types.CONT_CLOSED

	// Send messages
	character.SendMessage("Ok.\r\n")

	// Send message to room
	for _, ch := range character.InRoom.Characters {
		if ch != character {
			ch.SendMessage(fmt.Sprintf("%s closes %s.\r\n", character.Name, container.Prototype.ShortDesc))
		}
	}

	return nil
}

// closeDoor handles closing a door/exit
func (c *CloseCommand) closeDoor(character *types.Character, direction int) error {
	exit := character.InRoom.Exits[direction]
	if exit == nil {
		return fmt.Errorf("there is no exit in that direction.")
	}

	// Check if it's actually a door
	if (exit.Flags & types.EX_ISDOOR) == 0 {
		return fmt.Errorf("that's absurd.")
	}

	// Check if it's already closed
	if (exit.Flags & types.EX_CLOSED) != 0 {
		return fmt.Errorf("it's already closed!")
	}

	// Close the door by setting the CLOSED flag
	exit.Flags |= types.EX_CLOSED

	// Send messages
	character.SendMessage("Ok.\r\n")

	// Send message to room
	doorName := "door"
	if exit.Keywords != "" {
		doorName = exit.Keywords
	}

	for _, ch := range character.InRoom.Characters {
		if ch != character {
			ch.SendMessage(fmt.Sprintf("%s closes the %s.\r\n", character.Name, doorName))
		}
	}

	// Close the other side of the door if it exists
	c.closeOtherSide(character, direction, exit)

	return nil
}

// closeOtherSide closes the corresponding exit on the other side of the door
func (c *CloseCommand) closeOtherSide(character *types.Character, direction int, exit *types.Exit) {
	// Get the world interface to find the destination room
	world, ok := character.World.(interface {
		GetRoom(vnum int) *types.Room
	})
	if !ok {
		return
	}

	// Get the destination room
	destRoom := world.GetRoom(exit.DestVnum)
	if destRoom == nil {
		return
	}

	// Calculate reverse direction
	reverseDir := c.getReverseDirection(direction)
	if reverseDir < 0 || reverseDir >= 6 {
		return
	}

	// Get the reverse exit
	reverseExit := destRoom.Exits[reverseDir]
	if reverseExit == nil {
		return
	}

	// Check if the reverse exit leads back to this room
	if reverseExit.DestVnum != character.InRoom.VNUM {
		return
	}

	// Close the reverse exit
	reverseExit.Flags |= types.EX_CLOSED

	// Send message to the other room
	doorName := "door"
	if reverseExit.Keywords != "" {
		doorName = reverseExit.Keywords
	}

	for _, ch := range destRoom.Characters {
		ch.SendMessage(fmt.Sprintf("The %s closes quietly.\r\n", doorName))
	}
}

// getReverseDirection returns the opposite direction
func (c *CloseCommand) getReverseDirection(dir int) int {
	switch dir {
	case types.DIR_NORTH:
		return types.DIR_SOUTH
	case types.DIR_EAST:
		return types.DIR_WEST
	case types.DIR_SOUTH:
		return types.DIR_NORTH
	case types.DIR_WEST:
		return types.DIR_EAST
	case types.DIR_UP:
		return types.DIR_DOWN
	case types.DIR_DOWN:
		return types.DIR_UP
	default:
		return -1
	}
}

// Name returns the command name
func (c *CloseCommand) Name() string {
	return "close"
}

// Aliases returns command aliases
func (c *CloseCommand) Aliases() []string {
	return []string{}
}

// MinPosition returns the minimum position required
func (c *CloseCommand) MinPosition() int {
	return types.POS_SITTING
}

// Level returns the minimum level required
func (c *CloseCommand) Level() int {
	return 0
}

// LogCommand returns whether this command should be logged
func (c *CloseCommand) LogCommand() bool {
	return false
}
