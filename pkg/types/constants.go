package types

// Position constants
const (
	POS_DEAD     = 0
	POS_MORTALLY = 1
	POS_INCAP    = 2
	POS_STUNNED  = 3
	POS_SLEEPING = 4
	POS_RESTING  = 5
	POS_SITTING  = 6
	POS_FIGHTING = 7
	POS_STANDING = 8
)

// Room flag constants
const (
	ROOM_DARK = 1 << iota
	ROOM_DEATH
	ROOM_NO_MOB
	ROOM_INDOORS
	ROOM_LAWFUL
	ROOM_NEUTRAL
	ROOM_CHAOTIC
	ROOM_NO_MAGIC
	ROOM_TUNNEL
	ROOM_PRIVATE
	ROOM_NOSUMMON
)

// Direction constants
const (
	DIR_NORTH = 0
	DIR_EAST  = 1
	DIR_SOUTH = 2
	DIR_WEST  = 3
	DIR_UP    = 4
	DIR_DOWN  = 5
)

// Affect constants (AFF_XXX)
const (
	AFF_BLIND            = (1 << 0)
	AFF_INVISIBLE        = (1 << 1)
	AFF_DETECT_EVIL      = (1 << 2)
	AFF_DETECT_INVISIBLE = (1 << 3)
	AFF_DETECT_MAGIC     = (1 << 4)
	AFF_SENSE_LIFE       = (1 << 5)
	AFF_HOLD             = (1 << 6)
	AFF_SANCTUARY        = (1 << 7)
	AFF_GROUP            = (1 << 8)
	AFF_CURSE            = (1 << 10)
	AFF_FLAMING          = (1 << 11)
	AFF_POISON           = (1 << 12)
	AFF_PROTECT_EVIL     = (1 << 13)
	AFF_PARALYSIS        = (1 << 14)
	AFF_MORDEN_SWORD     = (1 << 15)
	AFF_FLAMING_SWORD    = (1 << 16)
	AFF_SLEEP            = (1 << 17)
	AFF_DODGE            = (1 << 18)
	AFF_SNEAK            = (1 << 19)
	AFF_HIDE             = (1 << 20)
	AFF_FEAR             = (1 << 21)
	AFF_CHARM            = (1 << 22)
	AFF_FOLLOW           = (1 << 23)
)

// Apply constants (APPLY_XXX)
const (
	APPLY_NONE          = 0
	APPLY_STR           = 1
	APPLY_DEX           = 2
	APPLY_INT           = 3
	APPLY_WIS           = 4
	APPLY_CON           = 5
	APPLY_SEX           = 6
	APPLY_CLASS         = 7
	APPLY_LEVEL         = 8
	APPLY_AGE           = 9
	APPLY_CHAR_WEIGHT   = 10
	APPLY_CHAR_HEIGHT   = 11
	APPLY_MANA          = 12
	APPLY_HIT           = 13
	APPLY_MOVE          = 14
	APPLY_GOLD          = 15
	APPLY_EXP           = 16
	APPLY_AC            = 17
	APPLY_ARMOR         = 17
	APPLY_HITROLL       = 18
	APPLY_DAMROLL       = 19
	APPLY_SAVING_PARA   = 20
	APPLY_SAVING_ROD    = 21
	APPLY_SAVING_PETRI  = 22
	APPLY_SAVING_BREATH = 23
	APPLY_SAVING_SPELL  = 24
)

// Weapon type constants
const (
	TYPE_HIT      = 0 // Default for bare hands
	TYPE_BLUDGEON = 1 // Blunt weapons
	TYPE_PIERCE   = 2 // Piercing weapons
	TYPE_SLASH    = 3 // Slashing weapons
	TYPE_WHIP     = 4 // Whips
	TYPE_CLAW     = 5 // Claws
	TYPE_BITE     = 6 // Bites
	TYPE_STING    = 7 // Stings
	TYPE_CRUSH    = 8 // Crushing attacks

	// Special damage types
	TYPE_SUFFERING = 100 // Damage from suffering (hunger, etc.)
	TYPE_POISON    = 101 // Damage from poison
)

// Sex constants
const (
	SEX_NEUTRAL = 0
	SEX_MALE    = 1
	SEX_FEMALE  = 2
)

// Race constants
const (
	RACE_HUMAN    = 0
	RACE_ELF      = 1
	RACE_DWARF    = 2
	RACE_GNOME    = 3
	RACE_HALFLING = 4
)

// Class constants - match original DikuMUD values
const (
	CLASS_UNDEFINED  = 0
	CLASS_MAGIC_USER = 1
	CLASS_CLERIC     = 2
	CLASS_THIEF      = 3
	CLASS_WARRIOR    = 4
	CLASS_MAGE       = CLASS_MAGIC_USER // Alias for backward compatibility
)

// Saving throw types
const (
	SAVING_PARA   = 0 // Paralyzation
	SAVING_ROD    = 1 // Rod/Staff/Wand
	SAVING_PETRI  = 2 // Petrification
	SAVING_BREATH = 3 // Breath weapon
	SAVING_SPELL  = 4 // Spell
)

// Condition types
const (
	COND_DRUNK  = 0
	COND_FULL   = 1
	COND_THIRST = 2
)

// Ability constants
const (
	ABILITY_STR = 0 // Strength
	ABILITY_DEX = 1 // Dexterity
	ABILITY_INT = 2 // Intelligence
	ABILITY_WIS = 3 // Wisdom
	ABILITY_CON = 4 // Constitution
	ABILITY_CHA = 5 // Charisma
)

// Object constants
const (
	MAX_OBJ_AFFECT = 6 // Maximum number of affects on an object
)

// Sector type constants
const (
	SECT_INSIDE       = 0
	SECT_CITY         = 1
	SECT_FIELD        = 2
	SECT_FOREST       = 3
	SECT_HILLS        = 4
	SECT_MOUNTAIN     = 5
	SECT_WATER_SWIM   = 6
	SECT_WATER_NOSWIM = 7
	SECT_UNDERWATER   = 8
	SECT_FLYING       = 9
)

// Liquid type constants
const (
	LIQ_WATER      = 0
	LIQ_BEER       = 1
	LIQ_WINE       = 2
	LIQ_ALE        = 3
	LIQ_DARKALE    = 4
	LIQ_WHISKY     = 5
	LIQ_LEMONADE   = 6
	LIQ_FIREBRT    = 7
	LIQ_LOCALSPC   = 8
	LIQ_SLIME      = 9
	LIQ_MILK       = 10
	LIQ_TEA        = 11
	LIQ_COFFEE     = 12
	LIQ_BLOOD      = 13
	LIQ_SALTWATER  = 14
	LIQ_CLEARWATER = 15
)
