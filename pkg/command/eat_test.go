package command

import (
	"strings"
	"testing"

	"github.com/wltechblog/DikuGo/pkg/types"
)

func TestEatCommand_Execute(t *testing.T) {
	// Create a test character
	character := &types.Character{
		Name:       "TestPlayer",
		Level:      5,
		Inventory:  make([]*types.ObjectInstance, 0),
		Conditions: [3]int{0, 0, 0},
	}

	// Create a test room
	room := &types.Room{
		VNUM:       1001,
		Name:       "Test Room",
		Characters: []*types.Character{character},
	}
	character.InRoom = room

	// Create a food item
	bread := &types.ObjectInstance{
		Prototype: &types.Object{
			VNUM:      1001,
			Name:      "bread",
			ShortDesc: "a loaf of bread",
			Type:      types.ITEM_FOOD,
			Value:     [4]int{5, 0, 0, 0}, // 5 food value
		},
		Value: [4]int{0, 0, 0, 0}, // Use prototype values
	}
	character.Inventory = append(character.Inventory, bread)

	// Test eating the bread
	eatCmd := &EatCommand{}
	err := eatCmd.Execute(character, "bread")
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	// Check that fullness increased
	if character.Conditions[types.COND_FULL] != 5 {
		t.Errorf("Expected fullness to be 5, got: %d", character.Conditions[types.COND_FULL])
	}

	// Check that bread was removed from inventory
	if len(character.Inventory) != 0 {
		t.Errorf("Expected inventory to be empty, got: %d items", len(character.Inventory))
	}
}

func TestEatCommand_PoisonedFood(t *testing.T) {
	// Create a test character
	character := &types.Character{
		Name:       "TestPlayer",
		Level:      5,
		Inventory:  make([]*types.ObjectInstance, 0),
		Conditions: [3]int{0, 0, 0},
		AffectedBy: 0,
	}

	// Create a test room
	room := &types.Room{
		VNUM:       1001,
		Name:       "Test Room",
		Characters: []*types.Character{character},
	}
	character.InRoom = room

	// Create a poisoned food item
	poisonedMeat := &types.ObjectInstance{
		Prototype: &types.Object{
			VNUM:      1002,
			Name:      "meat",
			ShortDesc: "a piece of meat",
			Type:      types.ITEM_FOOD,
			Value:     [4]int{3, 0, 0, 1}, // 3 food value, poisoned (Value[3] = 1)
		},
		Value: [4]int{0, 0, 0, 0}, // Use prototype values
	}
	character.Inventory = append(character.Inventory, poisonedMeat)

	// Test eating the poisoned meat
	eatCmd := &EatCommand{}
	err := eatCmd.Execute(character, "meat")
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	// Check that character is now poisoned
	if (character.AffectedBy & types.AFF_POISON) == 0 {
		t.Error("Expected character to be poisoned")
	}

	// Check that poison affect was added
	if character.Affected == nil {
		t.Error("Expected poison affect to be added")
	} else {
		if character.Affected.Type != 33 { // SPELL_POISON
			t.Errorf("Expected poison spell type 33, got: %d", character.Affected.Type)
		}
		if character.Affected.Duration != 6 { // 3 * 2
			t.Errorf("Expected poison duration 6, got: %d", character.Affected.Duration)
		}
	}
}

func TestEatCommand_TooFull(t *testing.T) {
	// Create a test character that's already full
	character := &types.Character{
		Name:       "TestPlayer",
		Level:      5,
		Inventory:  make([]*types.ObjectInstance, 0),
		Conditions: [3]int{0, 21, 0}, // Full condition > 20
	}

	// Create a test room
	room := &types.Room{
		VNUM:       1001,
		Name:       "Test Room",
		Characters: []*types.Character{character},
	}
	character.InRoom = room

	// Create a food item
	bread := &types.ObjectInstance{
		Prototype: &types.Object{
			VNUM:      1001,
			Name:      "bread",
			ShortDesc: "a loaf of bread",
			Type:      types.ITEM_FOOD,
			Value:     [4]int{5, 0, 0, 0},
		},
	}
	character.Inventory = append(character.Inventory, bread)

	// Test eating when too full
	eatCmd := &EatCommand{}
	err := eatCmd.Execute(character, "bread")
	if err == nil || !strings.Contains(err.Error(), "too full") {
		t.Errorf("Expected 'too full' error, got: %v", err)
	}

	// Check that bread is still in inventory
	if len(character.Inventory) != 1 {
		t.Errorf("Expected bread to remain in inventory, got: %d items", len(character.Inventory))
	}
}

