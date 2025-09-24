package world

import (
	"testing"

	"github.com/wltechblog/DikuGo/pkg/types"
)

func TestZoneResetMobLoading(t *testing.T) {
	// Create a mock storage
	storage := NewMockStorage()

	// Create test rooms
	room1 := &types.Room{
		VNUM:        3001,
		Name:        "Test Room 1",
		Description: "This is test room 1.",
		Characters:  make([]*types.Character, 0),
		Objects:     make([]*types.ObjectInstance, 0),
	}
	storage.rooms = append(storage.rooms, room1)

	// Create a test mobile prototype
	mobProto := &types.Mobile{
		VNUM:        1001,
		Name:        "test mob",
		ShortDesc:   "a test mob",
		LongDesc:    "A test mob is standing here.",
		Description: "This is a test mob for unit testing.",
		ActFlags:    types.ACT_ISNPC,
		Level:       5,
	}
	storage.mobiles = append(storage.mobiles, mobProto)

	// Create a test zone with reset commands
	zone := &types.Zone{
		VNUM:      1,
		Name:      "Test Zone",
		Lifespan:  30,
		Age:       0, // Not ready to reset yet
		ResetMode: 1,
		MinVNUM:   3000,
		MaxVNUM:   3999,
		Commands: []*types.ZoneCommand{
			{
				Command: 'M',  // Load mobile
				IfFlag:  0,    // Always execute
				Arg1:    1001, // Mobile VNUM
				Arg2:    1,    // Max number
				Arg3:    3001, // Room VNUM
			},
		},
	}
	storage.zones = append(storage.zones, zone)
	room1.Zone = zone

	// Create a test world
	world, err := NewWorld(nil, storage)
	if err != nil {
		t.Fatalf("Failed to create world: %v", err)
	}

	// Manually add the room to the world's rooms map
	world.rooms[room1.VNUM] = room1

	// We expect one mob in the room initially due to the initial zone reset
	if len(room1.Characters) != 1 {
		t.Errorf("Expected one mob in room initially due to initial zone reset, got %d", len(room1.Characters))
	}

	// Set the zone age to trigger a reset
	zone.Age = zone.Lifespan

	// Reset the zone
	world.ResetZones()

	// Verify a mob was loaded into the room
	if len(room1.Characters) == 0 {
		t.Errorf("Expected a mob to be loaded into the room after zone reset")
	} else {
		// Verify it's the correct mob
		mob := room1.Characters[0]
		if !mob.IsNPC {
			t.Errorf("Expected mob to be an NPC")
		}
		if mob.Prototype == nil || mob.Prototype.VNUM != 1001 {
			t.Errorf("Expected mob to have prototype VNUM 1001, got %v", mob.Prototype)
		}
	}

	// Reset the zone again, should not load another mob since we've reached the max
	zone.Age = zone.Lifespan
	world.ResetZones()

	// Verify still only one mob in the room
	if len(room1.Characters) != 1 {
		t.Errorf("Expected only one mob in room after second reset, got %d", len(room1.Characters))
	}

	// Kill the mob and verify it gets respawned
	mob := room1.Characters[0]
	world.RemoveCharacter(mob)

	// Verify no mobs in the room
	if len(room1.Characters) > 0 {
		t.Errorf("Expected no mobs in room after removal, got %d", len(room1.Characters))
	}

	// Reset the zone again
	zone.Age = zone.Lifespan
	world.ResetZones()

	// Verify a mob was loaded into the room again
	if len(room1.Characters) == 0 {
		t.Errorf("Expected a mob to be loaded into the room after third zone reset")
	} else {
		// Verify it's the correct mob
		mob := room1.Characters[0]
		if !mob.IsNPC {
			t.Errorf("Expected mob to be an NPC")
		}
		if mob.Prototype == nil || mob.Prototype.VNUM != 1001 {
			t.Errorf("Expected mob to have prototype VNUM 1001, got %v", mob.Prototype)
		}
	}
}

func TestZoneResetEquipment(t *testing.T) {
	// Create a mock storage
	storage := NewMockStorage()

	// Create a test mob prototype
	mobProto := &types.Mobile{
		VNUM:      7009,
		Name:      "grand knight paladin",
		ShortDesc: "the Grand Knight of Paladins",
		LongDesc:  "The Grand Knight of Paladins is standing here.",
		Level:     20,
		ActFlags:  types.ACT_ISNPC,
	}
	storage.mobiles = append(storage.mobiles, mobProto)

	// Create a test object prototype
	objProto := &types.Object{
		VNUM:        7212,
		Name:        "sword bastard",
		ShortDesc:   "a bastard sword",
		Description: "A bastard sword is lying here.",
	}
	storage.objects = append(storage.objects, objProto)

	// Create a test room
	room := &types.Room{
		VNUM:        7900,
		Name:        "Test Room",
		Description: "This is a test room.",
		Characters:  make([]*types.Character, 0),
		Objects:     make([]*types.ObjectInstance, 0),
	}
	storage.rooms = append(storage.rooms, room)

	// Create a test zone with commands
	zone := &types.Zone{
		VNUM: 79,
		Name: "Test Zone",
		Commands: []*types.ZoneCommand{
			{
				Command: 'M',
				Arg1:    mobProto.VNUM,
				Arg2:    1,
				Arg3:    room.VNUM,
			}, // Load mob
			{
				Command: 'E',
				Arg1:    objProto.VNUM,
				Arg2:    mobProto.VNUM,
				Arg3:    16,
			}, // Equip mob with object in position 16 (WIELD)
		},
	}

	storage.zones = append(storage.zones, zone)
	room.Zone = zone

	// Create a world with our mock storage
	world, err := NewWorld(nil, storage)
	if err != nil {
		t.Fatalf("Failed to create world: %v", err)
	}

	// Manually add the room to the world's rooms map
	world.rooms[room.VNUM] = room

	// Reset the zone
	world.resetZone(zone)

	// Find the mob in the room
	var mob *types.Character
	for _, ch := range room.Characters {
		if ch.IsNPC && ch.Prototype != nil && ch.Prototype.VNUM == mobProto.VNUM {
			mob = ch
			break
		}
	}

	// Check if the mob was loaded
	if mob == nil {
		t.Fatalf("Mob was not loaded")
	}

	// Check if the mob has the equipment
	if mob.Equipment[16] == nil {
		t.Fatalf("Mob does not have equipment in position 16")
	}

	// Check if the equipment is the correct object
	if mob.Equipment[16].Prototype.VNUM != objProto.VNUM {
		t.Fatalf("Mob has wrong equipment: expected VNUM %d, got VNUM %d",
			objProto.VNUM, mob.Equipment[16].Prototype.VNUM)
	}

	t.Logf("Mob %s (VNUM %d) has equipment %s (VNUM %d) in position %d",
		mob.Name, mob.Prototype.VNUM, mob.Equipment[16].Prototype.Name,
		mob.Equipment[16].Prototype.VNUM, 16)
}
