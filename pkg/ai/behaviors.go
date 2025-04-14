package ai

import (
	"log"
	"math/rand"

	"github.com/wltechblog/DikuGo/pkg/types"
)

// BehaviorFunc is a function that implements a specific behavior
type BehaviorFunc func(mobile *types.Character, world interface{}) bool

// Behaviors is a map of behavior names to behavior functions
var Behaviors = map[string]BehaviorFunc{
	"cityguard":     cityguardBehavior,
	"shopkeeper":    shopkeeperBehavior,
	"thief":         thiefBehavior,
	"magic_user":    magicUserBehavior,
	"snake":         snakeBehavior,
	"fido":          fidoBehavior,
	"janitor":       janitorBehavior,
	"guild_guard":   guildGuardBehavior,
	"receptionist":  receptionistBehavior,
	"mayor":         mayorBehavior,
}

// cityguardBehavior implements the behavior for cityguards
func cityguardBehavior(mobile *types.Character, world interface{}) bool {
	// Skip if the mobile is not in a room
	if mobile.InRoom == nil {
		return false
	}

	// Look for criminals in the room
	for _, character := range mobile.InRoom.Characters {
		// Skip NPCs and non-criminals
		if character.IsNPC || character.Alignment >= 0 {
			continue
		}

		// TODO: Implement attack command for mobiles
		log.Printf("Cityguard %s would attack criminal %s", mobile.Name, character.Name)
		return true
	}

	return false
}

// shopkeeperBehavior implements the behavior for shopkeepers
func shopkeeperBehavior(mobile *types.Character, world interface{}) bool {
	// Skip if the mobile is not in a room
	if mobile.InRoom == nil {
		return false
	}

	// Skip if the room is not a shop
	if mobile.InRoom.Shop == nil {
		return false
	}

	// Shopkeepers don't do much except stay in their shop
	return true
}

// thiefBehavior implements the behavior for thieves
func thiefBehavior(mobile *types.Character, world interface{}) bool {
	// Skip if the mobile is not in a room
	if mobile.InRoom == nil {
		return false
	}

	// Look for victims in the room
	for _, character := range mobile.InRoom.Characters {
		// Skip NPCs and high-level players
		if character.IsNPC || character.Level > 20 {
			continue
		}

		// Only steal occasionally
		if rand.Intn(100) > 10 {
			continue
		}

		// TODO: Implement steal command for mobiles
		log.Printf("Thief %s would steal from %s", mobile.Name, character.Name)
		return true
	}

	return false
}

// magicUserBehavior implements the behavior for magic users
func magicUserBehavior(mobile *types.Character, world interface{}) bool {
	// Skip if the mobile is not in a room
	if mobile.InRoom == nil {
		return false
	}

	// Skip if the mobile is not fighting
	if mobile.Fighting == nil {
		return false
	}

	// Cast a random spell
	spells := []string{"magic missile", "chill touch", "burning hands", "shocking grasp", "lightning bolt", "colour spray"}
	spell := spells[rand.Intn(len(spells))]

	// TODO: Implement cast command for mobiles
	log.Printf("Magic user %s would cast %s at %s", mobile.Name, spell, mobile.Fighting.Name)
	return true
}

// snakeBehavior implements the behavior for snakes
func snakeBehavior(mobile *types.Character, world interface{}) bool {
	// Skip if the mobile is not in a room
	if mobile.InRoom == nil {
		return false
	}

	// Skip if the mobile is not fighting
	if mobile.Fighting == nil {
		return false
	}

	// Snakes have a chance to poison their target
	if rand.Intn(100) > 20 {
		return false
	}

	// TODO: Implement poison effect for mobiles
	log.Printf("Snake %s would poison %s", mobile.Name, mobile.Fighting.Name)
	return true
}

// fidoBehavior implements the behavior for fidos (corpse eaters)
func fidoBehavior(mobile *types.Character, world interface{}) bool {
	// Skip if the mobile is not in a room
	if mobile.InRoom == nil {
		return false
	}

	// Look for corpses in the room
	// TODO: Implement corpse finding and eating
	return false
}

// janitorBehavior implements the behavior for janitors
func janitorBehavior(mobile *types.Character, world interface{}) bool {
	// Skip if the mobile is not in a room
	if mobile.InRoom == nil {
		return false
	}

	// Look for trash in the room
	// TODO: Implement trash finding and cleaning
	return false
}

// guildGuardBehavior implements the behavior for guild guards
func guildGuardBehavior(mobile *types.Character, world interface{}) bool {
	// Skip if the mobile is not in a room
	if mobile.InRoom == nil {
		return false
	}

	// Look for characters of the wrong class
	for _, character := range mobile.InRoom.Characters {
		// Skip NPCs
		if character.IsNPC {
			continue
		}

		// TODO: Implement guild guard behavior
		return false
	}

	return false
}

// receptionistBehavior implements the behavior for receptionists
func receptionistBehavior(mobile *types.Character, world interface{}) bool {
	// Skip if the mobile is not in a room
	if mobile.InRoom == nil {
		return false
	}

	// Receptionists don't do much except stay in their inn
	return true
}

// mayorBehavior implements the behavior for mayors
func mayorBehavior(mobile *types.Character, world interface{}) bool {
	// Skip if the mobile is not in a room
	if mobile.InRoom == nil {
		return false
	}

	// TODO: Implement mayor behavior (walking around town, etc.)
	return false
}
