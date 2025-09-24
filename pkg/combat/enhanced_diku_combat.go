package combat

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/wltechblog/DikuGo/pkg/types"
)

// EnhancedDikuCombatManager implements an enhanced version of the original DikuMUD combat system
type EnhancedDikuCombatManager struct {
	// Combats is a map of character names to their combat state
	Combats map[string]*CombatState
	// LastCombatTick is the last time combat was processed
	LastCombatTick time.Time
}

// CombatState represents the state of a character in combat
type CombatState struct {
	// Character is the character in combat
	Character *types.Character
	// Target is the character's target
	Target *types.Character
	// LastAttack is the last time the character attacked
	LastAttack time.Time
}

// NewEnhancedDikuCombatManager creates a new EnhancedDikuCombatManager
func NewEnhancedDikuCombatManager() *EnhancedDikuCombatManager {
	return &EnhancedDikuCombatManager{
		Combats:        make(map[string]*CombatState),
		LastCombatTick: time.Now(),
	}
}

// StartCombat starts combat between two characters
func (m *EnhancedDikuCombatManager) StartCombat(attacker, defender *types.Character) error {
	// Check if the attacker is already in combat
	if _, ok := m.Combats[attacker.Name]; ok {
		// Update the target
		m.Combats[attacker.Name].Target = defender
		return nil
	}

	// Create a new combat state
	m.Combats[attacker.Name] = &CombatState{
		Character:  attacker,
		Target:     defender,
		LastAttack: time.Now().Add(-1 * time.Second), // Allow immediate attack
	}

	// Set the character's position and fighting state
	attacker.Position = types.POS_FIGHTING
	attacker.Fighting = defender

	// If the defender is not already fighting the attacker, start combat for them too
	if _, ok := m.Combats[defender.Name]; !ok {
		m.Combats[defender.Name] = &CombatState{
			Character:  defender,
			Target:     attacker,
			LastAttack: time.Now().Add(-1 * time.Second), // Allow immediate attack
		}

		// Set the defender's position and fighting state
		defender.Position = types.POS_FIGHTING
		defender.Fighting = attacker
	}

	return nil
}

// StopCombat stops combat for a character
func (m *EnhancedDikuCombatManager) StopCombat(character *types.Character) {
	// Check if the character is in combat
	combatState, ok := m.Combats[character.Name]
	if !ok {
		return
	}

	// Get the target to clear their fighting state too
	target := combatState.Target

	// Remove the character from combat
	delete(m.Combats, character.Name)

	// Also remove the target from combat if they exist
	if target != nil {
		delete(m.Combats, target.Name)
	}

	// Clear the fighting pointers
	character.Fighting = nil
	if target != nil {
		target.Fighting = nil
	}

	// Set the character's position to standing if they were fighting
	if character.Position == types.POS_FIGHTING {
		character.Position = types.POS_STANDING
	}

	// Set the target's position to standing if they were fighting
	if target != nil && target.Position == types.POS_FIGHTING {
		target.Position = types.POS_STANDING
	}
}

// Update updates the combat manager
func (m *EnhancedDikuCombatManager) Update() {
	m.ProcessCombat()
}

