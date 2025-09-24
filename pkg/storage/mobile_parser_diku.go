package storage

import (
	"log"
	"math"
	"strconv"
	"strings"

	"github.com/wltechblog/DikuGo/pkg/types"
)

// setDefaultMobValues sets default values for a mobile
func setDefaultMobValues(mobile *types.Mobile) error {
	// Set default values for the mobile
	mobile.Level = 1
	mobile.HitRoll = 0
	mobile.DamRoll = 0
	mobile.Dice = [3]int{1, 6, 0}
	mobile.AC = [3]int{10, 10, 10}
	mobile.Gold = 10
	mobile.Experience = 100
	mobile.Position = 8
	mobile.DefaultPos = 8
	mobile.Sex = 0

	// Skip to the next mobile
	return nil
}

// parseDikuMobile parses a mobile in the original DikuMUD format
func parseDikuMobile(parser *Parser, mobile *types.Mobile) error {
	log.Printf("Parsing DikuMUD mobile #%d", mobile.VNUM)
	// Set default abilities
	// In the original code, simple mobiles have default ability scores of 11
	mobile.Abilities = [6]int{11, 11, 11, 11, 11, 11} // STR, INT, WIS, DEX, CON, CHA

	// In the original DikuMUD, there are two formats for mobs:
	// 1. Simple mobs (marked with 'S') - These have default ability scores of 11
	// 2. Detailed mobs (not marked with 'S') - These have explicit ability scores
	//
	// The original C code in read_mobile() function in old/db.c handles this by checking
	// for the letter 'S' after the alignment. If 'S' is found, it's a simple mob.
	// If not, it's a detailed mob with explicit ability scores.

	// Parse the act flags, affect flags, and alignment
	if !parser.NextLine() {
		log.Printf("Warning: unexpected end of file while parsing flags for mobile #%d, using default values", mobile.VNUM)
		mobile.ActFlags = 8 // Default to NPC flag
		mobile.AffectFlags = 0
		mobile.Alignment = 0
		return setDefaultMobValues(mobile)
	}

	// Parse the flags line: ActFlags AffectFlags Alignment S/E/D
	flagsLine := strings.TrimSpace(parser.Line())
	flagsParts := strings.Fields(flagsLine)
	if len(flagsParts) < 3 {
		log.Printf("Warning: invalid flags line on line %d for mobile #%d: %s, using default values", parser.LineNum(), mobile.VNUM, flagsLine)
		mobile.ActFlags = 8 // Default to NPC flag
		mobile.AffectFlags = 0
		mobile.Alignment = 0
		return setDefaultMobValues(mobile)
	}

	// Parse act flags
	actFlags, err := strconv.Atoi(flagsParts[0])
	if err != nil {
		log.Printf("Warning: invalid act flags on line %d for mobile #%d: %v, using default value", parser.LineNum(), mobile.VNUM, err)
		actFlags = 8 // Default to NPC flag
	}
	mobile.ActFlags = uint32(actFlags)

	// Parse affect flags
	affectFlags, err := strconv.Atoi(flagsParts[1])
	if err != nil {
		log.Printf("Warning: invalid affect flags on line %d for mobile #%d: %v, using default value", parser.LineNum(), mobile.VNUM, err)
		affectFlags = 0
	}
	mobile.AffectFlags = uint32(affectFlags)

	// Parse alignment
	alignment, err := strconv.Atoi(flagsParts[2])
	if err != nil {
		log.Printf("Warning: invalid alignment on line %d for mobile #%d: %v, using default value", parser.LineNum(), mobile.VNUM, err)
		alignment = 0
	}
	mobile.Alignment = alignment

	// Check if this is a simple or detailed mob
	var isSimpleMob bool
	if len(flagsParts) >= 4 {
		isSimpleMob = flagsParts[3] == "S"
	} else {
		// Default to simple mob if not specified
		isSimpleMob = true
	}

	// Debug logging
	log.Printf("DEBUG: Mobile #%d format type: '%s', flags: %v", mobile.VNUM, map[bool]string{true: "S", false: flagsParts[0]}[isSimpleMob], flagsParts)
	if isSimpleMob {
		log.Printf("DEBUG: Using simple parser for mobile #%d", mobile.VNUM)
		log.Printf("Parsing simple mobile #%d", mobile.VNUM)
		log.Printf("DEBUG: Starting to parse simple mobile #%d with line: %s", mobile.VNUM, strings.Join(flagsParts, " "))
	} else {
		log.Printf("DEBUG: Using DikuMUD parser for mobile #%d", mobile.VNUM)
	}

	log.Printf("Mobile #%d is a %s mob", mobile.VNUM, map[bool]string{true: "simple", false: "detailed"}[isSimpleMob])

	// If this is a detailed mob, parse the ability scores
	if !isSimpleMob {
		// Parse strength
		if !parser.NextLine() {
			log.Printf("Warning: unexpected end of file while parsing abilities for mobile #%d, using default values", mobile.VNUM)
			return setDefaultMobValues(mobile)
		}

		str, err := strconv.Atoi(strings.TrimSpace(parser.Line()))
		if err != nil {
			log.Printf("Warning: invalid strength on line %d for mobile #%d: %v, using default value", parser.LineNum(), mobile.VNUM, err)
			str = 11
		}
		mobile.Abilities[0] = str

		// Parse intelligence
		if !parser.NextLine() {
			log.Printf("Warning: unexpected end of file while parsing abilities for mobile #%d, using default values", mobile.VNUM)
			return setDefaultMobValues(mobile)
		}

		intel, err := strconv.Atoi(strings.TrimSpace(parser.Line()))
		if err != nil {
			log.Printf("Warning: invalid intelligence on line %d for mobile #%d: %v, using default value", parser.LineNum(), mobile.VNUM, err)
			intel = 11
		}
		mobile.Abilities[1] = intel

		// Parse wisdom
		if !parser.NextLine() {
			log.Printf("Warning: unexpected end of file while parsing abilities for mobile #%d, using default values", mobile.VNUM)
			return setDefaultMobValues(mobile)
		}

		wis, err := strconv.Atoi(strings.TrimSpace(parser.Line()))
		if err != nil {
			log.Printf("Warning: invalid wisdom on line %d for mobile #%d: %v, using default value", parser.LineNum(), mobile.VNUM, err)
			wis = 11
		}
		mobile.Abilities[2] = wis

		// Parse dexterity
		if !parser.NextLine() {
			log.Printf("Warning: unexpected end of file while parsing abilities for mobile #%d, using default values", mobile.VNUM)
			return setDefaultMobValues(mobile)
		}

		dex, err := strconv.Atoi(strings.TrimSpace(parser.Line()))
		if err != nil {
			log.Printf("Warning: invalid dexterity on line %d for mobile #%d: %v, using default value", parser.LineNum(), mobile.VNUM, err)
			dex = 11
		}
		mobile.Abilities[3] = dex

		// Parse constitution
		if !parser.NextLine() {
			log.Printf("Warning: unexpected end of file while parsing abilities for mobile #%d, using default values", mobile.VNUM)
			return setDefaultMobValues(mobile)
		}

		con, err := strconv.Atoi(strings.TrimSpace(parser.Line()))
		if err != nil {
			log.Printf("Warning: invalid constitution on line %d for mobile #%d: %v, using default value", parser.LineNum(), mobile.VNUM, err)
			con = 11
		}
		mobile.Abilities[4] = con
	}

	// Parse the level, hitroll, damroll, hit dice, and damage dice
	if !parser.NextLine() {
		log.Printf("Warning: unexpected end of file while parsing level/dice for mobile #%d, using default values", mobile.VNUM)
		return setDefaultMobValues(mobile)
	}

	// Parse differently based on mob type
	if isSimpleMob {
		// For simple mobs, the format is: Level HitRoll AC HitDice DamageDice
		levelLine := strings.TrimSpace(parser.Line())
		levelParts := strings.Fields(levelLine)
		if len(levelParts) < 5 {
			// Log warning but continue with default values
			log.Printf("Warning: invalid level line on line %d: %s, using default values", parser.LineNum(), levelLine)
			mobile.Level = 1
			mobile.HitRoll = 0
			mobile.DamRoll = 0
			mobile.Dice = [3]int{1, 6, 0}
		} else {
			// Parse level
			level, err := strconv.Atoi(levelParts[0])
			if err != nil {
				// Log warning but continue with default value
				log.Printf("Warning: invalid level on line %d: %v, using default value", parser.LineNum(), err)
				level = 1
			}
			mobile.Level = level
			log.Printf("Mobile #%d level: %d", mobile.VNUM, level)

			// Parse hitroll directly
			hitroll, err := strconv.Atoi(levelParts[1])
			if err != nil {
				// Log warning but continue with default value
				log.Printf("Warning: invalid hitroll on line %d: %v, using default value", parser.LineNum(), err)
				hitroll = 0
			}
			// Convert THAC0 to hitroll
			mobile.HitRoll = 20 - hitroll
			log.Printf("Mobile #%d hitroll: %d (from THAC0: %d)", mobile.VNUM, mobile.HitRoll, hitroll)

			// Parse AC (Armor Class)
			ac, err := strconv.Atoi(levelParts[2])
			if err != nil {
				// Log warning but continue with default value
				log.Printf("Warning: invalid AC on line %d: %v, using default value", parser.LineNum(), err)
				ac = 10
			}
			// Store AC directly
			mobile.AC = [3]int{ac, ac, ac} // Same AC for all positions
			log.Printf("Mobile #%d AC: %d", mobile.VNUM, ac)

			// Parse hit dice
			hitDiceStr := levelParts[3]
			hitDiceParts := strings.Split(hitDiceStr, "d")
			if len(hitDiceParts) != 2 {
				// Log warning but continue with default values
				log.Printf("Warning: invalid hit dice format on line %d: %s, using default values", parser.LineNum(), hitDiceStr)
				mobile.Dice[0] = 1 // Number of dice
				mobile.Dice[1] = 6 // Size of dice
				mobile.Dice[2] = 0 // Bonus
			} else {
				// Parse number of dice
				numDice, err := strconv.Atoi(hitDiceParts[0])
				if err != nil {
					// Log warning but continue with default value
					log.Printf("Warning: invalid number of hit dice on line %d: %v, using default value", parser.LineNum(), err)
					numDice = 1
				}

				// Parse size of dice and bonus
				sizeBonusParts := strings.Split(hitDiceParts[1], "+")
				sizeDice, err := strconv.Atoi(sizeBonusParts[0])
				if err != nil {
					// Log warning but continue with default value
					log.Printf("Warning: invalid size of hit dice on line %d: %v, using default value", parser.LineNum(), err)
					sizeDice = 6
				}

				var hitBonus int
				if len(sizeBonusParts) > 1 {
					hitBonus, err = strconv.Atoi(sizeBonusParts[1])
					if err != nil {
						// Log warning but continue with default value
						log.Printf("Warning: invalid hit bonus on line %d: %v, using default value", parser.LineNum(), err)
						hitBonus = 0
					}
				}

				// Store the hit dice values exactly as in the original DikuMUD
				// These will be used to roll actual HP when creating the mob
				mobile.Dice[0] = numDice  // Number of dice
				mobile.Dice[1] = sizeDice // Size of dice
				mobile.Dice[2] = hitBonus // Bonus
			}

			// Parse damage dice
			damDiceStr := levelParts[4]
			damDiceParts := strings.Split(damDiceStr, "d")
			if len(damDiceParts) != 2 {
				// Log warning but continue with default values
				log.Printf("Warning: invalid damage dice format on line %d: %s, using default values", parser.LineNum(), damDiceStr)
				mobile.DamageType = 1
				mobile.AttackType = 4 // punch
				mobile.DamRoll = 0
			} else {
				// Parse number of dice
				numDice, err := strconv.Atoi(damDiceParts[0])
				if err != nil {
					// Log warning but continue with default value
					log.Printf("Warning: invalid number of damage dice on line %d: %v, using default value", parser.LineNum(), err)
					numDice = 1
				}
				mobile.DamageType = numDice

				// Parse size of dice and bonus
				sizeBonusParts := strings.Split(damDiceParts[1], "+")
				sizeDice, err := strconv.Atoi(sizeBonusParts[0])
				if err != nil {
					// Log warning but continue with default value
					log.Printf("Warning: invalid size of damage dice on line %d: %v, using default value", parser.LineNum(), err)
					sizeDice = 4
				}
				mobile.AttackType = sizeDice

				// Parse damage bonus (damroll)
				if len(sizeBonusParts) > 1 {
					damBonus, err := strconv.Atoi(sizeBonusParts[1])
					if err != nil {
						// Log warning but continue with default value
						log.Printf("Warning: invalid damage bonus on line %d: %v, using default value", parser.LineNum(), err)
						damBonus = 0
					}
					mobile.DamRoll = damBonus
				} else {
					mobile.DamRoll = 0
				}
				log.Printf("Mobile #%d damage dice: %dd%d+%d", mobile.VNUM, mobile.DamageType, mobile.AttackType, mobile.DamRoll)
			}
		}
	} else {
		// For detailed mobs, the format is different
		// First line is hit points: min max
		hitLine := strings.TrimSpace(parser.Line())
		hitParts := strings.Fields(hitLine)
		if len(hitParts) < 2 {
			log.Printf("Warning: invalid hit points line on line %d for mobile #%d: %s, using default values", parser.LineNum(), mobile.VNUM, hitLine)
			mobile.Dice = [3]int{1, 6, 0}
		} else {
			// Parse hit points: min max
			min, err := strconv.Atoi(hitParts[0])
			if err != nil {
				log.Printf("Warning: invalid min hit points on line %d for mobile #%d: %v, using default value", parser.LineNum(), mobile.VNUM, err)
				min = 1
			}

			max, err := strconv.Atoi(hitParts[1])
			if err != nil {
				log.Printf("Warning: invalid max hit points on line %d for mobile #%d: %v, using default value", parser.LineNum(), mobile.VNUM, err)
				max = 10
			}

			// Set hit dice to approximate the hit points
			avg := (min + max) / 2
			mobile.Dice = [3]int{avg / 4, 4, 0}
			log.Printf("Mobile #%d hit points: %d-%d (avg: %d)", mobile.VNUM, min, max, avg)
		}

		// Parse AC
		if !parser.NextLine() {
			log.Printf("Warning: unexpected end of file while parsing AC for mobile #%d, using default values", mobile.VNUM)
			return setDefaultMobValues(mobile)
		}

		ac, err := strconv.Atoi(strings.TrimSpace(parser.Line()))
		if err != nil {
			log.Printf("Warning: invalid AC on line %d for mobile #%d: %v, using default value", parser.LineNum(), mobile.VNUM, err)
			ac = 10
		}
		mobile.AC = [3]int{ac, ac, ac}
		log.Printf("Mobile #%d AC: %d", mobile.VNUM, ac)

		// Skip the next 2 lines (mana and move)
		for i := 0; i < 2; i++ {
			if !parser.NextLine() {
				log.Printf("Warning: unexpected end of file while parsing mana/move for mobile #%d, using default values", mobile.VNUM)
				return setDefaultMobValues(mobile)
			}
		}

		// Parse gold
		if !parser.NextLine() {
			log.Printf("Warning: unexpected end of file while parsing gold for mobile #%d, using default values", mobile.VNUM)
			return setDefaultMobValues(mobile)
		}

		gold, err := strconv.Atoi(strings.TrimSpace(parser.Line()))
		if err != nil {
			log.Printf("Warning: invalid gold on line %d for mobile #%d: %v, using default value", parser.LineNum(), mobile.VNUM, err)
			gold = 0
		}
		mobile.Gold = gold
		log.Printf("Mobile #%d gold: %d", mobile.VNUM, gold)

		// Parse experience
		if !parser.NextLine() {
			log.Printf("Warning: unexpected end of file while parsing experience for mobile #%d, using default values", mobile.VNUM)
			return setDefaultMobValues(mobile)
		}

		exp, err := strconv.Atoi(strings.TrimSpace(parser.Line()))
		if err != nil {
			log.Printf("Warning: invalid experience on line %d for mobile #%d: %v, using default value", parser.LineNum(), mobile.VNUM, err)
			exp = 0
		}
		mobile.Experience = exp
		log.Printf("Mobile #%d experience: %d", mobile.VNUM, exp)

		// Parse position
		if !parser.NextLine() {
			log.Printf("Warning: unexpected end of file while parsing position for mobile #%d, using default values", mobile.VNUM)
			return setDefaultMobValues(mobile)
		}

		position, err := strconv.Atoi(strings.TrimSpace(parser.Line()))
		if err != nil {
			log.Printf("Warning: invalid position on line %d for mobile #%d: %v, using default value", parser.LineNum(), mobile.VNUM, err)
			position = 8 // Standing
		}
		mobile.Position = position
		log.Printf("Mobile #%d position: %d", mobile.VNUM, position)

		// Parse default position
		if !parser.NextLine() {
			log.Printf("Warning: unexpected end of file while parsing default position for mobile #%d, using default values", mobile.VNUM)
			return setDefaultMobValues(mobile)
		}

		defaultPos, err := strconv.Atoi(strings.TrimSpace(parser.Line()))
		if err != nil {
			log.Printf("Warning: invalid default position on line %d for mobile #%d: %v, using default value", parser.LineNum(), mobile.VNUM, err)
			defaultPos = 8 // Standing
		}
		mobile.DefaultPos = defaultPos
		log.Printf("Mobile #%d default position: %d", mobile.VNUM, defaultPos)

		// Parse sex
		if !parser.NextLine() {
			log.Printf("Warning: unexpected end of file while parsing sex for mobile #%d, using default values", mobile.VNUM)
			return setDefaultMobValues(mobile)
		}

		sex, err := strconv.Atoi(strings.TrimSpace(parser.Line()))
		if err != nil {
			log.Printf("Warning: invalid sex on line %d for mobile #%d: %v, using default value", parser.LineNum(), mobile.VNUM, err)
			sex = 0 // Neutral
		}
		mobile.Sex = sex
		log.Printf("Mobile #%d sex: %d", mobile.VNUM, sex)

		// Set level based on experience (approximation)
		mobile.Level = 1 // Default
		if mobile.Experience > 0 {
			// Rough approximation of level based on experience
			mobile.Level = int(math.Sqrt(float64(mobile.Experience) / 1000))
			if mobile.Level < 1 {
				mobile.Level = 1
			}
		}
		log.Printf("Mobile #%d level (approximated): %d", mobile.VNUM, mobile.Level)

		// Set hitroll and damroll based on level (approximation)
		mobile.HitRoll = mobile.Level / 2
		mobile.DamRoll = mobile.Level / 3
		log.Printf("Mobile #%d hitroll (approximated): %d", mobile.VNUM, mobile.HitRoll)
		log.Printf("Mobile #%d damroll (approximated): %d", mobile.VNUM, mobile.DamRoll)
	}

	// For simple mobs, we need to parse gold/exp and position/sex
	if isSimpleMob {
		// Parse the gold and experience
		if !parser.NextLine() {
			log.Printf("Warning: unexpected end of file while parsing gold/exp for mobile #%d, using default values", mobile.VNUM)
			return setDefaultMobValues(mobile)
		}

		// Parse the gold/exp line: Gold Experience
		goldExpLine := strings.TrimSpace(parser.Line())
		goldExpParts := strings.Fields(goldExpLine)
		if len(goldExpParts) < 2 {
			// Log warning but continue with default values
			log.Printf("Warning: invalid gold/exp line on line %d: %s, using default values", parser.LineNum(), goldExpLine)
			mobile.Gold = 0
			mobile.Experience = 0
		} else {
			// Parse gold
			gold, err := strconv.Atoi(goldExpParts[0])
			if err != nil {
				// Log warning but continue with default value
				log.Printf("Warning: invalid gold on line %d: %v, using default value", parser.LineNum(), err)
				gold = 0
			}
			mobile.Gold = gold
			log.Printf("Mobile #%d gold: %d", mobile.VNUM, gold)

			// Parse experience
			exp, err := strconv.Atoi(goldExpParts[1])
			if err != nil {
				// Log warning but continue with default value
				log.Printf("Warning: invalid experience on line %d: %v, using default value", parser.LineNum(), err)
				exp = 0
			}
			mobile.Experience = exp
			log.Printf("Mobile #%d experience: %d", mobile.VNUM, exp)
		}

		// Parse the position, default position, and sex
		if !parser.NextLine() {
			log.Printf("Warning: unexpected end of file while parsing position/sex for mobile #%d, using default values", mobile.VNUM)
			return setDefaultMobValues(mobile)
		}

		// Parse the position/sex line: Position DefaultPosition Sex
		positionSexLine := strings.TrimSpace(parser.Line())
		positionSexParts := strings.Fields(positionSexLine)
		if len(positionSexParts) < 3 {
			// Log warning but continue with default values
			log.Printf("Warning: invalid position/sex line on line %d: %s, using default values", parser.LineNum(), positionSexLine)
			mobile.Position = 8   // Standing
			mobile.DefaultPos = 8 // Standing
			mobile.Sex = 0        // Neutral
		} else {
			// Parse position
			position, err := strconv.Atoi(positionSexParts[0])
			if err != nil {
				// Log warning but continue with default value
				log.Printf("Warning: invalid position on line %d: %v, using default value", parser.LineNum(), err)
				position = 8 // Standing
			}
			mobile.Position = position

			// Parse default position
			defaultPos, err := strconv.Atoi(positionSexParts[1])
			if err != nil {
				// Log warning but continue with default value
				log.Printf("Warning: invalid default position on line %d: %v, using default value", parser.LineNum(), err)
				defaultPos = 8 // Standing
			}
			mobile.DefaultPos = defaultPos

			// Parse sex
			sex, err := strconv.Atoi(positionSexParts[2])
			if err != nil {
				// Log warning but continue with default value
				log.Printf("Warning: invalid sex on line %d: %v, using default value", parser.LineNum(), err)
				sex = 0 // Neutral
			}
			mobile.Sex = sex
		}
	}

	// Set default class to 0 (as in the original code)
	mobile.Class = 0

	// If AC hasn't been set yet, use default value
	if mobile.AC[0] == 0 && mobile.AC[1] == 0 && mobile.AC[2] == 0 {
		// Default AC is 100 in the original DikuMUD
		mobile.AC = [3]int{100, 100, 100} // Same AC for all positions
		log.Printf("Mobile #%d using default AC: %d", mobile.VNUM, 100)
	}

	// Log the final mobile stats
	log.Printf("Parsed DikuMUD mobile #%d with stats: Level=%d, HitRoll=%d, DamRoll=%d, AC=%v, Gold=%d, Exp=%d",
		mobile.VNUM, mobile.Level, mobile.HitRoll, mobile.DamRoll, mobile.AC, mobile.Gold, mobile.Experience)

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
