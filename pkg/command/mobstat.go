package command

import (
	"fmt"
	"strings"

	"github.com/wltechblog/DikuGo/pkg/types"
)

// MobstatCommand is a command that shows detailed information about a mob
type MobstatCommand struct{}

// Name returns the name of the command
func (c *MobstatCommand) Name() string {
	return "mobstat"
}

// Aliases returns the aliases of the command
func (c *MobstatCommand) Aliases() []string {
	return []string{"mstat"}
}

// MinPosition returns the minimum position required to execute the command
func (c *MobstatCommand) MinPosition() int {
	return types.POS_STANDING
}

// Level returns the minimum level required to execute the command
func (c *MobstatCommand) Level() int {
	return 1 // Available to everyone during development
}

// LogCommand returns whether the command should be logged
func (c *MobstatCommand) LogCommand() bool {
	return true
}

// Execute executes the mobstat command
func (c *MobstatCommand) Execute(character *types.Character, args string) error {
	// Check if the character is in a room
	if character.InRoom == nil {
		return fmt.Errorf("you are not in a room")
	}

	// Get the room
	room := character.InRoom

	// If no arguments, list all mobs in the room
	if args == "" {
		var sb strings.Builder
		sb.WriteString("\r\nMobs in this room:\r\n")

		mobCount := 0
		for _, ch := range room.Characters {
			if ch.IsNPC {
				mobCount++
				sb.WriteString(fmt.Sprintf("%d. %s\r\n", mobCount, ch.ShortDesc))
			}
		}

		if mobCount == 0 {
			sb.WriteString("No mobs in this room.\r\n")
		} else {
			sb.WriteString("\r\nUse 'mobstat <name>' to see details about a specific mob.\r\n")
		}

		character.SendMessage(sb.String())
		return nil
	}

	// Find the mob by name
	var targetMob *types.Character
	for _, ch := range room.Characters {
		if ch.IsNPC && strings.Contains(strings.ToLower(ch.Name), strings.ToLower(args)) {
			targetMob = ch
			break
		}
	}

	if targetMob == nil {
		return fmt.Errorf("no mob named '%s' found in this room", args)
	}

	// Build the mob stats
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("\r\n--- Mob Stats for %s ---\r\n", targetMob.Name))

	// Basic information
	sb.WriteString(fmt.Sprintf("VNUM: %d\r\n", targetMob.Prototype.VNUM))
	sb.WriteString(fmt.Sprintf("Name: %s\r\n", targetMob.Name))
	sb.WriteString(fmt.Sprintf("Short Description: %s\r\n", targetMob.ShortDesc))
	sb.WriteString(fmt.Sprintf("Long Description: %s\r\n", targetMob.LongDesc))
	sb.WriteString(fmt.Sprintf("Description: %s\r\n", targetMob.Description))

	// Stats
	sb.WriteString(fmt.Sprintf("\r\nLevel: %d\r\n", targetMob.Level))
	sb.WriteString(fmt.Sprintf("Position: %d (%s)\r\n", targetMob.Position, positionToString(targetMob.Position)))
	sb.WriteString(fmt.Sprintf("Sex: %d (%s)\r\n", targetMob.Sex, sexToString(targetMob.Sex)))
	sb.WriteString(fmt.Sprintf("Class: %d\r\n", targetMob.Class))
	sb.WriteString(fmt.Sprintf("Race: %d\r\n", targetMob.Race))

	// Combat stats
	sb.WriteString(fmt.Sprintf("\r\nHP: %d/%d\r\n", targetMob.HP, targetMob.MaxHitPoints))
	sb.WriteString(fmt.Sprintf("Mana: %d/%d\r\n", targetMob.ManaPoints, targetMob.MaxManaPoints))
	sb.WriteString(fmt.Sprintf("Move: %d/%d\r\n", targetMob.MovePoints, targetMob.MaxMovePoints))
	sb.WriteString(fmt.Sprintf("AC: %v\r\n", targetMob.ArmorClass))
	sb.WriteString(fmt.Sprintf("HitRoll: %d\r\n", targetMob.HitRoll))
	sb.WriteString(fmt.Sprintf("DamRoll: %d\r\n", targetMob.DamRoll))

	// Dice information (for damage and HP calculation)
	if targetMob.Prototype != nil {
		sb.WriteString(fmt.Sprintf("Hit Dice: %dd%d+%d\r\n",
			targetMob.Prototype.Dice[0],
			targetMob.Prototype.Dice[1],
			targetMob.Prototype.Dice[2]))

		// Add damage dice information
		sb.WriteString(fmt.Sprintf("Damage Dice: %dd%d+%d\r\n",
			targetMob.Prototype.DamageType,
			targetMob.Prototype.AttackType,
			targetMob.DamRoll))
	}

	// Abilities
	sb.WriteString(fmt.Sprintf("\r\nAbilities: STR:%d INT:%d WIS:%d DEX:%d CON:%d CHA:%d\r\n",
		targetMob.Abilities[0], targetMob.Abilities[1], targetMob.Abilities[2],
		targetMob.Abilities[3], targetMob.Abilities[4], targetMob.Abilities[5]))

	// Flags
	sb.WriteString(fmt.Sprintf("\r\nAct Flags: %032b\r\n", targetMob.ActFlags))

	// Equipment
	sb.WriteString("\r\nEquipment:\r\n")
	hasEquipment := false
	for pos, item := range targetMob.Equipment {
		if item != nil {
			hasEquipment = true
			sb.WriteString(fmt.Sprintf("  %s: %s\r\n", wearPositionToString(pos), item.Prototype.ShortDesc))
		}
	}
	if !hasEquipment {
		sb.WriteString("  None\r\n")
	}

	// Inventory
	sb.WriteString("\r\nInventory:\r\n")
	if len(targetMob.Inventory) == 0 {
		sb.WriteString("  None\r\n")
	} else {
		for _, item := range targetMob.Inventory {
			sb.WriteString(fmt.Sprintf("  %s\r\n", item.Prototype.ShortDesc))
		}
	}

	// Gold and XP
	sb.WriteString(fmt.Sprintf("\r\nGold: %d\r\n", targetMob.Gold))
	sb.WriteString(fmt.Sprintf("Experience: %d\r\n", targetMob.Experience))
	sb.WriteString(fmt.Sprintf("Alignment: %d\r\n", targetMob.Alignment))

	// Send the stats to the character
	character.SendMessage(sb.String())
	return nil
}

