package postman

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/finarfin/go-apiclient-tester/common"
)

func NewTester(path string) (*common.Tester, error) {
	categories, err := parse(path)
	if err != nil {
		return nil, err
	}

	return common.NewTester(categories), nil
}

func replaceVars(c PostmanCollection, value string) string {
	for _, v := range c.Variables {
		value = strings.ReplaceAll(value, fmt.Sprintf("{{%s}}", v.Key), v.Value)
	}

	return value
}

func parse(path string) (categories map[string]common.TestCategory, err error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}

	var c PostmanCollection
	if err = json.Unmarshal(b, &c); err != nil {
		return
	}

	categories = make(map[string]common.TestCategory)
	for _, v := range c.Items {
		scenarios := make(map[string]common.TestScenario)
		for _, s := range v.Response {
			reqHdr := make(map[string]string)
			for _, h := range s.OriginalRequest.Headers {
				reqHdr[h.Key] = replaceVars(c, h.Value)
			}

			resHdr := make(map[string]string)
			for _, h := range s.Headers {
				resHdr[h.Key] = replaceVars(c, h.Value)
			}

			scenarios[s.Name] = common.TestScenario{
				Request: common.TestRequest{
					Method:  s.OriginalRequest.Method,
					Path:    replaceVars(c, "/"+strings.Join(s.OriginalRequest.URL.Path, "/")),
					Headers: reqHdr,
					Body:    strings.NewReader(replaceVars(c, s.OriginalRequest.Body.Raw)),
				},
				Response: common.TestResponse{
					Code:    s.Code,
					Status:  s.Status,
					Headers: resHdr,

					Body: strings.NewReader(replaceVars(c, s.Body)),
				},
			}
		}

		categories[v.Name] = common.TestCategory{
			Scenarios: scenarios,
		}
	}

	return
}
