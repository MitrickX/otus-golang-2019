package entities

import (
	"time"
)

const layout = "02 Jan 2006 15:04"

// Time of event - wrapper on time.Time but more simple constructor
type DateTime struct {
	t time.Time
}

// Constructor
func NewDateTime(year, month, day, hour, minute int) DateTime {
	return DateTime{
		time.Date(year, time.Month(month), day, hour, minute, 0, 0, time.UTC),
	}
}

// Construct from time.Time
func ConvertFromTime(t time.Time) DateTime {
	return DateTime{
		time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), 0, 0, time.UTC),
	}
}

// Less method for compare 2 event times, will need for sorting
func (eventTime DateTime) Less(thatEventTime DateTime) bool {
	return eventTime.t.Unix() < thatEventTime.t.Unix()
}

// Less or equal for compare 2 event times, will need for query list of events from entities by period
func (eventTime DateTime) LessOrEqual(thatEventTime DateTime) bool {
	return eventTime.t.Unix() <= thatEventTime.t.Unix()
}

// String representation of event time
func (eventTime DateTime) String() string {
	return eventTime.t.Format(layout)
}

// Formatting
func (eventTime DateTime) Format(layout string) string {
	return eventTime.t.Format(layout)
}

// get time.Time
func (eventTime DateTime) Time() time.Time {
	return eventTime.t
}

// Minus minutes
func (eventTime DateTime) MinusMinutes(m int) DateTime {
	return ConvertFromTime(eventTime.t.Add(-time.Duration(m) * time.Minute))
}

// Plus minutes
func (eventTime DateTime) PlusMinutes(m int) DateTime {
	return ConvertFromTime(eventTime.t.Add(time.Duration(m) * time.Minute))
}
