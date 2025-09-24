package world

import (
	"fmt"
	"strings"

	"github.com/wltechblog/DikuGo/pkg/types"
)

// CastCreateFood casts the create food spell
func (w *World) CastCreateFood(level int, ch *types.Character, arg string, spellType int, victim *types.Character, obj *types.ObjectInstance) error {
	// Create a food object
	foodVnum := 19 // Waybread VNUM
	food := w.CreateObject(foodVnum)
	if food == nil {
		ch.SendMessage("You failed to create food.\r\n")
		return nil
	}

	// Add food to character's inventory
	w.ObjectToChar(food, ch)

	// Send messages
	w.Act("$n creates $p.", true, ch, food, nil, types.TO_ROOM)
	w.Act("You create $p.", false, ch, food, nil, types.TO_CHAR)

	return nil
}

// CastCreateWater casts the create water spell
func (w *World) CastCreateWater(level int, ch *types.Character, arg string, spellType int, victim *types.Character, obj *types.ObjectInstance) error {
	if obj == nil {
		return fmt.Errorf("you need to specify a drink container")
	}

	// Check if object is a drink container
	if obj.Prototype.Type != types.ITEM_DRINKCON {
		ch.SendMessage("You can only fill drink containers.\r\n")
		return nil
	}

	// Check if container is full
	if obj.Prototype.Value[1] >= obj.Prototype.Value[0] {
		ch.SendMessage("The container is already full.\r\n")
		return nil
	}

	// Check if container contains something other than water
	if obj.Prototype.Value[2] != 0 && obj.Prototype.Value[2] != types.LIQ_WATER {
		ch.SendMessage("The container contains some other liquid.\r\n")
		return nil
	}

	// Calculate water created (level * 2)
	water := level * 2

	// Add water to container
	obj.Prototype.Value[2] = types.LIQ_WATER // Set liquid type to water
	obj.Prototype.Value[1] += water          // Add water

	// Don't overfill
	if obj.Prototype.Value[1] > obj.Prototype.Value[0] {
		obj.Prototype.Value[1] = obj.Prototype.Value[0]
	}

	// Send messages
	w.Act("$p is filled with water.", false, ch, obj, nil, types.TO_CHAR)
	w.Act("$p is filled with water.", true, ch, obj, nil, types.TO_ROOM)

	return nil
}

// CastDetectEvil casts the detect evil spell
func (w *World) CastDetectEvil(level int, ch *types.Character, arg string, spellType int, victim *types.Character, obj *types.ObjectInstance) error {
	if victim == nil {
		victim = ch
	}

	// Check if already affected
	if w.AffectedBySpell(victim, types.SPELL_DETECT_EVIL) {
		if victim == ch {
			ch.SendMessage("You can already detect evil.\r\n")
		} else {
			ch.SendMessage("They can already detect evil.\r\n")
		}
		return nil
	}

	// Create affect
	affect := &types.Affect{
		Type:      types.SPELL_DETECT_EVIL,
		Duration:  level * 5,
		Modifier:  0,
		Location:  types.APPLY_NONE,
		Bitvector: types.AFF_DETECT_EVIL,
	}

	// Add affect to victim
	w.AffectToChar(victim, affect)

	// Send messages
	if victim == ch {
		ch.SendMessage("Your eyes tingle.\r\n")
	} else {
		victim.SendMessage("Your eyes tingle.\r\n")
		w.Act("You grant $N the ability to detect evil.", false, ch, nil, victim, types.TO_CHAR)
	}

	return nil
}

// CastDetectInvisible casts the detect invisible spell
func (w *World) CastDetectInvisible(level int, ch *types.Character, arg string, spellType int, victim *types.Character, obj *types.ObjectInstance) error {
	if victim == nil {
		victim = ch
	}

	// Check if already affected
	if w.AffectedBySpell(victim, types.SPELL_DETECT_INVISIBLE) {
		if victim == ch {
			ch.SendMessage("You can already detect invisible.\r\n")
		} else {
			ch.SendMessage("They can already detect invisible.\r\n")
		}
		return nil
	}

	// Create affect
	affect := &types.Affect{
		Type:      types.SPELL_DETECT_INVISIBLE,
		Duration:  level * 5,
		Modifier:  0,
		Location:  types.APPLY_NONE,
		Bitvector: types.AFF_DETECT_INVISIBLE,
	}

	// Add affect to victim
	w.AffectToChar(victim, affect)

	// Send messages
	if victim == ch {
		ch.SendMessage("Your eyes tingle.\r\n")
	} else {
		victim.SendMessage("Your eyes tingle.\r\n")
		w.Act("You grant $N the ability to detect invisible.", false, ch, nil, victim, types.TO_CHAR)
	}

	return nil
}

