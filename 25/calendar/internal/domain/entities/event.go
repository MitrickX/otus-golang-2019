package entities

import (
	"fmt"
)

// Simplest event struct, not support all day and repeat properties
type Event struct {
	id    int       // id of event, need for identify event in entities
	name  string    // name of event
	start EventTime // start event time
	end   EventTime // end event time
}

// Constructor
func NewEvent(name string, start EventTime, end EventTime) Event {
	event := Event{
		name:  name,
		start: start,
		end:   end,
	}
	return event
}

// Constructor for existing in entities events
func NewEventWithId(id int, name string, start EventTime, end EventTime) Event {
	event := Event{
		id:    id,
		name:  name,
		start: start,
		end:   end,
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
func (event Event) Start() EventTime {
	return event.start
}

// End of event getter
func (event Event) End() EventTime {
	return event.end
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
