package http

import (
	"fmt"
	"github.com/mitrickx/otus-golang-2019/21/calendar/internal/calendar"
	"testing"
)

func TestConvertFromCalendarEvent(t *testing.T) {
	calendarEvent := calendar.NewEvent("Do homework",
		calendar.NewEventTime(2019, 10, 15, 20, 0),
		calendar.NewEventTime(2019, 10, 15, 22, 0),
	)

	event := ConvertFromCalendarEvent(calendarEvent)
	if event.Name != calendarEvent.Name() {
		t.Errorf("http.Event.Name not equals calendar.Event.Name(): `%s` != `%s`\n", event.Name, calendarEvent.Name())
	}

	if event.Start != "2019-10-15 20:00" {
		t.Errorf("http.Event.Start expected be %s, instead of %s\n", "2019-10-15 20:00", event.Start)
	}

	if event.End != "2019-10-15 22:00" {
		t.Errorf("http.Event.Start expected be %s, instead of %s\n", "2019-10-15 22:00", event.End)
	}

	if event.Id != 0 {
		t.Errorf("http.Event.Id expected be 0, instead of %d\n", event.Id)
	}
}

func TestConvertToCalendarEvent(t *testing.T) {
	calendarEvent := calendar.NewEvent("Do homework",
		calendar.NewEventTime(2019, 10, 15, 20, 0),
		calendar.NewEventTime(2019, 10, 15, 22, 0),
	)

	event := ConvertFromCalendarEvent(calendarEvent)

	resultCalendarEvent, err := event.ConvertToCalendarEvent()
	if err != nil {
		t.Errorf("must not be error while converting %s", err)
	}

	if calendarEvent != *resultCalendarEvent {
		t.Errorf("\nexpect calendar.Event:\n`%#v`\ngot calendar.Event:\n`%#v`\n",
			calendarEvent,
			*resultCalendarEvent)
	}
}

func TestJsonUnmarshal(t *testing.T)  {
	jsonData := `{
		"name": "Do homework",
		"start": "2019-11-15 20:00",
		"end": "2019-11-15 22:00"
	}`

	event, err := JsonUnmarshal([]byte(jsonData))
	if err != nil {
		t.Errorf("must not be error %s", err)
	}

	expectedEvent := Event{
		Name: "Do homework",
		Start: "2019-11-15 20:00",
		End: "2019-11-15 22:00",
	}

	if expectedEvent != *event {
		t.Errorf("expect event %+v, got event %+v\n",
			expectedEvent,
			*event)
	}
}

func TestJsonMarshal(t *testing.T) {
	event := &Event{
		Name: "Do homework",
		Start: "2019-11-17 23:00",
		End: "2019-11-18 08:00",
	}

	result, err := event.JsonMarshall()

	if err != nil {
		t.Errorf("must not be error %s", err)
	}

	expected := `{"name":"Do homework","start":"2019-11-17 23:00","end":"2019-11-18 08:00"}`
	resultStr := string(result)

	if expected != resultStr {
		t.Errorf("expect result %s, got result %s\n",
			fmt.Sprintf("%s", expected),
			fmt.Sprintf("%s", resultStr))
	}
}
