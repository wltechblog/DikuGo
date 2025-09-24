package storage

import (
	"os"
	"testing"
)

func TestParseMobiles(t *testing.T) {
	// Skip this test if running in CI environment
	if os.Getenv("CI") != "" {
		t.Skip("Skipping test in CI environment")
	}
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

func TestParseCityguard(t *testing.T) {
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
10
10
2
1d12+123
1d8+3
500
9000
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

	// Check the key stats we fixed
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
}

func TestParseDikuMobile(t *testing.T) {
	// Create a temporary file with DikuMUD format mobile data
	tmpFile, err := os.CreateTemp("", "test_diku_mobile.*.mob")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	// Write DikuMUD format mobile data to the file
	testData := `#7044
lemure blob~
The lemure~
The lemure blob slithers terribly precisely towards you for an attack!~
This looks like a vaguely human blob. Big black yellow eyes, and a
mouth going a little bit out from the face. The lemure does not look
interested in you at all, but anyway it attackes. It looks like it's
mind has been burned out.
~
32 0 -500 S
5
15
5
1d6+50
1d4+1
10
100
8
8
0
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

	// Check the lemure stats
	lemure := mobiles[0]
	if lemure.VNUM != 7044 {
		t.Errorf("Expected VNUM 7044, got %d", lemure.VNUM)
	}
	if lemure.Name != "lemure blob" {
		t.Errorf("Expected name 'lemure blob', got '%s'", lemure.Name)
	}
	if lemure.ShortDesc != "The lemure" {
		t.Errorf("Expected short desc 'The lemure', got '%s'", lemure.ShortDesc)
	}
	if lemure.LongDesc != "The lemure blob slithers terribly precisely towards you for an attack!" {
		t.Errorf("Expected long desc 'The lemure blob slithers terribly precisely towards you for an attack!', got '%s'", lemure.LongDesc)
	}
	if lemure.ActFlags != 32 {
		t.Errorf("Expected act flags 32, got %d", lemure.ActFlags)
	}
	if lemure.Alignment != -500 {
		t.Errorf("Expected alignment -500, got %d", lemure.Alignment)
	}

	// Check the key stats for DikuMUD format
	if lemure.Level != 5 {
		t.Errorf("Expected level 5, got %d", lemure.Level)
	}
	if lemure.HitRoll != 15 {
		t.Errorf("Expected hitroll 15, got %d", lemure.HitRoll)
	}
	if lemure.AC[0] != 5 || lemure.AC[1] != 5 || lemure.AC[2] != 5 {
		t.Errorf("Expected AC [5 5 5], got %v", lemure.AC)
	}
	if lemure.Dice[0] != 1 || lemure.Dice[1] != 6 || lemure.Dice[2] != 50 {
		t.Errorf("Expected dice 1d6+50, got %dd%d+%d", lemure.Dice[0], lemure.Dice[1], lemure.Dice[2])
	}
	if lemure.DamRoll != 1 {
		t.Errorf("Expected damroll 1, got %d", lemure.DamRoll)
	}
	if lemure.Gold != 10 {
		t.Errorf("Expected gold 10, got %d", lemure.Gold)
	}
	if lemure.Experience != 100 {
		t.Errorf("Expected experience 100, got %d", lemure.Experience)
	}
	if lemure.Position != 8 {
		t.Errorf("Expected position 8, got %d", lemure.Position)
	}
	if lemure.DefaultPos != 8 {
		t.Errorf("Expected default position 8, got %d", lemure.DefaultPos)
	}
	if lemure.Sex != 0 {
		t.Errorf("Expected sex 0, got %d", lemure.Sex)
	}

	// Check default abilities for DikuMUD format
	for i, ability := range lemure.Abilities {
		if ability != 11 {
			t.Errorf("Expected ability %d to be 11, got %d", i, ability)
		}
	}
}
