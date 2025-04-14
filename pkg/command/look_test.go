package command

import (
	"strings"
	"testing"

	"github.com/wltechblog/DikuGo/pkg/types"
)

func TestLookAtExtraDescription(t *testing.T) {
	// Create a test room with extra descriptions
	room := &types.Room{
		VNUM:        3001,
		Name:        "Temple of Midgaard",
		Description: "The temple of Midgaard is a magnificent structure with tall marble columns.",
		ExtraDescs: []*types.ExtraDescription{
			{
				Keywords:    "fountain",
				Description: "A beautiful marble fountain with crystal clear water.",
			},
			{
				Keywords:    "altar",
				Description: "A large stone altar stands in the center of the room.",
			},
		},
	}

	// Create a test character
	character := &types.Character{
		Name:   "TestCharacter",
		InRoom: room,
	}

	// Create a look command
	lookCmd := &LookCommand{}

	// Test looking at the fountain
	err := lookCmd.lookAtTarget(character, "fountain")
	if err == nil {
		t.Errorf("Expected an error (which contains the description), got nil")
	} else {
		// The error should contain the fountain description
		if !strings.Contains(err.Error(), "A beautiful marble fountain with crystal clear water.") {
			t.Errorf("Expected error to contain fountain description, got: %s", err.Error())
		}
	}

	// Test looking at the altar
	err = lookCmd.lookAtTarget(character, "altar")
	if err == nil {
		t.Errorf("Expected an error (which contains the description), got nil")
	} else {
		// The error should contain the altar description
		if !strings.Contains(err.Error(), "A large stone altar stands in the center of the room.") {
			t.Errorf("Expected error to contain altar description, got: %s", err.Error())
		}
	}

	// Test looking at something that doesn't exist
	err = lookCmd.lookAtTarget(character, "nonexistent")
	if err == nil {
		t.Errorf("Expected an error for nonexistent target, got nil")
	} else {
		// The error should indicate that the target doesn't exist
		if !strings.Contains(err.Error(), "you don't see that here") {
			t.Errorf("Expected error to indicate target doesn't exist, got: %s", err.Error())
		}
	}
}

func TestLookInDirection(t *testing.T) {
	// Create a test room with an exit
	room := &types.Room{
		VNUM:        3001,
		Name:        "Temple of Midgaard",
		Description: "The temple of Midgaard is a magnificent structure with tall marble columns.",
		Exits: [6]*types.Exit{
			{
				Direction:   types.DIR_NORTH,
				Description: "The entrance to the temple leads north to the street.",
				DestVnum:    3014,
			},
			nil,
			nil,
			nil,
			nil,
			nil,
		},
	}

	// Create a test character
	character := &types.Character{
		Name:   "TestCharacter",
		InRoom: room,
	}

	// Create a look command
	lookCmd := &LookCommand{}

	// Test looking north
	err := lookCmd.lookInDirection(character, types.DIR_NORTH)
	if err == nil {
		t.Errorf("Expected an error (which contains the description), got nil")
	} else {
		// The error should contain the exit description
		if !strings.Contains(err.Error(), "The entrance to the temple leads north to the street.") {
			t.Errorf("Expected error to contain exit description, got: %s", err.Error())
		}
	}

	// Test looking in a direction with no exit
	err = lookCmd.lookInDirection(character, types.DIR_EAST)
	if err == nil {
		t.Errorf("Expected an error for nonexistent exit, got nil")
	} else {
		// The error should indicate that there's nothing special in that direction
		if !strings.Contains(err.Error(), "you see nothing special in that direction") {
			t.Errorf("Expected error to indicate nothing special in that direction, got: %s", err.Error())
		}
	}
}
