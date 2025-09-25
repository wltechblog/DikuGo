package command

import (
	"fmt"
	"strings"

	"github.com/wltechblog/DikuGo/pkg/types"
)

// WearCommand represents the wear command
type WearCommand struct{}

// Execute executes the wear command
func (c *WearCommand) Execute(character *types.Character, args string) error {
	// Check if there are arguments
	if args == "" {
		return fmt.Errorf("wear what?")
	}

	// Check if the target is "all"
	if strings.ToLower(args) == "all" {
		return c.wearAll(character)
	}

	// Find the target object
	obj := findObjectInInventory(character, args)
	if obj == nil {
		return fmt.Errorf("you don't have %s", args)
	}

	// Check if the object can be worn
	if obj.Prototype.WearFlags == 0 {
		return fmt.Errorf("you can't wear %s", obj.Prototype.ShortDesc)
	}

	// Check class restrictions
	if !canWearItem(character, obj.Prototype) {
		return fmt.Errorf("you are not the right class to use %s", obj.Prototype.ShortDesc)
	}

	// Find a position to wear the object (with multiple position support)
	position := findBestWearPosition(character, obj.Prototype.WearFlags)
	if position < 0 {
		return fmt.Errorf("you can't wear %s", obj.Prototype.ShortDesc)
	}

	// Check if the character is already wearing something in that position
	if character.Equipment[position] != nil {
		return fmt.Errorf("you're already wearing something on your %s", wearPositionName(position))
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
	obj.WornOn = position
	character.Equipment[position] = obj

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

	// Send a message to the character with position-specific messages
	return fmt.Errorf("%s", getWearMessage(obj, position))
}

// wearAll wears all wearable objects in the character's inventory
func (c *WearCommand) wearAll(character *types.Character) error {
	// Check if the character has any objects
	if len(character.Inventory) == 0 {
		return fmt.Errorf("you are not carrying anything.\r\n")
	}

	// Try to wear each object
	var sb strings.Builder
	for _, obj := range character.Inventory {
		// Check if the object can be worn
		if obj.Prototype.WearFlags == 0 {
			continue
		}

		// Check class restrictions
		if !canWearItem(character, obj.Prototype) {
			continue
		}

		// Find a position to wear the object (with multiple position support)
		position := findBestWearPosition(character, obj.Prototype.WearFlags)
		if position < 0 {
			continue
		}

		// Check if the character is already wearing something in that position
		if character.Equipment[position] != nil {
			continue
		}

		// Check for alignment restrictions
		if (obj.Prototype.ExtraFlags&types.ITEM_ANTI_EVIL != 0 && character.Alignment < 0) ||
			(obj.Prototype.ExtraFlags&types.ITEM_ANTI_GOOD != 0 && character.Alignment > 0) ||
			(obj.Prototype.ExtraFlags&types.ITEM_ANTI_NEUTRAL != 0 && character.Alignment == 0) {
			sb.WriteString(fmt.Sprintf("You are zapped by %s and instantly drop it.\r\n", obj.Prototype.ShortDesc))
			continue
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
		obj.WornOn = position
		character.Equipment[position] = obj

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

		// Add a message to the buffer
		sb.WriteString(getWearMessage(obj, position))
	}

	// If nothing was worn, return a message
	if sb.Len() == 0 {
		return fmt.Errorf("you have nothing you can wear.\r\n")
	}

	// Send the messages to the character
	return fmt.Errorf("%s", sb.String())
}

// Name returns the name of the command
func (c *WearCommand) Name() string {
	return "wear"
}

// Aliases returns the aliases of the command
func (c *WearCommand) Aliases() []string {
	return []string{"wield", "hold"}
}

// MinPosition returns the minimum position required to execute the command
func (c *WearCommand) MinPosition() int {
	return types.POS_STANDING
}

// Level returns the minimum level required to execute the command
func (c *WearCommand) Level() int {
	return 0
}

// LogCommand returns whether the command should be logged
func (c *WearCommand) LogCommand() bool {
	return false
}

// findBestWearPosition finds the best available wear position for the given wear flags
// This function handles multiple positions for fingers and wrists like the original DikuMUD
func findBestWearPosition(character *types.Character, wearFlags uint32) int {
	// Check each wear flag in priority order, handling multiple positions
	if wearFlags&types.ITEM_WEAR_FINGER != 0 {
		// Try left finger first, then right finger
		if character.Equipment[types.WEAR_FINGER_L] == nil {
			return types.WEAR_FINGER_L
		}
		if character.Equipment[types.WEAR_FINGER_R] == nil {
			return types.WEAR_FINGER_R
		}
		return -1 // Both fingers occupied
	}
	if wearFlags&types.ITEM_WEAR_NECK != 0 {
		// Try first neck position, then second
		if character.Equipment[types.WEAR_NECK_1] == nil {
			return types.WEAR_NECK_1
		}
		if character.Equipment[types.WEAR_NECK_2] == nil {
			return types.WEAR_NECK_2
		}
		return -1 // Both neck positions occupied
	}
	if wearFlags&types.ITEM_WEAR_WRIST != 0 {
		// Try left wrist first, then right wrist
		if character.Equipment[types.WEAR_WRIST_L] == nil {
			return types.WEAR_WRIST_L
		}
		if character.Equipment[types.WEAR_WRIST_R] == nil {
			return types.WEAR_WRIST_R
		}
		return -1 // Both wrists occupied
	}

	// Single position items
	if wearFlags&types.ITEM_WEAR_BODY != 0 && character.Equipment[types.WEAR_BODY] == nil {
		return types.WEAR_BODY
	}
	if wearFlags&types.ITEM_WEAR_HEAD != 0 && character.Equipment[types.WEAR_HEAD] == nil {
		return types.WEAR_HEAD
	}
	if wearFlags&types.ITEM_WEAR_LEGS != 0 && character.Equipment[types.WEAR_LEGS] == nil {
		return types.WEAR_LEGS
	}
	if wearFlags&types.ITEM_WEAR_FEET != 0 && character.Equipment[types.WEAR_FEET] == nil {
		return types.WEAR_FEET
	}
	if wearFlags&types.ITEM_WEAR_HANDS != 0 && character.Equipment[types.WEAR_HANDS] == nil {
		return types.WEAR_HANDS
	}
	if wearFlags&types.ITEM_WEAR_ARMS != 0 && character.Equipment[types.WEAR_ARMS] == nil {
		return types.WEAR_ARMS
	}
	if wearFlags&types.ITEM_WEAR_SHIELD != 0 && character.Equipment[types.WEAR_SHIELD] == nil {
		return types.WEAR_SHIELD
	}
	if wearFlags&types.ITEM_WEAR_ABOUT != 0 && character.Equipment[types.WEAR_ABOUT] == nil {
		return types.WEAR_ABOUT
	}
	if wearFlags&types.ITEM_WEAR_WAIST != 0 && character.Equipment[types.WEAR_WAIST] == nil {
		return types.WEAR_WAIST
	}
	if wearFlags&types.ITEM_WEAR_WIELD != 0 && character.Equipment[types.WEAR_WIELD] == nil {
		return types.WEAR_WIELD
	}
	if wearFlags&types.ITEM_WEAR_HOLD != 0 && character.Equipment[types.WEAR_HOLD] == nil {
		return types.WEAR_HOLD
	}

	// No valid position found
	return -1
}

// findWearPosition finds a position to wear an object
func findWearPosition(wearFlags uint32) int {
	// Check each wear position
	if wearFlags&types.ITEM_WEAR_FINGER != 0 {
		return types.WEAR_FINGER_R
	}
	if wearFlags&types.ITEM_WEAR_NECK != 0 {
		return types.WEAR_NECK_1
	}
	if wearFlags&types.ITEM_WEAR_BODY != 0 {
		return types.WEAR_BODY
	}
	if wearFlags&types.ITEM_WEAR_HEAD != 0 {
		return types.WEAR_HEAD
	}
	if wearFlags&types.ITEM_WEAR_LEGS != 0 {
		return types.WEAR_LEGS
	}
	if wearFlags&types.ITEM_WEAR_FEET != 0 {
		return types.WEAR_FEET
	}
	if wearFlags&types.ITEM_WEAR_HANDS != 0 {
		return types.WEAR_HANDS
	}
	if wearFlags&types.ITEM_WEAR_ARMS != 0 {
		return types.WEAR_ARMS
	}
	if wearFlags&types.ITEM_WEAR_SHIELD != 0 {
		return types.WEAR_SHIELD
	}
	if wearFlags&types.ITEM_WEAR_ABOUT != 0 {
		return types.WEAR_ABOUT
	}
	if wearFlags&types.ITEM_WEAR_WAIST != 0 {
		return types.WEAR_WAIST
	}
	if wearFlags&types.ITEM_WEAR_WRIST != 0 {
		return types.WEAR_WRIST_R
	}
	if wearFlags&types.ITEM_WEAR_WIELD != 0 {
		return types.WEAR_WIELD
	}
	if wearFlags&types.ITEM_WEAR_HOLD != 0 {
		return types.WEAR_HOLD
	}

	// No valid position found
	return -1
}

// canWearItem checks if a character can wear an item based on class restrictions
func canWearItem(character *types.Character, obj *types.Object) bool {
	// If the character is an NPC, they can wear anything
	if character.IsNPC {
		return true
	}

	// Use the CanClassUseItem function from the types package
	return types.CanClassUseItem(character.Class, obj)
}

// getWearMessage returns the appropriate wear message based on position
func getWearMessage(obj *types.ObjectInstance, position int) string {
	switch position {
	case types.WEAR_FINGER_L:
		return fmt.Sprintf("You put %s on your left finger.\r\n", obj.Prototype.ShortDesc)
	case types.WEAR_FINGER_R:
		return fmt.Sprintf("You put %s on your right finger.\r\n", obj.Prototype.ShortDesc)
	case types.WEAR_WRIST_L:
		return fmt.Sprintf("You wear %s around your left wrist.\r\n", obj.Prototype.ShortDesc)
	case types.WEAR_WRIST_R:
		return fmt.Sprintf("You wear %s around your right wrist.\r\n", obj.Prototype.ShortDesc)
	case types.WEAR_NECK_1, types.WEAR_NECK_2:
		return fmt.Sprintf("Ok.\r\n")
	case types.WEAR_SHIELD:
		return fmt.Sprintf("You start using %s.\r\n", obj.Prototype.ShortDesc)
	default:
		return fmt.Sprintf("You wear %s on your %s.\r\n", obj.Prototype.ShortDesc, wearPositionName(position))
	}
}

// wearPositionName returns the name of a wear position
func wearPositionName(position int) string {
	switch position {
	case types.WEAR_LIGHT:
		return "light"
	case types.WEAR_FINGER_R:
		return "right finger"
	case types.WEAR_FINGER_L:
		return "left finger"
	case types.WEAR_NECK_1:
		return "neck"
	case types.WEAR_NECK_2:
		return "neck"
	case types.WEAR_BODY:
		return "body"
	case types.WEAR_HEAD:
		return "head"
	case types.WEAR_LEGS:
		return "legs"
	case types.WEAR_FEET:
		return "feet"
	case types.WEAR_HANDS:
		return "hands"
	case types.WEAR_ARMS:
		return "arms"
	case types.WEAR_SHIELD:
		return "shield"
	case types.WEAR_ABOUT:
		return "about body"
	case types.WEAR_WAIST:
		return "waist"
	case types.WEAR_WRIST_R:
		return "right wrist"
	case types.WEAR_WRIST_L:
		return "left wrist"
	case types.WEAR_WIELD:
		return "wielded"
	case types.WEAR_HOLD:
		return "held"
	default:
		return "unknown"
	}
}
