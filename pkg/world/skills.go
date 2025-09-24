package world

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/wltechblog/DikuGo/pkg/types"
)

// InitializeNewCharacter initializes a new character with appropriate stats based on class
func (w *World) InitializeNewCharacter(ch *types.Character) {
	// Initialize abilities based on class
	ch.Abilities = RollAbilities(ch.Class)

	// Set up base hit points based on class and constitution
	baseHP := 10 // Base HP for all classes
	conBonus := getConstitutionHPBonus(ch.Abilities[types.ABILITY_CON])

	// Get starting stats from class data
	startingHP, startingMana, startingMove, startingGold, startingAlignment := types.GetClassStartingStats(ch.Class)

	// Apply constitution bonus to HP
	baseHP = startingHP + conBonus

	// Set hit points
	ch.HP = baseHP
	ch.MaxHitPoints = baseHP

	// Set mana points
	ch.ManaPoints = startingMana
	ch.MaxManaPoints = startingMana

	// Set move points based on class and dexterity
	dexBonus := ch.Abilities[types.ABILITY_DEX] - 10 // Simple bonus based on DEX above 10
	if dexBonus < 0 {
		dexBonus = 0
	}

	// Apply dexterity bonus to movement points
	baseMove := startingMove + dexBonus

	// Set move points
	ch.MovePoints = baseMove
	ch.MaxMovePoints = baseMove

	// Set starting gold
	ch.Gold = startingGold

	// Set alignment tendencies
	ch.Alignment = startingAlignment

	// Initialize skills
	w.InitializeCharacterSkills(ch)

	// Initialize spells for magic users and clerics
	w.InitializeCharacterSpells(ch)

	// Initialize conditions (hunger, thirst, drunk)
	ch.Conditions = [3]int{24, 24, 0} // Full, quenched, sober
}

// InitializeCharacterSkills initializes a character's skills based on class and level
func (w *World) InitializeCharacterSkills(ch *types.Character) {
	// Skip if not a player
	if ch.IsNPC {
		return
	}

	// Initialize skills map if it doesn't exist
	if ch.Skills == nil {
		ch.Skills = make(map[int]int)
	}

	// Initialize LastSkillTime map if it doesn't exist
	if ch.LastSkillTime == nil {
		ch.LastSkillTime = make(map[int]time.Time)
	}

	// Set default skills based on class and level
	for skill := 1; skill < types.MAX_SKILLS; skill++ {
		// Get default skill level
		skillLevel := types.GetDefaultSkill(ch.Class, ch.Level, skill)

		// Only set if the skill is applicable to this class
		if skillLevel > 0 {
			ch.Skills[skill] = skillLevel
		}
	}
}

// HasSkill checks if a character has a skill
func (w *World) HasSkill(ch *types.Character, skill int) bool {
	// Check if the character has the skill
	level, ok := ch.Skills[skill]
	return ok && level > 0
}

// GetSkillLevel returns a character's skill level
func (w *World) GetSkillLevel(ch *types.Character, skill int) int {
	// Check if the character has the skill
	level, ok := ch.Skills[skill]
	if !ok {
		return 0
	}
	return level
}

// CheckSkillSuccess checks if a skill use is successful
func (w *World) CheckSkillSuccess(ch *types.Character, skill int) bool {
	// Get the skill level
	level := w.GetSkillLevel(ch, skill)
	if level <= 0 {
		return false
	}

	// Roll percentile dice
	roll := rand.Intn(100) + 1

	// Check if the roll is less than or equal to the skill level
	return roll <= level
}

// CanUseSkill checks if a character can use a skill (cooldown, position, etc.)
func (w *World) CanUseSkill(ch *types.Character, skill int) (bool, string) {
	// Check if the character has the skill
	if !w.HasSkill(ch, skill) {
		return false, "You don't know that skill."
	}

	// Check if the character is in the right position
	if ch.Position < types.POS_FIGHTING {
		return false, "You're in no position to do that!"
	}

	// Check if the skill is on cooldown
	lastUse, ok := ch.LastSkillTime[skill]
	if ok {
		// Get the skill delay in combat rounds
		delay := types.GetSkillDelay(skill)

		// Convert to time duration (2 seconds per combat round)
		cooldown := time.Duration(delay) * 2 * time.Second

		// Check if enough time has passed
		if time.Since(lastUse) < cooldown {
			return false, "You're not ready to use that skill again yet."
		}
	}

	// Check if the class can use this skill
	if !types.ClassCanUseSkill(ch.Class, skill) {
		return false, fmt.Sprintf("Only %s can use that skill.", types.GetClassName(types.GetSkillClass(skill)))
	}

	return true, ""
}

