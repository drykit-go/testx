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
	if got := toBaseResults(res); !deq(got, exp) {
		failBadResults(t, got, exp)
	}
}

func failBadResults(t *testing.T, got, exp interface{}) {
	t.Errorf("bad results\nexp %#v\ngot %#v", exp, got)
}

func toBaseResults(res testx.Resulter) baseResults {
	return baseResults{
		checks:   res.Checks(),
		passed:   res.Passed(),
		failed:   res.Failed(),
		nPassed:  res.NPassed(),
		nFailed:  res.NFailed(),
		nChecks:  res.NChecks(),
		execTime: res.ExecTime(),
	}
}
