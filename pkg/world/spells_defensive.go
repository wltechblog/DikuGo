package world

import (
	"fmt"

	"github.com/wltechblog/DikuGo/pkg/types"
)

// CastArmor casts the armor spell
func (w *World) CastArmor(level int, ch *types.Character, arg string, spellType int, victim *types.Character, obj *types.ObjectInstance) error {
	// Handle different spell types
	switch spellType {
	case types.SPELL_TYPE_SPELL:
		// Check if target is already affected
		if w.AffectedBySpell(victim, types.SPELL_ARMOR) {
			if victim == ch {
				ch.SendMessage("You are already protected by armor.\r\n")
			} else {
				ch.SendMessage("They are already protected by armor.\r\n")
			}
			return nil
		}

		// Create affect
		affect := &types.Affect{
			Type:      types.SPELL_ARMOR,
			Duration:  24,
			Modifier:  -20,
			Location:  types.APPLY_AC,
			Bitvector: 0,
		}

		// Add affect to victim
		w.AffectToChar(victim, affect)

		// Send messages
		if victim == ch {
			ch.SendMessage("You feel protected.\r\n")
		} else {
			victim.SendMessage("You feel protected.\r\n")
			w.Act("$N is protected by your armor spell.", false, ch, nil, victim, types.TO_CHAR)
		}

	case types.SPELL_TYPE_POTION:
		// Check if already affected
		if w.AffectedBySpell(ch, types.SPELL_ARMOR) {
			return nil
		}

		// Create affect
		affect := &types.Affect{
			Type:      types.SPELL_ARMOR,
			Duration:  24,
			Modifier:  -20,
			Location:  types.APPLY_AC,
			Bitvector: 0,
		}

		// Add affect to drinker
		w.AffectToChar(ch, affect)
		ch.SendMessage("You feel protected.\r\n")

	case types.SPELL_TYPE_SCROLL:
		// If no target specified, target self
		if victim == nil {
			victim = ch
		}

		// Check if already affected
		if w.AffectedBySpell(victim, types.SPELL_ARMOR) {
			return nil
		}

		// Create affect
		affect := &types.Affect{
			Type:      types.SPELL_ARMOR,
			Duration:  24,
			Modifier:  -20,
			Location:  types.APPLY_AC,
			Bitvector: 0,
		}

		// Add affect to victim
		w.AffectToChar(victim, affect)
		victim.SendMessage("You feel protected.\r\n")

	case types.SPELL_TYPE_WAND:
		// If no target specified, fail
		if victim == nil {
			return nil
		}

		// Check if already affected
		if w.AffectedBySpell(victim, types.SPELL_ARMOR) {
			return nil
		}

		// Create affect
		affect := &types.Affect{
			Type:      types.SPELL_ARMOR,
			Duration:  24,
			Modifier:  -20,
			Location:  types.APPLY_AC,
			Bitvector: 0,
		}

		// Add affect to victim
		w.AffectToChar(victim, affect)
		victim.SendMessage("You feel protected.\r\n")
		w.Act("$N is surrounded by a magical armor.", true, ch, nil, victim, types.TO_ROOM)

	case types.SPELL_TYPE_STAFF:
		// Staff affects everyone in the room except the caster
		for _, v := range ch.InRoom.Characters {
			if v != ch {
				// Check if already affected
				if w.AffectedBySpell(v, types.SPELL_ARMOR) {
					continue
				}

				// Create affect
				affect := &types.Affect{
					Type:      types.SPELL_ARMOR,
					Duration:  24,
					Modifier:  -20,
					Location:  types.APPLY_AC,
					Bitvector: 0,
				}

				// Add affect to victim
				w.AffectToChar(v, affect)
				v.SendMessage("You feel protected.\r\n")
			}
		}
	}

	return nil
}

