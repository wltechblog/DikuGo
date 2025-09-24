package command

import (
	"fmt"
	"math/rand"

	"github.com/wltechblog/DikuGo/pkg/types"
)

// BackstabCommand represents the backstab command
type BackstabCommand struct {
	// CombatManager is the combat manager
	CombatManager CombatManagerInterface
}

// Execute executes the backstab command
func (c *BackstabCommand) Execute(character *types.Character, args string) error {
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

	// Check if the character can use the backstab skill
	if can, reason := world.CanUseSkill(character, types.SKILL_BACKSTAB); !can {
		return fmt.Errorf("%s", reason)
	}

	// Check if the character has a weapon
	weapon := character.Equipment[types.WEAR_WIELD]
	if weapon == nil {
		return fmt.Errorf("you need to wield a weapon to backstab someone")
	}

	// Check if the weapon is a piercing weapon
	if weapon.Prototype.Value[3] != types.TYPE_PIERCE {
		return fmt.Errorf("you need to wield a piercing weapon to backstab someone")
	}

	// Check if there's a target specified
	if args == "" {
		return fmt.Errorf("backstab whom?")
	}

	// Find the target in the room
	victim := findCharacterInRoom(character.InRoom, args)
	if victim == nil {
		return fmt.Errorf("they aren't here")
	}

	// Check if the target is the character
	if victim == character {
		return fmt.Errorf("how can you sneak up on yourself?")
	}

	// Check if the character is already fighting
	if character.Fighting != nil {
		return fmt.Errorf("you're too busy fighting to backstab")
	}

	// Check if the victim is already fighting
	if victim.Fighting != nil {
		return fmt.Errorf("you can't backstab someone who is fighting")
	}

	// Mark the skill as used (set cooldown)
	world.UseSkill(character, types.SKILL_BACKSTAB)

	// Check if the backstab is successful
	success := world.CheckSkillSuccess(character, types.SKILL_BACKSTAB)

	// Give a chance to improve the skill
	world.ImproveSkill(character, types.SKILL_BACKSTAB, success)

	if !success {
		// Backstab failed
		// Send messages
		character.SendMessage(fmt.Sprintf("You try to backstab %s, but they see you coming!\r\n", victim.ShortDesc))
		victim.SendMessage(fmt.Sprintf("%s tries to backstab you, but you see them coming!\r\n", character.Name))

		// Send a message to the room
		for _, ch := range character.InRoom.Characters {
			if ch != character && ch != victim {
				ch.SendMessage(fmt.Sprintf("%s tries to backstab %s, but is seen!\r\n", character.Name, victim.ShortDesc))
			}
		}

		// Start combat normally
		c.CombatManager.StartCombat(character, victim)

		return nil
	}

	// Backstab succeeded
	// Calculate damage multiplier based on level
	multiplier := 2
	if character.Level >= 10 {
		multiplier = 3
	}
	if character.Level >= 20 {
		multiplier = 4
	}

	// Calculate base damage
	var baseDamage int
	// Use the weapon's damage dice
	numDice := weapon.Prototype.Value[1]
	sizeDice := weapon.Prototype.Value[2]
	if numDice > 0 && sizeDice > 0 {
		// Roll the dice
		for i := 0; i < numDice; i++ {
			baseDamage += rand.Intn(sizeDice) + 1
		}
	}

	// Apply strength bonus
	strBonus := 0
	switch {
	case character.Abilities[0] <= 3:
		strBonus = -3
	case character.Abilities[0] <= 5:
		strBonus = -2
	case character.Abilities[0] <= 7:
		strBonus = -1
	case character.Abilities[0] <= 13:
		strBonus = 0
	case character.Abilities[0] <= 15:
		strBonus = 1
	case character.Abilities[0] <= 17:
		strBonus = 2
	case character.Abilities[0] <= 18:
		strBonus = 3
	default:
		strBonus = 4
	}

	// Calculate total damage
	damage := (baseDamage + strBonus + character.DamRoll) * multiplier

	// Ensure minimum damage
	if damage < 1 {
		damage = 1
	}

	// Apply damage
	victim.HP -= damage
	if victim.HP < 0 {
		victim.HP = 0
	}

	// Send messages
	character.SendMessage(fmt.Sprintf("You backstab %s for %d damage!\r\n", victim.ShortDesc, damage))
	victim.SendMessage(fmt.Sprintf("%s backstabs you for %d damage!\r\n", character.Name, damage))

	// Send a message to the room
	for _, ch := range character.InRoom.Characters {
		if ch != character && ch != victim {
			ch.SendMessage(fmt.Sprintf("%s backstabs %s with deadly precision!\r\n", character.Name, victim.ShortDesc))
		}
	}

	// Start combat
	c.CombatManager.StartCombat(character, victim)

	// Check if the victim died
	if victim.HP <= 0 {
		victim.Position = types.POS_DEAD

		// Send death messages
		character.SendMessage(fmt.Sprintf("You have slain %s!\r\n", victim.ShortDesc))
		victim.SendMessage(fmt.Sprintf("%s has slain you!\r\n", character.Name))

		// Send a message to the room
		for _, ch := range character.InRoom.Characters {
			if ch != character && ch != victim {
				ch.SendMessage(fmt.Sprintf("%s has slain %s!\r\n", character.Name, victim.ShortDesc))
			}
		}

		// Create a corpse
		w, ok := victim.World.(interface {
			MakeCorpse(*types.Character) *types.ObjectInstance
			RemoveCharacter(*types.Character)
			ScheduleMobRespawn(*types.Character)
		})
		if ok {
			// Create the corpse first
			w.MakeCorpse(victim)

			// If the defender is an NPC, schedule respawn
			if victim.IsNPC && victim.Prototype != nil {
				w.ScheduleMobRespawn(victim)
			} else {
				// For players or if scheduling fails, just remove from room
				w.RemoveCharacter(victim)
			}
		}

		// Stop combat
		c.CombatManager.StopCombat(character)
	}

	return nil
}

// Name returns the name of the command
func (c *BackstabCommand) Name() string {
	return "backstab"
}

// Aliases returns the aliases of the command
func (c *BackstabCommand) Aliases() []string {
	return []string{"bs"}
}

// MinPosition returns the minimum position required to execute the command
func (c *BackstabCommand) MinPosition() int {
	return types.POS_STANDING
}

// Level returns the minimum level required to execute the command
func (c *BackstabCommand) Level() int {
	return 0
}

// LogCommand returns whether the command should be logged
func (c *BackstabCommand) LogCommand() bool {
	return false
}
