package command

import (
	"strings"
	"testing"

	"github.com/wltechblog/DikuGo/pkg/types"
)

func TestOpenCommand_Container(t *testing.T) {
	// Create a test character
	character := &types.Character{
		Name:      "TestPlayer",
		Level:     5,
		Position:  types.POS_STANDING,
		Inventory: make([]*types.ObjectInstance, 0),
	}

	// Create a test room
	room := &types.Room{
		VNUM:       1001,
		Name:       "Test Room",
		Characters: []*types.Character{character},
		Objects:    make([]*types.ObjectInstance, 0),
	}
	character.InRoom = room

	// Create a closed, closeable container
	chest := &types.ObjectInstance{
		Prototype: &types.Object{
			VNUM:      2001,
			Name:      "chest",
			ShortDesc: "a wooden chest",
			Type:      types.ITEM_CONTAINER,
			Value:     [4]int{100, types.CONT_CLOSEABLE | types.CONT_CLOSED, 0, 0},
		},
		Value: [4]int{100, types.CONT_CLOSEABLE | types.CONT_CLOSED, 0, 0},
	}
	character.Inventory = append(character.Inventory, chest)

	// Test opening the chest
	openCmd := &OpenCommand{}
	err := openCmd.Execute(character, "chest")
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	// Check that the chest is now open (CLOSED flag removed)
	if (chest.Value[1] & types.CONT_CLOSED) != 0 {
		t.Error("Expected chest to be open, but CLOSED flag is still set")
	}
}

func TestOpenCommand_AlreadyOpen(t *testing.T) {
	// Create a test character
	character := &types.Character{
		Name:      "TestPlayer",
		Level:     5,
		Position:  types.POS_STANDING,
		Inventory: make([]*types.ObjectInstance, 0),
	}

	// Create a test room
	room := &types.Room{
		VNUM:       1001,
		Name:       "Test Room",
		Characters: []*types.Character{character},
	}
	character.InRoom = room

	// Create an already open container
	chest := &types.ObjectInstance{
		Prototype: &types.Object{
			VNUM:      2001,
			Name:      "chest",
			ShortDesc: "a wooden chest",
			Type:      types.ITEM_CONTAINER,
			Value:     [4]int{100, types.CONT_CLOSEABLE, 0, 0}, // Not closed
		},
		Value: [4]int{100, types.CONT_CLOSEABLE, 0, 0},
	}
	character.Inventory = append(character.Inventory, chest)

	// Test opening already open chest
	openCmd := &OpenCommand{}
	err := openCmd.Execute(character, "chest")
	if err == nil || !strings.Contains(err.Error(), "already open") {
		t.Errorf("Expected 'already open' error, got: %v", err)
	}
}

func TestOpenCommand_Locked(t *testing.T) {
	// Create a test character
	character := &types.Character{
		Name:      "TestPlayer",
		Level:     5,
		Position:  types.POS_STANDING,
		Inventory: make([]*types.ObjectInstance, 0),
	}

	// Create a test room
	room := &types.Room{
		VNUM:       1001,
		Name:       "Test Room",
		Characters: []*types.Character{character},
	}
	character.InRoom = room

	// Create a locked container
	chest := &types.ObjectInstance{
		Prototype: &types.Object{
			VNUM:      2001,
			Name:      "chest",
			ShortDesc: "a wooden chest",
			Type:      types.ITEM_CONTAINER,
			Value:     [4]int{100, types.CONT_CLOSEABLE | types.CONT_CLOSED | types.CONT_LOCKED, 0, 0},
		},
		Value: [4]int{100, types.CONT_CLOSEABLE | types.CONT_CLOSED | types.CONT_LOCKED, 0, 0},
	}
	character.Inventory = append(character.Inventory, chest)

	// Test opening locked chest
	openCmd := &OpenCommand{}
	err := openCmd.Execute(character, "chest")
	if err == nil || !strings.Contains(err.Error(), "locked") {
		t.Errorf("Expected 'locked' error, got: %v", err)
	}
}

