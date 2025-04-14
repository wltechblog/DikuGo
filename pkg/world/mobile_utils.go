package world

import (
	"github.com/wltechblog/DikuGo/pkg/types"
)

// GetMobilePrototypes returns all mobile prototypes
func (w *World) GetMobilePrototypes() []*types.Mobile {
	w.mutex.RLock()
	defer w.mutex.RUnlock()

	// Convert the map to a slice
	mobiles := make([]*types.Mobile, 0, len(w.mobiles))
	for _, mobile := range w.mobiles {
		mobiles = append(mobiles, mobile)
	}

	return mobiles
}

// GetMobile returns a mobile prototype by VNUM
func (w *World) GetMobile(vnum int) *types.Mobile {
	w.mutex.RLock()
	defer w.mutex.RUnlock()
	return w.mobiles[vnum]
}
