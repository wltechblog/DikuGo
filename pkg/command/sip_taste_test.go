package command

import (
	"strings"
	"testing"

	"github.com/wltechblog/DikuGo/pkg/types"
)

func TestSipCommand_Execute(t *testing.T) {
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

	// Create a drink container with water
	waterBottle := &types.ObjectInstance{
		Prototype: &types.Object{
			VNUM:      1001,
			Name:      "bottle",
			ShortDesc: "a water bottle",
			Type:      types.ITEM_DRINKCON,
			Value:     [4]int{10, 5, 0, 0}, // capacity 10, current 5, water (type 0)
		},
		Value: [4]int{0, 5, 0, 0}, // Use instance values for liquid amount
	}
	character.Inventory = append(character.Inventory, waterBottle)

	// Test sipping from the bottle
	sipCmd := &SipCommand{}
	err := sipCmd.Execute(character, "bottle")
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	// Check that liquid amount decreased by 1
	if waterBottle.Value[1] != 4 {
		t.Errorf("Expected liquid amount to be 4, got: %d", waterBottle.Value[1])
	}

	// Check that thirst increased (water gives +10 thirst per unit, /4 = 2.5, rounded to 2)
	expectedThirst := 10 / 4 // 2
	if character.Conditions[types.COND_THIRST] != expectedThirst {
		t.Errorf("Expected thirst to be %d, got: %d", expectedThirst, character.Conditions[types.COND_THIRST])
	}
}

func TestSipCommand_TooDrunk(t *testing.T) {
	// Create a test character that's too drunk
	character := &types.Character{
		Name:       "TestPlayer",
		Level:      5,
		Inventory:  make([]*types.ObjectInstance, 0),
		Conditions: [3]int{15, 0, 0}, // Drunk condition > 10
	}

	// Create a test room
	room := &types.Room{
		VNUM:       1001,
		Name:       "Test Room",
		Characters: []*types.Character{character},
	}
	character.InRoom = room

	// Create a drink container
	bottle := &types.ObjectInstance{
		Prototype: &types.Object{
			VNUM:      1001,
			Name:      "bottle",
			ShortDesc: "a bottle",
			Type:      types.ITEM_DRINKCON,
			Value:     [4]int{10, 5, 0, 0},
		},
		Value: [4]int{0, 5, 0, 0},
	}
	character.Inventory = append(character.Inventory, bottle)

	// Test sipping when too drunk
	sipCmd := &SipCommand{}
	err := sipCmd.Execute(character, "bottle")
	if err != nil {
		t.Errorf("Expected no error (should handle drunk state internally), got: %v", err)
	}

	// Check that liquid amount didn't change (failed to sip)
	if bottle.Value[1] != 5 {
		t.Errorf("Expected liquid amount to remain 5, got: %d", bottle.Value[1])
	}
}

func TestSipCommand_EmptyContainer(t *testing.T) {
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

	// Create an empty drink container
	emptyBottle := &types.ObjectInstance{
		Prototype: &types.Object{
			VNUM:      1001,
			Name:      "bottle",
			ShortDesc: "an empty bottle",
			Type:      types.ITEM_DRINKCON,
			Value:     [4]int{10, 0, 0, 0}, // Empty
		},
		Value: [4]int{0, 0, 0, 0},
	}
	character.Inventory = append(character.Inventory, emptyBottle)

	// Test sipping from empty container
	sipCmd := &SipCommand{}
	err := sipCmd.Execute(character, "bottle")
	if err == nil || !strings.Contains(err.Error(), "empty") {
		t.Errorf("Expected 'empty' error, got: %v", err)
	}
}

func TestTasteCommand_Food(t *testing.T) {
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
	apple := &types.ObjectInstance{
		Prototype: &types.Object{
			VNUM:      1001,
			Name:      "apple",
			ShortDesc: "a red apple",
			Type:      types.ITEM_FOOD,
			Value:     [4]int{5, 0, 0, 0}, // 5 food value
		},
		Value: [4]int{5, 0, 0, 0}, // Use instance values
	}
	character.Inventory = append(character.Inventory, apple)

	// Test tasting the apple
	tasteCmd := &TasteCommand{}
	err := tasteCmd.Execute(character, "apple")
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	// Check that fullness increased by 1
	if character.Conditions[types.COND_FULL] != 1 {
		t.Errorf("Expected fullness to be 1, got: %d", character.Conditions[types.COND_FULL])
	}

	// Check that food value decreased by 1
	if apple.Value[0] != 4 {
		t.Errorf("Expected food value to be 4, got: %d", apple.Value[0])
	}

	// Check that apple is still in inventory
	if len(character.Inventory) != 1 {
		t.Errorf("Expected apple to remain in inventory, got: %d items", len(character.Inventory))
	}
}

