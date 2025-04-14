package command

import (
	"strings"
	"testing"

	"github.com/wltechblog/DikuGo/pkg/types"
)

// Mock world for testing
type MockWorld struct {
	rooms []*types.Room
}

func (w *MockWorld) GetRooms() []*types.Room {
	return w.rooms
}

func TestValidateRoomsCommand(t *testing.T) {
	// Create test rooms
	room1 := &types.Room{
		VNUM:        3001,
		Name:        "Temple of Midgaard",
		Description: "The temple of Midgaard is a magnificent structure with tall marble columns.",
		Exits: [6]*types.Exit{
			{
				Direction: types.DIR_NORTH,
				DestVnum:  3014,
			},
			{
				Direction: types.DIR_EAST,
				DestVnum:  3002,
			},
			nil,
			nil,
			nil,
			nil,
		},
	}

	room2 := &types.Room{
		VNUM:        3002,
		Name:        "East of Temple",
		Description: "You are east of the temple.",
		Exits: [6]*types.Exit{
			nil,
			nil,
			nil,
			{
				Direction: types.DIR_WEST,
				DestVnum:  3001,
			},
			nil,
			nil,
		},
	}

	room3 := &types.Room{
		VNUM:        3003,
		Name:        "Invalid Exit Room",
		Description: "This room has an invalid exit.",
		Exits: [6]*types.Exit{
			{
				Direction: types.DIR_NORTH,
				DestVnum:  9999, // Invalid room VNUM
			},
			nil,
			nil,
			nil,
			nil,
			nil,
		},
	}

	// Create a mock world
	mockWorld := &MockWorld{
		rooms: []*types.Room{room1, room2, room3},
	}

	// Create a test character
	character := &types.Character{
		Name:   "TestCharacter",
		World:  mockWorld,
		InRoom: room1,
	}

	// Create a validaterooms command
	validateCmd := &ValidateRoomsCommand{}

	// Test the command
	err := validateCmd.Execute(character, "")
	if err == nil {
		t.Errorf("Expected an error (which contains the validation results), got nil")
	} else {
		// The error should contain the validation results
		if !strings.Contains(err.Error(), "Validating all room exits") {
			t.Errorf("Expected error to contain validation header, got: %s", err.Error())
		}
		if !strings.Contains(err.Error(), "ERROR: Room 3003 has exit in direction north with DestVnum 9999") {
			t.Errorf("Expected error to report invalid exit, got: %s", err.Error())
		}
		if !strings.Contains(err.Error(), "Total exits:") {
			t.Errorf("Expected error to report total exits, got: %s", err.Error())
		}
		if !strings.Contains(err.Error(), "Valid exits:") {
			t.Errorf("Expected error to report valid exits, got: %s", err.Error())
		}
		if !strings.Contains(err.Error(), "Invalid exits:") {
			t.Errorf("Expected error to report invalid exits, got: %s", err.Error())
		}
	}
}
