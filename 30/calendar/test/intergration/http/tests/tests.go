package tests

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/DATA-DOG/godog/gherkin"
	"github.com/mitrickx/otus-golang-2019/30/calendar/internal/domain/entities"
	serviceHttp "github.com/mitrickx/otus-golang-2019/30/calendar/internal/http"
)

// In this file bunch of helpers for tests module

const (
	dateTimeLayout      = "2006-01-02 15:04:05"
	dateTimeShortLayout = "2006-01-02 15:04"
	dateLayout          = "2006-01-02"
	timeLayout          = "15:04:05"
)

// Compiled regexp for parseAddDuration
var addDurationRegexp = regexp.MustCompile(`^([+\-])(\d+(?:s|m|h))$`)

// Convert table data to events
func convertGherkinTableToEvents(data *gherkin.DataTable) ([]entities.Event, error) {

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
						"conver from gherkin table failed, can't cast cell (%d, `%s`, `%s`) to int: %s", rowIndex,
						columnName,
						cellValue,
						err,
					)
				}
			case "name":
				name = cellValue
			case "start_time":
				t, err := parseStrToTime(cellValue)
				if err != nil {
					return nil, fmt.Errorf(
						"conver from gherkin table failed, can't cast cell (%d, `%s`, `%s`) to entities.DateTime: %s",
						rowIndex,
						columnName,
						cellValue,
						err,
					)
				}
				start = entities.ConvertFromTime(t)
			case "end_time":
				t, err := parseStrToTime(cellValue)
				if err != nil {
					return nil, fmt.Errorf(
						"conver from gherkin table failed, can't cast cell (%d, `%s`, `%s`) to entities.DateTime: %s",
						rowIndex,
						columnName,
						cellValue,
						err,
					)
				}
				end = entities.ConvertFromTime(t)
			case "before_minutes":
				if !isNull {
					isNotifyingEnabled = true
					beforeMinutes, err = strconv.Atoi(cellValue)
					if err != nil {
						return nil, fmt.Errorf(
							"conver from gherkin table failed, can't cast cell (%d, `%s`, `%s`) to int: %s",
							rowIndex,
							"before_minutes",
							cellValue,
							err,
						)
					}
				}
			case "notified_time":
				if !isNull {
					isNotified = true
					notifiedTime, err = parseStrToTime(cellValue)
					if err != nil {
						return nil, fmt.Errorf(
							"conver from gherkin table failed, can't cast cell (%d, `%s`, `%s`) to entities.DateTime: %s",
							rowIndex,
							"notified_time",
							cellValue,
							err,
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

// Parse strings like '+10s', '-13m', '+1h' into time.Time
// Adding applying to now time
func parseAddDuration(str string) (time.Time, error) {
	submatch := addDurationRegexp.FindStringSubmatch(str)
	if len(submatch) <= 0 {
		return time.Time{}, errors.New("not match")
	}
	d, err := time.ParseDuration(submatch[2])
	if err != nil {
		return time.Time{}, err
	}
	if submatch[1] == "+" {
		return time.Now().Add(d), nil
	} else {
		return time.Now().Add(-d), nil
	}
}

// parse string representation of date time into time.Time
func parseStrToTime(str string) (time.Time, error) {
	t, err := parseStrToTimeByLayouts(str)
	if err == nil {
		return t, err
	}

	// +someDuration, -someDuration
	t, err = parseAddDuration(str)
	if err == nil {
		return t, err
	}

	// Y-m-d H:i:s - it is now time
	// Mon, Tue, Wed, Thu, Fri, Sat, Sun - day of current week
	// ld - last day of month
	nowTime := time.Now()
	dateStr := nowTime.Format(dateLayout)
	timeStr := nowTime.Format(timeLayout)
	dateParts := strings.Split(dateStr, "-")
	timeParts := strings.Split(timeStr, ":")
	if len(dateParts) < 3 {
		return time.Time{}, fmt.Errorf("error when parse now time, date component is `%s`", dateStr)
	}
	if len(timeParts) < 3 {
		return time.Time{}, fmt.Errorf("error when parse now time, time component is `%s`", timeStr)
	}

	weekDay := nowTime.Weekday()
	shiftDay := weekDay - 1

	monTime := nowTime.Add(-24 * time.Duration(shiftDay) * time.Hour)
	tueTime := monTime.Add(24 * time.Hour)
	wedTime := tueTime.Add(24 * time.Hour)
	thuTime := wedTime.Add(24 * time.Hour)
	friTime := thuTime.Add(24 * time.Hour)
	satTime := friTime.Add(24 * time.Hour)
	sunTime := satTime.Add(24 * time.Hour)

	monStr := monTime.Format(dateLayout)
	tueStr := tueTime.Format(dateLayout)
	wedStr := wedTime.Format(dateLayout)
	thuStr := thuTime.Format(dateLayout)
	friStr := friTime.Format(dateLayout)
	satStr := satTime.Format(dateLayout)
	sunStr := sunTime.Format(dateLayout)

	ld := "31"
	switch nowTime.Month() {
	case time.January, time.March, time.May, time.July, time.August, time.October, time.December:
		ld = "31"
	case time.April, time.June, time.September, time.November:
		ld = "30"
	default:
		year := nowTime.Year()
		isLeap := year%4 == 0 && (year%100 != 0 || year%400 == 0)
		if isLeap {
			ld = "29"
		} else {
			ld = "28"
		}
	}

	replacePairs := []string{
		"Y", dateParts[0],
		"m", dateParts[1],
		"d", dateParts[2],
		"H", timeParts[0],
		"i", timeParts[1],
		"s", timeParts[2],
		"Mon", monStr,
		"Tue", tueStr,
		"Wed", wedStr,
		"Thu", thuStr,
		"Fri", friStr,
		"Sat", satStr,
		"Sun", sunStr,
		"ld", ld,
	}

	replacer := strings.NewReplacer(replacePairs...)
	str = replacer.Replace(str)

	return parseStrToTimeByLayouts(str)
}

// parse string representation of date time into time.Time by predefined layouts
func parseStrToTimeByLayouts(str string) (time.Time, error) {
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

// Is 2 int slice equals despite their orders
func isIntSliceEquals(a []int, b []int) bool {
	aCloned := append([]int(nil), a...)
	bCloned := append([]int(nil), b...)
	sort.Ints(aCloned)
	sort.Ints(bCloned)
	return reflect.DeepEqual(aCloned, bCloned)
}

// Open file for read
// If ignoreErrExist and file not exist it is not error returned to this case
// If seekToEnd passed as TRUE seek to end, otherwise not seek at all
func openFileForRead(filepath string, seekToEnd bool, ignoreErrExist bool) (*os.File, error) {

	file, err := os.OpenFile(filepath, os.O_RDONLY, 0666)

	if err != nil {
		if err != os.ErrExist || !ignoreErrExist {
			return nil, err
		}
	}

	if file != nil && seekToEnd {
		_, _ = file.Seek(0, io.SeekEnd)
	}

	return file, err
}

// unmarshal json string to string-string map
func jsonUnmarshalStringToStringMap(data string) (map[string]string, error) {
	result := make(map[string]string)
	err := json.Unmarshal([]byte(data), &result)
	if err != nil {
		return nil, fmt.Errorf("json unmarshal error %s", err)
	}

	return result, nil
}

// unmarshal json string to general map (string-empty interface map)
func jsonUnmarshalToMap(data string) (map[string]interface{}, error) {
	var result map[string]interface{}
	err := json.Unmarshal([]byte(data), &result)
	if err != nil {
		return nil, fmt.Errorf("json unmarshal error %s", err)
	}
	return result, nil
}

// unmarshal json string to specific type serviceHttp.EventListResponse
func jsonUnmarshalEventListResponse(data string) (serviceHttp.EventListResponse, error) {
	var result serviceHttp.EventListResponse
	err := json.Unmarshal([]byte(data), &result)
	if err != nil {
		return serviceHttp.EventListResponse{}, fmt.Errorf("json unmarshal error %s", err)
	}
	return result, nil
}

// marsha from specific type serviceHttp.EventListResponse to json
func jsonMarshalEventListResponse(response serviceHttp.EventListResponse) string {
	data, _ := json.Marshal(response)
	return string(data)
}

// read string value from general map
func readStringValueFromMap(key string, m map[string]interface{}) (string, bool) {
	value, ok := m[key]
	if !ok {
		return "", false
	}
	s, ok := value.(string)
	if !ok {
		return "", false
	}
	return s, true
}

// read int value from general map
func readIntValueFromMap(key string, m map[string]interface{}) (int, bool) {
	value, ok := m[key]
	if !ok {
		return 0, false
	}
	switch v := value.(type) {
	case int:
		return v, true
	case int64:
		return int(v), true
	case float64:
		return int(v), true
	default:
		return 0, false
	}
}

// Splice string like "1,2,3,4" to []int{1,2,3,4}
func splitStrIntoIntSlice(str string) ([]int, error) {
	parts := strings.Split(str, ",")
	var res []int
	for _, part := range parts {
		intVal, err := strconv.Atoi(strings.TrimSpace(part))
		if err != nil {
			return nil, fmt.Errorf("parse string failed when apply atoi to token `%s`", part)
		}
		res = append(res, intVal)
	}
	return res, nil
}
