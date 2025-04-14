package storage

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/wltechblog/DikuGo/pkg/types"
)

// parseSimpleMobile parses a simple mobile (type 'S')
func parseSimpleMobile(parser *Parser, mobile *types.Mobile) error {
	// Set default abilities
	// In the original code, simple mobiles have default ability scores of 11
	mobile.Abilities = [6]int{11, 11, 11, 11, 11, 11} // STR, INT, WIS, DEX, CON, CHA

	// Parse level
	if !parser.NextLine() {
		// Set default values and return
		setDefaultMobileValues(mobile)
		return nil
	}
	levelStr := strings.TrimSpace(parser.Line())
	level, err := strconv.Atoi(levelStr)
	if err != nil {
		// Log warning but continue with default value
		fmt.Printf("Warning: invalid level on line %d: %v, using default value\n", parser.LineNum(), err)
		level = 1
	}
	mobile.Level = level

	// Parse hitroll (THAC0)
	if !parser.NextLine() {
		// Set default values and return
		setDefaultMobileValues(mobile)
		return nil
	}
	hitrollStr := strings.TrimSpace(parser.Line())
	hitroll, err := strconv.Atoi(hitrollStr)
	if err != nil {
		// Log warning but continue with default value
		fmt.Printf("Warning: invalid hitroll on line %d: %v, using default value\n", parser.LineNum(), err)
		hitroll = 20
	}
	mobile.HitRoll = 20 - hitroll // Convert THAC0 to hitroll

	// Parse armor class
	if !parser.NextLine() {
		// Set default values and return
		setDefaultMobileValues(mobile)
		return nil
	}
	acStr := strings.TrimSpace(parser.Line())
	ac, err := strconv.Atoi(acStr)
	if err != nil {
		// Log warning but continue with default value
		fmt.Printf("Warning: invalid armor class on line %d: %v, using default value\n", parser.LineNum(), err)
		ac = 10
	}
	mobile.AC = [3]int{10 * ac, 10 * ac, 10 * ac} // Same AC for all positions

	// Parse hit dice (for max hit points)
	if !parser.NextLine() {
		// Set default values and return
		setDefaultMobileValues(mobile)
		return nil
	}
	hitDiceStr := strings.TrimSpace(parser.Line())
	hitDiceParts := strings.Split(hitDiceStr, "d")
	if len(hitDiceParts) != 2 {
		// Log warning but continue with default values
		fmt.Printf("Warning: invalid hit dice format on line %d: %s, using default values\n", parser.LineNum(), hitDiceStr)
		mobile.Dice[0] = 1 // Number of dice
		mobile.Dice[1] = 6 // Size of dice
		mobile.Dice[2] = 0 // Bonus
	} else {
		// Parse number of dice
		numDice, err := strconv.Atoi(hitDiceParts[0])
		if err != nil {
			// Log warning but continue with default value
			fmt.Printf("Warning: invalid number of hit dice on line %d: %v, using default value\n", parser.LineNum(), err)
			numDice = 1
		}

		// Parse size of dice and bonus
		sizeBonusParts := strings.Split(hitDiceParts[1], "+")
		sizeDice, err := strconv.Atoi(sizeBonusParts[0])
		if err != nil {
			// Log warning but continue with default value
			fmt.Printf("Warning: invalid size of hit dice on line %d: %v, using default value\n", parser.LineNum(), err)
			sizeDice = 6
		}

		var hitBonus int
		if len(sizeBonusParts) > 1 {
			hitBonus, err = strconv.Atoi(sizeBonusParts[1])
			if err != nil {
				// Log warning but continue with default value
				fmt.Printf("Warning: invalid hit bonus on line %d: %v, using default value\n", parser.LineNum(), err)
				hitBonus = 0
			}
		}

		// Store the hit dice values
		mobile.Dice[0] = numDice  // Number of dice
		mobile.Dice[1] = sizeDice // Size of dice
		mobile.Dice[2] = hitBonus // Bonus
	}

	// Parse damage dice
	if !parser.NextLine() {
		// Set default values and return
		setDefaultMobileValues(mobile)
		return nil
	}
	damDiceStr := strings.TrimSpace(parser.Line())
	damDiceParts := strings.Split(damDiceStr, "d")
	if len(damDiceParts) != 2 {
		// Log warning but continue with default values
		fmt.Printf("Warning: invalid damage dice format on line %d: %s, using default values\n", parser.LineNum(), damDiceStr)
		mobile.DamageType = 1
		mobile.AttackType = 4 // punch
		mobile.DamRoll = 0
	} else {
		// Parse number of dice
		numDice, err := strconv.Atoi(damDiceParts[0])
		if err != nil {
			// Log warning but continue with default value
			fmt.Printf("Warning: invalid number of damage dice on line %d: %v, using default value\n", parser.LineNum(), err)
			numDice = 1
		}
		mobile.DamageType = numDice

		// Parse size of dice and bonus
		sizeBonusParts := strings.Split(damDiceParts[1], "+")
		sizeDice, err := strconv.Atoi(sizeBonusParts[0])
		if err != nil {
			// Log warning but continue with default value
			fmt.Printf("Warning: invalid size of damage dice on line %d: %v, using default value\n", parser.LineNum(), err)
			sizeDice = 4
		}
		mobile.AttackType = sizeDice

		var damBonus int
		if len(sizeBonusParts) > 1 {
			damBonus, err = strconv.Atoi(sizeBonusParts[1])
			if err != nil {
				// Log warning but continue with default value
				fmt.Printf("Warning: invalid damage bonus on line %d: %v, using default value\n", parser.LineNum(), err)
				damBonus = 0
			}
		}

		// Set damage roll
		mobile.DamRoll = damBonus
	}

	// Parse gold
	if !parser.NextLine() {
		// Set default values and return
		setDefaultMobileValues(mobile)
		return nil
	}
	goldStr := strings.TrimSpace(parser.Line())
	gold, err := strconv.Atoi(goldStr)
	if err != nil {
		// Log warning but continue with default value
		fmt.Printf("Warning: invalid gold on line %d: %v, using default value\n", parser.LineNum(), err)
		gold = 0
	}
	mobile.Gold = gold

	// Parse experience
	if !parser.NextLine() {
		// Set default values and return
		setDefaultMobileValues(mobile)
		return nil
	}
	expStr := strings.TrimSpace(parser.Line())
	exp, err := strconv.Atoi(expStr)
	if err != nil {
		// Log warning but continue with default value
		fmt.Printf("Warning: invalid experience on line %d: %v, using default value\n", parser.LineNum(), err)
		exp = 0
	}
	mobile.Experience = exp

	// Parse position
	if !parser.NextLine() {
		// Set default values and return
		setDefaultMobileValues(mobile)
		return nil
	}
	posStr := strings.TrimSpace(parser.Line())
	pos, err := strconv.Atoi(posStr)
	if err != nil {
		// Log warning but continue with default value
		fmt.Printf("Warning: invalid position on line %d: %v, using default value\n", parser.LineNum(), err)
		pos = 8 // Standing
	}
	mobile.Position = pos

	// Parse default position
	if !parser.NextLine() {
		// Set default values and return
		setDefaultMobileValues(mobile)
		return nil
	}
	defPosStr := strings.TrimSpace(parser.Line())
	defPos, err := strconv.Atoi(defPosStr)
	if err != nil {
		// Log warning but continue with default value
		fmt.Printf("Warning: invalid default position on line %d: %v, using default value\n", parser.LineNum(), err)
		defPos = 8 // Standing
	}
	mobile.DefaultPos = defPos

	// Parse sex
	if !parser.NextLine() {
		// Set default values and return
		setDefaultMobileValues(mobile)
		return nil
	}
	sexStr := strings.TrimSpace(parser.Line())
	sex, err := strconv.Atoi(sexStr)
	if err != nil {
		// Log warning but continue with default value
		fmt.Printf("Warning: invalid sex on line %d: %v, using default value\n", parser.LineNum(), err)
		sex = 0 // Neutral
	}
	mobile.Sex = sex

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

// setDefaultMobileValues sets default values for a mobile
func setDefaultMobileValues(mobile *types.Mobile) {
	// Set default values for a mobile
	if mobile.Level == 0 {
		mobile.Level = 1
	}
	if mobile.HitRoll == 0 {
		mobile.HitRoll = 0
	}
	if mobile.AC[0] == 0 && mobile.AC[1] == 0 && mobile.AC[2] == 0 {
		mobile.AC = [3]int{100, 100, 100}
	}
	if mobile.Dice[0] == 0 && mobile.Dice[1] == 0 {
		mobile.Dice[0] = 1 // Number of dice
		mobile.Dice[1] = 6 // Size of dice
		mobile.Dice[2] = 0 // Bonus
	}
	if mobile.DamageType == 0 {
		mobile.DamageType = 1
	}
	if mobile.AttackType == 0 {
		mobile.AttackType = 4 // punch
	}
	if mobile.Position == 0 {
		mobile.Position = 8 // Standing
	}
	if mobile.DefaultPos == 0 {
		mobile.DefaultPos = 8 // Standing
	}
	if mobile.Sex == 0 {
		mobile.Sex = 0 // Neutral
	}
	mobile.Class = 0
}
