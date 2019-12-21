package memory

import (
	"github.com/mitrickx/otus-golang-2019/29/calendar/internal/domain/entities"
	"reflect"
	"testing"
	"time"
)

// Test adding new event in entities
func TestAddEvent(t *testing.T) {
	calendar := NewStorage()

	if getCalendarCount(calendar) > 0 {
		t.Error("empty entities must not has events")
	}

	event1 := entities.NewEvent("Do homework",
		entities.NewDateTime(2019, 10, 15, 20, 0),
		entities.NewDateTime(2019, 10, 15, 22, 0),
	)

	id, _ := calendar.AddEvent(event1)

	if id <= 0 {
		t.Errorf("id %d of new added event must be > 0", id)
	}

	if getCalendarCount(calendar) != 1 {
		t.Errorf("entities must has 1 event instead of %d\n", getCalendarCount(calendar))
	}

	event2 := entities.NewEvent("Watch movie",
		entities.NewDateTime(2019, 10, 15, 22, 0),
		entities.NewDateTime(2019, 10, 16, 1, 0),
	)

	_, _ = calendar.AddEvent(event2)

	if getCalendarCount(calendar) != 2 {
		t.Errorf("entities must has 2 event instead of %d\n", getCalendarCount(calendar))
	}
}

// Test get events from entities
func TestGetEvent(t *testing.T) {
	calendar := NewStorage()

	event1 := entities.NewEvent("Do homework",
		entities.NewDateTime(2019, 10, 15, 20, 0),
		entities.NewDateTime(2019, 10, 15, 22, 0),
	)

	id, _ := calendar.AddEvent(event1)

	event, err := calendar.GetEvent(id)
	if err == entities.StorageErrorEventNotFound {
		t.Error("Get Event must be ok")
	}

	if event.Name() != "Do homework" {
		t.Error("get Event return another event")
	}

	_, err = calendar.GetEvent(10000)
	if err != entities.StorageErrorEventNotFound {
		t.Error("get Event must not be ok")
	}

	_, err = calendar.GetEvent(0)
	if err != entities.StorageErrorEventNotFound {
		t.Error("get Event must not be ok")
	}
}

// Test of updating events in entities
func TestUpdateEvent(t *testing.T) {
	calendar := NewStorage()

	event1 := entities.NewEvent("Do homework",
		entities.NewDateTime(2019, 10, 15, 20, 0),
		entities.NewDateTime(2019, 10, 15, 22, 0),
	)

	id, _ := calendar.AddEvent(event1)

	event2 := entities.NewEvent("Watch movie",
		entities.NewDateTime(2019, 10, 15, 22, 0),
		entities.NewDateTime(2019, 10, 16, 1, 0),
	)

	err := calendar.UpdateEvent(id, event2)

	if err != nil {
		t.Errorf("update err %s must not be happened\n", err)
	}

	event, _ := calendar.GetEvent(id)

	if event.Name() != "Watch movie" {
		t.Error("get return another event, event actually not updated")
	}

	err = calendar.UpdateEvent(0, entities.Event{})
	if err == nil {
		t.Error("update by id = 0 must return error")
	}

	err = calendar.UpdateEvent(1000, entities.Event{})
	if err == nil {
		t.Error("update by id of not existed event must return error")
	}
}

// Test deleting events in entities
func TestDeleteEvent(t *testing.T) {
	calendar := NewStorage()

	event1 := entities.NewEvent("Do homework",
		entities.NewDateTime(2019, 10, 15, 20, 0),
		entities.NewDateTime(2019, 10, 15, 22, 0),
	)

	id1, _ := calendar.AddEvent(event1)

	event2 := entities.NewEvent("Watch movie",
		entities.NewDateTime(2019, 10, 15, 22, 0),
		entities.NewDateTime(2019, 10, 16, 1, 0),
	)

	id2, _ := calendar.AddEvent(event2)

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

	if getCalendarCount(calendar) != 1 {
		t.Error("delete actually not happened")
	}

	_ = calendar.DeleteEvent(id2)

	if getCalendarCount(calendar) != 0 {
		t.Error("delete actually not happened, after 2 delete calender must be empty")
	}
}

