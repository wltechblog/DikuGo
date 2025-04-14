package combat

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/wltechblog/DikuGo/pkg/types"
)

// Manager handles combat between characters
type Manager struct {
	// Combats is a map of character IDs to their combat state
	combats map[string]*Combat
}

// Combat represents a combat between two characters
type Combat struct {
	// Attacker is the character that initiated the combat
	Attacker *types.Character

	// Defender is the character that is being attacked
	Defender *types.Character

	// Round is the current combat round
	Round int

	// LastAttack is the time of the last attack
	LastAttack time.Time
}

// NewManager creates a new combat manager
func NewManager() *Manager {
	return &Manager{
		combats: make(map[string]*Combat),
	}
}

// StartCombat starts a combat between two characters
func (m *Manager) StartCombat(attacker, defender *types.Character) error {
	// Check if the attacker is already in combat
	if attacker.Fighting != nil {
		return fmt.Errorf("you are already fighting")
	}

	// Check if the defender is already in combat
	if defender.Fighting != nil {
		return fmt.Errorf("%s is already fighting", defender.ShortDesc)
	}

	// Check if the attacker and defender are in the same room
	if attacker.InRoom != defender.InRoom {
		return fmt.Errorf("%s is not here", defender.ShortDesc)
	}

	// Check if the attacker is in a position to fight
	if attacker.Position < types.POS_FIGHTING {
		return fmt.Errorf("you are in no position to fight")
	}

	// Check if the defender is in a position to fight
	if defender.Position < types.POS_FIGHTING {
		return fmt.Errorf("%s is in no position to fight", defender.ShortDesc)
	}

	// Create a new combat
	combat := &Combat{
		Attacker:   attacker,
		Defender:   defender,
		Round:      0,
		LastAttack: time.Now(),
	}

	// Set the fighting pointers
	attacker.Fighting = defender
	defender.Fighting = attacker

	// Set the positions
	attacker.Position = types.POS_FIGHTING
	defender.Position = types.POS_FIGHTING

	// Add the combat to the map
	m.combats[attacker.Name] = combat
	m.combats[defender.Name] = combat

	// Log the combat
	log.Printf("Combat started: %s vs %s", attacker.Name, defender.Name)

	return nil
}

// StopCombat stops a combat
func (m *Manager) StopCombat(character *types.Character) {
	// Check if the character is in combat
	if character.Fighting == nil {
		return
	}

	// Get the combat
	combat, ok := m.combats[character.Name]
	if !ok {
		// Character is not in combat
		character.Fighting = nil
		character.Position = types.POS_STANDING
		return
	}

	// Get the other character
	var other *types.Character
	if combat.Attacker == character {
		other = combat.Defender
	} else {
		other = combat.Attacker
	}

	// Remove the combat from the map
	delete(m.combats, character.Name)
	delete(m.combats, other.Name)

	// Reset the fighting pointers
	character.Fighting = nil
	other.Fighting = nil

	// Reset the positions
	character.Position = types.POS_STANDING
	other.Position = types.POS_STANDING

	// Log the combat
	log.Printf("Combat stopped: %s vs %s", combat.Attacker.Name, combat.Defender.Name)
}

// Update updates all combats
func (m *Manager) Update() {
	// Get the current time
	now := time.Now()

	// Create a list of combats to update
	var combats []*Combat
	for _, combat := range m.combats {
		// Check if we've already processed this combat
		alreadyProcessed := false
		for _, c := range combats {
			if c == combat {
				alreadyProcessed = true
				break
			}
		}
		if alreadyProcessed {
			continue
		}

		// Add the combat to the list
		combats = append(combats, combat)
	}

	// Update each combat
	for _, combat := range combats {
		// Check if it's time for the next round
		if now.Sub(combat.LastAttack) < 2*time.Second {
			continue
		}

		// Update the combat
		m.updateCombat(combat)
	}
}

// updateCombat updates a single combat
func (m *Manager) updateCombat(combat *Combat) {
	// Increment the round
	combat.Round++

	// Update the last attack time
	combat.LastAttack = time.Now()

	// Check if either character is dead
	if combat.Attacker.Position <= types.POS_DEAD || combat.Defender.Position <= types.POS_DEAD {
		// Stop the combat
		m.StopCombat(combat.Attacker)
		return
	}

	// Attacker attacks defender
	m.doAttack(combat.Attacker, combat.Defender)

	// Check if the defender is dead
	if combat.Defender.Position <= types.POS_DEAD {
		// Stop the combat
		m.StopCombat(combat.Attacker)
		return
	}

	// Defender attacks attacker
	m.doAttack(combat.Defender, combat.Attacker)

	// Check if the attacker is dead
	if combat.Attacker.Position <= types.POS_DEAD {
		// Stop the combat
		m.StopCombat(combat.Defender)
		return
	}
}

// doAttack performs an attack
func (m *Manager) doAttack(attacker, defender *types.Character) {
	// Calculate the hit chance
	hitChance := 50 + (attacker.Level - defender.Level) * 5

	// Roll the dice
	roll := rand.Intn(100)

	// Check if the attack hits
	if roll < hitChance {
		// Calculate the damage
		damage := rand.Intn(10) + 1 + attacker.Level

		// Apply the damage
		defender.HP -= damage

		// Send messages
		attacker.SendMessage(fmt.Sprintf("You hit %s for %d damage.\r\n", defender.ShortDesc, damage))
		defender.SendMessage(fmt.Sprintf("%s hits you for %d damage.\r\n", attacker.ShortDesc, damage))

		// Send a message to the room
		for _, ch := range attacker.InRoom.Characters {
			if ch != attacker && ch != defender {
				ch.SendMessage(fmt.Sprintf("%s hits %s.\r\n", attacker.ShortDesc, defender.ShortDesc))
			}
		}

		// Check if the defender is dead
		if defender.HP <= 0 {
			// Set the position to dead
			defender.Position = types.POS_DEAD

			// Send messages
			attacker.SendMessage(fmt.Sprintf("You have slain %s!\r\n", defender.ShortDesc))
			defender.SendMessage(fmt.Sprintf("%s has slain you!\r\n", attacker.ShortDesc))

			// Send a message to the room
			for _, ch := range attacker.InRoom.Characters {
				if ch != attacker && ch != defender {
					ch.SendMessage(fmt.Sprintf("%s has slain %s!\r\n", attacker.ShortDesc, defender.ShortDesc))
				}
			}
		}
	} else {
		// Send messages
		attacker.SendMessage(fmt.Sprintf("You miss %s.\r\n", defender.ShortDesc))
		defender.SendMessage(fmt.Sprintf("%s misses you.\r\n", attacker.ShortDesc))

		// Send a message to the room
		for _, ch := range attacker.InRoom.Characters {
			if ch != attacker && ch != defender {
				ch.SendMessage(fmt.Sprintf("%s misses %s.\r\n", attacker.ShortDesc, defender.ShortDesc))
			}
		}
	}
}
