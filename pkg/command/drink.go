package command

import (
	"fmt"
	"strings"

	"github.com/wltechblog/DikuGo/pkg/types"
)

// DrinkCommand represents the drink command
type DrinkCommand struct{}

// Execute executes the drink command
func (c *DrinkCommand) Execute(character *types.Character, args string) error {
	// Check if there are arguments
	if args == "" {
		return fmt.Errorf("drink what?")
	}

	// Check if the character is in a room
	if character.InRoom == nil {
		return fmt.Errorf("you are not in a room")
	}

	// Find the object first (any object)
	var targetObject *types.ObjectInstance
	args = strings.ToLower(strings.TrimSpace(args))

	// First check inventory
	for _, obj := range character.Inventory {
		if strings.Contains(strings.ToLower(obj.Prototype.Name), args) ||
			strings.Contains(strings.ToLower(obj.Prototype.ShortDesc), args) {
			targetObject = obj
			break
		}
	}

	// If not found in inventory, check room
	if targetObject == nil {
		for _, obj := range character.InRoom.Objects {
			if strings.Contains(strings.ToLower(obj.Prototype.Name), args) ||
				strings.Contains(strings.ToLower(obj.Prototype.ShortDesc), args) {
				targetObject = obj
				break
			}
		}
	}

	// Check if we found any object
	if targetObject == nil {
		return fmt.Errorf("you can't find it!")
	}

	// Check if it's actually a drink container or fountain
	if targetObject.Prototype.Type != types.ITEM_DRINKCON && targetObject.Prototype.Type != types.ITEM_FOUNTAIN {
		return fmt.Errorf("you can't drink from that!")
	}

	// Now we know it's a drink container
	drinkContainer := targetObject

	// Initialize instance values if not set (check if all values are zero)
	if drinkContainer.Value[0] == 0 && drinkContainer.Value[1] == 0 &&
		drinkContainer.Value[2] == 0 && drinkContainer.Value[3] == 0 {
		// Copy from prototype
		drinkContainer.Value = drinkContainer.Prototype.Value
	}

	// Check if the container has liquid
	liquidAmount := drinkContainer.Value[1]
	if liquidAmount <= 0 {
		return fmt.Errorf("it is empty.")
	}

	// Get liquid type
	liquidType := drinkContainer.Value[2]
	if liquidType < 0 || liquidType >= len(drinkNames) {
		liquidType = 0 // Default to water
	}

	// Check if character is too drunk
	if len(character.Conditions) > types.COND_DRUNK {
		if character.Conditions[types.COND_DRUNK] > 10 {
			return fmt.Errorf("you simply fail to reach your mouth!")
		}
	}

	// Calculate amount to drink (1 unit per drink, not all at once)
	amount := 1 // Drink 1 unit at a time
	if liquidAmount < amount {
		amount = liquidAmount
	}

	// Apply drink effects
	if len(character.Conditions) > types.COND_THIRST {
		// Apply thirst, hunger, and drunk effects based on liquid type
		drunkEffect := getDrinkEffect(liquidType, DRINK_EFFECT_DRUNK)
		fullEffect := getDrinkEffect(liquidType, DRINK_EFFECT_FULL)
		thirstEffect := getDrinkEffect(liquidType, DRINK_EFFECT_THIRST)

		character.Conditions[types.COND_DRUNK] += (drunkEffect * amount) / 4
		character.Conditions[types.COND_FULL] += (fullEffect * amount) / 4
		character.Conditions[types.COND_THIRST] += (thirstEffect * amount) / 4

		// Cap conditions at reasonable limits
		if character.Conditions[types.COND_DRUNK] > 24 {
			character.Conditions[types.COND_DRUNK] = 24
		}
		if character.Conditions[types.COND_FULL] > 24 {
			character.Conditions[types.COND_FULL] = 24
		}
		if character.Conditions[types.COND_THIRST] > 24 {
			character.Conditions[types.COND_THIRST] = 24
		}
	}

	// Reduce liquid in container (unless it's a fountain)
	if drinkContainer.Prototype.Type != types.ITEM_FOUNTAIN {
		drinkContainer.Value[1] -= amount
		if drinkContainer.Value[1] < 0 {
			drinkContainer.Value[1] = 0
		}
	}

	// Send messages
	drinkName := getDrinkName(liquidType)
	character.SendMessage(fmt.Sprintf("You drink the %s.\r\n", drinkName))

	// Send message to room
	for _, ch := range character.InRoom.Characters {
		if ch != character {
			ch.SendMessage(fmt.Sprintf("%s drinks %s from %s.\r\n",
				character.Name, drinkName, drinkContainer.Prototype.ShortDesc))
		}
	}

	// Check for drunk condition
	if len(character.Conditions) > types.COND_DRUNK {
		if character.Conditions[types.COND_DRUNK] > 10 {
			character.SendMessage("You feel drunk.\r\n")
		}
		if character.Conditions[types.COND_THIRST] > 20 {
			character.SendMessage("You do not feel thirsty.\r\n")
		}
		if character.Conditions[types.COND_FULL] > 20 {
			character.SendMessage("You are full.\r\n")
		}
	}

	return nil
}

// Name returns the name of the command
func (c *DrinkCommand) Name() string {
	return "drink"
}

// Aliases returns the aliases of the command
func (c *DrinkCommand) Aliases() []string {
	return []string{}
}

// MinPosition returns the minimum position required to execute the command
func (c *DrinkCommand) MinPosition() int {
	return types.POS_RESTING
}

// Level returns the minimum level required to execute the command
func (c *DrinkCommand) Level() int {
	return 1
}

// LogCommand returns whether the command should be logged
func (c *DrinkCommand) LogCommand() bool {
	return false
}

// Drink names from original DikuMUD
var drinkNames = []string{
	"water", "beer", "wine", "ale", "dark ale", "whisky", "lemonade",
	"firebreather", "local speciality", "slime mold juice", "milk",
	"tea", "coffee", "blood", "salt water", "coca cola",
}

// getDrinkName returns the name of a drink type
func getDrinkName(liquidType int) string {
	if liquidType >= 0 && liquidType < len(drinkNames) {
		return drinkNames[liquidType]
	}
	return "water"
}

// Drink effects constants
const (
	DRINK_EFFECT_DRUNK  = 0
	DRINK_EFFECT_FULL   = 1
	DRINK_EFFECT_THIRST = 2
)

// Drink effects table from original DikuMUD
var drinkEffects = [][]int{
	// DRUNK, FULL, THIRST
	{0, 1, 10}, // water
	{3, 2, 5},  // beer
	{5, 2, 5},  // wine
	{2, 2, 5},  // ale
	{1, 2, 5},  // dark ale
	{6, 1, 4},  // whisky
	{0, 1, 8},  // lemonade
	{10, 0, 0}, // firebreather
	{3, 3, 3},  // local speciality
	{0, 4, -8}, // slime mold juice
	{0, 3, 6},  // milk
	{0, 1, 6},  // tea
	{0, 1, 6},  // coffee
	{0, 2, -1}, // blood
	{0, 1, -2}, // salt water
	{0, 1, 5},  // coca cola
}

// getDrinkEffect returns the effect value for a drink type and effect type
func getDrinkEffect(liquidType, effectType int) int {
	if liquidType >= 0 && liquidType < len(drinkEffects) &&
		effectType >= 0 && effectType < len(drinkEffects[liquidType]) {
		return drinkEffects[liquidType][effectType]
	}
	return 0
}
