package entities

import (
	"fmt"
	"time"
)

// Simplest event struct, not support all day and repeat properties
type Event struct {
	id                 int       // id of event, need for identify event in entities
	name               string    // name of event
	start              DateTime  // start event time
	end                DateTime  // end event time
	isNotifyingEnabled bool      // need or not isNotifyingEnabled
	beforeMinutes      int       // isNotifyingEnabled before minutes
	isNotified         bool      // was notification enqueued
	notifiedTime       time.Time // when notification enqueued
}

// Constructor
func NewEvent(name string, start DateTime, end DateTime) Event {
	event := Event{
		name:  name,
		start: start,
		end:   end,
	}
	return event
}

// Clone constructor with setting ID
func WithId(event Event, id int) Event {
	event.id = id
	return event
}

// Constructor for existing in entities events
func NewEventWithId(id int, name string, start DateTime, end DateTime) Event {
	event := Event{
		id:    id,
		name:  name,
		start: start,
		end:   end,
	}
	return event
}

// Constructor for event all fields
func NewDetailedEvent(
	name string,
	start DateTime,
	end DateTime,
	isNotifyingEnable bool,
	beforeMinutes int,
	isNotified bool,
	notifiedTime time.Time,
) Event {

	return Event{
		name:               name,
		start:              start,
		end:                end,
		isNotifyingEnabled: isNotifyingEnable,
		beforeMinutes:      beforeMinutes,
		isNotified:         isNotified,
		notifiedTime:       notifiedTime,
	}
}

// Constructor for event all fields AND id
func NewDetailedEventWithId(
	id int,
	name string,
	start DateTime,
	end DateTime,
	isNotifyingEnable bool,
	beforeMinutes int,
	isNotified bool,
	notifiedTime time.Time,
) Event {

	return Event{
		id:                 id,
		name:               name,
		start:              start,
		end:                end,
		isNotifyingEnabled: isNotifyingEnable,
		beforeMinutes:      beforeMinutes,
		isNotified:         isNotified,
		notifiedTime:       notifiedTime,
	}
}

// Id of event getter
func (event Event) Id() int {
	return event.id
}

// Name of event getter
func (event Event) Name() string {
	return event.name
}

// Start of event getter
func (event Event) Start() DateTime {
	return event.start
}

// End of event getter
func (event Event) End() DateTime {
	return event.end
}

//
func (event Event) IsNotifyingEnabled() bool {
	return event.isNotifyingEnabled
}

//
func (event Event) BeforeMinutes() int {
	return event.beforeMinutes
}

//
func (event Event) IsNotified() bool {
	return event.isNotified
}

func (event Event) NotifiedTime() time.Time {
	return event.notifiedTime
}

func (event Event) Notified(t time.Time) Event {
	event.isNotified = true
	event.notifiedTime = t
	return event
}

// Less method for compare 2 event, will need for sorting in entities
func (event Event) Less(thatEvent Event) bool {
	if event.start != thatEvent.start {
		return event.start.Less(thatEvent.start)
	} else {
		return event.end.Less(thatEvent.end)
	}
}

// String representation of event
func (event Event) String() string {
	return fmt.Sprintf("%s: %s -> %s", event.name, event.start, event.end)
}