// CastBless casts the bless spell
func (w *World) CastBless(level int, ch *types.Character, arg string, spellType int, victim *types.Character, obj *types.ObjectInstance) error {
	// Handle different spell types
	switch spellType {
	case types.SPELL_TYPE_SPELL:
		if obj != nil {
			// Bless an object
			if obj.Prototype.ExtraFlags&types.ITEM_BLESS != 0 {
				ch.SendMessage("That item is already blessed.\r\n")
				return nil
			}

			// Add bless flag to object
			obj.Prototype.ExtraFlags |= types.ITEM_BLESS
			w.Act("$p glows with a holy aura.", false, ch, obj, nil, types.TO_CHAR)
			w.Act("$p glows with a holy aura.", true, ch, obj, nil, types.TO_ROOM)
			return nil
		}

		// Bless a character
		if victim == nil {
			victim = ch
		}

		// Check if already affected
		if w.AffectedBySpell(victim, types.SPELL_BLESS) {
			if victim == ch {
				ch.SendMessage("You are already blessed.\r\n")
			} else {
				ch.SendMessage("They are already blessed.\r\n")
			}
			return nil
		}

		// Create affect
		affect := &types.Affect{
			Type:      types.SPELL_BLESS,
			Duration:  6,
			Modifier:  1,
			Location:  types.APPLY_HITROLL,
			Bitvector: 0,
		}

		// Add affect to victim
		w.AffectToChar(victim, affect)

		// Create second affect for saving throw
		affect2 := &types.Affect{
			Type:      types.SPELL_BLESS,
			Duration:  6,
			Modifier:  -1,
			Location:  types.APPLY_SAVING_SPELL,
			Bitvector: 0,
		}

		// Add second affect to victim
		w.AffectToChar(victim, affect2)

		// Send messages
		if victim == ch {
			ch.SendMessage("You feel blessed.\r\n")
		} else {
			victim.SendMessage("You feel blessed.\r\n")
			w.Act("You bless $N.", false, ch, nil, victim, types.TO_CHAR)
		}

	case types.SPELL_TYPE_POTION:
		// Check if already affected
		if w.AffectedBySpell(ch, types.SPELL_BLESS) {
			return nil
		}

		// Create affect
		affect := &types.Affect{
			Type:      types.SPELL_BLESS,
			Duration:  6,
			Modifier:  1,
			Location:  types.APPLY_HITROLL,
			Bitvector: 0,
		}

		// Add affect to drinker
		w.AffectToChar(ch, affect)

		// Create second affect for saving throw
		affect2 := &types.Affect{
			Type:      types.SPELL_BLESS,
			Duration:  6,
			Modifier:  -1,
			Location:  types.APPLY_SAVING_SPELL,
			Bitvector: 0,
		}

		// Add second affect to drinker
		w.AffectToChar(ch, affect2)
		ch.SendMessage("You feel blessed.\r\n")

	case types.SPELL_TYPE_SCROLL:
		if obj != nil {
			// Bless an object
			if obj.Prototype.ExtraFlags&types.ITEM_BLESS != 0 {
				return nil
			}

			// Add bless flag to object
			obj.Prototype.ExtraFlags |= types.ITEM_BLESS
			w.Act("$p glows with a holy aura.", false, ch, obj, nil, types.TO_CHAR)
			w.Act("$p glows with a holy aura.", true, ch, obj, nil, types.TO_ROOM)
			return nil
		}

		// If no target specified, target self
		if victim == nil {
			victim = ch
		}

		// Check if already affected
		if w.AffectedBySpell(victim, types.SPELL_BLESS) {
			return nil
		}

		// Create affect
		affect := &types.Affect{
			Type:      types.SPELL_BLESS,
			Duration:  6,
			Modifier:  1,
			Location:  types.APPLY_HITROLL,
			Bitvector: 0,
		}

		// Add affect to victim
		w.AffectToChar(victim, affect)

		// Create second affect for saving throw
		affect2 := &types.Affect{
			Type:      types.SPELL_BLESS,
			Duration:  6,
			Modifier:  -1,
			Location:  types.APPLY_SAVING_SPELL,
			Bitvector: 0,
		}

		// Add second affect to victim
		w.AffectToChar(victim, affect2)
		victim.SendMessage("You feel blessed.\r\n")
	}

	return nil
}

