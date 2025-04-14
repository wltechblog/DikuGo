package world

import (
	"math/rand"
	"testing"

	"github.com/wltechblog/DikuGo/pkg/types"
)

// performTestCombatRound simulates a combat round between attacker and victim
func performTestCombatRound(attacker, victim *types.Character) {
	// Check if the attacker is still fighting
	if attacker.Fighting == nil || attacker.Position != types.POS_STANDING {
		return
	}

	// Check if the victim is still valid
	if victim.Position == types.POS_DEAD {
		attacker.Fighting = nil
		return
	}

	// Calculate hit chance
	hitChance := 50 + attacker.HitRoll
	if hitChance < 5 {
		hitChance = 5
	} else if hitChance > 95 {
		hitChance = 95
	}

	// Roll to hit
	roll := rand.Intn(100) + 1
	if roll <= hitChance {
		// Hit! Calculate damage
		damage := rand.Intn(attacker.DamRoll+1) + 1
		if damage < 1 {
			damage = 1
		}

		// Apply damage
		victim.HP -= damage

		// Check if the victim died
		if victim.HP <= 0 {
			victim.HP = 0
			victim.Position = types.POS_DEAD
			victim.Fighting = nil

			// Stop anyone fighting the victim
			if attacker.Fighting == victim {
				attacker.Fighting = nil
			}
		}
	}
}

func TestCombat(t *testing.T) {
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
		Level:        5,
		Position:     types.POS_STANDING,
		HP:           50,
		MaxHitPoints: 50,
		HitRoll:      5,
		DamRoll:      5,
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
		Level:        3,
		Position:     types.POS_STANDING,
		HP:           30,
		MaxHitPoints: 30,
		HitRoll:      3,
		DamRoll:      3,
	}
	room.Characters = append(room.Characters, mob)

	// Start combat
	player.Fighting = mob
	mob.Fighting = player

	// Run several combat rounds directly
	for i := 0; i < 10; i++ {
		// Simulate combat round for player
		performTestCombatRound(player, mob)

		// Simulate combat round for mob
		if mob.Position == types.POS_STANDING {
			performTestCombatRound(mob, player)
		}

		// If either character is dead, stop the test
		if player.Position == types.POS_DEAD || mob.Position == types.POS_DEAD {
			break
		}
	}

	// Check if combat had an effect
	if player.HP == 50 && mob.HP == 30 {
		t.Errorf("Expected combat to reduce hit points, but no damage was done")
	}

	// Check if the characters are still fighting
	if player.Position != types.POS_DEAD && mob.Position != types.POS_DEAD {
		if player.Fighting == nil {
			t.Errorf("Expected player to still be fighting")
		}
		if mob.Fighting == nil {
			t.Errorf("Expected mob to still be fighting")
		}
	}

	// If one of the characters died, check if combat was properly ended
	if player.Position == types.POS_DEAD {
		if player.Fighting != nil {
			t.Errorf("Expected dead player to not be fighting")
		}
		if mob.Fighting != nil {
			t.Errorf("Expected mob to stop fighting dead player")
		}
	}

	if mob.Position == types.POS_DEAD {
		if mob.Fighting != nil {
			t.Errorf("Expected dead mob to not be fighting")
		}
		if player.Fighting != nil {
			t.Errorf("Expected player to stop fighting dead mob")
		}
	}
}
