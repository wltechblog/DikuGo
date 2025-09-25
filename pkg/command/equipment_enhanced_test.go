package command

import (
	"strings"
	"testing"

	"github.com/wltechblog/DikuGo/pkg/types"
)

func TestHoldCommand_Light(t *testing.T) {
	// Create a test character
	character := &types.Character{
		Name:      "TestPlayer",
		Level:     5,
		Position:  types.POS_STANDING,
		Equipment: make([]*types.ObjectInstance, types.NUM_WEARS),
		Inventory: []*types.ObjectInstance{},
	}

	// Create a light object
	torch := &types.Object{
		VNUM:      1001,
		Name:      "torch",
		ShortDesc: "a torch",
		Type:      types.ITEM_LIGHT,
		WearFlags: types.ITEM_WEAR_TAKE,
		Value:     [4]int{0, 0, 10, 0}, // 10 hours of light
	}

	torchInstance := &types.ObjectInstance{
		Prototype: torch,
		CarriedBy: character,
	}

	// Add torch to character's inventory
	character.Inventory = append(character.Inventory, torchInstance)

	// Create a hold command
	holdCmd := &HoldCommand{}

	// Test holding the torch
	err := holdCmd.Execute(character, "torch")
	if err == nil {
		t.Errorf("Expected an error (which contains the success message), got nil")
	} else {
		// The error should contain a success message
		if !strings.Contains(err.Error(), "Ok") {
			t.Errorf("Expected error to contain success message, got: %s", err.Error())
		}
	}

	// Check if the torch is now in the character's equipment at WEAR_LIGHT position
	if character.Equipment[types.WEAR_LIGHT] != torchInstance {
		t.Errorf("Expected torch to be in character's equipment at WEAR_LIGHT position")
	}

	// Check if the torch is no longer in inventory
	if len(character.Inventory) != 0 {
		t.Errorf("Expected torch to be removed from inventory")
	}
}

func TestHoldCommand_HoldableItem(t *testing.T) {
	// Create a test character
	character := &types.Character{
		Name:      "TestPlayer",
		Level:     5,
		Position:  types.POS_STANDING,
		Equipment: make([]*types.ObjectInstance, types.NUM_WEARS),
		Inventory: []*types.ObjectInstance{},
	}

	// Create a holdable object
	wand := &types.Object{
		VNUM:      1002,
		Name:      "wand",
		ShortDesc: "a magic wand",
		Type:      types.ITEM_WAND,
		WearFlags: types.ITEM_WEAR_TAKE | types.ITEM_WEAR_HOLD,
	}

	wandInstance := &types.ObjectInstance{
		Prototype: wand,
		CarriedBy: character,
	}

	// Add wand to character's inventory
	character.Inventory = append(character.Inventory, wandInstance)

	// Create a hold command
	holdCmd := &HoldCommand{}

	// Test holding the wand
	err := holdCmd.Execute(character, "wand")
	if err == nil {
		t.Errorf("Expected an error (which contains the success message), got nil")
	} else {
		// The error should contain a success message
		if !strings.Contains(err.Error(), "Ok") {
			t.Errorf("Expected error to contain success message, got: %s", err.Error())
		}
	}

	// Check if the wand is now in the character's equipment at WEAR_HOLD position
	if character.Equipment[types.WEAR_HOLD] != wandInstance {
		t.Errorf("Expected wand to be in character's equipment at WEAR_HOLD position")
	}
}

func TestWieldCommand_Weapon(t *testing.T) {
	// Create a test character
	character := &types.Character{
		Name:      "TestPlayer",
		Level:     5,
		Position:  types.POS_STANDING,
		Equipment: make([]*types.ObjectInstance, types.NUM_WEARS),
		Inventory: []*types.ObjectInstance{},
		Class:     types.CLASS_WARRIOR,
	}

	// Create a weapon object
	sword := &types.Object{
		VNUM:      1003,
		Name:      "sword",
		ShortDesc: "a long sword",
		Type:      types.ITEM_WEAPON,
		WearFlags: types.ITEM_WEAR_TAKE | types.ITEM_WEAR_WIELD,
	}

	swordInstance := &types.ObjectInstance{
		Prototype: sword,
		CarriedBy: character,
	}

	// Add sword to character's inventory
	character.Inventory = append(character.Inventory, swordInstance)

	// Create a wield command
	wieldCmd := &WieldCommand{}

	// Test wielding the sword
	err := wieldCmd.Execute(character, "sword")
	if err == nil {
		t.Errorf("Expected an error (which contains the success message), got nil")
	} else {
		// The error should contain a success message
		if !strings.Contains(err.Error(), "Ok") {
			t.Errorf("Expected error to contain success message, got: %s", err.Error())
		}
	}

	// Check if the sword is now in the character's equipment at WEAR_WIELD position
	if character.Equipment[types.WEAR_WIELD] != swordInstance {
		t.Errorf("Expected sword to be in character's equipment at WEAR_WIELD position")
	}
}

