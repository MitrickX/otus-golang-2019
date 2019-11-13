package http

import (
	"fmt"
	"github.com/mitrickx/otus-golang-2019/21/calendar/internal/calendar"
)

type CalendarService struct {
	storage *calendar.Calendar
}

func NewCalendarService() *CalendarService {
	return &CalendarService{
		storage: calendar.NewCalendar(),
	}
}

func (calendarService *CalendarService) AddEvent(event *Event) (int, error) {
	calendarEvent, err := convertToCalendarEvent(event)
	if err != nil {
		return 0, err
	}

	return calendarService.storage.AddEvent(*calendarEvent), nil
}

func (calendarService *CalendarService) UpdateEvent(id int, event *Event) error {
	calendarEvent, err := convertToCalendarEvent(event)
	if err != nil {
		return err
	}

	err = calendarService.storage.UpdateEvent(id, *calendarEvent)
	if err != nil {
		return fmt.Errorf("couldn't update event in storage: %w", err)
	}
	return nil
}

func (calendarService *CalendarService) DeleteEvent(id int) error {
	err := calendarService.storage.DeleteEvent(id)
	if err != nil {
		return fmt.Errorf("couldn't delete event from storage: %w", err)
	}
	return nil
}

func (calendarService *CalendarService) GetEvent(id int) (*Event, bool) {
	if id <= 0 {
		return nil, false
	}

	calendarEvent, ok := calendarService.storage.GetEvent(id)
	if !ok {
		return nil, false
	}

	event := ConvertFromCalendarEvent(calendarEvent)
	return event, true

}

func (calendarService *CalendarService) getEventsTotalCount() int {
	return calendarService.storage.Count()
}

func convertToCalendarEvent(event *Event) (*calendar.Event, error) {
	calendarEvent, err := event.ConvertToCalendarEvent()
	if err != nil {
		return nil, fmt.Errorf("couldn't convert http.Event to calender.Event: %w", err)
	}
	return calendarEvent, nil
}