// ProcessCombat processes combat for all characters
func (m *EnhancedDikuCombatManager) ProcessCombat() {
	// Check if it's time to process combat
	now := time.Now()
	if now.Sub(m.LastCombatTick) < 2*time.Second {
		return
	}

	// Update the last combat tick
	m.LastCombatTick = now

	// Collect characters that need to stop combat (to avoid modifying map while iterating)
	var charactersToStopCombat []*types.Character
	var charactersToAttack []*CombatState

	// Process combat for each character
	for id, state := range m.Combats {
		// Skip if the character or target is nil
		if state.Character == nil || state.Target == nil {
			delete(m.Combats, id)
			continue
		}

		// Skip if the character is not in a room
		if state.Character.InRoom == nil {
			// Clear fighting state before removing
			if state.Character != nil {
				state.Character.Fighting = nil
				if state.Character.Position == types.POS_FIGHTING {
					state.Character.Position = types.POS_STANDING
				}
			}
			delete(m.Combats, id)
			continue
		}

		// Skip if the target is not in the same room
		if state.Target.InRoom != state.Character.InRoom {
			// Mark for stopping combat when characters are separated
			charactersToStopCombat = append(charactersToStopCombat, state.Character)
			continue
		}

		// Skip if the character is not in a fighting position
		if state.Character.Position != types.POS_FIGHTING {
			// Mark for stopping combat if character is no longer fighting
			charactersToStopCombat = append(charactersToStopCombat, state.Character)
			continue
		}

		// Skip if the target is dead
		if state.Target.Position == types.POS_DEAD {
			// Mark for stopping combat when target dies
			charactersToStopCombat = append(charactersToStopCombat, state.Character)
			continue
		}

		// Check if it's time for the character to attack
		if now.Sub(state.LastAttack) < 2*time.Second {
			continue
		}

		// Update the last attack time
		state.LastAttack = now

		// Mark for attack
		charactersToAttack = append(charactersToAttack, state)
	}

	// Stop combat for characters that need it
	for _, character := range charactersToStopCombat {
		m.StopCombat(character)
	}

	// Perform attacks for characters that can attack
	for _, state := range charactersToAttack {
		m.doAttack(state.Character, state.Target)
	}
}

// doAttack performs an attack from one character to another
func (m *EnhancedDikuCombatManager) doAttack(attacker, defender *types.Character) {
	// Check if the attacker is in a position to attack
	if attacker.Position != types.POS_FIGHTING {
		return
	}

	// Check if the defender is in a position to be attacked
	if defender.Position == types.POS_DEAD {
		return
	}

	// Calculate THAC0 for the attacker
	thac0 := calculateTHAC0(attacker)

	// Calculate AC for the defender
	ac := calculateAC(defender)

	// Roll to hit
	roll := rand.Intn(20) + 1 // 1-20

	// Check if the attack hits
	if roll < (thac0 - ac) {
		// Miss
		sendEnhancedCombatMessage(attacker, defender, 0, getEnhancedWeaponType(attacker))
		return
	}

	// Calculate damage
	damage := calculateDamage(attacker, defender)

	// Apply damage
	defender.HP -= damage

	// Send combat messages
	sendEnhancedCombatMessage(attacker, defender, damage, getEnhancedWeaponType(attacker))

	// Check if the defender died
	if defender.HP <= 0 {
		defender.HP = 0
		defender.Position = types.POS_DEAD

		// Stop combat for both characters
		m.StopCombat(attacker)

		// Send death messages
		attacker.SendMessage(fmt.Sprintf("You have slain %s!\r\n", defender.ShortDesc))
		defender.SendMessage(fmt.Sprintf("%s has slain you!\r\n", attacker.ShortDesc))

		// Send a message to the room
		for _, ch := range attacker.InRoom.Characters {
			if ch != attacker && ch != defender {
				ch.SendMessage(fmt.Sprintf("%s has slain %s!\r\n", attacker.ShortDesc, defender.ShortDesc))
			}
		}

		// Handle character death using the centralized death handler
		w, ok := defender.World.(interface {
			HandleCharacterDeath(*types.Character)
		})
		if ok {
			w.HandleCharacterDeath(defender)
		} else {
			log.Printf("Warning: Could not handle death for %s", defender.Name)
		}

		// Award experience
		if !defender.IsNPC {
			// PCs don't give experience
			return
		}

		// Calculate experience
		exp := calculateExperience(attacker, defender)

		// Award experience
		attacker.Experience += exp
		attacker.SendMessage(fmt.Sprintf("You gain %d experience points.\r\n", exp))
	}
}

