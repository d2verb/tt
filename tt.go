package tt

import (
	"flag"
	"fmt"
	"io"
	"math"
	"strings"
	"time"
)

// Run the tt
func Run(argv []string, outStream, errStream io.Writer) error {
	fs := flag.NewFlagSet(
		fmt.Sprintf("tt (v%s rev:%s", version, revision), flag.ContinueOnError)
	fs.SetOutput(errStream)

	v := fs.Bool("version", false, "display version")
	if err := fs.Parse(argv); err != nil {
		return err
	}

	if *v {
		return printVersion(outStream)
	}

	if len(fs.Args()) < 2 {
		return fmt.Errorf("startTime and endTime must be specified")
	}

	start := fs.Args()[0]
	end := fs.Args()[1]

	return run(start, end, outStream, errStream)
}

func printVersion(out io.Writer) error {
	_, err := fmt.Fprintf(out, "tt v%s (rev:%s)\n", version, revision)
	return err
}

func run(start, end string, outStream, errStream io.Writer) error {
	startTime, err := parseTime(start)
	if err != nil {
		return err
	}

	endTime, err := parseTime(end)
	if err != nil {
		return err
	}

	// if startTime >= endTime then error
	if !startTime.Before(*endTime) {
		return fmt.Errorf("\"%v\" must be older than \"%v\"", start, end)
	}

	regexp, err := genRegExp(startTime, endTime)
	if err != nil {
		return err
	}

	fmt.Println(*regexp)
	return nil
}

func parseTime(timeStr string) (*time.Time, error) {
	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		return nil, err
	}

	layout := "2006-01-02 15:04:05"

	t, err := time.ParseInLocation(layout, timeStr, jst)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func genRegExp(startTime *time.Time, endTime *time.Time) (*string, error) {
	yearPart, err := genRegExpYear(startTime.Year(), endTime.Year())
	if err != nil {
		return nil, err
	}

	var match bool

	match = startTime.Year() == endTime.Year()
	monthPart, err := genRegExpMonth(int(startTime.Month()), int(endTime.Month()), match)
	if err != nil {
		return nil, err
	}

	match = match && (startTime.Month() == endTime.Month())
	dayPart, err := genRegExpDay(startTime.Day(), endTime.Day(), match)
	if err != nil {
		return nil, err
	}

	match = match && (startTime.Day() == endTime.Day())
	hourPart, err := genRegExpHour(startTime.Hour(), endTime.Hour(), match)
	if err != nil {
		return nil, err
	}

	match = match && (startTime.Hour() == endTime.Hour())
	minutePart, err := genRegExpMinute(startTime.Minute(), endTime.Minute(), match)
	if err != nil {
		return nil, err
	}

	match = match && (startTime.Minute() == endTime.Minute())
	secondPart, err := genRegExpSecond(startTime.Second(), endTime.Second(), match)
	if err != nil {
		return nil, err
	}

	regexp := fmt.Sprintf("%v-%v-%v %v:%v:%v", *yearPart, *monthPart, *dayPart, *hourPart, *minutePart, *secondPart)
	regexp = simplify(regexp)

	return &regexp, nil
}

func simplify(regexp string) string {
	return strings.Replace(regexp, "[0-9]", "\\d", -1)
}

func genRegExpYear(startYear int, endYear int) (*string, error) {
	return genRegExpRange(genDigitsRange(startYear, endYear, 4))
}

func genRegExpMonth(startMonth int, endMonth int, matchTilYear bool) (*string, error) {
	var digitsRange []*valueRange

	if matchTilYear {
		if startMonth > endMonth {
			return nil, fmt.Errorf("start month (%v) is greater than end month (%v)", startMonth, endMonth)
		}
		digitsRange = genDigitsRange(startMonth, endMonth, 2)
	} else {
		digitsRange = genDigitsRange(startMonth, 12, 2)
		for d, vr := range genDigitsRange(1, endMonth, 2) {
			digitsRange[d].max = max(digitsRange[d].max, vr.max)
			digitsRange[d].min = min(digitsRange[d].min, vr.min)
		}
	}

	return genRegExpRange(digitsRange)
}

