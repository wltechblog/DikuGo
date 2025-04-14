package command

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/wltechblog/DikuGo/pkg/types"
)

// ResetZoneCommand is a command that forces a zone reset
type ResetZoneCommand struct{}

// Name returns the name of the command
func (c *ResetZoneCommand) Name() string {
	return "resetzones"
}

// Aliases returns the aliases of the command
func (c *ResetZoneCommand) Aliases() []string {
	return []string{"resetzone"}
}

// MinPosition returns the minimum position required to execute the command
func (c *ResetZoneCommand) MinPosition() int {
	return types.POS_STANDING
}

// Level returns the minimum level required to execute the command
func (c *ResetZoneCommand) Level() int {
	return 1 // Available to everyone during development
}

// LogCommand returns whether the command should be logged
func (c *ResetZoneCommand) LogCommand() bool {
	return true
}

// Execute executes the reset zone command
func (c *ResetZoneCommand) Execute(character *types.Character, args string) error {
	// Get the world from the character
	world, ok := character.World.(interface {
		ResetZones()
		GetZone(int) *types.Zone
	})
	if !ok {
		return fmt.Errorf("world interface not available")
	}

	// Parse the zone VNUM from args if provided
	vnumStr := strings.TrimSpace(args)
	if vnumStr != "" {
		// Parse the VNUM
		vnum, err := strconv.Atoi(vnumStr)
		if err != nil {
			return fmt.Errorf("invalid zone vnum: %s", vnumStr)
		}

		// Get the zone by VNUM
		zone := world.GetZone(vnum)
		if zone == nil {
			return fmt.Errorf("zone %d does not exist", vnum)
		}

		// Force the zone to be ready for reset
		zone.Age = zone.Lifespan
	}

	// Reset all zones
	world.ResetZones()

	character.SendMessage("Zone reset complete.")
	return nil
}
