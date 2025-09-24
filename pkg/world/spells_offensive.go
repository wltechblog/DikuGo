package world

import (
	"fmt"

	"github.com/wltechblog/DikuGo/pkg/types"
)

// CastMagicMissile casts the magic missile spell
func (w *World) CastMagicMissile(level int, ch *types.Character, arg string, spellType int, victim *types.Character, obj *types.ObjectInstance) error {
	// Magic missile only targets characters
	if victim == nil {
		return fmt.Errorf("you need to specify a target")
	}

	// Calculate damage
	damage := w.SpellDamage(level, types.SPELL_MAGIC_MISSILE)

	// Apply damage
	w.SpellDamageChar(ch, victim, damage, types.SPELL_MAGIC_MISSILE)

	return nil
}

// CastBurningHands casts the burning hands spell
func (w *World) CastBurningHands(level int, ch *types.Character, arg string, spellType int, victim *types.Character, obj *types.ObjectInstance) error {
	// Burning hands only targets characters
	if victim == nil {
		return fmt.Errorf("you need to specify a target")
	}

	// Calculate damage
	damage := w.SpellDamage(level, types.SPELL_BURNING_HANDS)

	// Apply damage
	w.SpellDamageChar(ch, victim, damage, types.SPELL_BURNING_HANDS)

	return nil
}

// CastChillTouch casts the chill touch spell
func (w *World) CastChillTouch(level int, ch *types.Character, arg string, spellType int, victim *types.Character, obj *types.ObjectInstance) error {
	// Chill touch only targets characters
	if victim == nil {
		return fmt.Errorf("you need to specify a target")
	}

	// Calculate damage
	damage := w.SpellDamage(level, types.SPELL_CHILL_TOUCH)

	// Apply damage
	w.SpellDamageChar(ch, victim, damage, types.SPELL_CHILL_TOUCH)

	// Apply strength reduction affect if victim doesn't save
	if !w.SavesSpell(victim, types.SAVING_SPELL) {
		// Create affect
		affect := &types.Affect{
			Type:      types.SPELL_CHILL_TOUCH,
			Duration:  6,
			Modifier:  -1,
			Location:  types.APPLY_STR,
			Bitvector: 0,
		}

		// Add affect to victim
		w.AffectJoin(victim, affect, true, false)
		victim.SendMessage("You feel your strength wither!\r\n")
	}

	return nil
}

// CastShockingGrasp casts the shocking grasp spell
func (w *World) CastShockingGrasp(level int, ch *types.Character, arg string, spellType int, victim *types.Character, obj *types.ObjectInstance) error {
	// Shocking grasp only targets characters
	if victim == nil {
		return fmt.Errorf("you need to specify a target")
	}

	// Calculate damage
	damage := w.SpellDamage(level, types.SPELL_SHOCKING_GRASP)

	// Apply damage
	w.SpellDamageChar(ch, victim, damage, types.SPELL_SHOCKING_GRASP)

	return nil
}

// CastLightningBolt casts the lightning bolt spell
func (w *World) CastLightningBolt(level int, ch *types.Character, arg string, spellType int, victim *types.Character, obj *types.ObjectInstance) error {
	// Lightning bolt only targets characters
	if victim == nil {
		return fmt.Errorf("you need to specify a target")
	}

	// Calculate damage
	damage := w.SpellDamage(level, types.SPELL_LIGHTNING_BOLT)

	// Apply damage
	w.SpellDamageChar(ch, victim, damage, types.SPELL_LIGHTNING_BOLT)

	return nil
}

// CastColorSpray casts the color spray spell
func (w *World) CastColorSpray(level int, ch *types.Character, arg string, spellType int, victim *types.Character, obj *types.ObjectInstance) error {
	// Color spray only targets characters
	if victim == nil {
		return fmt.Errorf("you need to specify a target")
	}

	// Calculate damage
	damage := w.SpellDamage(level, types.SPELL_COLOR_SPRAY)

	// Apply damage
	w.SpellDamageChar(ch, victim, damage, types.SPELL_COLOR_SPRAY)

	// Apply blindness if victim doesn't save and is not already blind
	if !w.SavesSpell(victim, types.SAVING_SPELL) && !w.AffectedBySpell(victim, types.SPELL_BLINDNESS) {
		// Create affect
		affect := &types.Affect{
			Type:      types.SPELL_BLINDNESS,
			Duration:  2,
			Modifier:  -4,
			Location:  types.APPLY_HITROLL,
			Bitvector: types.AFF_BLIND,
		}

		// Add affect to victim
		w.AffectToChar(victim, affect)
		victim.SendMessage("You are blinded by the spray of colors!\r\n")
		w.Act("$n appears to be blinded!", true, victim, nil, nil, types.TO_ROOM)
	}

	return nil
}

