package world

import (
	"github.com/wltechblog/DikuGo/pkg/types"
)

// PulsePointUpdate updates character points (HP, mana, move) and conditions
func (w *World) PulsePointUpdate() {
	// Make a copy of the characters to avoid deadlock
	w.mutex.RLock()
	characters := make([]*types.Character, 0, len(w.characters))
	for _, character := range w.characters {
		characters = append(characters, character)
	}
	w.mutex.RUnlock()

	// Process point updates for all characters
	for _, character := range characters {
		// Skip dead characters
		if character.Position <= types.POS_DEAD {
			continue
		}

		// Update hit points
		if character.Position > types.POS_STUNNED {
			// Normal regeneration
			hitGain := w.hitGain(character)
			character.HP = min(character.HP+hitGain, character.MaxHitPoints)

			// Update mana points
			manaGain := w.manaGain(character)
			character.ManaPoints = min(character.ManaPoints+manaGain, character.MaxManaPoints)

			// Update move points
			moveGain := w.moveGain(character)
			character.MovePoints = min(character.MovePoints+moveGain, character.MaxMovePoints)
		} else if character.Position == types.POS_STUNNED {
			// Stunned characters still regenerate, but may recover
			hitGain := w.hitGain(character)
			character.HP = min(character.HP+hitGain, character.MaxHitPoints)

			manaGain := w.manaGain(character)
			character.ManaPoints = min(character.ManaPoints+manaGain, character.MaxManaPoints)

			moveGain := w.moveGain(character)
			character.MovePoints = min(character.MovePoints+moveGain, character.MaxMovePoints)

			// Check if they should recover from being stunned
			w.updatePosition(character)
		} else if character.Position == types.POS_INCAP {
			// Incapacitated characters take damage
			w.damage(character, character, 1, types.TYPE_SUFFERING)
		} else if !character.IsNPC && character.Position == types.POS_MORTALLY {
			// Mortally wounded players take more damage
			w.damage(character, character, 2, types.TYPE_SUFFERING)
		}

		// Update conditions for players (hunger, thirst, drunk)
		if !character.IsNPC {
			w.gainCondition(character, types.COND_FULL, -1)
			w.gainCondition(character, types.COND_THIRST, -1)
			w.gainCondition(character, types.COND_DRUNK, -1)
		}
	}
}

// hitGain calculates how many hit points a character gains per tick
func (w *World) hitGain(ch *types.Character) int {
	var gain int

	if ch.IsNPC {
		// NPCs gain based on level
		gain = ch.Level
	} else {
		// Players gain based on constitution and position
		gain = 2 // Base gain

		// Adjust based on class
		if ch.Class == types.CLASS_WARRIOR {
			gain += 3
		} else if ch.Class == types.CLASS_THIEF {
			gain += 2
		} else if ch.Class == types.CLASS_CLERIC {
			gain += 1
		}

		// Position calculations
		switch ch.Position {
		case types.POS_SLEEPING:
			gain += gain / 2 // 50% bonus when sleeping
		case types.POS_RESTING:
			gain += gain / 4 // 25% bonus when resting
		case types.POS_SITTING:
			gain += gain / 8 // 12.5% bonus when sitting
		}

		// Magic users and clerics gain less hit points
		if ch.Class == types.CLASS_MAGE || ch.Class == types.CLASS_CLERIC {
			gain /= 2
		}
	}

	// Reduce gain if poisoned
	if ch.AffectedBy&types.AFF_POISON != 0 {
		gain /= 4
		w.damage(ch, ch, 2, types.TYPE_POISON)
	}

	// Reduce gain if hungry or thirsty
	if !ch.IsNPC && (ch.Conditions[types.COND_FULL] == 0 || ch.Conditions[types.COND_THIRST] == 0) {
		gain /= 4
	}

	return max(gain, 1) // Always gain at least 1 point
}

