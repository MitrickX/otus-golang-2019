package entities

import (
	"errors"
)

var StorageErrorEventNotFound = errors.New("event not found in storage")

type Storage interface {
	AddEvent(event Event) (int, error)
	UpdateEvent(id int, event Event) error
	DeleteEvent(id int) error
	GetEvent(id int) (Event, error)
	GetAllEvents() ([]Event, error)
	GetEventsByPeriod(startTime *DateTime, endTime *DateTime) ([]Event, error)
	GetEventsForNotification(startTime *DateTime, endTime *DateTime) ([]Event, error)
	Count() (int, error)
	ClearAll() error
}
