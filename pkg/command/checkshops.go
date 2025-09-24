package command

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/wltechblog/DikuGo/pkg/storage"
	"github.com/wltechblog/DikuGo/pkg/types"
)

// CheckShopsCommand represents the checkshops command
type CheckShopsCommand struct{}

// Name returns the name of the command
func (c *CheckShopsCommand) Name() string {
	return "checkshops"
}

// Aliases returns the aliases of the command
func (c *CheckShopsCommand) Aliases() []string {
	return []string{}
}

// MinPosition returns the minimum position required to execute the command
func (c *CheckShopsCommand) MinPosition() int {
	return types.POS_DEAD
}

// Level returns the minimum level required to execute the command
func (c *CheckShopsCommand) Level() int {
	return 1 // Available to everyone during development
}

// LogCommand returns whether the command should be logged
func (c *CheckShopsCommand) LogCommand() bool {
	return true
}

// Execute executes the checkshops command
func (c *CheckShopsCommand) Execute(ch *types.Character, args string) error {
	// Get the world from the character
	world, ok := ch.World.(interface {
		GetRoom(int) *types.Room
		GetShop(int) *types.Shop
	})
	if !ok {
		return fmt.Errorf("world interface not available")
	}

	// Parse the shop file directly
	shops, err := storage.ParseShops(filepath.Join("old/lib", "tinyworld.shp"))
	if err != nil {
		return fmt.Errorf("failed to parse shop file: %v", err)
	}

	// Build a report of the shops
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Parsed %d shops directly from file:\r\n", len(shops)))
	sb.WriteString("VNUM  Room   Keeper  Items\r\n")
	sb.WriteString("----  ----   ------  -----\r\n")

	for _, shop := range shops {
		items := ""
		for _, item := range shop.Producing {
			if item > 0 {
				if items != "" {
					items += ", "
				}
				items += fmt.Sprintf("%d", item)
			}
		}
		sb.WriteString(fmt.Sprintf("%-5d %-6d %-7d %s\r\n", shop.VNUM, shop.RoomVNUM, shop.MobileVNUM, items))

		// Check if the shop is in the world
		worldShop := world.GetShop(shop.VNUM)
		if worldShop == nil {
			sb.WriteString(fmt.Sprintf("  WARNING: Shop %d is not in the world\r\n", shop.VNUM))
		} else {
			// Check if the shop is linked to the correct room
			room := world.GetRoom(shop.RoomVNUM)
			if room == nil {
				sb.WriteString(fmt.Sprintf("  WARNING: Shop %d has invalid room VNUM %d\r\n", shop.VNUM, shop.RoomVNUM))
			} else if room.Shop == nil {
				sb.WriteString(fmt.Sprintf("  WARNING: Room %d has no shop linked to it\r\n", shop.RoomVNUM))
			} else if room.Shop.VNUM != shop.VNUM {
				sb.WriteString(fmt.Sprintf("  WARNING: Room %d has shop %d linked to it, not shop %d\r\n",
					shop.RoomVNUM, room.Shop.VNUM, shop.VNUM))
			}
		}
	}

	return fmt.Errorf("%s", sb.String())
}
