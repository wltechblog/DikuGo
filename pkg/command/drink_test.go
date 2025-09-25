package command

import (
	"testing"

	"github.com/wltechblog/DikuGo/pkg/types"
)

func TestDrinkCommand_Execute(t *testing.T) {
	// Create a test character
	character := &types.Character{
		Name:       "TestPlayer",
		Level:      1,
		Position:   types.POS_STANDING,
		Conditions: [3]int{0, 0, 0},
		Inventory:  make([]*types.ObjectInstance, 0),
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

	// Create a drink container (water fountain)
	fountain := &types.ObjectInstance{
		Prototype: &types.Object{
			VNUM:      2001,
			Name:      "fountain",
			ShortDesc: "a water fountain",
			Type:      types.ITEM_FOUNTAIN,
			Value:     [4]int{10, 10, 0, 0}, // capacity, current, liquid type, poisoned
		},
		Value: [4]int{10, 10, 0, 0}, // Instance values
	}
	room.Objects = append(room.Objects, fountain)

	// Create a drink container (water bottle)
	bottle := &types.ObjectInstance{
		Prototype: &types.Object{
			VNUM:      2002,
			Name:      "bottle",
			ShortDesc: "a water bottle",
			Type:      types.ITEM_DRINKCON,
			Value:     [4]int{5, 3, 0, 0}, // capacity, current, liquid type, poisoned
		},
		Value: [4]int{5, 3, 0, 0}, // Instance values
	}
	character.Inventory = append(character.Inventory, bottle)

	// Test drinking from fountain
	drinkCmd := &DrinkCommand{}
	err := drinkCmd.Execute(character, "fountain")
	if err != nil {
		t.Errorf("Expected no error drinking from fountain, got: %v", err)
	}

	// Test drinking from bottle in inventory
	t.Logf("Before drinking: bottle.Value[1] = %d", bottle.Value[1])
	err = drinkCmd.Execute(character, "bottle")
	if err != nil {
		t.Errorf("Expected no error drinking from bottle, got: %v", err)
	}
	t.Logf("After drinking: bottle.Value[1] = %d", bottle.Value[1])

	// Check that bottle liquid decreased (fountain should stay the same)
	if bottle.Value[1] != 2 { // Should have decreased from 3 to 2
		t.Errorf("Expected bottle liquid to decrease to 2, got: %d", bottle.Value[1])
	}

	if fountain.Value[1] != 10 { // Fountain should stay the same
		t.Errorf("Expected fountain liquid to stay at 10, got: %d", fountain.Value[1])
	}

	// Test drinking from empty container
	bottle.Value[1] = 0
	err = drinkCmd.Execute(character, "bottle")
	if err == nil || err.Error() != "it is empty." {
		t.Errorf("Expected 'it is empty.' error, got: %v", err)
	}

	// Test drinking from non-existent object
	err = drinkCmd.Execute(character, "nonexistent")
	if err == nil || err.Error() != "you can't find it!" {
		t.Errorf("Expected 'you can't find it!' error, got: %v", err)
	}

	// Test drinking from non-drink object
	sword := &types.ObjectInstance{
		Prototype: &types.Object{
			VNUM:      2003,
			Name:      "sword",
			ShortDesc: "a sword",
			Type:      types.ITEM_WEAPON,
		},
	}
	room.Objects = append(room.Objects, sword)

	err = drinkCmd.Execute(character, "sword")
	if err == nil || err.Error() != "you can't drink from that!" {
		t.Errorf("Expected 'you can't drink from that!' error, got: %v", err)
	}
}

func TestDrinkCommand_Properties(t *testing.T) {
	cmd := &DrinkCommand{}

	if cmd.Name() != "drink" {
		t.Errorf("Expected name 'drink', got: %s", cmd.Name())
	}

	if len(cmd.Aliases()) != 0 {
		t.Errorf("Expected no aliases, got: %v", cmd.Aliases())
	}

	if cmd.MinPosition() != types.POS_RESTING {
		t.Errorf("Expected min position %d, got: %d", types.POS_RESTING, cmd.MinPosition())
	}

	if cmd.Level() != 1 {
		t.Errorf("Expected level 1, got: %d", cmd.Level())
	}

	if cmd.LogCommand() != false {
		t.Errorf("Expected LogCommand false, got: %t", cmd.LogCommand())
	}
}

func TestDrinkEffects(t *testing.T) {
	// Test drink effect calculations
	tests := []struct {
		liquidType int
		effectType int
		expected   int
	}{
		{0, DRINK_EFFECT_DRUNK, 0},   // water - no drunk effect
		{0, DRINK_EFFECT_THIRST, 10}, // water - high thirst relief
		{1, DRINK_EFFECT_DRUNK, 3},   // beer - moderate drunk effect
		{5, DRINK_EFFECT_DRUNK, 6},   // whisky - high drunk effect
	}

	for _, test := range tests {
		result := getDrinkEffect(test.liquidType, test.effectType)
		if result != test.expected {
			t.Errorf("getDrinkEffect(%d, %d) = %d, expected %d",
				test.liquidType, test.effectType, result, test.expected)
		}
	}
}

func TestDrinkNames(t *testing.T) {
	// Test drink name lookup
	tests := []struct {
		liquidType int
		expected   string
	}{
		{0, "water"},
		{1, "beer"},
		{2, "wine"},
		{15, "coca cola"},
		{999, "water"}, // Invalid type should default to water
	}

	for _, test := range tests {
		result := getDrinkName(test.liquidType)
		if result != test.expected {
			t.Errorf("getDrinkName(%d) = %s, expected %s",
				test.liquidType, result, test.expected)
		}
	}
}