func TestEatCommand_NotFood(t *testing.T) {
	// Create a test character
	character := &types.Character{
		Name:      "TestPlayer",
		Level:     5,
		Inventory: make([]*types.ObjectInstance, 0),
	}

	// Create a test room
	room := &types.Room{
		VNUM:       1001,
		Name:       "Test Room",
		Characters: []*types.Character{character},
	}
	character.InRoom = room

	// Create a non-food item
	sword := &types.ObjectInstance{
		Prototype: &types.Object{
			VNUM:      1003,
			Name:      "sword",
			ShortDesc: "a steel sword",
			Type:      types.ITEM_WEAPON,
		},
	}
	character.Inventory = append(character.Inventory, sword)

	// Test eating non-food
	eatCmd := &EatCommand{}
	err := eatCmd.Execute(character, "sword")
	if err == nil || !strings.Contains(err.Error(), "stomach refuses") {
		t.Errorf("Expected 'stomach refuses' error, got: %v", err)
	}
}

func TestEatCommand_AdminCanEatAnything(t *testing.T) {
	// Create an admin character (level 22+)
	character := &types.Character{
		Name:       "AdminPlayer",
		Level:      25,
		Inventory:  make([]*types.ObjectInstance, 0),
		Conditions: [3]int{0, 0, 0},
	}

	// Create a test room
	room := &types.Room{
		VNUM:       1001,
		Name:       "Test Room",
		Characters: []*types.Character{character},
	}
	character.InRoom = room

	// Create a non-food item
	sword := &types.ObjectInstance{
		Prototype: &types.Object{
			VNUM:      1003,
			Name:      "sword",
			ShortDesc: "a steel sword",
			Type:      types.ITEM_WEAPON,
			Value:     [4]int{1, 0, 0, 0}, // 1 food value for admin
		},
	}
	character.Inventory = append(character.Inventory, sword)

	// Test admin eating non-food (should work)
	eatCmd := &EatCommand{}
	err := eatCmd.Execute(character, "sword")
	if err != nil {
		t.Errorf("Expected admin to be able to eat anything, got: %v", err)
	}

	// Check that sword was consumed
	if len(character.Inventory) != 0 {
		t.Errorf("Expected sword to be consumed, got: %d items", len(character.Inventory))
	}
}

func TestEatCommand_NotFound(t *testing.T) {
	// Create a test character
	character := &types.Character{
		Name:      "TestPlayer",
		Level:     5,
		Inventory: make([]*types.ObjectInstance, 0),
	}

	// Create a test room
	room := &types.Room{
		VNUM:       1001,
		Name:       "Test Room",
		Characters: []*types.Character{character},
	}
	character.InRoom = room

	// Test eating non-existent item
	eatCmd := &EatCommand{}
	err := eatCmd.Execute(character, "pizza")
	if err == nil || !strings.Contains(err.Error(), "can't find it") {
		t.Errorf("Expected 'can't find it' error, got: %v", err)
	}
}

func TestEatCommand_NoArgs(t *testing.T) {
	// Create a test character
	character := &types.Character{
		Name: "TestPlayer",
	}

	// Test eating with no arguments
	eatCmd := &EatCommand{}
	err := eatCmd.Execute(character, "")
	if err == nil || !strings.Contains(err.Error(), "eat what") {
		t.Errorf("Expected 'eat what' error, got: %v", err)
	}
}
