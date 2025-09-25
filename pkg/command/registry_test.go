package command

import (
	"testing"
)

func TestCommandRegistration(t *testing.T) {
	// Test that the new equipment commands can be created
	testCases := []struct {
		commandName string
		command     Command
	}{
		{"hold", &HoldCommand{}},
		{"wield", &WieldCommand{}},
		{"wear", &WearCommand{}},
		{"remove", &RemoveCommand{}},
		{"equipment", &EquipmentCommand{}},
	}

	for _, tc := range testCases {
		if tc.command.Name() != tc.commandName {
			t.Errorf("Expected command name to be '%s', got '%s'", tc.commandName, tc.command.Name())
		}
	}
}

func TestHoldCommandAliases(t *testing.T) {
	holdCmd := &HoldCommand{}

	// Test command name
	if holdCmd.Name() != "hold" {
		t.Errorf("Expected command name to be 'hold', got '%s'", holdCmd.Name())
	}

	// Test aliases
	aliases := holdCmd.Aliases()
	expectedAliases := []string{"grab"}

	if len(aliases) != len(expectedAliases) {
		t.Errorf("Expected %d aliases, got %d", len(expectedAliases), len(aliases))
	}

	for i, expected := range expectedAliases {
		if i >= len(aliases) || aliases[i] != expected {
			t.Errorf("Expected alias[%d] to be '%s', got '%s'", i, expected, aliases[i])
		}
	}
}

func TestWieldCommandProperties(t *testing.T) {
	wieldCmd := &WieldCommand{}

	// Test command name
	if wieldCmd.Name() != "wield" {
		t.Errorf("Expected command name to be 'wield', got '%s'", wieldCmd.Name())
	}

	// Test that it has no aliases (unlike original DikuMUD which had no wield aliases)
	aliases := wieldCmd.Aliases()
	if len(aliases) != 0 {
		t.Errorf("Expected no aliases for wield command, got %d", len(aliases))
	}

	// Test minimum position
	if wieldCmd.MinPosition() != 5 { // POS_RESTING
		t.Errorf("Expected minimum position to be 5 (POS_RESTING), got %d", wieldCmd.MinPosition())
	}

	// Test level requirement
	if wieldCmd.Level() != 0 {
		t.Errorf("Expected level requirement to be 0, got %d", wieldCmd.Level())
	}

	// Test log command
	if wieldCmd.LogCommand() != false {
		t.Errorf("Expected LogCommand to be false, got %t", wieldCmd.LogCommand())
	}
}