func TestGetEvents(t *testing.T) {
	calendar := NewStorage()

	_, _ = calendar.AddEvent(entities.NewEvent("Monday",
		entities.NewDateTime(2019, 11, 18, 8, 0),
		entities.NewDateTime(2019, 11, 18, 10, 0),
	))

	_, _ = calendar.AddEvent(entities.NewEvent("Tuesday",
		entities.NewDateTime(2019, 11, 19, 8, 0),
		entities.NewDateTime(2019, 11, 19, 10, 0),
	))

	_, _ = calendar.AddEvent(entities.NewEvent("Wednesday",
		entities.NewDateTime(2019, 11, 20, 8, 0),
		entities.NewDateTime(2019, 11, 20, 10, 0),
	))

	_, _ = calendar.AddEvent(entities.NewEvent("Thursday",
		entities.NewDateTime(2019, 11, 21, 8, 0),
		entities.NewDateTime(2019, 11, 21, 10, 0),
	))

	_, _ = calendar.AddEvent(entities.NewEvent("Friday",
		entities.NewDateTime(2019, 11, 22, 8, 0),
		entities.NewDateTime(2019, 11, 22, 10, 0),
	))

	_, _ = calendar.AddEvent(entities.NewEvent("Saturday",
		entities.NewDateTime(2019, 11, 23, 8, 0),
		entities.NewDateTime(2019, 11, 23, 10, 0),
	))

	_, _ = calendar.AddEvent(entities.NewEvent("Sunday",
		entities.NewDateTime(2019, 11, 24, 8, 0),
		entities.NewDateTime(2019, 11, 24, 10, 0),
	))

	if getCalendarCount(calendar) != 7 {
		t.Error("7 events must be in entities")
		return
	}

	allEvents, _ := calendar.GetAllEvents()
	if len(allEvents) != 7 {
		t.Error("7 events must be in entities and GetAllEvents must return all of them")
	}

	allEvents2, _ := calendar.GetEventsByPeriod(nil, nil)
	if len(allEvents2) != 7 {
		t.Error("7 events must be in entities and GetEventsByTimestampsPeriod(nil, nil) must return all of them")
	}

	if !reflect.DeepEqual(allEvents, allEvents2) {
		t.Errorf("Sorting of allEvents slices must be the same")
	}

	eventTime := entities.NewDateTime(2019, 11, 18, 8, 0)
	eventList, _ := calendar.GetEventsByPeriod(nil, &eventTime)

	if len(eventList) != 1 {
		t.Errorf("Must be returned one event")
	}

	if eventList[0].Name() != "Monday" {
		t.Errorf("Must be returned one event `Monday`")
	}

	eventTime = entities.NewDateTime(2019, 11, 24, 8, 0)
	eventList, _ = calendar.GetEventsByPeriod(&eventTime, nil)

	if len(eventList) != 1 {
		t.Errorf("Must be returned one event")
	}

	if eventList[0].Name() != "Sunday" {
		t.Errorf("Must be returned one event `Sunday`")
	}

	startEventTime := entities.NewDateTime(2019, 11, 20, 8, 0)
	endEventTime := entities.NewDateTime(2019, 11, 22, 8, 0)
	eventList, _ = calendar.GetEventsByPeriod(&startEventTime, &endEventTime)

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

	startEventTime = entities.NewDateTime(2019, 11, 20, 8, 1)
	endEventTime = entities.NewDateTime(2019, 11, 20, 9, 59)
	eventList, _ = calendar.GetEventsByPeriod(&startEventTime, &endEventTime)

	if len(eventList) != 0 {
		t.Errorf("Must be returned 0 events")
	}
}

func TestGetEventsForNotification1(t *testing.T) {

	calendar := NewStorage()

	_, err := calendar.AddEvent(
		entities.NewDetailedEvent(
			"TestEvent",
			entities.NewDateTime(2019, 11, 18, 8, 0),
			entities.NewDateTime(2019, 11, 18, 10, 0),
			true,
			30,
			false,
			time.Time{},
		),
	)

	if err != nil {
		t.Errorf("Error while insert: %s", err)
		return
	}

	_, err = calendar.AddEvent(
		entities.NewDetailedEvent(
			"TestEvent2",
			entities.NewDateTime(2019, 11, 18, 8, 0),
			entities.NewDateTime(2019, 11, 18, 10, 0),
			true,
			10,
			false,
			time.Time{},
		),
	)

	if err != nil {
		t.Errorf("Error while insert: %s", err)
		return
	}

	start := entities.NewDateTime(2019, 11, 18, 7, 0)
	end := entities.NewDateTime(2019, 11, 18, 8, 0)

	entites, err := calendar.GetEventsToNotify(&start, &end)

	if err != nil {
		t.Errorf("Error while getting events: %s", err)
		return
	}

	if len(entites) != 2 {
		t.Errorf("Must be 2 events for notification instread of %d", len(entites))
	}

}

func TestGetEventsForNotification2(t *testing.T) {

	calendar := NewStorage()

	event1Id, err := calendar.AddEvent(
		entities.NewDetailedEvent(
			"TestEvent",
			entities.NewDateTime(2019, 11, 18, 8, 0),
			entities.NewDateTime(2019, 11, 18, 10, 0),
			true,
			30,
			false,
			time.Time{},
		),
	)

	if err != nil {
		t.Errorf("Error while insert: %s", err)
		return
	}

	_, err = calendar.AddEvent(
		entities.NewDetailedEvent(
			"TestEvent2",
			entities.NewDateTime(2019, 11, 18, 8, 0),
			entities.NewDateTime(2019, 11, 18, 10, 0),
			true,
			10,
			false,
			time.Time{},
		),
	)

	if err != nil {
		t.Errorf("Error while insert: %s", err)
		return
	}

	start := entities.NewDateTime(2019, 11, 18, 7, 0)
	end := entities.NewDateTime(2019, 11, 18, 7, 49)

	events, err := calendar.GetEventsToNotify(&start, &end)

	if err != nil {
		t.Errorf("Error while getting events: %s", err)
		return
	}

	if len(events) != 1 {
		t.Errorf("Must be 1 events for notification instread of %d", len(events))
		return
	}

	if events[0].Id() != event1Id {
		t.Errorf("Must be event #1 instead of #%d", events[0].Id())
	}

}

