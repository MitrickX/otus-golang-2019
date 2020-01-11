package http

import (
	"encoding/json"
	"io/ioutil"
	"net/http/httptest"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/mitrickx/otus-golang-2019/30/calendar/internal/storage/memory"
)

func TestCreateEventOK(t *testing.T) {
	service := NewTestService()

	data := url.Values{}
	data.Set("name", "Do homework")
	data.Set("start", "2019-10-15 20:00")
	data.Set("end", "2019-10-15 22:00")
	data.Set("beforeMinutes", "10")
	body := strings.NewReader(data.Encode())

	req := httptest.NewRequest("POST", "http://test.com/create_event", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()

	service.CreateEvent(w, req)

	resp := w.Result()

	if resp.StatusCode != 200 {
		t.Errorf("must be status code 200 not %d", resp.StatusCode)
	}

	respBody, _ := ioutil.ReadAll(resp.Body)
	defer func() {
		_ = resp.Body.Close()
	}()

	okResp := &OkResponse{}
	err := json.Unmarshal(respBody, okResp)
	if err != nil {
		t.Errorf("failed on unmarshal json %s", err)
	}

	re := regexp.MustCompile(`created (\d+)`)
	match := re.FindStringSubmatch(okResp.Result)
	if match == nil {
		t.Errorf("unexpected OkResponse.Result value `%s`", okResp.Result)
		return
	}

	id, err := strconv.Atoi(match[1])
	if err != nil {
		t.Fatalf("strconv return error %s", err)
	}

	if id <= 0 {
		t.Errorf("unexpected event id %d, must be > 0", id)
	}

	if service.Calendar.getEventsTotalCount() != 1 {
		t.Errorf("unexpected count of events in entities, must be 1 instead of %d", service.Calendar.getEventsTotalCount())
		return
	}

	event, ok := service.Calendar.GetEvent(id)

	if !ok {
		t.Error("Expected event be present in calendar")
		return
	}

	expectedEvent := Event{
		Id:                 id,
		Name:               "Do homework",
		Start:              "2019-10-15 20:00",
		End:                "2019-10-15 22:00",
		IsNotifyingEnabled: true,
		BeforeMinutes:      10,
	}

	if *event != expectedEvent {
		t.Errorf("Expected\n`%+v`\ngot\n`%+v`", expectedEvent, *event)
	}
}

func TestCreateEventInvalidDate(t *testing.T) {
	service := NewTestService()

	data := url.Values{}
	data.Set("name", "Do homework")
	data.Set("start", "sdfasdf")
	data.Set("end", "2019-10-15 22:00")
	body := strings.NewReader(data.Encode())

	req := httptest.NewRequest("POST", "http://test.com/create_event", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()

	service.CreateEvent(w, req)

	resp := w.Result()

	if resp.StatusCode != 400 {
		t.Errorf("must be status code 400 not %d", resp.StatusCode)
	}

	respBody, _ := ioutil.ReadAll(resp.Body)
	defer func() {
		_ = resp.Body.Close()
	}()

	errResp := &ErrorResponse{}
	err := json.Unmarshal(respBody, errResp)
	if err != nil {
		t.Errorf("failed on unmarshal json %s", err)
	}

	if errResp.Error != DefaultErrorInvalidDatetime.Error() {
		t.Errorf("unexpected error `%s` instread of `%s`", errResp.Error, DefaultErrorInvalidDatetime.Error())
	}

	if service.Calendar.getEventsTotalCount() != 0 {
		t.Errorf("unexpected count of events in entities, must be 0 instead of %d", service.Calendar.getEventsTotalCount())
	}

}

func TestUpdateEventOK(t *testing.T) {
	service := NewTestService()

	event1 := &Event{
		Name:  "Do homework",
		Start: "2019-10-15 20:00",
		End:   "2019-10-15 22:00",
	}

	id := addEvent(t, &service.Calendar, event1, 1)
	if id <= 0 {
		return
	}

	event2 := &Event{
		Name:               "Watch movie",
		Start:              "2019-10-15 22:00",
		End:                "2019-10-16 01:00",
		IsNotifyingEnabled: true,
		BeforeMinutes:      5,
	}

	data := url.Values{}
	data.Set("id", strconv.Itoa(id))
	data.Set("name", event2.Name)
	data.Set("start", event2.Start)
	data.Set("end", event2.End)
	data.Set("beforeMinutes", strconv.Itoa(event2.BeforeMinutes))

	body := strings.NewReader(data.Encode())

	req := httptest.NewRequest("POST", "http://test.com/create_event", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()

	service.UpdateEvent(w, req)

	resp := w.Result()

	if resp.StatusCode != 200 {
		t.Errorf("must be status code 200 not %d", resp.StatusCode)
	}

	respBody, _ := ioutil.ReadAll(resp.Body)
	defer func() {
		_ = resp.Body.Close()
	}()

	okResp := &OkResponse{}
	err := json.Unmarshal(respBody, okResp)
	if err != nil {
		t.Errorf("failed on unmarshal json %s", err)
	}

	if okResp.Result != "updated" {
		t.Errorf("unexpected OkResponse.Result value `%s`", okResp.Result)
		return
	}

	event, found := service.GetEvent(id)
	if !found {
		t.Errorf("event with id = %d not found on entities service", id)
		return
	}

	expectedEvent := *event2
	expectedEvent.Id = id
	if expectedEvent != *event {
		t.Errorf("\nevent info not updated\nexpected be:\n%#v\ngot:\n%#v\n", expectedEvent, event)
	}

}

func TestUpdateEventInvalidId(t *testing.T) {
	service := NewTestService()

	event1 := &Event{
		Name:  "Do homework",
		Start: "2019-10-15 20:00",
		End:   "2019-10-15 22:00",
	}

	id := addEvent(t, &service.Calendar, event1, 1)
	if id <= 0 {
		return
	}

	event2 := &Event{
		Name:  "Watch movie",
		Start: "2019-10-15 22:00",
		End:   "2019-10-16 01:00",
	}

	data := url.Values{}
	data.Set("id", "sdfsdf")
	data.Set("name", event2.Name)
	data.Set("start", event2.Start)
	data.Set("end", event2.End)
	body := strings.NewReader(data.Encode())

	req := httptest.NewRequest("POST", "http://test.com/create_event", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()

	service.UpdateEvent(w, req)

	resp := w.Result()

	if resp.StatusCode != 400 {
		t.Errorf("must be status code 400 not %d", resp.StatusCode)
	}

	respBody, _ := ioutil.ReadAll(resp.Body)
	defer func() {
		_ = resp.Body.Close()
	}()

	errResp := &ErrorResponse{}
	err := json.Unmarshal(respBody, errResp)
	if err != nil {
		t.Errorf("failed on unmarshal json %s", err)
	}

	expectedErr := "invalid id parameter, must be int greater than 0"
	if errResp.Error != expectedErr {
		t.Errorf("unexpected error %s, must be %s", errResp.Error, expectedErr)
	}

	event, found := service.GetEvent(id)
	if !found {
		t.Errorf("event with id = %d not found on entities service", id)
		return
	}

	if event.Name != event1.Name || event.Start != event1.Start || event.End != event1.End {
		t.Errorf("\nevent info must be not updated\nexpected be:\n%+v\ngot:\n%+v\n", event1, event)
	}

}

func TestUpdateEventInvalidDate(t *testing.T) {
	service := NewTestService()

	event1 := &Event{
		Name:  "Do homework",
		Start: "2019-10-15 20:00",
		End:   "2019-10-15 22:00",
	}

	id := addEvent(t, &service.Calendar, event1, 1)
	if id <= 0 {
		return
	}

	event2 := &Event{
		Name:  "Watch movie",
		Start: "2019-10-15 22:00",
		End:   "2019-10-16 01:00",
	}

	data := url.Values{}
	data.Set("id", strconv.Itoa(id))
	data.Set("name", event2.Name)
	data.Set("start", "dafsdaf")
	data.Set("end", event2.End)
	body := strings.NewReader(data.Encode())

	req := httptest.NewRequest("POST", "http://test.com/create_event", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()

	service.UpdateEvent(w, req)

	resp := w.Result()

	if resp.StatusCode != 400 {
		t.Errorf("must be status code 400 not %d", resp.StatusCode)
	}

	respBody, _ := ioutil.ReadAll(resp.Body)
	defer func() {
		_ = resp.Body.Close()
	}()

	errResp := &ErrorResponse{}
	err := json.Unmarshal(respBody, errResp)
	if err != nil {
		t.Errorf("failed on unmarshal json %s", err)
	}

	if errResp.Error != DefaultErrorInvalidDatetime.Error() {
		t.Errorf("unexpected error `%s` instread of `%s`", errResp.Error, DefaultErrorInvalidDatetime.Error())
	}

	event, found := service.GetEvent(id)
	if !found {
		t.Errorf("event with id = %d not found on entities service", id)
		return
	}

	if event.Name != event1.Name || event.Start != event1.Start || event.End != event1.End {
		t.Errorf("\nevent info must be not updated\nexpected be:\n%+v\ngot:\n%+v\n", event1, event)
	}

}

func TestUpdateEventNotFound(t *testing.T) {
	service := NewTestService()

	event1 := &Event{
		Name:  "Do homework",
		Start: "2019-10-15 20:00",
		End:   "2019-10-15 22:00",
	}

	id := addEvent(t, &service.Calendar, event1, 1)
	if id <= 0 {
		return
	}

	event2 := &Event{
		Name:  "Watch movie",
		Start: "2019-10-15 22:00",
		End:   "2019-10-16 01:00",
	}

	data := url.Values{}
	data.Set("id", "100")
	data.Set("name", event2.Name)
	data.Set("start", event2.Start)
	data.Set("end", event2.End)
	body := strings.NewReader(data.Encode())

	req := httptest.NewRequest("POST", "http://test.com/create_event", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()

	service.UpdateEvent(w, req)

	resp := w.Result()

	if resp.StatusCode != 200 {
		t.Errorf("must be status code 200 not %d", resp.StatusCode)
	}

	respBody, _ := ioutil.ReadAll(resp.Body)
	defer func() {
		_ = resp.Body.Close()
	}()

	errResp := &ErrorResponse{}
	err := json.Unmarshal(respBody, errResp)
	if err != nil {
		t.Errorf("failed on unmarshal json %s", err)
	}

	if errResp.Error == "" {
		t.Errorf("unexpected empty error")
	}

}

func TestDeleteEventOK(t *testing.T) {
	service := NewTestService()

	event := &Event{
		Name:  "Do homework",
		Start: "2019-10-15 20:00",
		End:   "2019-10-15 22:00",
	}

	id := addEvent(t, &service.Calendar, event, 1)
	if id <= 0 {
		return
	}

	data := url.Values{}
	data.Set("id", strconv.Itoa(id))
	body := strings.NewReader(data.Encode())

	req := httptest.NewRequest("POST", "http://test.com/create_event", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()

	service.DeleteEvent(w, req)

	resp := w.Result()

	if resp.StatusCode != 200 {
		t.Errorf("must be status code 200 not %d", resp.StatusCode)
	}

	respBody, _ := ioutil.ReadAll(resp.Body)
	defer func() {
		_ = resp.Body.Close()
	}()

	okResp := &OkResponse{}
	err := json.Unmarshal(respBody, okResp)
	if err != nil {
		t.Errorf("failed on unmarshal json %s", err)
	}

	if okResp.Result != "deleted" {
		t.Errorf("unexpected OkResponse.Result value `%s`", okResp.Result)
		return
	}

	_, found := service.GetEvent(id)
	if found {
		t.Errorf("event with id = %d have not be deleted on entities service", id)
		return
	}

	if service.Calendar.getEventsTotalCount() != 0 {
		t.Errorf("unexpected count of events in entities, must be 0 instead of %d", service.Calendar.getEventsTotalCount())
	}

}

func TestDeleteEventInvalidId(t *testing.T) {
	service := NewTestService()

	event := &Event{
		Name:  "Do homework",
		Start: "2019-10-15 20:00",
		End:   "2019-10-15 22:00",
	}

	id := addEvent(t, &service.Calendar, event, 1)
	if id <= 0 {
		return
	}

	data := url.Values{}
	data.Set("id", "sdsdf")
	body := strings.NewReader(data.Encode())

	req := httptest.NewRequest("POST", "http://test.com/create_event", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()

	service.DeleteEvent(w, req)

	resp := w.Result()

	if resp.StatusCode != 400 {
		t.Errorf("must be status code 400 not %d", resp.StatusCode)
	}

	respBody, _ := ioutil.ReadAll(resp.Body)
	defer func() {
		_ = resp.Body.Close()
	}()

	errResp := &ErrorResponse{}
	err := json.Unmarshal(respBody, errResp)
	if err != nil {
		t.Errorf("failed on unmarshal json %s", err)
	}

	expectedErr := "invalid id parameter, must be int greater than 0"
	if errResp.Error != expectedErr {
		t.Errorf("unexpected error %s, must be %s", errResp.Error, expectedErr)
	}
}

func TestDeleteEventNotFound(t *testing.T) {
	service := NewTestService()

	event := &Event{
		Name:  "Do homework",
		Start: "2019-10-15 20:00",
		End:   "2019-10-15 22:00",
	}

	id := addEvent(t, &service.Calendar, event, 1)
	if id <= 0 {
		return
	}

	data := url.Values{}
	data.Set("id", "100")
	body := strings.NewReader(data.Encode())

	req := httptest.NewRequest("POST", "http://test.com/create_event", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()

	service.DeleteEvent(w, req)

	resp := w.Result()

	if resp.StatusCode != 200 {
		t.Errorf("must be status code 200 not %d", resp.StatusCode)
	}

	respBody, _ := ioutil.ReadAll(resp.Body)
	defer func() {
		_ = resp.Body.Close()
	}()

	errResp := &ErrorResponse{}
	err := json.Unmarshal(respBody, errResp)
	if err != nil {
		t.Errorf("failed on unmarshal json %s", err)
	}

	if errResp.Error == "" {
		t.Errorf("unexpected empty error")
	}
}

func TestGetEventsForDay(t *testing.T) {
	service := NewTestService()

	addFixedListOfEvents(t, &service.Calendar)

	req := httptest.NewRequest("GET", "http://test.com/events_for_day", nil)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()

	now := time.Date(2019, 11, 21, 8, 0, 0, 0, time.UTC)
	service.getEventsForDay(now, w, req)

	resp := w.Result()

	if resp.StatusCode != 200 {
		t.Errorf("must be status code 200 not %d", resp.StatusCode)
	}

	respBody, _ := ioutil.ReadAll(resp.Body)
	defer func() {
		_ = resp.Body.Close()
	}()

	eventListResp := &EventListResponse{}
	err := json.Unmarshal(respBody, eventListResp)
	if err != nil {
		t.Errorf("failed on unmarshal json %s", err)
	}

	if len(eventListResp.Result) != 1 {
		t.Errorf("event list must has one event instead of %d", len(eventListResp.Result))
	}

	event := *eventListResp.Result[0]
	if event.Name != "Thursday" {
		t.Error("Must be Thursday event")
	}
}

func TestGetEventsForWeek(t *testing.T) {
	service := NewTestService()

	addFixedListOfEvents(t, &service.Calendar)

	req := httptest.NewRequest("GET", "http://test.com/events_for_day", nil)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()

	now := time.Date(2019, 11, 21, 8, 0, 0, 0, time.UTC)
	service.getEventsForWeek(now, w, req)

	resp := w.Result()

	if resp.StatusCode != 200 {
		t.Errorf("must be status code 200 not %d", resp.StatusCode)
	}

	respBody, _ := ioutil.ReadAll(resp.Body)
	defer func() {
		_ = resp.Body.Close()
	}()

	eventListResp := &EventListResponse{}
	err := json.Unmarshal(respBody, eventListResp)
	if err != nil {
		t.Errorf("failed on unmarshal json %s", err)
	}

	if len(eventListResp.Result) != 7 {
		t.Errorf("event list must has 7 events instead of %d", len(eventListResp.Result))
	}
}

func TestGetEventsForMonth(t *testing.T) {
	service := NewTestService()

	addFixedListOfEvents(t, &service.Calendar)

	addEvent(t, &service.Calendar, &Event{
		Name:  "First day",
		Start: "2019-11-01 08:00",
		End:   "2019-11-01 10:00",
	}, 8)

	addEvent(t, &service.Calendar, &Event{
		Name:  "Last day",
		Start: "2019-11-30 08:00",
		End:   "2019-11-30 10:00",
	}, 9)

	req := httptest.NewRequest("POST", "http://test.com/events_for_day", nil)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()

	now := time.Date(2019, 11, 21, 8, 0, 0, 0, time.UTC)
	service.getEventsForMonth(now, w, req)

	resp := w.Result()

	if resp.StatusCode != 200 {
		t.Errorf("must be status code 200 not %d", resp.StatusCode)
	}

	respBody, _ := ioutil.ReadAll(resp.Body)
	defer func() {
		_ = resp.Body.Close()
	}()

	eventListResp := &EventListResponse{}
	err := json.Unmarshal(respBody, eventListResp)
	if err != nil {
		t.Errorf("failed on unmarshal json %s", err)
	}

	if len(eventListResp.Result) != 9 {
		t.Errorf("event list must has 9 events instead of %d", len(eventListResp.Result))
	}
}

func NewTestService() *Service {
	storage := memory.NewStorage()
	service, _ := NewService("", storage, nil, nil)
	return service
}
