package world

import (
	"fmt"
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
		if victim.Client != nil {
			if client, ok := victim.Client.(interface {
				Write(string)
				SetState(int)
			}); ok {
				client.Write(fmt.Sprintf("\r\nYou have been KILLED!!\r\n"))

				// Save the character
				w.SaveCharacter(victim)

				// Remove from world
				w.RemoveCharacter(victim)

				// Reset character for resurrection
				victim.HP = 1
				victim.Position = types.POS_STANDING

				// Return to menu
				client.Write("\r\n                                 DikuGo Main Menu\r\n\r\n1) Enter the game\r\n2) Enter description\r\n3) Read the background story\r\n4) Change password\r\n0) Exit from DikuGo\r\n\r\nMake your choice: ")
				client.SetState(3) // StateMainMenu = 3
			}
		}
	} else if victim.Prototype != nil {
		// Schedule respawn for NPCs
		w.ScheduleMobRespawn(victim)
	}
}
