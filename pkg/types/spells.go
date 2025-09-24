package types

// Spell constants
const (
	SPELL_UNDEFINED     = 0
	SPELL_ARMOR         = 1
	SPELL_TELEPORT      = 2
	SPELL_BLESS         = 3
	SPELL_BLINDNESS     = 4
	SPELL_BURNING_HANDS = 5
	SPELL_CALL_LIGHTNING = 6
	SPELL_CHARM_PERSON  = 7
	SPELL_CHILL_TOUCH   = 8
	SPELL_CLONE         = 9
	SPELL_COLOR_SPRAY   = 10
	SPELL_CONTROL_WEATHER = 11
	SPELL_CREATE_FOOD   = 12
	SPELL_CREATE_WATER  = 13
	SPELL_CURE_BLIND    = 14
	SPELL_CURE_CRITIC   = 15
	SPELL_CURE_LIGHT    = 16
	SPELL_CURSE         = 17
	SPELL_DETECT_EVIL   = 18
	SPELL_DETECT_INVISIBLE = 19
	SPELL_DETECT_MAGIC  = 20
	SPELL_DETECT_POISON = 21
	SPELL_DISPEL_EVIL   = 22
	SPELL_EARTHQUAKE    = 23
	SPELL_ENCHANT_WEAPON = 24
	SPELL_ENERGY_DRAIN  = 25
	SPELL_FIREBALL      = 26
	SPELL_HARM          = 27
	SPELL_HEAL          = 28
	SPELL_INVISIBLE     = 29
	SPELL_LIGHTNING_BOLT = 30
	SPELL_LOCATE_OBJECT = 31
	SPELL_MAGIC_MISSILE = 32
	SPELL_POISON        = 33
	SPELL_PROTECTION_FROM_EVIL = 34
	SPELL_REMOVE_CURSE  = 35
	SPELL_SANCTUARY     = 36
	SPELL_SHOCKING_GRASP = 37
	SPELL_SLEEP         = 38
	SPELL_STRENGTH      = 39
	SPELL_SUMMON        = 40
	SPELL_VENTRILOQUATE = 41
	SPELL_WORD_OF_RECALL = 42
	SPELL_REMOVE_POISON = 43
	SPELL_SENSE_LIFE    = 44
	SPELL_IDENTIFY      = 45
	
	// Breath weapons
	SPELL_FIRE_BREATH   = 46
	SPELL_GAS_BREATH    = 47
	SPELL_FROST_BREATH  = 48
	SPELL_ACID_BREATH   = 49
	SPELL_LIGHTNING_BREATH = 50
	
	MAX_SPELLS          = 51
)

// Spell types
const (
	SPELL_TYPE_SPELL  = 0
	SPELL_TYPE_POTION = 1
	SPELL_TYPE_WAND   = 2
	SPELL_TYPE_STAFF  = 3
	SPELL_TYPE_SCROLL = 4
)

// Spell target flags
const (
	TAR_IGNORE        = 1
	TAR_CHAR_ROOM     = 2
	TAR_CHAR_WORLD    = 4
	TAR_FIGHT_SELF    = 8
	TAR_FIGHT_VICT    = 16
	TAR_SELF_ONLY     = 32 // Only a check, use with TAR_CHAR_ROOM
	TAR_SELF_NONO     = 64 // Only a check, use with TAR_CHAR_ROOM
	TAR_OBJ_INV       = 128
	TAR_OBJ_ROOM      = 256
	TAR_OBJ_WORLD     = 512
	TAR_OBJ_EQUIP     = 1024
)

// SpellInfo contains information about a spell
type SpellInfo struct {
	Name           string // Spell name
	MinPosition    int    // Minimum position to cast
	MinMana        int    // Minimum mana required
	Beats          int    // Delay in combat rounds
	MinLevelMage   int    // Minimum level for magic users
	MinLevelCleric int    // Minimum level for clerics
	Targets        int    // Valid targets (TAR_XXX)
	Violent        bool   // Is this a violent spell
	Memorized      bool   // Is this spell memorized
}

