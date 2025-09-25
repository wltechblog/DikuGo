package command

import (
	"fmt"

	"github.com/wltechblog/DikuGo/pkg/types"
)

// WieldCommand represents the wield command
type WieldCommand struct{}

// Execute executes the wield command
func (c *WieldCommand) Execute(character *types.Character, args string) error {
	// Check if there are arguments
	if args == "" {
		return fmt.Errorf("wield what?")
	}

	// Find the target object
	obj := findObjectInInventory(character, args)
	if obj == nil {
		return fmt.Errorf("you don't have %s", args)
	}

	// Check if the object can be wielded
	if obj.Prototype.WearFlags&types.ITEM_WEAR_WIELD == 0 {
		return fmt.Errorf("you can't wield %s", obj.Prototype.ShortDesc)
	}

	// Check class restrictions
	if !canWearItem(character, obj.Prototype) {
		return fmt.Errorf("you are not the right class to use %s", obj.Prototype.ShortDesc)
	}

	// Check if already wielding something
	if character.Equipment[types.WEAR_WIELD] != nil {
		return fmt.Errorf("you are already wielding something")
	}

	// Check strength requirement for weapons
	if obj.Prototype.Type == types.ITEM_WEAPON {
		// TODO: Implement strength check based on weapon weight
		// For now, we'll skip this check
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
	obj.WornOn = types.WEAR_WIELD
	character.Equipment[types.WEAR_WIELD] = obj

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

// Name returns the name of the command
func (c *WieldCommand) Name() string {
	return "wield"
}

// Aliases returns the aliases of the command
func (c *WieldCommand) Aliases() []string {
	return []string{}
}

// MinPosition returns the minimum position required to execute the command
func (c *WieldCommand) MinPosition() int {
	return types.POS_RESTING
}

// Level returns the minimum level required to execute the command
func (c *WieldCommand) Level() int {
	return 0
}

// LogCommand returns whether the command should be logged
func (c *WieldCommand) LogCommand() bool {
	return false
}
