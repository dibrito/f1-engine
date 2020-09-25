package service

import (
	"fmt"
	"sort"
	"time"

	"github.com/dibrito/f1-engine/src/domain"
)

// GetRaceFinalResult should process race LapsMapper and give final result
func GetRaceFinalResult(result domain.LapsMapper) {
	raceResultM := ProcessFinalResult(result)
	// TODO maybe pass a pointer
	raceResultM = OrderRaceResult(raceResultM)
	ShowRaceResult(raceResultM)
}

func ShowRaceResult(result domain.RaceResultMapper) {
	fmt.Println("POSITION \t\t\t\t PILOT \t\t\t\t TOTAL LAPS \t\t\t\t TOTAL RACE TIME")
	for i := 1; i < len(result.Result)+1; i++ {
		v := result.Result[i]
		fmt.Printf("%v \t\t\t\t %d - %s \t\t\t\t %d \t\t\t\t %v\t\t\t\t \n",
			v.FinalPosition, v.PilotCode, v.PilotName, v.CompletedLaps, v.TotalRaceTime)
	}
}

// GetRaceDuration sums all Laps duration
func GetRaceDuration(laps []domain.LapResult) time.Duration {
	var raceDuration time.Duration

	for _, l := range laps {
		raceDuration += l.Duration
	}

	return raceDuration
}

// GetAvgSpeed get avg speed
func GetAvgSpeed(laps []domain.LapResult) float64 {
	avgSpeed := 0.0

	for _, l := range laps {
		avgSpeed += l.AvgSpeed
	}

	return avgSpeed / float64(len(laps))
}

// GetTotalCompletedLaps sums all Laps
func GetTotalCompletedLaps(laps []domain.LapResult) int {
	return len(laps)
}

// OrderRaceResult sets final position
func OrderRaceResult(result domain.RaceResultMapper) domain.RaceResultMapper {
	var orderedResult domain.RaceResultMapper
	orderedResult.Result = map[int]domain.RaceResult{}
	m := map[int]time.Duration{}

	for _, v := range result.Result {
		m[v.PilotCode] = v.TotalRaceTime
	}

	type kv struct {
		Key   int
		Value time.Duration
	}

	var ss []kv
	for k, v := range m {
		ss = append(ss, kv{k, v})
	}

	sort.Slice(ss, func(i, j int) bool {
		return ss[i].Value < ss[j].Value
	})

	for i, _ := range ss {
		pilotCode := ss[i].Key
		orderedResult.Result[i+1] = domain.RaceResult{
			FinalPosition: i + 1,
			PilotCode:     pilotCode,
			PilotName:     result.Result[pilotCode].PilotName,
			CompletedLaps: result.Result[pilotCode].CompletedLaps,
			TotalRaceTime: result.Result[pilotCode].TotalRaceTime,
		}
	}

	return orderedResult
}

// ProcessFinalResult sums all Laps
func ProcessFinalResult(mapper domain.LapsMapper) domain.RaceResultMapper {
	var finalResult domain.RaceResultMapper
	finalResult.Result = map[int]domain.RaceResult{}
	for _, laps := range mapper.Result {
		pilotCode := laps[0].PilotCode
		pilotName := laps[0].PilotName
		totalRaceTime := GetRaceDuration(laps)
		totalCompletedLaps := GetTotalCompletedLaps(laps)
		agvSpeeed := GetAvgSpeed(laps)

		if _, ok := finalResult.Result[pilotCode]; !ok {
			finalResult.Result[pilotCode] = domain.RaceResult{
				FinalPosition: 0,
				PilotCode:     pilotCode,
				PilotName:     pilotName,
				CompletedLaps: totalCompletedLaps,
				TotalRaceTime: totalRaceTime,
				AvgSpeed:      agvSpeeed,
			}
		}
	}

	return finalResult
}
