package http

import (
	"encoding/json"
	"fmt"
	"github.com/mitrickx/otus-golang-2019/21/calendar/internal/calendar"
	"time"
)

const layout = "2006-01-02 15:04"

type Event struct {
	Id    int       `json:",omitempty"`
	Name  string    `json:"name"`
	Start string	`json:"start"`	// Y-m-d H:i
	End   string	`json:"end"`	// Y-m-d H:i
}

func ConvertFromCalendarEvent(calendarEvent calendar.Event) *Event {
	event := &Event{
		Id:    calendarEvent.Id(),
		Name:  calendarEvent.Name(),
		Start: calendarEvent.Start().Format(layout),
		End:   calendarEvent.End().Format(layout),
	}
	return event
}

func (event *Event) ConvertToCalendarEvent() (*calendar.Event, error) {
	var startTime, endTime time.Time
	var err error

	startTime, err = time.Parse(layout, event.Start)
	if err != nil {
		return nil, fmt.Errorf("couldn't parse start datetime: %w", err)
	}

	endTime, err = time.Parse(layout, event.End)
	if err != nil {
		return nil, fmt.Errorf("couldn't parse end datetime: %w", err)
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