// SpellData contains information about all spells
var SpellData = map[int]*SpellInfo{
	SPELL_ARMOR: {
		Name:           "armor",
		MinPosition:    POS_STANDING,
		MinMana:        5,
		Beats:          12,
		MinLevelMage:   5,
		MinLevelCleric: 1,
		Targets:        TAR_CHAR_ROOM,
		Violent:        false,
	},
	SPELL_TELEPORT: {
		Name:           "teleport",
		MinPosition:    POS_FIGHTING,
		MinMana:        35,
		Beats:          12,
		MinLevelMage:   8,
		MinLevelCleric: 21, // Clerics don't get teleport
		Targets:        TAR_SELF_ONLY,
		Violent:        false,
	},
	SPELL_BLESS: {
		Name:           "bless",
		MinPosition:    POS_STANDING,
		MinMana:        5,
		Beats:          12,
		MinLevelMage:   21, // Mages don't get bless
		MinLevelCleric: 5,
		Targets:        TAR_OBJ_INV | TAR_CHAR_ROOM,
		Violent:        false,
	},
	SPELL_BLINDNESS: {
		Name:           "blindness",
		MinPosition:    POS_STANDING,
		MinMana:        5,
		Beats:          12,
		MinLevelMage:   8,
		MinLevelCleric: 6,
		Targets:        TAR_CHAR_ROOM,
		Violent:        true,
	},
	SPELL_BURNING_HANDS: {
		Name:           "burning hands",
		MinPosition:    POS_FIGHTING,
		MinMana:        15,
		Beats:          12,
		MinLevelMage:   5,
		MinLevelCleric: 21, // Clerics don't get burning hands
		Targets:        TAR_CHAR_ROOM | TAR_FIGHT_VICT,
		Violent:        true,
	},
	SPELL_CALL_LIGHTNING: {
		Name:           "call lightning",
		MinPosition:    POS_FIGHTING,
		MinMana:        15,
		Beats:          12,
		MinLevelMage:   21, // Mages don't get call lightning
		MinLevelCleric: 12,
		Targets:        TAR_CHAR_ROOM | TAR_FIGHT_VICT,
		Violent:        true,
	},
	SPELL_CHARM_PERSON: {
		Name:           "charm person",
		MinPosition:    POS_STANDING,
		MinMana:        5,
		Beats:          12,
		MinLevelMage:   14,
		MinLevelCleric: 21, // Clerics don't get charm person
		Targets:        TAR_CHAR_ROOM | TAR_SELF_NONO,
		Violent:        true,
	},
	SPELL_CHILL_TOUCH: {
		Name:           "chill touch",
		MinPosition:    POS_FIGHTING,
		MinMana:        15,
		Beats:          12,
		MinLevelMage:   3,
		MinLevelCleric: 21, // Clerics don't get chill touch
		Targets:        TAR_CHAR_ROOM | TAR_FIGHT_VICT,
		Violent:        true,
	},
	SPELL_COLOR_SPRAY: {
		Name:           "color spray",
		MinPosition:    POS_FIGHTING,
		MinMana:        15,
		Beats:          12,
		MinLevelMage:   11,
		MinLevelCleric: 21, // Clerics don't get color spray
		Targets:        TAR_CHAR_ROOM | TAR_FIGHT_VICT,
		Violent:        true,
	},
	SPELL_CURE_BLIND: {
		Name:           "cure blindness",
		MinPosition:    POS_STANDING,
		MinMana:        5,
		Beats:          12,
		MinLevelMage:   21, // Mages don't get cure blindness
		MinLevelCleric: 4,
		Targets:        TAR_CHAR_ROOM,
		Violent:        false,
	},
	SPELL_CURE_CRITIC: {
		Name:           "cure critical",
		MinPosition:    POS_FIGHTING,
		MinMana:        15,
		Beats:          12,
		MinLevelMage:   21, // Mages don't get cure critical
		MinLevelCleric: 9,
		Targets:        TAR_CHAR_ROOM,
		Violent:        false,
	},
	SPELL_CURE_LIGHT: {
		Name:           "cure light",
		MinPosition:    POS_FIGHTING,
		MinMana:        10,
		Beats:          12,
		MinLevelMage:   21, // Mages don't get cure light
		MinLevelCleric: 1,
		Targets:        TAR_CHAR_ROOM,
		Violent:        false,
	},
	SPELL_CURSE: {
		Name:           "curse",
		MinPosition:    POS_STANDING,
		MinMana:        5,
		Beats:          12,
		MinLevelMage:   12,
		MinLevelCleric: 14,
		Targets:        TAR_CHAR_ROOM | TAR_OBJ_INV,
		Violent:        true,
	},
	SPELL_DETECT_EVIL: {
		Name:           "detect evil",
		MinPosition:    POS_STANDING,
		MinMana:        5,
		Beats:          12,
		MinLevelMage:   21, // Mages don't get detect evil
		MinLevelCleric: 2,
		Targets:        TAR_CHAR_ROOM,
		Violent:        false,
	},
	SPELL_DETECT_INVISIBLE: {
		Name:           "detect invisible",
		MinPosition:    POS_STANDING,
		MinMana:        5,
		Beats:          12,
		MinLevelMage:   2,
		MinLevelCleric: 6,
		Targets:        TAR_CHAR_ROOM,
		Violent:        false,
	},
	SPELL_DETECT_MAGIC: {
		Name:           "detect magic",
		MinPosition:    POS_STANDING,
		MinMana:        5,
		Beats:          12,
		MinLevelMage:   2,
		MinLevelCleric: 4,
		Targets:        TAR_CHAR_ROOM,
		Violent:        false,
	},
	SPELL_DETECT_POISON: {
		Name:           "detect poison",
		MinPosition:    POS_STANDING,
		MinMana:        5,
		Beats:          12,
		MinLevelMage:   21, // Mages don't get detect poison
		MinLevelCleric: 3,
		Targets:        TAR_CHAR_ROOM | TAR_OBJ_INV,
		Violent:        false,
	},
	SPELL_DISPEL_EVIL: {
		Name:           "dispel evil",
		MinPosition:    POS_FIGHTING,
		MinMana:        15,
		Beats:          12,
		MinLevelMage:   21, // Mages don't get dispel evil
		MinLevelCleric: 10,
		Targets:        TAR_CHAR_ROOM | TAR_FIGHT_VICT,
		Violent:        true,
	},
	SPELL_EARTHQUAKE: {
		Name:           "earthquake",
		MinPosition:    POS_FIGHTING,
		MinMana:        15,
		Beats:          12,
		MinLevelMage:   21, // Mages don't get earthquake
		MinLevelCleric: 8,
		Targets:        TAR_IGNORE,
		Violent:        true,
	},
	SPELL_ENCHANT_WEAPON: {
		Name:           "enchant weapon",
		MinPosition:    POS_STANDING,
		MinMana:        30,
		Beats:          24,
		MinLevelMage:   12,
		MinLevelCleric: 21, // Clerics don't get enchant weapon
		Targets:        TAR_OBJ_INV,
		Violent:        false,
	},
	SPELL_ENERGY_DRAIN: {
		Name:           "energy drain",
		MinPosition:    POS_FIGHTING,
		MinMana:        40,
		Beats:          12,
		MinLevelMage:   13,
		MinLevelCleric: 19,
		Targets:        TAR_CHAR_ROOM | TAR_FIGHT_VICT,
		Violent:        true,
	},
	SPELL_FIREBALL: {
		Name:           "fireball",
		MinPosition:    POS_FIGHTING,
		MinMana:        15,
		Beats:          12,
		MinLevelMage:   15,
		MinLevelCleric: 21, // Clerics don't get fireball
		Targets:        TAR_CHAR_ROOM | TAR_FIGHT_VICT,
		Violent:        true,
	},
	SPELL_HARM: {
		Name:           "harm",
		MinPosition:    POS_FIGHTING,
		MinMana:        40,
		Beats:          12,
		MinLevelMage:   21, // Mages don't get harm
		MinLevelCleric: 15,
		Targets:        TAR_CHAR_ROOM | TAR_FIGHT_VICT,
		Violent:        true,
	},
	SPELL_HEAL: {
		Name:           "heal",
		MinPosition:    POS_FIGHTING,
		MinMana:        50,
		Beats:          12,
		MinLevelMage:   21, // Mages don't get heal
		MinLevelCleric: 16,
		Targets:        TAR_CHAR_ROOM,
		Violent:        false,
	},
	SPELL_INVISIBLE: {
		Name:           "invisible",
		MinPosition:    POS_STANDING,
		MinMana:        5,
		Beats:          12,
		MinLevelMage:   4,
		MinLevelCleric: 21, // Clerics don't get invisible
		Targets:        TAR_CHAR_ROOM | TAR_OBJ_INV | TAR_OBJ_ROOM,
		Violent:        false,
	},
	SPELL_LIGHTNING_BOLT: {
		Name:           "lightning bolt",
		MinPosition:    POS_FIGHTING,
		MinMana:        15,
		Beats:          12,
		MinLevelMage:   9,
		MinLevelCleric: 21, // Clerics don't get lightning bolt
		Targets:        TAR_CHAR_ROOM | TAR_FIGHT_VICT,
		Violent:        true,
	},
	SPELL_LOCATE_OBJECT: {
		Name:           "locate object",
		MinPosition:    POS_STANDING,
		MinMana:        20,
		Beats:          12,
		MinLevelMage:   6,
		MinLevelCleric: 8,
		Targets:        TAR_IGNORE,
		Violent:        false,
	},
	SPELL_MAGIC_MISSILE: {
		Name:           "magic missile",
		MinPosition:    POS_FIGHTING,
		MinMana:        15,
		Beats:          12,
		MinLevelMage:   1,
		MinLevelCleric: 21, // Clerics don't get magic missile
		Targets:        TAR_CHAR_ROOM | TAR_FIGHT_VICT,
		Violent:        true,
	},
	SPELL_POISON: {
		Name:           "poison",
		MinPosition:    POS_STANDING,
		MinMana:        5,
		Beats:          12,
		MinLevelMage:   21, // Mages don't get poison
		MinLevelCleric: 13,
		Targets:        TAR_CHAR_ROOM | TAR_OBJ_INV,
		Violent:        true,
	},
	SPELL_PROTECTION_FROM_EVIL: {
		Name:           "protection from evil",
		MinPosition:    POS_STANDING,
		MinMana:        5,
		Beats:          12,
		MinLevelMage:   21, // Mages don't get protection from evil
		MinLevelCleric: 8,
		Targets:        TAR_CHAR_ROOM,
		Violent:        false,
	},
	SPELL_REMOVE_CURSE: {
		Name:           "remove curse",
		MinPosition:    POS_STANDING,
		MinMana:        5,
		Beats:          12,
		MinLevelMage:   21, // Mages don't get remove curse
		MinLevelCleric: 12,
		Targets:        TAR_CHAR_ROOM | TAR_OBJ_INV,
		Violent:        false,
	},
	SPELL_SANCTUARY: {
		Name:           "sanctuary",
		MinPosition:    POS_STANDING,
		MinMana:        75,
		Beats:          12,
		MinLevelMage:   21, // Mages don't get sanctuary
		MinLevelCleric: 15,
		Targets:        TAR_CHAR_ROOM,
		Violent:        false,
	},
	SPELL_SHOCKING_GRASP: {
		Name:           "shocking grasp",
		MinPosition:    POS_FIGHTING,
		MinMana:        15,
		Beats:          12,
		MinLevelMage:   7,
		MinLevelCleric: 21, // Clerics don't get shocking grasp
		Targets:        TAR_CHAR_ROOM | TAR_FIGHT_VICT,
		Violent:        true,
	},
	SPELL_SLEEP: {
		Name:           "sleep",
		MinPosition:    POS_STANDING,
		MinMana:        15,
		Beats:          12,
		MinLevelMage:   10,
		MinLevelCleric: 21, // Clerics don't get sleep
		Targets:        TAR_CHAR_ROOM,
		Violent:        true,
	},
	SPELL_STRENGTH: {
		Name:           "strength",
		MinPosition:    POS_STANDING,
		MinMana:        5,
		Beats:          12,
		MinLevelMage:   6,
		MinLevelCleric: 7,
		Targets:        TAR_CHAR_ROOM,
		Violent:        false,
	},
	SPELL_SUMMON: {
		Name:           "summon",
		MinPosition:    POS_STANDING,
		MinMana:        50,
		Beats:          12,
		MinLevelMage:   21, // Mages don't get summon
		MinLevelCleric: 10,
		Targets:        TAR_CHAR_WORLD,
		Violent:        false,
	},
	SPELL_WORD_OF_RECALL: {
		Name:           "word of recall",
		MinPosition:    POS_FIGHTING,
		MinMana:        5,
		Beats:          12,
		MinLevelMage:   21, // Mages don't get word of recall
		MinLevelCleric: 12,
		Targets:        TAR_CHAR_ROOM,
		Violent:        false,
	},
	SPELL_REMOVE_POISON: {
		Name:           "remove poison",
		MinPosition:    POS_STANDING,
		MinMana:        5,
		Beats:          12,
		MinLevelMage:   21, // Mages don't get remove poison
		MinLevelCleric: 5,
		Targets:        TAR_CHAR_ROOM | TAR_OBJ_INV,
		Violent:        false,
	},
	SPELL_SENSE_LIFE: {
		Name:           "sense life",
		MinPosition:    POS_STANDING,
		MinMana:        5,
		Beats:          12,
		MinLevelMage:   16,
		MinLevelCleric: 9,
		Targets:        TAR_CHAR_ROOM,
		Violent:        false,
	},
}

