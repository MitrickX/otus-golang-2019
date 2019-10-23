package calendar

import "time"

// Time of event - wrapper on time.Time but more simple constructor
type EventTime struct {
	t time.Time
}

// Constructor
func NewEventTime(year, month, day, hour, minute int) EventTime {
	location, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		location = time.UTC
	}
	return EventTime{
		time.Date(year, time.Month(month), day, hour, minute, 0, 0, location),
	}
}

// Less method for compare 2 event times, will need for sorting
func (eventTime EventTime) Less(thatEventTime EventTime) bool {
	return eventTime.t.Unix() < thatEventTime.t.Unix()
}

// String representation of event time
func (eventTime EventTime) String() string {
	return eventTime.t.Format("02 Jan 2006 15:04")
}
