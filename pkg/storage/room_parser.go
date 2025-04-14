package storage

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/wltechblog/DikuGo/pkg/types"
)

// readString reads a string from the parser until it encounters a line with just ~
// This matches the behavior of fread_string in the original C code
func readString(parser *Parser) string {
	var builder strings.Builder
	line := parser.Line()

	// Check if the line is just ~
	if line == "~" {
		// Empty string
		return ""
	}

	// Check if the line ends with ~
	if strings.HasSuffix(line, "~") {
		// Remove the ~ and return the string
		return strings.TrimSuffix(line, "~")
	}

	// Add the first line
	builder.WriteString(line)

	// Read lines until we find one that is just ~
	for parser.NextLine() {
		line = parser.Line()
		if line == "~" {
			// End of string
			break
		}
		builder.WriteString("\n")
		builder.WriteString(line)
	}

	return strings.TrimSpace(builder.String())
}

// ParseRooms parses the room file and returns a slice of rooms
func ParseRooms(filename string) ([]*types.Room, error) {
	parser, err := NewParser(filename)
	if err != nil {
		return nil, err
	}

	var rooms []*types.Room
	for parser.NextLine() {
		line := parser.Line()
		if !strings.HasPrefix(line, "#") {
			continue
		}

		// Parse room number
		vnumStr := strings.TrimPrefix(line, "#")
		vnum, err := strconv.Atoi(vnumStr)
		if err != nil {
			return nil, fmt.Errorf("invalid room number on line %d: %w", parser.LineNum(), err)
		}

		// Parse room name
		if !parser.NextLine() {
			return nil, fmt.Errorf("unexpected end of file after room number on line %d", parser.LineNum())
		}

		// Check for end of file marker
		if parser.Line() == "$~" {
			// End of file
			break
		}

		name := readString(parser)

		// Parse room description
		if !parser.NextLine() {
			return nil, fmt.Errorf("unexpected end of file after room name on line %d", parser.LineNum())
		}
		description := readString(parser)

		// Parse room flags
		if !parser.NextLine() {
			return nil, fmt.Errorf("unexpected end of file after room description on line %d", parser.LineNum())
		}
		flagsLine := parser.Line()

		// Skip any lines that don't look like flags
		for !strings.Contains(flagsLine, " ") || strings.HasSuffix(flagsLine, "~") {
			if !parser.NextLine() {
				return nil, fmt.Errorf("unexpected end of file while looking for room flags on line %d", parser.LineNum())
			}
			flagsLine = parser.Line()
		}

		flagsParts := strings.Split(flagsLine, " ")
		if len(flagsParts) < 3 {
			return nil, fmt.Errorf("invalid room flags on line %d: %s", parser.LineNum(), flagsLine)
		}

		// Zone number is not used in the Room struct
		_, err = strconv.Atoi(flagsParts[0])
		if err != nil {
			return nil, fmt.Errorf("invalid zone number on line %d: %w", parser.LineNum(), err)
		}

		roomFlags, err := strconv.Atoi(flagsParts[1])
		if err != nil {
			return nil, fmt.Errorf("invalid room flags on line %d: %w", parser.LineNum(), err)
		}

		sectorType, err := strconv.Atoi(flagsParts[2])
		if err != nil {
			return nil, fmt.Errorf("invalid sector type on line %d: %w", parser.LineNum(), err)
		}

		// Create the room
		room := &types.Room{
			VNUM:        vnum,
			Name:        name,
			Description: description,
			Flags:       uint32(roomFlags),
			SectorType:  sectorType,
			// Exits array is initialized with nil values by default
		}

		// Parse exits and extra descriptions
		for parser.NextLine() {
			line := parser.Line()
			if line == "S" {
				// End of room
				break
			}

			if strings.HasPrefix(line, "D") {
				// Parse exit
				dirStr := strings.TrimPrefix(line, "D")
				dirStr = strings.TrimSpace(dirStr)
				dir, err := strconv.Atoi(dirStr)
				if err != nil {
					return nil, fmt.Errorf("invalid exit direction on line %d: %w", parser.LineNum(), err)
				}

				// Parse exit description (general description)
				var exitDesc string
				if !parser.NextLine() {
					return nil, fmt.Errorf("unexpected end of file after exit direction on line %d", parser.LineNum())
				}
				exitDesc = readString(parser)

				// Parse exit keywords (used for doors)
				var exitKeywords string
				if !parser.NextLine() {
					return nil, fmt.Errorf("unexpected end of file after exit description on line %d", parser.LineNum())
				}
				exitKeywords = readString(parser)

				// Parse exit flags and destination
				if !parser.NextLine() {
					return nil, fmt.Errorf("unexpected end of file after exit keywords on line %d", parser.LineNum())
				}

				exitFlagsLine := parser.Line()

				// If the line is just ~ or empty, read the next line for the exit flags
				if exitFlagsLine == "~" || exitFlagsLine == "" {
					if !parser.NextLine() {
						return nil, fmt.Errorf("unexpected end of file after exit keywords delimiter on line %d", parser.LineNum())
					}
					exitFlagsLine = parser.Line()

					// Skip any empty lines
					for exitFlagsLine == "" {
						if !parser.NextLine() {
							return nil, fmt.Errorf("unexpected end of file after exit keywords delimiter on line %d", parser.LineNum())
						}
						exitFlagsLine = parser.Line()
					}
				}

				// Check if this is actually a description line
				if strings.Contains(exitFlagsLine, " ") && !strings.HasPrefix(exitFlagsLine, "0 ") && !strings.HasPrefix(exitFlagsLine, "1 ") && !strings.HasPrefix(exitFlagsLine, "2 ") {
					// This is a description line, not a flags line
					// Skip until we find a line that looks like flags
					for parser.NextLine() {
						line := parser.Line()
						if strings.HasPrefix(line, "0 ") || strings.HasPrefix(line, "1 ") || strings.HasPrefix(line, "2 ") {
							exitFlagsLine = line
							break
						}
					}
				}

				// Parse the exit flags, key, and destination
				exitFlagsParts := strings.Fields(exitFlagsLine)
				if len(exitFlagsParts) < 3 {
					return nil, fmt.Errorf("invalid exit flags/key/destination on line %d: %s", parser.LineNum(), exitFlagsLine)
				}

				// Parse the exit flags (0 = no door, 1 = door, 2 = pickproof door)
				exitFlags, err := strconv.Atoi(exitFlagsParts[0])
				if err != nil {
					return nil, fmt.Errorf("invalid exit flags on line %d: %w", parser.LineNum(), err)
				}

				// Parse the exit key (item VNUM of the key, or -1 if no key)
				exitKey, err := strconv.Atoi(exitFlagsParts[1])
				if err != nil {
					return nil, fmt.Errorf("invalid exit key on line %d: %w", parser.LineNum(), err)
				}

				// Parse the destination room VNUM
				destRoomVnum, err := strconv.Atoi(exitFlagsParts[2])
				if err != nil {
					return nil, fmt.Errorf("invalid exit destination on line %d: %w", parser.LineNum(), err)
				}

				// In the original C code, -1 is used for NOWHERE, which means the exit doesn't lead anywhere
				// We'll keep this behavior for compatibility

				// Create the exit
				exit := &types.Exit{
					Direction:   dir,
					Description: exitDesc,
					Keywords:    exitKeywords,
					Flags:       uint32(exitFlags),
					Key:         exitKey,
					DestVnum:    destRoomVnum, // Store the destination room VNUM
				}

				log.Printf("DEBUG: Room %d, Exit direction %d, destRoomVnum: %d", room.VNUM, dir, destRoomVnum)

				// Add the exit to the room
				room.Exits[dir] = exit
			} else if strings.HasPrefix(line, "E") {
				// Extra description
				// Read the keywords
				if !parser.NextLine() {
					return nil, fmt.Errorf("unexpected end of file after extra description on line %d", parser.LineNum())
				}
				keywords := readString(parser)

				// Read the description
				if !parser.NextLine() {
					return nil, fmt.Errorf("unexpected end of file after extra description keywords on line %d", parser.LineNum())
				}
				description := readString(parser)

				// Create the extra description
				extraDesc := &types.ExtraDescription{
					Keywords:    keywords,
					Description: description,
				}

				// Add the extra description to the room
				room.ExtraDescs = append(room.ExtraDescs, extraDesc)

				// Check if we need to process another directive
				if parser.Line() != "" && parser.Line() != "~" {
					// We need to process this line again
					continue
				}
			}
		}

		// Add the room to the list
		rooms = append(rooms, room)
	}

	return rooms, nil
}
