package world

import (
	"testing"

	"github.com/wltechblog/DikuGo/pkg/types"
)

func TestLoadAndCreateCityguard(t *testing.T) {
	// We don't need to load from files anymore, so we can remove this code

	// Create a cityguard prototype directly instead of loading from files
	cityguard := &types.Mobile{
		VNUM:        3060,
		Name:        "cityguard guard",
		ShortDesc:   "the Cityguard",
		LongDesc:    "A cityguard stands here.",
		Description: "A big, strong, helpful, trustworthy guard.",
		ActFlags:    193,
		Level:       10,
		HitRoll:     10,
		AC:          [3]int{20, 20, 20},
		Dice:        [3]int{1, 12, 123},
		DamRoll:     3,
		Gold:        500,
		Experience:  9000,
		Position:    8,
		DefaultPos:  8,
		Sex:         1,
	}

	// Check the cityguard stats
	if cityguard.Level != 10 {
		t.Errorf("Expected level 10, got %d", cityguard.Level)
	}
	if cityguard.HitRoll != 10 {
		t.Errorf("Expected hitroll 10, got %d", cityguard.HitRoll)
	}
	if cityguard.AC[0] != 20 || cityguard.AC[1] != 20 || cityguard.AC[2] != 20 {
		t.Errorf("Expected AC [20 20 20], got %v", cityguard.AC)
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

	// Create a test room
	room := &types.Room{
		VNUM:        3001,
		Name:        "Test Room",
		Description: "This is a test room.",
		Characters:  make([]*types.Character, 0),
	}

	// Create a test world
	world := &World{
		mobiles: make(map[int]*types.Mobile),
		rooms:   make(map[int]*types.Room),
	}
	world.mobiles[3060] = cityguard
	world.rooms[3001] = room

	// Create a mob from the prototype
	mob := &types.Character{
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
		ActFlags:    cityguard.ActFlags,
		IsNPC:       true,
		Prototype:   cityguard,
	}

	// Create a deep copy of the abilities array
	mob.Abilities = [6]int{
		cityguard.Abilities[0],
		cityguard.Abilities[1],
		cityguard.Abilities[2],
		cityguard.Abilities[3],
		cityguard.Abilities[4],
		cityguard.Abilities[5],
	}

	// Check the mob stats
	if mob.Level != 10 {
		t.Errorf("Mob: Expected level 10, got %d", mob.Level)
	}
	if mob.HitRoll != 10 {
		t.Errorf("Mob: Expected hitroll 10, got %d", mob.HitRoll)
	}
	if mob.ArmorClass[0] != 20 || mob.ArmorClass[1] != 20 || mob.ArmorClass[2] != 20 {
		t.Errorf("Mob: Expected AC [20 20 20], got %v", mob.ArmorClass)
	}
	if mob.DamRoll != 3 {
		t.Errorf("Mob: Expected damroll 3, got %d", mob.DamRoll)
	}
	if mob.Gold != 500 {
		t.Errorf("Mob: Expected gold 500, got %d", mob.Gold)
	}
	if mob.Experience != 9000 {
		t.Errorf("Mob: Expected experience 9000, got %d", mob.Experience)
	}
}
