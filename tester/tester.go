package tester

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Tester provides methods to test a HTTP API client.
type Tester struct {
	Categories map[string]TestCategory
	URL        string

	mux    *http.ServeMux
	server *httptest.Server
}

/*
NewTester creates a new Tester
	categories := map[string]tester.TestCategory{
		"users": tester.TestCategory{
			Scenarios: map[string]tester.TestScenario{
				"create": TestScenario{
					Request: TestRequest{
						Method: "POST",
						Path:   "/user",
						Headers: map[string]string{"Content-Type": "application/json"},
						Body: strings.NewReader("{\"username\":\"user1\"}"),
					},
					Response: TestResponse{
						Code: 200,
						Headers: map[string]string{"Content-Type": "application/json"},
						Body: strings.NewReader("{\"id\":1,\"username\":\"user1\"}"),
					},
				},
			},
		}
	}

	api := tester.NewTester(categories)
*/
func NewTester(categories map[string]TestCategory) *Tester {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)

	t := Tester{
		Categories: categories,
		URL:        server.URL,
		mux:        mux,
		server:     server,
	}

	return &t
}

/*
Close stops the HTTP server.

Use defer api.Close() in each test case to ensure the server is closed and resources are released.

Example:
	api := tester.NewTester(categories)
	defer api.Close()
*/
func (tester *Tester) Close() {
	if tester.server != nil {
		tester.server.Close()
	}
}

// Scenario returns a configured TestScenario.
func (tester *Tester) Scenario(category string, scenario string) (*TestScenario, error) {
	c, ok := tester.Categories[category]
	if !ok {
		return nil, fmt.Errorf("category %s was not found", category)
	}

	s, ok := c.Scenarios[scenario]
	if !ok {
		return nil, fmt.Errorf("scenario %s was not found in category %s", scenario, category)
	}

	return &s, nil
}

/*
Setup configures a route for the TestScenario.

Example:
	api.Setup(t, "user", "online")

	err := c.WaitUserOnline("user1")
	if err != nil {
		t.Fatal(err)
	}
*/
func (tester *Tester) Setup(t *testing.T, category string, scenario string) {
	s, err := tester.Scenario(category, scenario)
	if err != nil {
		t.Fatal(err)
	}

	tester.mux.HandleFunc(s.Request.Path, func(w http.ResponseWriter, req *http.Request) {
		CompareRequests(t, &s.Request, req)
		WriteResponse(t, &s.Response, w)
	})
}

/*
Do creates a route for the handler function provided.

Example:
	offline, err := api.Scenario("user", "offline")
	if err != nil {
		t.Fatal(err)
	}

	online, err := api.Scenario("user", "online")
	if err != nil {
		t.Fatal(err)
	}

	scenarios := []*tester.TestScenario{offline, online}
	current := 0

	api.Do(offline.Path, func(w http.ResponseWriter, req *http.Request) {
		tester.WriteResponse(t, &scenarios[current].Response, w)
		current++
	})

	err := c.WaitUserOnline("user1")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 2, current)
*/
func (tester *Tester) Do(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	tester.mux.HandleFunc(pattern, handler)
}

// CompareRequests compare TestScenario to an actual HTTP request
func CompareRequests(t *testing.T, expected *TestRequest, actual *http.Request) {
	actualBody, err := ioutil.ReadAll(actual.Body)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expected.Method, actual.Method)
	for k, v := range expected.Headers {
		assert.Contains(t, actual.Header, k)
		assert.Equal(t, v, actual.Header[k][0])
	}

	if expected.Body != nil {
		// Rewind expected body
		_, err = expected.Body.Seek(0, 0)
		if err != nil {
			t.Fatal(err)
		}

		expectedBody, err := ioutil.ReadAll(expected.Body)
		if err != nil {
			t.Fatal(err)
		}

		if cType, ok := expected.Headers["Content-Type"]; ok && cType == "application/json" {
			var expectedData interface{}
			if len(expectedBody) > 0 {
				if err := json.Unmarshal(expectedBody, &expectedData); err != nil {
					t.Fatal(err)
				}
			}

			var actualData interface{}
			if len(actualBody) > 0 {
				if err := json.Unmarshal(actualBody, &actualData); err != nil {
					t.Fatal(err)
				}
			}

			assert.Equal(t, expectedData, actualData)
		} else {
			assert.Equal(t, expectedBody, actualBody)
		}
	}
}

// WriteResponse writes TestResponse to an actual HTTP ResponseWriter
func WriteResponse(t *testing.T, resp *TestResponse, w http.ResponseWriter) {
	for k, v := range resp.Headers {
		if k == "Content-Length" {
			continue
		}

		w.Header().Add(k, v)
	}

	if resp.Body != nil {
		// Rewind body
		_, err := resp.Body.Seek(0, 0)
		if err != nil {
			t.Fatal(err)
		}

		resBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}

		resLen := strconv.Itoa(len(resBody))
		w.Header().Add("Content-Length", resLen)
		w.WriteHeader(resp.Code)

		_, err = w.Write(resBody)
		if err != nil {
			t.Fatal(err)
		}
	} else {
		w.WriteHeader(resp.Code)
	}
}
