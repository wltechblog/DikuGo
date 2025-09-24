package world

import (
	"testing"

	"github.com/wltechblog/DikuGo/pkg/types"
)

// TestCharacterPlacementFix tests that new characters are placed in room 3001, not room 0
func TestCharacterPlacementFix(t *testing.T) {
	// Create a mock storage
	testStorage := NewMockStorage()

	// Create test rooms including room 0 (The Void) and room 3001 (Temple of Midgaard)
	room0 := &types.Room{
		VNUM:        0,
		Name:        "The Void",
		Description: "You are floating in nothing.",
		Characters:  make([]*types.Character, 0),
	}
	room3001 := &types.Room{
		VNUM:        3001,
		Name:        "Temple of Midgaard",
		Description: "This is the temple of Midgaard.",
		Characters:  make([]*types.Character, 0),
	}

	testStorage.rooms = append(testStorage.rooms, room0, room3001)

	// Create a test world
	world, err := NewWorld(nil, testStorage)
	if err != nil {
		t.Fatalf("Failed to create world: %v", err)
	}

	// Test 1: New character with RoomVNUM = -1 should go to room 3001
	newChar := &types.Character{
		Name:     "TestNewPlayer",
		IsNPC:    false,
		RoomVNUM: -1, // New character, no saved room
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

	// Test 2: Save the character and verify RoomVNUM is preserved correctly
	err = world.SaveCharacter(newChar)
	if err != nil {
		t.Fatalf("Failed to save character: %v", err)
	}

	// The character should now have RoomVNUM = 3001 (not 0)
	if newChar.RoomVNUM != 3001 {
		t.Errorf("After saving, expected character RoomVNUM to be 3001, got %d", newChar.RoomVNUM)
	}

	// Test 3: Create another new character and verify it doesn't get saved with RoomVNUM = 0
	anotherNewChar := &types.Character{
		Name:     "AnotherNewPlayer",
		IsNPC:    false,
		RoomVNUM: -1, // New character, no saved room
	}

	// Save the character before adding to world (simulating character creation process)
	err = world.SaveCharacter(anotherNewChar)
	if err != nil {
		t.Fatalf("Failed to save new character: %v", err)
	}

	// The character should still have RoomVNUM = -1 (not 0)
	if anotherNewChar.RoomVNUM != -1 {
		t.Errorf("After saving new character, expected RoomVNUM to remain -1, got %d", anotherNewChar.RoomVNUM)
	}

	// Now add to world and verify placement
	world.AddCharacter(anotherNewChar)

	if anotherNewChar.InRoom == nil {
		t.Fatalf("Another new character was not placed in any room")
	}

	if anotherNewChar.InRoom.VNUM != 3001 {
		t.Errorf("Expected another new character to be placed in room 3001, got room %d", anotherNewChar.InRoom.VNUM)
	}

	// Test 4: Character with RoomVNUM = 0 should be placed in room 0 (existing saved character)
	existingChar := &types.Character{
		Name:     "ExistingPlayer",
		IsNPC:    false,
		RoomVNUM: 0, // Saved in The Void
	}

	world.AddCharacter(existingChar)

	if existingChar.InRoom == nil {
		t.Fatalf("Existing character was not placed in any room")
	}

	if existingChar.InRoom.VNUM != 0 {
		t.Errorf("Expected existing character to be placed in room 0, got room %d", existingChar.InRoom.VNUM)
	}

	t.Log("All character placement tests passed!")
}
