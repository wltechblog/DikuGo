package world

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/wltechblog/DikuGo/pkg/types"
	"github.com/wltechblog/DikuGo/pkg/utils"
)

// SpellDamage calculates damage for a spell based on level
func (w *World) SpellDamage(level int, spell int) int {
	var baseDamage, minDamage, maxDamage int

	switch spell {
	case types.SPELL_MAGIC_MISSILE:
		// Magic missile damage: 1d4+1 per level, max 5d4+5
		numDice := level
		if numDice > 5 {
			numDice = 5
		}
		baseDamage = utils.Dice(numDice, 4) + numDice
	case types.SPELL_BURNING_HANDS:
		// Burning hands damage: 1d4+1 per level, max 5d4+5
		if level < 5 {
			return 0 // Minimum level 5
		}
		numDice := level - 4
		if numDice > 5 {
			numDice = 5
		}
		baseDamage = utils.Dice(numDice, 4) + numDice
	case types.SPELL_CHILL_TOUCH:
		// Chill touch damage: 1d6+1 per level, max 5d6+5
		if level < 3 {
			return 0 // Minimum level 3
		}
		numDice := level - 2
		if numDice > 5 {
			numDice = 5
		}
		baseDamage = utils.Dice(numDice, 6) + numDice
	case types.SPELL_SHOCKING_GRASP:
		// Shocking grasp damage: 1d8+1 per level, max 5d8+5
		if level < 7 {
			return 0 // Minimum level 7
		}
		numDice := level - 6
		if numDice > 5 {
			numDice = 5
		}
		baseDamage = utils.Dice(numDice, 8) + numDice
	case types.SPELL_LIGHTNING_BOLT:
		// Lightning bolt damage: 1d10+1 per level, max 10d10+10
		if level < 9 {
			return 0 // Minimum level 9
		}
		numDice := level - 8
		if numDice > 10 {
			numDice = 10
		}
		baseDamage = utils.Dice(numDice, 10) + numDice
	case types.SPELL_COLOR_SPRAY:
		// Color spray damage: 1d12+1 per level, max 10d12+10
		if level < 11 {
			return 0 // Minimum level 11
		}
		numDice := level - 10
		if numDice > 10 {
			numDice = 10
		}
		baseDamage = utils.Dice(numDice, 12) + numDice
	case types.SPELL_FIREBALL:
		// Fireball damage: 1d10+2 per level, max 10d10+20
		if level < 15 {
			return 0 // Minimum level 15
		}
		numDice := level - 14
		if numDice > 10 {
			numDice = 10
		}
		baseDamage = utils.Dice(numDice, 10) + (numDice * 2)
	case types.SPELL_HARM:
		// Harm damage: 1d8+2 per level, max 15d8+30
		if level < 15 {
			return 0 // Minimum level 15
		}
		numDice := level - 14
		if numDice > 15 {
			numDice = 15
		}
		baseDamage = utils.Dice(numDice, 8) + (numDice * 2)
	default:
		// Default damage calculation
		minDamage = level
		maxDamage = level * 3
		baseDamage = rand.Intn(maxDamage-minDamage+1) + minDamage
	}

	return baseDamage
}

// SavesSpell checks if a character saves against a spell
func (w *World) SavesSpell(ch *types.Character, saveType int) bool {
	// Base save value from character's saving throw
	save := ch.SavingThrow[saveType]

	// Adjust based on level and class
	if !ch.IsNPC {
		// TODO: Implement saving throw tables based on class and level
		// For now, use a simple formula
		save += 20 - ch.Level
	}

	// Roll 1d20, if roll is greater than save, the save fails
	roll := rand.Intn(20) + 1
	return roll > save
}

// DamageMessage returns a message for spell damage
func (w *World) DamageMessage(damage int, spell int) string {
	var message string

	switch spell {
	case types.SPELL_MAGIC_MISSILE:
		if damage < 10 {
			message = "The magic missile hits $N with a soft thud."
		} else if damage < 20 {
			message = "The magic missile strikes $N with force!"
		} else {
			message = "The magic missile BLASTS $N with incredible power!"
		}
	case types.SPELL_BURNING_HANDS:
		if damage < 15 {
			message = "Your burning hands singe $N."
		} else if damage < 30 {
			message = "Your burning hands scorch $N!"
		} else {
			message = "Your burning hands INCINERATE $N!"
		}
	case types.SPELL_CHILL_TOUCH:
		if damage < 15 {
			message = "Your chill touch causes $N to shiver."
		} else if damage < 30 {
			message = "Your chill touch freezes $N!"
		} else {
			message = "Your chill touch FREEZES $N solid!"
		}
	case types.SPELL_SHOCKING_GRASP:
		if damage < 20 {
			message = "Your shocking grasp zaps $N."
		} else if damage < 40 {
			message = "Your shocking grasp jolts $N!"
		} else {
			message = "Your shocking grasp ELECTROCUTES $N!"
		}
	case types.SPELL_LIGHTNING_BOLT:
		if damage < 30 {
			message = "Your lightning bolt zaps $N."
		} else if damage < 60 {
			message = "Your lightning bolt strikes $N with a loud CRACK!"
		} else {
			message = "Your lightning bolt DEVASTATES $N with electrical fury!"
		}
	case types.SPELL_COLOR_SPRAY:
		if damage < 40 {
			message = "Your color spray dazzles $N."
		} else if damage < 80 {
			message = "Your color spray blinds $N!"
		} else {
			message = "Your color spray OBLITERATES $N with psychedelic energy!"
		}
	case types.SPELL_FIREBALL:
		if damage < 50 {
			message = "Your fireball singes $N."
		} else if damage < 100 {
			message = "Your fireball engulfs $N in flames!"
		} else {
			message = "Your fireball IMMOLATES $N in a massive explosion!"
		}
	case types.SPELL_HARM:
		if damage < 60 {
			message = "Your harm spell hurts $N."
		} else if damage < 120 {
			message = "Your harm spell causes $N to writhe in agony!"
		} else {
			message = "Your harm spell DEVASTATES $N with divine fury!"
		}
	default:
		message = "Your spell hits $N."
	}

	return message
}

