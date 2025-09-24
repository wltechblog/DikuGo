package world

import (
	"github.com/wltechblog/DikuGo/pkg/types"
)

// ApplyObjectAffects applies or removes the magical effects of an object to a character
func (w *World) ApplyObjectAffects(ch *types.Character, obj *types.ObjectInstance, add bool) {
	if ch == nil || obj == nil {
		return
	}

	// Apply or remove bitvector flags from the object
	if add {
		ch.AffectedBy |= int64(obj.Prototype.ExtraFlags)
	} else {
		ch.AffectedBy &= ^int64(obj.Prototype.ExtraFlags)
	}

	// Apply or remove affects from the object
	for i := 0; i < types.MAX_OBJ_AFFECT; i++ {
		affect := obj.Affects[i]
		if affect.Location != types.APPLY_NONE && affect.Modifier != 0 {
			w.affectModify(ch, affect.Location, affect.Modifier, 0, add)
		}
	}

	// Apply or remove affects from the prototype
	for i := 0; i < types.MAX_OBJ_AFFECT; i++ {
		affect := obj.Prototype.Affects[i]
		if affect.Location != types.APPLY_NONE && affect.Modifier != 0 {
			w.affectModify(ch, affect.Location, affect.Modifier, 0, add)
		}
	}

	// Update the character's total affects
	w.affectTotal(ch)
}
