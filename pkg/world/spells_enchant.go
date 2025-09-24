package world

import (
	"fmt"

	"github.com/wltechblog/DikuGo/pkg/types"
)

// CastEnchantWeapon casts the enchant weapon spell
func (w *World) CastEnchantWeapon(level int, ch *types.Character, arg string, spellType int, victim *types.Character, obj *types.ObjectInstance) error {
	// Check if object is provided
	if obj == nil {
		return fmt.Errorf("what item do you want to enchant?")
	}

	// Check if object is a weapon
	if obj.Prototype.Type != types.ITEM_WEAPON {
		ch.SendMessage("That is not a weapon.\r\n")
		return nil
	}

	// Check if weapon is already enchanted
	if obj.Prototype.ExtraFlags&types.ITEM_MAGIC != 0 {
		ch.SendMessage("That weapon is already enchanted.\r\n")
		return nil
	}

	// Calculate chance of success
	chance := 50 + (level / 2)
	if chance > 95 {
		chance = 95
	}

	// Roll for success
	roll := w.Random(100)
	if roll > chance {
		// Spell failed
		ch.SendMessage("Your enchantment failed!\r\n")

		// Check for catastrophic failure
		if roll > 95 {
			// Weapon is destroyed
			w.Act("$p glows brightly and explodes!", false, ch, obj, nil, types.TO_CHAR)
			w.Act("$p glows brightly and explodes!", true, ch, obj, nil, types.TO_ROOM)
			w.ExtractObj(obj)
		}

		return nil
	}

	// Add magic flag
	obj.Prototype.ExtraFlags |= types.ITEM_MAGIC

	// Improve weapon
	damBonus := w.Random(2) + 1 // +1 or +2 to damage
	hitBonus := w.Random(2) + 1 // +1 or +2 to hit

	// Add affects
	// Find an empty affect slot for damage
	damSlotFound := false
	for i := 0; i < types.MAX_OBJ_AFFECT; i++ {
		if obj.Affects[i].Location == 0 && obj.Affects[i].Modifier == 0 {
			// Empty slot found
			obj.Affects[i].Location = types.APPLY_DAMROLL
			obj.Affects[i].Modifier = damBonus
			damSlotFound = true
			break
		}
	}

	// Find an empty affect slot for hit
	hitSlotFound := false
	for i := 0; i < types.MAX_OBJ_AFFECT; i++ {
		if obj.Affects[i].Location == 0 && obj.Affects[i].Modifier == 0 {
			// Empty slot found
			obj.Affects[i].Location = types.APPLY_HITROLL
			obj.Affects[i].Modifier = hitBonus
			hitSlotFound = true
			break
		}
	}

	// Check if we couldn't add both affects
	if !damSlotFound || !hitSlotFound {
		ch.SendMessage("The weapon already has too many magical properties.\r\n")
	}

	// Send messages
	w.Act("$p glows blue.", false, ch, obj, nil, types.TO_CHAR)
	w.Act("$p glows blue.", true, ch, obj, nil, types.TO_ROOM)

	return nil
}

// CastControlWeather casts the control weather spell
func (w *World) CastControlWeather(level int, ch *types.Character, arg string, spellType int, victim *types.Character, obj *types.ObjectInstance) error {
	// Check if arguments are provided
	if arg == "" {
		return fmt.Errorf("do you want the weather to get better or worse?")
	}

	// Check if character is outdoors
	if ch.InRoom.Flags&types.ROOM_INDOORS != 0 {
		ch.SendMessage("You need to be outdoors to control the weather.\r\n")
		return nil
	}

	// Determine direction of change
	var change int
	if arg == "better" {
		change = -1
	} else if arg == "worse" {
		change = 1
	} else {
		return fmt.Errorf("do you want the weather to get better or worse?")
	}

	// Change the weather
	w.time.Weather += change

	// Ensure weather stays within bounds
	if w.time.Weather < types.SKY_CLOUDLESS {
		w.time.Weather = types.SKY_CLOUDLESS
	} else if w.time.Weather > types.SKY_LIGHTNING {
		w.time.Weather = types.SKY_LIGHTNING
	}

	// Send messages based on new weather
	switch w.time.Weather {
	case types.SKY_CLOUDLESS:
		w.SendToAll("The clouds disappear.\r\n")
	case types.SKY_CLOUDY:
		w.SendToAll("The sky is getting cloudy.\r\n")
	case types.SKY_RAINY:
		w.SendToAll("It starts to rain.\r\n")
	case types.SKY_LIGHTNING:
		w.SendToAll("Lightning flashes in the sky.\r\n")
	}

	return nil
}

// SendToAll sends a message to all players in the game
func (w *World) SendToAll(message string) {
	for _, ch := range w.characters {
		if !ch.IsNPC {
			ch.SendMessage(message)
		}
	}
}