// calculateTHAC0 calculates the THAC0 for a character
func calculateTHAC0(ch *types.Character) int {
	// Base THAC0 is 20 - level
	thac0 := 20 - ch.Level

	// Apply hitroll bonus
	thac0 -= ch.HitRoll

	// Apply strength bonus
	thac0 -= getEnhancedStrengthHitBonus(ch.Abilities[0]) // STR is index 0

	return thac0
}

// calculateAC calculates the AC for a character
func calculateAC(ch *types.Character) int {
	// Base AC is the character's AC
	ac := ch.ArmorClass[0] / 10 // Convert from DikuMUD's AC*10 format

	// Apply dexterity bonus if awake
	if ch.Position > types.POS_SLEEPING {
		ac += getEnhancedDexterityACBonus(ch.Abilities[3]) // DEX is index 3
	}

	// Ensure AC is within bounds
	if ac < -10 {
		ac = -10 // -10 is the lowest AC in DikuMUD
	}

	return ac
}

// calculateDamage calculates the damage for an attack
func calculateDamage(attacker, defender *types.Character) int {
	// Start with weapon damage
	var weaponDamage int
	weapon := attacker.Equipment[types.WEAR_WIELD]
	if weapon != nil {
		// Use the weapon's damage dice
		numDice := weapon.Prototype.Value[1]
		sizeDice := weapon.Prototype.Value[2]
		if numDice > 0 && sizeDice > 0 {
			// Roll the dice
			for i := 0; i < numDice; i++ {
				weaponDamage += rand.Intn(sizeDice) + 1
			}
		}
	} else if attacker.IsNPC {
		// NPCs use their dice values
		if attacker.Prototype != nil {
			numDice := attacker.Prototype.Dice[0]
			sizeDice := attacker.Prototype.Dice[1]
			bonus := attacker.Prototype.Dice[2]
			if numDice > 0 && sizeDice > 0 {
				// Roll the dice
				for i := 0; i < numDice; i++ {
					weaponDamage += rand.Intn(sizeDice) + 1
				}
				// Add the bonus
				weaponDamage += bonus
			}
		}
	} else {
		// Bare hands for players - in original DikuMUD this is 0-2 damage
		weaponDamage = rand.Intn(3) // 0-2 damage
	}

	// Add strength damage bonus
	strDamBonus := getEnhancedStrengthDamageBonus(attacker.Abilities[0])

	// Add damage roll bonus
	damRollBonus := attacker.DamRoll

	// Calculate total damage
	damage := weaponDamage + strDamBonus + damRollBonus

	// Position modifier (more damage to non-fighting targets)
	if defender.Position < types.POS_FIGHTING {
		positionMult := 1.0 + float64(types.POS_FIGHTING-defender.Position)/3.0
		damage = int(float64(damage) * positionMult)
	}

	// Ensure minimum damage
	if damage < 1 {
		damage = 1
	}

	// Apply damage exactly as in original DikuMUD
	// In DikuMUD, damage is capped at 100 per hit
	if damage > 100 {
		damage = 100
	}

	return damage
}

// calculateExperience calculates the experience gained from killing a mob
func calculateExperience(attacker, defender *types.Character) int {
	// Base experience is 1/3 of the mob's experience
	exp := defender.Experience / 3

	// Apply level difference bonus
	levelDiff := defender.Level - attacker.Level
	if levelDiff > 0 {
		// Bonus for killing higher level mobs
		if attacker.IsNPC {
			// NPCs get less bonus
			exp += (exp * levelDiff) / 8
		} else {
			// Players get more bonus
			exp += (exp * levelDiff) / 4
		}
	}

	// Ensure minimum experience
	if exp < 1 {
		exp = 1
	}

	return exp
}

