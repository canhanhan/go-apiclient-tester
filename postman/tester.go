package postman

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/finarfin/go-apiclient-tester/tester"
)

/*
NewTester returns a Tester initialized from a Postman collection

Example:
	tester, err := postman.NewTester("testdata/sample_collection.json")
	if err != nil {
		t.Fatal(err)
	}
	defer tester.Close()
	tester.Setup("user", "create_success")

	_, err = c.CreateUser("user1")
	if err != nil {
		t.Fatal(err)
	}
*/
func NewTester(path string) (*tester.Tester, error) {
	categories, err := parse(path)
	if err != nil {
		return nil, err
	}

	return tester.NewTester(categories), nil
}

func replaceVars(c postmanCollection, value string) string {
	for _, v := range c.Variables {
		value = strings.ReplaceAll(value, fmt.Sprintf("{{%s}}", v.Key), v.Value)
	}

	return value
}

func parse(path string) (categories map[string]tester.TestCategory, err error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}

	var c postmanCollection
	if err = json.Unmarshal(b, &c); err != nil {
		return
	}

	categories = make(map[string]tester.TestCategory)
	for _, v := range c.Items {
		scenarios := make(map[string]tester.TestScenario)
		for _, s := range v.Response {
			reqHdr := make(map[string]string)
			for _, h := range s.OriginalRequest.Headers {
				reqHdr[h.Key] = replaceVars(c, h.Value)
			}

			resHdr := make(map[string]string)
			for _, h := range s.Headers {
				resHdr[h.Key] = replaceVars(c, h.Value)
			}

			scenarios[s.Name] = tester.TestScenario{
				Request: tester.TestRequest{
					Method:  s.OriginalRequest.Method,
					Path:    replaceVars(c, "/"+strings.Join(s.OriginalRequest.URL.Path, "/")),
					Headers: reqHdr,
					Body:    strings.NewReader(replaceVars(c, s.OriginalRequest.Body.Raw)),
				},
				Response: tester.TestResponse{
					Code:    s.Code,
					Status:  s.Status,
					Headers: resHdr,

					Body: strings.NewReader(replaceVars(c, s.Body)),
				},
			}
		}

		categories[v.Name] = tester.TestCategory{
			Scenarios: scenarios,
		}
	}

	return
}
