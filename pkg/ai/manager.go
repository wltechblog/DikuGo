package ai

import (
	"log"
	"math/rand"
	"time"

	"github.com/wltechblog/DikuGo/pkg/types"
)

// World interface for AI manager
type World interface {
	GetRooms() []*types.Room
	GetMobiles() []*types.Character
	FindObjectsInRoom(room *types.Room) []*types.ObjectInstance
	MoveCharacter(character *types.Character, room *types.Room) error
	GetRandomExitRoom(room *types.Room) (*types.Room, int)
}

// Manager handles AI for the game
type Manager struct {
	world    World
	lastTick time.Time
}

// NewManager creates a new AI manager
func NewManager(world World) *Manager {
	return &Manager{
		world:    world,
		lastTick: time.Now(),
	}
}

// Tick updates all AI entities
func (m *Manager) Tick() {
	// Get the current time
	now := time.Now()

	// Calculate the time since the last tick
	elapsed := now.Sub(m.lastTick)
	m.lastTick = now

	// Get all mobiles
	mobiles := m.world.GetMobiles()

	// Process the mobiles
	m.processMobiles(mobiles, elapsed)
}

// TickWithMobiles updates AI entities with a provided list of mobiles
// This is used to avoid deadlocks when the caller already has a lock
func (m *Manager) TickWithMobiles(mobiles []*types.Character) {
	// Get the current time
	now := time.Now()

	// Calculate the time since the last tick
	elapsed := now.Sub(m.lastTick)
	m.lastTick = now

	// Process the mobiles
	m.processMobiles(mobiles, elapsed)
}

// processMobiles processes AI behavior for a list of mobiles
func (m *Manager) processMobiles(mobiles []*types.Character, elapsed time.Duration) {
	// Process each mobile
	for _, mobile := range mobiles {
		// Skip if not an NPC
		if !mobile.IsNPC {
			continue
		}

		// Skip if the mobile is fighting
		if mobile.Fighting != nil {
			continue
		}

		// Process special procedures
		if m.processSpecialProcedures(mobile) {
			continue
		}

		// Process behaviors based on flags
		m.processBehaviors(mobile, elapsed)
	}
}

// processSpecialProcedures processes special procedures for a mobile
// Returns true if a special procedure was executed
func (m *Manager) processSpecialProcedures(mobile *types.Character) bool {
	// Skip if the mobile doesn't have a prototype
	if mobile.Prototype == nil {
		return false
	}

	// Skip if the mobile doesn't have special procedures
	if mobile.Prototype.Functions == nil || len(mobile.Prototype.Functions) == 0 {
		return false
	}

	// Execute each special procedure
	for _, fn := range mobile.Prototype.Functions {
		if fn != nil {
			if fn(mobile, "") {
				return true
			}
		}
	}

	return false
}

// processBehaviors processes behaviors for a mobile based on its flags
func (m *Manager) processBehaviors(mobile *types.Character, elapsed time.Duration) {
	// Skip if the mobile doesn't have a prototype
	if mobile.Prototype == nil {
		return
	}

	// Skip if the mobile is not in a room
	if mobile.InRoom == nil {
		return
	}

	// Process scavenger behavior
	if mobile.Prototype.ActFlags&types.ACT_SCAVENGER != 0 {
		m.processScavengerBehavior(mobile)
	}

	// Process wandering behavior (if not a sentinel)
	if mobile.Prototype.ActFlags&types.ACT_SENTINEL == 0 {
		m.processWanderingBehavior(mobile, elapsed)
	}

	// Process aggressive behavior
	if mobile.Prototype.ActFlags&types.ACT_AGGRESSIVE != 0 {
		m.processAggressiveBehavior(mobile)
	}
}

// processScavengerBehavior processes scavenger behavior for a mobile
func (m *Manager) processScavengerBehavior(mobile *types.Character) {
	// Find objects in the room
	objects := m.world.FindObjectsInRoom(mobile.InRoom)
	if len(objects) == 0 {
		return
	}

	// Pick a random object
	obj := objects[rand.Intn(len(objects))]

	// TODO: Implement get command for mobiles
	log.Printf("Mobile %s would pick up %s", mobile.Name, obj.Prototype.Name)
}

// processWanderingBehavior processes wandering behavior for a mobile
func (m *Manager) processWanderingBehavior(mobile *types.Character, elapsed time.Duration) {
	// Only move occasionally (about once every 5 minutes on average)
	if rand.Float64() > elapsed.Seconds()/(5*60) {
		return
	}

	// Skip if the mobile is in a shop
	if mobile.InRoom.Shop != nil {
		return
	}

	// Get a random exit
	nextRoom, _ := m.world.GetRandomExitRoom(mobile.InRoom)
	if nextRoom == nil {
		return
	}

	// Check if the mobile should stay in its zone
	if mobile.Prototype.ActFlags&types.ACT_STAY_ZONE != 0 {
		// Skip if the next room is not in the same zone
		if mobile.InRoom.Zone != nextRoom.Zone {
			return
		}
	}

	// Move the mobile
	err := m.world.MoveCharacter(mobile, nextRoom)
	if err != nil {
		log.Printf("Error moving mobile %s: %v", mobile.Name, err)
	}
}

// processAggressiveBehavior processes aggressive behavior for a mobile
func (m *Manager) processAggressiveBehavior(mobile *types.Character) {
	// Skip if the mobile is wimpy and injured
	if mobile.Prototype.ActFlags&types.ACT_WIMPY != 0 && mobile.HP < mobile.MaxHitPoints/2 {
		return
	}

	// Acquire read lock for the room before accessing its character list
	room := mobile.InRoom
	room.RLock()

	// Find a player to attack
	var target *types.Character
	for _, character := range room.Characters {
		// Skip NPCs and sleeping players if wimpy
		if character.IsNPC || (mobile.Prototype.ActFlags&types.ACT_WIMPY != 0 && character.Position <= types.POS_SLEEPING) {
			continue
		}
		// Found a potential target
		target = character
		break // Attack the first valid target found
	}

	// Release the read lock
	log.Printf("processAggressiveBehavior: Releasing RLock for room %d (Mobile: %s)", room.VNUM, mobile.Name)
	room.RUnlock()

	// If a target was found, initiate attack (outside the room lock)
	if target != nil {
		// TODO: Implement attack command for mobiles
		log.Printf("Mobile %s would attack %s", mobile.Name, target.Name)
	}
}