func TestTasteCommand_DrinkContainer(t *testing.T) {
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

	// Create a drink container
	bottle := &types.ObjectInstance{
		Prototype: &types.Object{
			VNUM:      1001,
			Name:      "bottle",
			ShortDesc: "a bottle",
			Type:      types.ITEM_DRINKCON,
			Value:     [4]int{10, 5, 0, 0},
		},
		Value: [4]int{0, 5, 0, 0},
	}
	character.Inventory = append(character.Inventory, bottle)

	// Test tasting drink container (should redirect to sip)
	tasteCmd := &TasteCommand{}
	err := tasteCmd.Execute(character, "bottle")
	if err != nil {
		t.Errorf("Expected no error (should redirect to sip), got: %v", err)
	}

	// Check that liquid amount decreased (sip was called)
	if bottle.Value[1] != 4 {
		t.Errorf("Expected liquid amount to be 4 (sip was called), got: %d", bottle.Value[1])
	}
}

func TestTasteCommand_FoodConsumed(t *testing.T) {
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

	// Create a food item with only 1 value left
	crumb := &types.ObjectInstance{
		Prototype: &types.Object{
			VNUM:      1001,
			Name:      "crumb",
			ShortDesc: "a bread crumb",
			Type:      types.ITEM_FOOD,
			Value:     [4]int{1, 0, 0, 0}, // 1 food value
		},
		Value: [4]int{1, 0, 0, 0}, // Use instance values
	}
	character.Inventory = append(character.Inventory, crumb)

	// Test tasting the crumb (should consume it completely)
	tasteCmd := &TasteCommand{}
	err := tasteCmd.Execute(character, "crumb")
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	// Check that crumb was removed from inventory
	if len(character.Inventory) != 0 {
		t.Errorf("Expected crumb to be consumed, got: %d items", len(character.Inventory))
	}
}

func TestTasteCommand_PoisonedFood(t *testing.T) {
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
	poisonedBerry := &types.ObjectInstance{
		Prototype: &types.Object{
			VNUM:      1002,
			Name:      "berry",
			ShortDesc: "a strange berry",
			Type:      types.ITEM_FOOD,
			Value:     [4]int{2, 0, 0, 1}, // 2 food value, poisoned
		},
		Value: [4]int{2, 0, 0, 0}, // Use prototype poison value
	}
	character.Inventory = append(character.Inventory, poisonedBerry)

	// Test tasting poisoned food
	tasteCmd := &TasteCommand{}
	err := tasteCmd.Execute(character, "berry")
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	// Check that character is now poisoned
	if (character.AffectedBy & types.AFF_POISON) == 0 {
		t.Error("Expected character to be poisoned")
	}

	// Check that poison affect was added with duration 2
	if character.Affected == nil {
		t.Error("Expected poison affect to be added")
	} else {
		if character.Affected.Duration != 2 {
			t.Errorf("Expected poison duration 2, got: %d", character.Affected.Duration)
		}
	}
}

func TestTasteCommand_NotFood(t *testing.T) {
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

	// Create a non-food, non-drink item
	rock := &types.ObjectInstance{
		Prototype: &types.Object{
			VNUM:      1003,
			Name:      "rock",
			ShortDesc: "a small rock",
			Type:      types.ITEM_TREASURE,
		},
	}
	character.Inventory = append(character.Inventory, rock)

	// Test tasting non-food
	tasteCmd := &TasteCommand{}
	err := tasteCmd.Execute(character, "rock")
	if err == nil || !strings.Contains(err.Error(), "stomach refuses") {
		t.Errorf("Expected 'stomach refuses' error, got: %v", err)
	}
}
