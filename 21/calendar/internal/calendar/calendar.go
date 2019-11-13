package calendar

import (
	"errors"
	"sort"
	"sync"
)

// Simplest calendar struct, not support all day and repeat properties inherent for more sophisticated calendar
type Calendar struct {
	events        map[int]Event // map of events indexed by id
	mx            sync.RWMutex  // rw mutex for safe concurrent read and modification of calendar
	autoincrement int           // autoincrement counter to generate next id on adding event in calendar
}

// Constructor
func NewCalendar() *Calendar {
	calendar := &Calendar{
		events: make(map[int]Event),
		mx:     sync.RWMutex{},
	}
	return calendar
}

// Add event in calendar, return new id for identify event in calendar
func (calendar *Calendar) AddEvent(event Event) int {
	calendar.mx.Lock()
	calendar.autoincrement++
	id := calendar.autoincrement
	calendar.events[id] = event
	calendar.mx.Unlock()
	return id
}

// Update event
// Get id and new event struct (inner id of event will be ignored)
// If not found returns error
func (calendar *Calendar) UpdateEvent(id int, event Event) error {

	if id <= 0 {
		return errors.New("event not found")
	}

	calendar.mx.RLock()
	_, ok := calendar.events[id]
	calendar.mx.RUnlock()

	if !ok {
		return errors.New("event not found")
	}

	calendar.mx.Lock()
	event.id = id
	calendar.events[id] = event
	calendar.mx.Unlock()

	return nil
}

// Delete event from calendar by id of event in calendar
// If not found returns error
func (calendar *Calendar) DeleteEvent(id int) error {
	if id <= 0 {
		return errors.New("event not found")
	}

	calendar.mx.RLock()
	_, ok := calendar.events[id]
	calendar.mx.RUnlock()

	if !ok {
		return errors.New("event not found")
	}

	calendar.mx.Lock()
	delete(calendar.events, id)
	calendar.mx.Unlock()

	return nil
}

// Get event by id of event in calendar
// 2d param says found or not
func (calendar *Calendar) GetEvent(id int) (Event, bool) {
	if id <= 0 {
		return Event{}, false
	}

	calendar.mx.RLock()
	event, ok := calendar.events[id]
	calendar.mx.RUnlock()

	if !ok {
		return Event{}, false
	}

	return event, true
}

// Get all events of calendar sorted by Less method of events
func (calendar *Calendar) GetAllEvents() []Event {
	calendar.mx.RLock()
	eventsMap := calendar.events
	calendar.mx.RUnlock()

	if len(eventsMap) <= 0 {
		return nil
	}

	events := make([]Event, 0, len(eventsMap))
	for _, event := range eventsMap {
		events = append(events, event)
	}

	sort.Slice(events, func(i, j int) bool {
		return events[i].Less(events[j])
	})

	return events
}

// Total number of events now in calendar
func (calendar *Calendar) Count() int {
	calendar.mx.RLock()
	count := len(calendar.events)
	calendar.mx.RUnlock()
	return count
}
