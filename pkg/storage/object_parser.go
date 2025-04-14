package storage

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/wltechblog/DikuGo/pkg/types"
)

// ParseObjects parses the object file and returns a slice of objects
func ParseObjects(filename string) ([]*types.Object, error) {
	// Open the file
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %w", filename, err)
	}
	defer file.Close()

	// Create a scanner to read the file line by line
	scanner := bufio.NewScanner(file)
	var objects []*types.Object
	var currentObject *types.Object
	var lineNum int
	var state string
	var skipUntilNextObject bool

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

		// Check for new object
		if strings.HasPrefix(line, "#") {
			// Save the current object if it exists
			if currentObject != nil {
				objects = append(objects, currentObject)
			}

			// Reset the skip flag
			skipUntilNextObject = false

			// Parse object number
			vnumStr := strings.TrimPrefix(line, "#")
			vnum, err := strconv.Atoi(vnumStr)
			if err != nil {
				return nil, fmt.Errorf("invalid object number on line %d: %w", lineNum, err)
			}

			// Create a new object
			currentObject = &types.Object{
				VNUM: vnum,
			}

			// Set the state to read keywords next
			state = "keywords"
			continue
		}

		// Skip lines until we find a new object if we're in skip mode
		if skipUntilNextObject {
			continue
		}

		// Process the line based on the current state
		switch state {
		case "keywords":
			currentObject.Name = strings.TrimSuffix(line, "~")
			state = "shortdesc"
		case "shortdesc":
			currentObject.ShortDesc = strings.TrimSuffix(line, "~")
			state = "longdesc"
		case "longdesc":
			currentObject.Description = strings.TrimSuffix(line, "~")
			state = "actiondesc"
		case "actiondesc":
			// Action description might be multi-line
			if strings.HasSuffix(line, "~") {
				currentObject.ActionDesc = strings.TrimSuffix(line, "~")
				state = "flags"
			} else {
				// Start collecting the action description
				currentObject.ActionDesc = line
				state = "actiondesc_cont"
			}
		case "actiondesc_cont":
			// Continue collecting the action description
			if strings.HasSuffix(line, "~") {
				currentObject.ActionDesc += "\n" + strings.TrimSuffix(line, "~")
				state = "flags"
			} else {
				currentObject.ActionDesc += "\n" + line
			}
		case "flags":
			// Check for extra descriptions or affects
			if strings.HasPrefix(line, "E") {
				// Extra description
				// Read the keywords
				if !scanner.Scan() {
					return nil, fmt.Errorf("unexpected end of file after extra description on line %d", lineNum)
				}
				lineNum++
				keywords := strings.TrimSuffix(scanner.Text(), "~")

				// Read the description
				var description strings.Builder
				for scanner.Scan() {
					lineNum++
					line = scanner.Text()
					if line == "~" {
						break
					}
					if description.Len() > 0 {
						description.WriteString("\n")
					}
					description.WriteString(line)
				}

				// Create the extra description
				extraDesc := &types.ExtraDescription{
					Keywords:    keywords,
					Description: description.String(),
				}

				// Add the extra description to the object
				currentObject.ExtraDescs = append(currentObject.ExtraDescs, extraDesc)

				// Continue to the next line
				continue
			} else if strings.HasPrefix(line, "A") {
				// Skip affects until we find a line that looks like flags
				for scanner.Scan() {
					lineNum++
					line = scanner.Text()

					// Skip empty lines
					if line == "" {
						continue
					}

					// Check for new object
					if strings.HasPrefix(line, "#") {
						// We've reached a new object, so we need to back up and process it
						// in the main loop
						skipUntilNextObject = true
						break
					}

					// Check if this line looks like flags
					if !strings.HasPrefix(line, "E") && !strings.HasPrefix(line, "A") && !strings.HasSuffix(line, "~") {
						// This might be a flags line
						flagsParts := strings.Split(line, " ")
						if len(flagsParts) >= 3 {
							// Try to parse the first part as an integer
							if _, err := strconv.Atoi(flagsParts[0]); err == nil {
								// This is probably a flags line
								break
							}
						}
					}
				}

				// If we're skipping until the next object, continue
				if skipUntilNextObject {
					continue
				}
			}

			// Parse object flags
			flagsParts := strings.Split(line, " ")
			if len(flagsParts) < 3 {
				// Skip this object
				skipUntilNextObject = true
				continue
			}

			objType, err := strconv.Atoi(flagsParts[0])
			if err != nil {
				// Skip this object
				skipUntilNextObject = true
				continue
			}

			extraFlags, err := strconv.Atoi(flagsParts[1])
			if err != nil {
				// Skip this object
				skipUntilNextObject = true
				continue
			}

			wearFlags, err := strconv.Atoi(flagsParts[2])
			if err != nil {
				// Skip this object
				skipUntilNextObject = true
				continue
			}

			currentObject.Type = objType
			currentObject.ExtraFlags = uint32(extraFlags)
			currentObject.WearFlags = uint32(wearFlags)

			state = "values"
		case "values":
			// Parse object values
			valuesParts := strings.Split(line, " ")
			if len(valuesParts) < 4 {
				// Skip this object
				skipUntilNextObject = true
				continue
			}

			var values [4]int
			for i := 0; i < 4; i++ {
				value, err := strconv.Atoi(valuesParts[i])
				if err != nil {
					// Skip this object
					skipUntilNextObject = true
					break
				}
				values[i] = value
			}

			if skipUntilNextObject {
				continue
			}

			currentObject.Value = values

			state = "weight"
		case "weight":
			// Parse object weight, cost, and rent
			weightParts := strings.Split(line, " ")
			if len(weightParts) < 2 {
				// Skip this object
				skipUntilNextObject = true
				continue
			}

			weight, err := strconv.Atoi(weightParts[0])
			if err != nil {
				// Skip this object
				skipUntilNextObject = true
				continue
			}

			cost, err := strconv.Atoi(weightParts[1])
			if err != nil {
				// Skip this object
				skipUntilNextObject = true
				continue
			}

			currentObject.Weight = weight
			currentObject.Cost = cost

			// We're done with this object, skip until we find a new one
			skipUntilNextObject = true
		}
	}

	// Add the last object if it exists and we're not skipping it
	if currentObject != nil && !skipUntilNextObject {
		objects = append(objects, currentObject)
	}

	return objects, nil
}
