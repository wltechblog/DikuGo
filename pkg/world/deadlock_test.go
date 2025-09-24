package world

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/wltechblog/DikuGo/pkg/types"
)

// TestMovementDeadlock tests that concurrent mob and player movement doesn't cause deadlocks
func TestMovementDeadlock(t *testing.T) {
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
	room3 := &types.Room{
		VNUM:        3,
		Name:        "Room 3",
		Description: "This is room 3.",
		Characters:  make([]*types.Character, 0),
	}

	testStorage.rooms = append(testStorage.rooms, room1, room2, room3)

	// Create a test world
	world, err := NewWorld(nil, testStorage)
	if err != nil {
		t.Fatalf("Failed to create world: %v", err)
	}

	// Create test characters (mobs and players)
	mob1 := &types.Character{
		Name:     "TestMob1",
		IsNPC:    true,
		InRoom:   room1,
		RoomVNUM: 1,
	}
	mob2 := &types.Character{
		Name:     "TestMob2",
		IsNPC:    true,
		InRoom:   room2,
		RoomVNUM: 2,
	}
	player1 := &types.Character{
		Name:     "TestPlayer1",
		IsNPC:    false,
		InRoom:   room3,
		RoomVNUM: 3,
	}
	player2 := &types.Character{
		Name:     "TestPlayer2",
		IsNPC:    false,
		InRoom:   room1,
		RoomVNUM: 1,
	}

	// Add characters to rooms
	room1.Characters = append(room1.Characters, mob1, player2)
	room2.Characters = append(room2.Characters, mob2)
	room3.Characters = append(room3.Characters, player1)

	// Test concurrent movement operations
	var wg sync.WaitGroup
	const numOperations = 100

	// Channel to signal completion or timeout
	done := make(chan bool, 1)

	// Start concurrent movement operations
	wg.Add(4)

	// Mob1 moving between room1 and room2
	go func() {
		defer wg.Done()
		for i := 0; i < numOperations; i++ {
			if i%2 == 0 {
				world.CharacterMove(mob1, room2)
			} else {
				world.CharacterMove(mob1, room1)
			}
			time.Sleep(1 * time.Millisecond) // Small delay to allow interleaving
		}
	}()

	// Mob2 moving between room2 and room3
	go func() {
		defer wg.Done()
		for i := 0; i < numOperations; i++ {
			if i%2 == 0 {
				world.CharacterMove(mob2, room3)
			} else {
				world.CharacterMove(mob2, room2)
			}
			time.Sleep(1 * time.Millisecond)
		}
	}()

	// Player1 moving between room3 and room1
	go func() {
		defer wg.Done()
		for i := 0; i < numOperations; i++ {
			if i%2 == 0 {
				world.CharacterMove(player1, room1)
			} else {
				world.CharacterMove(player1, room3)
			}
			time.Sleep(1 * time.Millisecond)
		}
	}()

	// Player2 moving between room1 and room2
	go func() {
		defer wg.Done()
		for i := 0; i < numOperations; i++ {
			if i%2 == 0 {
				world.CharacterMove(player2, room2)
			} else {
				world.CharacterMove(player2, room1)
			}
			time.Sleep(1 * time.Millisecond)
		}
	}()

	// Wait for all operations to complete or timeout
	go func() {
		wg.Wait()
		done <- true
	}()

	// Set a timeout to detect deadlocks
	select {
	case <-done:
		// All operations completed successfully
		t.Log("All movement operations completed without deadlock")
	case <-time.After(10 * time.Second):
		t.Fatal("Test timed out - likely deadlock detected")
	}
}

// TestConcurrentRoomAccess tests concurrent access to room character lists
func TestConcurrentRoomAccess(t *testing.T) {
	// Create a mock storage
	testStorage := NewMockStorage()

	// Create test room
	room := &types.Room{
		VNUM:        1,
		Name:        "Test Room",
		Description: "This is a test room.",
		Characters:  make([]*types.Character, 0),
	}

	testStorage.rooms = append(testStorage.rooms, room)

	// Create a test world
	world, err := NewWorld(nil, testStorage)
	if err != nil {
		t.Fatalf("Failed to create world: %v", err)
	}

	// Create multiple characters
	characters := make([]*types.Character, 10)
	for i := 0; i < 10; i++ {
		characters[i] = &types.Character{
			Name:     fmt.Sprintf("TestChar%d", i),
			IsNPC:    i%2 == 0, // Mix of NPCs and players
			InRoom:   room,
			RoomVNUM: 1,
		}
		room.Characters = append(room.Characters, characters[i])
	}

	var wg sync.WaitGroup
	const numOperations = 50

	// Channel to signal completion or timeout
	done := make(chan bool, 1)

	// Start concurrent operations that access room character lists
	wg.Add(3)

	// Goroutine 1: Move characters in and out of room
	go func() {
		defer wg.Done()
		for i := 0; i < numOperations; i++ {
			char := characters[i%len(characters)]
			// Move character out and back in
			world.CharacterMove(char, nil)  // Move to nil (limbo)
			world.CharacterMove(char, room) // Move back to room
			time.Sleep(1 * time.Millisecond)
		}
	}()

	// Goroutine 2: Read room character list
	go func() {
		defer wg.Done()
		for i := 0; i < numOperations*2; i++ {
			room.RLock()
			_ = len(room.Characters) // Read operation
			room.RUnlock()
			time.Sleep(500 * time.Microsecond)
		}
	}()

	// Goroutine 3: Use GetCharacterInRoom (which also reads the list)
	go func() {
		defer wg.Done()
		for i := 0; i < numOperations; i++ {
			world.GetCharacterInRoom(room, "TestChar0")
			time.Sleep(1 * time.Millisecond)
		}
	}()

	// Wait for all operations to complete or timeout
	go func() {
		wg.Wait()
		done <- true
	}()

	// Set a timeout to detect deadlocks
	select {
	case <-done:
		// All operations completed successfully
		t.Log("All concurrent room access operations completed without deadlock")
	case <-time.After(5 * time.Second):
		t.Fatal("Test timed out - likely deadlock detected in room access")
	}
}
