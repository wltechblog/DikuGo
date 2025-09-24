package world

import (
	"time"

	"github.com/wltechblog/DikuGo/pkg/types"
	"github.com/wltechblog/DikuGo/pkg/utils"
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
		LastSkillTime: make(map[int]time.Time),
		LastLogin:     time.Now(),
		Title:         " the newbie",
		Prompt:        "%h/%H hp %m/%M mana %v/%V mv> ",
	}

	return character
}

// NewNPC creates a new NPC from a mobile prototype
func NewNPC(prototype *types.Mobile) *types.Character {
	// Calculate HP based on dice exactly as in original DikuMUD
	// The format is XdY+Z where X is number of dice, Y is size of dice, Z is bonus
	baseHP := utils.Dice(prototype.Dice[0], prototype.Dice[1]) + prototype.Dice[2]

	// Calculate mana and move points based on level
	baseMana := prototype.Level * 10
	baseMove := prototype.Level * 10

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
		Gold:          prototype.Gold, // Gold is directly copied from prototype as in original DikuMUD
		Experience:    prototype.Experience,
		Alignment:     prototype.Alignment,
		HP:            baseHP,
		MaxHitPoints:  baseHP,
		ManaPoints:    baseMana,
		MaxManaPoints: baseMana,
		MovePoints:    baseMove,
		MaxMovePoints: baseMove,
		ArmorClass:    prototype.AC,
		HitRoll:       prototype.HitRoll,
		DamRoll:       prototype.DamRoll,
		Prototype:     prototype,
		Equipment:     make([]*types.ObjectInstance, types.NUM_WEARS),
		Inventory:     make([]*types.ObjectInstance, 0),
		Skills:        make(map[int]int),
		LastSkillTime: make(map[int]time.Time),
		Flags:         prototype.ActFlags,
	}

	// Create a deep copy of the abilities array
	character.Abilities = [6]int{
		prototype.Abilities[0],
		prototype.Abilities[1],
		prototype.Abilities[2],
		prototype.Abilities[3],
		prototype.Abilities[4],
		prototype.Abilities[5],
	}

	return character
}