func TestOpenCommand_NotCloseable(t *testing.T) {
	// Create a test character
	character := &types.Character{
		Name:      "TestPlayer",
		Level:     5,
		Position:  types.POS_STANDING,
		Inventory: make([]*types.ObjectInstance, 0),
	}

	// Create a test room
	room := &types.Room{
		VNUM:       1001,
		Name:       "Test Room",
		Characters: []*types.Character{character},
	}
	character.InRoom = room

	// Create a non-closeable container
	chest := &types.ObjectInstance{
		Prototype: &types.Object{
			VNUM:      2001,
			Name:      "chest",
			ShortDesc: "a wooden chest",
			Type:      types.ITEM_CONTAINER,
			Value:     [4]int{100, types.CONT_CLOSED, 0, 0}, // Closed but not closeable
		},
		Value: [4]int{100, types.CONT_CLOSED, 0, 0},
	}
	character.Inventory = append(character.Inventory, chest)

	// Test opening non-closeable chest
	openCmd := &OpenCommand{}
	err := openCmd.Execute(character, "chest")
	if err == nil || !strings.Contains(err.Error(), "can't do that") {
		t.Errorf("Expected 'can't do that' error, got: %v", err)
	}
}

func TestOpenCommand_NotContainer(t *testing.T) {
	// Create a test character
	character := &types.Character{
		Name:      "TestPlayer",
		Level:     5,
		Position:  types.POS_STANDING,
		Inventory: make([]*types.ObjectInstance, 0),
	}

	// Create a test room
	room := &types.Room{
		VNUM:       1001,
		Name:       "Test Room",
		Characters: []*types.Character{character},
	}
	character.InRoom = room

	// Create a non-container object
	sword := &types.ObjectInstance{
		Prototype: &types.Object{
			VNUM:      3001,
			Name:      "sword",
			ShortDesc: "a steel sword",
			Type:      types.ITEM_WEAPON,
		},
	}
	character.Inventory = append(character.Inventory, sword)

	// Test opening non-container
	openCmd := &OpenCommand{}
	err := openCmd.Execute(character, "sword")
	if err == nil || !strings.Contains(err.Error(), "not a container") {
		t.Errorf("Expected 'not a container' error, got: %v", err)
	}
}

func TestCloseCommand_Container(t *testing.T) {
	// Create a test character
	character := &types.Character{
		Name:      "TestPlayer",
		Level:     5,
		Position:  types.POS_STANDING,
		Inventory: make([]*types.ObjectInstance, 0),
	}

	// Create a test room
	room := &types.Room{
		VNUM:       1001,
		Name:       "Test Room",
		Characters: []*types.Character{character},
		Objects:    make([]*types.ObjectInstance, 0),
	}
	character.InRoom = room

	// Create an open, closeable container
	chest := &types.ObjectInstance{
		Prototype: &types.Object{
			VNUM:      2001,
			Name:      "chest",
			ShortDesc: "a wooden chest",
			Type:      types.ITEM_CONTAINER,
			Value:     [4]int{100, types.CONT_CLOSEABLE, 0, 0}, // Open
		},
		Value: [4]int{100, types.CONT_CLOSEABLE, 0, 0},
	}
	character.Inventory = append(character.Inventory, chest)

	// Test closing the chest
	closeCmd := &CloseCommand{}
	err := closeCmd.Execute(character, "chest")
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	// Check that the chest is now closed (CLOSED flag set)
	if (chest.Value[1] & types.CONT_CLOSED) == 0 {
		t.Error("Expected chest to be closed, but CLOSED flag is not set")
	}
}

func TestCloseCommand_AlreadyClosed(t *testing.T) {
	// Create a test character
	character := &types.Character{
		Name:      "TestPlayer",
		Level:     5,
		Position:  types.POS_STANDING,
		Inventory: make([]*types.ObjectInstance, 0),
	}

	// Create a test room
	room := &types.Room{
		VNUM:       1001,
		Name:       "Test Room",
		Characters: []*types.Character{character},
	}
	character.InRoom = room

	// Create an already closed container
	chest := &types.ObjectInstance{
		Prototype: &types.Object{
			VNUM:      2001,
			Name:      "chest",
			ShortDesc: "a wooden chest",
			Type:      types.ITEM_CONTAINER,
			Value:     [4]int{100, types.CONT_CLOSEABLE | types.CONT_CLOSED, 0, 0},
		},
		Value: [4]int{100, types.CONT_CLOSEABLE | types.CONT_CLOSED, 0, 0},
	}
	character.Inventory = append(character.Inventory, chest)

	// Test closing already closed chest
	closeCmd := &CloseCommand{}
	err := closeCmd.Execute(character, "chest")
	if err == nil || !strings.Contains(err.Error(), "already closed") {
		t.Errorf("Expected 'already closed' error, got: %v", err)
	}
}

