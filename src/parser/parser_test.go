package parser

import (
	"bufio"
	"log"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/dibrito/f1-engine/src/domain"
)

// lines
// TODO test all lines
var (
	Line1 = "23:49:08.277      038 – F.MASSA                           1\t\t1:02.852                        44,275"
	Line2 = "23:49:10.858      033 – R.BARRICHELLO                     1\t\t1:04.352                        43,243"
	Line3 = "23:49:11.075      002 – K.RAIKKONEN                       1             1:04.108                        43,408"
	//line4  = "23:49:12.667      023 – M.WEBBER                          1\t\t1:04.414                        43,202"
	//line5  = "23:49:30.976      015 – F.ALONSO                          1\t\t1:18.456\t\t\t35,47"
	Line6 = "23:50:11.447      038 – F.MASSA                           2\t\t1:03.170                        44,053"
	line7  = "23:50:14.860      033 – R.BARRICHELLO                     2\t\t1:04.002                        43,48"
	line8  = "23:50:15.057      002 – K.RAIKKONEN                       2             1:03.982                        43,493"
	//line9  = "23:50:17.472      023 – M.WEBBER                          2\t\t1:04.805                        42,941"
	//line10 = "23:50:37.987      015 – F.ALONSO                          2\t\t1:07.011\t\t\t41,528"
	//line11 = "23:51:14.216      038 – F.MASSA                           3\t\t1:02.769                        44,334"
	//line12 = "23:51:18.576      033 – R.BARRICHELLO\t\t          3\t\t1:03.716                        43,675"
	//line13 = "23:51:19.044      002 – K.RAIKKONEN                       3\t\t1:03.987                        43,49"
	//line14 = "23:51:21.759      023 – M.WEBBER                          3\t\t1:04.287                        43,287"
	//line15 = "23:51:46.691      015 – F.ALONSO                          3\t\t1:08.704\t\t\t40,504"
	//line16 = "23:52:01.796      011 – S.VETTEL                          1\t\t3:31.315\t\t\t13,169"
	//line17 = "23:52:17.003      038 – F.MASS                            4\t\t1:02.787                        44,321"
	//line18 = "23:52:22.586      033 – R.BARRICHELLO\t\t          4\t\t1:04.010                        43,474"
	//line19 = "23:52:22.120      002 – K.RAIKKONEN                       4\t\t1:03.076                        44,118"
	//line20 = "23:52:25.975      023 – M.WEBBER                          4\t\t1:04.216                        43,335"
	//line21 = "23:53:06.741      015 – F.ALONSO                          4\t\t1:20.050\t\t\t34,763"
	//line22 = "23:53:39.660      011 – S.VETTEL                          2\t\t1:37.864\t\t\t28,435"
)

// Laps
var (
	LapDuration1, _ = time.ParseDuration(FormatLapDuration("1:02.852"))
	LapDuration2, _ = time.ParseDuration(FormatLapDuration("1:04.352"))
	LapDuration3, _ = time.ParseDuration(FormatLapDuration("1:04.108"))
	LapDuration6, _ = time.ParseDuration(FormatLapDuration("1:03.170"))
)

// lap results
var (
	Loc, _ = time.LoadLocation("Local")
	Now    = time.Now()

	LapResults = []domain.LapResult{
		{
			Time:      time.Date(Now.Year(), Now.Month(), Now.Day(), 23, 49, 8, 277, Loc),
			PilotCode: 38,
			PilotName: "F.MASSA",
			Number:    1,
			Duration:  LapDuration1,
			AvgSpeed:  44.275,
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
			Time:      time.Date(Now.Year(), Now.Month(), Now.Day(), 23, 49, 11, 75, Loc),
			PilotCode: 2,
			PilotName: "K.RAIKKONEN",
			Number:    1,
			Duration:  LapDuration3,
			AvgSpeed:  43.408,
		},
	}
)

var (
	lines = []string{Line1, Line2, Line3}
)

