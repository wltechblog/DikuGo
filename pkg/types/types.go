package types

import (
	"sync"
	"time"
)

// ExtraDescription represents an extra description for a room or object
type ExtraDescription struct {
	Keywords    string
	Description string
}

// Room represents a location in the game world
type Room struct {
	VNUM        int
	Name        string
	Description string
	Flags       uint32
	SectorType  int
	Exits       [6]*Exit
	Characters  []*Character
	Objects     []*ObjectInstance
	ExtraDescs  []*ExtraDescription
	Functions   []func(*Character, string) bool // Special procedures
	Zone        *Zone
	Shop        *Shop
	mutex       sync.RWMutex // Internal mutex
}

// Lock acquires a write lock on the room
func (r *Room) Lock() {
	r.mutex.Lock()
}

// Unlock releases a write lock on the room
func (r *Room) Unlock() {
	r.mutex.Unlock()
}

// RLock acquires a read lock on the room
func (r *Room) RLock() {
	r.mutex.RLock()
}

// RUnlock releases a read lock on the room
func (r *Room) RUnlock() {
	r.mutex.RUnlock()
}

// Exit represents a connection between rooms
type Exit struct {
	Direction   int
	Description string
	Keywords    string
	Flags       uint32
	Key         int
	DestVnum    int // Destination room VNUM
}

// Object represents an object prototype
type Object struct {
	VNUM        int
	Name        string
	ShortDesc   string
	Description string
	ActionDesc  string
	Type        int
	ExtraFlags  uint32
	WearFlags   uint32
	Value       [4]int
	Weight      int
	Cost        int
	ExtraDescs  []*ExtraDescription
	Affects     []*Affect
}

// ObjectInstance represents an instance of an object in the game
type ObjectInstance struct {
	Prototype  *Object
	InRoom     *Room
	CarriedBy  *Character
	WornBy     *Character
	WornOn     int
	InObj      *ObjectInstance
	Contains   []*ObjectInstance
	Timer      int
	Affects    []*Affect
	CustomDesc string
	ExtraDescs []*ExtraDescription // Instance-specific extra descriptions
	mutex      sync.RWMutex
}

// Mobile represents a mobile (NPC) prototype
type Mobile struct {
	VNUM        int
	Name        string
	ShortDesc   string
	LongDesc    string
	Description string
	ActFlags    uint32
	AffectFlags uint32
	Alignment   int
	Level       int
	Sex         int
	Class       int
	Race        int
	Position    int
	DefaultPos  int
	Gold        int
	Experience  int
	LoadPos     int
	DamageType  int
	AttackType  int
	AC          [3]int
	HitRoll     int
	DamRoll     int
	Dice        [3]int                          // num, size, bonus
	Abilities   [6]int                          // STR, INT, WIS, DEX, CON, CHA
	Functions   []func(*Character, string) bool // Special procedures
	Equipment   []MobEquipment                  // Default equipment
}

// MobEquipment represents an equipment item for a mobile
type MobEquipment struct {
	ObjectVNUM int // VNUM of the object
	Position   int // Wear position
	Chance     int // Chance of being equipped (0-100)
}

// MobRespawn represents a mob scheduled for respawning
type MobRespawn struct {
	MobVNUM     int       // VNUM of the mob
	RoomVNUM    int       // VNUM of the room to respawn in
	RespawnTime time.Time // Time to respawn
}

// Character represents a character (player or NPC) in the game
type Character struct {
	Name          string
	ShortDesc     string // For NPCs
	LongDesc      string // For NPCs
	Description   string
	Level         int
	Sex           int
	Class         int
	Race          int
	Position      int
	Gold          int
	Experience    int
	Alignment     int
	HP            int // Current hit points
	HitPoints     int // Deprecated, use HP instead
	MaxHitPoints  int
	ManaPoints    int
	MaxManaPoints int
	MovePoints    int
	MaxMovePoints int
	ArmorClass    [3]int
	HitRoll       int
	DamRoll       int
	Abilities     [6]int // STR, INT, WIS, DEX, CON, CHA
	Affects       []*Affect
	Skills        map[int]int
	Spells        map[int]int
	Equipment     []*ObjectInstance
	Inventory     []*ObjectInstance
	InRoom        *Room
	RoomVNUM      int // VNUM of the room the character is in
	Fighting      *Character
	Following     *Character
	Followers     []*Character
	IsNPC         bool
	ActFlags      uint32                          // NPC behavior flags
	Prototype     *Mobile                         // If NPC
	Functions     []func(*Character, string) bool // Special procedures
	LastLogin     time.Time
	Password      string // Hashed
	Title         string
	Prompt        string
	Flags         uint32
	World         interface{} // Reference to the world
	mutex         sync.RWMutex
}

// IsNPCFlag returns true if the character is an NPC
func (c *Character) IsNPCFlag() bool {
	return c.IsNPC
}

// SendMessage sends a message to the character
func (c *Character) SendMessage(message string) {
	// TODO: Implement message sending
	// For now, just log the message
	// log.Printf("Message to %s: %s", c.Name, message)
}

// Zone represents a zone in the game world
type Zone struct {
	VNUM          int
	Name          string
	Lifespan      int
	Age           int
	ResetMode     int
	TopRoom       int
	BottomRoom    int
	Commands      []*ZoneCommand
	MinVNUM       int // Minimum VNUM in this zone
	MaxVNUM       int // Maximum VNUM in this zone
	ResetInterval int // Seconds between resets
	mutex         sync.RWMutex
}

// ShouldReset returns true if the zone should reset
func (z *Zone) ShouldReset() bool {
	return z.Age >= z.Lifespan
}

// ZoneCommand represents a command in a zone reset
type ZoneCommand struct {
	Command byte
	IfFlag  int
	Arg1    int
	Arg2    int
	Arg3    int
	Arg4    int
}

// Shop represents a shop in the game world
type Shop struct {
	VNUM       int
	RoomVNUM   int
	MobileVNUM int
	BuyTypes   []int
	Producing  []int // Items the shop produces (sells)
	ProfitBuy  float64
	ProfitSell float64
	OpenHour   int
	CloseHour  int
	Messages   []string                        // Shop messages
	Functions  []func(*Character, string) bool // Special procedures
	mutex      sync.RWMutex
}

// Affect represents a temporary effect on a character or object
type Affect struct {
	Type      int
	Duration  int
	Modifier  int
	Location  int
	BitVector uint32
}

// TimeWeather represents the game time and weather
type TimeWeather struct {
	Hour     int
	Day      int
	Month    int
	Year     int
	Sunlight int
	Weather  int
	Change   int
	mutex    sync.RWMutex
}
