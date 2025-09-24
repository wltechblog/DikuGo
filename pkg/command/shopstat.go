package command

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/wltechblog/DikuGo/pkg/types"
)

// ShopstatCommand represents the shopstat command
type ShopstatCommand struct{}

// Name returns the name of the command
func (c *ShopstatCommand) Name() string {
	return "shopstat"
}

// Aliases returns the aliases of the command
func (c *ShopstatCommand) Aliases() []string {
	return []string{}
}

// MinPosition returns the minimum position required to execute the command
func (c *ShopstatCommand) MinPosition() int {
	return types.POS_DEAD
}

// Level returns the minimum level required to execute the command
func (c *ShopstatCommand) Level() int {
	return 1 // Available to everyone during development
}

// LogCommand returns whether the command should be logged
func (c *ShopstatCommand) LogCommand() bool {
	return true
}

// Execute executes the shopstat command
func (c *ShopstatCommand) Execute(ch *types.Character, args string) error {
	// Get the world from the character
	world, ok := ch.World.(interface{
		GetRoom(int) *types.Room
		GetShop(int) *types.Shop
		GetRooms() []*types.Room
	})
	if !ok {
		return fmt.Errorf("world interface not available")
	}

	// If no arguments, show all shops
	if args == "" {
		return c.showAllShops(ch, world)
	}

	// Parse the shop VNUM from args
	vnumStr := strings.TrimSpace(args)
	vnum, err := strconv.Atoi(vnumStr)
	if err != nil {
		return fmt.Errorf("invalid shop vnum: %s", vnumStr)
	}

	// Get the shop
	shop := world.GetShop(vnum)
	if shop == nil {
		return fmt.Errorf("shop %d does not exist", vnum)
	}

	// Show the shop details
	return c.showShopDetails(ch, shop, world)
}

// showAllShops shows a list of all shops in the game
func (c *ShopstatCommand) showAllShops(ch *types.Character, world interface{
	GetRoom(int) *types.Room
	GetShop(int) *types.Shop
	GetRooms() []*types.Room
}) error {
	// Build a list of all shops
	var sb strings.Builder
	sb.WriteString("Shops in the game:\r\n")
	sb.WriteString("VNUM  Room   Keeper  Name\r\n")
	sb.WriteString("----  ----   ------  ----\r\n")

	// Get all rooms
	rooms := world.GetRooms()
	shopCount := 0

	// Find all rooms with shops
	for _, room := range rooms {
		if room.Shop != nil {
			shopCount++
			roomName := room.Name
			if len(roomName) > 30 {
				roomName = roomName[:27] + "..."
			}
			sb.WriteString(fmt.Sprintf("%-5d %-6d %-7d %s\r\n", room.Shop.VNUM, room.VNUM, room.Shop.MobileVNUM, roomName))
		}
	}

	if shopCount == 0 {
		sb.WriteString("No shops found.\r\n")
	} else {
		sb.WriteString(fmt.Sprintf("\r\nTotal shops: %d\r\n", shopCount))
	}

	return fmt.Errorf("%s", sb.String())
}

// showShopDetails shows the details of a specific shop
func (c *ShopstatCommand) showShopDetails(ch *types.Character, shop *types.Shop, world interface{
	GetRoom(int) *types.Room
}) error {
	// Get the room
	room := world.GetRoom(shop.RoomVNUM)
	if room == nil {
		return fmt.Errorf("shop %d has invalid room VNUM %d", shop.VNUM, shop.RoomVNUM)
	}

	// Build the shop details
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Shop VNUM: %d\r\n", shop.VNUM))
	sb.WriteString(fmt.Sprintf("Room VNUM: %d (%s)\r\n", room.VNUM, room.Name))
	sb.WriteString(fmt.Sprintf("Keeper VNUM: %d\r\n", shop.MobileVNUM))
	sb.WriteString(fmt.Sprintf("Profit Buy: %.2f\r\n", shop.ProfitBuy))
	sb.WriteString(fmt.Sprintf("Profit Sell: %.2f\r\n", shop.ProfitSell))
	sb.WriteString(fmt.Sprintf("Open Hour: %d\r\n", shop.OpenHour))
	sb.WriteString(fmt.Sprintf("Close Hour: %d\r\n", shop.CloseHour))

	// Check if the shopkeeper is in the room
	var keeper *types.Character
	for _, mob := range room.Characters {
		if mob.IsNPC && mob.Prototype != nil && mob.Prototype.VNUM == shop.MobileVNUM {
			keeper = mob
			break
		}
	}

	if keeper != nil {
		sb.WriteString(fmt.Sprintf("Shopkeeper: %s (present in room)\r\n", keeper.Name))
	} else {
		sb.WriteString("Shopkeeper: Not present in room\r\n")
		
		// Log this for debugging
		log.Printf("Shopstat: Shop %d has no shopkeeper in room %d", shop.VNUM, room.VNUM)
		log.Printf("Shopstat: Room %d has %d characters", room.VNUM, len(room.Characters))
		for _, mob := range room.Characters {
			if mob.IsNPC && mob.Prototype != nil {
				log.Printf("Shopstat: Room has NPC %s (VNUM %d)", mob.Name, mob.Prototype.VNUM)
			}
		}
	}

	// Show items the shop sells
	sb.WriteString("\r\nItems for sale:\r\n")
	if len(shop.Producing) == 0 {
		sb.WriteString("None\r\n")
	} else {
		for _, objVnum := range shop.Producing {
			if objVnum <= 0 {
				continue
			}
			sb.WriteString(fmt.Sprintf("  Object VNUM: %d\r\n", objVnum))
		}
	}

	// Show item types the shop buys
	sb.WriteString("\r\nBuys item types:\r\n")
	if len(shop.BuyTypes) == 0 {
		sb.WriteString("None\r\n")
	} else {
		for _, itemType := range shop.BuyTypes {
			sb.WriteString(fmt.Sprintf("  Item Type: %d\r\n", itemType))
		}
	}

	return fmt.Errorf("%s", sb.String())
}
