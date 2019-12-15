package memory

import (
	"errors"
	"github.com/mitrickx/otus-golang-2019/28/calendar/internal/domain/entities"
	"sort"
	"sync"
	"time"
)

// Simplest entities struct, not support all day and repeat properties inherent for more sophisticated entities
type Storage struct {
	events        map[int]entities.Event // map of events indexed by id
	mx            sync.RWMutex           // rw mutex for safe concurrent read and modification of entities
	autoincrement int                    // autoincrement counter to generate next id on adding event in entities
}

// Constructor
func NewStorage() *Storage {
	calendar := &Storage{
		events: make(map[int]entities.Event),
		mx:     sync.RWMutex{},
	}
	return calendar
}

// Add event in entities, return new id for identify event in entities
func (calendar *Storage) AddEvent(event entities.Event) (int, error) {
	calendar.mx.Lock()
	calendar.autoincrement++
	id := calendar.autoincrement
	calendar.events[id] = entities.WithId(event, id)
	calendar.mx.Unlock()
	return id, nil
}

// Update event
// Get id and new event struct (inner id of event will be ignored)
// If not found returns error
func (calendar *Storage) UpdateEvent(id int, event entities.Event) error {

	if id <= 0 {
		return errors.New("event not found")
	}

	calendar.mx.RLock()
	_, ok := calendar.events[id]
	calendar.mx.RUnlock()

	if !ok {
		return errors.New("event not found")
	}

	newEvent := entities.WithId(event, id)

	calendar.mx.Lock()
	calendar.events[id] = newEvent
	calendar.mx.Unlock()

	return nil
}

// Delete event from entities by id of event in entities
// If not found returns error
func (calendar *Storage) DeleteEvent(id int) error {
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

// Get event by id of event in entities
// 2d param says found or not
func (calendar *Storage) GetEvent(id int) (entities.Event, error) {
	if id <= 0 {
		return entities.Event{}, entities.StorageErrorEventNotFound
	}

	calendar.mx.RLock()
	event, ok := calendar.events[id]
	calendar.mx.RUnlock()

	if !ok {
		return entities.Event{}, entities.StorageErrorEventNotFound
	}

	return event, nil
}

// Get all events of entities sorted by Less method of events
func (calendar *Storage) GetAllEvents() ([]entities.Event, error) {
	calendar.mx.RLock()
	eventsMap := calendar.events
	calendar.mx.RUnlock()

	if len(eventsMap) <= 0 {
		return nil, nil
	}

	events := make([]entities.Event, 0, len(eventsMap))
	for _, event := range eventsMap {
		events = append(events, event)
	}

	sort.Slice(events, func(i, j int) bool {
		return events[i].Less(events[j])
	})

	return events, nil
}

// Get all events that started in period (boundary of period are included) sorted by Less method of events
// You also can pass nil for start or end times
// nil has special means - no boundary for range period
func (calendar *Storage) GetEventsByPeriod(startTime *entities.DateTime, endTime *entities.DateTime) ([]entities.Event, error) {
	calendar.mx.RLock()
	eventsMap := calendar.events
	calendar.mx.RUnlock()

	if len(eventsMap) <= 0 {
		return nil, nil
	}

	var events []entities.Event
	for _, event := range eventsMap {
		inPeriod := true
		if startTime != nil && !startTime.LessOrEqual(event.Start()) {
			inPeriod = false
		}
		if endTime != nil && !event.Start().LessOrEqual(*endTime) {
			inPeriod = false
		}
		if inPeriod {
			events = append(events, event)
		}
	}

	sort.Slice(events, func(i, j int) bool {
		return events[i].Less(events[j])
	})

	return events, nil
}

//
func (calendar *Storage) GetEventsToNotify(startTime *entities.DateTime, endTime *entities.DateTime) ([]entities.Event, error) {
	allEvents, err := calendar.GetAllEvents()

	if err != nil {
		return nil, err
	}

	if len(allEvents) == 0 {
		return nil, nil
	}

	var events []entities.Event
	for _, event := range allEvents {

		if !event.IsNotifyingEnabled() || event.IsNotified() {
			continue
		}

		beforeMinutes := event.BeforeMinutes()
		start := event.Start()

		notifyT := start.Time().Add(-time.Duration(beforeMinutes) * time.Minute)
		notifyTime := entities.ConvertFromTime(notifyT)

		inPeriod := true
		if startTime != nil && !startTime.LessOrEqual(notifyTime) {
			inPeriod = false
		}
		if endTime != nil && !notifyTime.LessOrEqual(*endTime) {
			inPeriod = false
		}

		if inPeriod {
			events = append(events, event)
		}
	}

	sort.Slice(events, func(i, j int) bool {
		return events[i].Less(events[j])
	})

	return events, nil
}

// Mark event as notified and set time when was notified
func (calendar *Storage) MarkEventAsNotified(id int, when time.Time) error {
	event, err := calendar.GetEvent(id)
	if err != nil {
		return err
	}
	newEvent := event.Notified(when)
	return calendar.UpdateEvent(id, newEvent)
}

// Total number of events now in entities
func (calendar *Storage) Count() (int, error) {
	calendar.mx.RLock()
	count := len(calendar.events)
	calendar.mx.RUnlock()
	return count, nil
}

func (calendar *Storage) ClearAll() error {
	calendar.mx.Lock()
	calendar.events = make(map[int]entities.Event)
	calendar.mx.Unlock()
	return nil
}
