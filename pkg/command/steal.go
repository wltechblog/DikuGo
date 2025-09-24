package command

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/wltechblog/DikuGo/pkg/types"
)

// StealCommand represents the steal command
type StealCommand struct {
	// CombatManager is the combat manager
	CombatManager CombatManagerInterface
}

// Execute executes the steal command
func (c *StealCommand) Execute(character *types.Character, args string) error {
	// Check if the character is in a room
	if character.InRoom == nil {
		return fmt.Errorf("you are not in a room")
	}

	// Get the world interface
	world, ok := character.World.(interface {
		HasSkill(*types.Character, int) bool
		CanUseSkill(*types.Character, int) (bool, string)
		CheckSkillSuccess(*types.Character, int) bool
		UseSkill(*types.Character, int)
		ImproveSkill(*types.Character, int, bool)
		AddObjectToCharacter(*types.Character, *types.ObjectInstance) error
		RemoveObjectFromCharacter(*types.Character, *types.ObjectInstance) error
	})
	if !ok {
		return fmt.Errorf("world interface not available")
	}

	// Check if the character can use the steal skill
	if can, reason := world.CanUseSkill(character, types.SKILL_STEAL); !can {
		return fmt.Errorf("%s", reason)
	}

	// Parse the arguments
	args = strings.TrimSpace(args)
	if args == "" {
		return fmt.Errorf("steal what from whom?")
	}

	// Split the arguments
	parts := strings.Split(args, " from ")
	if len(parts) != 2 {
		return fmt.Errorf("steal what from whom?")
	}

	itemName := strings.TrimSpace(parts[0])
	victimName := strings.TrimSpace(parts[1])

	// Find the victim in the room
	victim := findCharacterInRoom(character.InRoom, victimName)
	if victim == nil {
		return fmt.Errorf("they aren't here")
	}

	// Check if the victim is the character
	if victim == character {
		return fmt.Errorf("you can't steal from yourself")
	}

	// Check if the victim is a player
	if !victim.IsNPC {
		// Stealing from players is more difficult
		if rand.Intn(100) < 50 {
			// Failed attempt
			character.SendMessage(fmt.Sprintf("You fail to steal from %s.\r\n", victim.ShortDesc))

			// Alert the victim
			victim.SendMessage(fmt.Sprintf("%s tried to steal from you!\r\n", character.Name))

			// Alert the room
			for _, ch := range character.InRoom.Characters {
				if ch != character && ch != victim {
					ch.SendMessage(fmt.Sprintf("%s tried to steal from %s!\r\n", character.Name, victim.ShortDesc))
				}
			}

			// Start combat
			c.CombatManager.StartCombat(victim, character)

			return nil
		}
	}

	// Mark the skill as used (set cooldown)
	world.UseSkill(character, types.SKILL_STEAL)

	// Check if the steal is successful
	success := world.CheckSkillSuccess(character, types.SKILL_STEAL)

	// Give a chance to improve the skill
	world.ImproveSkill(character, types.SKILL_STEAL, success)

	if !success {
		// Steal failed
		character.SendMessage(fmt.Sprintf("You fail to steal from %s.\r\n", victim.ShortDesc))

		// There's a chance the victim notices
		if rand.Intn(100) < 25 {
			// Alert the victim
			victim.SendMessage(fmt.Sprintf("%s tried to steal from you!\r\n", character.Name))

			// Alert the room
			for _, ch := range character.InRoom.Characters {
				if ch != character && ch != victim {
					ch.SendMessage(fmt.Sprintf("%s tried to steal from %s!\r\n", character.Name, victim.ShortDesc))
				}
			}

			// If the victim is an NPC, they might attack
			if victim.IsNPC {
				c.CombatManager.StartCombat(victim, character)
			}
		}

		return nil
	}

	// Steal succeeded
	// Check if we're stealing gold
	if strings.ToLower(itemName) == "gold" || strings.ToLower(itemName) == "coins" {
		// Steal gold
		// Calculate how much gold to steal (10-25% of victim's gold)
		goldPercent := rand.Intn(16) + 10 // 10-25%
		goldAmount := (victim.Gold * goldPercent) / 100

		// Ensure minimum amount
		if goldAmount < 1 {
			goldAmount = 1
		}

		// Ensure we don't steal more than the victim has
		if goldAmount > victim.Gold {
			goldAmount = victim.Gold
		}

		// Transfer the gold
		victim.Gold -= goldAmount
		character.Gold += goldAmount

		// Send messages
		character.SendMessage(fmt.Sprintf("You steal %d gold coins from %s.\r\n", goldAmount, victim.ShortDesc))

		return nil
	}

	// Steal an item
	// Find the item in the victim's inventory
	var item *types.ObjectInstance
	for _, obj := range victim.Inventory {
		if strings.Contains(strings.ToLower(obj.Prototype.Name), strings.ToLower(itemName)) {
			item = obj
			break
		}
	}

	// Check if the item was found
	if item == nil {
		// Check equipment if not found in inventory
		for _, obj := range victim.Equipment {
			if obj != nil && strings.Contains(strings.ToLower(obj.Prototype.Name), strings.ToLower(itemName)) {
				// Can't steal worn equipment
				return fmt.Errorf("you can't steal worn equipment")
			}
		}

		return fmt.Errorf("you can't find it")
	}

	// Remove the item from the victim
	err := world.RemoveObjectFromCharacter(victim, item)
	if err != nil {
		return fmt.Errorf("you can't steal that")
	}

	// Add the item to the character
	err = world.AddObjectToCharacter(character, item)
	if err != nil {
		// If we can't add it to the character, put it back in the victim's inventory
		world.AddObjectToCharacter(victim, item)
		return fmt.Errorf("you can't carry any more")
	}

	// Send messages
	character.SendMessage(fmt.Sprintf("You steal %s from %s.\r\n", item.Prototype.ShortDesc, victim.ShortDesc))

	return nil
}

// Name returns the name of the command
func (c *StealCommand) Name() string {
	return "steal"
}

// Aliases returns the aliases of the command
func (c *StealCommand) Aliases() []string {
	return []string{}
}

// MinPosition returns the minimum position required to execute the command
func (c *StealCommand) MinPosition() int {
	return types.POS_STANDING
}

// Level returns the minimum level required to execute the command
func (c *StealCommand) Level() int {
	return 0
}

// LogCommand returns whether the command should be logged
func (c *StealCommand) LogCommand() bool {
	return false
}