// manaGain calculates how many mana points a character gains per tick
func (w *World) manaGain(ch *types.Character) int {
	var gain int

	if ch.IsNPC {
		// NPCs gain based on level
		gain = ch.Level
	} else {
		// Players gain based on intelligence/wisdom and position
		gain = 2 // Base gain

		// Adjust based on class
		if ch.Class == types.CLASS_MAGE {
			gain += 3
		} else if ch.Class == types.CLASS_CLERIC {
			gain += 3
		} else if ch.Class == types.CLASS_THIEF {
			gain += 1
		}

		// Position calculations
		switch ch.Position {
		case types.POS_SLEEPING:
			gain *= 2 // 100% bonus when sleeping
		case types.POS_RESTING:
			gain += gain / 2 // 50% bonus when resting
		case types.POS_SITTING:
			gain += gain / 4 // 25% bonus when sitting
		}

		// Magic users and clerics gain more mana
		if ch.Class == types.CLASS_MAGE || ch.Class == types.CLASS_CLERIC {
			gain *= 2
		}
	}

	// Reduce gain if poisoned
	if ch.AffectedBy&types.AFF_POISON != 0 {
		gain /= 4
	}

	// Reduce gain if hungry or thirsty
	if !ch.IsNPC && (ch.Conditions[types.COND_FULL] == 0 || ch.Conditions[types.COND_THIRST] == 0) {
		gain /= 4
	}

	return max(gain, 1) // Always gain at least 1 point
}

// moveGain calculates how many move points a character gains per tick
func (w *World) moveGain(ch *types.Character) int {
	var gain int

	if ch.IsNPC {
		// NPCs gain based on level
		gain = ch.Level
	} else {
		// Players gain based on dexterity and position
		gain = 5 // Base gain

		// Adjust based on class
		if ch.Class == types.CLASS_THIEF {
			gain += 5
		} else if ch.Class == types.CLASS_WARRIOR {
			gain += 3
		} else if ch.Class == types.CLASS_CLERIC {
			gain += 1
		}

		// Position calculations
		switch ch.Position {
		case types.POS_SLEEPING:
			gain += gain / 2 // 50% bonus when sleeping
		case types.POS_RESTING:
			gain += gain / 4 // 25% bonus when resting
		case types.POS_SITTING:
			gain += gain / 8 // 12.5% bonus when sitting
		}
	}

	// Reduce gain if poisoned
	if ch.AffectedBy&types.AFF_POISON != 0 {
		gain /= 4
	}

	// Reduce gain if hungry or thirsty
	if !ch.IsNPC && (ch.Conditions[types.COND_FULL] == 0 || ch.Conditions[types.COND_THIRST] == 0) {
		gain /= 4
	}

	return max(gain, 1) // Always gain at least 1 point
}

// updatePosition updates a character's position based on their hit points
func (w *World) updatePosition(ch *types.Character) {
	if ch.HP > 0 && ch.Position > types.POS_STUNNED {
		return
	} else if ch.HP > 0 {
		ch.Position = types.POS_STANDING
	} else if ch.HP <= -11 {
		ch.Position = types.POS_DEAD
	} else if ch.HP <= -6 {
		ch.Position = types.POS_MORTALLY
	} else if ch.HP <= -3 {
		ch.Position = types.POS_INCAP
	} else {
		ch.Position = types.POS_STUNNED
	}
}

// damage applies damage to a character
func (w *World) damage(ch, victim *types.Character, dam int, damageType int) {
	// Skip if victim is already dead
	if victim.Position <= types.POS_DEAD {
		return
	}

	// Apply damage
	victim.HP -= dam

	// Update position based on new hit points
	w.updatePosition(victim)

	// Check if the victim died
	if victim.Position == types.POS_DEAD {
		// Handle death
		w.HandleCharacterDeath(victim)
	}
}

// gainCondition updates a character's condition (hunger, thirst, drunk)
func (w *World) gainCondition(ch *types.Character, condition int, value int) {
	// Skip NPCs and immortals
	if ch.IsNPC || ch.Level >= 20 {
		return
	}

	// Skip if condition is disabled (-1)
	if ch.Conditions[condition] == -1 {
		return
	}

	// Remember if character was intoxicated
	wasIntoxicated := (condition == types.COND_DRUNK && ch.Conditions[condition] > 0)

	// Update condition
	ch.Conditions[condition] += value

	// Constrain condition value
	ch.Conditions[condition] = max(0, min(24, ch.Conditions[condition]))

	// Send messages if condition reaches 0
	if ch.Conditions[condition] == 0 {
		switch condition {
		case types.COND_FULL:
			ch.SendMessage("You are hungry.\r\n")
		case types.COND_THIRST:
			ch.SendMessage("You are thirsty.\r\n")
		case types.COND_DRUNK:
			if wasIntoxicated {
				ch.SendMessage("You are now sober.\r\n")
			}
		}
	}
}

// Helper functions
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
