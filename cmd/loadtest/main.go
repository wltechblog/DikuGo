package main

import (
	"log"
	"path/filepath"

	"github.com/wltechblog/DikuGo/pkg/config"
	"github.com/wltechblog/DikuGo/pkg/storage"
	"github.com/wltechblog/DikuGo/pkg/types"
	"github.com/wltechblog/DikuGo/pkg/world"
)

func main() {
	// Create a new config
	cfg := &config.Config{
		DataDir: "old/lib",
	}

	// Create a new storage
	fs := storage.NewFileStorage(cfg)

	// Create a new world
	w, _ := world.NewWorld(cfg, fs)

	// Load test mobile
	mobiles, err := storage.ParseMobiles(filepath.Join(cfg.DataDir, "test.mob"))
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
		character := &types.Character{
			Name:        mobile.Name,
			ShortDesc:   mobile.ShortDesc,
			LongDesc:    mobile.LongDesc,
			Description: mobile.Description,
			IsNPC:       true,
			InRoom:      room,
			RoomVNUM:    room.VNUM,
			Prototype:   mobile,
		}

		// Add character to room
		room.Characters = append(room.Characters, character)

		log.Printf("Created character %s in room %d", character.Name, room.VNUM)
	}

	// Initialize AI
	w.InitAI()

	// Tick AI a few times
	for i := 0; i < 10; i++ {
		log.Printf("Tick %d", i)
		w.TickAI()
	}
}
