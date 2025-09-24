package storage

import (
	"log"
	"strconv"
	"strings"

	"github.com/wltechblog/DikuGo/pkg/types"
)

// parseSimpleMobile parses a simple mobile (type 'S')
func parseSimpleMobile(parser *Parser, mobile *types.Mobile) error {
	log.Printf("Parsing simple mobile #%d", mobile.VNUM)
	log.Printf("DEBUG: Starting to parse simple mobile #%d with line: %s", mobile.VNUM, parser.Line())
	// Set default abilities
	// In the original code, simple mobiles have default ability scores of 11
	mobile.Abilities = [6]int{11, 11, 11, 11, 11, 11} // STR, INT, WIS, DEX, CON, CHA

	// Check if we need to parse the current line or get the next line
	currentLine := parser.Line()
	if strings.Contains(currentLine, "S") {
		// We're already on the flags line, so we need to get the next line
		if !parser.NextLine() {
			log.Printf("Warning: unexpected end of file while parsing stats for mobile #%d, using default values", mobile.VNUM)
			mobile.Level = 1
			mobile.HitRoll = 0
			mobile.AC = [3]int{10, 10, 10}
			mobile.Gold = 10
			mobile.Experience = 100
			return nil
		}
	}

	// Parse the stats line
	statsLine := strings.TrimSpace(parser.Line())
	log.Printf("DEBUG: Parsing stats for mobile #%d, line: '%s'", mobile.VNUM, statsLine)

	// Split the stats line into parts
	statsParts := strings.Fields(statsLine)
	if len(statsParts) < 5 {
		log.Printf("Warning: invalid stats line for mobile #%d: %s, using default values", mobile.VNUM, statsLine)
		mobile.Level = 1
		mobile.HitRoll = 0
		mobile.AC = [3]int{10, 10, 10}
		mobile.Gold = 10
		mobile.Experience = 100
		return nil
	}

	// Parse level
	level, err := strconv.Atoi(statsParts[0])
	if err != nil {
		log.Printf("Warning: invalid level for mobile #%d: %v, using default value", mobile.VNUM, err)
		level = 1
	}
	mobile.Level = level
	log.Printf("Mobile #%d level: %d", mobile.VNUM, level)

	// Parse hitroll (THAC0)
	if len(statsParts) > 1 {
		hitroll, err := strconv.Atoi(statsParts[1])
		if err != nil {
			log.Printf("Warning: invalid hitroll for mobile #%d: %v, using default value", mobile.VNUM, err)
			// Use default hitroll based on level
			hitroll = mobile.Level
			if hitroll < 0 {
				hitroll = 0
			}
		}
		// Convert THAC0 to hitroll
		mobile.HitRoll = 20 - hitroll
		log.Printf("Mobile #%d hitroll: %d (from THAC0: %d)", mobile.VNUM, mobile.HitRoll, hitroll)
	} else {
		// Use a default hitroll based on level
		if mobile.Level > 3 {
			mobile.HitRoll = mobile.Level - 3
		} else {
			mobile.HitRoll = 1
		}
	}

	// Parse armor class
	if len(statsParts) > 2 {
		ac, err := strconv.Atoi(statsParts[2])
		if err != nil {
			log.Printf("Warning: invalid armor class for mobile #%d: %v, using default value", mobile.VNUM, err)
			ac = 10 // Default AC value
		}
		// Store AC directly for all positions
		mobile.AC = [3]int{ac, ac, ac} // Same AC for all positions
		log.Printf("Mobile #%d AC: %d", mobile.VNUM, ac)
	} else {
		// Default AC is 10 in the original DikuMUD
		mobile.AC = [3]int{10, 10, 10}
	}

	// Parse hit dice and damage dice
	if len(statsParts) > 3 {
		// Parse hit dice (for max hit points)
		hitDiceStr := statsParts[3]
		hitDiceParts := strings.Split(hitDiceStr, "d")
		if len(hitDiceParts) != 2 {
			log.Printf("Warning: invalid hit dice format for mobile #%d: %s, using default values", mobile.VNUM, hitDiceStr)
			// Default hit dice is 1d6+0
			mobile.Dice[0] = 1 // Number of dice
			mobile.Dice[1] = 6 // Size of dice
			mobile.Dice[2] = 0 // Bonus
		} else {
			// Parse number of dice
			numDice, err := strconv.Atoi(hitDiceParts[0])
			if err != nil {
				log.Printf("Warning: invalid number of hit dice for mobile #%d: %v, using default value", mobile.VNUM, err)
				numDice = 1
			}

			// Parse size of dice and bonus
			sizeBonusParts := strings.Split(hitDiceParts[1], "+")
			sizeDice, err := strconv.Atoi(sizeBonusParts[0])
			if err != nil {
				log.Printf("Warning: invalid size of hit dice for mobile #%d: %v, using default value", mobile.VNUM, err)
				sizeDice = 6
			}

			var hitBonus int
			if len(sizeBonusParts) > 1 {
				hitBonus, err = strconv.Atoi(sizeBonusParts[1])
				if err != nil {
					log.Printf("Warning: invalid hit bonus for mobile #%d: %v, using default value", mobile.VNUM, err)
					hitBonus = 0
				}
			}

			// Store the hit dice values exactly as in the original DikuMUD
			// These will be used to roll actual HP when creating the mob
			mobile.Dice[0] = numDice  // Number of dice
			mobile.Dice[1] = sizeDice // Size of dice
			mobile.Dice[2] = hitBonus // Bonus
		}
		log.Printf("Mobile #%d hit dice: %dd%d+%d", mobile.VNUM, mobile.Dice[0], mobile.Dice[1], mobile.Dice[2])
	} else {
		// Default hit dice is 1d6+0
		mobile.Dice[0] = 1 // Number of dice
		mobile.Dice[1] = 6 // Size of dice
		mobile.Dice[2] = 0 // Bonus
	}

	// Parse damage dice
	if len(statsParts) > 4 {
		damDiceStr := statsParts[4]
		damDiceParts := strings.Split(damDiceStr, "d")
		if len(damDiceParts) != 2 {
			log.Printf("Warning: invalid damage dice format for mobile #%d: %s, using default values", mobile.VNUM, damDiceStr)
			// Default damage dice is 1d4+0
			mobile.DamageType = 1
			mobile.AttackType = 4 // punch
			mobile.DamRoll = 0
		} else {
			// Parse number of dice
			damNumDice, err := strconv.Atoi(damDiceParts[0])
			if err != nil {
				log.Printf("Warning: invalid number of damage dice for mobile #%d: %v, using default value", mobile.VNUM, err)
				damNumDice = 1
			}
			// In the original DikuMUD, this is stored in damnodice
			mobile.DamageType = damNumDice

			// Parse size of dice and bonus
			damSizeBonusParts := strings.Split(damDiceParts[1], "+")
			damSizeDice, err := strconv.Atoi(damSizeBonusParts[0])
			if err != nil {
				log.Printf("Warning: invalid size of damage dice for mobile #%d: %v, using default value", mobile.VNUM, err)
				damSizeDice = 4
			}
			// In the original DikuMUD, this is stored in damsizedice
			mobile.AttackType = damSizeDice

			var damBonus int
			if len(damSizeBonusParts) > 1 {
				damBonus, err = strconv.Atoi(damSizeBonusParts[1])
				if err != nil {
					log.Printf("Warning: invalid damage bonus for mobile #%d: %v, using default value", mobile.VNUM, err)
					damBonus = 0
				}
			}

			// Set damage roll to the bonus value, exactly as in the original DikuMUD
			// This is the "strength bonus" that applies to all damage
			mobile.DamRoll = damBonus
		}
		log.Printf("Mobile #%d damage dice: %dd%d+%d", mobile.VNUM, mobile.DamageType, mobile.AttackType, mobile.DamRoll)
	} else {
		// Default damage dice is 1d4+0
		mobile.DamageType = 1
		mobile.AttackType = 4 // punch
		mobile.DamRoll = 0
	}

	// Get the next line for gold and experience
	if !parser.NextLine() {
		log.Printf("Warning: unexpected end of file while parsing gold and experience for mobile #%d, using default values", mobile.VNUM)
		// Default gold is level * 10
		mobile.Gold = mobile.Level * 10
		// Default experience is level * 100
		mobile.Experience = mobile.Level * 100
		// Default position is 8 (Standing)
		mobile.Position = 8
		// Default default position is 8 (Standing)
		mobile.DefaultPos = 8
		// Default sex is 0 (Neutral)
		mobile.Sex = 0
		// Set default class to 0 (as in the original code)
		mobile.Class = 0
		return nil
	}

	// Parse gold and experience
	goldExpLine := strings.TrimSpace(parser.Line())
	goldExpParts := strings.Fields(goldExpLine)

	// Parse gold
	if len(goldExpParts) > 0 {
		gold, err := strconv.Atoi(goldExpParts[0])
		if err != nil {
			log.Printf("Warning: invalid gold for mobile #%d: %v, using default value", mobile.VNUM, err)
			gold = mobile.Level * 10
		}
		mobile.Gold = gold
		log.Printf("Mobile #%d gold: %d", mobile.VNUM, gold)
	} else {
		// Default gold is level * 10
		mobile.Gold = mobile.Level * 10
	}

	// Parse experience
	if len(goldExpParts) > 1 {
		exp, err := strconv.Atoi(goldExpParts[1])
		if err != nil {
			log.Printf("Warning: invalid experience for mobile #%d: %v, using default value", mobile.VNUM, err)
			exp = mobile.Level * 100
		}
		mobile.Experience = exp
		log.Printf("Mobile #%d experience: %d", mobile.VNUM, exp)
	} else {
		// Default experience is level * 100
		mobile.Experience = mobile.Level * 100
	}

	// Get the next line for position, default position, and sex
	if !parser.NextLine() {
		log.Printf("Warning: unexpected end of file while parsing position, default position, and sex for mobile #%d, using default values", mobile.VNUM)
		// Default position is 8 (Standing)
		mobile.Position = 8
		// Default default position is 8 (Standing)
		mobile.DefaultPos = 8
		// Default sex is 0 (Neutral)
		mobile.Sex = 0
		// Set default class to 0 (as in the original code)
		mobile.Class = 0
		return nil
	}

	// Parse position, default position, and sex
	posDefPosSexLine := strings.TrimSpace(parser.Line())
	posDefPosSexParts := strings.Fields(posDefPosSexLine)

	// Parse position
	if len(posDefPosSexParts) > 0 {
		pos, err := strconv.Atoi(posDefPosSexParts[0])
		if err != nil {
			log.Printf("Warning: invalid position for mobile #%d: %v, using default value", mobile.VNUM, err)
			pos = 8 // Standing
		}
		mobile.Position = pos
		log.Printf("Mobile #%d position: %d", mobile.VNUM, pos)
	} else {
		// Default position is 8 (Standing)
		mobile.Position = 8
	}

	// Parse default position
	if len(posDefPosSexParts) > 1 {
		defPos, err := strconv.Atoi(posDefPosSexParts[1])
		if err != nil {
			log.Printf("Warning: invalid default position for mobile #%d: %v, using default value", mobile.VNUM, err)
			defPos = 8 // Standing
		}
		mobile.DefaultPos = defPos
		log.Printf("Mobile #%d default position: %d", mobile.VNUM, defPos)
	} else {
		// Default default position is 8 (Standing)
		mobile.DefaultPos = 8
	}

	// Parse sex
	if len(posDefPosSexParts) > 2 {
		sex, err := strconv.Atoi(posDefPosSexParts[2])
		if err != nil {
			log.Printf("Warning: invalid sex for mobile #%d: %v, using default value", mobile.VNUM, err)
			sex = 0 // Neutral
		}
		mobile.Sex = sex
		log.Printf("Mobile #%d sex: %d", mobile.VNUM, sex)
	} else {
		// Default sex is 0 (Neutral)
		mobile.Sex = 0
	}

	// Set default class to 0 (as in the original code)
	mobile.Class = 0

	// Skip to the next mobile
	for parser.NextLine() {
		// Skip this line
		if strings.HasPrefix(parser.Line(), "#") {
			// We've reached a new mobile, so we need to back up
			parser.BackUp()
			break
		}
	}

	return nil
}
