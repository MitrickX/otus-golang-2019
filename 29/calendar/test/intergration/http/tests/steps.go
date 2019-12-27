package tests

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	httpService "github.com/mitrickx/otus-golang-2019/29/calendar/internal/http"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	"github.com/mitrickx/otus-golang-2019/29/calendar/internal/domain/entities"
)

type featureTest struct {
	r              *http.Response // http response
	responseBody   string
	subMatchResult []string // result of regexp sub-match searching
	eventId        int      // id of event to deal with in next step(s)
	eventIds       []int    // id of events to deal with in next step(s)
	vars           map[string]interface{}
}

func newFeatureTest() *featureTest {
	return new(featureTest)
}

func (t *featureTest) iSendRequestToWithParams(method, addr, contentType string, data *gherkin.DocString) error {
	var err error
	replacer := strings.NewReplacer("\n", "", "\t", "")
	query := replacer.Replace(data.Content)

	log.Printf("Send %s request with data `%s` of content type `%s` to addr `%s`", method, query, contentType, addr)

	req, err := http.NewRequest(method, addr, bytes.NewReader([]byte(query)))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", contentType)

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return fmt.Errorf("Do request error %s", err)
	}

	defer func() {
		_ = res.Body.Close()
	}()

	t.r = res

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	t.responseBody = string(body)

	return nil
}

func (t *featureTest) iSendRequestTo(method, addr string) error {
	data := &gherkin.DocString{
		Content: "",
	}
	return t.iSendRequestToWithParams(method, addr, "application/x-www-form-urlencoded", data)
}

func (t *featureTest) theResponseCodeShouldBe(code int) error {
	return assertStatusCode(t.r, code)
}

func (t *featureTest) theResponseContentTypeShouldBe(contentType string) error {
	return assertContentType(t.r, contentType)
}

func (t *featureTest) theResponseJsonShouldHasFieldWithValueMatch(field, pattern string) error {
	jsonResponse, err := jsonUnmarshalStringToStringMap(t.responseBody)
	if err != nil {
		return err
	}
	val, ok := jsonResponse[field]
	if !ok {
		return fmt.Errorf("expected that response json `%+v` should has field `%s`", jsonResponse, field)
	}
	re, err := regexp.Compile(pattern)
	if err != nil {
		return fmt.Errorf("compile regexp patter `%s` failed: %s", pattern, err)
	}

	m := re.Match([]byte(val))
	if !m {
		return fmt.Errorf("expected that value `%s` would match by pattern `%s`", val, pattern)
	}

	sb := re.FindStringSubmatch(val)
	t.subMatchResult = sb

	return nil
}

func (t *featureTest) extractedNumberIsEventId() error {
	if len(t.subMatchResult) < 2 {
		return fmt.Errorf("expected submatch search on previous step grap something")
	}
	eventId, err := strconv.Atoi(t.subMatchResult[1])
	if err != nil {
		return fmt.Errorf("error when convert captured group `%s` to real number `%s`", t.subMatchResult[1], err)
	}
	t.eventId = eventId
	if t.eventId <= 0 {
		return fmt.Errorf("expect event id %d is greater than 0", t.eventId)
	}
	return nil
}

func (t *featureTest) theRecordShouldMatch(eventsData *gherkin.DataTable) error {
	events, err := convertGherkinTableToEvents(eventsData)
	if err != nil {
		return err
	}
	if len(events) < 1 {
		return errors.New("event record not described in feature")
	}
	event := events[0]
	event = entities.WithId(event, t.eventId)

	cfg := GetTestConfig()

	dbEvent, err := cfg.DbStorage.GetEvent(event.Id())
	if err != nil {
		return fmt.Errorf("error when get event by id %d from db `%s`", event.Id(), err)
	}

	if dbEvent != event {
		return fmt.Errorf("expected event\n%#v\ninstread of event\n%#v", event, dbEvent)
	}
	return nil
}

func (t *featureTest) existingRecords(eventsData *gherkin.DataTable) error {

	events, err := convertGherkinTableToEvents(eventsData)
	if err != nil {
		return err
	}

	config := GetTestConfig()
	for _, event := range events {
		eventId, err := config.DbStorage.InsertEvent(event)
		if err != nil {
			return fmt.Errorf("preparation db, add event into is failed %s", err)
		}
		t.eventIds = append(t.eventIds, eventId)
	}

	return nil
}

