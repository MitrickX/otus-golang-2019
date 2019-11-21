package grpc

import (
	"errors"
	"fmt"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/mitrickx/otus-golang-2019/22/calendar/internal/calendar"
	"strings"
)

var ErrorNotFound = errors.New("event not found")

type ErrorEventListErrors struct {
	errs []error
}

func (e *ErrorEventListErrors) Error() string {
	buffer := strings.Builder{}
	buffer.WriteString("Some errors happened when get event list: ")
	first := true
	for _, er := range e.errs {
		if first {
			buffer.WriteString(fmt.Sprintf("%w", er))
			first = false
		} else {
			buffer.WriteString(fmt.Sprintf(", %w", er))
		}

	}
	return buffer.String()
}

// Calendar structure for work inside grpc package
// Clean architecture approach - not working with inner biz logic layer directly
type Calendar struct {
	storage *calendar.Calendar // for now it is inner biz entity itself, for future there will be storage interface
}

// Constructor
func NewCalendar(storage *calendar.Calendar) *Calendar {
	return &Calendar{
		storage: storage,
	}
}

// Default calender with in-memory storage
func NewDefaultCalendar() *Calendar {
	storage := calendar.NewCalendar()
	return NewCalendar(storage)
}

// Add Event
func (c *Calendar) AddEvent(event *Event) (int, error) {
	calendarEvent, err := convertToCalendarEvent(event)
	if err != nil {
		return 0, err
	}
	return c.storage.AddEvent(*calendarEvent), nil
}

// Update Event
func (c *Calendar) UpdateEvent(id int, event *Event) error {
	calendarEvent, err := convertToCalendarEvent(event)
	if err != nil {
		return err
	}

	err = c.storage.UpdateEvent(id, *calendarEvent)
	if err != nil {
		return fmt.Errorf("couldn't update event in storage: %w", err)
	}
	return nil
}

// Delete Event
func (c *Calendar) DeleteEvent(id int) error {
	err := c.storage.DeleteEvent(id)
	if err != nil {
		return fmt.Errorf("couldn't delete event from storage: %w", err)
	}
	return nil
}

// Get one event
func (c *Calendar) GetEvent(id int) (*Event, error) {
	if id <= 0 {
		return nil, ErrorNotFound
	}

	calendarEvent, ok := c.storage.GetEvent(id)
	if !ok {
		return nil, ErrorNotFound
	}

	event, err := convertFromCalendarEvent(calendarEvent)
	if err != nil {
		return nil, err
	}

	return event, nil

}

// Get all events
func (c *Calendar) GetAllEvents() ([]*Event, error) {
	calendarEvents := c.storage.GetAllEvents()
	if len(calendarEvents) == 0 {
		return nil, nil
	}

	var convertErrors []error
	var events []*Event

	for _, calendarEvent := range calendarEvents {
		event, err := convertFromCalendarEvent(calendarEvent)
		if err != nil {
			convertErrors = append(convertErrors, err)
		} else {
			events = append(events, event)
		}
	}
	listErr := &ErrorEventListErrors{
		errs: convertErrors,
	}
	return events, listErr
}

// Get all events that started in period (*Period struct) sorted by Less method of events
// Nils has special meaning - no boundary for range period
// Return slice of events and slice of errors
// Method try return max events that could be returned
func (c *Calendar) GetEventsByPeriod(period *Period) ([]*Event, error) {
	if period == nil {
		return c.GetEventsByTimestampsPeriod(nil, nil)
	} else {
		return c.GetEventsByTimestampsPeriod(period.start, period.end)
	}
}

// Get all events that started in period (boundary of period are included) sorted by Less method of events
// start/end are datetime values represented by string in format on this module (*timestamp.Timestamp)
// Nil has special meaning - no boundary for range period
// Return slice of events and slice of errors
// Method try return max events that could be returned
func (c *Calendar) GetEventsByTimestampsPeriod(start *timestamp.Timestamp, end *timestamp.Timestamp) ([]*Event, error) {
	var startTime, endTime *calendar.EventTime

	if start != nil {
		var err error
		startTime, err = convertToCalendarEventTime(start)
		if err != nil {
			return nil, err
		}
	}

	if end != nil {
		var err error
		endTime, err = convertToCalendarEventTime(end)
		if err != nil {
			return nil, err
		}
	}

	calendarEvents := c.storage.GetEventsByPeriod(startTime, endTime)
	if len(calendarEvents) == 0 {
		return nil, nil
	}

	var events []*Event
	var convertErrors []error
	var listErr error

	for _, calendarEvent := range calendarEvents {
		event, err := convertFromCalendarEvent(calendarEvent)
		if err != nil {
			convertErrors = append(convertErrors, err)
		} else {
			events = append(events, event)
		}
	}

	if len(convertErrors) > 0 {
		listErr = &ErrorEventListErrors{
			errs: convertErrors,
		}
	}

	return events, listErr
}

// Get total number of events in calendar
func (c *Calendar) getEventsTotalCount() int {
	return c.storage.Count()
}