// CastBlindness casts the blindness spell
func (w *World) CastBlindness(level int, ch *types.Character, arg string, spellType int, victim *types.Character, obj *types.ObjectInstance) error {
	if victim == nil {
		return fmt.Errorf("you need to specify a target")
	}

	// Check if already affected
	if w.AffectedBySpell(victim, types.SPELL_BLINDNESS) {
		ch.SendMessage("They are already blind.\r\n")
		return nil
	}

	// Check if victim saves against spell
	if w.SavesSpell(victim, types.SAVING_SPELL) {
		ch.SendMessage("They resist your blindness spell.\r\n")
		victim.SendMessage("You feel a slight tingling around your eyes, but it passes.\r\n")
		return nil
	}

	// Create affect
	affect := &types.Affect{
		Type:      types.SPELL_BLINDNESS,
		Duration:  level,
		Modifier:  -4,
		Location:  types.APPLY_HITROLL,
		Bitvector: types.AFF_BLIND,
	}

	// Add affect to victim
	w.AffectToChar(victim, affect)

	// Send messages
	victim.SendMessage("You have been blinded!\r\n")
	w.Act("$N appears to be blinded!", true, ch, nil, victim, types.TO_CHAR)
	w.Act("$N appears to be blinded!", true, ch, nil, victim, types.TO_NOTVICT)

	return nil
}

// CastCureBlindness casts the cure blindness spell
func (w *World) CastCureBlindness(level int, ch *types.Character, arg string, spellType int, victim *types.Character, obj *types.ObjectInstance) error {
	if victim == nil {
		victim = ch
	}

	// Check if affected by blindness
	if !w.AffectedBySpell(victim, types.SPELL_BLINDNESS) {
		if victim == ch {
			ch.SendMessage("You aren't blind.\r\n")
		} else {
			ch.SendMessage("They aren't blind.\r\n")
		}
		return nil
	}

	// Remove blindness
	w.AffectFromChar(victim, types.SPELL_BLINDNESS)

	// Send messages
	victim.SendMessage("Your vision returns!\r\n")
	if victim != ch {
		w.Act("You cure $N of blindness.", false, ch, nil, victim, types.TO_CHAR)
	}
	w.Act("$N's vision returns.", true, ch, nil, victim, types.TO_NOTVICT)

	return nil
}

// CastCureLight casts the cure light wounds spell
func (w *World) CastCureLight(level int, ch *types.Character, arg string, spellType int, victim *types.Character, obj *types.ObjectInstance) error {
	if victim == nil {
		victim = ch
	}

	// Calculate healing amount (1d8 + level/4)
	heal := w.Dice(1, 8) + (level / 4)

	// Apply healing
	victim.HP += heal
	if victim.HP > victim.MaxHitPoints {
		victim.HP = victim.MaxHitPoints
	}

	// Send messages
	if victim == ch {
		ch.SendMessage("You feel better!\r\n")
	} else {
		victim.SendMessage("You feel better!\r\n")
		w.Act("You heal $N.", false, ch, nil, victim, types.TO_CHAR)
	}
	w.Act("$N looks better.", true, ch, nil, victim, types.TO_NOTVICT)

	return nil
}

// CastCureCritic casts the cure critical wounds spell
func (w *World) CastCureCritic(level int, ch *types.Character, arg string, spellType int, victim *types.Character, obj *types.ObjectInstance) error {
	if victim == nil {
		victim = ch
	}

	// Calculate healing amount (3d8 + level/2)
	heal := w.Dice(3, 8) + (level / 2)

	// Apply healing
	victim.HP += heal
	if victim.HP > victim.MaxHitPoints {
		victim.HP = victim.MaxHitPoints
	}

	// Send messages
	if victim == ch {
		ch.SendMessage("You feel much better!\r\n")
	} else {
		victim.SendMessage("You feel much better!\r\n")
		w.Act("You heal $N.", false, ch, nil, victim, types.TO_CHAR)
	}
	w.Act("$N looks much better.", true, ch, nil, victim, types.TO_NOTVICT)

	return nil
}

