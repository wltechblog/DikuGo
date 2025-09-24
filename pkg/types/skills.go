package types

// Skill constants
const (
	// Basic skills
	SKILL_NONE = 0

	// Combat skills
	SKILL_BASH     = 1
	SKILL_RESCUE   = 2
	SKILL_BACKSTAB = 3
	SKILL_KICK     = 4

	// Thief skills
	SKILL_HIDE      = 5
	SKILL_STEAL     = 6
	SKILL_SNEAK     = 7
	SKILL_PICK_LOCK = 8

	// General skills
	SKILL_DODGE = 9

	// Maximum skill number
	MAX_SKILLS = 10
)

// Skill names
var SkillNames = map[int]string{
	SKILL_NONE:      "None",
	SKILL_BASH:      "Bash",
	SKILL_RESCUE:    "Rescue",
	SKILL_BACKSTAB:  "Backstab",
	SKILL_KICK:      "Kick",
	SKILL_HIDE:      "Hide",
	SKILL_STEAL:     "Steal",
	SKILL_SNEAK:     "Sneak",
	SKILL_PICK_LOCK: "Pick Lock",
	SKILL_DODGE:     "Dodge",
}

// GetSkillName returns the name of a skill
func GetSkillName(skill int) string {
	if name, ok := SkillNames[skill]; ok {
		return name
	}
	return "Unknown"
}

// Default skill levels by class and level
// These are the starting percentages for each skill
var DefaultSkills = map[int]map[int]map[int]int{
	// CLASS_MAGIC_USER - Mages don't have skills, they use spells
	CLASS_MAGIC_USER: {},

	// CLASS_CLERIC - Clerics don't have skills, they use spells
	CLASS_CLERIC: {},

	// CLASS_THIEF
	CLASS_THIEF: {
		// Level 1
		1: {
			SKILL_BACKSTAB:  10,
			SKILL_HIDE:      10,
			SKILL_STEAL:     5,
			SKILL_SNEAK:     10,
			SKILL_PICK_LOCK: 5,
		},
		// Level 5
		5: {
			SKILL_BACKSTAB:  30,
			SKILL_HIDE:      30,
			SKILL_STEAL:     25,
			SKILL_SNEAK:     30,
			SKILL_PICK_LOCK: 25,
		},
		// Level 10
		10: {
			SKILL_BACKSTAB:  50,
			SKILL_HIDE:      50,
			SKILL_STEAL:     45,
			SKILL_SNEAK:     50,
			SKILL_PICK_LOCK: 45,
		},
		// Level 15
		15: {
			SKILL_BACKSTAB:  70,
			SKILL_HIDE:      70,
			SKILL_STEAL:     65,
			SKILL_SNEAK:     70,
			SKILL_PICK_LOCK: 65,
		},
		// Level 20
		20: {
			SKILL_BACKSTAB:  90,
			SKILL_HIDE:      90,
			SKILL_STEAL:     85,
			SKILL_SNEAK:     90,
			SKILL_PICK_LOCK: 85,
		},
	},

	// CLASS_WARRIOR
	CLASS_WARRIOR: {
		// Level 1
		1: {
			SKILL_BASH:   10,
			SKILL_RESCUE: 5,
			SKILL_KICK:   10,
		},
		// Level 5
		5: {
			SKILL_BASH:   30,
			SKILL_RESCUE: 25,
			SKILL_KICK:   30,
		},
		// Level 10
		10: {
			SKILL_BASH:   50,
			SKILL_RESCUE: 45,
			SKILL_KICK:   50,
		},
		// Level 15
		15: {
			SKILL_BASH:   70,
			SKILL_RESCUE: 65,
			SKILL_KICK:   70,
		},
		// Level 20
		20: {
			SKILL_BASH:   90,
			SKILL_RESCUE: 85,
			SKILL_KICK:   90,
		},
	},
}

// GetDefaultSkill returns the default skill level for a character
func GetDefaultSkill(class, level, skill int) int {
	// Check if the class has skills
	classSkills, ok := DefaultSkills[class]
	if !ok {
		return 0
	}

	// Find the highest level bracket that's less than or equal to the character's level
	var highestBracket int
	for bracket := range classSkills {
		if bracket <= level && bracket > highestBracket {
			highestBracket = bracket
		}
	}

	// If no bracket was found, return 0
	if highestBracket == 0 {
		return 0
	}

	// Check if the skill exists in the bracket
	if skillLevel, ok := classSkills[highestBracket][skill]; ok {
		return skillLevel
	}

	return 0
}

// Skill delay in combat rounds (PULSE_VIOLENCE)
var SkillDelay = map[int]int{
	SKILL_BASH:     2, // 2 combat rounds
	SKILL_RESCUE:   2, // 2 combat rounds
	SKILL_BACKSTAB: 2, // 2 combat rounds
	SKILL_KICK:     3, // 3 combat rounds
}

// GetSkillDelay returns the delay for a skill in combat rounds
func GetSkillDelay(skill int) int {
	if delay, ok := SkillDelay[skill]; ok {
		return delay
	}
	return 1 // Default to 1 round
}

// GetSkillClass returns the class that can use a skill
func GetSkillClass(skill int) int {
	switch skill {
	case SKILL_BASH, SKILL_RESCUE, SKILL_KICK:
		return CLASS_WARRIOR
	case SKILL_BACKSTAB, SKILL_HIDE, SKILL_STEAL, SKILL_SNEAK, SKILL_PICK_LOCK:
		return CLASS_THIEF
	default:
		return CLASS_UNDEFINED
	}
}

// ClassCanUseSkill checks if a class can use a specific skill based on skill type
func ClassCanUseSkill(class, skill int) bool {
	// Get the class that can use this skill
	skillClass := GetSkillClass(skill)

	// If the skill doesn't have a specific class, anyone can use it
	if skillClass == CLASS_UNDEFINED {
		return true
	}

	// Check if the character's class matches the skill's class
	return class == skillClass
}