func TestOpenCommand_NoArgs(t *testing.T) {
	// Create a test character
	character := &types.Character{
		Name:     "TestPlayer",
		Level:    5,
		Position: types.POS_STANDING,
	}

	// Create a test room
	room := &types.Room{
		VNUM:       1001,
		Name:       "Test Room",
		Characters: []*types.Character{character},
	}
	character.InRoom = room

	// Test opening with no arguments
	openCmd := &OpenCommand{}
	err := openCmd.Execute(character, "")
	if err == nil || !strings.Contains(err.Error(), "open what") {
		t.Errorf("Expected 'open what' error, got: %v", err)
	}
}

func TestCloseCommand_NoArgs(t *testing.T) {
	// Create a test character
	character := &types.Character{
		Name:     "TestPlayer",
		Level:    5,
		Position: types.POS_STANDING,
	}

	// Create a test room
	room := &types.Room{
		VNUM:       1001,
		Name:       "Test Room",
		Characters: []*types.Character{character},
	}
	character.InRoom = room

	// Test closing with no arguments
	closeCmd := &CloseCommand{}
	err := closeCmd.Execute(character, "")
	if err == nil || !strings.Contains(err.Error(), "close what") {
		t.Errorf("Expected 'close what' error, got: %v", err)
	}
}

func TestOpenCommand_NotFound(t *testing.T) {
	// Create a test character
	character := &types.Character{
		Name:      "TestPlayer",
		Level:     5,
		Position:  types.POS_STANDING,
		Inventory: make([]*types.ObjectInstance, 0),
	}

	// Create a test room
	room := &types.Room{
		VNUM:       1001,
		Name:       "Test Room",
		Characters: []*types.Character{character},
		Objects:    make([]*types.ObjectInstance, 0),
	}
	character.InRoom = room

	// Test opening non-existent object
	openCmd := &OpenCommand{}
	err := openCmd.Execute(character, "nonexistent")
	if err == nil || !strings.Contains(err.Error(), "I see no") {
		t.Errorf("Expected 'I see no' error, got: %v", err)
	}
}

func TestOpenCommand_Door(t *testing.T) {
	// Create a test character
	character := &types.Character{
		Name:     "TestPlayer",
		Level:    5,
		Position: types.POS_STANDING,
	}

	// Create a test room with a closed door
	room := &types.Room{
		VNUM:       1001,
		Name:       "Test Room",
		Characters: []*types.Character{character},
		Exits:      [6]*types.Exit{},
	}
	character.InRoom = room

	// Create a closed door to the north
	room.Exits[types.DIR_NORTH] = &types.Exit{
		Direction:   types.DIR_NORTH,
		Description: "A wooden door leads north.",
		Keywords:    "door",
		Flags:       types.EX_ISDOOR | types.EX_CLOSED,
		Key:         -1,
		DestVnum:    1002,
	}

	// Test opening the door
	openCmd := &OpenCommand{}
	err := openCmd.Execute(character, "door")
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	// Check that the door is now open (CLOSED flag removed)
	if (room.Exits[types.DIR_NORTH].Flags & types.EX_CLOSED) != 0 {
		t.Error("Expected door to be open, but CLOSED flag is still set")
	}
}

func TestOpenCommand_DoorAlreadyOpen(t *testing.T) {
	// Create a test character
	character := &types.Character{
		Name:     "TestPlayer",
		Level:    5,
		Position: types.POS_STANDING,
	}

	// Create a test room with an open door
	room := &types.Room{
		VNUM:       1001,
		Name:       "Test Room",
		Characters: []*types.Character{character},
		Exits:      [6]*types.Exit{},
	}
	character.InRoom = room

	// Create an open door to the north
	room.Exits[types.DIR_NORTH] = &types.Exit{
		Direction:   types.DIR_NORTH,
		Description: "A wooden door leads north.",
		Keywords:    "door",
		Flags:       types.EX_ISDOOR, // Open (no CLOSED flag)
		Key:         -1,
		DestVnum:    1002,
	}

	// Test opening already open door
	openCmd := &OpenCommand{}
	err := openCmd.Execute(character, "door")
	if err == nil || !strings.Contains(err.Error(), "already open") {
		t.Errorf("Expected 'already open' error, got: %v", err)
	}
}

