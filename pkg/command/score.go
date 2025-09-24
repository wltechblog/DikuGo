package command

import (
	"fmt"
	"strings"

	"github.com/wltechblog/DikuGo/pkg/types"
)

// ScoreCommand displays the character's stats
type ScoreCommand struct{}

// Execute runs the score command
func (c *ScoreCommand) Execute(ch *types.Character, args string) error {
	Score(ch, args)
	return nil
}

// Name returns the name of the command
func (c *ScoreCommand) Name() string {
	return "score"
}

// Aliases returns the aliases for the command
func (c *ScoreCommand) Aliases() []string {
	return []string{"sc"}
}

// MinPosition returns the minimum position required to execute the command
func (c *ScoreCommand) MinPosition() int {
	return types.POS_RESTING
}

// Level returns the minimum level required to execute the command
func (c *ScoreCommand) Level() int {
	return 0
}

// LogCommand returns whether the command should be logged
func (c *ScoreCommand) LogCommand() bool {
	return false
}

// Score displays the character's stats
func Score(ch *types.Character, args string) bool {
	// Get the character's stats
	hp := ch.HP
	maxHP := ch.MaxHitPoints
	mana := ch.ManaPoints
	maxMana := ch.MaxManaPoints
	move := ch.MovePoints
	maxMove := ch.MaxMovePoints

	// Calculate HP, Mana, and Move percentages for colored output
	hpPercent := float64(hp) / float64(maxHP) * 100
	manaPercent := float64(mana) / float64(maxMana) * 100
	movePercent := float64(move) / float64(maxMove) * 100

	// Get color codes for HP, Mana, and Move
	hpColor := getColorCode(hpPercent)
	manaColor := getColorCode(manaPercent)
	moveColor := getColorCode(movePercent)

	// Build the score output
	var sb strings.Builder

	sb.WriteString("\r\n")
	sb.WriteString(fmt.Sprintf("You are %s, level %d %s %s.\r\n",
		ch.Name, ch.Level, getRaceName(ch.Race), getClassName(ch.Class)))

	sb.WriteString(fmt.Sprintf("%sHP: %d/%d%s  ", hpColor, hp, maxHP, "\033[0m"))
	sb.WriteString(fmt.Sprintf("%sMana: %d/%d%s  ", manaColor, mana, maxMana, "\033[0m"))
	sb.WriteString(fmt.Sprintf("%sMove: %d/%d%s\r\n", moveColor, move, maxMove, "\033[0m"))

	sb.WriteString(fmt.Sprintf("Str: %d  Int: %d  Wis: %d  Dex: %d  Con: %d  Cha: %d\r\n",
		ch.Abilities[0], ch.Abilities[1], ch.Abilities[2],
		ch.Abilities[3], ch.Abilities[4], ch.Abilities[5]))

	sb.WriteString(fmt.Sprintf("Hit Roll: %d  Damage Roll: %d\r\n", ch.HitRoll, ch.DamRoll))
	sb.WriteString(fmt.Sprintf("Armor Class: %d\r\n", ch.ArmorClass[0]))
	sb.WriteString(fmt.Sprintf("Gold: %d  Experience: %d\r\n", ch.Gold, ch.Experience))

	// Calculate experience needed for next level
	expForNextLevel := getExpForLevel(ch.Level + 1)
	expNeeded := expForNextLevel - ch.Experience
	if expNeeded < 0 {
		expNeeded = 0
	}

	sb.WriteString(fmt.Sprintf("Experience needed for next level: %d\r\n", expNeeded))
	sb.WriteString(fmt.Sprintf("Alignment: %d (%s)\r\n", ch.Alignment, getAlignmentString(ch.Alignment)))

	// Send the score to the character
	ch.SendMessage(sb.String())

	return true
}

// getColorCode returns ANSI color codes based on percentage
func getColorCode(percent float64) string {
	if percent < 25 {
		return "\033[31m" // Red
	} else if percent < 50 {
		return "\033[33m" // Yellow
	} else {
		return "\033[32m" // Green
	}
}

// getRaceName returns the name of a race
func getRaceName(race int) string {
	races := map[int]string{
		types.RACE_HUMAN:    "Human",
		types.RACE_ELF:      "Elf",
		types.RACE_DWARF:    "Dwarf",
		types.RACE_GNOME:    "Gnome",
		types.RACE_HALFLING: "Halfling",
	}

	if name, ok := races[race]; ok {
		return name
	}
	return "Unknown"
}

// getClassName returns the name of a class
func getClassName(class int) string {
	classes := map[int]string{
		types.CLASS_WARRIOR: "Warrior",
		types.CLASS_CLERIC:  "Cleric",
		types.CLASS_THIEF:   "Thief",
		types.CLASS_MAGE:    "Mage",
	}

	if name, ok := classes[class]; ok {
		return name
	}
	return "Unknown"
}

// getAlignmentString returns a string representation of an alignment
func getAlignmentString(alignment int) string {
	if alignment > 700 {
		return "Good"
	} else if alignment > 350 {
		return "Slightly Good"
	} else if alignment > 0 {
		return "Neutral Good"
	} else if alignment == 0 {
		return "Neutral"
	} else if alignment > -350 {
		return "Neutral Evil"
	} else if alignment > -700 {
		return "Slightly Evil"
	} else {
		return "Evil"
	}
}

// getExpForLevel returns the experience needed for a given level
func getExpForLevel(level int) int {
	// Simple exponential experience curve
	// Level 1: 0 exp
	// Level 2: 1000 exp
	// Level 3: 3000 exp
	// Level 4: 6000 exp
	// Level 5: 10000 exp
	// etc.
	if level <= 1 {
		return 0
	}

	return 1000 * (level - 1) * level / 2
}
