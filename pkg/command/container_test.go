package command

import (
	"strings"
	"testing"

	"github.com/wltechblog/DikuGo/pkg/types"
)

func TestGetFromContainer(t *testing.T) {
	// Create a test character
	character := &types.Character{
		Name:      "TestPlayer",
		Level:     1,
		Position:  types.POS_STANDING,
		Inventory: make([]*types.ObjectInstance, 0),
	}

	// Create a test room
	room := &types.Room{
		VNUM:        1001,
		Name:        "Test Room",
		Description: "A test room",
		Objects:     make([]*types.ObjectInstance, 0),
		Characters:  []*types.Character{character},
	}
	character.InRoom = room

	// Create a container (chest)
	chest := &types.ObjectInstance{
		Prototype: &types.Object{
			VNUM:      2001,
			Name:      "chest",
			ShortDesc: "a wooden chest",
			Type:      types.ITEM_CONTAINER,
			Value:     [4]int{100, 0, 0, 0}, // capacity, flags, key, corpse
		},
		Contains: make([]*types.ObjectInstance, 0),
	}
	room.Objects = append(room.Objects, chest)

	// Create an item to put in the container
	sword := &types.ObjectInstance{
		Prototype: &types.Object{
			VNUM:      2002,
			Name:      "sword",
			ShortDesc: "a steel sword",
			Type:      types.ITEM_WEAPON,
		},
		InObj: chest,
	}
	chest.Contains = append(chest.Contains, sword)

	// Test getting item from container
	getCmd := &GetCommand{}
	err := getCmd.Execute(character, "sword from chest")
	// In DikuMUD pattern, success messages are returned as errors
	if err == nil || !strings.Contains(err.Error(), "you get") {
		t.Errorf("Expected success message containing 'you get', got: %v", err)
	}

	// Check that the sword is now in the character's inventory
	if len(character.Inventory) != 1 {
		t.Errorf("Expected 1 item in inventory, got: %d", len(character.Inventory))
	}

	if character.Inventory[0] != sword {
		t.Errorf("Expected sword in inventory")
	}

	// Check that the sword is no longer in the chest
	if len(chest.Contains) != 0 {
		t.Errorf("Expected chest to be empty, got: %d items", len(chest.Contains))
	}

	// Test getting from non-existent container
	err = getCmd.Execute(character, "item from nonexistent")
	if err == nil || err.Error() != "you don't see nonexistent here" {
		t.Errorf("Expected 'you don't see nonexistent here' error, got: %v", err)
	}

	// Test getting non-existent item from container
	err = getCmd.Execute(character, "nonexistent from chest")
	if err == nil || err.Error() != "you don't see nonexistent in a wooden chest" {
		t.Errorf("Expected 'you don't see nonexistent in a wooden chest' error, got: %v", err)
	}
}

func TestGetFromCorpse(t *testing.T) {
	// Create a test character
	character := &types.Character{
		Name:      "TestPlayer",
		Level:     1,
		Position:  types.POS_STANDING,
		Inventory: make([]*types.ObjectInstance, 0),
	}

	// Create a test room
	room := &types.Room{
		VNUM:        1001,
		Name:        "Test Room",
		Description: "A test room",
		Objects:     make([]*types.ObjectInstance, 0),
		Characters:  []*types.Character{character},
	}
	character.InRoom = room

	// Create a corpse (container)
	corpse := &types.ObjectInstance{
		Prototype: &types.Object{
			VNUM:      -1, // Special VNUM for corpses
			Name:      "corpse",
			ShortDesc: "corpse of a goblin",
			Type:      types.ITEM_CONTAINER,
			Value:     [4]int{0, 0, 0, 1}, // Value[3] = 1 identifies it as a corpse
		},
		Contains: make([]*types.ObjectInstance, 0),
		Timer:    15, // Corpse decay timer
	}
	room.Objects = append(room.Objects, corpse)

	// Create items that would be in the corpse
	gold := &types.ObjectInstance{
		Prototype: &types.Object{
			VNUM:      1,
			Name:      "gold",
			ShortDesc: "some gold coins",
			Type:      types.ITEM_MONEY,
			Value:     [4]int{100, 0, 0, 0}, // 100 gold coins
		},
		InObj: corpse,
	}
	corpse.Contains = append(corpse.Contains, gold)

	dagger := &types.ObjectInstance{
		Prototype: &types.Object{
			VNUM:      2003,
			Name:      "dagger",
			ShortDesc: "a rusty dagger",
			Type:      types.ITEM_WEAPON,
		},
		InObj: corpse,
	}
	corpse.Contains = append(corpse.Contains, dagger)

	// Test getting gold from corpse
	initialGold := character.Gold
	getCmd := &GetCommand{}
	err := getCmd.Execute(character, "gold from corpse")
	if err == nil || !strings.Contains(err.Error(), "There were") {
		t.Errorf("Expected success message containing 'There were', got: %v", err)
	}

	// Check that gold was added to character's gold, not inventory
	expectedGold := initialGold + 100
	if character.Gold != expectedGold {
		t.Errorf("Expected gold to be %d, got: %d", expectedGold, character.Gold)
	}

	// Check that gold object is NOT in inventory (should be converted to currency)
	if len(character.Inventory) != 0 {
		t.Errorf("Expected 0 items in inventory after getting gold, got: %d", len(character.Inventory))
	}

	// Test getting dagger from corpse
	err = getCmd.Execute(character, "dagger from corpse")
	if err == nil || !strings.Contains(err.Error(), "you get") {
		t.Errorf("Expected success message containing 'you get', got: %v", err)
	}

	// Check that both items are now in inventory
	if len(character.Inventory) != 1 {
		t.Errorf("Expected 1 item in inventory (dagger only), got: %d", len(character.Inventory))
	}

	// Check that corpse is now empty
	if len(corpse.Contains) != 0 {
		t.Errorf("Expected corpse to be empty, got: %d items", len(corpse.Contains))
	}
}

