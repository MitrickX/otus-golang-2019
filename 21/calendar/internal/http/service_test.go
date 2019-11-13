package http

import (
	"testing"
)

func TestNewCalendarService(t *testing.T) {
	service := NewCalendarService()
	if service.storage == nil {
		t.Error("Storage must not be nil")
	}
}

func TestAddEvent(t *testing.T) {
	service := NewCalendarService()

	if service.getEventsTotalCount() > 0 {
		t.Error("new calendar service must not has events")
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
	service := NewCalendarService()

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
		t.Errorf("event with id = %d not found on calendar service", id)
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
	service := NewCalendarService()

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

// Helper to add Event and run through list of assets
// Need for reduce code duplication in these tests
// If expectedCount input argument is greater and equal 0 check getEventsTotalCount of service
// If adding is successful return int ID > 0, otherwise return int < 0
func addEvent(t *testing.T, service *CalendarService, event *Event, expectedCount int) int {
	id, err := service.AddEvent(event)

	if err != nil {
		t.Errorf("must not be error if add new event: %s", err)
		return -1
	} else if id <= 0 {
		t.Errorf("id %d of new added event must be > 0", id)
		return -2
	} else if expectedCount >= 0 && service.getEventsTotalCount() != expectedCount {
		t.Errorf("calendar service must has %d events instead of %d\n", expectedCount, service.getEventsTotalCount())
		return -3
	}

	return id
}
