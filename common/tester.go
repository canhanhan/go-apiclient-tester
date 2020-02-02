package common

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Tester struct {
	Categories map[string]TestCategory
	URL        string

	mux    *http.ServeMux
	server *httptest.Server
}

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

func (t *Tester) Close() {
	t.server.Close()
}

func (tester *Tester) Setup(t *testing.T, category string, scenario string) {
	c, ok := tester.Categories[category]
	if !ok {
		t.Fatalf("Category %s was not found.", category)
	}

	s, ok := c.Scenarios[scenario]
	if !ok {
		t.Fatalf("Scenario %s was not found in category %s.", scenario, category)
	}

	tester.mux.HandleFunc(s.Request.Path, func(w http.ResponseWriter, req *http.Request) {
		actualBody, err := ioutil.ReadAll(req.Body)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, s.Request.Method, req.Method)
		for k, v := range s.Request.Headers {
			assert.Contains(t, req.Header, k)
			assert.Equal(t, v, req.Header[k][0])
		}

		if s.Request.Body != nil {
			expectedBody, err := ioutil.ReadAll(s.Request.Body)
			if err != nil {
				t.Fatal(err)
			}

			if cType, ok := s.Request.Headers["Content-Type"]; ok && cType == "application/json" {
				expectedData := make(map[string]interface{})
				if err := json.Unmarshal(expectedBody, &expectedData); err != nil {
					t.Fatal(err)
				}

				actualData := make(map[string]interface{})
				if err := json.Unmarshal(actualBody, &actualData); err != nil {
					t.Fatal(err)
				}

				assert.Equal(t, expectedData, actualData)
			} else {
				assert.Equal(t, expectedBody, actualBody)
			}
		}

		for k, v := range s.Response.Headers {
			w.Header().Add(k, v)
		}

		w.WriteHeader(s.Response.Code)

		if s.Response.Body != nil {
			resBody, err := ioutil.ReadAll(s.Response.Body)
			if err != nil {
				t.Fatal(err)
			}

			_, err = w.Write(resBody)
			if err != nil {
				t.Fatal(err)
			}
		}
	})
}
