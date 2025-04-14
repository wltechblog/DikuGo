package command

import (
	"fmt"
	"strings"

	"github.com/wltechblog/DikuGo/pkg/types"
)

// ValidateRoomsCommand represents the validaterooms command
type ValidateRoomsCommand struct{}

// Execute executes the validaterooms command
func (c *ValidateRoomsCommand) Execute(character *types.Character, args string) error {
	// Check if the character has access to the world
	world, ok := character.World.(interface{ GetRooms() []*types.Room })
	if !ok {
		return fmt.Errorf("world interface not available")
	}

	// Get all rooms
	rooms := world.GetRooms()

	// Validate all room exits
	var sb strings.Builder
	sb.WriteString("Validating all room exits...\r\n")

	var totalExits, validExits, invalidExits int

	for _, room := range rooms {
		for dir := 0; dir < 6; dir++ {
			exit := room.Exits[dir]
			if exit == nil {
				continue
			}

			totalExits++

			// Skip exits that don't lead anywhere
			if exit.DestVnum == -1 {
				continue
			}

			// Check if the destination room exists
			destRoom := findRoomByVnum(rooms, exit.DestVnum)
			if destRoom == nil {
				invalidExits++
				sb.WriteString(fmt.Sprintf("ERROR: Room %d has exit in direction %s with DestVnum %d but that room does not exist\r\n",
					room.VNUM, directionName(dir), exit.DestVnum))
			} else {
				validExits++
			}
		}
	}

	// Add summary
	sb.WriteString(fmt.Sprintf("\r\nValidation complete.\r\n"))
	sb.WriteString(fmt.Sprintf("Total exits: %d\r\n", totalExits))
	sb.WriteString(fmt.Sprintf("Valid exits: %d\r\n", validExits))
	sb.WriteString(fmt.Sprintf("Invalid exits: %d\r\n", invalidExits))

	// Send the results to the character
	return fmt.Errorf("%s", sb.String())
}

// Name returns the name of the command
func (c *ValidateRoomsCommand) Name() string {
	return "validaterooms"
}

// Aliases returns the aliases of the command
func (c *ValidateRoomsCommand) Aliases() []string {
	return []string{"checkrooms"}
}

// MinPosition returns the minimum position required to execute the command
func (c *ValidateRoomsCommand) MinPosition() int {
	return types.POS_STANDING
}

// Level returns the minimum level required to execute the command
func (c *ValidateRoomsCommand) Level() int {
	return 1 // Admin only
}

// LogCommand returns whether the command should be logged
func (c *ValidateRoomsCommand) LogCommand() bool {
	return true
}

// findRoomByVnum finds a room by its VNUM
func findRoomByVnum(rooms []*types.Room, vnum int) *types.Room {
	for _, room := range rooms {
		if room.VNUM == vnum {
			return room
		}
	}
	return nil
}
