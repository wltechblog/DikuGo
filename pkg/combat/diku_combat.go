package combat

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/wltechblog/DikuGo/pkg/types"
)

// DikuCombatManager implements the original DikuMUD combat system
type DikuCombatManager struct {
	// Combats is a map of character IDs to their combat state
	combats map[string]*Combat
}

// NewDikuCombatManager creates a new DikuMUD-style combat manager
func NewDikuCombatManager() *DikuCombatManager {
	return &DikuCombatManager{
		combats: make(map[string]*Combat),
	}
}

// StartCombat starts a combat between two characters
func (m *DikuCombatManager) StartCombat(attacker, defender *types.Character) error {
	// Check if the attacker is already in combat
	if attacker.Fighting != nil {
		return fmt.Errorf("you are already fighting")
	}

	// Ensure Equipment slices are properly initialized
	if attacker.Equipment == nil {
		attacker.Equipment = make([]*types.ObjectInstance, types.NUM_WEARS)
	} else if len(attacker.Equipment) < types.NUM_WEARS {
		// Resize the slice if it's too small
		newEquipment := make([]*types.ObjectInstance, types.NUM_WEARS)
		copy(newEquipment, attacker.Equipment)
		attacker.Equipment = newEquipment
	}

	if defender.Equipment == nil {
		defender.Equipment = make([]*types.ObjectInstance, types.NUM_WEARS)
	} else if len(defender.Equipment) < types.NUM_WEARS {
		// Resize the slice if it's too small
		newEquipment := make([]*types.ObjectInstance, types.NUM_WEARS)
		copy(newEquipment, defender.Equipment)
		defender.Equipment = newEquipment
	}

	// Check if the defender is already in combat with someone else
	if defender.Fighting != nil && defender.Fighting != attacker {
		return fmt.Errorf("%s is already fighting", defender.ShortDesc)
	}

	// Check if the attacker and defender are in the same room
	if attacker.InRoom != defender.InRoom {
		return fmt.Errorf("%s is not here", defender.ShortDesc)
	}

	// Check if the attacker is in a position to fight
	if attacker.Position < types.POS_STANDING {
		return fmt.Errorf("you are in no position to fight")
	}

	// Check if the defender is in a position to fight
	if defender.Position < types.POS_SLEEPING {
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
	if defender.Position >= types.POS_STANDING {
		defender.Position = types.POS_FIGHTING
	}

	// Add the combat to the map
	m.combats[attacker.Name] = combat
	m.combats[defender.Name] = combat

	// Log the combat
	log.Printf("Combat started: %s vs %s", attacker.Name, defender.Name)

	// Perform the first attack immediately
	m.doAttack(attacker, defender)

	return nil
}

// StopCombat stops a combat
func (m *DikuCombatManager) StopCombat(character *types.Character) {
	// Check if the character is in combat
	if character.Fighting == nil {
		return
	}

	// Get the combat
	combat, ok := m.combats[character.Name]
	if !ok {
		// Character is not in combat
		character.Fighting = nil
		if character.Position == types.POS_FIGHTING {
			character.Position = types.POS_STANDING
		}
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

	// Clear the fighting pointers
	character.Fighting = nil
	other.Fighting = nil

	// Set the positions
	if character.Position == types.POS_FIGHTING {
		character.Position = types.POS_STANDING
	}
	if other.Position == types.POS_FIGHTING {
		other.Position = types.POS_STANDING
	}

	// Log the combat end
	log.Printf("Combat ended: %s vs %s", character.Name, other.Name)
}

// Update updates all combats
func (m *DikuCombatManager) Update() {
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
		// Check if it's time for the next round (2 seconds between rounds)
		if now.Sub(combat.LastAttack) < 2*time.Second {
			continue
		}

		// Update the combat
		m.updateCombat(combat)
	}
}

// updateCombat updates a single combat
func (m *DikuCombatManager) updateCombat(combat *Combat) {
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

	// Check if attacker can attack
	if combat.Attacker.Position >= types.POS_FIGHTING {
		// Attacker attacks defender
		m.doAttack(combat.Attacker, combat.Defender)
	}

	// Check if the defender is dead
	if combat.Defender.Position <= types.POS_DEAD {
		// Stop the combat
		m.StopCombat(combat.Attacker)
		return
	}

	// Check if defender can attack
	if combat.Defender.Position >= types.POS_FIGHTING {
		// Defender attacks attacker
		m.doAttack(combat.Defender, combat.Attacker)
	}

	// Check if the attacker is dead
	if combat.Attacker.Position <= types.POS_DEAD {
		// Stop the combat
		m.StopCombat(combat.Defender)
		return
	}
}

// doAttack performs an attack using DikuMUD combat mechanics
func (m *DikuCombatManager) doAttack(attacker, defender *types.Character) {
	// Calculate THAC0 (To Hit Armor Class 0)
	// In DikuMUD, lower THAC0 is better
	thac0 := 20 - attacker.Level
	if attacker.IsNPC {
		// For NPCs, use their HitRoll as a direct modifier to THAC0
		thac0 -= attacker.HitRoll
	} else {
		// For players, apply strength bonus
		thac0 -= getStrengthHitBonus(attacker.Abilities[0]) // STR is index 0
	}

	// Calculate the victim's Armor Class
	// In DikuMUD, lower AC is better
	victimAC := defender.ArmorClass[0] / 10 // Convert from DikuGo's AC system

	// Apply dexterity bonus if the defender is awake
	if defender.Position > types.POS_SLEEPING {
		victimAC -= getDexterityACBonus(defender.Abilities[3]) // DEX is index 3
	}

	// Ensure AC is within bounds (-10 to 10)
	if victimAC < -10 {
		victimAC = -10
	} else if victimAC > 10 {
		victimAC = 10
	}

	// Roll the dice (1d20)
	diceRoll := rand.Intn(20) + 1

	// Calculate the roll needed to hit
	// In DikuMUD, you need to roll >= (THAC0 - AC)
	rollNeeded := thac0 - victimAC

	// Check if the attack hits
	// Natural 20 always hits, natural 1 always misses
	hit := (diceRoll == 20) || (diceRoll != 1 && diceRoll >= rollNeeded)

	if !hit {
		// Miss
		sendCombatMessage(attacker, defender, 0, getWeaponType(attacker))
		return
	}

	// Calculate damage exactly as in original DikuMUD
	// Start with minimum damage of 1
	damage := 1

	// Get weapon damage if wielding a weapon
	weaponDamage := 0
	// Make sure Equipment is properly initialized and has the correct length
	if attacker.Equipment != nil && len(attacker.Equipment) > types.WEAR_WIELD && attacker.Equipment[types.WEAR_WIELD] != nil {
		weapon := attacker.Equipment[types.WEAR_WIELD]
		// Use weapon damage dice (value[1]d value[2])
		numDice := weapon.Prototype.Value[1]
		sizeDice := weapon.Prototype.Value[2]
		for i := 0; i < numDice; i++ {
			weaponDamage += rand.Intn(sizeDice) + 1
		}
	} else if attacker.IsNPC {
		// NPCs have natural damage dice
		if attacker.Prototype != nil {
			// Use the damage dice values (DamageType, AttackType)
			numDice := attacker.Prototype.DamageType
			sizeDice := attacker.Prototype.AttackType

			// If the damage dice values are not set, fall back to hit dice
			if numDice == 0 || sizeDice == 0 {
				numDice = attacker.Prototype.Dice[0]
				sizeDice = attacker.Prototype.Dice[1]
			}

			// Roll the dice
			for i := 0; i < numDice; i++ {
				weaponDamage += rand.Intn(sizeDice) + 1
			}

			// Don't add the DamRoll here, it will be added separately below
		} else {
			// Fallback for NPCs without prototype
			weaponDamage = rand.Intn(attacker.Level) + 1
		}
	} else {
		// Bare hands for players - in original DikuMUD this is 0-2 damage
		weaponDamage = rand.Intn(3) // 0-2 damage
	}

	// Add strength damage bonus
	strDamBonus := getStrengthDamageBonus(attacker.Abilities[0])

	// Add damage roll bonus
	damRollBonus := attacker.DamRoll

	// Calculate total damage
	damage = weaponDamage + strDamBonus + damRollBonus

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

	// Ensure at least 1 damage
	if damage < 1 {
		damage = 1
	}

	// Apply the damage
	defender.HP -= damage

	// Send combat messages
	sendCombatMessage(attacker, defender, damage, getWeaponType(attacker))

	// Check if the defender died
	if defender.HP <= 0 {
		defender.HP = 0
		defender.Position = types.POS_DEAD

		// Send death messages
		attacker.SendMessage(fmt.Sprintf("You have slain %s!\r\n", defender.ShortDesc))
		defender.SendMessage(fmt.Sprintf("%s has slain you!\r\n", attacker.ShortDesc))

		// Send a message to the room
		for _, ch := range attacker.InRoom.Characters {
			if ch != attacker && ch != defender {
				ch.SendMessage(fmt.Sprintf("%s has slain %s!\r\n", attacker.ShortDesc, defender.ShortDesc))
			}
		}

		// Create a corpse
		w, ok := defender.World.(interface {
			MakeCorpse(*types.Character) *types.ObjectInstance
			RemoveCharacter(*types.Character)
			ScheduleMobRespawn(*types.Character)
		})
		if ok {
			// Create the corpse first
			w.MakeCorpse(defender)

			// If the defender is an NPC, schedule respawn
			if defender.IsNPC && defender.Prototype != nil {
				w.ScheduleMobRespawn(defender)
			} else {
				// For players or if scheduling fails, just remove from room
				w.RemoveCharacter(defender)
			}
		} else {
			log.Printf("Warning: Could not create corpse for %s", defender.Name)
		}

		// Handle death
		log.Printf("%s has been killed by %s!", defender.Name, attacker.Name)
	}
}

// getWeaponType returns the weapon type for combat messages
func getWeaponType(ch *types.Character) int {
	if ch.Equipment[types.WEAR_WIELD] != nil {
		// Use the weapon's type (value[3])
		return ch.Equipment[types.WEAR_WIELD].Prototype.Value[3]
	}
	// Default to TYPE_HIT for bare hands
	return types.TYPE_HIT
}

// getStrengthHitBonus returns the hit bonus for a given strength
func getStrengthHitBonus(str int) int {
	// Simplified strength table
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

// getStrengthDamageBonus returns the damage bonus for a given strength
func getStrengthDamageBonus(str int) int {
	// Simplified strength table
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

// getDexterityACBonus returns the AC bonus for a given dexterity
func getDexterityACBonus(dex int) int {
	// Simplified dexterity table
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

// sendCombatMessage sends appropriate combat messages based on damage
func sendCombatMessage(attacker, defender *types.Character, damage int, weaponType int) {
	// Get the weapon verb
	verb := getWeaponVerb(weaponType, damage > 0)

	if damage == 0 {
		// Miss message
		attacker.SendMessage(fmt.Sprintf("You miss %s with your %s.\r\n", defender.ShortDesc, verb))
		defender.SendMessage(fmt.Sprintf("%s misses you with %s %s.\r\n", attacker.ShortDesc, getHisHer(attacker), verb))

		// Send a message to the room
		for _, ch := range attacker.InRoom.Characters {
			if ch != attacker && ch != defender {
				ch.SendMessage(fmt.Sprintf("%s misses %s with %s %s.\r\n", attacker.ShortDesc, defender.ShortDesc, getHisHer(attacker), verb))
			}
		}
		return
	}

	// Hit messages based on damage amount
	var attackerMsg, defenderMsg, roomMsg string

	if damage <= 2 {
		attackerMsg = fmt.Sprintf("You tickle %s as you %s %s.\r\n", defender.ShortDesc, verb, getHimHer(defender))
		defenderMsg = fmt.Sprintf("%s tickles you as %s %s you.\r\n", attacker.ShortDesc, getHeShe(attacker), verb)
		roomMsg = fmt.Sprintf("%s tickles %s with %s %s.\r\n", attacker.ShortDesc, defender.ShortDesc, getHisHer(attacker), verb)
	} else if damage <= 4 {
		attackerMsg = fmt.Sprintf("You barely %s %s.\r\n", verb, defender.ShortDesc)
		defenderMsg = fmt.Sprintf("%s barely %s you.\r\n", attacker.ShortDesc, verb)
		roomMsg = fmt.Sprintf("%s barely %s %s.\r\n", attacker.ShortDesc, verb, defender.ShortDesc)
	} else if damage <= 6 {
		attackerMsg = fmt.Sprintf("You %s %s.\r\n", verb, defender.ShortDesc)
		defenderMsg = fmt.Sprintf("%s %s you.\r\n", attacker.ShortDesc, verb)
		roomMsg = fmt.Sprintf("%s %s %s.\r\n", attacker.ShortDesc, verb, defender.ShortDesc)
	} else if damage <= 10 {
		attackerMsg = fmt.Sprintf("You %s %s hard.\r\n", verb, defender.ShortDesc)
		defenderMsg = fmt.Sprintf("%s %s you hard.\r\n", attacker.ShortDesc, verb)
		roomMsg = fmt.Sprintf("%s %s %s hard.\r\n", attacker.ShortDesc, verb, defender.ShortDesc)
	} else if damage <= 15 {
		attackerMsg = fmt.Sprintf("You %s %s very hard.\r\n", verb, defender.ShortDesc)
		defenderMsg = fmt.Sprintf("%s %s you very hard.\r\n", attacker.ShortDesc, verb)
		roomMsg = fmt.Sprintf("%s %s %s very hard.\r\n", attacker.ShortDesc, verb, defender.ShortDesc)
	} else if damage <= 20 {
		attackerMsg = fmt.Sprintf("You %s %s extremely hard.\r\n", verb, defender.ShortDesc)
		defenderMsg = fmt.Sprintf("%s %s you extremely hard.\r\n", attacker.ShortDesc, verb)
		roomMsg = fmt.Sprintf("%s %s %s extremely hard.\r\n", attacker.ShortDesc, verb, defender.ShortDesc)
	} else {
		attackerMsg = fmt.Sprintf("You massacre %s to small fragments with your %s.\r\n", defender.ShortDesc, verb)
		defenderMsg = fmt.Sprintf("%s massacres you to small fragments with %s %s.\r\n", attacker.ShortDesc, getHisHer(attacker), verb)
		roomMsg = fmt.Sprintf("%s massacres %s to small fragments with %s %s.\r\n", attacker.ShortDesc, defender.ShortDesc, getHisHer(attacker), verb)
	}

	// Send the messages
	attacker.SendMessage(attackerMsg)
	defender.SendMessage(defenderMsg)

	// Send a message to the room
	for _, ch := range attacker.InRoom.Characters {
		if ch != attacker && ch != defender {
			ch.SendMessage(roomMsg)
		}
	}
}

// getWeaponVerb returns the verb for a weapon type
func getWeaponVerb(weaponType int, hit bool) string {
	// Define weapon verbs (singular form)
	verbs := map[int]string{
		types.TYPE_HIT:      "hit",
		types.TYPE_BLUDGEON: "pound",
		types.TYPE_PIERCE:   "pierce",
		types.TYPE_SLASH:    "slash",
		types.TYPE_WHIP:     "whip",
		types.TYPE_CLAW:     "claw",
		types.TYPE_BITE:     "bite",
		types.TYPE_STING:    "sting",
		types.TYPE_CRUSH:    "crush",
	}

	// Get the verb for this weapon type
	verb, ok := verbs[weaponType]
	if !ok {
		verb = "hit" // Default to "hit" if weapon type not found
	}

	return verb
}

// getHeShe returns "he" or "she" based on character gender
func getHeShe(ch *types.Character) string {
	if ch.Sex == types.SEX_FEMALE {
		return "she"
	} else if ch.Sex == types.SEX_MALE {
		return "he"
	}
	return "it"
}

// getHimHer returns "him" or "her" based on character gender
func getHimHer(ch *types.Character) string {
	if ch.Sex == types.SEX_FEMALE {
		return "her"
	} else if ch.Sex == types.SEX_MALE {
		return "him"
	}
	return "it"
}

// getHisHer returns "his" or "her" based on character gender
func getHisHer(ch *types.Character) string {
	if ch.Sex == types.SEX_FEMALE {
		return "her"
	} else if ch.Sex == types.SEX_MALE {
		return "his"
	}
	return "its"
}
