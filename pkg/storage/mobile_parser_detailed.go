package storage

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/wltechblog/DikuGo/pkg/types"
)

// parseDetailedMobile parses a detailed mobile (type 'D')
func parseDetailedMobile(parser *Parser, mobile *types.Mobile) error {
	// Parse strength
	if !parser.NextLine() {
		return fmt.Errorf("unexpected end of file after mobile type")
	}
	strStr := strings.TrimSpace(parser.Line())
	str, err := strconv.Atoi(strStr)
	if err != nil {
		return fmt.Errorf("invalid strength on line %d: %w", parser.LineNum(), err)
	}
	mobile.Abilities[0] = str // STR

	// Parse intelligence
	if !parser.NextLine() {
		return fmt.Errorf("unexpected end of file after strength")
	}
	intStr := strings.TrimSpace(parser.Line())
	intel, err := strconv.Atoi(intStr)
	if err != nil {
		return fmt.Errorf("invalid intelligence on line %d: %w", parser.LineNum(), err)
	}
	mobile.Abilities[1] = intel // INT

	// Parse wisdom
	if !parser.NextLine() {
		return fmt.Errorf("unexpected end of file after intelligence")
	}
	wisStr := strings.TrimSpace(parser.Line())
	wis, err := strconv.Atoi(wisStr)
	if err != nil {
		return fmt.Errorf("invalid wisdom on line %d: %w", parser.LineNum(), err)
	}
	mobile.Abilities[2] = wis // WIS

	// Parse dexterity
	if !parser.NextLine() {
		return fmt.Errorf("unexpected end of file after wisdom")
	}
	dexStr := strings.TrimSpace(parser.Line())
	dex, err := strconv.Atoi(dexStr)
	if err != nil {
		return fmt.Errorf("invalid dexterity on line %d: %w", parser.LineNum(), err)
	}
	mobile.Abilities[3] = dex // DEX

	// Parse constitution
	if !parser.NextLine() {
		return fmt.Errorf("unexpected end of file after dexterity")
	}
	conStr := strings.TrimSpace(parser.Line())
	con, err := strconv.Atoi(conStr)
	if err != nil {
		return fmt.Errorf("invalid constitution on line %d: %w", parser.LineNum(), err)
	}
	mobile.Abilities[4] = con // CON

	// Parse hit points (min and max)
	if !parser.NextLine() {
		return fmt.Errorf("unexpected end of file after constitution")
	}
	hpStr := strings.TrimSpace(parser.Line())
	hpParts := strings.Fields(hpStr)
	if len(hpParts) < 2 {
		return fmt.Errorf("invalid hit points format on line %d: %s", parser.LineNum(), hpStr)
	}

	// We don't actually use these values directly in the Mobile struct
	// In the original code, these are used to generate a random value
	// For now, we'll just store the average in the Dice field
	minHP, err := strconv.Atoi(hpParts[0])
	if err != nil {
		return fmt.Errorf("invalid min hit points on line %d: %w", parser.LineNum(), err)
	}

	maxHP, err := strconv.Atoi(hpParts[1])
	if err != nil {
		return fmt.Errorf("invalid max hit points on line %d: %w", parser.LineNum(), err)
	}

	// Store the average in the Dice field
	avgHP := (minHP + maxHP) / 2
	mobile.Dice[0] = 1     // Number of dice
	mobile.Dice[1] = 1     // Size of dice
	mobile.Dice[2] = avgHP // Bonus

	// Parse armor class
	if !parser.NextLine() {
		return fmt.Errorf("unexpected end of file after hit points")
	}
	acStr := strings.TrimSpace(parser.Line())
	ac, err := strconv.Atoi(acStr)
	if err != nil {
		return fmt.Errorf("invalid armor class on line %d: %w", parser.LineNum(), err)
	}
	mobile.AC = [3]int{10 * ac, 10 * ac, 10 * ac} // Same AC for all positions

	// Parse mana
	if !parser.NextLine() {
		return fmt.Errorf("unexpected end of file after armor class")
	}
	manaStr := strings.TrimSpace(parser.Line())
	_, err = strconv.Atoi(manaStr)
	if err != nil {
		return fmt.Errorf("invalid mana on line %d: %w", parser.LineNum(), err)
	}
	// We don't store mana directly in the Mobile struct, but we could add it if needed

	// Parse move points
	if !parser.NextLine() {
		return fmt.Errorf("unexpected end of file after mana")
	}
	moveStr := strings.TrimSpace(parser.Line())
	_, err = strconv.Atoi(moveStr)
	if err != nil {
		return fmt.Errorf("invalid move points on line %d: %w", parser.LineNum(), err)
	}
	// We don't store move points directly in the Mobile struct, but we could add it if needed

	// Parse gold
	if !parser.NextLine() {
		return fmt.Errorf("unexpected end of file after move points")
	}
	goldStr := strings.TrimSpace(parser.Line())
	gold, err := strconv.Atoi(goldStr)
	if err != nil {
		return fmt.Errorf("invalid gold on line %d: %w", parser.LineNum(), err)
	}
	mobile.Gold = gold

	// Parse experience
	if !parser.NextLine() {
		return fmt.Errorf("unexpected end of file after gold")
	}
	expStr := strings.TrimSpace(parser.Line())
	exp, err := strconv.Atoi(expStr)
	if err != nil {
		return fmt.Errorf("invalid experience on line %d: %w", parser.LineNum(), err)
	}
	mobile.Experience = exp

	// Parse position
	if !parser.NextLine() {
		return fmt.Errorf("unexpected end of file after experience")
	}
	posStr := strings.TrimSpace(parser.Line())
	pos, err := strconv.Atoi(posStr)
	if err != nil {
		return fmt.Errorf("invalid position on line %d: %w", parser.LineNum(), err)
	}
	mobile.Position = pos

	// Parse default position
	if !parser.NextLine() {
		return fmt.Errorf("unexpected end of file after position")
	}
	defPosStr := strings.TrimSpace(parser.Line())
	defPos, err := strconv.Atoi(defPosStr)
	if err != nil {
		return fmt.Errorf("invalid default position on line %d: %w", parser.LineNum(), err)
	}
	mobile.DefaultPos = defPos

	// Parse sex
	if !parser.NextLine() {
		return fmt.Errorf("unexpected end of file after default position")
	}
	sexStr := strings.TrimSpace(parser.Line())
	sex, err := strconv.Atoi(sexStr)
	if err != nil {
		return fmt.Errorf("invalid sex on line %d: %w", parser.LineNum(), err)
	}
	mobile.Sex = sex

	// Parse class
	if !parser.NextLine() {
		return fmt.Errorf("unexpected end of file after sex")
	}
	classStr := strings.TrimSpace(parser.Line())
	class, err := strconv.Atoi(classStr)
	if err != nil {
		return fmt.Errorf("invalid class on line %d: %w", parser.LineNum(), err)
	}
	mobile.Class = class

	// Parse level
	if !parser.NextLine() {
		return fmt.Errorf("unexpected end of file after class")
	}
	levelStr := strings.TrimSpace(parser.Line())
	level, err := strconv.Atoi(levelStr)
	if err != nil {
		return fmt.Errorf("invalid level on line %d: %w", parser.LineNum(), err)
	}
	mobile.Level = level

	// Skip the rest of the detailed mobile data
	// In the original code, there are more fields like birth, weight, height, etc.
	// We'll skip these for now, but we could add them if needed
	for parser.NextLine() {
		// Skip this line
		if strings.HasPrefix(parser.Line(), "#") {
			// We've reached a new mobile, so we need to back up
			parser.BackUp()
			break
		}
	}

	// Set default damage dice
	mobile.DamRoll = 0

	// Set default hitroll based on level
	if mobile.Level > 3 {
		mobile.HitRoll = mobile.Level - 3
	} else {
		mobile.HitRoll = 1
	}

	return nil
}
