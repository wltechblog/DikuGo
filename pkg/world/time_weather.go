package world

import (
	"log"
)

// PulseTime updates the game time
func (w *World) PulseTime() {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	// Increment the hour
	w.time.Hour++

	// Check if we've reached the end of the day
	if w.time.Hour >= 24 {
		w.time.Hour = 0
		w.time.Day++

		// Check if we've reached the end of the month
		if w.time.Day > 30 {
			w.time.Day = 1
			w.time.Month++

			// Check if we've reached the end of the year
			if w.time.Month > 12 {
				w.time.Month = 1
				w.time.Year++
			}
		}
	}

	log.Printf("Game time: %d:%02d, Day %d, Month %d, Year %d",
		w.time.Hour, 0, w.time.Day, w.time.Month, w.time.Year)
}

// PulseWeather updates the game weather
func (w *World) PulseWeather() {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	// TODO: Implement weather changes
	// For now, just log the current weather
	log.Printf("Game weather: %d", w.time.Weather)
}
