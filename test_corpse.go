package main

import (
	"fmt"
	"log"
	"time"

	"github.com/wltechblog/DikuGo/pkg/combat"
	"github.com/wltechblog/DikuGo/pkg/config"
	"github.com/wltechblog/DikuGo/pkg/storage"
	"github.com/wltechblog/DikuGo/pkg/types"
	"github.com/wltechblog/DikuGo/pkg/world"
)

func main() {
	// Create a new config
	cfg := &config.Config{}
	cfg.Game.DataPath = "old/lib"
	cfg.Storage.PlayerDir = "data/players"

	// Create a new storage
	fs, err := storage.NewFileStorage(cfg)
	if err != nil {
		log.Fatalf("Failed to create storage: %v", err)
	}

	// Create a new world
	w, err := world.NewWorld(cfg, fs)
	if err != nil {
		log.Fatalf("Failed to create world: %v", err)
	}

	// Get a room
	room := w.GetRoom(3001) // Temple Square
	if room == nil {
		log.Fatalf("Room 3001 not found")
	}

	// Print room info
	log.Printf("Using room %d: %s", room.VNUM, room.Name)

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
		DamRoll:      15,                             // Higher damage to ensure mob dies quickly
		Abilities:    [6]int{15, 10, 10, 10, 10, 10}, // STR, INT, WIS, DEX, CON, CHA
		Sex:          types.SEX_MALE,
		Equipment:    make([]*types.ObjectInstance, types.NUM_WEARS),
		World:        w,
	}

	// Add player to room directly instead of using AddCharacter
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
		World:        w,
		Gold:         100, // Give the mob some gold
	}

	// Add mob to room directly
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

	// Create an item for the mob
	item := &types.ObjectInstance{
		Prototype: &types.Object{
			VNUM:      2002,
			Name:      "ring",
			ShortDesc: "a gold ring",
			Type:      types.ITEM_TREASURE,
		},
	}
	mob.Inventory = append(mob.Inventory, item)
	item.CarriedBy = mob

	// Create a combat manager
	combatManager := combat.NewDikuCombatManager()

	// Start combat
	fmt.Println("Starting combat between", player.Name, "and", mob.Name)
	err = combatManager.StartCombat(player, mob)
	if err != nil {
		log.Fatalf("Failed to start combat: %v", err)
	}

	// Run several combat rounds until the mob dies
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

	// Check for corpse
	fmt.Println("\n--- Room Contents ---")
	for _, obj := range room.Objects {
		fmt.Printf("Object: %s\n", obj.Prototype.ShortDesc)
		if obj.Prototype.Type == types.ITEM_CONTAINER && obj.Prototype.Value[3] == 1 {
			fmt.Println("  This is a corpse!")
			fmt.Printf("  Timer: %d\n", obj.Timer)
			fmt.Println("  Contents:")
			for _, item := range obj.Contains {
				fmt.Printf("    %s\n", item.Prototype.ShortDesc)
			}
		}
	}

	// Simulate corpse decay
	fmt.Println("\n--- Simulating Corpse Decay ---")
	for i := 0; i < 3; i++ {
		fmt.Printf("\nDecay pulse %d\n", i+1)
		w.PulseCorpses()

		// Check room contents after each pulse
		fmt.Println("Room contents:")
		for _, obj := range room.Objects {
			if obj.Prototype.Type == types.ITEM_CONTAINER && obj.Prototype.Value[3] == 1 {
				fmt.Printf("Corpse (Timer: %d)\n", obj.Timer)
			} else {
				fmt.Printf("%s\n", obj.Prototype.ShortDesc)
			}
		}

		time.Sleep(1 * time.Second)
	}
}
