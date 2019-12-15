package http

import (
	"github.com/mitrickx/otus-golang-2019/28/calendar/internal/storage/memory"
	"reflect"
	"testing"
)

func TestNewCalendar(t *testing.T) {
	service := NewTestCalendar()
	if service.storage == nil {
		t.Error("Storage must not be nil")
	}
}

func TestAddEvent(t *testing.T) {
	service := NewTestCalendar()

	if service.getEventsTotalCount() > 0 {
		t.Error("new entities service must not has events")
	}

	event1 := &Event{
		Name:  "Do homework",
		Start: "2019-10-15 20:00",
		End:   "2019-10-15 22:00",
	}

	id := addEvent(t, service, event1, 1)
	if id <= 0 {
		return
	}

	event2 := &Event{
		Name:  "Watch movie",
		Start: "2019-10-15 22:00",
		End:   "2019-10-16 01:00",
	}

	id = addEvent(t, service, event2, 2)
	if id <= 0 {
		return
	}
}

func TestUpdateEvent(t *testing.T) {
	service := NewTestCalendar()

	event1 := &Event{
		Name:  "Do homework",
		Start: "2019-10-15 20:00",
		End:   "2019-10-15 22:00",
	}

	id := addEvent(t, service, event1, 1)
	if id <= 0 {
		return
	}

	event2 := &Event{
		Name:  "Watch movie",
		Start: "2019-10-15 22:00",
		End:   "2019-10-16 01:00",
	}

	err := service.UpdateEvent(id, event2)

	if err != nil {
		t.Errorf("must not be happened error on update: %s\n", err)
		return
	}

	event, found := service.GetEvent(id)
	if !found {
		t.Errorf("event with id = %d not found on entities service", id)
		return
	}

	if event.Name != "Watch movie" || event.Start != "2019-10-15 22:00" || event.End != "2019-10-16 01:00" {
		t.Errorf("\nevent info not updated\nexpected be:\n%+v\ngot:\n%+v\n", event2, event)
	}

	err = service.UpdateEvent(0, &Event{})
	if err == nil {
		t.Error("update by id = 0 must return error")
	}

	err = service.UpdateEvent(1000, &Event{})
	if err == nil {
		t.Error("update by id of not existed event must return error")
	}
}

func TestDeleteEvent(t *testing.T) {
	service := NewTestCalendar()

	event1 := &Event{
		Name:  "Do homework",
		Start: "2019-10-15 20:00",
		End:   "2019-10-15 22:00",
	}

	id1 := addEvent(t, service, event1, 1)
	if id1 <= 0 {
		return
	}

	event2 := &Event{
		Name:  "Watch movie",
		Start: "2019-10-15 22:00",
		End:   "2019-10-16 01:00",
	}

	id2 := addEvent(t, service, event2, 2)
	if id2 <= 0 {
		return
	}

	err := service.DeleteEvent(0)
	if err == nil {
		t.Error("delete by id = 0 must return error")
	}

	err = service.DeleteEvent(1000)
	if err == nil {
		t.Error("delete by id of not existed event must return error")
	}

	err = service.DeleteEvent(id1)
	if err != nil {
		t.Errorf("delete by id = %d must not return error: %s", id1, err)
	}

	if service.getEventsTotalCount() != 1 {
		t.Error("delete actually not happened")
	}

	err = service.DeleteEvent(id2)
	if err != nil {
		t.Errorf("delete by id = %d must not return error: %s", id2, err)
	}

	if service.getEventsTotalCount() != 0 {
		t.Error("delete actually not happened, after 2 delete calender must be empty")
	}
}

