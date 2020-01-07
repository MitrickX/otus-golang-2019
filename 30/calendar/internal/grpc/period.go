package grpc

import (
	"github.com/golang/protobuf/ptypes/timestamp"
	"time"
)

// Period struct for get event lists by periods
// start and end could be nil - means no boundary (-ies) of period
type Period struct {
	start *timestamp.Timestamp
	end   *timestamp.Timestamp
}

// Constructor
func NewPeriod(start *timestamp.Timestamp, end *timestamp.Timestamp) *Period {
	return &Period{
		start: start,
		end:   end,
	}
}

//  Construct period of current day
func NewDayPeriod(now time.Time) (*Period, error) {
	start, err := NewTimestamp(now.Year(), int(now.Month()), now.Day(), 0, 0)
	if err != nil {
		return nil, err
	}
	end, err := NewTimestamp(now.Year(), int(now.Month()), now.Day(), 23, 59)
	if err != nil {
		return nil, err
	}
	return NewPeriod(start, end), nil
}

// Construct period of current week
func NewWeekPeriod(now time.Time) (*Period, error) {
	nowWeek := now.Weekday()

	// shift to monday
	shiftDays := 0
	if nowWeek == time.Sunday {
		shiftDays = 6
	} else {
		shiftDays = int(nowWeek) - int(time.Monday)
	}

	monday := now.AddDate(0, 0, -shiftDays)
	sunday := monday.AddDate(0, 0, 6)

	start, err := NewTimestamp(monday.Year(), int(monday.Month()), monday.Day(), 0, 0)
	if err != nil {
		return nil, err
	}
	end, err := NewTimestamp(sunday.Year(), int(sunday.Month()), sunday.Day(), 23, 59)
	if err != nil {
		return nil, err
	}

	return NewPeriod(start, end), nil
}

// // Construct period of current month
func NewMonthPeriod(now time.Time) (*Period, error) {
	firstDayInMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	nextMonth := firstDayInMonth.AddDate(0, 1, 0)
	lastDayInMonth := nextMonth.AddDate(0, 0, -1)

	startTime, err := NewTimestamp(firstDayInMonth.Year(), int(firstDayInMonth.Month()), firstDayInMonth.Day(), 0, 0)
	if err != nil {
		return nil, err
	}

	endTime, err := NewTimestamp(lastDayInMonth.Year(), int(lastDayInMonth.Month()), lastDayInMonth.Day(), 23, 59)
	if err != nil {
		return nil, err
	}

	return NewPeriod(startTime, endTime), nil
}
