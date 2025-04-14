package types

// Item type constants
const (
	ITEM_LIGHT      = 1
	ITEM_SCROLL     = 2
	ITEM_WAND       = 3
	ITEM_STAFF      = 4
	ITEM_WEAPON     = 5
	ITEM_FIREWEAPON = 6
	ITEM_MISSILE    = 7
	ITEM_TREASURE   = 8
	ITEM_ARMOR      = 9
	ITEM_POTION     = 10
	ITEM_WORN       = 11
	ITEM_OTHER      = 12
	ITEM_TRASH      = 13
	ITEM_TRAP       = 14
	ITEM_CONTAINER  = 15
	ITEM_NOTE       = 16
	ITEM_DRINKCON   = 17
	ITEM_KEY        = 18
	ITEM_FOOD       = 19
	ITEM_MONEY      = 20
	ITEM_PEN        = 21
	ITEM_BOAT       = 22
	ITEM_FOUNTAIN   = 23
)

// Item extra flag constants
const (
	ITEM_GLOW       = (1 << 0)
	ITEM_HUM        = (1 << 1)
	ITEM_NORENT     = (1 << 2)
	ITEM_NODONATE   = (1 << 3)
	ITEM_NOINVIS    = (1 << 4)
	ITEM_INVISIBLE  = (1 << 5)
	ITEM_MAGIC      = (1 << 6)
	ITEM_NODROP     = (1 << 7)
	ITEM_BLESS      = (1 << 8)
	ITEM_ANTI_GOOD  = (1 << 9)
	ITEM_ANTI_EVIL  = (1 << 10)
	ITEM_ANTI_NEUTRAL = (1 << 11)
	ITEM_ANTI_MAGIC = (1 << 12)
	ITEM_ANTI_CLERIC = (1 << 13)
	ITEM_ANTI_THIEF = (1 << 14)
	ITEM_ANTI_WARRIOR = (1 << 15)
	ITEM_NOSELL     = (1 << 16)
	ITEM_NOPICK     = (1 << 17)
)

// Item wear flag constants
const (
	ITEM_WEAR_TAKE    = (1 << 0)
	ITEM_WEAR_FINGER  = (1 << 1)
	ITEM_WEAR_NECK    = (1 << 2)
	ITEM_WEAR_BODY    = (1 << 3)
	ITEM_WEAR_HEAD    = (1 << 4)
	ITEM_WEAR_LEGS    = (1 << 5)
	ITEM_WEAR_FEET    = (1 << 6)
	ITEM_WEAR_HANDS   = (1 << 7)
	ITEM_WEAR_ARMS    = (1 << 8)
	ITEM_WEAR_SHIELD  = (1 << 9)
	ITEM_WEAR_ABOUT   = (1 << 10)
	ITEM_WEAR_WAIST   = (1 << 11)
	ITEM_WEAR_WRIST   = (1 << 12)
	ITEM_WEAR_WIELD   = (1 << 13)
	ITEM_WEAR_HOLD    = (1 << 14)
)

// Wear position constants
const (
	WEAR_LIGHT   = 0
	WEAR_FINGER_R = 1
	WEAR_FINGER_L = 2
	WEAR_NECK_1  = 3
	WEAR_NECK_2  = 4
	WEAR_BODY    = 5
	WEAR_HEAD    = 6
	WEAR_LEGS    = 7
	WEAR_FEET    = 8
	WEAR_HANDS   = 9
	WEAR_ARMS    = 10
	WEAR_SHIELD  = 11
	WEAR_ABOUT   = 12
	WEAR_WAIST   = 13
	WEAR_WRIST_R = 14
	WEAR_WRIST_L = 15
	WEAR_WIELD   = 16
	WEAR_HOLD    = 17
	NUM_WEARS    = 18
)
