package command

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/wltechblog/DikuGo/pkg/types"
)

// ExamineMobCommand represents the examinemob command
type ExamineMobCommand struct{}

// Name returns the name of the command
func (c *ExamineMobCommand) Name() string {
	return "examinemob"
}

// Aliases returns the aliases of the command
func (c *ExamineMobCommand) Aliases() []string {
	return []string{"mobexamine"}
}

// MinPosition returns the minimum position required to execute the command
func (c *ExamineMobCommand) MinPosition() int {
	return types.POS_DEAD
}

// Level returns the minimum level required to execute the command
func (c *ExamineMobCommand) Level() int {
	return 1 // Available to everyone during development
}

// LogCommand returns whether the command should be logged
func (c *ExamineMobCommand) LogCommand() bool {
	return true
}

// Execute executes the examinemob command
func (c *ExamineMobCommand) Execute(ch *types.Character, args string) error {
	// Check if a VNUM was provided
	if args == "" {
		return fmt.Errorf("Usage: examinemob <vnum>\r\nExample: examinemob 3001")
	}

	// Parse the VNUM
	vnum, err := strconv.Atoi(args)
	if err != nil {
		return fmt.Errorf("Invalid VNUM: %s", args)
	}

	// Get the world from the character
	world, ok := ch.World.(interface{
		GetMobile(int) *types.Mobile
	})
	if !ok {
		return fmt.Errorf("world interface not available")
	}

	// Get the mobile prototype
	mob := world.GetMobile(vnum)
	if mob == nil {
		return fmt.Errorf("Mobile VNUM %d not found in the world", vnum)
	}

	// Build a report of the mobile prototype
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Mobile VNUM %d:\r\n", vnum))
	sb.WriteString(fmt.Sprintf("Name: %s\r\n", mob.Name))
	sb.WriteString(fmt.Sprintf("Short Description: %s\r\n", mob.ShortDesc))
	sb.WriteString(fmt.Sprintf("Long Description: %s\r\n", mob.LongDesc))
	sb.WriteString(fmt.Sprintf("Description: %s\r\n", mob.Description))
	sb.WriteString(fmt.Sprintf("Level: %d\r\n", mob.Level))
	sb.WriteString(fmt.Sprintf("HitRoll: %d\r\n", mob.HitRoll))
	sb.WriteString(fmt.Sprintf("DamRoll: %d\r\n", mob.DamRoll))
	sb.WriteString(fmt.Sprintf("AC: %v\r\n", mob.AC))
	sb.WriteString(fmt.Sprintf("Gold: %d\r\n", mob.Gold))
	sb.WriteString(fmt.Sprintf("Experience: %d\r\n", mob.Experience))
	sb.WriteString(fmt.Sprintf("Position: %d\r\n", mob.Position))
	sb.WriteString(fmt.Sprintf("Default Position: %d\r\n", mob.DefaultPos))
	sb.WriteString(fmt.Sprintf("Sex: %d\r\n", mob.Sex))
	sb.WriteString(fmt.Sprintf("Abilities: %v\r\n", mob.Abilities))
	sb.WriteString(fmt.Sprintf("Dice: %v\r\n", mob.Dice))
	sb.WriteString(fmt.Sprintf("Act Flags: %d\r\n", mob.ActFlags))
	sb.WriteString(fmt.Sprintf("Affect Flags: %d\r\n", mob.AffectFlags))
	sb.WriteString(fmt.Sprintf("Alignment: %d\r\n", mob.Alignment))

	// Now let's examine the raw mobile file to see what's in it
	sb.WriteString("\r\nRaw mobile file data:\r\n")
	
	// Open the mobile file
	file, err := os.Open(filepath.Join("old/lib", "tinyworld.mob"))
	if err != nil {
		return fmt.Errorf("Failed to open mobile file: %v", err)
	}
	defer file.Close()

	// Create a scanner to read the file line by line
	scanner := bufio.NewScanner(file)
	
	// Find the mobile in the file
	found := false
	var rawLines []string
	for scanner.Scan() {
		line := scanner.Text()
		
		// Check if this is the start of the mobile we're looking for
		if line == fmt.Sprintf("#%d", vnum) {
			found = true
			rawLines = append(rawLines, line)
			continue
		}
		
		// If we've found the mobile, collect its lines until we reach the next mobile or end of file
		if found {
			rawLines = append(rawLines, line)
			
			// Check if this is the start of the next mobile
			if strings.HasPrefix(line, "#") && line != fmt.Sprintf("#%d", vnum) {
				break
			}
		}
	}
	
	if !found {
		sb.WriteString(fmt.Sprintf("Mobile VNUM %d not found in the raw file\r\n", vnum))
	} else {
		for _, line := range rawLines {
			sb.WriteString(fmt.Sprintf("%s\r\n", line))
		}
	}

	return fmt.Errorf("%s", sb.String())
}
