package grpc

import (
	"github.com/golang/protobuf/ptypes"
	"github.com/mitrickx/otus-golang-2019/22/calendar/internal/calendar"
	"testing"
)

const datetimeFormat = "2006-01-02 15:04"

func TestConvertToCalendarEventTime(t *testing.T) {
	tspb := ts(2019, 11, 15, 22, 0)
	eventTime, err := convertToCalendarEventTime(tspb)
	if err != nil {
		t.Errorf("Must not be error on converting %s", err)
	}

	strVal := eventTime.Format(datetimeFormat)
	expected := "2019-11-15 22:00"
	if strVal != expected {
		t.Errorf("Must be %s instead of %s", expected, strVal)
	}
}

func TestConvertToCalendarEvent(t *testing.T) {
	event := &Event{
		Id:    100,
		Name:  "Test",
		Start: ts(2019, 11, 15, 22, 0),
		End:   ts(2019, 11, 16, 1, 0),
	}

	calendarEvent, err := convertToCalendarEvent(event)

	if err != nil {
		t.Errorf("Must not be error on converting %s", err)
	}

	if calendarEvent.Id() != int(event.Id) {
		t.Errorf("Id must be %d instead of %d", event.Id, calendarEvent.Id())
	}

	if calendarEvent.Name() != event.Name {
		t.Errorf("Name must be %s instead of %s", event.Name, calendarEvent.Name())
	}

	expectedStart := "2019-11-15 22:00"
	start := calendarEvent.Start().Format(datetimeFormat)
	if start != expectedStart {
		t.Errorf("Start must be %s instead of %s", expectedStart, start)
	}

	expectedEnd := "2019-11-16 01:00"
	end := calendarEvent.End().Format(datetimeFormat)
	if end != expectedEnd {
		t.Errorf("Start must be %s instead of %s", expectedEnd, end)
	}
}

func TestConvertFromCalendarEvent(t *testing.T) {
	calendarEvent := calendar.NewEventWithId(
		100,
		"Test",
		calendar.NewEventTime(2019, 11, 19, 0, 12),
		calendar.NewEventTime(2019, 11, 19, 0, 50),
	)
	event, err := convertFromCalendarEvent(calendarEvent)
	if err != nil {
		t.Errorf("Must not be error on converting %s", err)
	}

	if calendarEvent.Id() != int(event.Id) {
		t.Errorf("Id must be %d instead of %d", calendarEvent.Id(), event.Id)
	}

	if event.Name != calendarEvent.Name() {
		t.Errorf("Name must be %s instead of %s", calendarEvent.Name(), event.Name)
	}

	start, err := ptypes.Timestamp(event.Start)
	if err != nil {
		t.Errorf("Must not be error on time  %s", err)
	}
	startTime := start.Format(datetimeFormat)
	expectedStart := "2019-11-19 00:12"
	if startTime != expectedStart {
		t.Errorf("Name must be %s instead of %s", expectedStart, startTime)
	}

	end, err := ptypes.Timestamp(event.End)
	if err != nil {
		t.Errorf("Must not be error on time  %s", err)
	}
	endTime := end.Format(datetimeFormat)
	expectedEnd := "2019-11-19 00:50"
	if endTime != expectedEnd {
		t.Errorf("Name must be %s instead of %s", expectedEnd, endTime)
	}
}

func TestIsEventEquals(t *testing.T) {
	event1 := &Event{
		Id:    100,
		Name:  "Test",
		Start: ts(2019, 11, 15, 22, 0),
		End:   ts(2019, 11, 16, 1, 0),
	}

	event2 := &Event{
		Id:    100,
		Name:  "Test",
		Start: ts(2019, 11, 15, 22, 0),
		End:   ts(2019, 11, 16, 1, 0),
	}

	event3 := &Event{
		Id:    101,
		Name:  "Test",
		Start: ts(2019, 11, 15, 22, 0),
		End:   ts(2019, 11, 16, 1, 0),
	}

	if !isEventEquals(event1, event1, false) {
		t.Error("Even1 must be equal event1")
	}

	if !isEventEquals(event1, event2, true) {
		t.Error("Even1 must be equal event1")
	}

	if !isEventEquals(event1, event2, false) {
		t.Error("Even1 must be equal event2")
	}

	if !isEventEquals(event1, event2, true) {
		t.Error("Even1 must be equal event2 (deep)")
	}

	if !isEventEquals(event2, event3, false) {
		t.Error("Even2 must be equal event3")
	}

	if isEventEquals(event2, event3, true) {
		t.Error("Even2 must NOT be equal event3 (deep)")
	}
}
