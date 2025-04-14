package ai

import (
	"testing"
	"time"

	"github.com/wltechblog/DikuGo/pkg/types"
)

// MockWorld is a mock implementation of the world interface for testing
type MockWorld struct {
	rooms    []*types.Room
	mobiles  []*types.Character
	objects  map[*types.Room][]*types.ObjectInstance
	moveLog  []string
	exitRoom *types.Room
}

func (m *MockWorld) GetRooms() []*types.Room {
	return m.rooms
}

func (m *MockWorld) GetMobiles() []*types.Character {
	return m.mobiles
}

func (m *MockWorld) FindObjectsInRoom(room *types.Room) []*types.ObjectInstance {
	return m.objects[room]
}

func (m *MockWorld) MoveCharacter(character *types.Character, room *types.Room) error {
	m.moveLog = append(m.moveLog, character.Name+" moved to room "+room.Name)
	character.InRoom = room
	return nil
}

func (m *MockWorld) GetRandomExitRoom(room *types.Room) (*types.Room, int) {
	return m.exitRoom, 0
}

func TestMobAI(t *testing.T) {
	// Create a mock world
	world := &MockWorld{
		rooms:   make([]*types.Room, 0),
		mobiles: make([]*types.Character, 0),
		objects: make(map[*types.Room][]*types.ObjectInstance),
	}

	// Create a room
	room := &types.Room{
		VNUM:       1,
		Name:       "Test Room",
		Characters: make([]*types.Character, 0),
	}
	world.rooms = append(world.rooms, room)

	// Create an exit room
	exitRoom := &types.Room{
		VNUM:       2,
		Name:       "Exit Room",
		Characters: make([]*types.Character, 0),
	}
	world.rooms = append(world.rooms, exitRoom)
	world.exitRoom = exitRoom

	// Create a mobile prototype
	mobileProto := &types.Mobile{
		VNUM:      1,
		Name:      "test mob",
		ShortDesc: "a test mob",
		LongDesc:  "A test mob is standing here.",
		ActFlags:  types.ACT_SENTINEL | types.ACT_AGGRESSIVE,
	}

	// Create a mobile instance
	mobile := &types.Character{
		Name:      "Test Mob",
		ShortDesc: "a test mob",
		IsNPC:     true,
		Prototype: mobileProto,
		InRoom:    room,
	}
	world.mobiles = append(world.mobiles, mobile)
	room.Characters = append(room.Characters, mobile)

	// Create an AI manager
	manager := NewManager(world)

	// Test the Tick method
	manager.Tick()

	// The mobile should not have moved (it's a sentinel)
	if len(world.moveLog) > 0 {
		t.Errorf("Expected mobile to stay in place (sentinel), but it moved: %v", world.moveLog)
	}

	// Test with a non-sentinel mobile
	mobileProto.ActFlags = types.ACT_AGGRESSIVE
	manager.lastTick = time.Now().Add(-10 * time.Minute) // Force a move
	manager.Tick()

	// The mobile should have moved
	if len(world.moveLog) == 0 {
		t.Errorf("Expected mobile to move, but it didn't")
	}
}