func TestProcessRaceResultInput(t *testing.T) {
	t.Run("should parse all pilot laps", func(t *testing.T) {
		for i, v := range lines {
			lap := BuildLap(v)
			require.EqualValues(t, LapResults[i], lap)
		}
	})

	t.Run("should remove all tabs from a string", func(t *testing.T) {
		given := "23:49:08.277      038 – F.MASSA                           1\t\t1:02.852                        44,275"
		want := "23:49:08.277      038 – F.MASSA                           1  1:02.852                        44,275"
		var got string
		got = RemoveTabsFromLine(given)
		require.EqualValues(t, want, got)
	})

	t.Run("should create a mapper of laps from pilot laps", func(t *testing.T) {
		LapResults = append(LapResults, domain.LapResult{
			Time:      time.Date(Now.Year(), Now.Month(), Now.Day(), 23, 50, 11, 447, Loc),
			PilotCode: 38,
			PilotName: "F.MASSA",
			Number:    2,
			Duration:  LapDuration6,
			AvgSpeed:  44.053,
		})

		lines = append(lines, Line6)

		laps := []domain.LapResult{}

		for i, v := range lines {
			lap := BuildLap(v)
			require.EqualValues(t, LapResults[i], lap)
			laps = append(laps, lap)
		}

		m := BuildMapper(laps)
		code := LapResults[0].PilotCode
		_, ok := m.Result[code]
		require.True(t, ok)

		code = LapResults[1].PilotCode
		_, ok = m.Result[code]
		require.True(t, ok)

		code = LapResults[2].PilotCode
		_, ok = m.Result[code]
		require.True(t, ok)

		code = LapResults[3].PilotCode
		_, ok = m.Result[code]
		require.True(t, ok)
		require.Len(t, m.Result[code], 2)
	})

	t.Run("should split a no tab string into fields", func(t *testing.T) {
		given := "23:49:08.277      038 – F.MASSA                           1  1:02.852                        44,275"
		want := []string{
			"23:49:08.277",
			"038",
			"–",
			"F.MASSA",
			"1",
			"1:02.852",
			"44,275",
		}
		var got []string
		got = RemoveEmptySpacesAndSplit(given)
		require.Len(t, got, 7)
		require.EqualValues(t, want, got)
	})

	t.Run("should split a no tab string into fields and remove - field", func(t *testing.T) {
		given := "23:49:08.277      038 – F.MASSA                           1  1:02.852                        44,275"
		want := []string{
			"23:49:08.277",
			"038",
			"–",
			"F.MASSA",
			"1",
			"1:02.852",
			"44,275",
		}

		wantNoDashField := []string{
			"23:49:08.277",
			"038",
			"F.MASSA",
			"1",
			"1:02.852",
			"44,275",
		}
		var got []string
		got = RemoveEmptySpacesAndSplit(given)

		gotDashField := RemoveDashFieldFromSplitedFields(got)
		require.Len(t, got, 7)
		require.Len(t, gotDashField, 6)
		require.EqualValues(t, want, got)
		require.EqualValues(t, wantNoDashField, gotDashField)
	})

	t.Run("should parse time", func(t *testing.T) {
		given := "23:49:08.277"
		h := 23
		m := 49
		s := 8
		ns := 277
		hGot, mGot, sGot, nsGot := BreakTime(given)
		require.EqualValues(t, hGot, h)
		require.EqualValues(t, mGot, m)
		require.EqualValues(t, sGot, s)
		require.EqualValues(t, nsGot, ns)

		given = "1:02.852"
		h = 0
		m = 1
		s = 2
		ns = 852
		hGot, mGot, sGot, nsGot = BreakTime(given)
		require.EqualValues(t, hGot, h)
		require.EqualValues(t, mGot, m)
		require.EqualValues(t, sGot, s)
		require.EqualValues(t, nsGot, ns)
	})

	t.Run("should format lap duration", func(t *testing.T) {
		given := "1:02.852"
		want := "1m02.852s"
		got := FormatLapDuration(given)
		require.EqualValues(t, want, got)
	})

	t.Run("should parse lap duration", func(t *testing.T) {
		given := "1m02.852s"
		want := 62852000000
		got, err := time.ParseDuration(given)

		require.Nil(t, err)
		require.EqualValues(t, want, got)
	})

	t.Run("should parse lap duration and convert to min seconds and nano seconds", func(t *testing.T) {
		given := "1m02.852s"
		m, s, ns := BreakLapTime(given)
		require.EqualValues(t, m, 1)
		require.EqualValues(t, s, 2)
		require.EqualValues(t, ns, 852)
	})

	t.Run("should fill LapResult struct", func(t *testing.T) {
		given := "23:49:08.277      038 – F.MASSA                           1  1:02.852                        44,275"
		loc, _ := time.LoadLocation("Local")
		now := time.Now()
		lapDuration, _ := time.ParseDuration("1m02.852s")
		want := domain.LapResult{
			Time:      time.Date(now.Year(), now.Month(), now.Day(), 23, 49, 8, 277, loc),
			PilotCode: 38,
			PilotName: "F.MASSA",
			Number:    1,
			Duration:  lapDuration,
			AvgSpeed:  44.275,
		}

		lap := BuildLap(given)
		require.EqualValues(t, want, lap)
	})

	t.Run("should fill LapResult struct", func(t *testing.T) {
		given := "44,275"
		want := 44.275

		got := ParseAvgSpeed(given)
		require.EqualValues(t, want, got)
	})

	t.Run("should parse all files lines", func(t *testing.T) {
		file, err := os.Open("/Users/pedro.brito/go/src/github.com/dibrito/f1-engine/src/run-result.txt")
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			noTabsLine := RemoveTabsFromLine(scanner.Text())
			sliceOfFields := RemoveEmptySpacesAndSplit(noTabsLine)
			require.Len(t, sliceOfFields, 7)
		}
	})
}
