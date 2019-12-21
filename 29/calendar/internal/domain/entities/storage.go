package entities

import (
	"errors"
	"time"
)

var StorageErrorEventNotFound = errors.New("event not found in storage")

type Storage interface {

	// Add event
	AddEvent(event Event) (int, error)

	// Update event
	UpdateEvent(id int, event Event) error

	// Delete event
	DeleteEvent(id int) error

	// Get one event by id
	GetEvent(id int) (Event, error)

	// Get all events
	GetAllEvents() ([]Event, error)

	// Get events by period. start and end is inclusive
	GetEventsByPeriod(startTime *DateTime, endTime *DateTime) ([]Event, error)

	// Get events to notify, start and end is inclusive
	GetEventsToNotify(startTime *DateTime, endTime *DateTime) ([]Event, error)

	// Mark event as notified, when is time when event is mark as notified
	MarkEventAsNotified(id int, when time.Time) error

	// Count of all events
	Count() (int, error)

	// Delete all events
	ClearAll() error
}