// CastHeal casts the heal spell
func (w *World) CastHeal(level int, ch *types.Character, arg string, spellType int, victim *types.Character, obj *types.ObjectInstance) error {
	switch spellType {
	case types.SPELL_TYPE_SPELL, types.SPELL_TYPE_POTION, types.SPELL_TYPE_SCROLL, types.SPELL_TYPE_WAND:
		if victim == nil {
			victim = ch
		}

		// Cure blindness first
		if w.AffectedBySpell(victim, types.SPELL_BLINDNESS) {
			w.AffectFromChar(victim, types.SPELL_BLINDNESS)
			victim.SendMessage("Your vision returns!\r\n")
		}

		// Calculate healing amount (100 + 1d8 for randomness)
		heal := 100 + w.Dice(1, 8)

		// Apply healing
		victim.HP += heal
		if victim.HP > victim.MaxHitPoints {
			victim.HP = victim.MaxHitPoints
		}

		// Send messages
		if victim == ch {
			ch.SendMessage("A warm feeling fills your body!\r\n")
		} else {
			victim.SendMessage("A warm feeling fills your body!\r\n")
			if spellType == types.SPELL_TYPE_SPELL {
				w.Act("You heal $N.", false, ch, nil, victim, types.TO_CHAR)
			}
		}
		w.Act("$N looks completely healed.", true, ch, nil, victim, types.TO_NOTVICT)

	case types.SPELL_TYPE_STAFF:
		// Staff affects everyone in the room except the caster
		for _, v := range ch.InRoom.Characters {
			if v != ch {
				// Cure blindness first
				if w.AffectedBySpell(v, types.SPELL_BLINDNESS) {
					w.AffectFromChar(v, types.SPELL_BLINDNESS)
					v.SendMessage("Your vision returns!\r\n")
				}

				// Calculate healing amount (100 + 1d8 for randomness)
				heal := 100 + w.Dice(1, 8)

				// Apply healing
				v.HP += heal
				if v.HP > v.MaxHitPoints {
					v.HP = v.MaxHitPoints
				}

				// Send messages
				v.SendMessage("A warm feeling fills your body!\r\n")
				w.Act("$N looks completely healed.", true, ch, nil, v, types.TO_NOTVICT)
			}
		}
	}

	return nil
}

// CastProtectionFromEvil casts the protection from evil spell
func (w *World) CastProtectionFromEvil(level int, ch *types.Character, arg string, spellType int, victim *types.Character, obj *types.ObjectInstance) error {
	if victim == nil {
		victim = ch
	}

	// Check if already affected
	if w.AffectedBySpell(victim, types.SPELL_PROTECTION_FROM_EVIL) {
		if victim == ch {
			ch.SendMessage("You are already protected from evil.\r\n")
		} else {
			ch.SendMessage("They are already protected from evil.\r\n")
		}
		return nil
	}

	// Create affect
	affect := &types.Affect{
		Type:      types.SPELL_PROTECTION_FROM_EVIL,
		Duration:  24,
		Modifier:  0,
		Location:  types.APPLY_NONE,
		Bitvector: types.AFF_PROTECT_EVIL,
	}

	// Add affect to victim
	w.AffectToChar(victim, affect)

	// Send messages
	if victim == ch {
		ch.SendMessage("You feel protected from evil.\r\n")
	} else {
		victim.SendMessage("You feel protected from evil.\r\n")
		w.Act("$N is surrounded by a pale blue aura.", false, ch, nil, victim, types.TO_CHAR)
	}
	w.Act("$n is surrounded by a pale blue aura.", true, victim, nil, nil, types.TO_ROOM)

	return nil
}

