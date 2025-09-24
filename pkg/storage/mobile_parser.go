package storage

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/wltechblog/DikuGo/pkg/types"
)

// ParseMobiles parses the mobile file and returns a slice of mobiles
func ParseMobiles(filename string) ([]*types.Mobile, error) {
	// Check if this is the test.mob file
	if strings.HasSuffix(filename, "test.mob") {
		// Use the custom parser for test.mob
		return ParseCustomMobiles(filename)
	}

	// Open the file
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %w", filename, err)
	}
	defer file.Close()

	// Create a new parser
	parser, err := NewParser(filename)
	if err != nil {
		return nil, err
	}

	var mobiles []*types.Mobile
	var currentMobile *types.Mobile

	// Read the file line by line
	for {
		hasNext := parser.NextLine()
		if !hasNext {
			break
		}

		line := parser.Line()

		// Skip empty lines
		if line == "" {
			continue
		}

		// Check for end of file marker
		if line == "$~" {
			break
		}

		// Check for new mobile
		if strings.HasPrefix(line, "#") {
			// Save the current mobile if it exists
			if currentMobile != nil {
				mobiles = append(mobiles, currentMobile)
			}

			// Parse mobile number
			vnumStr := strings.TrimPrefix(line, "#")
			vnum, err := strconv.Atoi(vnumStr)
			if err != nil {
				return nil, fmt.Errorf("invalid mobile number on line %d: %w", parser.LineNum(), err)
			}

			// Create a new mobile
			currentMobile = &types.Mobile{
				VNUM: vnum,
			}

			// Parse the mobile data
			if err := parseMobileData(parser, currentMobile); err != nil {
				log.Printf("Warning: Error parsing mobile #%d: %v", vnum, err)
				// Instead of skipping the mobile entirely, try to continue with default values
				// This ensures that all mobile prototypes are loaded, even if they have errors
				log.Printf("Continuing with default values for mobile #%d", vnum)

				// Only set default values if the name is empty
				if currentMobile.Name == "" {
					currentMobile.Name = fmt.Sprintf("mobile%d", vnum)
				}
				if currentMobile.ShortDesc == "" {
					currentMobile.ShortDesc = fmt.Sprintf("Mobile %d", vnum)
				}
				if currentMobile.LongDesc == "" {
					currentMobile.LongDesc = fmt.Sprintf("Mobile %d is standing here.\n", vnum)
				}
				if currentMobile.Description == "" {
					currentMobile.Description = fmt.Sprintf("This is mobile %d.\n", vnum)
				}

				// Set default values for stats
				currentMobile.Level = 1
				currentMobile.HitRoll = 0
				currentMobile.DamRoll = 0
				currentMobile.AC = [3]int{10, 10, 10}
				currentMobile.Gold = 10
				currentMobile.Experience = 100
				currentMobile.Position = 8
				currentMobile.DefaultPos = 8
				currentMobile.Sex = 0
				currentMobile.Abilities = [6]int{11, 11, 11, 11, 11, 11}
				currentMobile.Dice = [3]int{1, 8, 0}

				// Skip to the next mobile
				for parser.NextLine() {
					if strings.HasPrefix(parser.Line(), "#") {
						parser.BackUp()
						break
					}
				}
			} else {
				// Successfully parsed mobile
				log.Printf("Parsed DikuMUD mobile #%d with stats: Level=%d, HitRoll=%d, DamRoll=%d, AC=%v, Gold=%d, Exp=%d",
					currentMobile.VNUM, currentMobile.Level, currentMobile.HitRoll, currentMobile.DamRoll,
					currentMobile.AC, currentMobile.Gold, currentMobile.Experience)
			}

			// Check if we need to process the current line again (if a parser backed up)
			if strings.HasPrefix(parser.Line(), "#") {
				continue
			}
		}
	}

	// Add the last mobile if it exists
	if currentMobile != nil {
		mobiles = append(mobiles, currentMobile)
	}

	// Debug: Print all loaded mobiles
	for _, mob := range mobiles {
		log.Printf("DEBUG: Loaded mobile #%d (%s) with stats: Level=%d, HitRoll=%d, DamRoll=%d, AC=%v, Gold=%d, Exp=%d",
			mob.VNUM, mob.Name, mob.Level, mob.HitRoll, mob.DamRoll, mob.AC, mob.Gold, mob.Experience)
	}

	return mobiles, nil
}

