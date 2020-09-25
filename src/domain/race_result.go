package domain

import (
	"time"
)

// LapsMapper defines the struct to hold alll pilot's laps
type LapsMapper struct {
	// the key should be the pilot code
	Result map[int][]LapResult
}

type LapResult struct {
	Time      time.Time
	PilotCode int
	PilotName string
	Number    int
	Duration  time.Duration
	AvgSpeed  float64
}


type RaceResultMapper struct {
	// the key should be the pilot code
	Result map[int]RaceResult
}

// RaceResult represents the final race Result for each pilot
type RaceResult struct {
	FinalPosition int
	PilotCode     int
	PilotName     string
	CompletedLaps int
	TotalRaceTime time.Duration
}
