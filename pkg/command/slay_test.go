package command

import (
	"strings"
	"testing"

	"github.com/wltechblog/DikuGo/pkg/types"
)

// MockSlayWorld implements the world interface needed for slay command testing
type MockSlayWorld struct {
	deathHandled  bool
	deadCharacter *types.Character
}

func (w *MockSlayWorld) HandleCharacterDeath(victim *types.Character) {
	w.deathHandled = true
	w.deadCharacter = victim
}

// MockSlayTarget is a mock character that can receive messages
type MockSlayTarget struct {
	*types.Character
	messages []string
}

func (m *MockSlayTarget) SendMessage(msg string) {
	m.messages = append(m.messages, msg)
}

func (m *MockSlayTarget) HasMessage(msg string) bool {
	for _, message := range m.messages {
		if strings.Contains(message, msg) {
			return true
		}
	}
	return false
}

// MockSlayCombatManager implements CombatManagerInterface for testing
type MockSlayCombatManager struct {
	combatStarted bool
	attacker      *types.Character
	defender      *types.Character
}

func (m *MockSlayCombatManager) StartCombat(attacker, defender *types.Character) error {
	m.combatStarted = true
	m.attacker = attacker
	m.defender = defender

	// Set fighting pointers
	attacker.Fighting = defender
	defender.Fighting = attacker

	// Simulate the combat system dealing lethal damage
	// The slay command should have boosted damage high enough to kill
	damage := attacker.DamRoll
	defender.HP -= damage

	if defender.HP <= 0 {
		defender.HP = 0
		defender.Position = types.POS_DEAD

		// Award experience (simulate combat system behavior)
		if defender.IsNPC {
			exp := defender.Experience / 3
			if exp < 1 {
				exp = 1
			}
			attacker.Experience += exp
		}

		// Handle death if world interface is available
		if w, ok := defender.World.(*MockSlayWorld); ok {
			w.HandleCharacterDeath(defender)
		}
	}

	return nil
}

func (m *MockSlayCombatManager) StopCombat(character *types.Character) {
	character.Fighting = nil
}

func (m *MockSlayCombatManager) Update() {
	// No-op for testing
}

func TestSlayCommand_Execute(t *testing.T) {
	// Create a test room
	room := &types.Room{
		VNUM:        3001,
		Name:        "Test Room",
		Description: "A test room for slay command testing.",
		Characters:  make([]*types.Character, 0),
	}

	// Create a mock world
	mockWorld := &MockSlayWorld{}

	// Create a test player (admin)
	player := &MockSlayTarget{
		Character: &types.Character{
			Name:         "admin",
			ShortDesc:    "an admin",
			LongDesc:     "An admin is standing here.",
			Description:  "This is a test admin for unit testing.",
			IsNPC:        false,
			InRoom:       room,
			Level:        25, // High level admin
			Position:     types.POS_STANDING,
			HP:           100,
			MaxHitPoints: 100,
			HitRoll:      5,
			DamRoll:      10,
			Experience:   50000,
			World:        mockWorld,
		},
		messages: make([]string, 0),
	}
	room.Characters = append(room.Characters, player.Character)

	// Create a test mob
	mob := &MockSlayTarget{
		Character: &types.Character{
			Name:         "testmob",
			ShortDesc:    "a test mob",
			LongDesc:     "A test mob is standing here.",
			Description:  "This is a test mob for unit testing.",
			IsNPC:        true,
			ActFlags:     types.ACT_ISNPC,
			InRoom:       room,
			Level:        10,
			Position:     types.POS_STANDING,
			HP:           80,
			MaxHitPoints: 80,
			HitRoll:      5,
			DamRoll:      8,
			Experience:   3000,
			World:        mockWorld,
		},
		messages: make([]string, 0),
	}
	room.Characters = append(room.Characters, mob.Character)

	// Create a slay command with the mock combat manager
	mockCombat := &MockSlayCombatManager{}
	slayCmd := &SlayCommand{
		CombatManager: mockCombat,
	}

	// Test 1: Successful slay
	originalExp := player.Experience
	err := slayCmd.Execute(player.Character, "testmob")

	// Should return a success message as error (command pattern)
	if err == nil || !strings.Contains(err.Error(), "slay") {
		t.Errorf("Expected slay success message, got: %v", err)
	}

	// Verify combat was started
	if !mockCombat.combatStarted {
		t.Error("Expected combat to be started")
	}

	// Verify the mob was killed
	if mob.HP > 0 {
		t.Errorf("Expected mob to be dead (HP=0), got HP=%d", mob.HP)
	}

	if mob.Position != types.POS_DEAD {
		t.Errorf("Expected mob position to be DEAD (%d), got %d", types.POS_DEAD, mob.Position)
	}

	// Verify death was handled
	if !mockWorld.deathHandled {
		t.Error("Expected death to be handled by world")
	}

	if mockWorld.deadCharacter != mob.Character {
		t.Error("Expected the mob to be the dead character")
	}

	// Verify experience was awarded
	if player.Experience <= originalExp {
		t.Errorf("Expected player to gain experience, original: %d, current: %d", originalExp, player.Experience)
	}

	// Test 2: No target specified
	err = slayCmd.Execute(player.Character, "")
	if err == nil || !strings.Contains(err.Error(), "slay whom") {
		t.Errorf("Expected 'slay whom?' error, got: %v", err)
	}

	// Test 3: Target not found
	err = slayCmd.Execute(player.Character, "nonexistent")
	if err == nil || !strings.Contains(err.Error(), "aren't here") {
		t.Errorf("Expected 'aren't here' error, got: %v", err)
	}

	// Test 4: Try to slay self
	err = slayCmd.Execute(player.Character, "admin")
	if err == nil || !strings.Contains(err.Error(), "can't slay yourself") {
		t.Errorf("Expected 'can't slay yourself' error, got: %v", err)
	}
}