// Helper functions to convert numeric values to strings
func positionToString(pos int) string {
	positions := map[int]string{
		types.POS_DEAD:     "DEAD",
		types.POS_MORTALLY: "MORTALLY WOUNDED",
		types.POS_INCAP:    "INCAPACITATED",
		types.POS_STUNNED:  "STUNNED",
		types.POS_SLEEPING: "SLEEPING",
		types.POS_RESTING:  "RESTING",
		types.POS_SITTING:  "SITTING",
		types.POS_FIGHTING: "FIGHTING",
		types.POS_STANDING: "STANDING",
	}

	if name, ok := positions[pos]; ok {
		return name
	}
	return "UNKNOWN"
}

func sexToString(sex int) string {
	sexes := map[int]string{
		0: "NEUTRAL",
		1: "MALE",
		2: "FEMALE",
	}

	if name, ok := sexes[sex]; ok {
		return name
	}
	return "UNKNOWN"
}

func wearPositionToString(pos int) string {
	positions := map[int]string{
		types.WEAR_LIGHT:    "Light",
		types.WEAR_FINGER_R: "Right Finger",
		types.WEAR_FINGER_L: "Left Finger",
		types.WEAR_NECK_1:   "Neck 1",
		types.WEAR_NECK_2:   "Neck 2",
		types.WEAR_BODY:     "Body",
		types.WEAR_HEAD:     "Head",
		types.WEAR_LEGS:     "Legs",
		types.WEAR_FEET:     "Feet",
		types.WEAR_HANDS:    "Hands",
		types.WEAR_ARMS:     "Arms",
		types.WEAR_SHIELD:   "Shield",
		types.WEAR_ABOUT:    "About Body",
		types.WEAR_WAIST:    "Waist",
		types.WEAR_WRIST_R:  "Right Wrist",
		types.WEAR_WRIST_L:  "Left Wrist",
		types.WEAR_WIELD:    "Wielded",
		types.WEAR_HOLD:     "Held",
	}

	if name, ok := positions[pos]; ok {
		return name
	}
	return fmt.Sprintf("Position %d", pos)
}
