package command

import (
	"fmt"
	"strings"

	"github.com/wltechblog/DikuGo/pkg/types"
)

// OpenCommand represents the open command
type OpenCommand struct{}

// Execute executes the open command
func (c *OpenCommand) Execute(character *types.Character, args string) error {
	// Check if there are arguments
	if args == "" {
		return fmt.Errorf("open what?")
	}

	// Check if the character is in a room
	if character.InRoom == nil {
		return fmt.Errorf("you are not in a room")
	}

	// First try to find an object (container) using the entire argument string
	// This matches the original DikuMUD behavior of using generic_find with the full argument
	targetObject := c.findObject(character, strings.TrimSpace(strings.ToLower(args)))
	if targetObject != nil {
		return c.openContainer(character, targetObject)
	}

	// Parse arguments using DikuMUD's argument_interpreter logic for door finding
	// This splits the argument into two parts: object/door name and direction
	objectName, direction := c.argumentInterpreter(args)

	// If no object found, try to find a door/exit
	exitDir := c.findDoor(character, objectName, direction)
	if exitDir >= 0 {
		return c.openDoor(character, exitDir)
	}

	// Nothing found
	return fmt.Errorf("I see no %s here.", objectName)
}

// argumentInterpreter splits arguments like the original DikuMUD argument_interpreter
func (c *OpenCommand) argumentInterpreter(args string) (string, string) {
	// Trim and convert to lowercase
	args = strings.TrimSpace(strings.ToLower(args))
	if args == "" {
		return "", ""
	}

	// Split into words
	words := strings.Fields(args)
	if len(words) == 0 {
		return "", ""
	}

	// Find the first non-fill word for the first argument
	firstArg := ""
	firstIndex := 0
	for i, word := range words {
		if !c.isFillWord(word) {
			firstArg = word
			firstIndex = i
			break
		}
	}

	// If we didn't find a non-fill word, use the first word
	if firstArg == "" && len(words) > 0 {
		firstArg = words[0]
		firstIndex = 0
	}

	// Find the second non-fill word for the second argument
	secondArg := ""
	for i := firstIndex + 1; i < len(words); i++ {
		if !c.isFillWord(words[i]) {
			secondArg = words[i]
			break
		}
	}

	return firstArg, secondArg
}

// isFillWord checks if a word is a fill word that should be ignored
func (c *OpenCommand) isFillWord(word string) bool {
	fillWords := []string{"the", "a", "an", "in", "on", "at", "to", "from", "with", "by"}
	for _, fillWord := range fillWords {
		if word == fillWord {
			return true
		}
	}
	return false
}

