package world

import (
	"log"
	"math/rand"

	"github.com/wltechblog/DikuGo/pkg/types"
)

// PulseViolence handles combat updates
func (w *World) PulseViolence() {
	// Make a copy of the characters to avoid deadlock
	w.mutex.RLock()
	characters := make([]*types.Character, 0, len(w.characters))
	for _, character := range w.characters {
		characters = append(characters, character)
	}
	w.mutex.RUnlock()

	// Process combat for all characters
	for _, character := range characters {
		if character.Fighting != nil && character.Position == types.POS_STANDING {
			w.performCombatRound(character)
		}
	}
}

// performCombatRound handles a single round of combat for a character
func (w *World) performCombatRound(ch *types.Character) {
	// Check if the character is still fighting
	if ch.Fighting == nil || ch.Position != types.POS_STANDING {
		return
	}

	// Check if the target is still valid
	victim := ch.Fighting
	if victim.Position == types.POS_DEAD {
		ch.Fighting = nil
		return
	}

	// Calculate hit chance
	hitChance := 50 + ch.HitRoll
	if hitChance < 5 {
		hitChance = 5
	} else if hitChance > 95 {
		hitChance = 95
	}

	// Roll to hit
	roll := rand.Intn(100) + 1
	if roll <= hitChance {
		// Hit! Calculate damage
		damage := rand.Intn(ch.DamRoll+1) + 1
		if damage < 1 {
			damage = 1
		}

		// Apply damage
		victim.HP -= damage

		// Check if the victim died
		if victim.HP <= 0 {
			victim.HP = 0
			victim.Position = types.POS_DEAD
			victim.Fighting = nil

			// Stop anyone fighting the victim
			for _, character := range w.characters {
				if character.Fighting == victim {
					character.Fighting = nil
				}
			}

			log.Printf("%s has been killed by %s!", victim.Name, ch.Name)
		}
	}
}
