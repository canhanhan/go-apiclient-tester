package common

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type SampleRequest struct {
	TestParameter1 string `json:"test1"`
	TestParameter2 int    `json:"test2"`
}

type SampleResponse struct {
	TestParameter1 string `json:"test1"`
	TestParameter2 int    `json:"test2"`
}

func TestMissingPath(t *testing.T) {
	tester := setup()

	tester.Setup(t, "test_category", "test_json")

	_, err := getResponseObject(t, "POST", fmt.Sprintf("%s/sample_wrong_path", tester.URL), nil)

	assert.Error(t, err)
}

func TestResponseHeader(t *testing.T) {
	tester := setup()

	tester.Setup(t, "test_category", "test_header")

	res, err := getResponseObject(t, "POST", fmt.Sprintf("%s/sample_path", tester.URL), nil)

	assert.NoError(t, err)
	assert.Equal(t, "SampleValue", res.Header.Get("SampleHeader"))
}

func TestResponseBodyJson(t *testing.T) {
	tester := setup()
	data := SampleRequest{
		TestParameter1: "req1",
		TestParameter2: 3,
	}

	tester.Setup(t, "test_category", "test_json")
	body, err := getResponseBody(t, "POST", fmt.Sprintf("%s/sample_path", tester.URL), data)
	var result SampleResponse
	if err = json.Unmarshal(body, &result); err != nil {
		t.Fatal(err)
	}

	expected := SampleResponse{
		TestParameter1: "test1",
		TestParameter2: 2,
	}
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestResponseBodyText(t *testing.T) {
	tester := setup()

	tester.Setup(t, "test_category", "test_text")
	result, err := getResponseBody(t, "GET", fmt.Sprintf("%s/sample_path", tester.URL), nil)

	assert.NoError(t, err)
	assert.Equal(t, "Hello World", string(result))
}

func TestResponseNil(t *testing.T) {
	tester := setup()

	tester.Setup(t, "test_category", "test_nil_response")
	_, err := getResponseObject(t, "GET", fmt.Sprintf("%s/sample_path", tester.URL), nil)

	assert.NoError(t, err)
}

func setup() *Tester {
	categories := make(map[string]TestCategory)
	categories["test_category"] = TestCategory{
		Scenarios: map[string]TestScenario{
			"test_json": TestScenario{
				Request: TestRequest{
					Method: "POST",
					Path:   "/sample_path",
					Headers: map[string]string{
						"Content-Type": "application/json",
					},
					Body: strings.NewReader("{\"test1\":\"req1\",\"test2\":3}"),
				},
				Response: TestResponse{
					Code: 200,
					Headers: map[string]string{
						"Content-Type": "application/json",
					},
					Body: strings.NewReader("{\"test1\":\"test1\",\"test2\":2}"),
				},
			},
			"test_header": TestScenario{
				Request: TestRequest{
					Method: "POST",
					Path:   "/sample_path",
				},
				Response: TestResponse{
					Code: 200,
					Headers: map[string]string{
						"SampleHeader": "SampleValue",
					},
				},
			},
			"test_text": TestScenario{
				Request: TestRequest{
					Method: "GET",
					Path:   "/sample_path",
				},
				Response: TestResponse{
					Code: 200,
					Body: strings.NewReader("Hello World"),
				},
			},
			"test_nil_response": TestScenario{
				Request: TestRequest{
					Method: "GET",
					Path:   "/sample_path",
				},
				Response: TestResponse{
					Code: 200,
				},
			},
		},
	}

	return NewTester(categories)
}

func getResponseBody(t *testing.T, method string, url string, data interface{}) ([]byte, error) {
	res, err := getResponseObject(t, method, url, data)

	if err != nil {
		return nil, err
	}

	return ioutil.ReadAll(res.Body)
}

func getResponseObject(t *testing.T, method string, url string, data interface{}) (*http.Response, error) {
	res, err := sendRequest(t, method, url, data)

	if err != nil {
		t.Fatal(err)
	}

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("Server returned: %d (%s)", res.StatusCode, res.Status)
	}

	return res, nil
}

func sendRequest(t *testing.T, method string, url string, data interface{}) (*http.Response, error) {
	reqJSON, err := json.Marshal(data)
	if err != nil {
		t.Fatal(err)
	}

	c := http.Client{}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(reqJSON))
	req.Header.Add("Content-Type", "application/json")
	if err != nil {
		t.Fatal(err)
	}

	return c.Do(req)
}
