package command

import (
	"strings"
	"testing"

	"github.com/wltechblog/DikuGo/pkg/types"
)

// MockWorldForWho implements the interface needed by WhoCommand
type MockWorldForWho struct {
	characters map[string]*types.Character
}

func (m *MockWorldForWho) GetCharacters() map[string]*types.Character {
	return m.characters
}

func TestWhoCommand(t *testing.T) {
	// Create a mock world with some characters
	mockWorld := &MockWorldForWho{
		characters: make(map[string]*types.Character),
	}

	// Create test characters
	player1 := &types.Character{
		Name:  "TestPlayer1",
		Level: 5,
		Title: " the Warrior",
		IsNPC: false,
		World: mockWorld,
	}

	player2 := &types.Character{
		Name:  "TestPlayer2",
		Level: 10,
		Title: " the Mage",
		IsNPC: false,
		World: mockWorld,
	}

	// Create an NPC (should not appear in who list)
	npc := &types.Character{
		Name:  "TestNPC",
		Level: 1,
		Title: "",
		IsNPC: true,
		World: mockWorld,
	}

	// Add characters to the mock world
	mockWorld.characters["TestPlayer1"] = player1
	mockWorld.characters["TestPlayer2"] = player2
	mockWorld.characters["TestNPC"] = npc

	// Create the who command
	whoCmd := &WhoCommand{}

	// Execute the command
	err := whoCmd.Execute(player1, "")
	if err == nil {
		t.Fatalf("Expected who command to return output as error, got nil")
	}

	output := err.Error()

	// Check that the output contains the expected header
	if !strings.Contains(output, "Players online:") {
		t.Errorf("Expected output to contain 'Players online:', got: %s", output)
	}

	// Check that both players are listed
	if !strings.Contains(output, "TestPlayer1") {
		t.Errorf("Expected output to contain 'TestPlayer1', got: %s", output)
	}

	if !strings.Contains(output, "TestPlayer2") {
		t.Errorf("Expected output to contain 'TestPlayer2', got: %s", output)
	}

	// Check that the NPC is NOT listed
	if strings.Contains(output, "TestNPC") {
		t.Errorf("Expected output to NOT contain 'TestNPC', got: %s", output)
	}

	// Check that levels and titles are displayed
	if !strings.Contains(output, "[ 5] TestPlayer1 the Warrior") {
		t.Errorf("Expected output to contain player1 with level and title, got: %s", output)
	}

	if !strings.Contains(output, "[10] TestPlayer2 the Mage") {
		t.Errorf("Expected output to contain player2 with level and title, got: %s", output)
	}
}

func TestWhoCommandEmptyWorld(t *testing.T) {
	// Create a mock world with no characters
	mockWorld := &MockWorldForWho{
		characters: make(map[string]*types.Character),
	}

	// Create a test character
	player := &types.Character{
		Name:  "TestPlayer",
		Level: 1,
		Title: " the newbie",
		IsNPC: false,
		World: mockWorld,
	}

	// Create the who command
	whoCmd := &WhoCommand{}

	// Execute the command
	err := whoCmd.Execute(player, "")
	if err == nil {
		t.Fatalf("Expected who command to return output as error, got nil")
	}

	output := err.Error()

	// Check that the output contains the expected header
	if !strings.Contains(output, "Players online:") {
		t.Errorf("Expected output to contain 'Players online:', got: %s", output)
	}

	// Check that no players are listed (except the header)
	lines := strings.Split(output, "\r\n")
	playerLines := 0
	for _, line := range lines {
		if strings.Contains(line, "[") && strings.Contains(line, "]") {
			playerLines++
		}
	}

	if playerLines != 0 {
		t.Errorf("Expected no player lines in empty world, got %d lines: %s", playerLines, output)
	}
}

func TestWhoCommandNoWorldInterface(t *testing.T) {
	// Create a character with no world interface
	player := &types.Character{
		Name:  "TestPlayer",
		Level: 1,
		Title: " the newbie",
		IsNPC: false,
		World: nil, // No world interface
	}

	// Create the who command
	whoCmd := &WhoCommand{}

	// Execute the command
	err := whoCmd.Execute(player, "")
	if err == nil {
		t.Fatalf("Expected who command to return error for missing world interface, got nil")
	}

	if !strings.Contains(err.Error(), "world interface not available") {
		t.Errorf("Expected error about missing world interface, got: %s", err.Error())
	}
}
