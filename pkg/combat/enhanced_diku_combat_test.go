package combat

import (
	"testing"
	"time"

	"github.com/wltechblog/DikuGo/pkg/types"
)

// TestFightingStateClearedOnDeath tests that the Fighting field is properly cleared when a character dies
func TestFightingStateClearedOnDeath(t *testing.T) {
	// Create a combat manager
	manager := NewEnhancedDikuCombatManager()

	// Create a test room
	room := &types.Room{
		VNUM:       3001,
		Name:       "Test Room",
		Characters: make([]*types.Character, 0),
	}

	// Create a player character
	player := &types.Character{
		Name:         "TestPlayer",
		Level:        5,
		Position:     types.POS_STANDING,
		HP:           50,
		MaxHitPoints: 50,
		HitRoll:      10,
		DamRoll:      5,
		InRoom:       room,
		IsNPC:        false,
		Equipment:    make([]*types.ObjectInstance, types.NUM_WEARS),
		Inventory:    make([]*types.ObjectInstance, 0),
	}

	// Create a weak mob that will die quickly
	mob := &types.Character{
		Name:         "TestMob",
		ShortDesc:    "a test mob",
		Level:        1,
		Position:     types.POS_STANDING,
		HP:           1, // Very low HP so it dies in one hit
		MaxHitPoints: 1,
		HitRoll:      0,
		DamRoll:      0,
		InRoom:       room,
		IsNPC:        true,
		Equipment:    make([]*types.ObjectInstance, types.NUM_WEARS),
		Inventory:    make([]*types.ObjectInstance, 0),
	}

	// Add characters to room
	room.Characters = append(room.Characters, player, mob)

	// Start combat
	err := manager.StartCombat(player, mob)
	if err != nil {
		t.Fatalf("Failed to start combat: %v", err)
	}

	// Verify fighting state is set
	if player.Fighting != mob {
		t.Error("Player should be fighting the mob")
	}
	if mob.Fighting != player {
		t.Error("Mob should be fighting the player")
	}

	// Verify positions are set to fighting
	if player.Position != types.POS_FIGHTING {
		t.Error("Player position should be POS_FIGHTING")
	}
	if mob.Position != types.POS_FIGHTING {
		t.Error("Mob position should be POS_FIGHTING")
	}

	// Simulate the mob dying by directly setting HP to 0 and triggering death handling
	mob.HP = 0
	mob.Position = types.POS_DEAD

	// Call StopCombat to simulate what should happen when a character dies
	manager.StopCombat(player)

	// Verify fighting state is cleared
	if player.Fighting != nil {
		t.Error("Player Fighting field should be nil after mob dies")
	}
	if mob.Fighting != nil {
		t.Error("Mob Fighting field should be nil after death")
	}

	// Verify positions are reset
	if player.Position != types.POS_STANDING {
		t.Error("Player position should be POS_STANDING after combat ends")
	}
	if mob.Position != types.POS_DEAD {
		t.Error("Mob position should be POS_DEAD after death")
	}

	// Verify combat states are removed from manager
	if _, exists := manager.Combats[player.Name]; exists {
		t.Error("Player should be removed from combat manager")
	}
	if _, exists := manager.Combats[mob.Name]; exists {
		t.Error("Mob should be removed from combat manager")
	}
}

