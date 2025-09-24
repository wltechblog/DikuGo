package world

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/wltechblog/DikuGo/pkg/types"
)

// Weekday names
var weekdays = []string{
	"the Day of the Moon",
	"the Day of the Bull",
	"the Day of the Deception",
	"the Day of Thunder",
	"the Day of Freedom",
	"the day of the Great Gods",
	"the Day of the Sun",
}

// Month names
var monthNames = []string{
	"Month of Winter",
	"Month of the Winter Wolf",
	"Month of the Frost Giant",
	"Month of the Old Forces",
	"Month of the Grand Struggle",
	"Month of the Spring",
	"Month of Nature",
	"Month of Futility",
	"Month of the Dragon",
	"Month of the Sun",
	"Month of the Heat",
	"Month of the Battle",
	"Month of the Dark Shades",
	"Month of the Shadows",
	"Month of the Long Shadows",
	"Month of the Ancient Darkness",
	"Month of the Great Evil",
}

// Sky descriptions
var skyDescriptions = []string{
	"cloudless",
	"cloudy",
	"rainy",
	"lit by flashes of lightning",
}

// InitializeTime initializes the game time based on real time
func (w *World) InitializeTime() {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	// Use a fixed beginning of time (similar to DikuMUD's approach)
	beginningOfTime := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	now := time.Now()

	// Calculate seconds passed since beginning of time
	secondsPassed := now.Unix() - beginningOfTime.Unix()

	// Convert to mud time
	w.time.Hours = int((secondsPassed / types.SECS_PER_MUD_HOUR) % 24)
	w.time.Day = int((secondsPassed / types.SECS_PER_MUD_DAY) % types.DAYS_PER_MONTH)
	w.time.Month = int((secondsPassed / types.SECS_PER_MUD_MONTH) % types.MONTHS_PER_YEAR)
	w.time.Year = int(secondsPassed / types.SECS_PER_MUD_YEAR)

	// Initialize sunlight based on time of day
	if w.time.Hours >= 5 && w.time.Hours < 21 {
		if w.time.Hours == 5 {
			w.time.Sunlight = types.SUN_RISE
		} else if w.time.Hours == 21 {
			w.time.Sunlight = types.SUN_SET
		} else {
			w.time.Sunlight = types.SUN_LIGHT
		}
	} else {
		w.time.Sunlight = types.SUN_DARK
	}

	// Initialize weather (default to cloudless)
	w.time.Weather = types.SKY_CLOUDLESS
	w.time.Change = 0

	log.Printf("Game time initialized: %d:%02d, Day %d, Month %d, Year %d, Sunlight %d, Weather %d",
		w.time.Hours, 0, w.time.Day+1, w.time.Month+1, w.time.Year, w.time.Sunlight, w.time.Weather)
}

// PulseTime updates the game time
func (w *World) PulseTime() {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	// Increment the hour
	w.time.Hours++

	// Update sunlight based on time of day
	if w.time.Hours == 5 {
		w.time.Sunlight = types.SUN_RISE
		w.SendToOutdoor("The sun rises in the east.\r\n")
	} else if w.time.Hours == 6 {
		w.time.Sunlight = types.SUN_LIGHT
		w.SendToOutdoor("The day has begun.\r\n")
	} else if w.time.Hours == 21 {
		w.time.Sunlight = types.SUN_SET
		w.SendToOutdoor("The sun slowly disappears in the west.\r\n")
	} else if w.time.Hours == 22 {
		w.time.Sunlight = types.SUN_DARK
		w.SendToOutdoor("The night has begun.\r\n")
	}

	// Check if we've reached the end of the day
	if w.time.Hours >= 24 {
		w.time.Hours = 0
		w.time.Day++

		// Check if we've reached the end of the month
		if w.time.Day >= types.DAYS_PER_MONTH {
			w.time.Day = 0
			w.time.Month++

			// Check if we've reached the end of the year
			if w.time.Month >= types.MONTHS_PER_YEAR {
				w.time.Month = 0
				w.time.Year++
			}
		}
	}

	log.Printf("Game time: %d:%02d, Day %d, Month %d, Year %d, Sunlight %d",
		w.time.Hours, 0, w.time.Day+1, w.time.Month+1, w.time.Year, w.time.Sunlight)
}

