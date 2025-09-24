package world

import (
	"testing"

	"github.com/wltechblog/DikuGo/pkg/types"
)

func TestHandleCharacterDeath_Player(t *testing.T) {
	// Create a test world
	world := &World{
		rooms:      make(map[int]*types.Room),
		characters: make(map[string]*types.Character),
	}

	// Create a test room
	room := &types.Room{
		VNUM:       3001,
		Name:       "Test Room",
		Characters: make([]*types.Character, 0),
		Objects:    make([]*types.ObjectInstance, 0),
	}
	world.rooms[3001] = room

	// Create a test player character
	player := &types.Character{
		Name:      "TestPlayer",
		IsNPC:     false,
		HP:        0,
		Position:  types.POS_DEAD,
		InRoom:    room,
		RoomVNUM:  3001,
		Inventory: make([]*types.ObjectInstance, 0),
		Equipment: make([]*types.ObjectInstance, types.NUM_WEARS),
		Messages:  make([]string, 0),
	}

	// Add player to room
	room.Characters = append(room.Characters, player)
	world.characters[player.Name] = player

	// Handle character death
	world.HandleCharacterDeath(player)

	// Verify player was reset for resurrection
	if player.HP != 1 {
		t.Errorf("Expected player HP to be reset to 1, got %d", player.HP)
	}

	if player.Position != types.POS_STANDING {
		t.Errorf("Expected player position to be STANDING, got %d", player.Position)
	}

	// Verify RETURN_TO_MENU message was sent
	if !player.HasMessage("RETURN_TO_MENU") {
		t.Error("Expected player to have RETURN_TO_MENU message")
	}
}

func TestHandleCharacterDeath_NPC(t *testing.T) {
	// Create a test world
	world := &World{
		rooms:       make(map[int]*types.Room),
		characters:  make(map[string]*types.Character),
		mobRespawns: make([]*types.MobRespawn, 0),
		zones:       make(map[int]*types.Zone),
	}

	// Create a test zone
	zone := &types.Zone{
		VNUM:     30,
		MinVNUM:  3000,
		MaxVNUM:  3099,
		Lifespan: 15,
	}
	world.zones[30] = zone

	// Create a test room
	room := &types.Room{
		VNUM:       3001,
		Name:       "Test Room",
		Characters: make([]*types.Character, 0),
		Objects:    make([]*types.ObjectInstance, 0),
		Zone:       zone,
	}
	world.rooms[3001] = room

	// Create a mobile prototype
	prototype := &types.Mobile{
		VNUM:      3001,
		Name:      "test mob",
		ShortDesc: "a test mob",
		Level:     5,
	}

	// Create a test NPC character
	npc := &types.Character{
		Name:      "test mob",
		IsNPC:     true,
		HP:        0,
		Position:  types.POS_DEAD,
		InRoom:    room,
		RoomVNUM:  3001,
		Prototype: prototype,
		Inventory: make([]*types.ObjectInstance, 0),
		Equipment: make([]*types.ObjectInstance, types.NUM_WEARS),
	}

	// Add NPC to room and world
	room.Characters = append(room.Characters, npc)
	world.characters[npc.Name] = npc

	// Handle character death
	world.HandleCharacterDeath(npc)

	// Verify NPC was scheduled for respawn
	if len(world.mobRespawns) != 1 {
		t.Errorf("Expected 1 mob respawn scheduled, got %d", len(world.mobRespawns))
	}

	// Verify NPC was removed from world
	if _, exists := world.characters[npc.Name]; exists {
		t.Error("Expected NPC to be removed from world characters map")
	}
}

func TestHandleCharacterDeath_NilCharacter(t *testing.T) {
	// Create a test world
	world := &World{}

	// Handle death of nil character - should not panic
	world.HandleCharacterDeath(nil)

	// Test passes if no panic occurs
}
