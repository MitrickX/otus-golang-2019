package http

import (
	"errors"
	"fmt"
	"github.com/mitrickx/otus-golang-2019/28/calendar/internal/domain/entities"
)

// Calendar structure for work inside http package
// Clean architecture approach - not working with inner biz logic layer directly
type Calendar struct {
	storage entities.Storage // for now it is inner biz entity itself, for future there will be storage interface
}

// Constructor
func NewCalendar(storage entities.Storage) (*Calendar, error) {
	if storage == nil {
		return nil, errors.New("storage must be defined (not nil)")
	}
	return &Calendar{
		storage: storage,
	}, nil
}

// Add Event
func (thisCalendar *Calendar) AddEvent(event *Event) (int, error) {
	calendarEvent, err := convertToCalendarEvent(event)
	if err != nil {
		return 0, err
	}

	return thisCalendar.storage.AddEvent(*calendarEvent)
}

// Update Event
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

// Delete Event
func (thisCalendar *Calendar) DeleteEvent(id int) error {
	err := thisCalendar.storage.DeleteEvent(id)
	if err != nil {
		return fmt.Errorf("couldn't delete event from storage: %w", err)
	}
	return nil
}

// Get one event
func (thisCalendar *Calendar) GetEvent(id int) (*Event, bool) {
	if id <= 0 {
		return nil, false
	}

	calendarEvent, err := thisCalendar.storage.GetEvent(id)
	if err == entities.StorageErrorEventNotFound {
		return nil, false
	}

	event := ConvertFromCalendarEvent(calendarEvent)
	return event, true

}

// Get all events
func (thisCalendar *Calendar) GetAllEvents() ([]*Event, error) {
	calendarEvents, err := thisCalendar.storage.GetAllEvents()
	if err != nil {
		return nil, nil
	}
	if len(calendarEvents) == 0 {
		return nil, nil
	}
	var events []*Event
	for _, calendarEvent := range calendarEvents {
		events = append(events, ConvertFromCalendarEvent(calendarEvent))
	}
	return events, nil
}

// Get all events that started in period (boundary of period are included) sorted by Less method of events
// start/end are datetime values represented by string in format on this module (see http.dateTimeLayout)
// Empty string has special meaning - no boundary for range period
func (thisCalendar *Calendar) GetEventsByPeriod(start string, end string) ([]*Event, error) {
	var startTime, endTime *entities.DateTime
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

	calendarEvents, err := thisCalendar.storage.GetEventsByPeriod(startTime, endTime)
	if len(calendarEvents) == 0 {
		return nil, err
	}
	var events []*Event
	for _, calendarEvent := range calendarEvents {
		events = append(events, ConvertFromCalendarEvent(calendarEvent))
	}
	return events, nil
}

// Get total number of events in entities
func (thisCalendar *Calendar) getEventsTotalCount() int {
	cnt, _ := thisCalendar.storage.Count()
	return cnt
}

// Inner Helper that helps convert http.Event to entities.Event
func convertToCalendarEvent(event *Event) (*entities.Event, error) {
	calendarEvent, err := event.ConvertToCalendarEvent()
	if err != nil {
		return nil, fmt.Errorf("couldn't convert http.Event to calender.Event: %w", err)
	}
	return calendarEvent, nil
}
