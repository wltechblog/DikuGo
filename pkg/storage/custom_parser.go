package storage

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/wltechblog/DikuGo/pkg/types"
)

// ParseCustomMobiles parses a mobile file and returns a slice of mobiles
func ParseCustomMobiles(filename string) ([]*types.Mobile, error) {
	// Open the file
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %w", filename, err)
	}
	defer file.Close()

	// Read the file content
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", filename, err)
	}

	// Split the content into lines
	lines := strings.Split(string(content), "\n")

	var mobiles []*types.Mobile
	var currentMobile *types.Mobile
	var lineIndex int

	// Process each line
	for lineIndex < len(lines) {
		line := lines[lineIndex]
		lineIndex++

		// Skip empty lines
		if line == "" {
			continue
		}

		// Check for new mobile
		if strings.HasPrefix(line, "#") {
			// Parse mobile number
			vnumStr := strings.TrimPrefix(line, "#")
			vnum, err := strconv.Atoi(vnumStr)
			if err != nil {
				log.Printf("Warning: invalid mobile number on line %d: %v, skipping", lineIndex, err)
				continue
			}

			// Create a new mobile
			currentMobile = &types.Mobile{
				VNUM: vnum,
				Dice: [3]int{0, 0, 0},
				AC:   [3]int{0, 0, 0},
			}

			// Parse name
			if lineIndex < len(lines) {
				currentMobile.Name = strings.TrimSuffix(lines[lineIndex], "~")
				lineIndex++
			} else {
				log.Printf("Warning: unexpected end of file after mobile number")
				break
			}

			// Parse short description
			if lineIndex < len(lines) {
				currentMobile.ShortDesc = strings.TrimSuffix(lines[lineIndex], "~")
				lineIndex++
			} else {
				log.Printf("Warning: unexpected end of file after mobile name")
				break
			}

			// Parse long description
			if lineIndex < len(lines) {
				currentMobile.LongDesc = strings.TrimSuffix(lines[lineIndex], "~")
				lineIndex++
			} else {
				log.Printf("Warning: unexpected end of file after mobile short description")
				break
			}

			// Parse description (multi-line)
			var description strings.Builder
			for lineIndex < len(lines) {
				descLine := lines[lineIndex]
				lineIndex++
				if descLine == "~" {
					break
				}
				description.WriteString(descLine)
				description.WriteString("\n")
			}
			currentMobile.Description = description.String()

			// Parse flags
			if lineIndex < len(lines) {
				flagsLine := lines[lineIndex]
				lineIndex++
				flagsParts := strings.Fields(flagsLine)
				if len(flagsParts) >= 3 {
					// Parse act flags
					actFlags, err := strconv.Atoi(flagsParts[0])
					if err == nil {
						currentMobile.ActFlags = uint32(actFlags)
					}

					// Parse affect flags
					affectFlags, err := strconv.Atoi(flagsParts[1])
					if err == nil {
						currentMobile.AffectFlags = uint32(affectFlags)
					}

					// Parse alignment
					alignment, err := strconv.Atoi(flagsParts[2])
					if err == nil {
						currentMobile.Alignment = alignment
					}
				}
			} else {
				log.Printf("Warning: unexpected end of file after mobile description")
				break
			}

			// Parse level line
			if lineIndex < len(lines) {
				levelLine := lines[lineIndex]
				lineIndex++
				levelParts := strings.Fields(levelLine)
				if len(levelParts) >= 5 {
					// Parse level
					level, err := strconv.Atoi(levelParts[0])
					if err == nil {
						currentMobile.Level = level
					}

					// Parse hitroll
					hitroll, err := strconv.Atoi(levelParts[1])
					if err == nil {
						currentMobile.HitRoll = 20 - hitroll // Convert THAC0 to hitroll
					}

					// Parse AC
					ac, err := strconv.Atoi(levelParts[2])
					if err == nil {
						currentMobile.AC = [3]int{10 * ac, 10 * ac, 10 * ac}
					}

					// Parse hit dice
					hitDice := levelParts[3]
					parts := strings.Split(hitDice, "+")
					if len(parts) == 2 {
						// Parse the dice part (e.g., "5d10")
						diceParts := strings.Split(parts[0], "d")
						if len(diceParts) == 2 {
							numDice, err := strconv.Atoi(diceParts[0])
							if err == nil {
								currentMobile.Dice[0] = numDice
							}

							sizeDice, err := strconv.Atoi(diceParts[1])
							if err == nil {
								currentMobile.Dice[1] = sizeDice
							}
						}

						// Parse the bonus part
						bonus, err := strconv.Atoi(parts[1])
						if err == nil {
							currentMobile.Dice[2] = bonus
						}
					}

					// Parse damage dice
					if len(levelParts) > 4 {
						damageDice := levelParts[4]
						parts = strings.Split(damageDice, "+")
						if len(parts) == 2 {
							// Parse the dice part (e.g., "1d8")
							diceParts := strings.Split(parts[0], "d")
							if len(diceParts) == 2 {
								numDice, err := strconv.Atoi(diceParts[0])
								if err == nil {
									currentMobile.DamageType = numDice
								}

								sizeDice, err := strconv.Atoi(diceParts[1])
								if err == nil {
									currentMobile.AttackType = sizeDice
								}
							}

							// Parse the bonus part
							bonus, err := strconv.Atoi(parts[1])
							if err == nil {
								currentMobile.DamRoll = bonus
							}
						}
					}
				}
			} else {
				log.Printf("Warning: unexpected end of file after flags")
				break
			}

			// Parse gold and experience
			if lineIndex < len(lines) {
				goldExpLine := lines[lineIndex]
				lineIndex++
				goldExpParts := strings.Fields(goldExpLine)
				if len(goldExpParts) >= 2 {
					// Parse gold
					gold, err := strconv.Atoi(goldExpParts[0])
					if err == nil {
						currentMobile.Gold = gold
					}

					// Parse experience
					exp, err := strconv.Atoi(goldExpParts[1])
					if err == nil {
						currentMobile.Experience = exp
					}
				}
			} else {
				log.Printf("Warning: unexpected end of file after level line")
				break
			}

			// Parse position, default position, and sex
			if lineIndex < len(lines) {
				posLine := lines[lineIndex]
				lineIndex++
				posParts := strings.Fields(posLine)
				if len(posParts) >= 3 {
					// Parse position
					position, err := strconv.Atoi(posParts[0])
					if err == nil {
						currentMobile.Position = position
					}

					// Parse default position
					defaultPos, err := strconv.Atoi(posParts[1])
					if err == nil {
						currentMobile.DefaultPos = defaultPos
					}

					// Parse sex
					sex, err := strconv.Atoi(posParts[2])
					if err == nil {
						currentMobile.Sex = sex
					}
				}
			} else {
				log.Printf("Warning: unexpected end of file after gold/exp line")
				break
			}

			// Add the mobile to the list
			mobiles = append(mobiles, currentMobile)
		}
	}

	return mobiles, nil
}
