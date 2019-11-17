package http

import (
	"fmt"
	"github.com/mitrickx/otus-golang-2019/21/calendar/internal/calendar"
)

type Calendar struct {
	storage *calendar.Calendar
}

func NewCalendar() *Calendar {
	return &Calendar{
		storage: calendar.NewCalendar(),
	}
}

func (thisCalendar *Calendar) AddEvent(event *Event) (int, error) {
	calendarEvent, err := convertToCalendarEvent(event)
	if err != nil {
		return 0, err
	}

	return thisCalendar.storage.AddEvent(*calendarEvent), nil
}

func (thisCalendar *Calendar) UpdateEvent(id int, event *Event) error {
	calendarEvent, err := convertToCalendarEvent(event)
	if err != nil {
		return err
	}

	err = thisCalendar.storage.UpdateEvent(id, *calendarEvent)
	if err != nil {
		return fmt.Errorf("couldn't update event in storage: %w", err)
	}
	return nil
}

func (thisCalendar *Calendar) DeleteEvent(id int) error {
	err := thisCalendar.storage.DeleteEvent(id)
	if err != nil {
		return fmt.Errorf("couldn't delete event from storage: %w", err)
	}
	return nil
}

func (thisCalendar *Calendar) GetEvent(id int) (*Event, bool) {
	if id <= 0 {
		return nil, false
	}

	calendarEvent, ok := thisCalendar.storage.GetEvent(id)
	if !ok {
		return nil, false
	}

	event := ConvertFromCalendarEvent(calendarEvent)
	return event, true

}

func (thisCalendar *Calendar) GetAllEvents() []*Event {
	calendarEvents := thisCalendar.storage.GetAllEvents()
	if len(calendarEvents) == 0 {
		return nil
	}
	var events []*Event
	for _, calendarEvent := range calendarEvents {
		events = append(events, ConvertFromCalendarEvent(calendarEvent))
	}
	return events
}

// Get all events that started in period (boundary of period are included) sorted by Less method of events
// start/end are datetime values represented by string in format on this module (see http.dateTimeLayout)
// Empty string has special meaning - no boundary for range period
func (thisCalendar *Calendar) GetEventsByPeriod(start string, end string) ([]*Event, error) {
	var startTime, endTime *calendar.EventTime
	var err error

	if start != "" {
		startTime, err = ConvertToCalendarEventTime(start)
		if err != nil {
			return nil, err
		}
	}

	if end != "" {
		endTime, err = ConvertToCalendarEventTime(end)
		if err != nil {
			return nil, err
		}
	}

	calendarEvents := thisCalendar.storage.GetEventsByPeriod(startTime, endTime)
	if len(calendarEvents) == 0 {
		return nil, nil
	}
	var events []*Event
	for _, calendarEvent := range calendarEvents {
		events = append(events, ConvertFromCalendarEvent(calendarEvent))
	}
	return events, nil
}

func (thisCalendar *Calendar) getEventsTotalCount() int {
	return thisCalendar.storage.Count()
}

func convertToCalendarEvent(event *Event) (*calendar.Event, error) {
	calendarEvent, err := event.ConvertToCalendarEvent()
	if err != nil {
		return nil, fmt.Errorf("couldn't convert http.Event to calender.Event: %w", err)
	}
	return calendarEvent, nil
}
