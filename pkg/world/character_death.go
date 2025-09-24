package world

import (
	"log"

	"github.com/wltechblog/DikuGo/pkg/types"
)

// HandleCharacterDeath handles the death of a character
func (w *World) HandleCharacterDeath(victim *types.Character) {
	if victim == nil {
		return
	}

	log.Printf("HandleCharacterDeath: Character %s has died", victim.Name)

	// Create a corpse with the victim's items
	w.MakeCorpse(victim)

	// If the victim is a player, handle player death
	if !victim.IsNPC {
		// Send death message to the player
		victim.SendMessage("\r\nYou have been KILLED!!\r\n")
		victim.SendMessage("You find yourself back at the temple...\r\n")

		// Reset character for resurrection
		victim.HP = 1
		victim.Position = types.POS_STANDING

		// Set the special flag to indicate the player should return to the menu
		victim.SendMessage("RETURN_TO_MENU")
	} else if victim.Prototype != nil {
		// Schedule respawn for NPCs
		w.ScheduleMobRespawn(victim)
	}
}