func TestGetFromClosedContainer(t *testing.T) {
	// Create a test character
	character := &types.Character{
		Name:      "TestPlayer",
		Level:     1,
		Position:  types.POS_STANDING,
		Inventory: make([]*types.ObjectInstance, 0),
	}

	// Create a test room
	room := &types.Room{
		VNUM:        1001,
		Name:        "Test Room",
		Description: "A test room",
		Objects:     make([]*types.ObjectInstance, 0),
		Characters:  []*types.Character{character},
	}
	character.InRoom = room

	// Create a closed container
	chest := &types.ObjectInstance{
		Prototype: &types.Object{
			VNUM:      2001,
			Name:      "chest",
			ShortDesc: "a wooden chest",
			Type:      types.ITEM_CONTAINER,
			Value:     [4]int{100, types.CONT_CLOSED, 0, 0}, // Closed container
		},
		Contains: make([]*types.ObjectInstance, 0),
	}
	room.Objects = append(room.Objects, chest)

	// Create an item in the closed container
	sword := &types.ObjectInstance{
		Prototype: &types.Object{
			VNUM:      2002,
			Name:      "sword",
			ShortDesc: "a steel sword",
			Type:      types.ITEM_WEAPON,
		},
		InObj: chest,
	}
	chest.Contains = append(chest.Contains, sword)

	// Test getting from closed container - should fail
	getCmd := &GetCommand{}
	err := getCmd.Execute(character, "sword from chest")
	if err == nil || err.Error() != "a wooden chest is closed" {
		t.Errorf("Expected 'a wooden chest is closed' error, got: %v", err)
	}

	// Check that sword is still in chest
	if len(chest.Contains) != 1 {
		t.Errorf("Expected 1 item in chest, got: %d", len(chest.Contains))
	}

	// Check that character inventory is empty
	if len(character.Inventory) != 0 {
		t.Errorf("Expected 0 items in inventory, got: %d", len(character.Inventory))
	}
}

