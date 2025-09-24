package command

import (
	"fmt"
	"log"
	"strings"

	"github.com/wltechblog/DikuGo/pkg/types"
)

// MovementCommand represents a movement command
type MovementCommand struct {
	direction int
}

// Execute executes the movement command
func (c *MovementCommand) Execute(character *types.Character, args string) error {
	// Check if the character is in a room
	if character.InRoom == nil {
		return fmt.Errorf("you are not in a room")
	}

	// Check if there is an exit in the specified direction
	log.Printf("Movement: Character %s in room %d trying to move in direction %d", character.Name, character.InRoom.VNUM, c.direction)
	log.Printf("Movement: Room %d has exits: %v", character.InRoom.VNUM, character.InRoom.Exits)
	exit := character.InRoom.Exits[c.direction]
	if exit == nil {
		log.Printf("Movement: No exit in direction %d", c.direction)
		return fmt.Errorf("you cannot go that way")
	}

	// Check if the exit leads to a room
	log.Printf("Movement: Exit in direction %d, DestVnum: %d", c.direction, exit.DestVnum)
	if exit.DestVnum == -1 {
		log.Printf("Movement: Exit in direction %d has no destination room", c.direction)
		return fmt.Errorf("you cannot go that way")
	}

	// Get the destination room from the world
	var destRoom *types.Room
	if world, ok := character.World.(interface{ GetRoom(int) *types.Room }); ok {
		destRoom = world.GetRoom(exit.DestVnum)
		if destRoom == nil {
			log.Printf("Movement: Could not find destination room %d", exit.DestVnum)
			return fmt.Errorf("you cannot go that way")
		}
	} else {
		log.Printf("Movement: Character %s has no World field", character.Name)
		return fmt.Errorf("you cannot go that way")
	}

	// Check if the exit is closed
	if exit.IsClosed() {
		if exit.Keywords != "" {
			return fmt.Errorf("the %s is closed", exit.Keywords)
		} else {
			return fmt.Errorf("the door is closed")
		}
	}

	// Move the character to the new room using the world's proper movement method
	if worldInterface, ok := character.World.(interface {
		CharacterMove(*types.Character, *types.Room)
	}); ok {
		worldInterface.CharacterMove(character, destRoom)
	} else {
		log.Printf("Movement: Character %s has no World interface with CharacterMove method", character.Name)
		return fmt.Errorf("you cannot go that way")
	}

	// Show the new room to the character
	var sb strings.Builder

	// Room name
	sb.WriteString(fmt.Sprintf("\r\n%s\r\n", destRoom.Name))

	// Room description
	sb.WriteString(fmt.Sprintf("%s\r\n", destRoom.Description))

	// Exits
	sb.WriteString("\r\nExits: ")
	var exits []string
	for dir := 0; dir < 6; dir++ {
		e := destRoom.Exits[dir]
		if e != nil && e.DestVnum != -1 {
			exits = append(exits, directionName(dir))
		}
	}
	if len(exits) > 0 {
		sb.WriteString(strings.Join(exits, ", "))
	} else {
		sb.WriteString("none")
	}
	sb.WriteString("\r\n")

	// Characters in the room
	for _, ch := range destRoom.Characters {
		if ch != character {
			sb.WriteString(fmt.Sprintf("%s is here.\r\n", ch.ShortDesc))
		}
	}

	// Objects in the room
	for _, obj := range destRoom.Objects {
		sb.WriteString(fmt.Sprintf("%s is here.\r\n", obj.Prototype.ShortDesc))
	}

	// Send the description to the character
	return fmt.Errorf("%s", sb.String())
}

// Name returns the name of the command
func (c *MovementCommand) Name() string {
	return directionName(c.direction)
}

// Aliases returns the aliases of the command
func (c *MovementCommand) Aliases() []string {
	switch c.direction {
	case types.DIR_NORTH:
		return []string{"n"}
	case types.DIR_EAST:
		return []string{"e"}
	case types.DIR_SOUTH:
		return []string{"s"}
	case types.DIR_WEST:
		return []string{"w"}
	case types.DIR_UP:
		return []string{"u"}
	case types.DIR_DOWN:
		return []string{"d"}
	default:
		return []string{}
	}
}

// MinPosition returns the minimum position required to execute the command
func (c *MovementCommand) MinPosition() int {
	return types.POS_STANDING
}

// Level returns the minimum level required to execute the command
func (c *MovementCommand) Level() int {
	return 0
}

// LogCommand returns whether the command should be logged
func (c *MovementCommand) LogCommand() bool {
	return false
}