// CastDetectMagic casts the detect magic spell
func (w *World) CastDetectMagic(level int, ch *types.Character, arg string, spellType int, victim *types.Character, obj *types.ObjectInstance) error {
	if victim == nil {
		victim = ch
	}

	// Check if already affected
	if w.AffectedBySpell(victim, types.SPELL_DETECT_MAGIC) {
		if victim == ch {
			ch.SendMessage("You can already detect magic.\r\n")
		} else {
			ch.SendMessage("They can already detect magic.\r\n")
		}
		return nil
	}

	// Create affect
	affect := &types.Affect{
		Type:      types.SPELL_DETECT_MAGIC,
		Duration:  level * 5,
		Modifier:  0,
		Location:  types.APPLY_NONE,
		Bitvector: types.AFF_DETECT_MAGIC,
	}

	// Add affect to victim
	w.AffectToChar(victim, affect)

	// Send messages
	if victim == ch {
		ch.SendMessage("Your eyes tingle.\r\n")
	} else {
		victim.SendMessage("Your eyes tingle.\r\n")
		w.Act("You grant $N the ability to detect magic.", false, ch, nil, victim, types.TO_CHAR)
	}

	return nil
}

// CastDetectPoison casts the detect poison spell
func (w *World) CastDetectPoison(level int, ch *types.Character, arg string, spellType int, victim *types.Character, obj *types.ObjectInstance) error {
	if obj != nil {
		// Check if object is poisoned
		isPoisoned := false
		if obj.Prototype.Type == types.ITEM_DRINKCON || obj.Prototype.Type == types.ITEM_FOOD {
			if obj.Prototype.Value[3] != 0 {
				isPoisoned = true
			}
		}

		// Send messages
		if isPoisoned {
			w.Act("You detect poison in $p.", false, ch, obj, nil, types.TO_CHAR)
		} else {
			w.Act("You detect no poison in $p.", false, ch, obj, nil, types.TO_CHAR)
		}
		return nil
	}

	if victim == nil {
		victim = ch
	}

	// Check if victim is poisoned
	isPoisoned := w.AffectedBySpell(victim, types.SPELL_POISON)

	// Send messages
	if victim == ch {
		if isPoisoned {
			ch.SendMessage("You detect poison in your blood!\r\n")
		} else {
			ch.SendMessage("You detect no poison in your blood.\r\n")
		}
	} else {
		if isPoisoned {
			w.Act("You detect poison in $N's blood!", false, ch, nil, victim, types.TO_CHAR)
		} else {
			w.Act("You detect no poison in $N's blood.", false, ch, nil, victim, types.TO_CHAR)
		}
	}

	return nil
}

// CastInvisible casts the invisible spell
func (w *World) CastInvisible(level int, ch *types.Character, arg string, spellType int, victim *types.Character, obj *types.ObjectInstance) error {
	if obj != nil {
		// Make object invisible
		if obj.Prototype.ExtraFlags&types.ITEM_INVISIBLE != 0 {
			ch.SendMessage("That object is already invisible.\r\n")
			return nil
		}

		// Add invisible flag to object
		obj.Prototype.ExtraFlags |= types.ITEM_INVISIBLE

		// Send messages
		w.Act("$p vanishes.", false, ch, obj, nil, types.TO_CHAR)
		w.Act("$p vanishes.", true, ch, obj, nil, types.TO_ROOM)
		return nil
	}

	if victim == nil {
		victim = ch
	}

	// Check if already affected
	if w.AffectedBySpell(victim, types.SPELL_INVISIBLE) {
		if victim == ch {
			ch.SendMessage("You are already invisible.\r\n")
		} else {
			ch.SendMessage("They are already invisible.\r\n")
		}
		return nil
	}

	// Create affect
	affect := &types.Affect{
		Type:      types.SPELL_INVISIBLE,
		Duration:  24,
		Modifier:  0,
		Location:  types.APPLY_NONE,
		Bitvector: types.AFF_INVISIBLE,
	}

	// Add affect to victim
	w.AffectToChar(victim, affect)

	// Send messages
	if victim == ch {
		ch.SendMessage("You vanish.\r\n")
	} else {
		victim.SendMessage("You vanish.\r\n")
		w.Act("$N vanishes.", false, ch, nil, victim, types.TO_CHAR)
	}
	w.Act("$n slowly fades out of existence.", true, victim, nil, nil, types.TO_ROOM)

	return nil
}

