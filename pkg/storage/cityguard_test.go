package storage

import (
	"os"
	"testing"

	"github.com/wltechblog/DikuGo/pkg/types"
)

func TestParseCityguardAndCreateMob(t *testing.T) {
	// Create a temporary file with cityguard data
	tmpFile, err := os.CreateTemp("", "test_cityguard.*.mob")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	// Write cityguard data to the file (copied from tinyworld.mob)
	testData := `#3060
cityguard guard~
the Cityguard~
A cityguard stands here.~
A big, strong, helpful, trustworthy guard.
~
193 0 1000 S
10 10 2 1d12+123 1d8+3
500 9000
8 8 1
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

	// Check the cityguard stats
	cityguard := mobiles[0]
	if cityguard.VNUM != 3060 {
		t.Errorf("Expected VNUM 3060, got %d", cityguard.VNUM)
	}
	if cityguard.Name != "cityguard guard" {
		t.Errorf("Expected name 'cityguard guard', got '%s'", cityguard.Name)
	}
	if cityguard.ShortDesc != "the Cityguard" {
		t.Errorf("Expected short desc 'the Cityguard', got '%s'", cityguard.ShortDesc)
	}
	if cityguard.LongDesc != "A cityguard stands here." {
		t.Errorf("Expected long desc 'A cityguard stands here.', got '%s'", cityguard.LongDesc)
	}
	if cityguard.Description != "A big, strong, helpful, trustworthy guard." {
		t.Errorf("Expected description 'A big, strong, helpful, trustworthy guard.', got '%s'", cityguard.Description)
	}
	if cityguard.ActFlags != 193 {
		t.Errorf("Expected act flags 193, got %d", cityguard.ActFlags)
	}

	// Check the key stats based on the test data
	if cityguard.Level != 10 {
		t.Errorf("Expected level 10, got %d", cityguard.Level)
	}
	if cityguard.HitRoll != 10 {
		t.Errorf("Expected hitroll 10, got %d", cityguard.HitRoll)
	}
	if cityguard.AC[0] != 2 || cityguard.AC[1] != 2 || cityguard.AC[2] != 2 {
		t.Errorf("Expected AC [2 2 2], got %v", cityguard.AC)
	}
	if cityguard.Dice[0] != 1 || cityguard.Dice[1] != 12 || cityguard.Dice[2] != 123 {
		t.Errorf("Expected dice 1d12+123, got %dd%d+%d", cityguard.Dice[0], cityguard.Dice[1], cityguard.Dice[2])
	}
	if cityguard.DamRoll != 3 {
		t.Errorf("Expected damroll 3, got %d", cityguard.DamRoll)
	}
	if cityguard.Gold != 500 {
		t.Errorf("Expected gold 500, got %d", cityguard.Gold)
	}
	if cityguard.Experience != 9000 {
		t.Errorf("Expected experience 9000, got %d", cityguard.Experience)
	}
	if cityguard.Position != 8 {
		t.Errorf("Expected position 8, got %d", cityguard.Position)
	}
	if cityguard.DefaultPos != 8 {
		t.Errorf("Expected default position 8, got %d", cityguard.DefaultPos)
	}
	if cityguard.Sex != 1 {
		t.Errorf("Expected sex 1, got %d", cityguard.Sex)
	}

	// Now create a character from the mobile prototype
	// This simulates what happens in the game when a mob is created
	character := &types.Character{
		Name:        cityguard.Name,
		ShortDesc:   cityguard.ShortDesc,
		LongDesc:    cityguard.LongDesc,
		Description: cityguard.Description,
		Level:       cityguard.Level,
		Sex:         cityguard.Sex,
		Class:       cityguard.Class,
		Race:        cityguard.Race,
		Gold:        cityguard.Gold,
		Experience:  cityguard.Experience,
		Alignment:   cityguard.Alignment,
		Position:    cityguard.Position,
		ArmorClass:  cityguard.AC,
		HitRoll:     cityguard.HitRoll,
		DamRoll:     cityguard.DamRoll,
		Abilities:   cityguard.Abilities,
		ActFlags:    cityguard.ActFlags,
		IsNPC:       true,
		Prototype:   cityguard,
	}

	// Check that the character has the correct stats
	if character.Level != 10 {
		t.Errorf("Mob: Expected level 10, got %d", character.Level)
	}
	if character.HitRoll != 10 {
		t.Errorf("Mob: Expected hitroll 10, got %d", character.HitRoll)
	}
	if character.ArmorClass[0] != 2 || character.ArmorClass[1] != 2 || character.ArmorClass[2] != 2 {
		t.Errorf("Mob: Expected AC [2 2 2], got %v", character.ArmorClass)
	}
	if character.DamRoll != 3 {
		t.Errorf("Mob: Expected damroll 3, got %d", character.DamRoll)
	}
	if character.Gold != 500 {
		t.Errorf("Mob: Expected gold 500, got %d", character.Gold)
	}
	if character.Experience != 9000 {
		t.Errorf("Mob: Expected experience 9000, got %d", character.Experience)
	}
}
