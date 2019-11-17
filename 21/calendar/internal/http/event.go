package http

import (
	"encoding/json"
	"fmt"
	"github.com/mitrickx/otus-golang-2019/21/calendar/internal/calendar"
	"time"
)

const dateTimeLayout = "2006-01-02 15:04"
const dateLayout = "2006-01-02"

type ErrorInvalidDatetime struct {
	err error
}

func (e *ErrorInvalidDatetime) Error() string {
	return e.err.Error()
}

var DefaultErrorInvalidDatetime = &ErrorInvalidDatetime{
	fmt.Errorf("invalid format of datetime - must be Y-m-d H:i (e.g %s)", dateTimeLayout),
}

type Event struct {
	Id    int    `json:",omitempty"`
	Name  string `json:"name"`
	Start string `json:"start"` // Y-m-d H:i
	End   string `json:"end"`   // Y-m-d H:i
}

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

func ConvertFromCalendarEvent(calendarEvent calendar.Event) *Event {
	event := &Event{
		Id:    calendarEvent.Id(),
		Name:  calendarEvent.Name(),
		Start: calendarEvent.Start().Format(dateTimeLayout),
		End:   calendarEvent.End().Format(dateTimeLayout),
	}
	return event
}

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

func JsonUnmarshal(data []byte) (*Event, error) {
	event := &Event{}
	err := event.JsonUnmarshal(data)
	return event, err
}

func (event *Event) JsonUnmarshal(data []byte) error {
	return json.Unmarshal(data, event)
}

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
