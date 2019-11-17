package calendar

import (
	"reflect"
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

func TestGetEvents(t *testing.T) {
	calendar := NewCalendar()

	calendar.AddEvent(NewEvent("Monday",
		NewEventTime(2019, 11, 18, 8, 0),
		NewEventTime(2019, 11, 18, 10, 0),
	))

	calendar.AddEvent(NewEvent("Tuesday",
		NewEventTime(2019, 11, 19, 8, 0),
		NewEventTime(2019, 11, 19, 10, 0),
	))

	calendar.AddEvent(NewEvent("Wednesday",
		NewEventTime(2019, 11, 20, 8, 0),
		NewEventTime(2019, 11, 20, 10, 0),
	))

	calendar.AddEvent(NewEvent("Thursday",
		NewEventTime(2019, 11, 21, 8, 0),
		NewEventTime(2019, 11, 21, 10, 0),
	))

	calendar.AddEvent(NewEvent("Friday",
		NewEventTime(2019, 11, 22, 8, 0),
		NewEventTime(2019, 11, 22, 10, 0),
	))

	calendar.AddEvent(NewEvent("Saturday",
		NewEventTime(2019, 11, 23, 8, 0),
		NewEventTime(2019, 11, 23, 10, 0),
	))

	calendar.AddEvent(NewEvent("Sunday",
		NewEventTime(2019, 11, 24, 8, 0),
		NewEventTime(2019, 11, 24, 10, 0),
	))

	if calendar.Count() != 7 {
		t.Error("7 events must be in calendar")
		return
	}

	allEvents := calendar.GetAllEvents()
	if len(allEvents) != 7 {
		t.Error("7 events must be in calendar and GetAllEvents must return all of them")
	}

	allEvents2 := calendar.GetEventsByPeriod(nil, nil)
	if len(allEvents2) != 7 {
		t.Error("7 events must be in calendar and GetEventsByPeriod(nil, nil) must return all of them")
	}

	if !reflect.DeepEqual(allEvents, allEvents2) {
		t.Errorf("Sorting of allEvents slices must be the same")
	}

	eventTime := NewEventTime(2019, 11, 18, 8, 0)
	eventList := calendar.GetEventsByPeriod(nil, &eventTime)

	if len(eventList) != 1 {
		t.Errorf("Must be returned one event")
	}

	if eventList[0].Name() != "Monday" {
		t.Errorf("Must be returned one event `Monday`")
	}

	eventTime = NewEventTime(2019, 11, 24, 8, 0)
	eventList = calendar.GetEventsByPeriod(&eventTime, nil)

	if len(eventList) != 1 {
		t.Errorf("Must be returned one event")
	}

	if eventList[0].Name() != "Sunday" {
		t.Errorf("Must be returned one event `Sunday`")
	}

	startEventTime := NewEventTime(2019, 11, 20, 8, 0)
	endEventTime := NewEventTime(2019, 11, 22, 8, 0)
	eventList = calendar.GetEventsByPeriod(&startEventTime, &endEventTime)

	if len(eventList) != 3 {
		t.Errorf("Must be returned 3 events")
	}

	if eventList[0].Name() != "Wednesday" {
		t.Errorf("First event must be `Wednesday`")
	}
	if eventList[1].Name() != "Thursday" {
		t.Errorf("First event must be `Thursday`")
	}
	if eventList[2].Name() != "Friday" {
		t.Errorf("First event must be `Friday`")
	}

	startEventTime = NewEventTime(2019, 11, 20, 8, 1)
	endEventTime = NewEventTime(2019, 11, 20, 9, 59)
	eventList = calendar.GetEventsByPeriod(&startEventTime, &endEventTime)

	if len(eventList) != 0 {
		t.Errorf("Must be returned 0 events")
	}
}
