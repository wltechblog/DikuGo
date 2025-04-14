package command

import (
	"fmt"
	"strings"

	"github.com/wltechblog/DikuGo/pkg/types"
)

// SellCommand represents the sell command
type SellCommand struct{}

// Name returns the name of the command
func (c *SellCommand) Name() string {
	return "sell"
}

// Aliases returns the aliases of the command
func (c *SellCommand) Aliases() []string {
	return []string{}
}

// MinPosition returns the minimum position required to execute the command
func (c *SellCommand) MinPosition() int {
	return types.POS_STANDING
}

// Level returns the minimum level required to execute the command
func (c *SellCommand) Level() int {
	return 0
}

// LogCommand returns whether the command should be logged
func (c *SellCommand) LogCommand() bool {
	return false
}

// Execute executes the sell command
func (c *SellCommand) Execute(ch *types.Character, args string) error {
	// Check if the character is in a shop
	shop := ch.InRoom.Shop
	if shop == nil {
		return fmt.Errorf("You are not in a shop.")
	}

	// Find the shopkeeper
	var keeper *types.Character
	for _, mob := range ch.InRoom.Characters {
		if mob.IsNPC && mob.Prototype != nil && mob.Prototype.VNUM == shop.MobileVNUM {
			keeper = mob
			break
		}
	}

	if keeper == nil {
		return fmt.Errorf("The shopkeeper is not here.")
	}

	// Check if the shop is open
	if !isShopOpen(shop) {
		return fmt.Errorf("The shop is closed.")
	}

	// Parse arguments
	args = strings.TrimSpace(args)
	if args == "" {
		return fmt.Errorf("Sell what?")
	}

	// Extract the item name
	itemName := args

	// Find the item in the character's inventory
	var foundObj *types.ObjectInstance
	var foundIndex int
	for i, obj := range ch.Inventory {
		if strings.Contains(strings.ToLower(obj.Prototype.Name), strings.ToLower(itemName)) {
			foundObj = obj
			foundIndex = i
			break
		}
	}

	// If no item was found, return an error
	if foundObj == nil {
		return fmt.Errorf("You don't have that item.")
	}

	// Check if the shop buys this type of item
	canBuy := false
	for _, buyType := range shop.BuyTypes {
		if buyType == foundObj.Prototype.Type {
			canBuy = true
			break
		}
	}

	if !canBuy {
		return fmt.Errorf("%s says, 'Sorry, I don't buy that type of item.'", keeper.ShortDesc)
	}

	// Calculate the price
	price := int(float64(foundObj.Prototype.Cost) * shop.ProfitBuy)

	// Remove the object from the character's inventory
	ch.Inventory = append(ch.Inventory[:foundIndex], ch.Inventory[foundIndex+1:]...)
	foundObj.CarriedBy = nil

	// Add the gold
	ch.Gold += price

	// Store the object description before destroying it
	objDesc := foundObj.Prototype.ShortDesc

	// Destroy the object (in a real implementation, we might want to add it to the shop's inventory)
	foundObj = nil

	// Send messages
	return fmt.Errorf("You sell %s for %d gold.\n\r", objDesc, price)
}
