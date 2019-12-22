package tests

import (
	"encoding/json"
	"fmt"
	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	"github.com/mitrickx/otus-golang-2019/29/calendar/internal/domain/entities"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"
)

const (
	dateTimeLayout      = "2006-01-02 15:04:05"
	dateTimeShortLayout = "2006-01-02 15:04"
	dateLayout          = "2006-01-02"
)

var (
	indexToIdRegExp = regexp.MustCompile(`id=(\d+)`)
)

// replace number in id=\d+ by value in eventIds found by that number as index
// if something wrong return original string
func replaceIndexToEventId(query string, eventIds []int) string {
	re := indexToIdRegExp
	m := re.Match([]byte(query))
	if !m {
		return query
	}
	sb := re.FindStringSubmatch(query)
	if len(sb) < 2 {
		return query
	}
	index, err := strconv.Atoi(sb[1])
	if err != nil {
		return query
	}
	if len(eventIds) <= index {
		return query
	}
	return strings.Replace(query, "id="+sb[1], "id="+strconv.Itoa(eventIds[index]), 1)
}

// Convert table data to events
func convertGherkinTableEvents(data *gherkin.DataTable) ([]entities.Event, error) {

	if len(data.Rows) <= 1 {
		return nil, nil
	}

	var columns []string

	var rows = data.Rows
	var firstRow = rows[0]

	for _, cell := range firstRow.Cells {
		columns = append(columns, cell.Value)
	}

	columnsCount := len(columns)

	var events []entities.Event

	rows = rows[1:]

	for rowIndex, row := range rows {

		var eventId, beforeMinutes int
		var start, end entities.DateTime
		var isNotifyingEnabled, isNotified bool
		var notifiedTime time.Time
		var name string

		for cellIndex, cell := range row.Cells {

			if cellIndex >= columnsCount {
				return nil, fmt.Errorf("convert from gherkin table failed, unknown column by cellIndex %d", cellIndex)
			}

			var err error

			columnName := columns[cellIndex]

			cellValue := strings.TrimSpace(cell.Value)
			val := strings.ToLower(cellValue)
			isNull := val == "" || val == "nil" || val == "null"

			switch columnName {
			case "id":
				eventId, err = strconv.Atoi(cellValue)
				if err != nil {
					return nil, fmt.Errorf(
						"conver from gherkin table failed, can't cast cell (%d, `%s`) to int", rowIndex,
						columnName,
					)
				}
			case "name":
				name = cellValue
			case "start_time":
				t, err := parseStrToTime(cellValue)
				if err != nil {
					return nil, fmt.Errorf(
						"conver from gherkin table failed, can't cast cell (%d, `%s`) to entities.DateTime",
						rowIndex,
						columnName,
					)
				}
				start = entities.ConvertFromTime(t)
			case "end_time":
				t, err := parseStrToTime(cellValue)
				if err != nil {
					return nil, fmt.Errorf(
						"conver from gherkin table failed, can't cast cell (%d, `%s`) to entities.DateTime",
						rowIndex,
						columnName,
					)
				}
				end = entities.ConvertFromTime(t)
			case "before_minutes":
				if !isNull {
					isNotifyingEnabled = true
					beforeMinutes, err = strconv.Atoi(cellValue)
					if err != nil {
						return nil, fmt.Errorf(
							"conver from gherkin table failed, can't cast cell (%d, `%s`) to int",
							rowIndex,
							"before_minutes",
						)
					}
				}
			case "notified_time":
				if !isNull {
					isNotified = true
					notifiedTime, err = parseStrToTime(cellValue)
					if err != nil {
						return nil, fmt.Errorf(
							"conver from gherkin table failed, can't cast cell (%d, `%s`) to entities.DateTime",
							rowIndex,
							"notified_time",
						)
					}
				}
			}
		}

		event := entities.NewDetailedEventWithId(
			eventId,
			name,
			start,
			end,
			isNotifyingEnabled,
			beforeMinutes,
			isNotified,
			notifiedTime,
		)
		events = append(events, event)
	}

	return events, nil
}

// parse string representation of date time into time.Time
func parseStrToTime(str string) (time.Time, error) {
	var t time.Time
	var err error
	t, err = time.Parse(dateTimeLayout, str)
	if err != nil {
		t, err = time.Parse(dateTimeShortLayout, str)
	}
	if err != nil {
		t, err = time.Parse(dateLayout, str)
	}
	return t, err
}

// Assert statue code in response equal passed
func assertStatusCode(r *http.Response, code int) error {
	if r.StatusCode != code {
		return fmt.Errorf("unexpected status code: %d != %d", r.StatusCode, code)
	}
	return nil
}

// Assert content type header in response equal passed
func assertContentType(r *http.Response, contentType string) error {
	respContentType := r.Header.Get("Content-Type")
	if respContentType != contentType {
		return fmt.Errorf("unexpected content type: `%s` != `%s`", respContentType, contentType)
	}
	return nil
}

// Treat response body as json and convert it to map[string]string
func readStringToStringMapFromJsonBody(r *http.Response) (map[string]string, error) {
	var err error
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	jsonResponse := make(map[string]string)
	err = json.Unmarshal(body, &jsonResponse)
	if err != nil {
		return nil, fmt.Errorf("json unmarshal error %s", err)
	}

	return jsonResponse, nil
}

// Test entry point
func TestMain(m *testing.M) {

	status := godog.RunWithOptions("integration", func(s *godog.Suite) {
		t := newFeatureTest()
		FeatureContext(s, t)
	}, godog.Options{
		Format:    "pretty", // progress, pretty
		Paths:     []string{"../features"},
		Randomize: 0,
	})

	if st := m.Run(); st > status {
		status = st
	}

	os.Exit(status)
}
