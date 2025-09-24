package command

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/wltechblog/DikuGo/pkg/storage"
	"github.com/wltechblog/DikuGo/pkg/types"
)

// TestMobParserCommand represents the testmobparser command
type TestMobParserCommand struct{}

// Name returns the name of the command
func (c *TestMobParserCommand) Name() string {
	return "testmobparser"
}

// Aliases returns the aliases of the command
func (c *TestMobParserCommand) Aliases() []string {
	return []string{}
}

// MinPosition returns the minimum position required to execute the command
func (c *TestMobParserCommand) MinPosition() int {
	return types.POS_DEAD
}

// Level returns the minimum level required to execute the command
func (c *TestMobParserCommand) Level() int {
	return 1 // Available to everyone during development
}

// LogCommand returns whether the command should be logged
func (c *TestMobParserCommand) LogCommand() bool {
	return true
}

// Execute executes the testmobparser command
func (c *TestMobParserCommand) Execute(ch *types.Character, args string) error {
	// Check if a VNUM range was provided
	var startVnum, endVnum int
	var err error
	
	if args == "" {
		// Default to testing all mobs
		startVnum = 0
		endVnum = 10000
	} else {
		// Parse the VNUM range
		parts := strings.Split(args, "-")
		if len(parts) == 1 {
			// Single VNUM
			startVnum, err = strconv.Atoi(parts[0])
			if err != nil {
				return fmt.Errorf("Invalid VNUM: %s", parts[0])
			}
			endVnum = startVnum
		} else if len(parts) == 2 {
			// VNUM range
			startVnum, err = strconv.Atoi(parts[0])
			if err != nil {
				return fmt.Errorf("Invalid start VNUM: %s", parts[0])
			}
			endVnum, err = strconv.Atoi(parts[1])
			if err != nil {
				return fmt.Errorf("Invalid end VNUM: %s", parts[1])
			}
		} else {
			return fmt.Errorf("Usage: testmobparser [vnum | vnum-vnum]")
		}
	}

	// Parse the mobile file directly
	mobiles, err := storage.ParseMobiles(filepath.Join("old/lib", "tinyworld.mob"))
	if err != nil {
		return fmt.Errorf("Failed to parse mobile file: %v", err)
	}

	// Filter mobiles by VNUM range
	var filteredMobiles []*types.Mobile
	for _, mob := range mobiles {
		if mob.VNUM >= startVnum && mob.VNUM <= endVnum {
			filteredMobiles = append(filteredMobiles, mob)
		}
	}

	// Build a report of the parsed mobiles
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Parsed %d mobiles in range %d-%d:\r\n", len(filteredMobiles), startVnum, endVnum))
	sb.WriteString("VNUM  Name                 Level  HitRoll  DamRoll  Gold    Exp\r\n")
	sb.WriteString("----  -------------------- ------ -------- -------- ------- -------\r\n")

	for _, mob := range filteredMobiles {
		name := mob.Name
		if len(name) > 20 {
			name = name[:17] + "..."
		} else {
			name = fmt.Sprintf("%-20s", name)
		}
		
		sb.WriteString(fmt.Sprintf("%-5d %s %-6d %-8d %-8d %-7d %-7d\r\n", 
			mob.VNUM, name, mob.Level, mob.HitRoll, mob.DamRoll, mob.Gold, mob.Experience))
	}

	// Get the world from the character
	world, ok := ch.World.(interface{
		GetMobile(int) *types.Mobile
	})
	if !ok {
		return fmt.Errorf("world interface not available")
	}

	// Compare with the mobiles in the world
	sb.WriteString("\r\nComparing with mobiles in the world:\r\n")
	sb.WriteString("VNUM  Name                 Parser Name            World Name\r\n")
	sb.WriteString("----  -------------------- ---------------------- ----------------------\r\n")

	for _, parsedMob := range filteredMobiles {
		worldMob := world.GetMobile(parsedMob.VNUM)
		
		parsedName := parsedMob.Name
		if len(parsedName) > 20 {
			parsedName = parsedName[:17] + "..."
		} else {
			parsedName = fmt.Sprintf("%-20s", parsedName)
		}
		
		worldName := "Not found"
		if worldMob != nil {
			worldName = worldMob.Name
			if len(worldName) > 20 {
				worldName = worldName[:17] + "..."
			}
		}
		
		sb.WriteString(fmt.Sprintf("%-5d %s %-22s %-22s\r\n", 
			parsedMob.VNUM, parsedName, parsedMob.Name, worldName))
	}

	return fmt.Errorf("%s", sb.String())
}