// UseSkill marks a skill as used (sets cooldown)
func (w *World) UseSkill(ch *types.Character, skill int) {
	// Set the last use time
	ch.LastSkillTime[skill] = time.Now()
}

// RollAbilities rolls abilities for a new character based on class
// This follows the original DikuMUD method:
// 1. Roll 4d6, drop the lowest die, sum the remaining 3 dice
// 2. Do this 5 times to get 5 ability scores
// 3. Assign the scores to abilities based on class priority
func RollAbilities(class int) [6]int {
	// Roll 5 ability scores
	scores := rollAbilityScores()

	// Create a new abilities array with default values
	abilities := [6]int{9, 9, 9, 9, 9, 9}

	// Get the ability priorities for this class
	abilityPriorities := types.GetClassAbilities(class)

	// Assign scores based on class priority
	for i := 0; i < 5 && i < len(abilityPriorities); i++ {
		abilities[abilityPriorities[i]] = scores[i]
	}

	// CHA is always the last priority for all classes
	abilities[types.ABILITY_CHA] = 9 // Default CHA

	return abilities
}

// rollAbilityScores rolls 5 ability scores using the 4d6-drop-lowest method
// and returns them in descending order
func rollAbilityScores() []int {
	scores := make([]int, 5)

	for i := 0; i < 5; i++ {
		// Roll 4d6
		dice := []int{
			rand.Intn(6) + 1,
			rand.Intn(6) + 1,
			rand.Intn(6) + 1,
			rand.Intn(6) + 1,
		}

		// Find the lowest die
		lowestIndex := 0
		for j := 1; j < 4; j++ {
			if dice[j] < dice[lowestIndex] {
				lowestIndex = j
			}
		}

		// Sum the remaining 3 dice
		sum := 0
		for j := 0; j < 4; j++ {
			if j != lowestIndex {
				sum += dice[j]
			}
		}

		scores[i] = sum
	}

	// Sort scores in descending order
	for i := 0; i < len(scores)-1; i++ {
		for j := i + 1; j < len(scores); j++ {
			if scores[j] > scores[i] {
				scores[i], scores[j] = scores[j], scores[i]
			}
		}
	}

	return scores
}

// getConstitutionHPBonus returns the hit point bonus for a given constitution score
func getConstitutionHPBonus(con int) int {
	switch {
	case con <= 3:
		return -3
	case con <= 5:
		return -2
	case con <= 7:
		return -1
	case con <= 13:
		return 0
	case con <= 15:
		return 1
	case con <= 17:
		return 2
	default:
		return 3
	}
}

// ImproveSkill gives a chance to improve a skill after use
func (w *World) ImproveSkill(ch *types.Character, skill int, success bool) {
	// Skip if not a player
	if ch.IsNPC {
		return
	}

	// Get the current skill level
	level := w.GetSkillLevel(ch, skill)
	if level <= 0 || level >= 95 {
		return // Can't improve if you don't have the skill or it's already maxed
	}

	// Calculate chance to improve
	// Higher chance on failure, lower chance at higher skill levels
	var improveChance int
	if success {
		improveChance = 10 - (level / 10)
	} else {
		improveChance = 20 - (level / 5)
	}

	// Ensure minimum chance
	if improveChance < 1 {
		improveChance = 1
	}

	// Roll for improvement
	roll := rand.Intn(100) + 1
	if roll <= improveChance {
		// Improve the skill
		ch.Skills[skill]++

		// Notify the character
		ch.SendMessage(fmt.Sprintf("You have become better at %s!\r\n", types.GetSkillName(skill)))

		// Log the improvement
		log.Printf("%s improved %s to %d%%", ch.Name, types.GetSkillName(skill), ch.Skills[skill])
	}
}

// InitializeCharacterSpells initializes a character's spells based on class and level
func (w *World) InitializeCharacterSpells(ch *types.Character) {
	// Skip if not a player
	if ch.IsNPC {
		return
	}

	// Only magic users and clerics get spells
	if ch.Class != types.CLASS_MAGIC_USER && ch.Class != types.CLASS_CLERIC {
		return
	}

	// Initialize spells map if it doesn't exist
	if ch.Spells == nil {
		ch.Spells = make(map[int]int)
	}

	// Set initial spell levels based on class and level
	for spellID := 1; spellID < types.MAX_SPELLS; spellID++ {
		// Get minimum level for this spell
		minLevel := types.GetSpellMinLevel(spellID, ch.Class)

		// If character is high enough level, give them the spell
		if minLevel <= ch.Level && minLevel < 21 {
			// Initial spell level is 50% + 2% per level above minimum
			spellLevel := 50 + ((ch.Level - minLevel) * 2)
			if spellLevel > 95 {
				spellLevel = 95 // Cap at 95%
			}
			ch.Spells[spellID] = spellLevel
		}
	}
}
