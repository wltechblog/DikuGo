package main

import (
	"fmt"
	"log"
	"time"

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

	// Print initial room contents
	fmt.Println("\n--- Initial Room Contents ---")
	for _, obj := range room.Objects {
		fmt.Printf("Object: %s\n", obj.Prototype.ShortDesc)
	}

	// Kill the mob
	fmt.Println("\n--- Killing the mob ---")
	mob.HP = 0
	mob.Position = types.POS_DEAD

	// Create a corpse
	fmt.Println("\n--- Creating corpse ---")
	corpse := w.MakeCorpse(mob)
	if corpse != nil {
		fmt.Printf("Corpse created: %s\n", corpse.Prototype.ShortDesc)
	} else {
		fmt.Println("Failed to create corpse!")
	}

	// Print room contents after corpse creation
	fmt.Println("\n--- Room Contents After Corpse Creation ---")
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
