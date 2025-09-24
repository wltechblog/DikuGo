package world

import (
	"testing"

	"github.com/wltechblog/DikuGo/pkg/types"
)

func TestNewCharacterPlacement(t *testing.T) {
	// Create a mock storage
	testStorage := NewMockStorage()

	// Create room 3001 (Temple of Midgaard)
	room3001 := &types.Room{
		VNUM:        3001,
		Name:        "Temple of Midgaard",
		Description: "The temple of Midgaard is a magnificent structure.",
		SectorType:  0,
		Flags:       0,
		Characters:  make([]*types.Character, 0),
	}
	testStorage.rooms = append(testStorage.rooms, room3001)

	// Create room 0 (The Void) as fallback
	room0 := &types.Room{
		VNUM:        0,
		Name:        "The Void",
		Description: "You are floating in nothing.",
		SectorType:  1,
		Flags:       8,
		Characters:  make([]*types.Character, 0),
	}
	testStorage.rooms = append(testStorage.rooms, room0)

	// Create a test world
	world, err := NewWorld(nil, testStorage)
	if err != nil {
		t.Fatalf("Failed to create world: %v", err)
	}

	// Test 1: New character with no saved room should go to room 3001
	newChar := &types.Character{
		Name:     "TestPlayer",
		IsNPC:    false,
		RoomVNUM: -1, // No saved room
	}

	world.AddCharacter(newChar)

	if newChar.InRoom == nil {
		t.Fatalf("New character was not placed in any room")
	}

	if newChar.InRoom.VNUM != 3001 {
		t.Errorf("Expected new character to be placed in room 3001, got room %d", newChar.InRoom.VNUM)
	}

	if newChar.RoomVNUM != 3001 {
		t.Errorf("Expected character RoomVNUM to be 3001, got %d", newChar.RoomVNUM)
	}

	// Test 2: Character with saved room should go to that room
	savedChar := &types.Character{
		Name:     "SavedPlayer",
		IsNPC:    false,
		RoomVNUM: 0, // Saved in room 0
	}

	world.AddCharacter(savedChar)

	if savedChar.InRoom == nil {
		t.Fatalf("Saved character was not placed in any room")
	}

	if savedChar.InRoom.VNUM != 0 {
		t.Errorf("Expected saved character to be placed in room 0, got room %d", savedChar.InRoom.VNUM)
	}

	// Test 3: Character should be in room's character list
	found := false
	for _, ch := range newChar.InRoom.Characters {
		if ch == newChar {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("New character not found in room 3001's character list")
	}

	// Test 4: Test fallback when room 3001 doesn't exist
	testStorage2 := NewMockStorage()
	testStorage2.rooms = append(testStorage2.rooms, room0) // Only room 0 exists

	world2, err := NewWorld(nil, testStorage2)
	if err != nil {
		t.Fatalf("Failed to create world2: %v", err)
	}

	fallbackChar := &types.Character{
		Name:     "FallbackPlayer",
		IsNPC:    false,
		RoomVNUM: -1, // No saved room
	}

	world2.AddCharacter(fallbackChar)

	if fallbackChar.InRoom == nil {
		t.Fatalf("Fallback character was not placed in any room")
	}

	if fallbackChar.InRoom.VNUM != 0 {
		t.Errorf("Expected fallback character to be placed in room 0, got room %d", fallbackChar.InRoom.VNUM)
	}
}

func TestCharacterMovement(t *testing.T) {
	// Create a mock storage
	testStorage := NewMockStorage()

	// Create test rooms
	room1 := &types.Room{
		VNUM:        1,
		Name:        "Room 1",
		Description: "This is room 1.",
		Characters:  make([]*types.Character, 0),
	}
	room2 := &types.Room{
		VNUM:        2,
		Name:        "Room 2",
		Description: "This is room 2.",
		Characters:  make([]*types.Character, 0),
	}

	testStorage.rooms = append(testStorage.rooms, room1, room2)

	// Create a test world
	world, err := NewWorld(nil, testStorage)
	if err != nil {
		t.Fatalf("Failed to create world: %v", err)
	}

	// Create a test character
	character := &types.Character{
		Name:     "TestPlayer",
		IsNPC:    false,
		InRoom:   room1,
		RoomVNUM: 1,
	}
	room1.Characters = append(room1.Characters, character)

	// Test movement from room 1 to room 2
	world.CharacterMove(character, room2)

	// Check that character is now in room 2
	if character.InRoom != room2 {
		t.Errorf("Expected character to be in room 2, but InRoom is %v", character.InRoom)
	}

	if character.RoomVNUM != 2 {
		t.Errorf("Expected character RoomVNUM to be 2, got %d", character.RoomVNUM)
	}

	// Check that character is in room 2's character list
	found := false
	for _, ch := range room2.Characters {
		if ch == character {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Character not found in room 2's character list")
	}

	// Check that character is not in room 1's character list
	for _, ch := range room1.Characters {
		if ch == character {
			t.Errorf("Character still found in room 1's character list after move")
		}
	}
}
