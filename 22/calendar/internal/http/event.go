package http

import (
	"encoding/json"
	"fmt"
	"github.com/mitrickx/otus-golang-2019/22/calendar/internal/calendar"
	"time"
)

// Inside http package and for communication with outer world by http we deal with "Y-m-d H:i" and "Y-m-d" date/datetime strings
const dateTimeLayout = "2006-01-02 15:04"
const dateLayout = "2006-01-02"

// Typed error about invalid datetime
type ErrorInvalidDatetime struct {
	err error
}

// Error interface
func (e *ErrorInvalidDatetime) Error() string {
	return e.err.Error()
}

// Default invalid datetime error
var DefaultErrorInvalidDatetime = &ErrorInvalidDatetime{
	fmt.Errorf("invalid format of datetime - must be Y-m-d H:i (e.g %s)", dateTimeLayout),
}

// Event structure for work inside http package
// Clean architecture approach - not working with inner biz logic layer directly
type Event struct {
	Id    int    `json:",omitempty"`
	Name  string `json:"name"`
	Start string `json:"start"` // Y-m-d H:i
	End   string `json:"end"`   // Y-m-d H:i
}

// Constructor
func NewEvent(name, start, end string) (*Event, error) {
	_, err := time.Parse(dateTimeLayout, start)
	if err != nil {
		return nil, DefaultErrorInvalidDatetime
	}

	_, err = time.Parse(dateTimeLayout, end)
	if err != nil {
		return nil, DefaultErrorInvalidDatetime
	}

	event := &Event{
		Name:  name,
		Start: start,
		End:   end,
	}

	return event, nil
}

// Convert from inner Event entity (calendar.Event) to http.Event
func ConvertFromCalendarEvent(calendarEvent calendar.Event) *Event {
	event := &Event{
		Id:    calendarEvent.Id(),
		Name:  calendarEvent.Name(),
		Start: calendarEvent.Start().Format(dateTimeLayout),
		End:   calendarEvent.End().Format(dateTimeLayout),
	}
	return event
}

// Convert from http.Entity to inner Event entity (calendar.Event)
func (event *Event) ConvertToCalendarEvent() (*calendar.Event, error) {
	var startTime, endTime time.Time
	var err error

	startTime, err = time.Parse(dateTimeLayout, event.Start)
	if err != nil {
		return nil, &ErrorInvalidDatetime{
			fmt.Errorf("couldn't parse start datetime: %w", err),
		}
	}

	endTime, err = time.Parse(dateTimeLayout, event.End)
	if err != nil {
		return nil, &ErrorInvalidDatetime{
			fmt.Errorf("couldn't parse end datetime: %w", err),
		}
	}

	calendarEvent := calendar.NewEvent(event.Name, calendar.ConvertFromTime(startTime), calendar.ConvertFromTime(endTime))
	return &calendarEvent, nil
}

// Json unmarshal function for event
func JsonUnmarshal(data []byte) (*Event, error) {
	event := &Event{}
	err := event.JsonUnmarshal(data)
	return event, err
}

// Json unmarshal method for event
func (event *Event) JsonUnmarshal(data []byte) error {
	return json.Unmarshal(data, event)
}

// Json marshal method for event
func (event *Event) JsonMarshall() ([]byte, error) {
	return json.Marshal(event)
}

// Datetime in format on current package (dateTimeLayout)
func ConvertToCalendarEventTime(datetime string) (*calendar.EventTime, error) {
	t, err := time.Parse(dateTimeLayout, datetime)
	if err != nil {
		return nil, DefaultErrorInvalidDatetime
	}
	eventTime := calendar.ConvertFromTime(t)
	return &eventTime, nil
}

//  Helper that calculated period for day
func GetDayPeriod(now time.Time) (string, string) {
	startTime := now.Format(dateLayout) + " 00:00"
	endTime := now.Format(dateLayout) + " 23:59"
	return startTime, endTime
}

// Helper that calculated period for week
func GetWeekPeriod(now time.Time) (string, string) {
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

	startTime := monday.Format(dateLayout) + " 00:00"
	endTime := sunday.Format(dateLayout) + " 23:59"

	return startTime, endTime
}

// Helper that calculated period for month
func GetMonthPeriod(now time.Time) (string, string) {
	firstDayInMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	nextMonth := firstDayInMonth.AddDate(0, 1, 0)
	lastDayInMonth := nextMonth.AddDate(0, 0, -1)

	startTime := firstDayInMonth.Format(dateTimeLayout)
	endTime := lastDayInMonth.Format(dateLayout) + " 23:59"

	return startTime, endTime
}