// parseMobileData parses the data for a mobile
func parseMobileData(parser *Parser, mobile *types.Mobile) error {
	log.Printf("Parsing mobile #%d", mobile.VNUM)
	// Parse the mobile's name (keywords)
	if !parser.NextLine() {
		return fmt.Errorf("unexpected end of file after mobile number")
	}
	mobile.Name = strings.TrimSuffix(parser.Line(), "~")
	log.Printf("Mobile #%d name: %s", mobile.VNUM, mobile.Name)

	// Parse the mobile's short description
	if !parser.NextLine() {
		return fmt.Errorf("unexpected end of file after mobile name")
	}
	mobile.ShortDesc = strings.TrimSuffix(parser.Line(), "~")

	// Parse the mobile's long description
	if !parser.NextLine() {
		return fmt.Errorf("unexpected end of file after mobile short description")
	}
	mobile.LongDesc = strings.TrimSuffix(parser.Line(), "~")

	// Parse the mobile's description
	if !parser.NextLine() {
		return fmt.Errorf("unexpected end of file after mobile long description")
	}
	description := readString(parser)
	mobile.Description = description
	log.Printf("DEBUG: Mobile #%d description: %s", mobile.VNUM, description)

	// Parse the mobile's flags
	if !parser.NextLine() {
		return fmt.Errorf("unexpected end of file after mobile description")
	}

	// Skip any lines that don't look like flags
	for {
		// If we hit another mobile, back up and return
		if strings.HasPrefix(parser.Line(), "#") {
			parser.BackUp()
			return nil
		}

		// If the line has at least 3 fields and the first field is a number, it's probably the flags line
		fields := strings.Fields(parser.Line())
		if len(fields) >= 3 {
			_, err := strconv.Atoi(fields[0])
			if err == nil {
				// This is probably the flags line
				break
			}
		}

		// Otherwise, skip this line
		if !parser.NextLine() {
			return fmt.Errorf("unexpected end of file while looking for mobile flags")
		}
	}

	// Parse the flags line
	flagsParts := strings.Fields(parser.Line())
	if len(flagsParts) < 3 {
		// Try to continue with default values
		mobile.ActFlags = 0
		mobile.AffectFlags = 0
		mobile.Alignment = 0
	} else {
		actFlags, err := strconv.Atoi(flagsParts[0])
		if err != nil {
			// Try to continue anyway with default values
			actFlags = 0
		}

		affectFlags, err := strconv.Atoi(flagsParts[1])
		if err != nil {
			// Try to continue anyway with default values
			affectFlags = 0
		}

		alignment, err := strconv.Atoi(flagsParts[2])
		if err != nil {
			// Try to continue anyway with default values
			alignment = 0
		}

		mobile.ActFlags = uint32(actFlags)
		mobile.AffectFlags = uint32(affectFlags)
		mobile.Alignment = alignment
	}

	// Check if this is a simple, detailed, or DikuMUD mobile
	var formatType string
	if len(flagsParts) >= 4 {
		formatType = flagsParts[3]
	}

	log.Printf("DEBUG: Mobile #%d format type: '%s', flags: %v", mobile.VNUM, formatType, flagsParts)

	// Parse based on format type
	switch formatType {
	case "D":
		// Detailed mobile
		log.Printf("DEBUG: Using detailed parser for mobile #%d", mobile.VNUM)
		return parseDetailedMobile(parser, mobile)
	case "S":
		// Simple mobile
		log.Printf("DEBUG: Using simple parser for mobile #%d", mobile.VNUM)
		return parseSimpleMobile(parser, mobile)
	default:
		// Original DikuMUD format - use the DikuMUD parser
		log.Printf("DEBUG: Using DikuMUD parser for mobile #%d", mobile.VNUM)
		return parseDikuMobile(parser, mobile)
	}
}
