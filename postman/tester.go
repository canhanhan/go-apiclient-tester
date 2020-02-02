package postman

import (
	"encoding/json"
	"github.com/finarfin/go-apiclient-tester/common"
	"io/ioutil"
	"strings"
)

func NewTester(path string) (*common.Tester, error) {
	categories, err := parse(path)
	if err != nil {
		return nil, err
	}

	return &common.Tester{
		Categories: categories,
	}, nil
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
				reqHdr[h.Key] = h.Value
			}

			resHdr := make(map[string]string)
			for _, h := range s.Headers {
				resHdr[h.Key] = h.Value
			}

			scenarios[s.Name] = common.TestScenario{
				Request: common.TestRequest{
					Method:  s.OriginalRequest.Method,
					Path:    s.OriginalRequest.URL.Path[0],
					Headers: reqHdr,
					Body:    strings.NewReader(s.OriginalRequest.Body.Raw),
				},
				Response: common.TestResponse{
					Code:    s.Code,
					Status:  s.Status,
					Headers: resHdr,
					Body:    strings.NewReader(s.Body),
				},
			}
		}

		categories[v.Name] = common.TestCategory{
			Scenarios: scenarios,
		}
	}

	return
}