// sendEnhancedCombatMessage sends combat messages to the attacker, defender, and room
func sendEnhancedCombatMessage(attacker, defender *types.Character, damage int, weaponType int) {
	// Get the attack messages
	var attackerMsg, defenderMsg, roomMsg string

	if damage == 0 {
		// Miss
		attackerMsg = fmt.Sprintf("You miss %s.\r\n", defender.ShortDesc)
		defenderMsg = fmt.Sprintf("%s misses you.\r\n", attacker.ShortDesc)
		roomMsg = fmt.Sprintf("%s misses %s.\r\n", attacker.ShortDesc, defender.ShortDesc)
	} else {
		// Hit
		var verb string
		switch weaponType {
		case types.TYPE_HIT:
			verb = "hit"
		case types.TYPE_SLASH:
			verb = "slash"
		case types.TYPE_PIERCE:
			verb = "pierce"
		case types.TYPE_BLUDGEON:
			verb = "pound"
		case types.TYPE_CRUSH:
			verb = "crush"
		default:
			verb = "hit"
		}

		// Damage level messages
		var damageDesc string
		if damage <= 2 {
			damageDesc = "barely"
		} else if damage <= 4 {
			damageDesc = "slightly"
		} else if damage <= 6 {
			damageDesc = "fairly hard"
		} else if damage <= 10 {
			damageDesc = "hard"
		} else if damage <= 15 {
			damageDesc = "very hard"
		} else if damage <= 20 {
			damageDesc = "extremely hard"
		} else {
			damageDesc = "EXTREMELY hard"
		}

		attackerMsg = fmt.Sprintf("You %s %s %s.\r\n", verb, defender.ShortDesc, damageDesc)
		defenderMsg = fmt.Sprintf("%s %ss you %s.\r\n", attacker.ShortDesc, verb, damageDesc)
		roomMsg = fmt.Sprintf("%s %ss %s %s.\r\n", attacker.ShortDesc, verb, defender.ShortDesc, damageDesc)
	}

	// Send the messages
	attacker.SendMessage(attackerMsg)
	defender.SendMessage(defenderMsg)
	for _, ch := range attacker.InRoom.Characters {
		if ch != attacker && ch != defender {
			ch.SendMessage(roomMsg)
		}
	}
}

// getEnhancedWeaponType returns the weapon type for a character
func getEnhancedWeaponType(ch *types.Character) int {
	// Check if equipment array is properly initialized and has enough slots
	if ch.Equipment == nil || len(ch.Equipment) <= types.WEAR_WIELD {
		return types.TYPE_HIT
	}

	weapon := ch.Equipment[types.WEAR_WIELD]
	if weapon != nil && weapon.Prototype != nil {
		// Use the weapon's type
		return weapon.Prototype.Value[3]
	}
	// Default to TYPE_HIT for bare hands
	return types.TYPE_HIT
}

// getEnhancedStrengthHitBonus returns the hit bonus for a given strength
func getEnhancedStrengthHitBonus(str int) int {
	// Apply strength bonus to hit
	switch {
	case str <= 3:
		return -3
	case str <= 5:
		return -2
	case str <= 7:
		return -1
	case str <= 13:
		return 0
	case str <= 15:
		return 1
	case str <= 17:
		return 2
	case str <= 18:
		return 3
	default:
		return 4
	}
}

// getEnhancedStrengthDamageBonus returns the damage bonus for a given strength
func getEnhancedStrengthDamageBonus(str int) int {
	// Apply strength bonus to damage
	switch {
	case str <= 3:
		return -3
	case str <= 5:
		return -2
	case str <= 7:
		return -1
	case str <= 13:
		return 0
	case str <= 15:
		return 1
	case str <= 17:
		return 2
	case str <= 18:
		return 3
	default:
		return 4
	}
}

// getEnhancedDexterityACBonus returns the AC bonus for a given dexterity
func getEnhancedDexterityACBonus(dex int) int {
	// Apply dexterity bonus to AC
	switch {
	case dex <= 3:
		return -3
	case dex <= 5:
		return -2
	case dex <= 7:
		return -1
	case dex <= 13:
		return 0
	case dex <= 15:
		return 1
	case dex <= 17:
		return 2
	case dex <= 18:
		return 3
	default:
		return 4
	}
}