// CastFireball casts the fireball spell
func (w *World) CastFireball(level int, ch *types.Character, arg string, spellType int, victim *types.Character, obj *types.ObjectInstance) error {
	switch spellType {
	case types.SPELL_TYPE_SPELL, types.SPELL_TYPE_SCROLL, types.SPELL_TYPE_WAND:
		// Fireball only targets characters
		if victim == nil {
			return fmt.Errorf("you need to specify a target")
		}

		// Calculate damage
		damage := w.SpellDamage(level, types.SPELL_FIREBALL)

		// Apply damage
		w.SpellDamageChar(ch, victim, damage, types.SPELL_FIREBALL)

	case types.SPELL_TYPE_STAFF:
		// Staff affects everyone in the room except the caster
		w.Act("The room is filled with flames!", false, ch, nil, nil, types.TO_CHAR)
		w.Act("The room is filled with flames!", false, ch, nil, nil, types.TO_ROOM)

		// Calculate damage
		damage := w.SpellDamage(level, types.SPELL_FIREBALL)

		// Apply damage to everyone in the room except the caster
		for _, v := range ch.InRoom.Characters {
			if v != ch && v.Position > types.POS_DEAD {
				// Apply damage
				w.SpellDamageChar(ch, v, damage, types.SPELL_FIREBALL)
			}
		}
	}

	return nil
}

// CastHarm casts the harm spell
func (w *World) CastHarm(level int, ch *types.Character, arg string, spellType int, victim *types.Character, obj *types.ObjectInstance) error {
	// Harm only targets characters
	if victim == nil {
		return fmt.Errorf("you need to specify a target")
	}

	// Calculate damage
	damage := w.SpellDamage(level, types.SPELL_HARM)

	// Apply damage
	w.SpellDamageChar(ch, victim, damage, types.SPELL_HARM)

	return nil
}

// CastEnergyDrain casts the energy drain spell
func (w *World) CastEnergyDrain(level int, ch *types.Character, arg string, spellType int, victim *types.Character, obj *types.ObjectInstance) error {
	// Energy drain only targets characters
	if victim == nil {
		return fmt.Errorf("you need to specify a target")
	}

	// Calculate damage
	damage := w.SpellDamage(level, types.SPELL_ENERGY_DRAIN)

	// Apply damage
	w.SpellDamageChar(ch, victim, damage, types.SPELL_ENERGY_DRAIN)

	// Drain experience if victim doesn't save
	if !w.SavesSpell(victim, types.SAVING_SPELL) {
		// Drain 5% of experience
		expDrain := victim.Experience / 20
		if expDrain > 0 {
			victim.Experience -= expDrain
			victim.SendMessage("You feel less experienced!\r\n")
			ch.SendMessage("You feel more powerful as you drain their life force!\r\n")
		}
	}

	return nil
}

// CastDispelEvil casts the dispel evil spell
func (w *World) CastDispelEvil(level int, ch *types.Character, arg string, spellType int, victim *types.Character, obj *types.ObjectInstance) error {
	// Dispel evil only targets characters
	if victim == nil {
		return fmt.Errorf("you need to specify a target")
	}

	// Check alignment
	if ch.Alignment < 0 {
		victim = ch // Target self if caster is evil
		ch.SendMessage("The spell backfires!\r\n")
	}

	if victim.Alignment >= 0 {
		w.Act("$N doesn't seem to be affected.", false, ch, nil, victim, types.TO_CHAR)
		return nil
	}

	// Calculate damage
	damage := w.SpellDamage(level, types.SPELL_DISPEL_EVIL)

	// Apply damage
	w.SpellDamageChar(ch, victim, damage, types.SPELL_DISPEL_EVIL)

	return nil
}

// CastEarthquake casts the earthquake spell
func (w *World) CastEarthquake(level int, ch *types.Character, arg string, spellType int, victim *types.Character, obj *types.ObjectInstance) error {
	// Earthquake affects everyone in the room except the caster
	w.Act("The earth trembles beneath your feet!", false, ch, nil, nil, types.TO_CHAR)
	w.Act("The earth trembles beneath your feet!", false, ch, nil, nil, types.TO_ROOM)

	// Calculate damage
	damage := w.SpellDamage(level, types.SPELL_EARTHQUAKE)

	// Apply damage to everyone in the room except the caster
	for _, victim := range ch.InRoom.Characters {
		if victim != ch && victim.Position > types.POS_DEAD {
			// Apply damage
			w.SpellDamageChar(ch, victim, damage, types.SPELL_EARTHQUAKE)
		}
	}

	return nil
}

// CastCallLightning casts the call lightning spell
func (w *World) CastCallLightning(level int, ch *types.Character, arg string, spellType int, victim *types.Character, obj *types.ObjectInstance) error {
	// Call lightning only targets characters
	if victim == nil {
		return fmt.Errorf("you need to specify a target")
	}

	// Check if outdoors
	if ch.InRoom.SectorType == types.SECT_INSIDE || (ch.InRoom.Flags&types.ROOM_INDOORS) != 0 {
		ch.SendMessage("You must be outdoors to call lightning.\r\n")
		return nil
	}

	// Check weather
	if w.time.Weather < types.SKY_RAINY {
		ch.SendMessage("You need bad weather to call lightning.\r\n")
		return nil
	}

	// Calculate damage
	damage := w.SpellDamage(level, types.SPELL_CALL_LIGHTNING)

	// Apply damage
	w.SpellDamageChar(ch, victim, damage, types.SPELL_CALL_LIGHTNING)

	return nil
}
