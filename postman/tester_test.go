package postman

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSomething(t *testing.T) {
	tester, err := NewTester("testdata/sample.json")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 1, len(tester.Categories))
	assert.Contains(t, tester.Categories, "test_category")
	assert.Equal(t, 2, len(tester.Categories["test_category"].Scenarios))
	assert.Contains(t, tester.Categories["test_category"].Scenarios, "test_scenario1")
	assert.Contains(t, tester.Categories["test_category"].Scenarios, "test_scenario2")
}
