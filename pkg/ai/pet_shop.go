package ai

import (
	"fmt"
	"log"
	"strings"

	"github.com/wltechblog/DikuGo/pkg/types"
)

// petShopProc is the special procedure for pet shops
func petShopProc(ch *types.Character, argument string) bool {
	// Skip if the character is not an NPC
	if !ch.IsNPC {
		return false
	}

	// Skip if the character is fighting
	if ch.Fighting != nil {
		return false
	}

	// This is a pet shop keeper - handle pet shop commands
	return true
}

// PetShopCommand handles pet shop interactions
func PetShopCommand(ch *types.Character, cmd string, argument string) bool {
	// Find the pet shop room (usually one room east of the shop)
	petRoom := findPetRoom(ch.InRoom)
	if petRoom == nil {
		return false
	}

	switch strings.ToLower(cmd) {
	case "list":
		return listPets(ch, petRoom)
	case "buy":
		return buyPet(ch, petRoom, argument)
	default:
		return false
	}
}

// findPetRoom finds the room where pets are kept
func findPetRoom(shopRoom *types.Room) *types.Room {
	// For now, just return the shop room itself as the pet room
	// In a full implementation, we would need access to the world to resolve DestVnum
	// This is a simplified version for the pet shop functionality
	return shopRoom
}

// isPet determines if a character is a pet that can be bought
func isPet(ch *types.Character) bool {
	if !ch.IsNPC || ch.Prototype == nil {
		return false
	}

	// Pets are usually low-level, non-aggressive NPCs
	// In original DikuMUD, pets have specific characteristics
	if ch.Level <= 10 && ch.ActFlags&types.ACT_AGGRESSIVE == 0 {
		// Check if it's a typical pet (dog, cat, bird, etc.)
		name := strings.ToLower(ch.Name)
		if strings.Contains(name, "dog") || strings.Contains(name, "cat") ||
			strings.Contains(name, "bird") || strings.Contains(name, "rabbit") ||
			strings.Contains(name, "puppy") || strings.Contains(name, "kitten") {
			return true
		}
	}

	return false
}

// listPets shows available pets and their prices
func listPets(ch *types.Character, petRoom *types.Room) bool {
	ch.SendMessage("Available pets are:\r\n")

	found := false
	for _, pet := range petRoom.Characters {
		if pet.IsNPC && isPet(pet) {
			// Price is 3 times the pet's experience value (original DikuMUD formula)
			price := pet.Experience * 3
			if price <= 0 {
				price = 100 // Minimum price
			}

			ch.SendMessage(fmt.Sprintf("%8d - %s\r\n", price, pet.ShortDesc))
			found = true
		}
	}

	if !found {
		ch.SendMessage("No pets are currently available.\r\n")
	}

	return true
}

