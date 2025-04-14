package storage

import (
	"os"
	"testing"

	"github.com/wltechblog/DikuGo/pkg/types"
)

func TestRoomParserExtraDescriptions(t *testing.T) {
	// Create a sample room file with extra descriptions
	roomData := `#3001
Temple of Midgaard~
The temple of Midgaard is a magnificent structure with tall marble columns.
A large altar stands in the center of the room, and a fountain bubbles nearby.
~
30 8 0
D0
The entrance to the temple leads north to the street.
~
~
0 0 3014
E
fountain~
A beautiful marble fountain with crystal clear water.
~
E
altar~
A large stone altar stands in the center of the room.
~
S
`

	// Create a temporary file for testing
	tmpFile, err := os.CreateTemp("", "room_test*.wld")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	// Write the test data to the file
	if _, err := tmpFile.WriteString(roomData); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	// Parse the room
	rooms, err := ParseRooms(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to parse room: %v", err)
	}

	// Check if we got the expected number of rooms
	if len(rooms) != 1 {
		t.Fatalf("Expected 1 room, got %d", len(rooms))
	}

	// Get the room
	room := rooms[0]

	// Check if the room has the expected VNUM
	if room.VNUM != 3001 {
		t.Errorf("Expected room VNUM 3001, got %d", room.VNUM)
	}

	// Check if the room has the expected name
	if room.Name != "Temple of Midgaard" {
		t.Errorf("Expected room name 'Temple of Midgaard', got '%s'", room.Name)
	}

	// Check if the room has the expected number of extra descriptions
	if len(room.ExtraDescs) != 2 {
		t.Errorf("Expected 2 extra descriptions, got %d", len(room.ExtraDescs))
	}

	// Check the first extra description
	if room.ExtraDescs[0].Keywords != "fountain" {
		t.Errorf("Expected first extra description keywords 'fountain', got '%s'", room.ExtraDescs[0].Keywords)
	}
	if room.ExtraDescs[0].Description != "A beautiful marble fountain with crystal clear water." {
		t.Errorf("Expected first extra description text 'A beautiful marble fountain with crystal clear water.', got '%s'", room.ExtraDescs[0].Description)
	}

	// Check the second extra description
	if room.ExtraDescs[1].Keywords != "altar" {
		t.Errorf("Expected second extra description keywords 'altar', got '%s'", room.ExtraDescs[1].Keywords)
	}
	if room.ExtraDescs[1].Description != "A large stone altar stands in the center of the room." {
		t.Errorf("Expected second extra description text 'A large stone altar stands in the center of the room.', got '%s'", room.ExtraDescs[1].Description)
	}

	// Check if the room has the expected exit
	if room.Exits[types.DIR_NORTH] == nil {
		t.Errorf("Expected north exit, got nil")
	} else {
		if room.Exits[types.DIR_NORTH].DestVnum != 3014 {
			t.Errorf("Expected north exit to lead to room 3014, got %d", room.Exits[types.DIR_NORTH].DestVnum)
		}
		if room.Exits[types.DIR_NORTH].Description != "The entrance to the temple leads north to the street." {
			t.Errorf("Expected north exit description 'The entrance to the temple leads north to the street.', got '%s'", room.Exits[types.DIR_NORTH].Description)
		}
	}
}
