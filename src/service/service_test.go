package service

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/dibrito/f1-engine/src/domain"
	"github.com/dibrito/f1-engine/src/parser"
)

// TODO fix issue with adding laps multiple times
// could not reference var even thought they are exported
// Laps
var (
	LapDuration1, _ = time.ParseDuration(parser.FormatLapDuration("1:02.852"))
	LapDuration2, _ = time.ParseDuration(parser.FormatLapDuration("1:04.352"))
	LapDuration3, _ = time.ParseDuration(parser.FormatLapDuration("1:04.108"))
	LapDuration6, _ = time.ParseDuration(parser.FormatLapDuration("1:03.170"))
	LapDuration7, _ = time.ParseDuration(parser.FormatLapDuration("1:04.002"))
	LapDuration8, _ = time.ParseDuration(parser.FormatLapDuration("1:03.982"))
)

var (
	Loc, _ = time.LoadLocation("Local")
	Now    = time.Now()

	Laps = []domain.LapResult{
		{
			Time:      time.Date(Now.Year(), Now.Month(), Now.Day(), 23, 49, 8, 277, Loc),
			PilotCode: 38,
			PilotName: "F.MASSA",
			Number:    1,
			Duration:  LapDuration1,
			AvgSpeed:  44.275,
		},
		{
			Time:      time.Date(Now.Year(), Now.Month(), Now.Day(), 23, 50, 11, 447, Loc),
			PilotCode: 38,
			PilotName: "F.MASSA",
			Number:    2,
			Duration:  LapDuration6,
			AvgSpeed:  44.053,
		},
		{
			Time:      time.Date(Now.Year(), Now.Month(), Now.Day(), 23, 49, 10, 858, Loc),
			PilotCode: 33,
			PilotName: "R.BARRICHELLO",
			Number:    1,
			Duration:  LapDuration2,
			AvgSpeed:  43.243,
		},
		{
			Time:      time.Date(Now.Year(), Now.Month(), Now.Day(), 23, 50, 14, 860, Loc),
			PilotCode: 33,
			PilotName: "R.BARRICHELLO",
			Number:    2,
			Duration:  LapDuration7,
			AvgSpeed:  43.48,
		},
		{
			Time:      time.Date(Now.Year(), Now.Month(), Now.Day(), 23, 49, 11, 75, Loc),
			PilotCode: 2,
			PilotName: "K.RAIKKONEN",
			Number:    1,
			Duration:  LapDuration3,
			AvgSpeed:  43.243,
		},
		{
			Time:      time.Date(Now.Year(), Now.Month(), Now.Day(), 23, 50, 15, 57, Loc),
			PilotCode: 2,
			PilotName: "K.RAIKKONEN",
			Number:    2,
			Duration:  LapDuration8,
			AvgSpeed:  43.48,
		},
	}

	m = parser.BuildMapper(Laps)

	wantFinalResult = domain.RaceResultMapper{
		Result: map[int]domain.RaceResult{},
	}
)

func setFinalResult() {
	wantFinalResult.Result[38] = domain.RaceResult{
		FinalPosition: 1,
		PilotCode:     38,
		PilotName:     "F.MASSA",
		CompletedLaps: 2,
		TotalRaceTime: GetRaceDuration([]domain.LapResult{
			{
				Duration: LapDuration1,
			},
			{
				Duration: LapDuration6,
			},
		}),
	}
	wantFinalResult.Result[33] = domain.RaceResult{
		FinalPosition: 2,
		PilotCode:     33,
		PilotName:     "R.BARRICHELLO",
		CompletedLaps: 2,
		TotalRaceTime: GetRaceDuration([]domain.LapResult{
			{
				Duration: LapDuration2,
			},
			{
				Duration: LapDuration7,
			},
		}),
	}
	wantFinalResult.Result[2] = domain.RaceResult{
		FinalPosition: 3,
		PilotCode:     2,
		PilotName:     "K.RAIKKONEN",
		CompletedLaps: 2,
		TotalRaceTime: GetRaceDuration([]domain.LapResult{
			{
				Duration: LapDuration3,
			},
			{
				Duration: LapDuration8,
			},
		}),
	}
}

func TestGetRaceTotalTime(t *testing.T) {
	t.Run("should get race total time give pilot's Laps", func(t *testing.T) {
		want, err := time.ParseDuration(parser.FormatLapDuration("2:06.022"))
		require.NoError(t, err)

		totalRaceTime := GetRaceDuration(Laps)
		require.EqualValues(t, totalRaceTime, want)
	})
}

func TestGetTotalCompletedLaps(t *testing.T) {
	t.Run("should get race total Laps", func(t *testing.T) {
		want := 2
		got := GetTotalCompletedLaps(Laps)
		require.EqualValues(t, want, got)
	})
}

func TestOrderRaceResult(t *testing.T) {
	t.Run("should ordered result", func(t *testing.T) {
		setFinalResult()
		ordered:=OrderRaceResult(wantFinalResult)
		// TODO way of comparing this without put pilot code hard coded
		require.EqualValues(t,ordered.Result[38].FinalPosition,1)
		require.EqualValues(t,ordered.Result[2].FinalPosition,2)
		require.EqualValues(t,ordered.Result[33].FinalPosition,3)
	})
}

func TestProcessFinalResult(t *testing.T) {
	t.Run("should final position for a pilot", func(t *testing.T) {
		setFinalResult()
		got := ProcessFinalResult(m)
		require.EqualValues(t, wantFinalResult, got)
	})
}
