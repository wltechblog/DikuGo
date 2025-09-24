package command

import (
	"fmt"
	"log"

	"github.com/wltechblog/DikuGo/pkg/types"
)

// AddBakerCommand represents the addbaker command
type AddBakerCommand struct{}

// Name returns the name of the command
func (c *AddBakerCommand) Name() string {
	return "addbaker"
}

// Aliases returns the aliases of the command
func (c *AddBakerCommand) Aliases() []string {
	return []string{}
}

// MinPosition returns the minimum position required to execute the command
func (c *AddBakerCommand) MinPosition() int {
	return types.POS_DEAD
}

// Level returns the minimum level required to execute the command
func (c *AddBakerCommand) Level() int {
	return 1 // Available to everyone during development
}

// LogCommand returns whether the command should be logged
func (c *AddBakerCommand) LogCommand() bool {
	return true
}

// Execute executes the addbaker command
func (c *AddBakerCommand) Execute(ch *types.Character, args string) error {
	// Get the world from the character
	world, ok := ch.World.(interface{
		AddMobilePrototype(*types.Mobile)
		GetMobile(int) *types.Mobile
	})
	if !ok {
		return fmt.Errorf("world interface not available")
	}

	// Check if the baker already exists
	existingBaker := world.GetMobile(3001)
	if existingBaker != nil {
		return fmt.Errorf("Baker (VNUM 3001) already exists in the world.")
	}

	// Create the baker mobile prototype
	baker := &types.Mobile{
		VNUM:        3001,
		Name:        "baker",
		ShortDesc:   "the Baker",
		LongDesc:    "The Baker looks at you calmly, wiping flour from his face with one hand.\n",
		Description: "A fat, nice looking baker. But you can see that he has many scars on his body.\n",
		Level:       23,
		HitRoll:     18,
		DamRoll:     12,
		AC:          [3]int{2, 2, 2},
		Gold:        2000,
		Experience:  80000,
		Position:    8,
		DefaultPos:  8,
		Sex:         1,
		ActFlags:    10, // NPC + SENTINEL
		AffectFlags: 0,
		Alignment:   900,
		Abilities:   [6]int{18, 18, 18, 18, 18, 18},
		Dice:        [3]int{6, 10, 390},
	}

	// Add the baker to the world
	world.AddMobilePrototype(baker)

	log.Printf("Added baker (VNUM 3001) to the world")
	return fmt.Errorf("Baker (VNUM 3001) has been created and added to the world.")
}