// CastRemoveCurse casts the remove curse spell
func (w *World) CastRemoveCurse(level int, ch *types.Character, arg string, spellType int, victim *types.Character, obj *types.ObjectInstance) error {
	// Handle different spell types
	switch spellType {
	case types.SPELL_TYPE_SPELL:
		if obj != nil {
			// Remove curse from an object
			if obj.Prototype.ExtraFlags&types.ITEM_NODROP == 0 {
				ch.SendMessage("That item is not cursed.\r\n")
				return nil
			}

			// Remove cursed flag
			obj.Prototype.ExtraFlags &= ^uint32(types.ITEM_NODROP)
			w.Act("$p briefly glows blue.", false, ch, obj, nil, types.TO_CHAR)
			w.Act("$p briefly glows blue.", true, ch, obj, nil, types.TO_ROOM)
			return nil
		}

		// Remove curse from a character
		if victim == nil {
			victim = ch
		}

		// Check if affected by curse
		if !w.AffectedBySpell(victim, types.SPELL_CURSE) {
			if victim == ch {
				ch.SendMessage("You aren't cursed.\r\n")
			} else {
				ch.SendMessage("They aren't cursed.\r\n")
			}
			return nil
		}

		// Remove curse
		w.AffectFromChar(victim, types.SPELL_CURSE)

		// Send messages
		victim.SendMessage("You feel better.\r\n")
		if victim != ch {
			w.Act("You remove the curse from $N.", false, ch, nil, victim, types.TO_CHAR)
		}

	case types.SPELL_TYPE_POTION:
		// Check if affected by curse
		if !w.AffectedBySpell(ch, types.SPELL_CURSE) {
			return nil
		}

		// Remove curse
		w.AffectFromChar(ch, types.SPELL_CURSE)
		ch.SendMessage("You feel better.\r\n")

	case types.SPELL_TYPE_SCROLL:
		if obj != nil {
			// Remove curse from an object
			if obj.Prototype.ExtraFlags&types.ITEM_NODROP == 0 {
				return nil
			}

			// Remove cursed flag
			obj.Prototype.ExtraFlags &= ^uint32(types.ITEM_NODROP)
			w.Act("$p briefly glows blue.", false, ch, obj, nil, types.TO_CHAR)
			w.Act("$p briefly glows blue.", true, ch, obj, nil, types.TO_ROOM)
			return nil
		}

		// If no target specified, target self
		if victim == nil {
			victim = ch
		}

		// Check if affected by curse
		if !w.AffectedBySpell(victim, types.SPELL_CURSE) {
			return nil
		}

		// Remove curse
		w.AffectFromChar(victim, types.SPELL_CURSE)
		victim.SendMessage("You feel better.\r\n")
	}

	return nil
}

