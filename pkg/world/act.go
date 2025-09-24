package world

import (
	"strings"

	"github.com/wltechblog/DikuGo/pkg/types"
)

// Act sends a message to characters in a room
// ch is the character performing the action
// msg is the message to send
// hide is whether to hide the action from others
// obj is an optional object involved in the action
// vict is an optional victim of the action
// type is the type of message (TO_CHAR, TO_ROOM, TO_VICT, TO_NOTVICT)
func (w *World) Act(msg string, hide bool, ch *types.Character, obj *types.ObjectInstance, vict *types.Character, msgType int) {
	// Check for invalid parameters
	if msg == "" || ch == nil {
		return
	}

	// Get the room
	room := ch.InRoom
	if room == nil {
		return
	}

	// Process the message
	processedMsg := processActMessage(msg, ch, obj, vict)

	// Send the message to the appropriate recipients
	switch msgType {
	case types.TO_CHAR:
		// Send to the character
		if ch != nil {
			ch.SendMessage(processedMsg)
		}
	case types.TO_ROOM:
		// Send to everyone in the room except the character
		for _, rch := range room.Characters {
			if rch != ch && (rch.Position > types.POS_SLEEPING || !hide) {
				rch.SendMessage(processedMsg)
			}
		}
	case types.TO_VICT:
		// Send to the victim
		if vict != nil {
			vict.SendMessage(processedMsg)
		}
	case types.TO_NOTVICT:
		// Send to everyone in the room except the character and victim
		for _, rch := range room.Characters {
			if rch != ch && rch != vict && (rch.Position > types.POS_SLEEPING || !hide) {
				rch.SendMessage(processedMsg)
			}
		}
	case types.TO_ALL:
		// Send to everyone in the room
		for _, rch := range room.Characters {
			if rch.Position > types.POS_SLEEPING || !hide {
				rch.SendMessage(processedMsg)
			}
		}
	}
}

// processActMessage replaces placeholders in the message with actual names
func processActMessage(msg string, ch *types.Character, obj *types.ObjectInstance, vict *types.Character) string {
	// Replace $n with the character's name
	msg = strings.ReplaceAll(msg, "$n", ch.Name)

	// Replace $N with the victim's name
	if vict != nil {
		msg = strings.ReplaceAll(msg, "$N", vict.Name)
	}

	// Replace $p with the object's name
	if obj != nil {
		msg = strings.ReplaceAll(msg, "$p", obj.Prototype.ShortDesc)
	}

	// Replace $m with him/her
	if ch != nil {
		var himHer string
		if ch.Sex == types.SEX_MALE {
			himHer = "him"
		} else if ch.Sex == types.SEX_FEMALE {
			himHer = "her"
		} else {
			himHer = "it"
		}
		msg = strings.ReplaceAll(msg, "$m", himHer)
	}

	// Replace $s with his/her
	if ch != nil {
		var hisHer string
		if ch.Sex == types.SEX_MALE {
			hisHer = "his"
		} else if ch.Sex == types.SEX_FEMALE {
			hisHer = "her"
		} else {
			hisHer = "its"
		}
		msg = strings.ReplaceAll(msg, "$s", hisHer)
	}

	// Replace $e with he/she
	if ch != nil {
		var heShe string
		if ch.Sex == types.SEX_MALE {
			heShe = "he"
		} else if ch.Sex == types.SEX_FEMALE {
			heShe = "she"
		} else {
			heShe = "it"
		}
		msg = strings.ReplaceAll(msg, "$e", heShe)
	}

	// Replace $T with the third argument (string)
	if vict != nil && vict.Name == "" {
		msg = strings.ReplaceAll(msg, "$T", vict.Description)
	}

	return msg + "\r\n"
}
