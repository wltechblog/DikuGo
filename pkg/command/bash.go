package command

import (
	"fmt"

	"github.com/wltechblog/DikuGo/pkg/types"
)

// BashCommand represents the bash command
type BashCommand struct {
	// CombatManager is the combat manager
	CombatManager CombatManagerInterface
}

// Execute executes the bash command
func (c *BashCommand) Execute(character *types.Character, args string) error {
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
	})
	if !ok {
		return fmt.Errorf("world interface not available")
	}

	// Check if the character can use the bash skill
	if can, reason := world.CanUseSkill(character, types.SKILL_BASH); !can {
		return fmt.Errorf("%s", reason)
	}

	// Check if the character has a weapon
	if character.Equipment[types.WEAR_WIELD] == nil {
		return fmt.Errorf("you need to wield a weapon to bash someone")
	}

	// Find the target
	var victim *types.Character
	if args == "" {
		// If no target specified, use current opponent
		if character.Fighting == nil {
			return fmt.Errorf("bash whom?")
		}
		victim = character.Fighting
	} else {
		// Find the target in the room
		victim = findCharacterInRoom(character.InRoom, args)
		if victim == nil {
			return fmt.Errorf("they aren't here")
		}
	}

	// Check if the target is the character
	if victim == character {
		return fmt.Errorf("aren't we funny today...")
	}

	// Mark the skill as used (set cooldown)
	world.UseSkill(character, types.SKILL_BASH)

	// Check if the bash is successful
	success := world.CheckSkillSuccess(character, types.SKILL_BASH)

	// Give a chance to improve the skill
	world.ImproveSkill(character, types.SKILL_BASH, success)

	if !success {
		// Bash failed
		// Send messages
		character.SendMessage(fmt.Sprintf("You try to bash %s, but fall flat on your face!\r\n", victim.ShortDesc))
		victim.SendMessage(fmt.Sprintf("%s tries to bash you, but falls flat on their face!\r\n", character.Name))

		// Send a message to the room
		for _, ch := range character.InRoom.Characters {
			if ch != character && ch != victim {
				ch.SendMessage(fmt.Sprintf("%s tries to bash %s, but falls flat on their face!\r\n", character.Name, victim.ShortDesc))
			}
		}

		// Character falls down
		character.Position = types.POS_SITTING

		// Wait state (2 combat rounds)
		// This is handled by the skill cooldown system

		return nil
	}

	// Bash succeeded
	// Send messages
	character.SendMessage(fmt.Sprintf("You bash %s, sending them sprawling!\r\n", victim.ShortDesc))
	victim.SendMessage(fmt.Sprintf("%s bashes you, sending you sprawling!\r\n", character.Name))

	// Send a message to the room
	for _, ch := range character.InRoom.Characters {
		if ch != character && ch != victim {
			ch.SendMessage(fmt.Sprintf("%s bashes %s, sending them sprawling!\r\n", character.Name, victim.ShortDesc))
		}
	}

	// Deal minor damage (1 point)
	if victim.HP > 0 {
		victim.HP -= 1
		if victim.HP < 0 {
			victim.HP = 0
		}
	}

	// Victim falls down
	victim.Position = types.POS_SITTING

	// Start combat if not already fighting
	if character.Fighting == nil || character.Fighting != victim {
		c.CombatManager.StartCombat(character, victim)
	}

	// Wait state for victim (2 combat rounds)
	// This is handled by setting the position to sitting

	return nil
}

// Name returns the name of the command
func (c *BashCommand) Name() string {
	return "bash"
}

// Aliases returns the aliases of the command
func (c *BashCommand) Aliases() []string {
	return []string{}
}

// MinPosition returns the minimum position required to execute the command
func (c *BashCommand) MinPosition() int {
	return types.POS_FIGHTING
}

// Level returns the minimum level required to execute the command
func (c *BashCommand) Level() int {
	return 0
}

// LogCommand returns whether the command should be logged
func (c *BashCommand) LogCommand() bool {
	return false
}
