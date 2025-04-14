package command

import (
	"fmt"
	"strings"

	"github.com/wltechblog/DikuGo/pkg/types"
)

// TestExitsCommand represents the testexits command
type TestExitsCommand struct{}

// Execute executes the testexits command
func (c *TestExitsCommand) Execute(character *types.Character, args string) error {
	// Check if the character has access to the world
	world, ok := character.World.(interface{ GetRooms() []*types.Room })
	if !ok {
		return fmt.Errorf("world interface not available")
	}

	// Get all rooms
	rooms := world.GetRooms()

	// Test all room exits
	var sb strings.Builder
	sb.WriteString("Testing all room exits...\r\n")

	var totalExits, validExits, invalidExits, extraDescCount int

	for _, room := range rooms {
		// Count extra descriptions
		extraDescCount += len(room.ExtraDescs)

		// Test exits
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
	sb.WriteString(fmt.Sprintf("\r\nTest complete.\r\n"))
	sb.WriteString(fmt.Sprintf("Total rooms: %d\r\n", len(rooms)))
	sb.WriteString(fmt.Sprintf("Total exits: %d\r\n", totalExits))
	sb.WriteString(fmt.Sprintf("Valid exits: %d\r\n", validExits))
	sb.WriteString(fmt.Sprintf("Invalid exits: %d\r\n", invalidExits))
	sb.WriteString(fmt.Sprintf("Extra descriptions: %d\r\n", extraDescCount))

	// Send the results to the character
	return fmt.Errorf("%s", sb.String())
}

// Name returns the name of the command
func (c *TestExitsCommand) Name() string {
	return "testexits"
}

// Aliases returns the aliases of the command
func (c *TestExitsCommand) Aliases() []string {
	return []string{"testrooms"}
}

// MinPosition returns the minimum position required to execute the command
func (c *TestExitsCommand) MinPosition() int {
	return types.POS_STANDING
}

// Level returns the minimum level required to execute the command
func (c *TestExitsCommand) Level() int {
	return 0 // Available to all during development
}

// LogCommand returns whether the command should be logged
func (c *TestExitsCommand) LogCommand() bool {
	return true
}
