package sql

import (
	"fmt"
	"github.com/mitrickx/otus-golang-2019/25/calendar/internal/domain/entities"
	"github.com/spf13/viper"
	"reflect"
	"testing"
)

// These tests test operation on actual DB
// So we need DB connection settings to test DB
// These tests also must be run from directory of current file, separately from other tests: go test -v -race ./

const cfgDefaultFilePath = "../../../configs/test.yaml"

var config testsConfig

type testsConfig struct {
	skip     bool
	dbConfig *Config
}

func NewTestsConfig() testsConfig {

	config := testsConfig{}
	v := viper.New()
	v.SetConfigFile(cfgDefaultFilePath)
	err := v.ReadInConfig()
	if err != nil {
		config.skip = true
		return config
	}

	db := v.GetStringMapString("db")
	if db == nil {
		config.skip = true
		return config
	}

	dbConfig, err := NewConfig(db)
	if err != nil {
		config.skip = true
		return config
	}

	config.dbConfig = dbConfig

	return config
}

func init() {
	config = NewTestsConfig()
}

// Test adding new event in entities
func TestAddEvent(t *testing.T) {

	if config.skip {
		t.SkipNow()
	}

	calendar := NewTestStorage(t, &config)

	if getCalendarCount(calendar) > 0 {
		t.Error("empty entities must not has events")
	}

	event1 := entities.NewEvent("Do homework",
		entities.NewEventTime(2019, 10, 15, 20, 0),
		entities.NewEventTime(2019, 10, 15, 22, 0),
	)

	id, _ := calendar.AddEvent(event1)

	if id <= 0 {
		t.Errorf("id %d of new added event must be > 0", id)
	}

	if getCalendarCount(calendar) != 1 {
		t.Errorf("entities must has 1 event instead of %d\n", getCalendarCount(calendar))
	}

	event2 := entities.NewEvent("Watch movie",
		entities.NewEventTime(2019, 10, 15, 22, 0),
		entities.NewEventTime(2019, 10, 16, 1, 0),
	)

	calendar.AddEvent(event2)

	if getCalendarCount(calendar) != 2 {
		t.Errorf("entities must has 2 event instead of %d\n", getCalendarCount(calendar))
	}
}

// Test get events from entities
func TestGetEvent(t *testing.T) {

	if config.skip {
		t.SkipNow()
	}

	calendar := NewTestStorage(t, &config)

	event1 := entities.NewEvent("Do homework",
		entities.NewEventTime(2019, 10, 15, 20, 0),
		entities.NewEventTime(2019, 10, 15, 22, 0),
	)

	id, _ := calendar.AddEvent(event1)

	event, err := calendar.GetEvent(id)
	if err != nil {
		t.Errorf("Get Event must be ok, instread of error: %s\n", err)
		return
	}

	expectedName := "Do homework"
	if event.Name() != expectedName {
		t.Errorf("event.Name() must be `%s` instread of `%s`\n", expectedName, event.Name())
		fmt.Println(event)
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

	if config.skip {
		t.SkipNow()
	}

	calendar := NewTestStorage(t, &config)

	event1 := entities.NewEvent("Do homework",
		entities.NewEventTime(2019, 10, 15, 20, 0),
		entities.NewEventTime(2019, 10, 15, 22, 0),
	)

	id, _ := calendar.AddEvent(event1)

	event2 := entities.NewEvent("Watch movie",
		entities.NewEventTime(2019, 10, 15, 22, 0),
		entities.NewEventTime(2019, 10, 16, 1, 0),
	)

	err := calendar.UpdateEvent(id, event2)

	if err != nil {
		t.Errorf("update err `%s` must not be happened\n", err)
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

	if config.skip {
		t.SkipNow()
	}

	calendar := NewTestStorage(t, &config)

	event1 := entities.NewEvent("Do homework",
		entities.NewEventTime(2019, 10, 15, 20, 0),
		entities.NewEventTime(2019, 10, 15, 22, 0),
	)

	id1, _ := calendar.AddEvent(event1)

	event2 := entities.NewEvent("Watch movie",
		entities.NewEventTime(2019, 10, 15, 22, 0),
		entities.NewEventTime(2019, 10, 16, 1, 0),
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

	err = calendar.DeleteEvent(id2)

	if getCalendarCount(calendar) != 0 {
		t.Error("delete actually not happened, after 2 delete calender must be empty")
	}
}

func TestGetEvents(t *testing.T) {

	if config.skip {
		t.SkipNow()
	}

	calendar := NewTestStorage(t, &config)

	calendar.AddEvent(entities.NewEvent("Monday",
		entities.NewEventTime(2019, 11, 18, 8, 0),
		entities.NewEventTime(2019, 11, 18, 10, 0),
	))

	calendar.AddEvent(entities.NewEvent("Tuesday",
		entities.NewEventTime(2019, 11, 19, 8, 0),
		entities.NewEventTime(2019, 11, 19, 10, 0),
	))

	calendar.AddEvent(entities.NewEvent("Wednesday",
		entities.NewEventTime(2019, 11, 20, 8, 0),
		entities.NewEventTime(2019, 11, 20, 10, 0),
	))

	calendar.AddEvent(entities.NewEvent("Thursday",
		entities.NewEventTime(2019, 11, 21, 8, 0),
		entities.NewEventTime(2019, 11, 21, 10, 0),
	))

	calendar.AddEvent(entities.NewEvent("Friday",
		entities.NewEventTime(2019, 11, 22, 8, 0),
		entities.NewEventTime(2019, 11, 22, 10, 0),
	))

	calendar.AddEvent(entities.NewEvent("Saturday",
		entities.NewEventTime(2019, 11, 23, 8, 0),
		entities.NewEventTime(2019, 11, 23, 10, 0),
	))

	calendar.AddEvent(entities.NewEvent("Sunday",
		entities.NewEventTime(2019, 11, 24, 8, 0),
		entities.NewEventTime(2019, 11, 24, 10, 0),
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

	eventTime := entities.NewEventTime(2019, 11, 18, 8, 0)
	eventList, _ := calendar.GetEventsByPeriod(nil, &eventTime)

	if len(eventList) != 1 {
		t.Errorf("Must be returned one event")
	}

	if eventList[0].Name() != "Monday" {
		t.Errorf("Must be returned one event `Monday`")
	}

	eventTime = entities.NewEventTime(2019, 11, 24, 8, 0)
	eventList, _ = calendar.GetEventsByPeriod(&eventTime, nil)

	if len(eventList) != 1 {
		t.Errorf("Must be returned one event")
	}

	if eventList[0].Name() != "Sunday" {
		t.Errorf("Must be returned one event `Sunday`")
	}

	startEventTime := entities.NewEventTime(2019, 11, 20, 8, 0)
	endEventTime := entities.NewEventTime(2019, 11, 22, 8, 0)
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

	startEventTime = entities.NewEventTime(2019, 11, 20, 8, 1)
	endEventTime = entities.NewEventTime(2019, 11, 20, 9, 59)
	eventList, _ = calendar.GetEventsByPeriod(&startEventTime, &endEventTime)

	if len(eventList) != 0 {
		t.Errorf("Must be returned 0 events")
	}
}

func NewTestStorage(t *testing.T, config *testsConfig) *Storage {
	storage, err := NewStorage(*config.dbConfig)
	if err != nil {
		t.Fatalf("fail on create storage instalce %s", err)
	}
	_ = storage.ClearAll()
	return storage
}

func getCalendarCount(storage *Storage) int {
	cnt, _ := storage.Count()
	return cnt
}