// SpellDamage applies damage from a spell
func (w *World) SpellDamageChar(ch *types.Character, victim *types.Character, damage int, spell int) {
	// Check if the victim saves against the spell
	if w.SavesSpell(victim, types.SAVING_SPELL) {
		damage = damage / 2
		victim.SendMessage("You partially resist the spell.\r\n")
	}

	// Get damage message
	message := w.DamageMessage(damage, spell)

	// Send messages
	w.Act(message, false, ch, nil, victim, types.TO_CHAR)
	w.Act("$n's spell hits you!", false, ch, nil, victim, types.TO_VICT)
	w.Act("$n's spell hits $N!", false, ch, nil, victim, types.TO_NOTVICT)

	// Apply damage
	w.Damage(ch, victim, damage, spell)
}

// SaySpell makes a character say the spell words
func (w *World) SaySpell(ch *types.Character, spell int) {
	// Skip for NPCs
	if ch.IsNPC {
		return
	}

	// Get spell name
	spellName := types.GetSpellName(spell)
	if spellName == "unknown" {
		return
	}

	// Create a garbled version for non-class members
	garbled := ""
	for _, c := range spellName {
		if rand.Intn(3) == 0 {
			garbled += string('a' + rune(rand.Intn(26)))
		} else {
			garbled += string(c)
		}
	}

	// Send messages
	w.Act("$n utters the words, '"+spellName+"'", false, ch, nil, nil, types.TO_ROOM)
}

// GetSpellTarget finds a target for a spell
func (w *World) GetSpellTarget(ch *types.Character, arg string, spell int) (*types.Character, *types.ObjectInstance, error) {
	// Get spell info
	targets := types.GetSpellTargets(spell)
	if targets == 0 {
		return nil, nil, fmt.Errorf("spell has no valid targets")
	}

	// If no argument and spell can target self, target self
	if arg == "" && (targets&types.TAR_SELF_ONLY != 0 || targets&types.TAR_CHAR_ROOM != 0) {
		// Check if spell can't target self
		if targets&types.TAR_SELF_NONO != 0 {
			return nil, nil, fmt.Errorf("you cannot cast this spell on yourself")
		}
		return ch, nil, nil
	}

	// If no argument and spell requires a target
	if arg == "" {
		if targets < types.TAR_OBJ_INV {
			return nil, nil, fmt.Errorf("who should the spell be cast upon?")
		} else {
			return nil, nil, fmt.Errorf("what should the spell be cast upon?")
		}
	}

	// Try to find a character target
	if targets&types.TAR_CHAR_ROOM != 0 {
		victim := w.GetCharacterInRoom(ch.InRoom, arg)
		if victim != nil {
			// Check if spell can only target self
			if targets&types.TAR_SELF_ONLY != 0 && victim != ch {
				return nil, nil, fmt.Errorf("you can only cast this spell upon yourself")
			}
			// Check if spell can't target self
			if targets&types.TAR_SELF_NONO != 0 && victim == ch {
				return nil, nil, fmt.Errorf("you cannot cast this spell on yourself")
			}
			return victim, nil, nil
		}
	}

	// Try to find an object target in inventory
	if targets&types.TAR_OBJ_INV != 0 {
		obj := w.GetObjectInList(ch, arg, ch.Inventory)
		if obj != nil {
			return nil, obj, nil
		}
	}

	// Try to find an object target in room
	if targets&types.TAR_OBJ_ROOM != 0 {
		obj := w.GetObjectInRoom(ch.InRoom, arg)
		if obj != nil {
			return nil, obj, nil
		}
	}

	// Try to find an object target in equipment
	if targets&types.TAR_OBJ_EQUIP != 0 {
		obj := w.GetObjectInEquipment(ch, arg)
		if obj != nil {
			return nil, obj, nil
		}
	}

	// No target found
	return nil, nil, fmt.Errorf("nothing by that name here")
}

// GetObjectInList finds an object in a list by name
func (w *World) GetObjectInList(ch *types.Character, name string, list []*types.ObjectInstance) *types.ObjectInstance {
	name = strings.ToLower(name)
	for _, obj := range list {
		if strings.Contains(strings.ToLower(obj.Prototype.Name), name) {
			return obj
		}
	}
	return nil
}

// GetObjectInRoom finds an object in a room by name
func (w *World) GetObjectInRoom(room *types.Room, name string) *types.ObjectInstance {
	name = strings.ToLower(name)
	for _, obj := range room.Objects {
		if strings.Contains(strings.ToLower(obj.Prototype.Name), name) {
			return obj
		}
	}
	return nil
}

// GetObjectInEquipment finds an object in equipment by name
func (w *World) GetObjectInEquipment(ch *types.Character, name string) *types.ObjectInstance {
	name = strings.ToLower(name)
	for _, obj := range ch.Equipment {
		if obj != nil && strings.Contains(strings.ToLower(obj.Prototype.Name), name) {
			return obj
		}
	}
	return nil
}
