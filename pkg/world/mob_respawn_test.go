package world

import (
	"testing"
	"time"

	"github.com/wltechblog/DikuGo/pkg/types"
)

func TestMobRespawn(t *testing.T) {
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

	// Create a test mobile prototype
	mobProto := &types.Mobile{
		VNUM:        1001,
		Name:        "testmob",
		ShortDesc:   "a test mob",
		LongDesc:    "A test mob is standing here.",
		Description: "This is a test mob for unit testing.",
		ActFlags:    types.ACT_ISNPC,
		Level:       5,
		Position:    types.POS_STANDING,
		Gold:        100,
		Experience:  500,
		AC:          [3]int{5, 5, 5},
		HitRoll:     2,
		DamRoll:     2,
		Dice:        [3]int{2, 8, 5}, // 2d8+5 hit points
	}
	world.mobiles[1001] = mobProto

	// Create a mob from the prototype
	mob := world.CreateMobFromPrototype(1001, room)
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

	// Manually remove the mob from the room for testing purposes
	for i, ch := range room.Characters {
		if ch == mob {
			room.Characters = append(room.Characters[:i], room.Characters[i+1:]...)
			break
		}
	}

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
		if respawn.MobVNUM != 1001 {
			t.Errorf("Expected respawn mob VNUM to be 1001, got %d", respawn.MobVNUM)
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
		if respawnedMob.Prototype.VNUM != 1001 {
			t.Errorf("Expected respawned mob VNUM to be 1001, got %d", respawnedMob.Prototype.VNUM)
		}
	}

	// Verify the respawn entry has been removed
	if len(world.mobRespawns) > 0 {
		t.Errorf("Expected respawn entry to be removed")
	}
}
