package game

import (
	"log"
	"strings"

	"github.com/wltechblog/DikuGo/pkg/ai"
	"github.com/wltechblog/DikuGo/pkg/types"
	"github.com/wltechblog/DikuGo/pkg/world"
)

// registerSpecialProcs registers special procedures for mobiles
func registerSpecialProcs(w *world.World) {
	// Get all mobile prototypes
	mobiles := w.GetMobilePrototypes()

	// Register special procedures based on keywords
	for _, mobile := range mobiles {
		// Check for special procedures based on keywords
		for name, proc := range ai.SpecialProcs {
			if strings.Contains(strings.ToLower(mobile.Name), name) {
				log.Printf("Registering special procedure %s for mobile %s", name, mobile.Name)
				mobile.Functions = append(mobile.Functions, proc)
			}
		}

		// Register special procedures based on flags
		if mobile.ActFlags&types.ACT_SCAVENGER != 0 {
			log.Printf("Registering scavenger behavior for mobile %s", mobile.Name)
			// Scavenger behavior is handled by the AI system
		}

		if mobile.ActFlags&types.ACT_AGGRESSIVE != 0 {
			log.Printf("Registering aggressive behavior for mobile %s", mobile.Name)
			// Aggressive behavior is handled by the AI system
		}
	}
}
