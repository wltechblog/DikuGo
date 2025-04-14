package command

import (
	"strings"
	"testing"

	"github.com/wltechblog/DikuGo/pkg/types"
)

func TestEquipmentCommand(t *testing.T) {
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

	// Create an equipment command
	equipCmd := &EquipmentCommand{}

	// Test with no equipment
	err := equipCmd.Execute(character, "")
	if err == nil {
		t.Errorf("Expected an error for no equipment, got nil")
	} else {
		// The error should indicate that the character has no equipment
		if !strings.Contains(err.Error(), "you are not using any equipment") {
			t.Errorf("Expected error to indicate no equipment, got: %s", err.Error())
		}
	}

	// Add objects to character's equipment
	character.Equipment[types.WEAR_WIELD] = swordInstance
	character.Equipment[types.WEAR_SHIELD] = shieldInstance

	// Test with equipment
	err = equipCmd.Execute(character, "")
	if err == nil {
		t.Errorf("Expected an error (which contains the equipment list), got nil")
	} else {
		// The error should contain the equipment list
		if !strings.Contains(err.Error(), "You are using:") {
			t.Errorf("Expected error to contain equipment list header, got: %s", err.Error())
		}
		if !strings.Contains(err.Error(), "wielded:") && !strings.Contains(err.Error(), "a long sword") {
			t.Errorf("Expected error to list sword, got: %s", err.Error())
		}
		if !strings.Contains(err.Error(), "shield:") && !strings.Contains(err.Error(), "a wooden shield") {
			t.Errorf("Expected error to list shield, got: %s", err.Error())
		}
	}
}
