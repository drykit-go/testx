package testx_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/drykit-go/testx"
)

var deq = reflect.DeepEqual

type baseResults struct {
	checks                    []testx.CheckResult
	passed, failed            bool
	nPassed, nFailed, nChecks int
	execTime                  time.Duration
}

func assertEqualBaseResults(t *testing.T, res testx.Resulter, exp baseResults) {
	t.Helper()
	if got := toBaseResults(res); !deq(got, exp) {
		failBadResults(t, "baseResults", got, exp)
	}
}

func failBadResults(t *testing.T, label string, got, exp interface{}) {
	t.Helper()
	t.Errorf("bad results: %s\nexp %#v\ngot %#v", label, exp, got)
}

func toBaseResults(res testx.Resulter) baseResults {
	withLabelRemoved := func(checks []testx.CheckResult) []testx.CheckResult {
		newChecks := make([]testx.CheckResult, len(checks))
		for i, c := range checks {
			newChecks[i] = testx.CheckResult{
				Passed: c.Passed,
				Reason: c.Reason,
			}
		}
		return newChecks
	}

	return baseResults{
		checks:   withLabelRemoved(res.Checks()),
		passed:   res.Passed(),
		failed:   res.Failed(),
		nPassed:  res.NPassed(),
		nFailed:  res.NFailed(),
		nChecks:  res.NChecks(),
		execTime: res.ExecTime(),
	}
}