func TestOpenCommand_DoorLocked(t *testing.T) {
	// Create a test character
	character := &types.Character{
		Name:     "TestPlayer",
		Level:    5,
		Position: types.POS_STANDING,
	}

	// Create a test room with a locked door
	room := &types.Room{
		VNUM:       1001,
		Name:       "Test Room",
		Characters: []*types.Character{character},
		Exits:      [6]*types.Exit{},
	}
	character.InRoom = room

	// Create a locked door to the north
	room.Exits[types.DIR_NORTH] = &types.Exit{
		Direction:   types.DIR_NORTH,
		Description: "A wooden door leads north.",
		Keywords:    "door",
		Flags:       types.EX_ISDOOR | types.EX_CLOSED | types.EX_LOCKED,
		Key:         1001, // Has a key
		DestVnum:    1002,
	}

	// Test opening locked door
	openCmd := &OpenCommand{}
	err := openCmd.Execute(character, "door")
	if err == nil || !strings.Contains(err.Error(), "locked") {
		t.Errorf("Expected 'locked' error, got: %v", err)
	}
}

func TestOpenCommand_NotDoor(t *testing.T) {
	// Create a test character
	character := &types.Character{
		Name:     "TestPlayer",
		Level:    5,
		Position: types.POS_STANDING,
	}

	// Create a test room with a non-door exit
	room := &types.Room{
		VNUM:       1001,
		Name:       "Test Room",
		Characters: []*types.Character{character},
		Exits:      [6]*types.Exit{},
	}
	character.InRoom = room

	// Create a non-door exit to the north
	room.Exits[types.DIR_NORTH] = &types.Exit{
		Direction:   types.DIR_NORTH,
		Description: "A passage leads north.",
		Keywords:    "passage",
		Flags:       0, // No EX_ISDOOR flag
		Key:         -1,
		DestVnum:    1002,
	}

	// Test opening non-door
	openCmd := &OpenCommand{}
	err := openCmd.Execute(character, "passage")
	if err == nil || !strings.Contains(err.Error(), "impossible") {
		t.Errorf("Expected 'impossible' error, got: %v", err)
	}
}

func TestCloseCommand_Door(t *testing.T) {
	// Create a test character
	character := &types.Character{
		Name:     "TestPlayer",
		Level:    5,
		Position: types.POS_STANDING,
	}

	// Create a test room with an open door
	room := &types.Room{
		VNUM:       1001,
		Name:       "Test Room",
		Characters: []*types.Character{character},
		Exits:      [6]*types.Exit{},
	}
	character.InRoom = room

	// Create an open door to the north
	room.Exits[types.DIR_NORTH] = &types.Exit{
		Direction:   types.DIR_NORTH,
		Description: "A wooden door leads north.",
		Keywords:    "door",
		Flags:       types.EX_ISDOOR, // Open (no CLOSED flag)
		Key:         -1,
		DestVnum:    1002,
	}

	// Test closing the door
	closeCmd := &CloseCommand{}
	err := closeCmd.Execute(character, "door")
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	// Check that the door is now closed (CLOSED flag set)
	if (room.Exits[types.DIR_NORTH].Flags & types.EX_CLOSED) == 0 {
		t.Error("Expected door to be closed, but CLOSED flag is not set")
	}
}

func TestCloseCommand_DoorAlreadyClosed(t *testing.T) {
	// Create a test character
	character := &types.Character{
		Name:     "TestPlayer",
		Level:    5,
		Position: types.POS_STANDING,
	}

	// Create a test room with a closed door
	room := &types.Room{
		VNUM:       1001,
		Name:       "Test Room",
		Characters: []*types.Character{character},
		Exits:      [6]*types.Exit{},
	}
	character.InRoom = room

	// Create a closed door to the north
	room.Exits[types.DIR_NORTH] = &types.Exit{
		Direction:   types.DIR_NORTH,
		Description: "A wooden door leads north.",
		Keywords:    "door",
		Flags:       types.EX_ISDOOR | types.EX_CLOSED,
		Key:         -1,
		DestVnum:    1002,
	}

	// Test closing already closed door
	closeCmd := &CloseCommand{}
	err := closeCmd.Execute(character, "door")
	if err == nil || !strings.Contains(err.Error(), "already closed") {
		t.Errorf("Expected 'already closed' error, got: %v", err)
	}
}

