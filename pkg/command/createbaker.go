package command

import (
	"fmt"

	"github.com/wltechblog/DikuGo/pkg/types"
)

// CreateBakerCommand represents the createbaker command
type CreateBakerCommand struct{}

// Name returns the name of the command
func (c *CreateBakerCommand) Name() string {
	return "createbaker"
}

// Aliases returns the aliases of the command
func (c *CreateBakerCommand) Aliases() []string {
	return []string{}
}

// MinPosition returns the minimum position required to execute the command
func (c *CreateBakerCommand) MinPosition() int {
	return types.POS_DEAD
}

// Level returns the minimum level required to execute the command
func (c *CreateBakerCommand) Level() int {
	return 1 // Available to everyone during development
}

// LogCommand returns whether the command should be logged
func (c *CreateBakerCommand) LogCommand() bool {
	return true
}

// Execute executes the createbaker command
func (c *CreateBakerCommand) Execute(ch *types.Character, args string) error {
	// Get the world from the character
	world, ok := ch.World.(interface{
		AddMobilePrototype(*types.Mobile)
	})
	if !ok {
		return fmt.Errorf("world interface not available")
	}

	// Create the baker mobile prototype
	baker := &types.Mobile{
		VNUM:        3001,
		Name:        "baker",
		ShortDesc:   "the Baker",
		LongDesc:    "The Baker looks at you calmly, wiping flour from his face with one hand.",
		Description: "A fat, nice looking baker. But you can see that he has many scars on his body.",
		Level:       23,
		HitRoll:     18,
		DamRoll:     12,
		AC:          [3]int{2, 2, 2},
		Gold:        2000,
		Experience:  80000,
		Position:    8,
		DefaultPos:  8,
		Sex:         1,
		ActFlags:    2,
		AffectFlags: 0,
		Alignment:   900,
		Abilities:   [6]int{11, 11, 11, 11, 11, 11},
		Dice:        [3]int{6, 10, 390},
	}

	// Add the baker to the world
	world.AddMobilePrototype(baker)

	return fmt.Errorf("Baker (VNUM 3001) has been created and added to the world.")
}
