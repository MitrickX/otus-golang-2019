package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

type createEventTest struct {
	responseStatusCode int
	responseBody       []byte
	contentType        string
}

func (t *createEventTest) iSendRequestToWithParams(httpMethod, addr, contentType string, data *gherkin.DocString) error {

	var r *http.Response
	var err error

	switch httpMethod {
	case http.MethodPost:
		replacer := strings.NewReplacer("\n", "", "\t", "")
		cleanData := replacer.Replace(data.Content)
		r, err = http.Post(addr, contentType, bytes.NewReader([]byte(cleanData)))
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown method: %s", httpMethod)
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	t.responseStatusCode = r.StatusCode
	t.responseBody = body
	t.contentType = r.Header.Get("Content-Type")

	return nil
}

func (t *createEventTest) theResponseCodeShouldBe(code int) error {
	if t.responseStatusCode != code {
		return fmt.Errorf("unexpected status code: %d != %d", t.responseStatusCode, code)
	}
	return nil
}

func (t *createEventTest) theResponseContentTypeShouldBe(contentType string) error {
	if t.contentType != contentType {
		return fmt.Errorf("unexpected content type: `%s` != `%s`", t.contentType, contentType)
	}
	return nil
}

func (t *createEventTest) theResponseJsonShouldHasFieldWithValueMatch(field, pattern string) error {
	jsonResponse := make(map[string]string)
	err := json.Unmarshal(t.responseBody, &jsonResponse)
	if err != nil {
		return fmt.Errorf("json unmarshal error %s", err)
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

	return nil
}

func FeatureContext(s *godog.Suite) {
	t := new(createEventTest)
	s.Step(`^I send "([^"]*)" request to "([^"]*)" with "([^"]*)" params:$`, t.iSendRequestToWithParams)
	s.Step(`^The response code should be (\d+)$`, t.theResponseCodeShouldBe)
	s.Step(`^The response contentType should be "([^"]*)"$`, t.theResponseContentTypeShouldBe)
	s.Step(`^The response json should has field "([^"]*)" with value match "([^"]*)"$`, t.theResponseJsonShouldHasFieldWithValueMatch)

}
