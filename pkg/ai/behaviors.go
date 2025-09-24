package ai

import (
	"fmt"
	"log"
	"math/rand"

	"github.com/wltechblog/DikuGo/pkg/types"
)

// BehaviorFunc is a function that implements a specific behavior
type BehaviorFunc func(mobile *types.Character, world interface{}) bool

// Behaviors is a map of behavior names to behavior functions
var Behaviors = map[string]BehaviorFunc{
	"cityguard":    cityguardBehavior,
	"shopkeeper":   shopkeeperBehavior,
	"thief":        thiefBehavior,
	"magic_user":   magicUserBehavior,
	"snake":        snakeBehavior,
	"fido":         fidoBehavior,
	"janitor":      janitorBehavior,
	"guild_guard":  guildGuardBehavior,
	"receptionist": receptionistBehavior,
	"mayor":        mayorBehavior,
}

// cityguardBehavior implements the behavior for cityguards
func cityguardBehavior(mobile *types.Character, world interface{}) bool {
	// Skip if the mobile is not in a room
	if mobile.InRoom == nil {
		return false
	}

	// Skip if the mobile is not awake or is already fighting
	if mobile.Position < types.POS_RESTING || mobile.Fighting != nil {
		return false
	}

	// Look for fighting in the room
	var target *types.Character
	worstAlignment := 1000 // Start with a high alignment threshold

	for _, character := range mobile.InRoom.Characters {
		// Skip if not fighting
		if character.Fighting == nil {
			continue
		}

		// Check if this character is more evil than our current target
		if character.Alignment < worstAlignment && (character.IsNPC || character.Fighting.IsNPC) {
			worstAlignment = character.Alignment
			target = character
		}
	}

	// If we found a target and they're fighting someone good
	if target != nil && target.Fighting != nil && target.Fighting.Alignment >= 0 {
		// Cityguard found a criminal, now attack them
		log.Printf("%s screams 'PROTECT THE INNOCENT! BANZAI!!! CHARGE!!! ARARARAGGGHH!'", mobile.Name)

		// Broadcast the action to the room
		for _, ch := range mobile.InRoom.Characters {
			ch.SendMessage(fmt.Sprintf("%s screams 'PROTECT THE INNOCENT! BANZAI!!! CHARGE!!! ARARARAGGGHH!'\r\n", mobile.Name))
		}

		// Start fighting the target
		mobile.Fighting = target
		target.Fighting = mobile
		mobile.Position = types.POS_FIGHTING

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

	// Apply poison effect to the target
	target := mobile.Fighting

	// Create a poison affect
	poisonAffect := &types.Affect{
		Type:      types.AFF_POISON,
		Duration:  rand.Intn(3) + 2, // 2-4 rounds
		Modifier:  -2,
		Location:  types.APPLY_STR, // Reduce strength
		Bitvector: types.AFF_POISON,
	}

	// Add the poison affect to the target
	if world, ok := mobile.World.(interface {
		AffectToChar(*types.Character, *types.Affect)
	}); ok {
		world.AffectToChar(target, poisonAffect)
	}

	// Broadcast the action to the room
	for _, ch := range mobile.InRoom.Characters {
		if ch == target {
			ch.SendMessage(fmt.Sprintf("You feel poison coursing through your veins!\r\n"))
		} else {
			ch.SendMessage(fmt.Sprintf("%s has been poisoned by %s!\r\n", target.Name, mobile.Name))
		}
	}

	log.Printf("Snake %s poisoned %s", mobile.Name, target.Name)
	return true
}

// fidoBehavior implements the behavior for fidos (corpse eaters)
func fidoBehavior(mobile *types.Character, world interface{}) bool {
	// Skip if the mobile is not in a room
	if mobile.InRoom == nil {
		return false
	}

	// Skip if the mobile is not awake
	if mobile.Position < types.POS_RESTING {
		return false
	}

	// Look for corpses in the room
	var corpse *types.ObjectInstance
	for _, obj := range mobile.InRoom.Objects {
		// Check if the object is a container (corpse)
		if obj.Prototype != nil && obj.Prototype.Type == types.ITEM_CONTAINER {
			// Check if the container is a corpse (value[3] is set for corpses)
			if obj.Prototype.Value[3] > 0 {
				corpse = obj
				break
			}
		}
	}

	// If no corpse found, return
	if corpse == nil {
		return false
	}

	// Fido found a corpse, now eat it
	log.Printf("%s savagely devours a corpse.", mobile.Name)

	// Move the contents of the corpse to the room
	for _, item := range corpse.Contains {
		// Remove the item from the corpse
		item.InObj = nil

		// Add the item to the room
		item.InRoom = mobile.InRoom
		mobile.InRoom.Objects = append(mobile.InRoom.Objects, item)
	}

	// Clear the corpse's contents
	corpse.Contains = nil

	// Remove the corpse from the room
	for i, obj := range mobile.InRoom.Objects {
		if obj == corpse {
			// Remove the corpse from the room's objects
			mobile.InRoom.Objects = append(mobile.InRoom.Objects[:i], mobile.InRoom.Objects[i+1:]...)
			break
		}
	}

	// Broadcast the action to the room
	for _, ch := range mobile.InRoom.Characters {
		if ch != mobile {
			ch.SendMessage(fmt.Sprintf("%s savagely devours a corpse.\r\n", mobile.Name))
		}
	}

	return true
}

// janitorBehavior implements the behavior for janitors
func janitorBehavior(mobile *types.Character, world interface{}) bool {
	// Skip if the mobile is not in a room
	if mobile.InRoom == nil {
		return false
	}

	// Skip if the mobile is not awake
	if mobile.Position < types.POS_RESTING {
		return false
	}

	// Look for trash in the room
	var trash *types.ObjectInstance
	for _, obj := range mobile.InRoom.Objects {
		// Check if the object is trash
		if obj.Prototype != nil && (obj.Prototype.Type == types.ITEM_TRASH ||
			(obj.Prototype.Cost <= 10 && obj.Prototype.Weight <= 10)) {
			trash = obj
			break
		}
	}

	// If no trash found, return
	if trash == nil {
		return false
	}

	// Janitor found trash, now pick it up
	log.Printf("%s picks up some trash.", mobile.Name)

	// Remove the trash from the room
	for i, obj := range mobile.InRoom.Objects {
		if obj == trash {
			// Remove the trash from the room's objects
			mobile.InRoom.Objects = append(mobile.InRoom.Objects[:i], mobile.InRoom.Objects[i+1:]...)
			break
		}
	}

	// Add the trash to the janitor's inventory
	trash.InRoom = nil
	trash.CarriedBy = mobile
	mobile.Inventory = append(mobile.Inventory, trash)

	// Broadcast the action to the room
	for _, ch := range mobile.InRoom.Characters {
		if ch != mobile {
			ch.SendMessage(fmt.Sprintf("%s picks up some trash.\r\n", mobile.Name))
		}
	}

	return true
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
