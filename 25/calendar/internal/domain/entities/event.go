package entities

import (
	"fmt"
)

// Simplest event struct, not support all day and repeat properties
type Event struct {
	id            int      // id of event, need for identify event in entities
	name          string   // name of event
	start         DateTime // start event time
	end           DateTime // end event time
	notify        bool     // need or not notify
	beforeMinutes int      // notify before minutes
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

// Constructor for event with notification parameters
func NewNotifiedEvent(name string, start DateTime, end DateTime, beforeMinutes int) Event {
	event := Event{
		name:          name,
		start:         start,
		end:           end,
		notify:        true,
		beforeMinutes: beforeMinutes,
	}
	return event
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
	return event.notify
}

//
func (event Event) BeforeMinutes() int {
	return event.beforeMinutes
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