func TestSlayCommand_Properties(t *testing.T) {
	slayCmd := &SlayCommand{}

	// Test command name
	if slayCmd.Name() != "slay" {
		t.Errorf("Expected command name 'slay', got '%s'", slayCmd.Name())
	}

	// Test aliases (should be empty)
	aliases := slayCmd.Aliases()
	if len(aliases) != 0 {
		t.Errorf("Expected no aliases, got %v", aliases)
	}

	// Test minimum position (should be standing)
	if slayCmd.MinPosition() != types.POS_STANDING {
		t.Errorf("Expected minimum position POS_STANDING (%d), got %d", types.POS_STANDING, slayCmd.MinPosition())
	}

	// Test minimum level (should be admin level)
	if slayCmd.Level() != 1 {
		t.Errorf("Expected minimum level 1, got %d", slayCmd.Level())
	}

	// Test logging (should be true for admin commands)
	if !slayCmd.LogCommand() {
		t.Error("Expected slay command to be logged")
	}
}

func TestSlayCommand_AlreadyDead(t *testing.T) {
	// Create a test room
	room := &types.Room{
		VNUM:        3001,
		Name:        "Test Room",
		Description: "A test room for slay command testing.",
		Characters:  make([]*types.Character, 0),
	}

	// Create a test player (admin)
	player := &types.Character{
		Name:         "admin",
		ShortDesc:    "an admin",
		IsNPC:        false,
		InRoom:       room,
		Level:        25,
		Position:     types.POS_STANDING,
		HP:           100,
		MaxHitPoints: 100,
	}
	room.Characters = append(room.Characters, player)

	// Create a dead mob
	deadMob := &types.Character{
		Name:         "deadmob",
		ShortDesc:    "a dead mob",
		IsNPC:        true,
		InRoom:       room,
		Level:        10,
		Position:     types.POS_DEAD, // Already dead
		HP:           0,
		MaxHitPoints: 80,
	}
	room.Characters = append(room.Characters, deadMob)

	// Create a slay command
	slayCmd := &SlayCommand{
		CombatManager: &MockSlayCombatManager{},
	}

	// Test slaying already dead target
	err := slayCmd.Execute(player, "deadmob")
	if err == nil || !strings.Contains(err.Error(), "already dead") {
		t.Errorf("Expected 'already dead' error, got: %v", err)
	}
}
