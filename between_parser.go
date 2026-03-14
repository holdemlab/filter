package filter

import (
	"fmt"
	"strings"
	"time"
)

var (
	formatDate     = "2006-01-02"
	formatDateTime = "2006-01-02T15:04:05"
)

// ParseBetweenOperator splits value by "*" and parses both parts using format.
// It returns a *[Between] holding the parsed start and end times, or a
// *[ParseError] if the value is malformed or either date fails to parse.
func ParseBetweenOperator(format, value string) (*Between, error) {
	var start, end time.Time
	split := strings.Split(value, "*")
	if len(split) != 2 {
		return nil, &ParseError{
			Value:   value,
			Message: fmt.Sprintf("invalid between value %q: expected format 'start*end'", value),
		}
	}
	start, err := time.Parse(format, split[0])
	if err != nil {
		return nil, &ParseError{Value: split[0], Message: "error parsing start date", Err: err}
	}
	end, err = time.Parse(format, split[1])
	if err != nil {
		return nil, &ParseError{Value: split[1], Message: "error parsing end date", Err: err}
	}
	return &Between{start: start, end: end}, nil
}