func TestGetEvents(t *testing.T) {
	calendar := NewTestCalendar()

	addFixedListOfEvents(t, calendar)

	if calendar.getEventsTotalCount() != 7 {
		t.Error("7 events must be in entities")
		return
	}

	allEvents, _ := calendar.GetAllEvents()
	if len(allEvents) != 7 {
		t.Error("7 events must be in entities and GetAllEvents must return all of them")
	}

	allEvents2, _ := calendar.GetEventsByPeriod("", "")
	if len(allEvents2) != 7 {
		t.Error("7 events must be in entities and GetEventsByTimestampsPeriod(nil, nil) must return all of them")
	}

	if !reflect.DeepEqual(allEvents, allEvents2) {
		t.Errorf("Sorting of allEvents slices must be the same")
	}

	eventTime := "2019-11-18 08:00"
	eventList, _ := calendar.GetEventsByPeriod("", eventTime)

	if len(eventList) != 1 {
		t.Errorf("Must be returned one event")
	}

	if eventList[0].Name != "Monday" {
		t.Errorf("Must be returned one event `Monday`")
	}

	eventTime = "2019-11-24 08:00"
	eventList, _ = calendar.GetEventsByPeriod(eventTime, "")

	if len(eventList) != 1 {
		t.Errorf("Must be returned one event")
	}

	if eventList[0].Name != "Sunday" {
		t.Errorf("Must be returned one event `Sunday`")
	}

	startEventTime := "2019-11-20 08:00"
	endEventTime := "2019-11-22 08:00"
	eventList, _ = calendar.GetEventsByPeriod(startEventTime, endEventTime)

	if len(eventList) != 3 {
		t.Errorf("Must be returned 3 events")
	}

	if eventList[0].Name != "Wednesday" {
		t.Errorf("First event must be `Wednesday`")
	}
	if eventList[1].Name != "Thursday" {
		t.Errorf("First event must be `Thursday`")
	}
	if eventList[2].Name != "Friday" {
		t.Errorf("First event must be `Friday`")
	}

	startEventTime = "2019-11-20 08:01"
	endEventTime = "2019-11-20 09:59"
	eventList, _ = calendar.GetEventsByPeriod(startEventTime, endEventTime)

	if len(eventList) != 0 {
		t.Errorf("Must be returned 0 events")
	}
}

// Helper to add Event and run through list of assets
// Need for reduce code duplication in these tests
// If expectedCount input argument is greater and equal 0 check getEventsTotalCount of service
// If adding is successful return int ID > 0, otherwise return int < 0
func addEvent(t *testing.T, calendar *Calendar, event *Event, expectedCount int) int {
	id, err := calendar.AddEvent(event)

	if err != nil {
		t.Errorf("must not be error if add new event: %s", err)
		return -1
	} else if id <= 0 {
		t.Errorf("id %d of new added event must be > 0", id)
		return -2
	} else if expectedCount >= 0 && calendar.getEventsTotalCount() != expectedCount {
		t.Errorf("entities entities must has %d events instead of %d\n", expectedCount, calendar.getEventsTotalCount())
		return -3
	}

	return id
}

// Helper for tests that tests get event list method
func addFixedListOfEvents(t *testing.T, calendar *Calendar) {
	addEvent(t, calendar, &Event{
		Name:  "Monday",
		Start: "2019-11-18 08:00",
		End:   "2019-11-18 10:00",
	}, 1)

	addEvent(t, calendar, &Event{
		Name:  "Tuesday",
		Start: "2019-11-19 08:00",
		End:   "2019-11-19 10:00",
	}, 2)

	addEvent(t, calendar, &Event{
		Name:  "Wednesday",
		Start: "2019-11-20 08:00",
		End:   "2019-11-20 10:00",
	}, 3)

	addEvent(t, calendar, &Event{
		Name:  "Thursday",
		Start: "2019-11-21 08:00",
		End:   "2019-11-21 10:00",
	}, 4)

	addEvent(t, calendar, &Event{
		Name:  "Friday",
		Start: "2019-11-22 08:00",
		End:   "2019-11-22 10:00",
	}, 5)

	addEvent(t, calendar, &Event{
		Name:  "Saturday",
		Start: "2019-11-23 08:00",
		End:   "2019-11-23 10:00",
	}, 6)

	addEvent(t, calendar, &Event{
		Name:  "Sunday",
		Start: "2019-11-24 08:00",
		End:   "2019-11-24 10:00",
	}, 7)

	if calendar.getEventsTotalCount() != 7 {
		t.Error("7 events must be in entities")
		return
	}
}

func NewTestCalendar() *Calendar {
	storage := memory.NewStorage()
	calendar, _ := NewCalendar(storage)
	return calendar
}
