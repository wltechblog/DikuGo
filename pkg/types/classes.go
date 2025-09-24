package types

// ClassInfo contains information about a character class
type ClassInfo struct {
	Name        string   // Class name
	Description string   // Class description
	Abilities   []int    // Priority of abilities (indexes into ABILITY_* constants)
	Skills      []int    // Skills available to this class
	HitDie      int      // Hit points gained per level
	ManaDie     int      // Mana points gained per level
	MoveDie     int      // Movement points gained per level
	StartingHP  int      // Starting hit points
	StartingMana int     // Starting mana points
	StartingMove int     // Starting movement points
	StartingGold int     // Starting gold
	Alignment   int      // Starting alignment tendency
}

// ClassData contains information about all character classes
var ClassData = map[int]*ClassInfo{
	CLASS_MAGIC_USER: {
		Name:        "Magic User",
		Description: "Magic Users are the masters of arcane magic. They can cast powerful spells but are physically weak and cannot use most weapons or armor.",
		Abilities:   []int{ABILITY_INT, ABILITY_WIS, ABILITY_DEX, ABILITY_STR, ABILITY_CON, ABILITY_CHA},
		Skills:      []int{}, // Magic users rely on spells, not skills
		HitDie:      4,
		ManaDie:     10,
		MoveDie:     2,
		StartingHP:  10,
		StartingMana: 150,
		StartingMove: 100,
		StartingGold: 100,
		Alignment:   0, // Neutral tendency
	},
	CLASS_CLERIC: {
		Name:        "Cleric",
		Description: "Clerics are the masters of divine magic. They can heal wounds, cure afflictions, and smite their enemies with divine power.",
		Abilities:   []int{ABILITY_WIS, ABILITY_INT, ABILITY_STR, ABILITY_DEX, ABILITY_CON, ABILITY_CHA},
		Skills:      []int{}, // Clerics rely on spells, not skills
		HitDie:      6,
		ManaDie:     8,
		MoveDie:     2,
		StartingHP:  12,
		StartingMana: 125,
		StartingMove: 100,
		StartingGold: 150,
		Alignment:   350, // Good tendency
	},
	CLASS_THIEF: {
		Name:        "Thief",
		Description: "Thieves are masters of stealth and trickery. They can pick locks, disarm traps, hide in shadows, and backstab unwary opponents.",
		Abilities:   []int{ABILITY_DEX, ABILITY_STR, ABILITY_CON, ABILITY_INT, ABILITY_WIS, ABILITY_CHA},
		Skills:      []int{SKILL_SNEAK, SKILL_HIDE, SKILL_STEAL, SKILL_BACKSTAB, SKILL_PICK_LOCK},
		HitDie:      6,
		ManaDie:     4,
		MoveDie:     3,
		StartingHP:  12,
		StartingMana: 50,
		StartingMove: 120,
		StartingGold: 250,
		Alignment:   -350, // Evil tendency
	},
	CLASS_WARRIOR: {
		Name:        "Warrior",
		Description: "Warriors are masters of combat. They can use all weapons and armor, and are the most physically powerful class.",
		Abilities:   []int{ABILITY_STR, ABILITY_DEX, ABILITY_CON, ABILITY_WIS, ABILITY_INT, ABILITY_CHA},
		Skills:      []int{SKILL_BASH, SKILL_RESCUE, SKILL_KICK},
		HitDie:      10,
		ManaDie:     2,
		MoveDie:     2,
		StartingHP:  15,
		StartingMana: 50,
		StartingMove: 110,
		StartingGold: 200,
		Alignment:   0, // Neutral tendency
	},
}

// GetClassName returns the name of a class
func GetClassName(class int) string {
	if info, ok := ClassData[class]; ok {
		return info.Name
	}
	return "Unknown"
}

// GetClassDescription returns the description of a class
func GetClassDescription(class int) string {
	if info, ok := ClassData[class]; ok {
		return info.Description
	}
	return "Unknown class"
}

// GetClassAbilities returns the ability priorities for a class
func GetClassAbilities(class int) []int {
	if info, ok := ClassData[class]; ok {
		return info.Abilities
	}
	return []int{ABILITY_STR, ABILITY_DEX, ABILITY_CON, ABILITY_INT, ABILITY_WIS, ABILITY_CHA}
}

// GetClassSkills returns the skills available to a class
func GetClassSkills(class int) []int {
	if info, ok := ClassData[class]; ok {
		return info.Skills
	}
	return []int{}
}

// GetClassStartingStats returns the starting stats for a class
func GetClassStartingStats(class int) (hp, mana, move, gold, alignment int) {
	if info, ok := ClassData[class]; ok {
		return info.StartingHP, info.StartingMana, info.StartingMove, info.StartingGold, info.Alignment
	}
	return 10, 100, 100, 100, 0
}

// CanUseSkill checks if a class can use a specific skill
func CanUseSkill(class, skill int) bool {
	if info, ok := ClassData[class]; ok {
		for _, s := range info.Skills {
			if s == skill {
				return true
			}
		}
	}
	return false
}

// CanWearItem checks if a class can wear a specific item
func CanWearItem(class int, itemType, itemFlags int) bool {
	// Check class-specific restrictions
	switch class {
	case CLASS_MAGIC_USER:
		// Magic users can't wear heavy armor
		if itemType == ITEM_ARMOR && itemFlags > 2 {
			return false
		}
		// Magic users can only use daggers, staves, and wands
		if itemType == ITEM_WEAPON && itemFlags != 0 && itemFlags != 3 && itemFlags != 4 {
			return false
		}
	case CLASS_CLERIC:
		// Clerics can't use edged weapons
		if itemType == ITEM_WEAPON && (itemFlags == 1 || itemFlags == 2) {
			return false
		}
	case CLASS_THIEF:
		// Thieves can't wear heavy armor
		if itemType == ITEM_ARMOR && itemFlags > 3 {
			return false
		}
	case CLASS_WARRIOR:
		// Warriors can use all items
		return true
	}
	
	return true
}

// GetClassHitDie returns the hit die for a class
func GetClassHitDie(class int) int {
	if info, ok := ClassData[class]; ok {
		return info.HitDie
	}
	return 6
}

// GetClassManaDie returns the mana die for a class
func GetClassManaDie(class int) int {
	if info, ok := ClassData[class]; ok {
		return info.ManaDie
	}
	return 4
}

// GetClassMoveDie returns the movement die for a class
func GetClassMoveDie(class int) int {
	if info, ok := ClassData[class]; ok {
		return info.MoveDie
	}
	return 2
}