// CastCurse casts the curse spell
func (w *World) CastCurse(level int, ch *types.Character, arg string, spellType int, victim *types.Character, obj *types.ObjectInstance) error {
	// Handle different spell types
	switch spellType {
	case types.SPELL_TYPE_SPELL:
		if obj != nil {
			// Curse an object
			if obj.Prototype.ExtraFlags&types.ITEM_NODROP != 0 {
				ch.SendMessage("That item is already cursed.\r\n")
				return nil
			}

			// Add cursed flag
			obj.Prototype.ExtraFlags |= types.ITEM_NODROP
			w.Act("$p glows red.", false, ch, obj, nil, types.TO_CHAR)
			w.Act("$p glows red.", true, ch, obj, nil, types.TO_ROOM)
			return nil
		}

		// Curse a character
		if victim == nil {
			return fmt.Errorf("whom do you wish to curse?")
		}

		// Check if already affected
		if w.AffectedBySpell(victim, types.SPELL_CURSE) {
			ch.SendMessage("They are already cursed.\r\n")
			return nil
		}

		// Check if victim saves against spell
		if w.SavesSpell(victim, types.SAVING_SPELL) {
			ch.SendMessage("They resist your curse spell.\r\n")
			victim.SendMessage("You feel a brief chill, but it passes.\r\n")
			return nil
		}

		// Create affect for hitroll
		affect1 := &types.Affect{
			Type:      types.SPELL_CURSE,
			Duration:  level / 2,
			Modifier:  -1,
			Location:  types.APPLY_HITROLL,
			Bitvector: types.AFF_CURSE,
		}

		// Add affect to victim
		w.AffectToChar(victim, affect1)

		// Create affect for saving throw
		affect2 := &types.Affect{
			Type:      types.SPELL_CURSE,
			Duration:  level / 2,
			Modifier:  1,
			Location:  types.APPLY_SAVING_SPELL,
			Bitvector: types.AFF_CURSE,
		}

		// Add affect to victim
		w.AffectToChar(victim, affect2)

		// Send messages
		victim.SendMessage("You feel very uncomfortable.\r\n")
		w.Act("$n briefly glows red!", true, victim, nil, nil, types.TO_ROOM)

	case types.SPELL_TYPE_POTION:
		// Check if already affected
		if w.AffectedBySpell(ch, types.SPELL_CURSE) {
			return nil
		}

		// Create affect for hitroll
		affect1 := &types.Affect{
			Type:      types.SPELL_CURSE,
			Duration:  level / 2,
			Modifier:  -1,
			Location:  types.APPLY_HITROLL,
			Bitvector: types.AFF_CURSE,
		}

		// Add affect to drinker
		w.AffectToChar(ch, affect1)

		// Create affect for saving throw
		affect2 := &types.Affect{
			Type:      types.SPELL_CURSE,
			Duration:  level / 2,
			Modifier:  1,
			Location:  types.APPLY_SAVING_SPELL,
			Bitvector: types.AFF_CURSE,
		}

		// Add affect to drinker
		w.AffectToChar(ch, affect2)
		ch.SendMessage("You feel very uncomfortable.\r\n")

	case types.SPELL_TYPE_SCROLL:
		if obj != nil {
			// Curse an object
			if obj.Prototype.ExtraFlags&types.ITEM_NODROP != 0 {
				return nil
			}

			// Add cursed flag
			obj.Prototype.ExtraFlags |= types.ITEM_NODROP
			w.Act("$p glows red.", false, ch, obj, nil, types.TO_CHAR)
			w.Act("$p glows red.", true, ch, obj, nil, types.TO_ROOM)
			return nil
		}

		// If no target specified, target self
		if victim == nil {
			victim = ch
		}

		// Check if already affected
		if w.AffectedBySpell(victim, types.SPELL_CURSE) {
			return nil
		}

		// Create affect for hitroll
		affect1 := &types.Affect{
			Type:      types.SPELL_CURSE,
			Duration:  level / 2,
			Modifier:  -1,
			Location:  types.APPLY_HITROLL,
			Bitvector: types.AFF_CURSE,
		}

		// Add affect to victim
		w.AffectToChar(victim, affect1)

		// Create affect for saving throw
		affect2 := &types.Affect{
			Type:      types.SPELL_CURSE,
			Duration:  level / 2,
			Modifier:  1,
			Location:  types.APPLY_SAVING_SPELL,
			Bitvector: types.AFF_CURSE,
		}

		// Add affect to victim
		w.AffectToChar(victim, affect2)
		victim.SendMessage("You feel very uncomfortable.\r\n")
	}

	return nil
}

// Dice rolls dice
func (w *World) Dice(num, size int) int {
	result := 0
	for i := 0; i < num; i++ {
		result += (w.Random(size) + 1)
	}
	return result
}

// Random returns a random number between 0 and n-1
func (w *World) Random(n int) int {
	if n <= 0 {
		return 0
	}
	return w.rand.Intn(n)
}