// CastPoison casts the poison spell
func (w *World) CastPoison(level int, ch *types.Character, arg string, spellType int, victim *types.Character, obj *types.ObjectInstance) error {
	if obj != nil {
		// Poison object
		if obj.Prototype.Type == types.ITEM_DRINKCON || obj.Prototype.Type == types.ITEM_FOOD {
			if obj.Prototype.Value[3] != 0 {
				ch.SendMessage("That is already poisoned.\r\n")
				return nil
			}
			obj.Prototype.Value[3] = 1 // Set poisoned flag
			w.Act("$p is now poisoned.", false, ch, obj, nil, types.TO_CHAR)
			return nil
		} else {
			ch.SendMessage("You can only poison food or drink.\r\n")
			return nil
		}
	}

	if victim == nil {
		return fmt.Errorf("you need to specify a target")
	}

	// Check if already affected
	if w.AffectedBySpell(victim, types.SPELL_POISON) {
		ch.SendMessage("They are already poisoned.\r\n")
		return nil
	}

	// Check if victim saves against spell
	if w.SavesSpell(victim, types.SAVING_PARA) {
		ch.SendMessage("They resist your poison spell.\r\n")
		victim.SendMessage("You feel a slight tingling, but it passes.\r\n")
		return nil
	}

	// Create affect
	affect := &types.Affect{
		Type:      types.SPELL_POISON,
		Duration:  level * 2,
		Modifier:  -2,
		Location:  types.APPLY_STR,
		Bitvector: types.AFF_POISON,
	}

	// Add affect to victim
	w.AffectJoin(victim, affect, false, false)

	// Send messages
	victim.SendMessage("You feel very sick.\r\n")
	w.Act("$N looks very ill.", true, ch, nil, victim, types.TO_CHAR)

	return nil
}

// CastRemovePoison casts the remove poison spell
func (w *World) CastRemovePoison(level int, ch *types.Character, arg string, spellType int, victim *types.Character, obj *types.ObjectInstance) error {
	if obj != nil {
		// Remove poison from object
		if obj.Prototype.Type == types.ITEM_DRINKCON || obj.Prototype.Type == types.ITEM_FOOD {
			if obj.Prototype.Value[3] == 0 {
				ch.SendMessage("That isn't poisoned.\r\n")
				return nil
			}
			obj.Prototype.Value[3] = 0 // Remove poisoned flag
			w.Act("$p is no longer poisoned.", false, ch, obj, nil, types.TO_CHAR)
			return nil
		} else {
			ch.SendMessage("You can only remove poison from food or drink.\r\n")
			return nil
		}
	}

	if victim == nil {
		victim = ch
	}

	// Check if affected by poison
	if !w.AffectedBySpell(victim, types.SPELL_POISON) {
		if victim == ch {
			ch.SendMessage("You aren't poisoned.\r\n")
		} else {
			ch.SendMessage("They aren't poisoned.\r\n")
		}
		return nil
	}

	// Remove poison
	w.AffectFromChar(victim, types.SPELL_POISON)

	// Send messages
	victim.SendMessage("A warm feeling runs through your body!\r\n")
	if victim != ch {
		w.Act("You cure $N of poison.", false, ch, nil, victim, types.TO_CHAR)
	}
	w.Act("$N looks better.", true, ch, nil, victim, types.TO_NOTVICT)

	return nil
}

