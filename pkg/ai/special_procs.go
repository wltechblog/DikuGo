package ai

import (
	"log"
	"strings"

	"github.com/wltechblog/DikuGo/pkg/types"
)

// SpecialProcs is a map of special procedure names to functions
var SpecialProcs = map[string]func(*types.Character, string) bool{
	"cityguard":     cityguardProc,
	"shopkeeper":    shopkeeperProc,
	"thief":         thiefProc,
	"magic_user":    magicUserProc,
	"snake":         snakeProc,
	"fido":          fidoProc,
	"janitor":       janitorProc,
	"guild_guard":   guildGuardProc,
	"receptionist":  receptionistProc,
	"mayor":         mayorProc,
}

// cityguardProc is the special procedure for cityguards
func cityguardProc(ch *types.Character, argument string) bool {
	// Skip if the character is not an NPC
	if !ch.IsNPC {
		return false
	}

	// Skip if the character is fighting
	if ch.Fighting != nil {
		return false
	}

	// Use the cityguard behavior
	return Behaviors["cityguard"](ch, nil)
}

// shopkeeperProc is the special procedure for shopkeepers
func shopkeeperProc(ch *types.Character, argument string) bool {
	// Skip if the character is not an NPC
	if !ch.IsNPC {
		return false
	}

	// Skip if the character is fighting
	if ch.Fighting != nil {
		return false
	}

	// Handle shop commands
	if strings.HasPrefix(argument, "list") {
		// TODO: Implement list command
		return true
	} else if strings.HasPrefix(argument, "buy") {
		// TODO: Implement buy command
		return true
	} else if strings.HasPrefix(argument, "sell") {
		// TODO: Implement sell command
		return true
	}

	// Use the shopkeeper behavior
	return Behaviors["shopkeeper"](ch, nil)
}

// thiefProc is the special procedure for thieves
func thiefProc(ch *types.Character, argument string) bool {
	// Skip if the character is not an NPC
	if !ch.IsNPC {
		return false
	}

	// Skip if the character is fighting
	if ch.Fighting != nil {
		return false
	}

	// Use the thief behavior
	return Behaviors["thief"](ch, nil)
}

// magicUserProc is the special procedure for magic users
func magicUserProc(ch *types.Character, argument string) bool {
	// Skip if the character is not an NPC
	if !ch.IsNPC {
		return false
	}

	// Use the magic user behavior
	return Behaviors["magic_user"](ch, nil)
}

// snakeProc is the special procedure for snakes
func snakeProc(ch *types.Character, argument string) bool {
	// Skip if the character is not an NPC
	if !ch.IsNPC {
		return false
	}

	// Use the snake behavior
	return Behaviors["snake"](ch, nil)
}

// fidoProc is the special procedure for fidos (corpse eaters)
func fidoProc(ch *types.Character, argument string) bool {
	// Skip if the character is not an NPC
	if !ch.IsNPC {
		return false
	}

	// Use the fido behavior
	return Behaviors["fido"](ch, nil)
}

// janitorProc is the special procedure for janitors
func janitorProc(ch *types.Character, argument string) bool {
	// Skip if the character is not an NPC
	if !ch.IsNPC {
		return false
	}

	// Use the janitor behavior
	return Behaviors["janitor"](ch, nil)
}

// guildGuardProc is the special procedure for guild guards
func guildGuardProc(ch *types.Character, argument string) bool {
	// Skip if the character is not an NPC
	if !ch.IsNPC {
		return false
	}

	// Use the guild guard behavior
	return Behaviors["guild_guard"](ch, nil)
}

// receptionistProc is the special procedure for receptionists
func receptionistProc(ch *types.Character, argument string) bool {
	// Skip if the character is not an NPC
	if !ch.IsNPC {
		return false
	}

	// Handle rent commands
	if strings.HasPrefix(argument, "rent") {
		// TODO: Implement rent command
		return true
	} else if strings.HasPrefix(argument, "offer") {
		// TODO: Implement offer command
		return true
	}

	// Use the receptionist behavior
	return Behaviors["receptionist"](ch, nil)
}

// mayorProc is the special procedure for mayors
func mayorProc(ch *types.Character, argument string) bool {
	// Skip if the character is not an NPC
	if !ch.IsNPC {
		return false
	}

	// Use the mayor behavior
	return Behaviors["mayor"](ch, nil)
}

// RegisterSpecialProcs registers special procedures with mobile prototypes
func RegisterSpecialProcs(mobiles []*types.Mobile) {
	for _, mobile := range mobiles {
		// Check for special procedures based on keywords
		for name, proc := range SpecialProcs {
			if strings.Contains(strings.ToLower(mobile.Name), name) {
				log.Printf("Registering special procedure %s for mobile %s", name, mobile.Name)
				mobile.Functions = append(mobile.Functions, proc)
			}
		}
	}
}
