package command

import (
	"fmt"
	"strings"

	"github.com/wltechblog/DikuGo/pkg/types"
)

// BuyCommand represents the buy command
type BuyCommand struct{}

// Name returns the name of the command
func (c *BuyCommand) Name() string {
	return "buy"
}

// Aliases returns the aliases of the command
func (c *BuyCommand) Aliases() []string {
	return []string{}
}

// MinPosition returns the minimum position required to execute the command
func (c *BuyCommand) MinPosition() int {
	return types.POS_STANDING
}

// Level returns the minimum level required to execute the command
func (c *BuyCommand) Level() int {
	return 0
}

// LogCommand returns whether the command should be logged
func (c *BuyCommand) LogCommand() bool {
	return false
}

// Execute executes the buy command
func (c *BuyCommand) Execute(ch *types.Character, args string) error {
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
		return fmt.Errorf("Buy what?")
	}

	// Extract the item name
	itemName := args

	// Find the item
	var foundObj *types.Object
	var objVnum int
	for _, vnum := range shop.Producing {
		// Skip invalid items
		if vnum <= 0 {
			continue
		}

		// Get the object prototype
		world, ok := ch.World.(interface{ GetObjectPrototype(int) *types.Object })
		if !ok {
			return fmt.Errorf("world interface not available")
		}
		obj := world.GetObjectPrototype(vnum)
		if obj == nil {
			continue
		}

		// Check if the item matches the name
		if strings.Contains(strings.ToLower(obj.Name), strings.ToLower(itemName)) {
			foundObj = obj
			objVnum = vnum
			break
		}
	}

	// If no item was found, return an error
	if foundObj == nil {
		return fmt.Errorf("%s says, 'Sorry, I don't sell that item.'", keeper.ShortDesc)
	}

	// Calculate the price
	price := int(float64(foundObj.Cost) * shop.ProfitSell)

	// Check if the character has enough gold
	if ch.Gold < price {
		return fmt.Errorf("%s says, 'You don't have enough gold for that!'", keeper.ShortDesc)
	}

	// Create the object
	world, ok := ch.World.(interface {
		CreateObjectFromPrototype(int) *types.ObjectInstance
	})
	if !ok {
		return fmt.Errorf("world interface not available")
	}
	obj := world.CreateObjectFromPrototype(objVnum)
	if obj == nil {
		return fmt.Errorf("Error creating object.")
	}

	// Deduct the gold
	ch.Gold -= price

	// Add the object to the character's inventory
	obj.CarriedBy = ch
	ch.Inventory = append(ch.Inventory, obj)

	// Send messages
	return fmt.Errorf("You buy %s for %d gold.\n\r", obj.Prototype.ShortDesc, price)
}
