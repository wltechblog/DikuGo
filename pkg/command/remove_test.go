package command

import (
	"strings"
	"testing"

	"github.com/wltechblog/DikuGo/pkg/types"
)

func TestRemoveCommand(t *testing.T) {
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
		WornBy:    character,
		WornOn:    types.WEAR_WIELD,
	}

	shieldInstance := &types.ObjectInstance{
		Prototype: shield,
		WornBy:    character,
		WornOn:    types.WEAR_SHIELD,
	}

	// Add objects to character's equipment
	character.Equipment[types.WEAR_WIELD] = swordInstance
	character.Equipment[types.WEAR_SHIELD] = shieldInstance

	// Create a remove command
	removeCmd := &RemoveCommand{}

	// Test removing the sword
	err := removeCmd.Execute(character, "sword")
	if err == nil {
		t.Errorf("Expected an error (which contains the success message), got nil")
	} else {
		// The error should contain a success message
		if !strings.Contains(err.Error(), "you remove a long sword") {
			t.Errorf("Expected error to contain success message, got: %s", err.Error())
		}
	}

	// Check if the sword is now in the character's inventory
	found := false
	for _, obj := range character.Inventory {
		if obj == swordInstance {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Expected sword to be added to character's inventory")
	}

	// Check if the sword is no longer in the character's equipment
	if character.Equipment[types.WEAR_WIELD] != nil {
		t.Errorf("Expected sword to be removed from character's equipment")
	}

	// Test removing something that isn't worn
	err = removeCmd.Execute(character, "sword")
	if err == nil {
		t.Errorf("Expected an error for item not worn, got nil")
	} else {
		// The error should indicate that the item isn't worn
		if !strings.Contains(err.Error(), "you're not wearing sword") {
			t.Errorf("Expected error to indicate item isn't worn, got: %s", err.Error())
		}
	}
}

func TestRemoveAll(t *testing.T) {
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
		WornBy:    character,
		WornOn:    types.WEAR_WIELD,
	}

	shieldInstance := &types.ObjectInstance{
		Prototype: shield,
		WornBy:    character,
		WornOn:    types.WEAR_SHIELD,
	}

	helmetInstance := &types.ObjectInstance{
		Prototype: helmet,
		WornBy:    character,
		WornOn:    types.WEAR_HEAD,
	}

	// Add objects to character's equipment
	character.Equipment[types.WEAR_WIELD] = swordInstance
	character.Equipment[types.WEAR_SHIELD] = shieldInstance
	character.Equipment[types.WEAR_HEAD] = helmetInstance

	// Create a remove command
	removeCmd := &RemoveCommand{}

	// Test removing all
	err := removeCmd.Execute(character, "all")
	if err == nil {
		t.Errorf("Expected an error (which contains the success messages), got nil")
	} else {
		// The error should contain success messages for all items
		if !strings.Contains(err.Error(), "You remove a long sword") {
			t.Errorf("Expected error to contain success message for sword, got: %s", err.Error())
		}
		if !strings.Contains(err.Error(), "You remove a wooden shield") {
			t.Errorf("Expected error to contain success message for shield, got: %s", err.Error())
		}
		if !strings.Contains(err.Error(), "You remove a steel helmet") {
			t.Errorf("Expected error to contain success message for helmet, got: %s", err.Error())
		}
	}

	// Check if all items are now in the character's inventory
	if len(character.Inventory) != 3 {
		t.Errorf("Expected character's inventory to have 3 items, got %d", len(character.Inventory))
	}

	// Check if all equipment slots are now empty
	for i, obj := range character.Equipment {
		if obj != nil {
			t.Errorf("Expected equipment slot %d to be empty, got %v", i, obj)
		}
	}
}
