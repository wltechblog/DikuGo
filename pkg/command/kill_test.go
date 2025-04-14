package command

import (
	"fmt"
	"strings"
	"testing"

	"github.com/wltechblog/DikuGo/pkg/types"
)

// MockCombatManager is a mock implementation of the combat manager for testing
type MockCombatManager struct{}

// StartCombat is a mock implementation that just sets the fighting pointers
func (m *MockCombatManager) StartCombat(attacker, defender *types.Character) error {
	// Set the fighting pointers
	attacker.Fighting = defender
	defender.Fighting = attacker

	// Return a combat message
	return fmt.Errorf("You hit %s!", defender.ShortDesc)
}

func TestKillCommand(t *testing.T) {
	// Create a test room
	room := &types.Room{
		VNUM:        3001,
		Name:        "Test Room",
		Description: "This is a test room.",
		Characters:  make([]*types.Character, 0),
	}

	// Create a test player
	player := &types.Character{
		Name:         "TestPlayer",
		ShortDesc:    "",
		IsNPC:        false,
		InRoom:       room,
		Level:        1,
		Position:     types.POS_STANDING,
		HP:           20,
		MaxHitPoints: 20,
		HitRoll:      1,
		DamRoll:      1,
	}
	room.Characters = append(room.Characters, player)

	// Create a test mob
	mob := &types.Character{
		Name:         "testmob",
		ShortDesc:    "a test mob",
		LongDesc:     "A test mob is standing here.",
		Description:  "This is a test mob for unit testing.",
		IsNPC:        true,
		ActFlags:     types.ACT_ISNPC,
		InRoom:       room,
		Level:        1,
		Position:     types.POS_STANDING,
		HP:           10,
		MaxHitPoints: 10,
		HitRoll:      0,
		DamRoll:      0,
	}
	room.Characters = append(room.Characters, mob)

	// Create a kill command with the mock combat manager
	killCmd := &KillCommand{
		CombatManager: &MockCombatManager{},
	}

	// Test killing the mob
	err := killCmd.Execute(player, "testmob")
	if err == nil {
		t.Errorf("Expected an error (which contains the combat message), got nil")
	} else {
		// The error should contain a combat message
		if !strings.Contains(err.Error(), "You hit") {
			t.Errorf("Expected error to contain combat message, got: %s", err.Error())
		}
	}

	// Check if the player is fighting the mob
	if player.Fighting != mob {
		t.Errorf("Expected player to be fighting the mob")
	}

	// Check if the mob is fighting the player
	if mob.Fighting != player {
		t.Errorf("Expected mob to be fighting the player")
	}

	// Test killing a non-existent target
	err = killCmd.Execute(player, "nonexistent")
	if err == nil {
		t.Errorf("Expected an error for non-existent target, got nil")
	} else {
		// The error should indicate that the target doesn't exist
		if !strings.Contains(err.Error(), "they aren't here") {
			t.Errorf("Expected error to indicate target doesn't exist, got: %s", err.Error())
		}
	}

	// Test killing a player (should not be allowed)
	player2 := &types.Character{
		Name:         "TestPlayer2",
		ShortDesc:    "",
		IsNPC:        false,
		InRoom:       room,
		Level:        1,
		Position:     types.POS_STANDING,
		HP:           20,
		MaxHitPoints: 20,
	}
	room.Characters = append(room.Characters, player2)

	err = killCmd.Execute(player, "TestPlayer2")
	if err == nil {
		t.Errorf("Expected an error for killing a player, got nil")
	} else {
		// The error should indicate that killing players is not allowed
		if !strings.Contains(err.Error(), "you can't kill other players") {
			t.Errorf("Expected error to indicate killing players is not allowed, got: %s", err.Error())
		}
	}
}
