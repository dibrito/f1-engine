package service

import (
	"fmt"
	"sort"
	"time"

	"github.com/dibrito/f1-engine/src/domain"
)

// GetRaceDuration sums all Laps duration
func GetRaceDuration(laps []domain.LapResult) time.Duration {
	var raceDuration time.Duration

	for _, l := range laps {
		raceDuration += l.Duration
	}

	return raceDuration
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

	for _, kv := range ss {
		fmt.Printf("%d, %v\n", kv.Key, kv.Value)
	}

	for i, _ := range ss {
		pilotCode := ss[i].Key
		orderedResult.Result[pilotCode] = domain.RaceResult{
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

		if _, ok := finalResult.Result[pilotCode]; !ok {
			finalResult.Result[pilotCode] = domain.RaceResult{
				FinalPosition: 0,
				PilotCode:     pilotCode,
				PilotName:     pilotName,
				CompletedLaps: totalCompletedLaps,
				TotalRaceTime: totalRaceTime,
			}
		}
	}

	return finalResult
}