func TestGetEventsForNotification3(t *testing.T) {

	calendar := NewStorage()

	_, err := calendar.AddEvent(
		entities.NewDetailedEvent(
			"TestEvent",
			entities.NewDateTime(2019, 11, 18, 8, 0),
			entities.NewDateTime(2019, 11, 18, 10, 0),
			true,
			30,
			false,
			time.Time{},
		),
	)

	if err != nil {
		t.Errorf("Error while insert: %s", err)
		return
	}

	_, err = calendar.AddEvent(
		entities.NewDetailedEvent(
			"TestEvent2",
			entities.NewDateTime(2019, 11, 18, 8, 0),
			entities.NewDateTime(2019, 11, 18, 10, 0),
			true,
			10,
			false,
			time.Time{},
		),
	)

	if err != nil {
		t.Errorf("Error while insert: %s", err)
		return
	}

	start := entities.NewDateTime(2019, 11, 18, 8, 1)
	end := entities.NewDateTime(2019, 11, 18, 8, 59)

	events, err := calendar.GetEventsToNotify(&start, &end)

	if err != nil {
		t.Errorf("Error while getting events: %s", err)
		return
	}

	if len(events) != 0 {
		t.Errorf("Must be 1 events for notification instread of %d", len(events))
		return
	}

}

func TestGetEventsForNotification4(t *testing.T) {

	calendar := NewStorage()

	originalEvents := []entities.Event{
		entities.NewDetailedEvent(
			"A",
			entities.NewDateTime(2019, 11, 18, 8, 0),
			entities.NewDateTime(2019, 11, 18, 10, 0),
			true,
			30,
			false,
			time.Time{},
		),
		entities.NewDetailedEvent(
			"B",
			entities.NewDateTime(2019, 11, 19, 8, 0),
			entities.NewDateTime(2019, 11, 19, 10, 0),
			true,
			10,
			false,
			time.Time{},
		),
	}

	for index, event := range originalEvents {
		id, err := calendar.AddEvent(event)
		if err != nil {
			t.Errorf("Error while insert %s: %s", event.Name(), err)
			return
		}
		originalEvents[index] = entities.WithId(event, id)
	}

	start := originalEvents[0].Start().MinusMinutes(originalEvents[0].BeforeMinutes())
	end := originalEvents[1].Start().MinusMinutes(originalEvents[1].BeforeMinutes())

	events, err := calendar.GetEventsToNotify(&start, &end)

	if err != nil {
		t.Errorf("Error while getting events: %s", err)
		return
	}

	var names []string
	for _, event := range events {
		names = append(names, event.Name())
	}

	expectedNames := []string{"A", "B"}

	if !reflect.DeepEqual(expectedNames, names) {
		t.Errorf("Expected names %v insteadof %v", expectedNames, names)
	}

}

func TestMarkEventAsNotified(t *testing.T) {

	calendar := NewStorage()

	id, err := calendar.AddEvent(
		entities.NewDetailedEvent(
			"A",
			entities.NewDateTime(2019, 11, 18, 8, 0),
			entities.NewDateTime(2019, 11, 18, 10, 0),
			true,
			30,
			false,
			time.Time{},
		),
	)

	if err != nil {
		t.Errorf("Error while insert: %s", err)
		return
	}

	event, err := calendar.GetEvent(id)

	if err != nil {
		t.Errorf("Error while get by id %d: %s", id, err)
		return
	}

	if event.IsNotified() {
		t.Errorf("Event with id %d must not be marked as notified yet", id)
		return
	}

	dt := time.Date(2000, time.Month(10), 13, 14, 32, 18, 0, time.UTC)

	_ = calendar.MarkEventAsNotified(id, dt)

	event, err = calendar.GetEvent(id)

	if err != nil {
		t.Errorf("Error while get by id %d: %s", id, err)
		return
	}

	if !event.IsNotified() {
		t.Errorf("Event with id %d must be marked as notified", id)
		return
	}

	nTime := event.NotifiedTime()

	if !dt.Equal(nTime) {
		t.Errorf("Notified time of event with id %d not equal to original time", id)
	}
}

func getCalendarCount(storage *Storage) int {
	cnt, _ := storage.Count()
	return cnt
}
