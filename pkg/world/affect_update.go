package world

import (
	"log"

	"github.com/wltechblog/DikuGo/pkg/types"
)

// PulseAffectUpdate updates all character affects (spells, skills, etc.)
func (w *World) PulseAffectUpdate() {
	// Make a copy of the characters to avoid deadlock
	w.mutex.RLock()
	characters := make([]*types.Character, 0, len(w.characters))
	for _, character := range w.characters {
		characters = append(characters, character)
	}
	w.mutex.RUnlock()

	// Process affect updates for all characters
	for _, character := range characters {
		// Skip dead characters
		if character.Position <= types.POS_DEAD {
			continue
		}

		// Process each affect
		var nextAffect *types.Affect
		for affect := character.Affected; affect != nil; affect = nextAffect {
			// Store next affect in case this one is removed
			nextAffect = affect.Next

			// Skip permanent affects
			if affect.Duration < 0 {
				continue
			}

			// Decrement duration
			affect.Duration--

			// Check if the affect has expired
			if affect.Duration <= 0 {
				// Send wear-off message if applicable
				if affect.Type > 0 && affect.Type < len(types.SpellWearOffMessages) {
					message := types.SpellWearOffMessages[affect.Type]
					if message != "" {
						character.SendMessage(message + "\r\n")
					}
				}

				// Remove the affect
				w.affectRemove(character, affect)
			}
		}
	}
}

// affectRemove removes an affect from a character
func (w *World) affectRemove(ch *types.Character, af *types.Affect) {
	// Remove the affect's modifiers
	w.affectModify(ch, af.Location, af.Modifier, af.Bitvector, false)

	// Remove the affect from the character's list
	if ch.Affected == af {
		// Remove from head of list
		ch.Affected = af.Next
	} else {
		// Find the affect in the list
		for prev := ch.Affected; prev != nil; prev = prev.Next {
			if prev.Next == af {
				// Remove from middle/end of list
				prev.Next = af.Next
				break
			}
		}
	}

	// Update the character's total affects
	w.affectTotal(ch)
}

// affectModify applies or removes an affect's modifiers to a character
func (w *World) affectModify(ch *types.Character, location int, modifier int, bitvector int64, add bool) {
	// Apply or remove bitvector flags
	if add {
		ch.AffectedBy |= bitvector
	} else {
		ch.AffectedBy &= ^bitvector
		modifier = -modifier // Invert modifier for removal
	}

	// Apply modifier based on location
	switch location {
	case types.APPLY_NONE:
		// No effect
	case types.APPLY_STR:
		ch.Abilities[types.ABILITY_STR] += modifier
	case types.APPLY_DEX:
		ch.Abilities[types.ABILITY_DEX] += modifier
	case types.APPLY_INT:
		ch.Abilities[types.ABILITY_INT] += modifier
	case types.APPLY_WIS:
		ch.Abilities[types.ABILITY_WIS] += modifier
	case types.APPLY_CON:
		ch.Abilities[types.ABILITY_CON] += modifier
	case types.APPLY_CHAR_HEIGHT: // Using as CHA for now
		ch.Abilities[types.ABILITY_CHA] += modifier
	case types.APPLY_HIT:
		ch.MaxHitPoints += modifier
	case types.APPLY_MANA:
		ch.MaxManaPoints += modifier
	case types.APPLY_MOVE:
		ch.MaxMovePoints += modifier
	case types.APPLY_AC:
		// Apply to all AC types
		for i := range ch.ArmorClass {
			ch.ArmorClass[i] += modifier
		}
	case types.APPLY_HITROLL:
		ch.HitRoll += modifier
	case types.APPLY_DAMROLL:
		ch.DamRoll += modifier
	case types.APPLY_SAVING_PARA:
		ch.SavingThrow[types.SAVING_PARA] += modifier
	case types.APPLY_SAVING_ROD:
		ch.SavingThrow[types.SAVING_ROD] += modifier
	case types.APPLY_SAVING_PETRI:
		ch.SavingThrow[types.SAVING_PETRI] += modifier
	case types.APPLY_SAVING_BREATH:
		ch.SavingThrow[types.SAVING_BREATH] += modifier
	case types.APPLY_SAVING_SPELL:
		ch.SavingThrow[types.SAVING_SPELL] += modifier
	default:
		log.Printf("WARNING: Unknown affect location %d in affectModify", location)
	}
}

// affectTotal recalculates all affects on a character
func (w *World) affectTotal(ch *types.Character) {
	// Store original abilities
	originalAbilities := ch.Abilities

	// Reset all affected abilities to original values
	for i := range ch.Abilities {
		ch.Abilities[i] = originalAbilities[i]
	}

	// Apply equipment affects
	for pos := 0; pos < types.NUM_WEARS; pos++ {
		if ch.Equipment[pos] != nil {
			for j := 0; j < types.MAX_OBJ_AFFECT; j++ {
				affect := ch.Equipment[pos].Affects[j]
				if affect.Location != types.APPLY_NONE {
					w.affectModify(ch, affect.Location, affect.Modifier, 0, true)
				}
			}
		}
	}

	// Apply spell/skill affects
	for affect := ch.Affected; affect != nil; affect = affect.Next {
		w.affectModify(ch, affect.Location, affect.Modifier, affect.Bitvector, true)
	}
}

// affectToChar adds an affect to a character
func (w *World) AffectToChar(ch *types.Character, af *types.Affect) {
	// Create a new affect
	newAffect := &types.Affect{
		Type:      af.Type,
		Duration:  af.Duration,
		Modifier:  af.Modifier,
		Location:  af.Location,
		Bitvector: af.Bitvector,
	}

	// Add to head of list
	newAffect.Next = ch.Affected
	ch.Affected = newAffect

	// Apply the affect's modifiers
	w.affectModify(ch, af.Location, af.Modifier, af.Bitvector, true)

	// Update total affects
	w.affectTotal(ch)
}

// AffectJoin adds an affect to a character, combining with existing affects of the same type
func (w *World) AffectJoin(ch *types.Character, af *types.Affect, avgDuration bool, avgModifier bool) {
	var found bool

	// Look for an existing affect of the same type
	for existing := ch.Affected; existing != nil && !found; existing = existing.Next {
		if existing.Type == af.Type {
			// Found an existing affect of the same type
			af.Duration += existing.Duration
			if avgDuration {
				af.Duration /= 2
			}

			af.Modifier += existing.Modifier
			if avgModifier {
				af.Modifier /= 2
			}

			// Remove the existing affect
			w.affectRemove(ch, existing)

			// Add the new combined affect
			w.AffectToChar(ch, af)
			found = true
		}
	}

	// If no existing affect was found, just add the new one
	if !found {
		w.AffectToChar(ch, af)
	}
}

// AffectedBySpell checks if a character is affected by a specific spell
func (w *World) AffectedBySpell(ch *types.Character, spellType int) bool {
	for affect := ch.Affected; affect != nil; affect = affect.Next {
		if affect.Type == spellType {
			return true
		}
	}
	return false
}

// AffectFromChar removes all affects of a specific type from a character
func (w *World) AffectFromChar(ch *types.Character, spellType int) {
	var nextAffect *types.Affect

	for affect := ch.Affected; affect != nil; affect = nextAffect {
		nextAffect = affect.Next
		if affect.Type == spellType {
			w.affectRemove(ch, affect)
		}
	}
}