func TestWearCommand_MultipleFingers(t *testing.T) {
	// Create a test character
	character := &types.Character{
		Name:      "TestPlayer",
		Level:     5,
		Position:  types.POS_STANDING,
		Equipment: make([]*types.ObjectInstance, types.NUM_WEARS),
		Inventory: []*types.ObjectInstance{},
		Class:     types.CLASS_WARRIOR,
	}

	// Create ring objects
	ring1 := &types.Object{
		VNUM:      1004,
		Name:      "ring gold",
		ShortDesc: "a gold ring",
		Type:      types.ITEM_WORN,
		WearFlags: types.ITEM_WEAR_TAKE | types.ITEM_WEAR_FINGER,
	}

	ring2 := &types.Object{
		VNUM:      1005,
		Name:      "ring silver",
		ShortDesc: "a silver ring",
		Type:      types.ITEM_WORN,
		WearFlags: types.ITEM_WEAR_TAKE | types.ITEM_WEAR_FINGER,
	}

	ring1Instance := &types.ObjectInstance{
		Prototype: ring1,
		CarriedBy: character,
	}

	ring2Instance := &types.ObjectInstance{
		Prototype: ring2,
		CarriedBy: character,
	}

	// Add rings to character's inventory
	character.Inventory = append(character.Inventory, ring1Instance, ring2Instance)

	// Create a wear command
	wearCmd := &WearCommand{}

	// Test wearing the first ring
	err := wearCmd.Execute(character, "gold")
	if err == nil {
		t.Errorf("Expected an error (which contains the success message), got nil")
	} else {
		// The error should contain a success message about left finger
		if !strings.Contains(err.Error(), "left finger") {
			t.Errorf("Expected error to mention left finger, got: %s", err.Error())
		}
	}

	// Check if the first ring is on the left finger
	if character.Equipment[types.WEAR_FINGER_L] != ring1Instance {
		t.Errorf("Expected first ring to be on left finger")
	}

	// Test wearing the second ring
	err = wearCmd.Execute(character, "silver")
	if err == nil {
		t.Errorf("Expected an error (which contains the success message), got nil")
	} else {
		// The error should contain a success message about right finger
		if !strings.Contains(err.Error(), "right finger") {
			t.Errorf("Expected error to mention right finger, got: %s", err.Error())
		}
	}

	// Check if the second ring is on the right finger
	if character.Equipment[types.WEAR_FINGER_R] != ring2Instance {
		t.Errorf("Expected second ring to be on right finger")
	}
}

func TestWearCommand_MultipleWrists(t *testing.T) {
	// Create a test character
	character := &types.Character{
		Name:      "TestPlayer",
		Level:     5,
		Position:  types.POS_STANDING,
		Equipment: make([]*types.ObjectInstance, types.NUM_WEARS),
		Inventory: []*types.ObjectInstance{},
		Class:     types.CLASS_WARRIOR,
	}

	// Create bracelet objects
	bracelet1 := &types.Object{
		VNUM:      1006,
		Name:      "bracelet gold",
		ShortDesc: "a gold bracelet",
		Type:      types.ITEM_WORN,
		WearFlags: types.ITEM_WEAR_TAKE | types.ITEM_WEAR_WRIST,
	}

	bracelet2 := &types.Object{
		VNUM:      1007,
		Name:      "bracelet silver",
		ShortDesc: "a silver bracelet",
		Type:      types.ITEM_WORN,
		WearFlags: types.ITEM_WEAR_TAKE | types.ITEM_WEAR_WRIST,
	}

	bracelet1Instance := &types.ObjectInstance{
		Prototype: bracelet1,
		CarriedBy: character,
	}

	bracelet2Instance := &types.ObjectInstance{
		Prototype: bracelet2,
		CarriedBy: character,
	}

	// Add bracelets to character's inventory
	character.Inventory = append(character.Inventory, bracelet1Instance, bracelet2Instance)

	// Create a wear command
	wearCmd := &WearCommand{}

	// Test wearing the first bracelet
	err := wearCmd.Execute(character, "gold")
	if err == nil {
		t.Errorf("Expected an error (which contains the success message), got nil")
	} else {
		// The error should contain a success message about left wrist
		if !strings.Contains(err.Error(), "left wrist") {
			t.Errorf("Expected error to mention left wrist, got: %s", err.Error())
		}
	}

	// Check if the first bracelet is on the left wrist
	if character.Equipment[types.WEAR_WRIST_L] != bracelet1Instance {
		t.Errorf("Expected first bracelet to be on left wrist")
	}

	// Test wearing the second bracelet
	err = wearCmd.Execute(character, "silver")
	if err == nil {
		t.Errorf("Expected an error (which contains the success message), got nil")
	} else {
		// The error should contain a success message about right wrist
		if !strings.Contains(err.Error(), "right wrist") {
			t.Errorf("Expected error to mention right wrist, got: %s", err.Error())
		}
	}

	// Check if the second bracelet is on the right wrist
	if character.Equipment[types.WEAR_WRIST_R] != bracelet2Instance {
		t.Errorf("Expected second bracelet to be on right wrist")
	}
}
