package ai

import (
	"testing"
	"time"

	"github.com/wltechblog/DikuGo/pkg/types"
)

// TestManagerTick tests the AI manager tick function with additional behaviors
func TestManagerTick(t *testing.T) {
	// Create a mock world
	world := &MockWorld{
		rooms:   make([]*types.Room, 0),
		mobiles: make([]*types.Character, 0),
		objects: make(map[*types.Room][]*types.ObjectInstance),
	}

	// Create rooms
	room1 := &types.Room{
		VNUM:        1,
		Name:        "Room 1",
		Description: "This is room 1",
		Characters:  make([]*types.Character, 0),
	}

	room2 := &types.Room{
		VNUM:        2,
		Name:        "Room 2",
		Description: "This is room 2",
		Characters:  make([]*types.Character, 0),
	}

	world.rooms = append(world.rooms, room1, room2)
	world.exitRoom = room2

	// Create a mobile prototype
	mobileProto := &types.Mobile{
		VNUM:      1,
		Name:      "test mob",
		ShortDesc: "a test mobile",
		LongDesc:  "A test mobile is here.",
		ActFlags:  types.ACT_SENTINEL,
	}

	// Create a mobile
	mobile := &types.Character{
		Name:        "Test Mobile",
		ShortDesc:   "a test mobile",
		LongDesc:    "A test mobile is here.",
		Description: "This is a test mobile for AI testing.",
		IsNPC:       true,
		InRoom:      room1,
		Prototype:   mobileProto,
	}

	// Add mobile to room
	room1.Characters = append(room1.Characters, mobile)
	world.mobiles = append(world.mobiles, mobile)

	// Create an AI manager
	manager := NewManager(world)

	// Test sentinel mobile (shouldn't move)
	manager.Tick()

	if len(world.moveLog) > 0 {
		t.Errorf("Sentinel mobile moved: %v", world.moveLog)
	}

	// Change mobile to wanderer
	mobileProto.ActFlags = 0                             // Clear sentinel flag
	manager.lastTick = time.Now().Add(-10 * time.Minute) // Force a move

	// Test wandering mobile (should move)
	manager.Tick()

	if len(world.moveLog) == 0 {
		t.Errorf("Wandering mobile didn't move after tick")
	}
}

// TestAggressiveBehavior tests the aggressive behavior
func TestAggressiveBehavior(t *testing.T) {
	// Create a mock world
	world := &MockWorld{
		rooms:   make([]*types.Room, 0),
		mobiles: make([]*types.Character, 0),
		objects: make(map[*types.Room][]*types.ObjectInstance),
	}

	// Create a room
	room := &types.Room{
		VNUM:        1,
		Name:        "Room 1",
		Description: "This is room 1",
		Characters:  make([]*types.Character, 0),
	}

	world.rooms = append(world.rooms, room)

	// Create a player
	player := &types.Character{
		Name:        "Test Player",
		ShortDesc:   "a test player",
		LongDesc:    "A test player is here.",
		Description: "This is a test player.",
		IsNPC:       false,
		InRoom:      room,
		Position:    types.POS_STANDING,
	}

	// Create a mobile prototype
	mobileProto := &types.Mobile{
		VNUM:      2,
		Name:      "aggressive mob",
		ShortDesc: "an aggressive mobile",
		LongDesc:  "An aggressive mobile is here.",
		ActFlags:  types.ACT_AGGRESSIVE,
	}

	// Create an aggressive mobile
	mobile := &types.Character{
		Name:        "Aggressive Mobile",
		ShortDesc:   "an aggressive mobile",
		LongDesc:    "An aggressive mobile is here.",
		Description: "This is an aggressive mobile for AI testing.",
		IsNPC:       true,
		InRoom:      room,
		Prototype:   mobileProto,
	}

	// Add characters to room
	room.Characters = append(room.Characters, player, mobile)
	world.mobiles = append(world.mobiles, mobile)

	// Create an AI manager
	manager := NewManager(world)

	// Test aggressive behavior
	manager.Tick()

	// TODO: Add assertions when attack command is implemented
}

// TestScavengerBehavior tests the scavenger behavior
func TestScavengerBehavior(t *testing.T) {
	// Create a mock world
	world := &MockWorld{
		rooms:   make([]*types.Room, 0),
		mobiles: make([]*types.Character, 0),
		objects: make(map[*types.Room][]*types.ObjectInstance),
	}

	// Create a room
	room := &types.Room{
		VNUM:        1,
		Name:        "Room 1",
		Description: "This is room 1",
		Characters:  make([]*types.Character, 0),
	}

	world.rooms = append(world.rooms, room)

	// Create a mobile prototype
	mobileProto := &types.Mobile{
		VNUM:      3,
		Name:      "scavenger mob",
		ShortDesc: "a scavenger mobile",
		LongDesc:  "A scavenger mobile is here.",
		ActFlags:  types.ACT_SCAVENGER,
	}

	// Create a scavenger mobile
	mobile := &types.Character{
		Name:        "Scavenger Mobile",
		ShortDesc:   "a scavenger mobile",
		LongDesc:    "A scavenger mobile is here.",
		Description: "This is a scavenger mobile for AI testing.",
		IsNPC:       true,
		InRoom:      room,
		Prototype:   mobileProto,
	}

	// Add mobile to room
	room.Characters = append(room.Characters, mobile)
	world.mobiles = append(world.mobiles, mobile)

	// Create an object
	obj := &types.ObjectInstance{
		Prototype: &types.Object{
			Name:        "test object",
			ShortDesc:   "a test object",
			Description: "A test object is lying here.",
		},
		InRoom: room,
	}

	// Add object to room
	world.objects[room] = append(world.objects[room], obj)

	// Create an AI manager
	manager := NewManager(world)

	// Test scavenger behavior
	manager.Tick()

	// TODO: Add assertions when get command is implemented
}
