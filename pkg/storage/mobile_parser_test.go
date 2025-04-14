package storage

import (
	"os"
	"testing"
)

func TestParseMobiles(t *testing.T) {
	// Create a temporary file with test data
	tmpFile, err := os.CreateTemp("", "test_mobiles.*.mob")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	// Write test data to the file
	testData := `#100
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
1d4+1
100
1000
8
8
1
#101
detailed mob~
a detailed mob~
A detailed mob is standing here.~
This is a detailed mob for unit testing.
~
8 0 0 D
18
12
12
12
12
100 150
5
100
50
200
2000
8
8
1
0
10
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
	if len(mobiles) != 2 {
		t.Fatalf("Expected 2 mobiles, got %d", len(mobiles))
	}

	// Check the simple mobile
	simpleMob := mobiles[0]
	if simpleMob.VNUM != 100 {
		t.Errorf("Expected VNUM 100, got %d", simpleMob.VNUM)
	}
	if simpleMob.Name != "test mob" {
		t.Errorf("Expected name 'test mob', got '%s'", simpleMob.Name)
	}
	if simpleMob.ShortDesc != "a test mob" {
		t.Errorf("Expected short desc 'a test mob', got '%s'", simpleMob.ShortDesc)
	}
	if simpleMob.LongDesc != "A test mob is standing here." {
		t.Errorf("Expected long desc 'A test mob is standing here.', got '%s'", simpleMob.LongDesc)
	}
	if simpleMob.Description != "This is a test mob for unit testing." {
		t.Errorf("Expected description 'This is a test mob for unit testing.', got '%s'", simpleMob.Description)
	}
	if simpleMob.ActFlags != 8 {
		t.Errorf("Expected act flags 8, got %d", simpleMob.ActFlags)
	}
	if simpleMob.Level != 5 {
		t.Errorf("Expected level 5, got %d", simpleMob.Level)
	}
	if simpleMob.Gold != 100 {
		t.Errorf("Expected gold 100, got %d", simpleMob.Gold)
	}
	if simpleMob.Experience != 1000 {
		t.Errorf("Expected experience 1000, got %d", simpleMob.Experience)
	}
	if simpleMob.Position != 8 {
		t.Errorf("Expected position 8, got %d", simpleMob.Position)
	}
	if simpleMob.DefaultPos != 8 {
		t.Errorf("Expected default position 8, got %d", simpleMob.DefaultPos)
	}
	if simpleMob.Sex != 1 {
		t.Errorf("Expected sex 1, got %d", simpleMob.Sex)
	}
	if simpleMob.Dice[0] != 2 || simpleMob.Dice[1] != 6 || simpleMob.Dice[2] != 50 {
		t.Errorf("Expected dice 2d6+50, got %dd%d+%d", simpleMob.Dice[0], simpleMob.Dice[1], simpleMob.Dice[2])
	}
	if simpleMob.DamRoll != 1 {
		t.Errorf("Expected damroll 1, got %d", simpleMob.DamRoll)
	}

	// Check the detailed mobile
	detailedMob := mobiles[1]
	if detailedMob.VNUM != 101 {
		t.Errorf("Expected VNUM 101, got %d", detailedMob.VNUM)
	}
	if detailedMob.Name != "detailed mob" {
		t.Errorf("Expected name 'detailed mob', got '%s'", detailedMob.Name)
	}
	if detailedMob.ShortDesc != "a detailed mob" {
		t.Errorf("Expected short desc 'a detailed mob', got '%s'", detailedMob.ShortDesc)
	}
	if detailedMob.LongDesc != "A detailed mob is standing here." {
		t.Errorf("Expected long desc 'A detailed mob is standing here.', got '%s'", detailedMob.LongDesc)
	}
	if detailedMob.Description != "This is a detailed mob for unit testing." {
		t.Errorf("Expected description 'This is a detailed mob for unit testing.', got '%s'", detailedMob.Description)
	}
	if detailedMob.ActFlags != 8 {
		t.Errorf("Expected act flags 8, got %d", detailedMob.ActFlags)
	}
	if detailedMob.Level != 10 {
		t.Errorf("Expected level 10, got %d", detailedMob.Level)
	}
	if detailedMob.Gold != 200 {
		t.Errorf("Expected gold 200, got %d", detailedMob.Gold)
	}
	if detailedMob.Experience != 2000 {
		t.Errorf("Expected experience 2000, got %d", detailedMob.Experience)
	}
	if detailedMob.Position != 8 {
		t.Errorf("Expected position 8, got %d", detailedMob.Position)
	}
	if detailedMob.DefaultPos != 8 {
		t.Errorf("Expected default position 8, got %d", detailedMob.DefaultPos)
	}
	if detailedMob.Sex != 1 {
		t.Errorf("Expected sex 1, got %d", detailedMob.Sex)
	}
	if detailedMob.Class != 0 {
		t.Errorf("Expected class 0, got %d", detailedMob.Class)
	}
	if detailedMob.Abilities[0] != 18 {
		t.Errorf("Expected strength 18, got %d", detailedMob.Abilities[0])
	}
	if detailedMob.Abilities[1] != 12 {
		t.Errorf("Expected intelligence 12, got %d", detailedMob.Abilities[1])
	}
	if detailedMob.Abilities[2] != 12 {
		t.Errorf("Expected wisdom 12, got %d", detailedMob.Abilities[2])
	}
	if detailedMob.Abilities[3] != 12 {
		t.Errorf("Expected dexterity 12, got %d", detailedMob.Abilities[3])
	}
	if detailedMob.Abilities[4] != 12 {
		t.Errorf("Expected constitution 12, got %d", detailedMob.Abilities[4])
	}
	if detailedMob.HitRoll != 7 {
		t.Errorf("Expected hitroll 7, got %d", detailedMob.HitRoll)
	}
}