// buyPet handles buying a pet
func buyPet(ch *types.Character, petRoom *types.Room, argument string) bool {
	args := strings.Fields(argument)
	if len(args) == 0 {
		ch.SendMessage("Buy which pet?\r\n")
		return true
	}

	petName := args[0]
	var newPetName string
	if len(args) > 1 {
		newPetName = args[1]
	}

	// Find the pet
	var pet *types.Character
	for _, character := range petRoom.Characters {
		if character.IsNPC && isPet(character) {
			if strings.Contains(strings.ToLower(character.Name), strings.ToLower(petName)) ||
				strings.Contains(strings.ToLower(character.ShortDesc), strings.ToLower(petName)) {
				pet = character
				break
			}
		}
	}

	if pet == nil {
		ch.SendMessage("There is no such pet!\r\n")
		return true
	}

	// Calculate price
	price := pet.Experience * 3
	if price <= 0 {
		price = 100
	}

	// Check if player has enough gold
	if ch.Gold < price {
		ch.SendMessage("You don't have enough gold!\r\n")
		return true
	}

	// Check if player already has a pet
	if hasFollower(ch) {
		ch.SendMessage("You already have a pet following you!\r\n")
		return true
	}

	// Create a copy of the pet for the player
	newPet := createPetCopy(pet, ch)
	if newPet == nil {
		ch.SendMessage("Sorry, that pet is not available right now.\r\n")
		return true
	}

	// Set custom name if provided
	if newPetName != "" {
		newPet.Name = newPetName
		newPet.ShortDesc = fmt.Sprintf("%s the %s", newPetName, getBasePetType(pet))
	}

	// Deduct gold
	ch.Gold -= price

	// Move pet to player's room
	if newPet.InRoom != nil {
		// Remove from current room
		for i, char := range newPet.InRoom.Characters {
			if char == newPet {
				newPet.InRoom.Characters = append(newPet.InRoom.Characters[:i], newPet.InRoom.Characters[i+1:]...)
				break
			}
		}
	}

	// Add to player's room
	newPet.InRoom = ch.InRoom
	ch.InRoom.Characters = append(ch.InRoom.Characters, newPet)

	// Set up following relationship
	newPet.Master = ch
	newPet.ActFlags |= types.ACT_FOLLOWER

	// Messages
	ch.SendMessage(fmt.Sprintf("You buy %s for %d gold coins.\r\n", newPet.ShortDesc, price))
	ch.SendMessage(fmt.Sprintf("%s starts following you.\r\n", newPet.ShortDesc))

	// Message to room
	for _, character := range ch.InRoom.Characters {
		if character != ch && character != newPet {
			character.SendMessage(fmt.Sprintf("%s buys %s.\r\n", ch.Name, newPet.ShortDesc))
		}
	}

	log.Printf("Player %s bought pet %s for %d gold", ch.Name, newPet.ShortDesc, price)

	return true
}

// hasFollower checks if a character already has a follower
func hasFollower(ch *types.Character) bool {
	if ch.InRoom == nil {
		return false
	}

	for _, character := range ch.InRoom.Characters {
		if character.IsNPC && character.Master == ch {
			return true
		}
	}

	return false
}

// createPetCopy creates a copy of a pet for a player
func createPetCopy(original *types.Character, owner *types.Character) *types.Character {
	if original.Prototype == nil {
		return nil
	}

	// Create a new character based on the original
	pet := &types.Character{
		Name:          original.Name,
		ShortDesc:     original.ShortDesc,
		LongDesc:      original.LongDesc,
		Description:   original.Description,
		Level:         original.Level,
		Sex:           original.Sex,
		Class:         original.Class,
		Race:          original.Race,
		Position:      types.POS_STANDING,
		Gold:          0, // Pets don't carry gold
		Experience:    original.Experience,
		Alignment:     original.Alignment,
		HP:            original.HP,
		MaxHitPoints:  original.MaxHitPoints,
		ManaPoints:    original.ManaPoints,
		MaxManaPoints: original.MaxManaPoints,
		MovePoints:    original.MovePoints,
		MaxMovePoints: original.MaxMovePoints,
		ArmorClass:    original.ArmorClass,
		HitRoll:       original.HitRoll,
		DamRoll:       original.DamRoll,
		Abilities:     original.Abilities,
		IsNPC:         true,
		ActFlags:      original.ActFlags | types.ACT_FOLLOWER,
		Prototype:     original.Prototype,
		Master:        owner,
		Equipment:     make([]*types.ObjectInstance, types.NUM_WEARS),
		Inventory:     make([]*types.ObjectInstance, 0),
		Skills:        make(map[int]int),
		Spells:        make(map[int]int),
	}

	return pet
}

// getBasePetType extracts the base type of pet (dog, cat, etc.)
func getBasePetType(pet *types.Character) string {
	name := strings.ToLower(pet.Name)

	if strings.Contains(name, "dog") || strings.Contains(name, "puppy") {
		return "dog"
	}
	if strings.Contains(name, "cat") || strings.Contains(name, "kitten") {
		return "cat"
	}
	if strings.Contains(name, "bird") {
		return "bird"
	}
	if strings.Contains(name, "rabbit") {
		return "rabbit"
	}

	// Default to the first word of the name
	words := strings.Fields(name)
	if len(words) > 0 {
		return words[0]
	}

	return "pet"
}

// RegisterPetShop registers the pet shop special procedure
func RegisterPetShop() {
	SpecialProcs["pet_shop"] = petShopProc
	SpecialProcs["petshop"] = petShopProc
}
