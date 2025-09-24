package command

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/wltechblog/DikuGo/pkg/types"
)

// ResetShopsCommand represents the resetshops command
type ResetShopsCommand struct{}

// Name returns the name of the command
func (c *ResetShopsCommand) Name() string {
	return "resetshops"
}

// Aliases returns the aliases of the command
func (c *ResetShopsCommand) Aliases() []string {
	return []string{}
}

// MinPosition returns the minimum position required to execute the command
func (c *ResetShopsCommand) MinPosition() int {
	return types.POS_DEAD
}

// Level returns the minimum level required to execute the command
func (c *ResetShopsCommand) Level() int {
	return 1 // Available to everyone during development
}

// LogCommand returns whether the command should be logged
func (c *ResetShopsCommand) LogCommand() bool {
	return true
}

// Execute executes the resetshops command
func (c *ResetShopsCommand) Execute(ch *types.Character, args string) error {
	// Get the world from the character
	world, ok := ch.World.(interface {
		GetRooms() []*types.Room
		GetRoom(int) *types.Room
		GetShop(int) *types.Shop
		GetMobile(int) *types.Mobile
		CreateMobFromPrototype(int, *types.Room) *types.Character
	})
	if !ok {
		return fmt.Errorf("world interface not available")
	}

	// Check if a specific shop VNUM was provided
	specificShop := -1
	if args != "" {
		var err error
		specificShop, err = strconv.Atoi(args)
		if err != nil {
			return fmt.Errorf("Invalid shop VNUM: %s", args)
		}
	}

	// Get all rooms
	rooms := world.GetRooms()
	shopCount := 0
	shopkeeperCount := 0
	missingMobProtos := make(map[int]bool)

	// Find all rooms with shops
	for _, room := range rooms {
		if room.Shop != nil {
			// Skip if a specific shop was requested and this isn't it
			if specificShop != -1 && room.Shop.VNUM != specificShop {
				continue
			}

			shopCount++
			log.Printf("ResetShops: Found shop %d in room %d with keeper VNUM %d",
				room.Shop.VNUM, room.VNUM, room.Shop.MobileVNUM)

			// Check if the mobile prototype exists
			mobProto := world.GetMobile(room.Shop.MobileVNUM)
			if mobProto == nil {
				log.Printf("ResetShops: Warning - Mobile prototype %d not found for shop %d",
					room.Shop.MobileVNUM, room.Shop.VNUM)
				missingMobProtos[room.Shop.MobileVNUM] = true
				continue
			}

			// Check if the shopkeeper is already in the room
			shopkeeperExists := false
			for _, mob := range room.Characters {
				if mob.IsNPC && mob.Prototype != nil && mob.Prototype.VNUM == room.Shop.MobileVNUM {
					shopkeeperExists = true
					log.Printf("ResetShops: Shopkeeper %d (%s) already exists in room %d",
						room.Shop.MobileVNUM, mob.ShortDesc, room.VNUM)
					break
				}
			}

			// If the shopkeeper doesn't exist, create it
			if !shopkeeperExists {
				// Create the shopkeeper
				mob := world.CreateMobFromPrototype(room.Shop.MobileVNUM, room)
				if mob != nil {
					shopkeeperCount++
					log.Printf("ResetShops: Created shopkeeper %s (VNUM %d) in room %d",
						mob.ShortDesc, room.Shop.MobileVNUM, room.VNUM)
				} else {
					log.Printf("ResetShops: Failed to create shopkeeper %d in room %d",
						room.Shop.MobileVNUM, room.VNUM)
				}
			}
		}
	}

	// Build the result message
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Reset %d shops, created %d shopkeepers.\r\n", shopCount, shopkeeperCount))

	// Report any missing mobile prototypes
	if len(missingMobProtos) > 0 {
		sb.WriteString("\r\nWARNING: The following mobile prototypes are missing:\r\n")
		for mobVnum := range missingMobProtos {
			sb.WriteString(fmt.Sprintf("- Mobile VNUM %d\r\n", mobVnum))
		}
		sb.WriteString("\r\nThese shopkeepers could not be created.\r\n")
	}

	return fmt.Errorf("%s", sb.String())
}
