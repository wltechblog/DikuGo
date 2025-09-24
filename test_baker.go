package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/wltechblog/DikuGo/pkg/storage"
)

func main() {
	// Parse the mobile file directly
	mobiles, err := storage.ParseMobiles(filepath.Join("old/lib", "tinyworld.mob"))
	if err != nil {
		log.Fatalf("Failed to parse mobile file: %v", err)
	}

	// Print all mobiles
	fmt.Printf("Loaded %d mobiles\n", len(mobiles))
	
	// Check for baker (VNUM 3001)
	for _, mob := range mobiles {
		if mob.VNUM == 3001 {
			fmt.Printf("Found baker (VNUM 3001):\n")
			fmt.Printf("  Name: %s\n", mob.Name)
			fmt.Printf("  ShortDesc: %s\n", mob.ShortDesc)
			fmt.Printf("  Level: %d\n", mob.Level)
			fmt.Printf("  HitRoll: %d\n", mob.HitRoll)
			fmt.Printf("  DamRoll: %d\n", mob.DamRoll)
			fmt.Printf("  AC: %v\n", mob.AC)
			fmt.Printf("  Gold: %d\n", mob.Gold)
			fmt.Printf("  Experience: %d\n", mob.Experience)
			os.Exit(0)
		}
	}
	
	fmt.Println("Baker (VNUM 3001) not found!")
	os.Exit(1)
}