// CastSanctuary casts the sanctuary spell
func (w *World) CastSanctuary(level int, ch *types.Character, arg string, spellType int, victim *types.Character, obj *types.ObjectInstance) error {
	if victim == nil {
		victim = ch
	}

	// Check if already affected
	if w.AffectedBySpell(victim, types.SPELL_SANCTUARY) {
		if victim == ch {
			ch.SendMessage("You are already in sanctuary.\r\n")
		} else {
			ch.SendMessage("They are already in sanctuary.\r\n")
		}
		return nil
	}

	// Create affect
	affect := &types.Affect{
		Type:      types.SPELL_SANCTUARY,
		Duration:  3 + (level / 6),
		Modifier:  0,
		Location:  types.APPLY_NONE,
		Bitvector: types.AFF_SANCTUARY,
	}

	// Add affect to victim
	w.AffectToChar(victim, affect)

	// Send messages
	if victim == ch {
		ch.SendMessage("A white aura surrounds you.\r\n")
	} else {
		victim.SendMessage("A white aura surrounds you.\r\n")
		w.Act("$N is surrounded by a white aura.", false, ch, nil, victim, types.TO_CHAR)
	}
	w.Act("$n is surrounded by a white aura.", true, victim, nil, nil, types.TO_ROOM)

	return nil
}

// CastStrength casts the strength spell
func (w *World) CastStrength(level int, ch *types.Character, arg string, spellType int, victim *types.Character, obj *types.ObjectInstance) error {
	if victim == nil {
		victim = ch
	}

	// Check if already affected
	if w.AffectedBySpell(victim, types.SPELL_STRENGTH) {
		if victim == ch {
			ch.SendMessage("You are already strengthened.\r\n")
		} else {
			ch.SendMessage("They are already strengthened.\r\n")
		}
		return nil
	}

	// Create affect
	modifier := 1
	if level >= 18 {
		modifier = 2
	}

	affect := &types.Affect{
		Type:      types.SPELL_STRENGTH,
		Duration:  level,
		Modifier:  modifier,
		Location:  types.APPLY_STR,
		Bitvector: 0,
	}

	// Add affect to victim
	w.AffectToChar(victim, affect)

	// Send messages
	if victim == ch {
		ch.SendMessage("You feel stronger.\r\n")
	} else {
		victim.SendMessage("You feel stronger.\r\n")
		w.Act("$N looks stronger.", false, ch, nil, victim, types.TO_CHAR)
	}

	return nil
}

// CastSummon casts the summon spell
func (w *World) CastSummon(level int, ch *types.Character, arg string, spellType int, victim *types.Character, obj *types.ObjectInstance) error {
	// Summon requires a target name
	if arg == "" {
		return fmt.Errorf("who do you want to summon?")
	}

	// Find the victim in the world
	victim = w.FindCharacterInWorld(arg)
	if victim == nil {
		return fmt.Errorf("no one by that name in the world")
	}

	// Check if victim is self
	if victim == ch {
		return fmt.Errorf("you are already here")
	}

	// Check if victim is in same room
	if victim.InRoom == ch.InRoom {
		return fmt.Errorf("they are already here")
	}

	// Check if victim is an NPC
	if victim.IsNPC {
		return fmt.Errorf("you failed")
	}

	// Check if victim is fighting
	if victim.Fighting != nil {
		return fmt.Errorf("you cannot summon someone who is fighting")
	}

	// Check if victim's level is too high
	if victim.Level > ch.Level+3 {
		return fmt.Errorf("you failed")
	}

	// Check if victim's room has no-summon flag
	if victim.InRoom.Flags&types.ROOM_NOSUMMON != 0 {
		return fmt.Errorf("magical forces prevent your summon")
	}

	// Check if caster's room has no-summon flag
	if ch.InRoom.Flags&types.ROOM_NOSUMMON != 0 {
		return fmt.Errorf("magical forces prevent your summon")
	}

	// Move victim to caster's room
	w.Act("$n disappears suddenly.", true, victim, nil, nil, types.TO_ROOM)
	w.CharFromRoom(victim)
	w.CharToRoom(victim, ch.InRoom)
	w.Act("$n arrives suddenly.", true, victim, nil, nil, types.TO_ROOM)
	w.Act("$n has summoned you!", false, ch, nil, victim, types.TO_VICT)
	victim.SendMessage("You look around...\r\n")

	return nil
}

// FindCharacterInWorld finds a character in the world by name
func (w *World) FindCharacterInWorld(name string) *types.Character {
	name = strings.ToLower(name)
	for _, ch := range w.characters {
		if strings.Contains(strings.ToLower(ch.Name), name) {
			return ch
		}
	}
	return nil
}