func TestOpenCommand_DoorWithDirection(t *testing.T) {
	// Create a test character
	character := &types.Character{
		Name:     "TestPlayer",
		Level:    5,
		Position: types.POS_STANDING,
	}

	// Create a test room with multiple doors
	room := &types.Room{
		VNUM:       1001,
		Name:       "Test Room",
		Characters: []*types.Character{character},
		Exits:      [6]*types.Exit{},
	}
	character.InRoom = room

	// Create doors in multiple directions with same keyword
	room.Exits[types.DIR_NORTH] = &types.Exit{
		Direction:   types.DIR_NORTH,
		Description: "A wooden gate leads north.",
		Keywords:    "gate wooden",
		Flags:       types.EX_ISDOOR | types.EX_CLOSED,
		Key:         -1,
		DestVnum:    1002,
	}

	room.Exits[types.DIR_SOUTH] = &types.Exit{
		Direction:   types.DIR_SOUTH,
		Description: "An iron gate leads south.",
		Keywords:    "gate iron",
		Flags:       types.EX_ISDOOR | types.EX_CLOSED,
		Key:         -1,
		DestVnum:    1003,
	}

	// Test opening specific gate by direction
	openCmd := &OpenCommand{}
	err := openCmd.Execute(character, "gate north")
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	// Check that only the north gate is open
	if (room.Exits[types.DIR_NORTH].Flags & types.EX_CLOSED) != 0 {
		t.Error("Expected north gate to be open, but CLOSED flag is still set")
	}
	if (room.Exits[types.DIR_SOUTH].Flags & types.EX_CLOSED) == 0 {
		t.Error("Expected south gate to remain closed, but CLOSED flag was removed")
	}
}

func TestOpenCommand_KeywordMatching(t *testing.T) {
	// Create a test character
	character := &types.Character{
		Name:     "TestPlayer",
		Level:    5,
		Position: types.POS_STANDING,
	}

	// Create a test room with a door with multiple keywords
	room := &types.Room{
		VNUM:       1001,
		Name:       "Test Room",
		Characters: []*types.Character{character},
		Exits:      [6]*types.Exit{},
	}
	character.InRoom = room

	// Create a door with multiple keywords
	room.Exits[types.DIR_EAST] = &types.Exit{
		Direction:   types.DIR_EAST,
		Description: "A heavy wooden door with iron hinges leads east.",
		Keywords:    "door wooden heavy iron",
		Flags:       types.EX_ISDOOR | types.EX_CLOSED,
		Key:         -1,
		DestVnum:    1002,
	}

	// Test opening with different keyword matches
	testCases := []string{"door", "wooden", "heavy", "iron"}

	for i, keyword := range testCases {
		// Reset door to closed state
		room.Exits[types.DIR_EAST].Flags |= types.EX_CLOSED

		openCmd := &OpenCommand{}
		err := openCmd.Execute(character, keyword)
		if err != nil {
			t.Errorf("Test case %d: Expected no error for keyword '%s', got: %v", i, keyword, err)
		}

		// Check that the door is open
		if (room.Exits[types.DIR_EAST].Flags & types.EX_CLOSED) != 0 {
			t.Errorf("Test case %d: Expected door to be open for keyword '%s', but CLOSED flag is still set", i, keyword)
		}
	}
}

func TestOpenCommand_ArgumentParsing(t *testing.T) {
	// Create a test character
	character := &types.Character{
		Name:     "TestPlayer",
		Level:    5,
		Position: types.POS_STANDING,
	}

	// Create a test room
	room := &types.Room{
		VNUM:       1001,
		Name:       "Test Room",
		Characters: []*types.Character{character},
		Exits:      [6]*types.Exit{},
	}
	character.InRoom = room

	// Create a door
	room.Exits[types.DIR_WEST] = &types.Exit{
		Direction:   types.DIR_WEST,
		Description: "A large gate leads west.",
		Keywords:    "gate large",
		Flags:       types.EX_ISDOOR | types.EX_CLOSED,
		Key:         -1,
		DestVnum:    1002,
	}

	// Test various argument formats with fill words
	testCases := []string{
		"gate west",
		"the gate west",
		"a gate west",
		"gate to the west",
		"large west", // This should work since "large" is in the keywords
	}

	for i, args := range testCases {
		// Reset door to closed state
		room.Exits[types.DIR_WEST].Flags |= types.EX_CLOSED

		openCmd := &OpenCommand{}
		err := openCmd.Execute(character, args)
		if err != nil {
			t.Errorf("Test case %d: Expected no error for args '%s', got: %v", i, args, err)
		}

		// Check that the door is open
		if (room.Exits[types.DIR_WEST].Flags & types.EX_CLOSED) != 0 {
			t.Errorf("Test case %d: Expected door to be open for args '%s', but CLOSED flag is still set", i, args)
		}
	}
}
