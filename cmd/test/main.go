package main

import (
	"log"

	"github.com/wltechblog/DikuGo/pkg/types"
)

func main() {
	// Parse the mobile file using our custom parser
	mobiles, err := parseMobiles("old/lib/test.mob")
	if err != nil {
		log.Fatalf("Error parsing mobiles: %v", err)
	}

	log.Printf("Parsed %d mobile prototypes", len(mobiles))

	// Create a room
	room := &types.Room{
		VNUM:        1,
		Name:        "Test Room",
		Description: "This is a test room.",
		Characters:  make([]*types.Character, 0),
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
}
