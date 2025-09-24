package world

import (
	"testing"
	"time"

	"github.com/wltechblog/DikuGo/pkg/types"
)

func TestCityguardRespawn(t *testing.T) {
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

	// Create a test zone
	zone := &types.Zone{
		VNUM:     1,
		Name:     "Test Zone",
		Lifespan: 1, // 1 minute respawn time
		Age:      0,
		MinVNUM:  3000,
		MaxVNUM:  3999,
	}
	room.Zone = zone

	// Create a test world
	world, err := NewWorld(nil, storage)
	if err != nil {
		t.Fatalf("Failed to create world: %v", err)
	}
	world.zones[1] = zone
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

	// Ensure the mob's RoomVNUM is set correctly
	mob.RoomVNUM = room.VNUM

	// Verify the mob is in the room
	if len(room.Characters) == 0 || room.Characters[0] != mob {
		t.Errorf("Expected mob to be in the room")
	}

	// Schedule the mob for respawning
	world.ScheduleMobRespawn(mob)

	// Verify the mob is no longer in the room
	mobFound := false
	for _, ch := range room.Characters {
		if ch == mob {
			mobFound = true
			break
		}
	}
	if mobFound {
		t.Errorf("Expected mob to be removed from the room")
	}

	// Verify the mob is scheduled for respawning
	if len(world.mobRespawns) == 0 {
		t.Errorf("Expected mob to be scheduled for respawning")
	} else {
		respawn := world.mobRespawns[0]
		if respawn.MobVNUM != 3060 {
			t.Errorf("Expected respawn mob VNUM to be 3060, got %d", respawn.MobVNUM)
		}
		if respawn.RoomVNUM != 3001 {
			t.Errorf("Expected respawn room VNUM to be 3001, got %d", respawn.RoomVNUM)
		}
	}

	// Set the respawn time to the past to trigger respawning
	if len(world.mobRespawns) > 0 {
		world.mobRespawns[0].RespawnTime = time.Now().Add(-time.Minute)
	} else {
		t.Fatalf("No mob respawn entries found")
	}

	// Process respawns
	world.ProcessMobRespawns()

	// Verify the mob has been respawned
	if len(room.Characters) == 0 {
		t.Errorf("Expected mob to be respawned in the room")
	} else {
		respawnedMob := room.Characters[0]
		if respawnedMob.Prototype.VNUM != 3060 {
			t.Errorf("Expected respawned mob VNUM to be 3060, got %d", respawnedMob.Prototype.VNUM)
		}
		
		// Check the respawned mob stats
		if respawnedMob.Level != 10 {
			t.Errorf("Expected level 10, got %d", respawnedMob.Level)
		}
		if respawnedMob.HitRoll != 10 {
			t.Errorf("Expected hitroll 10, got %d", respawnedMob.HitRoll)
		}
		if respawnedMob.ArmorClass[0] != 20 || respawnedMob.ArmorClass[1] != 20 || respawnedMob.ArmorClass[2] != 20 {
			t.Errorf("Expected AC [20 20 20], got %v", respawnedMob.ArmorClass)
		}
		if respawnedMob.DamRoll != 3 {
			t.Errorf("Expected damroll 3, got %d", respawnedMob.DamRoll)
		}
		if respawnedMob.Gold != 500 {
			t.Errorf("Expected gold 500, got %d", respawnedMob.Gold)
		}
		if respawnedMob.Experience != 9000 {
			t.Errorf("Expected experience 9000, got %d", respawnedMob.Experience)
		}
	}

	// Verify the respawn entry has been removed
	if len(world.mobRespawns) > 0 {
		t.Errorf("Expected respawn entry to be removed")
	}
}