// PulseWeather updates the game weather
func (w *World) PulseWeather() {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	// Change the weather based on pressure and season
	var diff int
	if w.time.Month >= 9 && w.time.Month <= 16 {
		// Winter
		if w.time.Change > 985 {
			diff = -2
		} else {
			diff = 2
		}
	} else {
		// Summer
		if w.time.Change > 1015 {
			diff = -2
		} else {
			diff = 2
		}
	}

	// Calculate change in pressure
	w.time.Change += diff * (rand.Intn(4) - 2) // -2 to +2 range

	// Constrain pressure
	if w.time.Change < 960 {
		w.time.Change = 960
	} else if w.time.Change > 1040 {
		w.time.Change = 1040
	}

	// Update weather based on pressure
	if w.time.Change < 990 {
		if w.time.Weather != types.SKY_LIGHTNING {
			w.SendToOutdoor("Lightning flashes in the sky.\r\n")
			w.time.Weather = types.SKY_LIGHTNING
		}
	} else if w.time.Change < 1010 {
		if w.time.Weather != types.SKY_RAINY {
			w.SendToOutdoor("It starts to rain.\r\n")
			w.time.Weather = types.SKY_RAINY
		}
	} else if w.time.Change < 1030 {
		if w.time.Weather != types.SKY_CLOUDY {
			w.SendToOutdoor("The sky is getting cloudy.\r\n")
			w.time.Weather = types.SKY_CLOUDY
		}
	} else {
		if w.time.Weather != types.SKY_CLOUDLESS {
			w.SendToOutdoor("The clouds disappear.\r\n")
			w.time.Weather = types.SKY_CLOUDLESS
		}
	}

	log.Printf("Game weather: %d (Change: %d)", w.time.Weather, w.time.Change)
}

// SendToOutdoor sends a message to all characters who are outside
func (w *World) SendToOutdoor(message string) {
	w.mutex.RLock()
	defer w.mutex.RUnlock()

	// Get all characters
	for _, ch := range w.characters {
		// Check if character is in a room
		if ch.InRoom != nil {
			// Check if the room is outdoors (not INDOORS flag)
			if ch.InRoom.Flags&types.ROOM_INDOORS == 0 {
				// Send message to character
				if w.messageHandler != nil {
					w.messageHandler(ch, message)
				}
			}
		}
	}
}

// GetTimeString returns a string representation of the current game time
func (w *World) GetTimeString() string {
	w.mutex.RLock()
	defer w.mutex.RUnlock()

	hour := w.time.Hours

	// Convert to 12-hour format with AM/PM
	suffix := "AM"
	if hour >= 12 {
		suffix = "PM"
		if hour > 12 {
			hour -= 12
		}
	}
	if hour == 0 {
		hour = 12
	}

	return fmt.Sprintf("%d%s", hour, suffix)
}

// GetWeekdayName returns the name of the current weekday
func (w *World) GetWeekdayName() string {
	w.mutex.RLock()
	defer w.mutex.RUnlock()

	// Calculate weekday (similar to original DikuMUD)
	weekday := ((types.DAYS_PER_MONTH * w.time.Month) + w.time.Day + 1) % types.DAYS_PER_WEEK
	return weekdays[weekday]
}

// GetMonthName returns the name of the current month
func (w *World) GetMonthName() string {
	w.mutex.RLock()
	defer w.mutex.RUnlock()

	return monthNames[w.time.Month]
}

// GetSkyDescription returns a description of the current sky
func (w *World) GetSkyDescription() string {
	w.mutex.RLock()
	defer w.mutex.RUnlock()

	return skyDescriptions[w.time.Weather]
}

// IsShopOpen checks if a shop is open based on the current time
func (w *World) IsShopOpen(shop *types.Shop) bool {
	w.mutex.RLock()
	defer w.mutex.RUnlock()

	// If no specific hours are set, shop is always open
	if shop.OpenHour == 0 && shop.CloseHour == 0 {
		return true
	}

	// Check if current hour is within shop hours
	if shop.OpenHour <= shop.CloseHour {
		// Normal hours (e.g., 9 to 17)
		return w.time.Hours >= shop.OpenHour && w.time.Hours < shop.CloseHour
	} else {
		// Overnight hours (e.g., 21 to 9)
		return w.time.Hours >= shop.OpenHour || w.time.Hours < shop.CloseHour
	}
}

// GetTimeInfo returns the current time information
func (w *World) GetTimeInfo() (hours, day, month, year int) {
	w.mutex.RLock()
	defer w.mutex.RUnlock()

	return w.time.Hours, w.time.Day, w.time.Month, w.time.Year
}
