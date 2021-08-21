package testx

import (
	"testing"
	"time"

	"github.com/drykit-go/testx/check"
	"github.com/drykit-go/testx/check/checkconv"
)

type (
	gotType interface{}
	getFunc func() gotType

	testCheck struct {
		label string
		get   func() gotType
		check check.UntypedChecker
	}
)

type baseRunner struct {
	checks      []testCheck
	baseResults baseResults
}

func (r *baseRunner) addCheck(c testCheck) {
	r.checks = append(r.checks, c)
}

func (r *baseRunner) addChecks(label string, get getFunc, checks []interface{}) {
	for _, c := range checks {
		if !checkconv.IsChecker(c) {
			panic("invalid checker provided to MustPass")
		}
		r.addCheck(testCheck{label: label, get: get, check: checkconv.UntypedChecker(c)})
	}
}

func (r *baseRunner) addIntChecks(label string, get getFunc, checks []check.IntChecker) {
	for _, c := range checks {
		r.addCheck(testCheck{label: label, get: get, check: checkconv.FromInt(c)})
	}
}

func (r *baseRunner) addBytesChecks(label string, get getFunc, checks []check.BytesChecker) {
	for _, c := range checks {
		r.addCheck(testCheck{label: label, get: get, check: checkconv.FromBytes(c)})
	}
}

func (r *baseRunner) addStringChecks(label string, get getFunc, checks []check.StringChecker) {
	for _, c := range checks {
		r.addCheck(testCheck{label: label, get: get, check: checkconv.FromString(c)})
	}
}

func (r *baseRunner) addDurationChecks(label string, get getFunc, checks []check.DurationChecker) {
	for _, c := range checks {
		r.addCheck(testCheck{label: label, get: get, check: checkconv.FromDuration(c)})
	}
}

func (r *baseRunner) addHTTPHeaderChecks(label string, get getFunc, checks []check.HTTPHeaderChecker) {
	for _, c := range checks {
		r.addCheck(testCheck{label: label, get: get, check: checkconv.FromHTTPHeader(c)})
	}
}

func (r *baseRunner) addUntypedChecks(label string, get getFunc, checks []check.UntypedChecker) {
	for _, c := range checks {
		r.addCheck(testCheck{label: label, get: get, check: c})
	}
}

func (r *baseRunner) run(t *testing.T) {
	for _, current := range r.checks {
		got := current.get()
		if !current.check.Pass(got) {
			r.fail(t, current.check.Explain(current.label, got))
		}
	}
}

func (r *baseRunner) fail(t *testing.T, msg string) {
	t.Error(msg)
}

func (r *baseRunner) updateBaseResults(c testCheck) {
	got := c.get()
	passed := c.check.Pass(got)
	reason := condString("", c.check.Explain(c.label, got), passed)

	r.baseResults.checks = append(r.baseResults.checks, CheckResult{
		Passed: passed,
		Reason: reason,
	})
	if !passed {
		r.baseResults.nFailed++
	}
}

type baseResults struct {
	checks   []CheckResult
	nFailed  int
	execTime time.Duration
}

func (r baseResults) Checks() []CheckResult {
	return r.checks
}

func (r baseResults) Passed() bool {
	return r.nFailed == 0
}

func (r baseResults) Failed() bool {
	return !r.Passed()
}

func (r baseResults) NPassed() int {
	return r.NChecks() - r.nFailed
}

func (r baseResults) NFailed() int {
	return r.nFailed
}

func (r baseResults) NChecks() int {
	return len(r.checks)
}

func (r baseResults) ExecTime() time.Duration {
	return r.execTime
}
