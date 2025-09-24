package world

import (
	"testing"

	"github.com/wltechblog/DikuGo/pkg/types"
)

func TestCreateCityguardFromPrototype(t *testing.T) {
	// Create a mock storage
	storage := NewMockStorage()

	// Create a test room
	room := &types.Room{
		VNUM:        3001,
		Name:        "Test Room",
		Description: "This is a test room.",
		Characters:  make([]*types.Character, 0),
	}
	storage.rooms = append(storage.rooms, room)

	// Create a test world
	world, err := NewWorld(nil, storage)
	if err != nil {
		t.Fatalf("Failed to create world: %v", err)
	}
	world.rooms[3001] = room

	// Create a cityguard prototype
	cityguard := &types.Mobile{
		VNUM:        3060,
		Name:        "cityguard guard",
		ShortDesc:   "the Cityguard",
		LongDesc:    "A cityguard stands here.",
		Description: "A big, strong, helpful, trustworthy guard.",
		ActFlags:    193,
		Level:       10,
		HitRoll:     10,
		AC:          [3]int{20, 20, 20},
		Dice:        [3]int{1, 12, 123},
		DamRoll:     3,
		Gold:        500,
		Experience:  9000,
		Position:    8,
		DefaultPos:  8,
		Sex:         1,
	}
	world.mobiles[3060] = cityguard

	// Create a mob from the prototype
	mob := world.CreateMobFromPrototype(3060, room)
	if mob == nil {
		t.Fatalf("Failed to create mob from prototype")
	}

	// Check the mob stats
	if mob.Level != 10 {
		t.Errorf("Expected level 10, got %d", mob.Level)
	}
	if mob.HitRoll != 10 {
		t.Errorf("Expected hitroll 10, got %d", mob.HitRoll)
	}
	if mob.ArmorClass[0] != 20 || mob.ArmorClass[1] != 20 || mob.ArmorClass[2] != 20 {
		t.Errorf("Expected AC [20 20 20], got %v", mob.ArmorClass)
	}
	if mob.DamRoll != 3 {
		t.Errorf("Expected damroll 3, got %d", mob.DamRoll)
	}
	if mob.Gold != 500 {
		t.Errorf("Expected gold 500, got %d", mob.Gold)
	}
	if mob.Experience != 9000 {
		t.Errorf("Expected experience 9000, got %d", mob.Experience)
	}
}
