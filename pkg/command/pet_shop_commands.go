package command

import (
	"fmt"
	"strings"

	"github.com/wltechblog/DikuGo/pkg/ai"
	"github.com/wltechblog/DikuGo/pkg/types"
)

// ListPetsCommand handles listing pets in a pet shop
type ListPetsCommand struct{}

// Execute executes the list pets command
func (c *ListPetsCommand) Execute(character *types.Character, args string) error {
	// Check if we're in a pet shop
	if !isInPetShop(character) {
		return fmt.Errorf("you are not in a pet shop")
	}

	// Use the pet shop AI to handle the listing
	if ai.PetShopCommand(character, "list", args) {
		return nil
	}

	return fmt.Errorf("no pets are available")
}

// Name returns the name of the command
func (c *ListPetsCommand) Name() string {
	return "list"
}

// Aliases returns the aliases of the command
func (c *ListPetsCommand) Aliases() []string {
	return []string{}
}

// MinPosition returns the minimum position required to execute the command
func (c *ListPetsCommand) MinPosition() int {
	return types.POS_STANDING
}

// Level returns the minimum level required to execute the command
func (c *ListPetsCommand) Level() int {
	return 1
}

// LogCommand returns whether the command should be logged
func (c *ListPetsCommand) LogCommand() bool {
	return false
}

// BuyPetCommand handles buying pets in a pet shop
type BuyPetCommand struct{}

// Execute executes the buy pet command
func (c *BuyPetCommand) Execute(character *types.Character, args string) error {
	// Check if we're in a pet shop
	if !isInPetShop(character) {
		return fmt.Errorf("you are not in a pet shop")
	}

	if args == "" {
		return fmt.Errorf("buy which pet?")
	}

	// Use the pet shop AI to handle the purchase
	if ai.PetShopCommand(character, "buy", args) {
		return nil
	}

	return fmt.Errorf("you can't buy that")
}

// Name returns the name of the command
func (c *BuyPetCommand) Name() string {
	return "buy"
}

// Aliases returns the aliases of the command
func (c *BuyPetCommand) Aliases() []string {
	return []string{}
}

// MinPosition returns the minimum position required to execute the command
func (c *BuyPetCommand) MinPosition() int {
	return types.POS_STANDING
}

// Level returns the minimum level required to execute the command
func (c *BuyPetCommand) Level() int {
	return 1
}

// LogCommand returns whether the command should be logged
func (c *BuyPetCommand) LogCommand() bool {
	return true // Log pet purchases
}

// isInPetShop checks if the character is in a pet shop
func isInPetShop(character *types.Character) bool {
	if character.InRoom == nil {
		return false
	}

	// Check if there's a pet shop keeper in the room
	for _, ch := range character.InRoom.Characters {
		if ch.IsNPC && ch.Prototype != nil {
			// Check if this NPC has pet shop special procedure
			if len(ch.Functions) > 0 {
				// This is a bit tricky to check, but we can look for pet shop characteristics
				// For now, check if the room has an east exit (typical pet shop layout)
				if character.InRoom.Exits[types.DIR_EAST] != nil {
					return true
				}
			}
		}
	}

	// Alternative check: look for keywords in room description
	roomDesc := strings.ToLower(character.InRoom.Description)
	if strings.Contains(roomDesc, "pet") || strings.Contains(roomDesc, "animal") {
		return true
	}

	// Check room name
	roomName := strings.ToLower(character.InRoom.Name)
	if strings.Contains(roomName, "pet") || strings.Contains(roomName, "animal") {
		return true
	}

	return false
}
