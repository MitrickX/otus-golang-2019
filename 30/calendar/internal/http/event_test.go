package http

import (
	"github.com/mitrickx/otus-golang-2019/30/calendar/internal/domain/entities"
	"testing"
	"time"
)

func TestNewEvent(t *testing.T) {
	_, err := NewEvent("Do homework", "2019-10-15 20:00", "2019-10-15 22:00", false, 0)
	if err != nil {
		t.Errorf("http.Event construction must be ok, not failed because of `%s`", err)
	}

	_, err = NewEvent("Do homework", "17897897", "2019-10-15 22:00", false, 0)

	if err == nil {
		t.Error("http.Event construction must return ErrorInvalidDatetime, not nil")
	} else if _, ok := err.(*ErrorInvalidDatetime); !ok {
		t.Errorf("http.Event construction must return ErrorInvalidDatetime, not `%+v`", err)
	}

	_, err = NewEvent("Do homework", "2019-10-15 20:00", "02 Jan 06 15:04 MST", false, 0)

	if err == nil {
		t.Error("http.Event construction must return ErrorInvalidDatetime, not nil")
	} else if _, ok := err.(*ErrorInvalidDatetime); !ok {
		t.Errorf("http.Event construction must return ErrorInvalidDatetime, not `%+v`", err)
	}

}

func TestConvertFromCalendarEvent(t *testing.T) {
	calendarEvent := entities.NewDetailedEvent("Do homework",
		entities.NewDateTime(2019, 10, 15, 20, 0),
		entities.NewDateTime(2019, 10, 15, 22, 0),
		true,
		10,
		false,
		time.Time{},
	)

	event := ConvertFromCalendarEvent(calendarEvent)
	if event.Name != calendarEvent.Name() {
		t.Errorf("http.Event.Name not equals entities.Event.Name(): `%s` != `%s`\n", event.Name, calendarEvent.Name())
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

	if event.IsNotifyingEnabled != true {
		t.Errorf("http.Event.IsNotifiedEnable expected be true")
	}

	if event.BeforeMinutes != 10 {
		t.Errorf("http.Event.BeforeMinutes expected to be 10 instread of %d\n", event.BeforeMinutes)
	}
}

func TestConvertToCalendarEvent(t *testing.T) {
	calendarEvent := entities.NewDetailedEvent("Do homework",
		entities.NewDateTime(2019, 10, 15, 20, 0),
		entities.NewDateTime(2019, 10, 15, 22, 0),
		true,
		11,
		false,
		time.Time{},
	)

	event := ConvertFromCalendarEvent(calendarEvent)

	resultCalendarEvent, err := event.ConvertToCalendarEvent()
	if err != nil {
		t.Errorf("must not be error while converting %s", err)
	}

	if calendarEvent != *resultCalendarEvent {
		t.Errorf("\nexpect entities.Event:\n`%#v`\ngot entities.Event:\n`%#v`\n",
			calendarEvent,
			*resultCalendarEvent)
	}
}

func TestJsonUnmarshal1(t *testing.T) {
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
		Name:  "Do homework",
		Start: "2019-11-15 20:00",
		End:   "2019-11-15 22:00",
	}

	if expectedEvent != *event {
		t.Errorf("expect event %+v, got event %+v\n",
			expectedEvent,
			*event)
	}
}

func TestJsonUnmarshal2(t *testing.T) {
	jsonData := `{
		"name": "Do homework",
		"start": "2019-11-15 20:00",
		"end": "2019-11-15 22:00",
		"isNotifyingEnabled": true,
		"beforeMinutes": 5
	}`

	event, err := JsonUnmarshal([]byte(jsonData))
	if err != nil {
		t.Errorf("must not be error %s", err)
	}

	expectedEvent := Event{
		Name:               "Do homework",
		Start:              "2019-11-15 20:00",
		End:                "2019-11-15 22:00",
		IsNotifyingEnabled: true,
		BeforeMinutes:      5,
	}

	if expectedEvent != *event {
		t.Errorf("expect event %+v, got event %+v\n",
			expectedEvent,
			*event)
	}
}

func TestJsonMarshal1(t *testing.T) {
	event := &Event{
		Name:  "Do homework",
		Start: "2019-11-17 23:00",
		End:   "2019-11-18 08:00",
	}

	result, err := event.JsonMarshall()

	if err != nil {
		t.Errorf("must not be error %s", err)
	}

	expected := `{"name":"Do homework","start":"2019-11-17 23:00","end":"2019-11-18 08:00"}`
	resultStr := string(result)

	if expected != resultStr {
		t.Errorf("expect result %s, got result %s\n", expected, resultStr)
	}
}