func (t *featureTest) theRecordsShouldMatch(eventsData *gherkin.DataTable) error {
	expectedEvents, err := convertGherkinTableToEvents(eventsData)
	if err != nil {
		return err
	}

	config := GetTestConfig()

	var errList []string

	for index, event := range expectedEvents {
		dbEventId := t.eventIds[index]
		expectedEvent := entities.WithId(event, dbEventId)
		dbEvent, err := config.DbStorage.GetEvent(dbEventId)
		if err != nil {
			return fmt.Errorf("error when get event by id %d from db `%s`", dbEventId, err)
		}
		if expectedEvent != dbEvent {
			err := fmt.Sprintf("expected event\n%#v\ninstread of event\n%#v", expectedEvent, dbEvent)
			errList = append(errList, err)
		}
	}

	if len(errList) > 0 {
		errText := strings.Join(errList, "\n")
		return errors.New(errText)
	}

	return nil

}

func (t *featureTest) cleanDB() error {
	config := GetTestConfig()
	err := config.DbStorage.ClearAll()
	if err != nil {
		return fmt.Errorf("cleaning error return error %s", err)
	}
	return nil
}

func (t *featureTest) theDBShouldBeClean() error {
	config := GetTestConfig()
	count, err := config.DbStorage.Count()
	if err != nil {
		return fmt.Errorf("counting all items cause error `%s`", err)
	}
	if count > 0 {
		return fmt.Errorf("db not clean (number of events is %d)", count)
	}
	return nil
}

func (t *featureTest) theResponseJsonShouldMatch(data *gherkin.DocString) error {
	replacer := strings.NewReplacer("\n", "", "\t", "")
	expectedData := replacer.Replace(data.Content)
	strings.Contains(expectedData, "")
	expectedDataJson, err := jsonUnmarshalToMap(expectedData)
	if err != nil {
		return err
	}
	responseDataJson, err := jsonUnmarshalToMap(t.responseBody)
	if err != nil {
		return err
	}
	if !reflect.DeepEqual(expectedDataJson, responseDataJson) {
		return fmt.Errorf("expected:\n`%#v`\ninstead of:\n`%#v`", expectedDataJson, responseDataJson)
	}
	return nil
}

func (t *featureTest) theResponseJsonIsEventListResponseFilledWithEventsOfIds(ids string) error {
	eventListResponse, err := jsonUnmarshalEventListResponse(t.responseBody)
	if err != nil {
		return err
	}

	eventIds, err := splitStrIntoIntSlice(ids)
	if err != nil {
		return err
	}

	if len(eventIds) <= 0 {
		return fmt.Errorf("unexpected empty event ids list after parse ids `%s`", ids)
	}

	config := GetTestConfig()

	expectedEventListResponse := httpService.EventListResponse{}
	for _, eventId := range eventIds {
		event, err := config.DbStorage.GetEvent(eventId)
		if err != nil {
			return fmt.Errorf("error when get event by id %d from db `%s`", eventId, err)
		}
		httpServiceEvent := httpService.ConvertFromCalendarEvent(event)
		expectedEventListResponse.Result = append(expectedEventListResponse.Result, httpServiceEvent)
	}

	if !reflect.DeepEqual(expectedEventListResponse, eventListResponse) {
		return fmt.Errorf("expected reponse\n%s\n instead of\n%s\n",
			jsonMarshalEventListResponse(expectedEventListResponse),
			jsonMarshalEventListResponse(eventListResponse),
		)
	}

	return nil
}

