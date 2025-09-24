package main

import (
	"log"
	"path/filepath"

	"github.com/wltechblog/DikuGo/pkg/config"
	"github.com/wltechblog/DikuGo/pkg/storage"
	"github.com/wltechblog/DikuGo/pkg/world"
)

func main() {
	// Create a new config
	cfg := &config.Config{}
	cfg.Game.DataPath = "old/lib"
	cfg.Storage.PlayerDir = "old/lib/players"

	// Create a new storage
	fs, err := storage.NewFileStorage(cfg)
	if err != nil {
		log.Fatalf("Error creating file storage: %v", err)
	}

	// Create a new world
	w, _ := world.NewWorld(cfg, fs)

	// Load test mobile
	mobiles, err := storage.ParseMobiles(filepath.Join(cfg.Game.DataPath, "test.mob"))
	if err != nil {
		log.Fatalf("Error loading mobiles: %v", err)
	}

	log.Printf("Loaded %d mobile prototypes", len(mobiles))

	// Add mobiles to world
	for _, mobile := range mobiles {
		w.AddMobile(mobile)
	}

	// Get a room
	room := w.GetRoom(3001) // Temple Square
	if room == nil {
		log.Fatalf("Room 3001 not found")
	}

	// Create a character from the mobile prototype
	if len(mobiles) > 0 {
		mobile := mobiles[0]
		// Use the CreateMobFromPrototype function to create a character with all the correct stats
		character := w.CreateMobFromPrototype(mobile.VNUM, room)
		if character == nil {
			log.Fatalf("Failed to create character from prototype")
		}

		log.Printf("Created character %s in room %d with stats: Level=%d, HitRoll=%d, DamRoll=%d, AC=%v, Gold=%d, Exp=%d",
			character.Name, room.VNUM, character.Level, character.HitRoll, character.DamRoll, character.ArmorClass, character.Gold, character.Experience)
	}

	// Initialize AI
	w.InitAI()

	// Tick AI a few times
	for i := 0; i < 10; i++ {
		log.Printf("Tick %d", i)
		w.TickAI()
	}
}
