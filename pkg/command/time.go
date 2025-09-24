package command

import (
	"fmt"

	"github.com/wltechblog/DikuGo/pkg/types"
)

// TimeCommand represents the time command
type TimeCommand struct{}

// Execute executes the time command
func (c *TimeCommand) Execute(character *types.Character, args string) error {
	// Get the world from the character
	world, ok := character.World.(interface {
		GetTimeString() string
		GetWeekdayName() string
		GetMonthName() string
	})
	if !ok {
		return fmt.Errorf("world interface not available")
	}

	// Get the time string
	timeStr := world.GetTimeString()
	weekdayName := world.GetWeekdayName()
	monthName := world.GetMonthName()

	// Get the time info from the world
	timeInfo, ok := character.World.(interface {
		GetTimeInfo() (hours, day, month, year int)
	})
	if !ok {
		return fmt.Errorf("time info not available")
	}

	// Get the time info
	_, day, _, year := timeInfo.GetTimeInfo()

	// Format the day with the appropriate suffix
	daySuffix := "th"
	dayNum := day + 1 // Convert from 0-based to 1-based
	if dayNum == 1 || (dayNum > 20 && dayNum%10 == 1) {
		daySuffix = "st"
	} else if dayNum == 2 || (dayNum > 20 && dayNum%10 == 2) {
		daySuffix = "nd"
	} else if dayNum == 3 || (dayNum > 20 && dayNum%10 == 3) {
		daySuffix = "rd"
	}

	// Send the time to the character
	character.SendMessage(fmt.Sprintf("It is %s, on %s.\r\n", timeStr, weekdayName))
	character.SendMessage(fmt.Sprintf("The %d%s Day of the %s, Year %d.\r\n",
		dayNum, daySuffix, monthName, year))

	return nil
}

// Name returns the name of the command
func (c *TimeCommand) Name() string {
	return "time"
}

// Aliases returns the aliases of the command
func (c *TimeCommand) Aliases() []string {
	return []string{}
}

// MinPosition returns the minimum position required to use the command
func (c *TimeCommand) MinPosition() int {
	return types.POS_RESTING
}

// Level returns the required level to use the command
func (c *TimeCommand) Level() int {
	return 0
}

// LogCommand returns whether the command should be logged
func (c *TimeCommand) LogCommand() bool {
	return false
}
