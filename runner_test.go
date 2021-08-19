package testx_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/drykit-go/testx"
)

var deq = reflect.DeepEqual

type results struct {
	checks                    []testx.CheckResult
	passed, failed            bool
	nPassed, nFailed, nChecks int
	execTime                  time.Duration
}

func assertEqualResults(t *testing.T, res testx.Resulter, exp results) {
	if got := toTestResults(res); !deq(got, exp) {
		t.Errorf("bad results\nexp %#v\ngot %#v", exp, got)
	}
}

func toTestResults(res testx.Resulter) results {
	return results{
		checks:   res.Checks(),
		passed:   res.Passed(),
		failed:   res.Failed(),
		nPassed:  res.NPassed(),
		nFailed:  res.NFailed(),
		nChecks:  res.NChecks(),
		execTime: res.ExecTime(),
	}
}