func genRegExpDay(startDay int, endDay int, matchTilMonth bool) (*string, error) {
	var digitsRange []*valueRange

	if matchTilMonth {
		if startDay > endDay {
			return nil, fmt.Errorf("start day (%v) is greater than end day (%v)", startDay, endDay)
		}
		digitsRange = genDigitsRange(startDay, endDay, 2)
	} else {
		digitsRange = genDigitsRange(startDay, 31, 2)
		for d, vr := range genDigitsRange(1, endDay, 2) {
			digitsRange[d].max = max(digitsRange[d].max, vr.max)
			digitsRange[d].min = min(digitsRange[d].min, vr.min)
		}
	}

	return genRegExpRange(digitsRange)
}

func genRegExpHour(startHour int, endHour int, matchTilDay bool) (*string, error) {
	var digitsRange []*valueRange

	if matchTilDay {
		if startHour > endHour {
			return nil, fmt.Errorf("start hour (%v) is greater than end hour (%v)", startHour, endHour)
		}
		digitsRange = genDigitsRange(startHour, endHour, 2)
	} else {
		digitsRange = genDigitsRange(startHour, 23, 2)
		for d, vr := range genDigitsRange(0, endHour, 2) {
			digitsRange[d].max = max(digitsRange[d].max, vr.max)
			digitsRange[d].min = min(digitsRange[d].min, vr.min)
		}
	}

	return genRegExpRange(digitsRange)
}

func genRegExpMinute(startMinute int, endMinute int, matchTilHour bool) (*string, error) {
	var digitsRange []*valueRange

	if matchTilHour {
		if startMinute > endMinute {
			return nil, fmt.Errorf("start minute (%v) is greater than end minute (%v)", startMinute, endMinute)
		}
		digitsRange = genDigitsRange(startMinute, endMinute, 2)
	} else {
		digitsRange = genDigitsRange(startMinute, 59, 2)
		for d, vr := range genDigitsRange(0, endMinute, 2) {
			digitsRange[d].max = max(digitsRange[d].max, vr.max)
			digitsRange[d].min = min(digitsRange[d].min, vr.min)
		}
	}

	return genRegExpRange(digitsRange)
}

func genRegExpSecond(startSecond int, endSecond int, matchTilMinute bool) (*string, error) {
	var digitsRange []*valueRange

	if matchTilMinute {
		if startSecond > endSecond {
			return nil, fmt.Errorf("start second (%v) is greater than end second (%v)", startSecond, endSecond)
		}
		digitsRange = genDigitsRange(startSecond, endSecond, 2)
	} else {
		digitsRange = genDigitsRange(startSecond, 59, 2)
		for d, vr := range genDigitsRange(0, endSecond, 2) {
			digitsRange[d].max = max(digitsRange[d].max, vr.max)
			digitsRange[d].min = min(digitsRange[d].min, vr.min)
		}
	}

	return genRegExpRange(digitsRange)
}

type valueRange struct {
	max int
	min int
}

func newValueRange() *valueRange {
	return &valueRange{
		max: math.MinInt32,
		min: math.MaxInt32,
	}
}

func genRegExpRange(digitsRange []*valueRange) (*string, error) {
	regexp := ""
	for _, vr := range digitsRange {
		if vr.max == vr.min {
			regexp = fmt.Sprintf("%v", vr.max) + regexp
		} else {
			regexp = fmt.Sprintf("[%v-%v]", vr.min, vr.max) + regexp
		}
	}
	return &regexp, nil
}

func genDigitsRange(start, end, nd int) []*valueRange {
	digits := make([]*valueRange, nd)
	for i := 0; i < nd; i++ {
		digits[i] = newValueRange()
	}
	for i := start; i <= end; i++ {
		value := i
		for j := 0; j < nd; j++ {
			digit := value % 10
			digits[j].max = max(digits[j].max, digit)
			digits[j].min = min(digits[j].min, digit)
			value /= 10
		}
	}
	return digits
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a < b {
		return b
	}
	return a
}
