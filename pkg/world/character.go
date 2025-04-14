package world

import (
	"time"

	"github.com/wltechblog/DikuGo/pkg/types"
)

// NewCharacter creates a new character
func NewCharacter(name string, isNPC bool) *types.Character {
	character := &types.Character{
		Name:          name,
		IsNPC:         isNPC,
		Position:      types.POS_STANDING,
		Level:         1,
		Gold:          0,
		Experience:    0,
		Alignment:     0,
		HP:            20,
		MaxHitPoints:  20,
		ManaPoints:    100,
		MaxManaPoints: 100,
		MovePoints:    100,
		MaxMovePoints: 100,
		ArmorClass:    [3]int{100, 100, 100},
		HitRoll:       0,
		DamRoll:       0,
		Abilities:     [6]int{10, 10, 10, 10, 10, 10},
		Skills:        make(map[int]int),
		Spells:        make(map[int]int),
		Equipment:     make([]*types.ObjectInstance, types.NUM_WEARS),
		Inventory:     make([]*types.ObjectInstance, 0),
		LastLogin:     time.Now(),
		Title:         " the newbie",
		Prompt:        "%h/%H hp %m/%M mana %v/%V mv> ",
	}

	return character
}

// NewNPC creates a new NPC from a mobile prototype
func NewNPC(prototype *types.Mobile) *types.Character {
	character := &types.Character{
		Name:          prototype.Name,
		ShortDesc:     prototype.ShortDesc,
		LongDesc:      prototype.LongDesc,
		Description:   prototype.Description,
		IsNPC:         true,
		Level:         prototype.Level,
		Sex:           prototype.Sex,
		Class:         prototype.Class,
		Race:          prototype.Race,
		Position:      prototype.Position,
		Gold:          prototype.Gold,
		Experience:    prototype.Experience,
		Alignment:     prototype.Alignment,
		HP:            prototype.Dice[0] * prototype.Dice[1] + prototype.Dice[2],
		MaxHitPoints:  prototype.Dice[0] * prototype.Dice[1] + prototype.Dice[2],
		ArmorClass:    prototype.AC,
		HitRoll:       prototype.HitRoll,
		DamRoll:       prototype.DamRoll,
		Prototype:     prototype,
		Equipment:     make([]*types.ObjectInstance, types.NUM_WEARS),
		Inventory:     make([]*types.ObjectInstance, 0),
		Flags:         prototype.ActFlags,
	}

	return character
}