// CastWordOfRecall casts the word of recall spell
func (w *World) CastWordOfRecall(level int, ch *types.Character, arg string, spellType int, victim *types.Character, obj *types.ObjectInstance) error {
	if victim == nil {
		victim = ch
	}

	// Get recall room
	recallRoom := w.GetRoom(3001) // Temple of Midgaard
	if recallRoom == nil {
		return fmt.Errorf("you failed")
	}

	// Send messages
	victim.SendMessage("You feel a white aura surround you.\r\n")
	w.Act("$n disappears.", true, victim, nil, nil, types.TO_ROOM)

	// Move victim to recall room
	w.CharFromRoom(victim)
	w.CharToRoom(victim, recallRoom)
	w.Act("$n appears in the middle of the room.", true, victim, nil, nil, types.TO_ROOM)
	victim.SendMessage("You look around...\r\n")

	return nil
}

// CastLocateObject casts the locate object spell
func (w *World) CastLocateObject(level int, ch *types.Character, arg string, spellType int, victim *types.Character, obj *types.ObjectInstance) error {
	if arg == "" {
		return fmt.Errorf("what do you want to locate?")
	}

	// Search for objects in the world
	count := 0
	arg = strings.ToLower(arg)

	// Search player inventories
	for _, player := range w.characters {
		if player.IsNPC {
			continue
		}

		// Check inventory
		for _, obj := range player.Inventory {
			if strings.Contains(strings.ToLower(obj.Prototype.Name), arg) {
				count++
				if count <= 15 { // Limit to 15 results
					ch.SendMessage(fmt.Sprintf("%s is carrying %s.\r\n", player.Name, obj.Prototype.ShortDesc))
				}
			}
		}

		// Check equipment
		for _, obj := range player.Equipment {
			if obj != nil && strings.Contains(strings.ToLower(obj.Prototype.Name), arg) {
				count++
				if count <= 15 { // Limit to 15 results
					ch.SendMessage(fmt.Sprintf("%s is using %s.\r\n", player.Name, obj.Prototype.ShortDesc))
				}
			}
		}
	}

	// Search rooms
	for _, room := range w.rooms {
		for _, obj := range room.Objects {
			if strings.Contains(strings.ToLower(obj.Prototype.Name), arg) {
				count++
				if count <= 15 { // Limit to 15 results
					ch.SendMessage(fmt.Sprintf("%s is in %s.\r\n", obj.Prototype.ShortDesc, room.Name))
				}
			}
		}
	}

	if count == 0 {
		ch.SendMessage("You cannot locate any such object.\r\n")
	} else if count > 15 {
		ch.SendMessage("There are too many objects to list them all.\r\n")
	}

	return nil
}

// CreateObject creates a new object instance from a prototype
func (w *World) CreateObject(vnum int) *types.ObjectInstance {
	// Find the prototype
	prototype := w.GetObjectPrototype(vnum)
	if prototype == nil {
		return nil
	}

	// Create a new instance
	obj := &types.ObjectInstance{
		Prototype: prototype,
		Timer:     -1, // Permanent
	}

	return obj
}

// CastSenseLife casts the sense life spell
func (w *World) CastSenseLife(level int, ch *types.Character, arg string, spellType int, victim *types.Character, obj *types.ObjectInstance) error {
	if victim == nil {
		victim = ch
	}

	// Check if already affected
	if w.AffectedBySpell(victim, types.SPELL_SENSE_LIFE) {
		if victim == ch {
			ch.SendMessage("You can already sense life.\r\n")
		} else {
			ch.SendMessage("They can already sense life.\r\n")
		}
		return nil
	}

	// Create affect
	affect := &types.Affect{
		Type:      types.SPELL_SENSE_LIFE,
		Duration:  level * 5,
		Modifier:  0,
		Location:  types.APPLY_NONE,
		Bitvector: types.AFF_SENSE_LIFE,
	}

	// Add affect to victim
	w.AffectToChar(victim, affect)

	// Send messages
	if victim == ch {
		ch.SendMessage("You feel your awareness improve.\r\n")
	} else {
		victim.SendMessage("You feel your awareness improve.\r\n")
		w.Act("You grant $N the ability to sense life.", false, ch, nil, victim, types.TO_CHAR)
	}

	return nil
}
