package parser

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/dibrito/f1-engine/src/domain"
)

// ProcessRaceResultInput process the race input result files and parses the data to the race result struct
func ProcessRaceResultInput() {
	// 1. read entry file
	readInputFile()
}

func readInputFile() {
	// TODO make a way to have a global logger and not defining it at every file
	//logger, _ := zap.NewProduction()
	//logger.Info(string(dat))

	// TODO make file path configurable
	file, err := os.Open("/Users/pedro.brito/go/src/github.com/dibrito/f1-engine/src/run-result.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	BuildLapMapper(file)
}

func BuildLapMapper(file *os.File) {
	var m domain.LapsMapper
	m.Result = map[int][]domain.LapResult{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		//line := scanner.Text()
		//lap := BuildLap(line)

	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

// BuildMapper build a mapper of pilot code with its laps
func BuildMapper(laps []domain.LapResult) domain.LapsMapper {
	var m domain.LapsMapper
	m.Result = map[int][]domain.LapResult{}
	for _, lap := range laps {
		if _, ok := m.Result[lap.PilotCode]; !ok {
			m.Result[lap.PilotCode] = []domain.LapResult{}
		}
		m.Result[lap.PilotCode] = append(m.Result[lap.PilotCode], lap)
	}
	return m
}

func BuildLap(line string) domain.LapResult {
	//var m domain.LapsMapper
	//m.Result = map[int][]domain.LapResult{}
	line = RemoveTabsFromLine(line)
	fields := RemoveEmptySpacesAndSplit(line)
	fields = RemoveDashFieldFromSplitedFields(fields)
	lap := BuildLapResult(fields)
	//m.Result[lap.PilotCode] = append(m.Result[lap.PilotCode], lap)
	return lap
}

// RemoveDashFieldFromSplitedFields removes – field it got when parsing the line
func RemoveDashFieldFromSplitedFields(splitedFields []string) (result []string) {
	for _, v := range splitedFields {
		if v != "–" {
			result = append(result, v)
		}
	}

	return result
}

// BuildLapResult builds LapResult struct
func BuildLapResult(fields []string) domain.LapResult {
	var lapResult domain.LapResult

	if len(fields) != 6 {
		panic("wrong total field per line")
	}

	// 1. time
	now := time.Now()
	loc, _ := time.LoadLocation("Local")
	hour, minute, second, nsecond := BreakTime(fields[0])
	lapResult.Time = time.Date(now.Year(), now.Month(), now.Day(), hour, minute, second, nsecond, loc)

	// 2. pilot code
	lapResult.PilotCode, _ = strconv.Atoi(fields[1])

	// 3. pilot name
	lapResult.PilotName = fields[2]

	// 4. lap number
	lapResult.Number, _ = strconv.Atoi(fields[3])

	// 5. lap time
	lapResult.Duration, _ = time.ParseDuration(FormatLapDuration(fields[4]))

	// 6. avg lap speed
	lapResult.AvgSpeed = ParseAvgSpeed(fields[5])

	return lapResult
}

// ParseAvgSpeed replaces , speed to . speed
func ParseAvgSpeed(s string) float64 {
	if strings.Contains(s, ",") {
		s = strings.Replace(s, ",", ".", 1)
	}

	sf, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0.0
	}

	return sf
}

// RemoveTabsFromLine should remove all tabs from a given string
func RemoveTabsFromLine(line string) (result string) {
	var sb strings.Builder
	for _, char := range line {
		str := fmt.Sprintf("%c", char)
		if str != "\t" {
			sb.WriteString(str)
		} else {
			sb.WriteString(" ")
		}
	}

	return sb.String()
}

// RemoveEmptySpacesAndSplit should remove all empty spaces
// and return a slice of splited string by empty spaces
func RemoveEmptySpacesAndSplit(line string) (result []string) {
	fields := strings.Split(line, " ")
	for _, f := range fields {
		if f != "" {
			result = append(result, f)
		}
	}
	return result
}

// FormatLapDuration should add m instead of : on duration string
func FormatLapDuration(lapDuration string) string {
	if strings.Contains(lapDuration, ":") {
		lapDuration = strings.Replace(lapDuration, ":", "m", 1)
		lapDuration = fmt.Sprintf("%ss", lapDuration)
		return lapDuration
	}
	return lapDuration
}

// BreakLapTime breaks a give string lap time and returns its minute second and n second
func BreakLapTime(lapTime string) (m, s, ns time.Duration) {
	got, err := time.ParseDuration(lapTime)
	if err != nil {
		panic("error parsing lap time duration ")
	}
	m = got / time.Minute
	rest := got % time.Minute
	s = rest / time.Second
	restN := rest % time.Second
	ns = restN / time.Millisecond

	return m, s, ns
}

// BreakTime given a string time return it's hour minute second and mili second
func BreakTime(time string) (h, m, s, ns int) {
	hasSeparator := strings.Contains(time, ":")
	lenSplitSeparator := len(strings.Split(time, ":"))

	if hasSeparator {
		splitedTime := strings.Split(time, ":")
		if lenSplitSeparator == 2 {
			m, _ = strconv.Atoi(splitedTime[0])
			seconds := strings.Split(splitedTime[1], ".")
			s, _ = strconv.Atoi(seconds[0])
			ns, _ = strconv.Atoi(seconds[1])
			return 0, m, s, ns
		} else if lenSplitSeparator == 3 {
			h, _ = strconv.Atoi(splitedTime[0])
			m, _ = strconv.Atoi(splitedTime[1])
			seconds := strings.Split(splitedTime[2], ".")
			s, _ = strconv.Atoi(seconds[0])
			ns, _ = strconv.Atoi(seconds[1])
		} else {
			panic("wrong time format")
		}
	}

	return h, m, s, ns
}