func TestJsonMarshal2(t *testing.T) {
	event := &Event{
		Name:               "Do homework",
		Start:              "2019-11-17 23:00",
		End:                "2019-11-18 08:00",
		IsNotifyingEnabled: true,
		BeforeMinutes:      7,
	}

	result, err := event.JsonMarshall()

	if err != nil {
		t.Errorf("must not be error %s", err)
	}

	expected := `{"name":"Do homework","start":"2019-11-17 23:00","end":"2019-11-18 08:00","isNotifyingEnabled":true,"beforeMinutes":7}`
	resultStr := string(result)

	if expected != resultStr {
		t.Errorf("expect result %s, got result %s\n", expected, resultStr)
	}
}

func TestGetDayPeriod(t *testing.T) {
	now := time.Date(2019, 11, 17, 14, 33, 12, 0, time.UTC)
	start, end := GetDayPeriod(now)
	expectedStart := "2019-11-17 00:00"
	expectedEnd := "2019-11-17 23:59"
	if start != expectedStart {
		t.Errorf("start must be %s insteadof %s", start, expectedStart)
	}
	if end != expectedEnd {
		t.Errorf("end must be %s insteadof %s", end, expectedEnd)
	}
}

func TestGetWeekPeriod1(t *testing.T) {
	now := time.Date(2019, 11, 17, 14, 33, 12, 0, time.UTC)
	start, end := GetWeekPeriod(now)
	expectedStart := "2019-11-11 00:00"
	expectedEnd := "2019-11-17 23:59"
	if start != expectedStart {
		t.Errorf("start must be %s insteadof %s", start, expectedStart)
	}
	if end != expectedEnd {
		t.Errorf("end must be %s insteadof %s", end, expectedEnd)
	}
}

func TestGetWeekPeriod2(t *testing.T) {
	now := time.Date(2019, 11, 13, 14, 33, 12, 0, time.UTC)
	start, end := GetWeekPeriod(now)
	expectedStart := "2019-11-11 00:00"
	expectedEnd := "2019-11-17 23:59"
	if start != expectedStart {
		t.Errorf("start must be %s insteadof %s", start, expectedStart)
	}
	if end != expectedEnd {
		t.Errorf("end must be %s insteadof %s", end, expectedEnd)
	}
}

func TestGetWeekPeriod3(t *testing.T) {
	now := time.Date(2019, 11, 11, 14, 33, 12, 0, time.UTC)
	start, end := GetWeekPeriod(now)
	expectedStart := "2019-11-11 00:00"
	expectedEnd := "2019-11-17 23:59"
	if start != expectedStart {
		t.Errorf("start must be %s insteadof %s", start, expectedStart)
	}
	if end != expectedEnd {
		t.Errorf("end must be %s insteadof %s", end, expectedEnd)
	}
}

func TestGetMonthPeriod1(t *testing.T) {
	now := time.Date(2019, 11, 11, 14, 33, 12, 0, time.UTC)
	start, end := GetMonthPeriod(now)
	expectedStart := "2019-11-01 00:00"
	expectedEnd := "2019-11-30 23:59"
	if start != expectedStart {
		t.Errorf("start must be %s insteadof %s", start, expectedStart)
	}
	if end != expectedEnd {
		t.Errorf("end must be %s insteadof %s", end, expectedEnd)
	}
}

func TestGetMonthPeriod2(t *testing.T) {
	now := time.Date(2019, 12, 11, 14, 33, 12, 0, time.UTC)
	start, end := GetMonthPeriod(now)
	expectedStart := "2019-12-01 00:00"
	expectedEnd := "2019-12-31 23:59"
	if start != expectedStart {
		t.Errorf("start must be %s insteadof %s", start, expectedStart)
	}
	if end != expectedEnd {
		t.Errorf("end must be %s insteadof %s", end, expectedEnd)
	}
}

func TestGetMonthPeriod3(t *testing.T) {
	now := time.Date(2020, 2, 11, 14, 33, 12, 0, time.UTC)
	start, end := GetMonthPeriod(now)
	expectedStart := "2020-02-01 00:00"
	expectedEnd := "2020-02-29 23:59"
	if start != expectedStart {
		t.Errorf("start must be %s insteadof %s", start, expectedStart)
	}
	if end != expectedEnd {
		t.Errorf("end must be %s insteadof %s", end, expectedEnd)
	}
}

func TestGetMonthPeriod4(t *testing.T) {
	now := time.Date(2019, 2, 11, 14, 33, 12, 0, time.UTC)
	start, end := GetMonthPeriod(now)
	expectedStart := "2019-02-01 00:00"
	expectedEnd := "2019-02-28 23:59"
	if start != expectedStart {
		t.Errorf("start must be %s insteadof %s", start, expectedStart)
	}
	if end != expectedEnd {
		t.Errorf("end must be %s insteadof %s", end, expectedEnd)
	}
}