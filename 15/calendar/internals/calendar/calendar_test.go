package calendar

import (
	"testing"
)

// Test adding new event in calendar
func TestAddEvent(t *testing.T) {
	calendar := NewCalendar()

	if calendar.Count() > 0 {
		t.Error("empty calendar must not has events")
	}

	event1 := NewEvent("Do homework",
		NewEventTime(2019, 10, 15, 20, 0),
		NewEventTime(2019, 10, 15, 22, 0),
	)

	id := calendar.AddEvent(event1)

	if id <= 0 {
		t.Errorf("id %d of new added event must be > 0", id)
	}

	if calendar.Count() != 1 {
		t.Errorf("calendar must has 1 event instead of %d\n", calendar.Count())
	}

	event2 := NewEvent("Watch movie",
		NewEventTime(2019, 10, 15, 22, 0),
		NewEventTime(2019, 10, 16, 1, 0),
	)

	calendar.AddEvent(event2)

	if calendar.Count() != 2 {
		t.Errorf("calendar must has 2 event instead of %d\n", calendar.Count())
	}
}

// Test get events from calendar
func TestGetEvent(t *testing.T) {
	calendar := NewCalendar()

	event1 := NewEvent("Do homework",
		NewEventTime(2019, 10, 15, 20, 0),
		NewEventTime(2019, 10, 15, 22, 0),
	)

	id := calendar.AddEvent(event1)

	event, ok := calendar.GetEvent(id)
	if !ok {
		t.Error("Get Event must be ok")
	}

	if event.Name() != "Do homework" {
		t.Error("get Event return another event")
	}

	_, ok = calendar.GetEvent(10000)
	if ok {
		t.Error("get Event must not be ok")
	}

	_, ok = calendar.GetEvent(0)
	if ok {
		t.Error("get Event must not be ok")
	}
}

// Test of updating events in calendar
func TestUpdateEvent(t *testing.T) {
	calendar := NewCalendar()

	event1 := NewEvent("Do homework",
		NewEventTime(2019, 10, 15, 20, 0),
		NewEventTime(2019, 10, 15, 22, 0),
	)

	id := calendar.AddEvent(event1)

	event2 := NewEvent("Watch movie",
		NewEventTime(2019, 10, 15, 22, 0),
		NewEventTime(2019, 10, 16, 1, 0),
	)

	err := calendar.UpdateEvent(id, event2)

	if err != nil {
		t.Errorf("update err %s must not be happened\n", err)
	}

	event, _ := calendar.GetEvent(id)

	if event.Name() != "Watch movie" {
		t.Error("get return another event, event actually not updated")
	}

	err = calendar.UpdateEvent(0, Event{})
	if err == nil {
		t.Error("update by id = 0 must return error")
	}

	err = calendar.UpdateEvent(1000, Event{})
	if err == nil {
		t.Error("update by id of not existed event must return error")
	}
}

// Test deleting events in calendar
func TestDeleteEvent(t *testing.T) {
	calendar := NewCalendar()

	event1 := NewEvent("Do homework",
		NewEventTime(2019, 10, 15, 20, 0),
		NewEventTime(2019, 10, 15, 22, 0),
	)

	id1 := calendar.AddEvent(event1)

	event2 := NewEvent("Watch movie",
		NewEventTime(2019, 10, 15, 22, 0),
		NewEventTime(2019, 10, 16, 1, 0),
	)

	id2 := calendar.AddEvent(event2)

	err := calendar.DeleteEvent(0)
	if err == nil {
		t.Error("delete by id = 0 must return error")
	}

	err = calendar.DeleteEvent(1000)
	if err == nil {
		t.Error("delete by id of not existed event must return error")
	}

	err = calendar.DeleteEvent(id1)
	if err != nil {
		t.Errorf("delete by id = %d must not return error", id1)
	}

	if calendar.Count() != 1 {
		t.Error("delete actually not happened")
	}

	err = calendar.DeleteEvent(id2)

	if calendar.Count() != 0 {
		t.Error("delete actually not happened, after 2 delete calender must be empty")
	}
}
