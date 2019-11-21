package calendar

import (
	"time"
)

var location = time.UTC

// TODO: For now we not taking into timezones for simplicity
/*
func init() {
	l, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		l = time.UTC
	}
	location = l
}*/

// Time of event - wrapper on time.Time but more simple constructor
type EventTime struct {
	t time.Time
}

// Constructor
func NewEventTime(year, month, day, hour, minute int) EventTime {
	return EventTime{
		time.Date(year, time.Month(month), day, hour, minute, 0, 0, location),
	}
}

// Construct from time.Time
func ConvertFromTime(t time.Time) EventTime {
	return EventTime{
		time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), 0, 0, location),
	}
}

// Less method for compare 2 event times, will need for sorting
func (eventTime EventTime) Less(thatEventTime EventTime) bool {
	return eventTime.t.Unix() < thatEventTime.t.Unix()
}

// Less or equal for compare 2 event times, will need for query list of events from calendar by period
func (eventTime EventTime) LessOrEqual(thatEventTime EventTime) bool {
	return eventTime.t.Unix() <= thatEventTime.t.Unix()
}

// String representation of event time
func (eventTime EventTime) String() string {
	return eventTime.t.Format("02 Jan 2006 15:04")
}

// Formatting
func (eventTime EventTime) Format(layout string) string {
	return eventTime.t.Format(layout)
}

// get time.Time
func (eventTime EventTime) Time() time.Time {
	return eventTime.t
}
