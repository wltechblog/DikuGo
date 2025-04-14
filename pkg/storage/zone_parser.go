package storage

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/wltechblog/DikuGo/pkg/types"
)

// ParseZones parses the zone file and returns a slice of zones
func ParseZones(filename string) ([]*types.Zone, error) {
	// Open the file
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %w", filename, err)
	}
	defer file.Close()

	// Create a scanner to read the file line by line
	scanner := bufio.NewScanner(file)
	var zones []*types.Zone
	var currentZone *types.Zone
	var lineNum int
	var state string
	var skipUntilNextZone bool

	// Read the file line by line
	for scanner.Scan() {
		lineNum++
		line := scanner.Text()

		// Skip empty lines
		if line == "" {
			continue
		}

		// Check for end of file marker
		if line == "$~" {
			break
		}

		// Check for new zone
		if strings.HasPrefix(line, "#") {
			// Save the current zone if it exists
			if currentZone != nil {
				zones = append(zones, currentZone)
			}

			// Reset the skip flag
			skipUntilNextZone = false

			// Parse zone number
			vnumStr := strings.TrimPrefix(line, "#")
			vnum, err := strconv.Atoi(vnumStr)
			if err != nil {
				return nil, fmt.Errorf("invalid zone number on line %d: %w", lineNum, err)
			}

			// Create a new zone
			currentZone = &types.Zone{
				VNUM:     vnum,
				Commands: []*types.ZoneCommand{},
			}

			// Set the state to read name next
			state = "name"
			continue
		}

		// Skip lines until we find a new zone if we're in skip mode
		if skipUntilNextZone {
			continue
		}

		// Process the line based on the current state
		switch state {
		case "name":
			currentZone.Name = strings.TrimSuffix(line, "~")
			state = "flags"
		case "flags":
			// Parse zone flags
			flagsParts := strings.Split(line, " ")
			if len(flagsParts) < 3 {
				// Skip this zone
				skipUntilNextZone = true
				continue
			}
			
			topRoom, err := strconv.Atoi(flagsParts[0])
			if err != nil {
				// Skip this zone
				skipUntilNextZone = true
				continue
			}
			
			lifespan, err := strconv.Atoi(flagsParts[1])
			if err != nil {
				// Skip this zone
				skipUntilNextZone = true
				continue
			}
			
			resetMode, err := strconv.Atoi(flagsParts[2])
			if err != nil {
				// Skip this zone
				skipUntilNextZone = true
				continue
			}
			
			currentZone.TopRoom = topRoom
			currentZone.Lifespan = lifespan
			currentZone.ResetMode = resetMode
			
			state = "commands"
		case "commands":
			// Check for end of zone
			if line == "S" {
				// End of zone
				state = "name"
				continue
			}
			
			// Skip comments
			if strings.HasPrefix(line, "*") {
				continue
			}
			
			// Parse zone command
			cmdParts := strings.Fields(line)
			if len(cmdParts) < 5 {
				// Skip this command
				continue
			}
			
			if len(cmdParts[0]) == 0 {
				continue
			}
			
			cmd := cmdParts[0][0]
			
			ifFlag, err := strconv.Atoi(cmdParts[1])
			if err != nil {
				// Skip this command
				continue
			}
			
			arg1, err := strconv.Atoi(cmdParts[2])
			if err != nil {
				// Skip this command
				continue
			}
			
			arg2, err := strconv.Atoi(cmdParts[3])
			if err != nil {
				// Skip this command
				continue
			}
			
			arg3, err := strconv.Atoi(cmdParts[4])
			if err != nil {
				// Skip this command
				continue
			}
			
			// Create the zone command
			zoneCmd := &types.ZoneCommand{
				Command: cmd,
				IfFlag:  ifFlag,
				Arg1:    arg1,
				Arg2:    arg2,
				Arg3:    arg3,
			}
			
			// Add the zone command to the zone
			currentZone.Commands = append(currentZone.Commands, zoneCmd)
		}
	}

	// Add the last zone if it exists and we're not skipping it
	if currentZone != nil && !skipUntilNextZone {
		zones = append(zones, currentZone)
	}

	return zones, nil
}
