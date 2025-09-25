package command

import (
	"fmt"

	"github.com/wltechblog/DikuGo/pkg/types"
)

// HoldCommand represents the hold command
type HoldCommand struct{}

// Execute executes the hold command
func (c *HoldCommand) Execute(character *types.Character, args string) error {
	// Check if there are arguments
	if args == "" {
		return fmt.Errorf("hold what?")
	}

	// Find the target object
	obj := findObjectInInventory(character, args)
	if obj == nil {
		return fmt.Errorf("you don't have %s", args)
	}

	// Check class restrictions
	if !canWearItem(character, obj.Prototype) {
		return fmt.Errorf("you are not the right class to use %s", obj.Prototype.ShortDesc)
	}

	// Check if it's a light source
	if obj.Prototype.Type == types.ITEM_LIGHT {
		return c.holdLight(character, obj)
	}

	// Check if it can be held
	if obj.Prototype.WearFlags&types.ITEM_WEAR_HOLD == 0 {
		return fmt.Errorf("you can't hold %s", obj.Prototype.ShortDesc)
	}

	// Check if already holding something
	if character.Equipment[types.WEAR_HOLD] != nil {
		return fmt.Errorf("you are already holding something")
	}

	// Check for alignment restrictions
	if (obj.Prototype.ExtraFlags&types.ITEM_ANTI_EVIL != 0 && character.Alignment < 0) ||
		(obj.Prototype.ExtraFlags&types.ITEM_ANTI_GOOD != 0 && character.Alignment > 0) ||
		(obj.Prototype.ExtraFlags&types.ITEM_ANTI_NEUTRAL != 0 && character.Alignment == 0) {
		return fmt.Errorf("you are zapped by %s and instantly drop it.\r\n", obj.Prototype.ShortDesc)
	}

	// Remove the object from the character's inventory
	for i, o := range character.Inventory {
		if o == obj {
			character.Inventory = append(character.Inventory[:i], character.Inventory[i+1:]...)
			break
		}
	}

	// Add the object to the character's equipment
	obj.CarriedBy = nil
	obj.WornBy = character
	obj.WornOn = types.WEAR_HOLD
	character.Equipment[types.WEAR_HOLD] = obj

	// Apply magical effects from the object
	world := character.World
	if world != nil {
		// Try to use the ApplyObjectAffects method if it exists
		if applier, ok := world.(interface {
			ApplyObjectAffects(*types.Character, *types.ObjectInstance, bool)
		}); ok {
			applier.ApplyObjectAffects(character, obj, true)
		}
	}

	// Send a message to the character
	return fmt.Errorf("Ok.\r\n")
}

// holdLight handles holding a light source
func (c *HoldCommand) holdLight(character *types.Character, obj *types.ObjectInstance) error {
	// Check if already holding a light
	if character.Equipment[types.WEAR_LIGHT] != nil {
		return fmt.Errorf("you are already holding a light source")
	}

	// Check for alignment restrictions
	if (obj.Prototype.ExtraFlags&types.ITEM_ANTI_EVIL != 0 && character.Alignment < 0) ||
		(obj.Prototype.ExtraFlags&types.ITEM_ANTI_GOOD != 0 && character.Alignment > 0) ||
		(obj.Prototype.ExtraFlags&types.ITEM_ANTI_NEUTRAL != 0 && character.Alignment == 0) {
		return fmt.Errorf("you are zapped by %s and instantly drop it.\r\n", obj.Prototype.ShortDesc)
	}

	// Remove the object from the character's inventory
	for i, o := range character.Inventory {
		if o == obj {
			character.Inventory = append(character.Inventory[:i], character.Inventory[i+1:]...)
			break
		}
	}

	// Add the object to the character's equipment
	obj.CarriedBy = nil
	obj.WornBy = character
	obj.WornOn = types.WEAR_LIGHT
	character.Equipment[types.WEAR_LIGHT] = obj

	// Apply magical effects from the object
	world := character.World
	if world != nil {
		// Try to use the ApplyObjectAffects method if it exists
		if applier, ok := world.(interface {
			ApplyObjectAffects(*types.Character, *types.ObjectInstance, bool)
		}); ok {
			applier.ApplyObjectAffects(character, obj, true)
		}
	}

	// Increase room light if the light source is lit
	if obj.Prototype.Value[2] > 0 && character.InRoom != nil {
		// TODO: Implement room lighting system
		// character.InRoom.Light++
	}

	// Send a message to the character
	return fmt.Errorf("Ok.\r\n")
}

// Name returns the name of the command
func (c *HoldCommand) Name() string {
	return "hold"
}

// Aliases returns the aliases of the command
func (c *HoldCommand) Aliases() []string {
	return []string{"grab"}
}

// MinPosition returns the minimum position required to execute the command
func (c *HoldCommand) MinPosition() int {
	return types.POS_RESTING
}

// Level returns the minimum level required to execute the command
func (c *HoldCommand) Level() int {
	return 0
}

// LogCommand returns whether the command should be logged
func (c *HoldCommand) LogCommand() bool {
	return false
}