// TestStopCombatClearsFightingState tests that StopCombat properly clears fighting state
func TestStopCombatClearsFightingState(t *testing.T) {
	// Create a combat manager
	manager := NewEnhancedDikuCombatManager()

	// Create a test room
	room := &types.Room{
		VNUM:       3001,
		Name:       "Test Room",
		Characters: make([]*types.Character, 0),
	}

	// Create two characters
	char1 := &types.Character{
		Name:      "Character1",
		Level:     5,
		Position:  types.POS_STANDING,
		InRoom:    room,
		IsNPC:     false,
		Equipment: make([]*types.ObjectInstance, types.NUM_WEARS),
		Inventory: make([]*types.ObjectInstance, 0),
	}

	char2 := &types.Character{
		Name:      "Character2",
		ShortDesc: "character two",
		Level:     5,
		Position:  types.POS_STANDING,
		InRoom:    room,
		IsNPC:     true,
		Equipment: make([]*types.ObjectInstance, types.NUM_WEARS),
		Inventory: make([]*types.ObjectInstance, 0),
	}

	// Add characters to room
	room.Characters = append(room.Characters, char1, char2)

	// Start combat
	err := manager.StartCombat(char1, char2)
	if err != nil {
		t.Fatalf("Failed to start combat: %v", err)
	}

	// Verify fighting state is set
	if char1.Fighting != char2 {
		t.Error("Character1 should be fighting Character2")
	}
	if char2.Fighting != char1 {
		t.Error("Character2 should be fighting Character1")
	}

	// Stop combat
	manager.StopCombat(char1)

	// Verify fighting state is cleared for both characters
	if char1.Fighting != nil {
		t.Error("Character1 Fighting field should be nil after StopCombat")
	}
	if char2.Fighting != nil {
		t.Error("Character2 Fighting field should be nil after StopCombat")
	}

	// Verify positions are reset
	if char1.Position != types.POS_STANDING {
		t.Error("Character1 position should be POS_STANDING after StopCombat")
	}
	if char2.Position != types.POS_STANDING {
		t.Error("Character2 position should be POS_STANDING after StopCombat")
	}

	// Verify combat states are removed from manager
	if _, exists := manager.Combats[char1.Name]; exists {
		t.Error("Character1 should be removed from combat manager")
	}
	if _, exists := manager.Combats[char2.Name]; exists {
		t.Error("Character2 should be removed from combat manager")
	}
}

// TestProcessCombatClearsDeadTargets tests that ProcessCombat properly handles dead targets
func TestProcessCombatClearsDeadTargets(t *testing.T) {
	// Create a combat manager
	manager := NewEnhancedDikuCombatManager()

	// Create a test room
	room := &types.Room{
		VNUM:       3001,
		Name:       "Test Room",
		Characters: make([]*types.Character, 0),
	}

	// Create two characters
	attacker := &types.Character{
		Name:      "Attacker",
		Level:     5,
		Position:  types.POS_FIGHTING,
		InRoom:    room,
		IsNPC:     false,
		Fighting:  nil, // Will be set by combat manager
		Equipment: make([]*types.ObjectInstance, types.NUM_WEARS),
		Inventory: make([]*types.ObjectInstance, 0),
	}

	victim := &types.Character{
		Name:      "Victim",
		ShortDesc: "a victim",
		Level:     5,
		Position:  types.POS_DEAD, // Already dead
		InRoom:    room,
		IsNPC:     true,
		Fighting:  nil, // Will be set by combat manager
		Equipment: make([]*types.ObjectInstance, types.NUM_WEARS),
		Inventory: make([]*types.ObjectInstance, 0),
	}

	// Set the last combat tick to allow processing
	manager.LastCombatTick = time.Now().Add(-3 * time.Second)

	// Manually set up combat state (simulating a situation where victim died but combat wasn't cleaned up)
	manager.Combats[attacker.Name] = &CombatState{
		Character:  attacker,
		Target:     victim,
		LastAttack: manager.LastCombatTick,
	}
	manager.Combats[victim.Name] = &CombatState{
		Character:  victim,
		Target:     attacker,
		LastAttack: manager.LastCombatTick,
	}
	attacker.Fighting = victim
	victim.Fighting = attacker

	// Verify initial state
	if len(manager.Combats) != 2 {
		t.Errorf("Expected 2 combat states, got %d", len(manager.Combats))
	}

	// Process combat - this should clean up the dead target
	manager.ProcessCombat()

	// Verify fighting state is cleared
	if attacker.Fighting != nil {
		t.Errorf("Attacker Fighting field should be nil after ProcessCombat with dead target, got %v", attacker.Fighting)
	}
	if victim.Fighting != nil {
		t.Errorf("Victim Fighting field should be nil after ProcessCombat, got %v", victim.Fighting)
	}

	// Verify attacker position is reset
	if attacker.Position != types.POS_STANDING {
		t.Errorf("Attacker position should be POS_STANDING after ProcessCombat with dead target, got %d", attacker.Position)
	}

	// Verify combat states are removed from manager
	if _, exists := manager.Combats[attacker.Name]; exists {
		t.Error("Attacker should be removed from combat manager")
	}
	if _, exists := manager.Combats[victim.Name]; exists {
		t.Error("Victim should be removed from combat manager")
	}

	// Debug: Check final combat state count
	if len(manager.Combats) != 0 {
		t.Errorf("Expected 0 combat states after ProcessCombat, got %d", len(manager.Combats))
	}
}
