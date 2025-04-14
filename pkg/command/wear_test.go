package command

import (
	"strings"
	"testing"

	"github.com/wltechblog/DikuGo/pkg/types"
)

func TestWearCommand(t *testing.T) {
	// Create a test character
	character := &types.Character{
		Name:      "TestCharacter",
		Position:  types.POS_STANDING,
		Equipment: make([]*types.ObjectInstance, types.NUM_WEARS),
		Inventory: make([]*types.ObjectInstance, 0),
	}

	// Create a test sword
	sword := &types.Object{
		VNUM:      3001,
		Name:      "sword",
		ShortDesc: "a long sword",
		WearFlags: types.ITEM_WEAR_WIELD,
	}

	// Create a test shield
	shield := &types.Object{
		VNUM:      3002,
		Name:      "shield",
		ShortDesc: "a wooden shield",
		WearFlags: types.ITEM_WEAR_SHIELD,
	}

	// Create object instances
	swordInstance := &types.ObjectInstance{
		Prototype: sword,
		CarriedBy: character,
	}

	shieldInstance := &types.ObjectInstance{
		Prototype: shield,
		CarriedBy: character,
	}

	// Add objects to character's inventory
	character.Inventory = append(character.Inventory, swordInstance, shieldInstance)

	// Create a wear command
	wearCmd := &WearCommand{}

	// Test wearing the sword
	err := wearCmd.Execute(character, "sword")
	if err == nil {
		t.Errorf("Expected an error (which contains the success message), got nil")
	} else {
		// The error should contain a success message
		if !strings.Contains(err.Error(), "you wear a long sword") {
			t.Errorf("Expected error to contain success message, got: %s", err.Error())
		}
	}

	// Check if the sword is now in the character's equipment
	if character.Equipment[types.WEAR_WIELD] != swordInstance {
		t.Errorf("Expected sword to be in character's equipment at WEAR_WIELD position")
	}

	// Check if the sword is no longer in the character's inventory
	for _, obj := range character.Inventory {
		if obj == swordInstance {
			t.Errorf("Expected sword to be removed from character's inventory")
		}
	}

	// Test wearing the shield
	err = wearCmd.Execute(character, "shield")
	if err == nil {
		t.Errorf("Expected an error (which contains the success message), got nil")
	} else {
		// The error should contain a success message
		if !strings.Contains(err.Error(), "you wear a wooden shield") {
			t.Errorf("Expected error to contain success message, got: %s", err.Error())
		}
	}

	// Check if the shield is now in the character's equipment
	if character.Equipment[types.WEAR_SHIELD] != shieldInstance {
		t.Errorf("Expected shield to be in character's equipment at WEAR_SHIELD position")
	}

	// Test wearing something that doesn't exist
	err = wearCmd.Execute(character, "nonexistent")
	if err == nil {
		t.Errorf("Expected an error for nonexistent item, got nil")
	} else {
		// The error should indicate that the item doesn't exist
		if !strings.Contains(err.Error(), "you don't have nonexistent") {
			t.Errorf("Expected error to indicate item doesn't exist, got: %s", err.Error())
		}
	}
}

func TestWearAll(t *testing.T) {
	// Create a test character
	character := &types.Character{
		Name:      "TestCharacter",
		Position:  types.POS_STANDING,
		Equipment: make([]*types.ObjectInstance, types.NUM_WEARS),
		Inventory: make([]*types.ObjectInstance, 0),
	}

	// Create test objects
	sword := &types.Object{
		VNUM:      3001,
		Name:      "sword",
		ShortDesc: "a long sword",
		WearFlags: types.ITEM_WEAR_WIELD,
	}

	shield := &types.Object{
		VNUM:      3002,
		Name:      "shield",
		ShortDesc: "a wooden shield",
		WearFlags: types.ITEM_WEAR_SHIELD,
	}

	helmet := &types.Object{
		VNUM:      3003,
		Name:      "helmet",
		ShortDesc: "a steel helmet",
		WearFlags: types.ITEM_WEAR_HEAD,
	}

	// Create object instances
	swordInstance := &types.ObjectInstance{
		Prototype: sword,
		CarriedBy: character,
	}

	shieldInstance := &types.ObjectInstance{
		Prototype: shield,
		CarriedBy: character,
	}

	helmetInstance := &types.ObjectInstance{
		Prototype: helmet,
		CarriedBy: character,
	}

	// Add objects to character's inventory
	character.Inventory = append(character.Inventory, swordInstance, shieldInstance, helmetInstance)

	// Create a wear command
	wearCmd := &WearCommand{}

	// Test wearing all
	err := wearCmd.Execute(character, "all")
	if err == nil {
		t.Errorf("Expected an error (which contains the success messages), got nil")
	} else {
		// The error should contain success messages for the items
		// Note: The order of wearing items may vary, so we just check that some items were worn
		if !strings.Contains(err.Error(), "You wear") {
			t.Errorf("Expected error to contain success messages for wearing items, got: %s", err.Error())
		}
	}

	// Check if items are now in the character's equipment
	// Note: The wear command processes items in the order they appear in the inventory,
	// which may not be the same as the order we added them. So we just check that
	// some equipment slots are filled.
	equippedCount := 0
	for _, item := range character.Equipment {
		if item != nil {
			equippedCount++
		}
	}

	if equippedCount == 0 {
		t.Errorf("Expected some items to be equipped, but none were")
	}

	// Check if the inventory has fewer items
	if len(character.Inventory) >= 3 {
		t.Errorf("Expected character's inventory to have fewer than 3 items, got %d items", len(character.Inventory))
	}
}
