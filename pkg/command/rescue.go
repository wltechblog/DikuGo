package command

import (
	"fmt"

	"github.com/wltechblog/DikuGo/pkg/types"
)

// RescueCommand represents the rescue command
type RescueCommand struct {
	// CombatManager is the combat manager
	CombatManager CombatManagerInterface
}

// Execute executes the rescue command
func (c *RescueCommand) Execute(character *types.Character, args string) error {
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

	// Check if the character can use the rescue skill
	if can, reason := world.CanUseSkill(character, types.SKILL_RESCUE); !can {
		return fmt.Errorf("%s", reason)
	}

	// Check if there's a target specified
	if args == "" {
		return fmt.Errorf("rescue whom?")
	}

	// Find the target in the room
	victim := findCharacterInRoom(character.InRoom, args)
	if victim == nil {
		return fmt.Errorf("they aren't here")
	}

	// Check if the target is the character
	if victim == character {
		return fmt.Errorf("what about fleeing instead?")
	}

	// Check if the character is fighting the victim
	if character.Fighting == victim {
		return fmt.Errorf("how can you rescue someone you are trying to kill?")
	}

	// Find someone who is fighting the victim
	var attacker *types.Character
	for _, ch := range character.InRoom.Characters {
		if ch.Fighting == victim {
			attacker = ch
			break
		}
	}

	// Check if anyone is fighting the victim
	if attacker == nil {
		return fmt.Errorf("but nobody is fighting %s!", victim.ShortDesc)
	}

	// Mark the skill as used (set cooldown)
	world.UseSkill(character, types.SKILL_RESCUE)

	// Check if the rescue is successful
	success := world.CheckSkillSuccess(character, types.SKILL_RESCUE)

	// Give a chance to improve the skill
	world.ImproveSkill(character, types.SKILL_RESCUE, success)

	if !success {
		// Rescue failed
		character.SendMessage(fmt.Sprintf("You fail to rescue %s!\r\n", victim.ShortDesc))
		return nil
	}

	// Rescue succeeded
	// Send messages
	character.SendMessage(fmt.Sprintf("Banzai! To the rescue of %s!\r\n", victim.ShortDesc))
	victim.SendMessage(fmt.Sprintf("%s heroically rescues you!\r\n", character.Name))

	// Send a message to the room
	for _, ch := range character.InRoom.Characters {
		if ch != character && ch != victim {
			ch.SendMessage(fmt.Sprintf("%s heroically rescues %s!\r\n", character.Name, victim.ShortDesc))
		}
	}

	// Stop the current fights
	if victim.Fighting == attacker {
		c.CombatManager.StopCombat(victim)
	}

	if attacker.Fighting == victim {
		c.CombatManager.StopCombat(attacker)
	}

	if character.Fighting != nil {
		c.CombatManager.StopCombat(character)
	}

	// Start a new fight between the character and the attacker
	c.CombatManager.StartCombat(character, attacker)

	// Wait state for victim (2 combat rounds)
	// This is handled by the skill cooldown system

	return nil
}

// Name returns the name of the command
func (c *RescueCommand) Name() string {
	return "rescue"
}

// Aliases returns the aliases of the command
func (c *RescueCommand) Aliases() []string {
	return []string{}
}

// MinPosition returns the minimum position required to execute the command
func (c *RescueCommand) MinPosition() int {
	return types.POS_FIGHTING
}

// Level returns the minimum level required to execute the command
func (c *RescueCommand) Level() int {
	return 0
}

// LogCommand returns whether the command should be logged
func (c *RescueCommand) LogCommand() bool {
	return false
}
