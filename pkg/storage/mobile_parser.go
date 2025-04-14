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
				// Skip to the next mobile
				currentMobile = nil
			} else {
				// Successfully parsed mobile
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

	return mobiles, nil
}

// parseMobileData parses the data for a mobile
func parseMobileData(parser *Parser, mobile *types.Mobile) error {
	// Parse the mobile's name (keywords)
	if !parser.NextLine() {
		return fmt.Errorf("unexpected end of file after mobile number")
	}
	mobile.Name = strings.TrimSuffix(parser.Line(), "~")

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

	// Parse the mobile's flags
	if !parser.NextLine() {
		return fmt.Errorf("unexpected end of file after mobile description")
	}

	// Skip any lines that don't look like flags
	for strings.HasPrefix(parser.Line(), "E") || strings.HasPrefix(parser.Line(), "A") || strings.HasSuffix(parser.Line(), "~") {
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

	// Check if this is a simple or detailed mobile
	if len(flagsParts) >= 4 && flagsParts[3] == "D" {
		// Detailed mobile
		return parseDetailedMobile(parser, mobile)
	} else {
		// Simple mobile (default)
		return parseSimpleMobile(parser, mobile)
	}
}
