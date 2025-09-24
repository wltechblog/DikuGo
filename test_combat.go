package main

import (
	"fmt"
	"log"
	"time"

	"github.com/wltechblog/DikuGo/pkg/combat"
	"github.com/wltechblog/DikuGo/pkg/types"
)

func main() {
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
		ShortDesc:    "a test player",
		IsNPC:        false,
		InRoom:       room,
		Level:        5,
		Position:     types.POS_STANDING,
		HP:           50,
		MaxHitPoints: 50,
		HitRoll:      5,
		DamRoll:      5,
		Abilities:    [6]int{15, 10, 10, 10, 10, 10}, // STR, INT, WIS, DEX, CON, CHA
		Sex:          types.SEX_MALE,
		Equipment:    make([]*types.ObjectInstance, types.NUM_WEARS),
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
		Abilities:    [6]int{12, 10, 10, 10, 10, 10}, // STR, INT, WIS, DEX, CON, CHA
		Sex:          types.SEX_MALE,
		Equipment:    make([]*types.ObjectInstance, types.NUM_WEARS),
	}
	room.Characters = append(room.Characters, mob)

	// Create a weapon for the player
	weapon := &types.ObjectInstance{
		Prototype: &types.Object{
			VNUM:      2001,
			Name:      "sword",
			ShortDesc: "a sharp sword",
			Value:     [4]int{0, 2, 4, types.TYPE_SLASH}, // 2d4 slashing damage
		},
		WornBy: player,
		WornOn: types.WEAR_WIELD,
	}
	player.Equipment[types.WEAR_WIELD] = weapon

	// Create a combat manager
	combatManager := combat.NewEnhancedDikuCombatManager()

	// Start combat
	fmt.Println("Starting combat between", player.Name, "and", mob.Name)
	err := combatManager.StartCombat(player, mob)
	if err != nil {
		log.Fatalf("Failed to start combat: %v", err)
	}

	// Run several combat rounds
	for i := 0; i < 10; i++ {
		fmt.Printf("\n--- Round %d ---\n", i+1)
		fmt.Printf("%s HP: %d/%d\n", player.Name, player.HP, player.MaxHitPoints)
		fmt.Printf("%s HP: %d/%d\n", mob.Name, mob.HP, mob.MaxHitPoints)

		// Update combat
		combatManager.Update()

		// Check if combat is over
		if player.Fighting == nil || mob.Fighting == nil {
			fmt.Println("Combat has ended!")
			break
		}

		// Sleep for a bit to simulate time passing
		time.Sleep(500 * time.Millisecond)
	}

	// Print final status
	fmt.Printf("\n--- Final Status ---\n")
	fmt.Printf("%s HP: %d/%d\n", player.Name, player.HP, player.MaxHitPoints)
	fmt.Printf("%s HP: %d/%d\n", mob.Name, mob.HP, mob.MaxHitPoints)

	if player.HP <= 0 {
		fmt.Printf("%s has been defeated!\n", player.Name)
	}
	if mob.HP <= 0 {
		fmt.Printf("%s has been defeated!\n", mob.Name)
	}
}