// findObject looks for an object in inventory and room
func (c *OpenCommand) findObject(character *types.Character, name string) *types.ObjectInstance {
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
func (c *OpenCommand) matchesName(obj *types.ObjectInstance, name string) bool {
	return strings.Contains(strings.ToLower(obj.Prototype.Name), name) ||
		strings.Contains(strings.ToLower(obj.Prototype.ShortDesc), name)
}

// findDoor looks for a door/exit by name and optional direction (following original DikuMUD logic)
func (c *OpenCommand) findDoor(character *types.Character, doorName string, direction string) int {
	// If a direction was specified, check that specific exit
	if direction != "" {
		exitDir := c.parseDirection(direction)
		if exitDir < 0 || exitDir >= 6 {
			// Invalid direction
			return -1
		}

		exit := character.InRoom.Exits[exitDir]
		if exit == nil {
			// No exit in that direction
			return -1
		}

		// If the exit has keywords, check if doorName matches
		if exit.Keywords != "" {
			if c.isName(doorName, exit.Keywords) {
				return exitDir
			}
			// Door exists but name doesn't match
			return -1
		} else {
			// No keywords, so any door name matches this exit
			return exitDir
		}
	}

	// No direction specified, search all exits for matching keyword
	for dir := 0; dir < 6; dir++ {
		exit := character.InRoom.Exits[dir]
		if exit != nil && exit.Keywords != "" {
			if c.isName(doorName, exit.Keywords) {
				return dir
			}
		}
	}

	return -1
}

// isName checks if a name matches any of the keywords (following original DikuMUD isname function)
func (c *OpenCommand) isName(name string, keywords string) bool {
	if keywords == "" {
		return false
	}

	// Split keywords by spaces and check each one
	keywordList := strings.Fields(strings.ToLower(keywords))
	name = strings.ToLower(name)

	for _, keyword := range keywordList {
		// Check for exact match or if name is a prefix of keyword
		if keyword == name || strings.HasPrefix(keyword, name) {
			return true
		}
	}
	return false
}

// matchesDoorName checks if an exit matches the given door name (legacy function, kept for compatibility)
func (c *OpenCommand) matchesDoorName(exit *types.Exit, name string) bool {
	if exit.Keywords == "" {
		// If no keywords, match common door names
		return name == "door"
	}
	return c.isName(name, exit.Keywords)
}

// parseDirection converts direction string to direction constant
func (c *OpenCommand) parseDirection(dir string) int {
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

// openContainer handles opening a container object
func (c *OpenCommand) openContainer(character *types.Character, container *types.ObjectInstance) error {
	// Check if it's actually a container
	if container.Prototype.Type != types.ITEM_CONTAINER {
		return fmt.Errorf("that's not a container.")
	}

	// Check if it's already open
	if (container.Prototype.Value[1] & types.CONT_CLOSED) == 0 {
		return fmt.Errorf("but it's already open!")
	}

	// Check if it can be opened
	if (container.Prototype.Value[1] & types.CONT_CLOSEABLE) == 0 {
		return fmt.Errorf("you can't do that.")
	}

	// Check if it's locked
	if (container.Prototype.Value[1] & types.CONT_LOCKED) != 0 {
		return fmt.Errorf("it seems to be locked.")
	}

	// Open the container by removing the CLOSED flag
	container.Value[1] = container.Prototype.Value[1] &^ types.CONT_CLOSED

	// Send messages
	character.SendMessage("Ok.\r\n")

	// Send message to room
	for _, ch := range character.InRoom.Characters {
		if ch != character {
			ch.SendMessage(fmt.Sprintf("%s opens %s.\r\n", character.Name, container.Prototype.ShortDesc))
		}
	}

	return nil
}

// openDoor handles opening a door/exit
func (c *OpenCommand) openDoor(character *types.Character, direction int) error {
	exit := character.InRoom.Exits[direction]
	if exit == nil {
		return fmt.Errorf("there is no exit in that direction.")
	}

	// Check if it's actually a door
	if (exit.Flags & types.EX_ISDOOR) == 0 {
		return fmt.Errorf("that's impossible, I'm afraid.")
	}

	// Check if it's already open
	if (exit.Flags & types.EX_CLOSED) == 0 {
		return fmt.Errorf("it's already open!")
	}

	// Check if it's locked
	if (exit.Flags & types.EX_LOCKED) != 0 {
		return fmt.Errorf("it seems to be locked.")
	}

	// Open the door by removing the CLOSED flag
	exit.Flags &^= types.EX_CLOSED

	// Send messages
	character.SendMessage("Ok.\r\n")

	// Send message to room
	doorName := "door"
	if exit.Keywords != "" {
		doorName = exit.Keywords
	}

	for _, ch := range character.InRoom.Characters {
		if ch != character {
			ch.SendMessage(fmt.Sprintf("%s opens the %s.\r\n", character.Name, doorName))
		}
	}

	// Open the other side of the door if it exists
	c.openOtherSide(character, direction, exit)

	return nil
}

// openOtherSide opens the corresponding exit on the other side of the door
func (c *OpenCommand) openOtherSide(character *types.Character, direction int, exit *types.Exit) {
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

	// Open the reverse exit
	reverseExit.Flags &^= types.EX_CLOSED

	// Send message to the other room
	doorName := "door"
	if reverseExit.Keywords != "" {
		doorName = reverseExit.Keywords
	}

	for _, ch := range destRoom.Characters {
		ch.SendMessage(fmt.Sprintf("The %s is opened from the other side.\r\n", doorName))
	}
}

// getReverseDirection returns the opposite direction
func (c *OpenCommand) getReverseDirection(dir int) int {
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
func (c *OpenCommand) Name() string {
	return "open"
}

// Aliases returns command aliases
func (c *OpenCommand) Aliases() []string {
	return []string{}
}

// MinPosition returns the minimum position required
func (c *OpenCommand) MinPosition() int {
	return types.POS_SITTING
}

// Level returns the minimum level required
func (c *OpenCommand) Level() int {
	return 0
}

// LogCommand returns whether this command should be logged
func (c *OpenCommand) LogCommand() bool {
	return false
}
