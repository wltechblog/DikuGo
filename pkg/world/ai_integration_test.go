package world

import (
	"testing"

	"github.com/wltechblog/DikuGo/pkg/types"
)

func TestAIIntegration(t *testing.T) {
	// Create a mock storage
	storage := NewMockStorage()

	// Create a test room
	room := &types.Room{
		VNUM:        1,
		Name:        "Test Room",
		Description: "This is a test room.",
		Exits:       [6]*types.Exit{},
	}
	storage.rooms = append(storage.rooms, room)

	// Create a test mobile prototype
	mob := &types.Mobile{
		VNUM:        1,
		Name:        "test mob",
		ShortDesc:   "A test mob is here.",
		LongDesc:    "A test mob is standing here.",
		Description: "This is a test mob for AI integration.",
		Level:       1,
	}
	storage.mobiles = append(storage.mobiles, mob)

	// Create a world
	world, err := NewWorld(nil, storage)
	if err != nil {
		t.Fatalf("Failed to create world: %v", err)
	}

	// Initialize the AI system
	world.InitAI()

	// Create a mob in the room
	mobChar := world.CreateMobFromPrototype(1, room)
	if mobChar == nil {
		t.Fatal("Failed to create mob from prototype")
	}

	// Tick the AI system
	world.TickAI()

	// No assertions needed, just make sure it doesn't crash
}
