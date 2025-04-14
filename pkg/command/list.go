package command

import (
	"fmt"
	"strings"

	"github.com/wltechblog/DikuGo/pkg/types"
)

// ListCommand represents the list command
type ListCommand struct{}

// Name returns the name of the command
func (c *ListCommand) Name() string {
	return "list"
}

// Aliases returns the aliases of the command
func (c *ListCommand) Aliases() []string {
	return []string{}
}

// MinPosition returns the minimum position required to execute the command
func (c *ListCommand) MinPosition() int {
	return types.POS_STANDING
}

// Level returns the minimum level required to execute the command
func (c *ListCommand) Level() int {
	return 0
}

// LogCommand returns whether the command should be logged
func (c *ListCommand) LogCommand() bool {
	return false
}

// Execute executes the list command
func (c *ListCommand) Execute(ch *types.Character, args string) error {
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

	// If no arguments, list all items
	if args == "" {
		return listAllItems(ch, shop, keeper)
	}

	// List specific item
	return listSpecificItem(ch, shop, keeper, args)
}

// isShopOpen checks if the shop is open
func isShopOpen(shop *types.Shop) bool {
	// TODO: Implement time-based shop hours
	return true
}

// listAllItems lists all items in the shop
func listAllItems(ch *types.Character, shop *types.Shop, keeper *types.Character) error {
	// Check if the shop has any items
	if len(shop.Producing) == 0 {
		return fmt.Errorf("%s says, 'Sorry, I don't have anything to sell.'", keeper.ShortDesc)
	}

	// Build the list of items
	var itemList strings.Builder
	itemList.WriteString(fmt.Sprintf("%s's Inventory:\n\r", keeper.ShortDesc))
	itemList.WriteString("Item                                      Price\n\r")
	itemList.WriteString("----------------------------------------------\n\r")

	// Add each item to the list
	for _, objVnum := range shop.Producing {
		// Skip invalid items
		if objVnum <= 0 {
			continue
		}

		// Get the object prototype
		world, ok := ch.World.(interface{ GetObjectPrototype(int) *types.Object })
		if !ok {
			return fmt.Errorf("world interface not available")
		}
		obj := world.GetObjectPrototype(objVnum)
		if obj == nil {
			continue
		}

		// Calculate the price
		price := int(float64(obj.Cost) * shop.ProfitSell)

		// Add the item to the list
		itemList.WriteString(fmt.Sprintf("%-40s %d gold\n\r", obj.ShortDesc, price))
	}

	// Send the list to the character
	return fmt.Errorf("%s", itemList.String())
}

// listSpecificItem lists a specific item in the shop
func listSpecificItem(ch *types.Character, shop *types.Shop, keeper *types.Character, itemName string) error {
	// Check if the shop has any items
	if len(shop.Producing) == 0 {
		return fmt.Errorf("%s says, 'Sorry, I don't have anything to sell.'", keeper.ShortDesc)
	}

	// Find the item
	var foundObj *types.Object
	for _, objVnum := range shop.Producing {
		// Skip invalid items
		if objVnum <= 0 {
			continue
		}

		// Get the object prototype
		world, ok := ch.World.(interface{ GetObjectPrototype(int) *types.Object })
		if !ok {
			return fmt.Errorf("world interface not available")
		}
		obj := world.GetObjectPrototype(objVnum)
		if obj == nil {
			continue
		}

		// Check if the item matches the name
		if strings.Contains(strings.ToLower(obj.Name), strings.ToLower(itemName)) {
			foundObj = obj
			break
		}
	}

	// If no item was found, return an error
	if foundObj == nil {
		return fmt.Errorf("%s says, 'Sorry, I don't sell that item.'", keeper.ShortDesc)
	}

	// Calculate the price
	price := int(float64(foundObj.Cost) * shop.ProfitSell)

	// Send the item info to the character
	return fmt.Errorf("%s costs %d gold.\n\r", foundObj.ShortDesc, price)
}
