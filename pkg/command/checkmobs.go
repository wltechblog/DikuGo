package command

import (
	"fmt"
	"sort"
	"strings"

	"github.com/wltechblog/DikuGo/pkg/types"
)

// CheckMobsCommand represents the checkmobs command
type CheckMobsCommand struct{}

// Name returns the name of the command
func (c *CheckMobsCommand) Name() string {
	return "checkmobs"
}

// Aliases returns the aliases of the command
func (c *CheckMobsCommand) Aliases() []string {
	return []string{}
}

// MinPosition returns the minimum position required to execute the command
func (c *CheckMobsCommand) MinPosition() int {
	return types.POS_DEAD
}

// Level returns the minimum level required to execute the command
func (c *CheckMobsCommand) Level() int {
	return 1 // Available to everyone during development
}

// LogCommand returns whether the command should be logged
func (c *CheckMobsCommand) LogCommand() bool {
	return true
}

// Execute executes the checkmobs command
func (c *CheckMobsCommand) Execute(ch *types.Character, args string) error {
	// Get the world from the character
	world, ok := ch.World.(interface{
		GetMobilePrototypes() []*types.Mobile
	})
	if !ok {
		return fmt.Errorf("world interface not available")
	}

	// Get all mobile prototypes
	mobilePrototypes := world.GetMobilePrototypes()

	// Sort mobile prototypes by VNUM
	sort.Slice(mobilePrototypes, func(i, j int) bool {
		return mobilePrototypes[i].VNUM < mobilePrototypes[j].VNUM
	})

	// Build a report of the mobile prototypes
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Found %d mobile prototypes:\r\n", len(mobilePrototypes)))
	sb.WriteString("VNUM  Name\r\n")
	sb.WriteString("----  ----\r\n")

	for _, mob := range mobilePrototypes {
		sb.WriteString(fmt.Sprintf("%-5d %s\r\n", mob.VNUM, mob.Name))
	}

	// Check for specific VNUMs
	missingVNUMs := []int{3000, 3001, 3004, 3006, 3046, 3100}
	sb.WriteString("\r\nChecking for specific VNUMs:\r\n")
	for _, vnum := range missingVNUMs {
		found := false
		for _, mob := range mobilePrototypes {
			if mob.VNUM == vnum {
				found = true
				sb.WriteString(fmt.Sprintf("VNUM %d: Found - %s\r\n", vnum, mob.Name))
				break
			}
		}
		if !found {
			sb.WriteString(fmt.Sprintf("VNUM %d: Not found\r\n", vnum))
		}
	}

	return fmt.Errorf("%s", sb.String())
}