func TestGetAllFromContainerWithGold(t *testing.T) {
	// Create a test character
	character := &types.Character{
		Name:      "TestPlayer",
		Gold:      50,
		Inventory: make([]*types.ObjectInstance, 0),
	}

	// Create a test room
	room := &types.Room{
		VNUM:    1001,
		Name:    "Test Room",
		Objects: make([]*types.ObjectInstance, 0),
	}
	character.InRoom = room

	// Create a corpse container
	corpse := &types.ObjectInstance{
		Prototype: &types.Object{
			VNUM:      2001,
			Name:      "corpse",
			ShortDesc: "a corpse",
			Type:      types.ITEM_CONTAINER,
		},
		Contains: make([]*types.ObjectInstance, 0),
	}
	room.Objects = append(room.Objects, corpse)

	// Add gold coins to the corpse
	gold1 := &types.ObjectInstance{
		Prototype: &types.Object{
			VNUM:      -2,
			Name:      "coins gold",
			ShortDesc: "50 gold coins",
			Type:      types.ITEM_MONEY,
			Value:     [4]int{50, 0, 0, 0},
		},
		InObj: corpse,
	}
	corpse.Contains = append(corpse.Contains, gold1)

	// Add more gold coins to the corpse
	gold2 := &types.ObjectInstance{
		Prototype: &types.Object{
			VNUM:      -2,
			Name:      "coins gold",
			ShortDesc: "25 gold coins",
			Type:      types.ITEM_MONEY,
			Value:     [4]int{25, 0, 0, 0},
		},
		InObj: corpse,
	}
	corpse.Contains = append(corpse.Contains, gold2)

	// Add a sword to the corpse
	sword := &types.ObjectInstance{
		Prototype: &types.Object{
			VNUM:      2003,
			Name:      "sword",
			ShortDesc: "a steel sword",
			Type:      types.ITEM_WEAPON,
		},
		InObj: corpse,
	}
	corpse.Contains = append(corpse.Contains, sword)

	// Add a potion to the corpse
	potion := &types.ObjectInstance{
		Prototype: &types.Object{
			VNUM:      2004,
			Name:      "potion",
			ShortDesc: "a healing potion",
			Type:      types.ITEM_POTION,
		},
		InObj: corpse,
	}
	corpse.Contains = append(corpse.Contains, potion)

	// Test getting all items from corpse
	initialGold := character.Gold
	getCmd := &GetCommand{}
	err := getCmd.Execute(character, "all from corpse")
	if err == nil {
		t.Errorf("Expected success message, got: %v", err)
	}

	// Check that gold was added to character's gold (50 + 25 = 75)
	expectedGold := initialGold + 75
	if character.Gold != expectedGold {
		t.Errorf("Expected gold to be %d, got: %d", expectedGold, character.Gold)
	}

	// Check that only regular items are in inventory (sword and potion)
	if len(character.Inventory) != 2 {
		t.Errorf("Expected 2 items in inventory (sword and potion), got: %d", len(character.Inventory))
	}

	// Verify the message contains both gold and item information
	message := err.Error()
	if !strings.Contains(message, "There were 75 coins") {
		t.Errorf("Expected message to contain 'There were 75 coins', got: %s", message)
	}
	if !strings.Contains(message, "You get a steel sword") {
		t.Errorf("Expected message to contain 'You get a steel sword', got: %s", message)
	}
	if !strings.Contains(message, "You get a healing potion") {
		t.Errorf("Expected message to contain 'You get a healing potion', got: %s", message)
	}
}

func TestGetSingleGoldCoin(t *testing.T) {
	// Create a test character
	character := &types.Character{
		Name:      "TestPlayer",
		Gold:      10,
		Inventory: make([]*types.ObjectInstance, 0),
	}

	// Create a test room
	room := &types.Room{
		VNUM:    1001,
		Name:    "Test Room",
		Objects: make([]*types.ObjectInstance, 0),
	}
	character.InRoom = room

	// Create a single gold coin on the ground
	goldCoin := &types.ObjectInstance{
		Prototype: &types.Object{
			VNUM:      -2,
			Name:      "coin gold",
			ShortDesc: "a gold coin",
			Type:      types.ITEM_MONEY,
			Value:     [4]int{1, 0, 0, 0}, // 1 gold coin
		},
	}
	room.Objects = append(room.Objects, goldCoin)

	// Test getting the single gold coin
	initialGold := character.Gold
	getCmd := &GetCommand{}
	err := getCmd.Execute(character, "coin")
	if err == nil || !strings.Contains(err.Error(), "There was 1 coin") {
		t.Errorf("Expected success message containing 'There was 1 coin', got: %v", err)
	}

	// Check that gold was added to character's gold
	expectedGold := initialGold + 1
	if character.Gold != expectedGold {
		t.Errorf("Expected gold to be %d, got: %d", expectedGold, character.Gold)
	}

	// Check that gold coin is NOT in inventory
	if len(character.Inventory) != 0 {
		t.Errorf("Expected 0 items in inventory after getting gold coin, got: %d", len(character.Inventory))
	}

	// Check that gold coin is no longer in room
	if len(room.Objects) != 0 {
		t.Errorf("Expected 0 objects in room after getting gold coin, got: %d", len(room.Objects))
	}
}