// GetSpellName returns the name of a spell
func GetSpellName(spell int) string {
	if info, ok := SpellData[spell]; ok {
		return info.Name
	}
	return "unknown"
}

// GetSpellByName returns the spell ID for a given name
func GetSpellByName(name string) int {
	for id, info := range SpellData {
		if info.Name == name {
			return id
		}
	}
	return SPELL_UNDEFINED
}

// GetSpellMinLevel returns the minimum level required to cast a spell for a given class
func GetSpellMinLevel(spell int, class int) int {
	if info, ok := SpellData[spell]; ok {
		if class == CLASS_MAGIC_USER {
			return info.MinLevelMage
		} else if class == CLASS_CLERIC {
			return info.MinLevelCleric
		}
	}
	return 99 // Unreachable level
}

// GetSpellMana returns the mana cost for a spell
func GetSpellMana(spell int) int {
	if info, ok := SpellData[spell]; ok {
		return info.MinMana
	}
	return 0
}

// GetSpellPosition returns the minimum position required to cast a spell
func GetSpellPosition(spell int) int {
	if info, ok := SpellData[spell]; ok {
		return info.MinPosition
	}
	return POS_STANDING
}

// GetSpellTargets returns the valid targets for a spell
func GetSpellTargets(spell int) int {
	if info, ok := SpellData[spell]; ok {
		return info.Targets
	}
	return 0
}

// IsSpellViolent returns whether a spell is violent
func IsSpellViolent(spell int) bool {
	if info, ok := SpellData[spell]; ok {
		return info.Violent
	}
	return false
}

// GetSpellDelay returns the delay for a spell in combat rounds
func GetSpellDelay(spell int) int {
	if info, ok := SpellData[spell]; ok {
		return info.Beats
	}
	return 12 // Default delay
}

// SpellNames returns a list of all spell names
var SpellNames []string

func init() {
	// Initialize SpellNames
	SpellNames = make([]string, MAX_SPELLS)
	for i := 0; i < MAX_SPELLS; i++ {
		if info, ok := SpellData[i]; ok {
			SpellNames[i] = info.Name
		} else {
			SpellNames[i] = "unknown"
		}
	}
}
