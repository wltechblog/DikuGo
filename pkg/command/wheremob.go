package command

import (
	"fmt"
	"strings"

	"github.com/wltechblog/DikuGo/pkg/types"
)

// WhereMobCommand represents the wheremob command
type WhereMobCommand struct{}

// Name returns the name of the command
func (c *WhereMobCommand) Name() string {
	return "wheremob"
}

// Aliases returns the aliases of the command
func (c *WhereMobCommand) Aliases() []string {
	return []string{"findmob"}
}

// MinPosition returns the minimum position required to execute the command
func (c *WhereMobCommand) MinPosition() int {
	return types.POS_DEAD
}

// Level returns the minimum level required to execute the command
func (c *WhereMobCommand) Level() int {
	return 1 // Available to everyone during development
}

// LogCommand returns whether the command should be logged
func (c *WhereMobCommand) LogCommand() bool {
	return true
}

// Execute executes the wheremob command
func (c *WhereMobCommand) Execute(ch *types.Character, args string) error {
	// Check if a search term was provided
	if args == "" {
		return fmt.Errorf("Usage: wheremob <name or keyword>\r\nExample: wheremob grocer")
	}

	// Get the world from the character
	world, ok := ch.World.(interface {
		GetCharacters() map[string]*types.Character
		GetMobilePrototypes() []*types.Mobile
	})
	if !ok {
		return fmt.Errorf("world interface not available")
	}

	// Get all characters in the world
	characters := world.GetCharacters()

	// Get all mobile prototypes for reference
	mobilePrototypes := world.GetMobilePrototypes()

	// Build a map of VNUM to prototype name for quick lookup
	prototypeNames := make(map[int]string)
	for _, proto := range mobilePrototypes {
		prototypeNames[proto.VNUM] = proto.Name
	}

	// Search term (case insensitive)
	searchTerm := strings.ToLower(args)

	// Find mobs that match the search term
	var matchedMobs []*types.Character
	for _, character := range characters {
		if !character.IsNPC {
			continue // Skip players
		}

		// Check if the mob name contains the search term
		if strings.Contains(strings.ToLower(character.Name), searchTerm) {
			matchedMobs = append(matchedMobs, character)
			continue
		}

		// Check if the mob short description contains the search term
		if strings.Contains(strings.ToLower(character.ShortDesc), searchTerm) {
			matchedMobs = append(matchedMobs, character)
			continue
		}

		// Check if the mob prototype name contains the search term (if available)
		if character.Prototype != nil {
			protoVnum := character.Prototype.VNUM
			if protoName, ok := prototypeNames[protoVnum]; ok {
				if strings.Contains(strings.ToLower(protoName), searchTerm) {
					matchedMobs = append(matchedMobs, character)
					continue
				}
			}
		}
	}

	// If no mobs were found, check if any prototypes match the search term
	// This helps find mobs that aren't currently in the game
	if len(matchedMobs) == 0 {
		var protoMatches []string
		for _, proto := range mobilePrototypes {
			if strings.Contains(strings.ToLower(proto.Name), searchTerm) ||
				strings.Contains(strings.ToLower(proto.ShortDesc), searchTerm) {
				protoMatches = append(protoMatches, fmt.Sprintf("- %s (VNUM: %d) - not currently in game", 
					proto.ShortDesc, proto.VNUM))
			}
		}

		if len(protoMatches) > 0 {
			result := "No active mobs found, but these prototypes match your search:\r\n"
			result += strings.Join(protoMatches, "\r\n")
			return fmt.Errorf("%s", result)
		}

		return fmt.Errorf("No mobs found matching '%s'", args)
	}

	// Build the result string
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Mobs matching '%s':\r\n", args))
	
	for _, mob := range matchedMobs {
		roomName := "Unknown"
		roomVnum := 0
		if mob.InRoom != nil {
			roomName = mob.InRoom.Name
			roomVnum = mob.InRoom.VNUM
		}

		protoVnum := 0
		if mob.Prototype != nil {
			protoVnum = mob.Prototype.VNUM
		}

		sb.WriteString(fmt.Sprintf("- %s (VNUM: %d) in room %d (%s)\r\n", 
			mob.ShortDesc, protoVnum, roomVnum, roomName))
	}

	return fmt.Errorf("%s", sb.String())
}
