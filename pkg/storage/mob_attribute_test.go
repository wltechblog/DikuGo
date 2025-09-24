package storage

import (
	"os"
	"testing"

	"github.com/wltechblog/DikuGo/pkg/types"
	"github.com/wltechblog/DikuGo/pkg/utils"
)

func TestMobAttributesParsing(t *testing.T) {
	// Skip this test if running in CI environment
	if os.Getenv("CI") != "" {
		t.Skip("Skipping test in CI environment")
	}

	// Create a temporary file with test data for a simple mob
	tmpFile, err := os.CreateTemp("", "test_simple_mob.*.mob")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	// Write test data to the file - this is a typical simple mob format
	testData := `#1234
test mob~
a test mob~
A test mob is standing here.~
This is a test mob for unit testing.
~
8 0 0 S
5
15
5
2d6+50
1d4+3
100
1000
8
8
1
$~
`
	if _, err := tmpFile.Write([]byte(testData)); err != nil {
		t.Fatalf("Failed to write to temporary file: %v", err)
	}
	if err := tmpFile.Close(); err != nil {
		t.Fatalf("Failed to close temporary file: %v", err)
	}

	// Parse the test data
	mobiles, err := ParseMobiles(tmpFile.Name())
	if err != nil {
		t.Fatalf("ParseMobiles failed: %v", err)
	}

	// Check that we got the expected number of mobiles
	if len(mobiles) != 1 {
		t.Fatalf("Expected 1 mobile, got %d", len(mobiles))
	}

	// Check the simple mobile attributes
	mob := mobiles[0]
	if mob.VNUM != 1234 {
		t.Errorf("Expected VNUM 1234, got %d", mob.VNUM)
	}
	if mob.Name != "test mob" {
		t.Errorf("Expected name 'test mob', got '%s'", mob.Name)
	}
	if mob.Level != 5 {
		t.Errorf("Expected level 5, got %d", mob.Level)
	}

	// Check hitroll directly
	expectedHitroll := 15 // Hitroll is 15
	if mob.HitRoll != expectedHitroll {
		t.Errorf("Expected hitroll %d, got %d", expectedHitroll, mob.HitRoll)
	}

	// Check armor class directly
	expectedAC := 5 // AC is 5
	if mob.AC[0] != expectedAC {
		t.Errorf("Expected AC %d, got %d", expectedAC, mob.AC[0])
	}

	// Check hit dice
	if mob.Dice[0] != 2 || mob.Dice[1] != 6 || mob.Dice[2] != 50 {
		t.Errorf("Expected hit dice 2d6+50, got %dd%d+%d", mob.Dice[0], mob.Dice[1], mob.Dice[2])
	}

	// Check damage roll (strength bonus)
	expectedDamRoll := 3 // From 1d4+3
	if mob.DamRoll != expectedDamRoll {
		t.Errorf("Expected damroll %d, got %d", expectedDamRoll, mob.DamRoll)
	}

	// Check gold and experience
	if mob.Gold != 100 {
		t.Errorf("Expected gold 100, got %d", mob.Gold)
	}
	if mob.Experience != 1000 {
		t.Errorf("Expected experience 1000, got %d", mob.Experience)
	}

	// Now create a character from the mobile prototype
	// This simulates what happens in the game when a mob is created
	character := &types.Character{
		Name:        mob.Name,
		ShortDesc:   mob.ShortDesc,
		LongDesc:    mob.LongDesc,
		Description: mob.Description,
		Level:       mob.Level,
		Sex:         mob.Sex,
		Class:       mob.Class,
		Race:        mob.Race,
		Gold:        mob.Gold,
		Experience:  mob.Experience,
		Alignment:   mob.Alignment,
		Position:    mob.Position,
		ArmorClass:  mob.AC,
		HitRoll:     mob.HitRoll,
		DamRoll:     mob.DamRoll,
		Abilities:   mob.Abilities,
		ActFlags:    mob.ActFlags,
		IsNPC:       true,
		Prototype:   mob,
	}

	// Calculate HP based on dice
	baseHP := utils.Dice(mob.Dice[0], mob.Dice[1]) + mob.Dice[2]
	character.HP = baseHP
	character.MaxHitPoints = baseHP

	// Check that the character has the correct stats
	if character.Level != 5 {
		t.Errorf("Character: Expected level 5, got %d", character.Level)
	}
	if character.HitRoll != expectedHitroll {
		t.Errorf("Character: Expected hitroll %d, got %d", expectedHitroll, character.HitRoll)
	}
	if character.ArmorClass[0] != expectedAC {
		t.Errorf("Character: Expected AC %d, got %d", expectedAC, character.ArmorClass[0])
	}
	if character.DamRoll != expectedDamRoll {
		t.Errorf("Character: Expected damroll %d, got %d", expectedDamRoll, character.DamRoll)
	}
	if character.Gold != 100 {
		t.Errorf("Character: Expected gold 100, got %d", character.Gold)
	}
	if character.Experience != 1000 {
		t.Errorf("Character: Expected experience 1000, got %d", character.Experience)
	}
}