func (t *featureTest) iAfterWaitShouldReceiveNotificationAboutEventsOfIds(duration, ids string) error {
	d, err := time.ParseDuration(duration)
	if err != nil {
		return fmt.Errorf("error when parse duration `%s`", err)
	}

	config := GetTestConfig()
	if config.SenderOutputPath == "" {
		return errors.New("sender output path is requried for this test, set SENDER_OUTPUT_PATH env var")
	}

	file, err := openFileForRead(config.SenderOutputPath, true, true)

	if err != nil {
		return fmt.Errorf("error when open file %s for read-only %s", config.SenderOutputPath, err)
	}

	defer func() {
		if file != nil {
			_ = file.Close()
		}
	}()

	log.Printf("Wait for %s...", d)

	// wait
	time.Sleep(d)

	if file == nil {
		file, err = openFileForRead(config.SenderOutputPath, false, false)
		if err != nil {
			return fmt.Errorf("error when open file %s for read-only %s", config.SenderOutputPath, err)
		}
		// closer already defined
	}

	data, err := ioutil.ReadAll(file)

	if err != nil {
		return fmt.Errorf("error while read file %s: %s", config.SenderOutputPath, err)
	}

	// simple parsing, suppose that each log record in one line
	logRecords := strings.Split(string(data), "\n")

	var receivedEventIds []int
	for _, logRecord := range logRecords {

		logRecord := strings.TrimSpace(logRecord)

		// skip empty strings
		if logRecord == "" {
			continue
		}

		logRecordJson, err := jsonUnmarshalToMap(logRecord)
		if err != nil {
			return fmt.Errorf("error when unmarchal logRecord `%s`: %s", logRecord, err)
		}
		msg, ok := readStringValueFromMap("msg", logRecordJson)
		if !ok || msg != "LogSender.Send" {
			return fmt.Errorf("unexpected value of key `msg`: `%s`", msg)
		}
		eventSerialized, ok := readStringValueFromMap("eventSerialized", logRecordJson)
		if !ok {
			return errors.New("expected readable key `eventSerialized` in log record")
		}
		eventSerializedJson, err := jsonUnmarshalToMap(eventSerialized)
		if err != nil {
			return fmt.Errorf("error when unmarchal `eventSerialized` from log record `%s`: %s", eventSerialized, err)
		}
		eventId, ok := readIntValueFromMap("id", eventSerializedJson)
		if !ok || eventId <= 0 {
			return fmt.Errorf("unexpected value of key `id`: `%d`", eventId)
		}
		receivedEventIds = append(receivedEventIds, eventId)
	}

	expectedEventIds, err := splitStrIntoIntSlice(ids)

	if !isIntSliceEquals(expectedEventIds, receivedEventIds) {
		return fmt.Errorf("expected slice `%v` not equals received `%v`", expectedEventIds, receivedEventIds)
	}

	return nil
}

func (t *featureTest) fieldOfRecordsMustBeNotNil(field, ids string) error {

	eventIds, err := splitStrIntoIntSlice(ids)
	if err != nil {
		return err
	}

	if len(eventIds) <= 0 {
		return fmt.Errorf("unexpected empty event ids list after parse ids `%s`", ids)
	}

	config := GetTestConfig()

	for _, eventId := range eventIds {
		event, err := config.DbStorage.GetEvent(eventId)
		if err != nil {
			return fmt.Errorf("error when get event by id %d from db `%s`", eventId, err)
		}
		if field == "notified_time" {
			if !event.IsNotified() {
				return fmt.Errorf("looks like mark about notifying is not set, `notified_time` is NULL (nil) for event\n`%#v\n`", event)
			}
		} else {
			return fmt.Errorf("unknown field `%s`", field)
		}
	}

	return nil
}

func FeatureContext(s *godog.Suite, t *featureTest) {
	s.Step(`^I send "([^"]*)" request to "([^"]*)" with "([^"]*)" params:$`, t.iSendRequestToWithParams)
	s.Step(`^I send "([^"]*)" request to "([^"]*)"$`, t.iSendRequestTo)
	s.Step(`^The response code should be (\d+)$`, t.theResponseCodeShouldBe)
	s.Step(`^The response contentType should be "([^"]*)"$`, t.theResponseContentTypeShouldBe)
	s.Step(`^The response json should has field "([^"]*)" with value match "([^"]*)"$`, t.theResponseJsonShouldHasFieldWithValueMatch)
	s.Step(`^Extracted number is event id$`, t.extractedNumberIsEventId)
	s.Step(`^The record should match:$`, t.theRecordShouldMatch)
	s.Step(`^Clean DB$`, t.cleanDB)
	s.Step(`^Existing records:$`, t.existingRecords)
	s.Step(`^The records should match:$`, t.theRecordsShouldMatch)
	s.Step(`^The response json should match:$`, t.theResponseJsonShouldMatch)
	s.Step(`^The DB should be clean$`, t.theDBShouldBeClean)
	s.Step(`^The response json is EventListResponse filled with events of ids "([^"]*)"$`, t.theResponseJsonIsEventListResponseFilledWithEventsOfIds)
	s.Step(`^I after wait "([^"]*)" should receive notification about events of ids "([^"]*)"$`, t.iAfterWaitShouldReceiveNotificationAboutEventsOfIds)
	s.Step(`^Field "([^"]*)" of records "([^"]*)" must be not nil$`, t.fieldOfRecordsMustBeNotNil)

}
