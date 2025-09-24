package world

import (
	"fmt"

	"github.com/wltechblog/DikuGo/pkg/types"
)

// CastSleep casts the sleep spell
func (w *World) CastSleep(level int, ch *types.Character, arg string, spellType int, victim *types.Character, obj *types.ObjectInstance) error {
	// Handle different spell types
	switch spellType {
	case types.SPELL_TYPE_SPELL:
		if victim == nil {
			return fmt.Errorf("whom do you wish to sleep?")
		}

		// Check if victim is already asleep
		if victim.Position == types.POS_SLEEPING {
			ch.SendMessage("They are already asleep.\r\n")
			return nil
		}

		// Check if victim is too high level
		if victim.Level > ch.Level+2 {
			ch.SendMessage("Your spell is not powerful enough.\r\n")
			return nil
		}

		// Check if victim is fighting
		if victim.Fighting != nil {
			ch.SendMessage("You cannot sleep someone who is fighting!\r\n")
			return nil
		}

		// Check if victim saves against spell
		if w.SavesSpell(victim, types.SAVING_SPELL) {
			ch.SendMessage("They resist your sleep spell.\r\n")
			victim.SendMessage("You feel drowsy for a moment, but it passes.\r\n")
			return nil
		}

		// Create affect
		affect := &types.Affect{
			Type:      types.SPELL_SLEEP,
			Duration:  4 + (level / 4),
			Modifier:  0,
			Location:  types.APPLY_NONE,
			Bitvector: types.AFF_SLEEP,
		}

		// Add affect to victim
		w.AffectToChar(victim, affect)

		// Set victim's position to sleeping
		if victim.Position > types.POS_SLEEPING {
			victim.Position = types.POS_SLEEPING
		}

		// Send messages
		victim.SendMessage("You feel very sleepy...zzzzzz\r\n")
		w.Act("$n goes to sleep.", true, victim, nil, nil, types.TO_ROOM)

	case types.SPELL_TYPE_POTION:
		// Check if already asleep
		if ch.Position == types.POS_SLEEPING {
			return nil
		}

		// Create affect
		affect := &types.Affect{
			Type:      types.SPELL_SLEEP,
			Duration:  4 + (level / 4),
			Modifier:  0,
			Location:  types.APPLY_NONE,
			Bitvector: types.AFF_SLEEP,
		}

		// Add affect to drinker
		w.AffectToChar(ch, affect)

		// Set position to sleeping
		if ch.Position > types.POS_SLEEPING {
			ch.Position = types.POS_SLEEPING
		}

		ch.SendMessage("You feel very sleepy...zzzzzz\r\n")
		w.Act("$n goes to sleep.", true, ch, nil, nil, types.TO_ROOM)

	case types.SPELL_TYPE_SCROLL:
		// If no target specified, target self
		if victim == nil {
			victim = ch
		}

		// Check if victim is already asleep
		if victim.Position == types.POS_SLEEPING {
			return nil
		}

		// Check if victim saves against spell
		if victim != ch && w.SavesSpell(victim, types.SAVING_SPELL) {
			return nil
		}

		// Create affect
		affect := &types.Affect{
			Type:      types.SPELL_SLEEP,
			Duration:  4 + (level / 4),
			Modifier:  0,
			Location:  types.APPLY_NONE,
			Bitvector: types.AFF_SLEEP,
		}

		// Add affect to victim
		w.AffectToChar(victim, affect)

		// Set victim's position to sleeping
		if victim.Position > types.POS_SLEEPING {
			victim.Position = types.POS_SLEEPING
		}

		victim.SendMessage("You feel very sleepy...zzzzzz\r\n")
		w.Act("$n goes to sleep.", true, victim, nil, nil, types.TO_ROOM)
	}

	return nil
}

// CastCharmPerson casts the charm person spell
func (w *World) CastCharmPerson(level int, ch *types.Character, arg string, spellType int, victim *types.Character, obj *types.ObjectInstance) error {
	if victim == nil {
		return fmt.Errorf("whom do you wish to charm?")
	}

	// Check if victim is self
	if victim == ch {
		ch.SendMessage("You like yourself even better!\r\n")
		return nil
	}

	// Check if victim is already charmed
	if w.AffectedBySpell(victim, types.SPELL_CHARM_PERSON) {
		ch.SendMessage("They are already charmed.\r\n")
		return nil
	}

	// Check if victim is an NPC
	if !victim.IsNPC {
		ch.SendMessage("You can only charm monsters.\r\n")
		return nil
	}

	// Check if victim is too high level
	if victim.Level > ch.Level {
		ch.SendMessage("Your spell is not powerful enough.\r\n")
		return nil
	}

	// Check if victim saves against spell
	if w.SavesSpell(victim, types.SAVING_PARA) {
		ch.SendMessage("They resist your charm spell.\r\n")
		return nil
	}

	// Check if victim is already following someone
	if victim.Following != nil {
		ch.SendMessage("They are already following someone else.\r\n")
		return nil
	}

	// Create affect
	affect := &types.Affect{
		Type:      types.SPELL_CHARM_PERSON,
		Duration:  24,
		Modifier:  0,
		Location:  types.APPLY_NONE,
		Bitvector: types.AFF_CHARM,
	}

	// Add affect to victim
	w.AffectToChar(victim, affect)

	// Set victim to follow caster
	victim.Following = ch
	ch.Followers = append(ch.Followers, victim)

	// Stop victim from fighting
	if victim.Fighting != nil {
		victim.Fighting = nil
	}

	// Send messages
	w.Act("Isn't $n just such a nice fellow?", false, ch, nil, victim, types.TO_VICT)
	ch.SendMessage("They are now charmed.\r\n")

	return nil
}
