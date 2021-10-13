package testx

import (
	"testing"

	"github.com/drykit-go/testx/check"
	"github.com/drykit-go/testx/checkconv"
)

type (
	gottype interface{}
	getfunc func() gottype

	baseCheck struct {
		get      getfunc
		getLabel func() string
		label    string
		checker  check.ValueChecker
	}
)

type baseRunner struct {
	checks []baseCheck
}

func (r *baseRunner) addCheck(c baseCheck) {
	r.checks = append(r.checks, c)
}

func (r *baseRunner) addChecks(label string, get getfunc, checkers []check.ValueChecker, safe bool) {
	for _, c := range checkers {
		if !safe && !checkconv.IsChecker(c) {
			panic("invalid checker provided to MustPass")
		}
		r.addCheck(baseCheck{label: label, get: get, checker: c})
	}
}

func (r *baseRunner) addIntChecks(label string, get getfunc, checks []check.IntChecker) {
	for _, c := range checks {
		r.addCheck(baseCheck{label: label, get: get, checker: checkconv.FromInt(c)})
	}
}

func (r *baseRunner) addBytesChecks(label string, get getfunc, checks []check.BytesChecker) {
	for _, c := range checks {
		r.addCheck(baseCheck{label: label, get: get, checker: checkconv.FromBytes(c)})
	}
}

func (r *baseRunner) addStringChecks(label string, get getfunc, checks []check.StringChecker) {
	for _, c := range checks {
		r.addCheck(baseCheck{label: label, get: get, checker: checkconv.FromString(c)})
	}
}

func (r *baseRunner) addDurationChecks(label string, get getfunc, checks []check.DurationChecker) {
	for _, c := range checks {
		r.addCheck(baseCheck{label: label, get: get, checker: checkconv.FromDuration(c)})
	}
}

func (r *baseRunner) addHTTPHeaderChecks(label string, get getfunc, checks []check.HTTPHeaderChecker) {
	for _, c := range checks {
		r.addCheck(baseCheck{label: label, get: get, checker: checkconv.FromHTTPHeader(c)})
	}
}

func (r *baseRunner) addUntypedChecks(label string, get getfunc, checks []check.ValueChecker) {
	for _, c := range checks {
		r.addCheck(baseCheck{label: label, get: get, checker: c})
	}
}

func (r *baseRunner) run(t *testing.T) {
	t.Helper()
	for _, current := range r.checks {
		got := current.get()
		if !current.checker.Pass(got) {
			r.fail(t, r.explainCheck(current, got, false))
		}
	}
}

func (r *baseRunner) fail(t *testing.T, msg string) {
	t.Helper()
	t.Error(msg)
}

func (r *baseRunner) explainCheck(bc baseCheck, got interface{}, passed bool) string {
	if passed {
		return ""
	}
	var label string
	if bc.getLabel != nil {
		label = bc.getLabel()
	} else {
		label = bc.label
	}
	return bc.checker.Explain(label, got)
}

func (r *baseRunner) baseResults() baseResults {
	res := baseResults{}
	for _, bc := range r.checks {
		got := bc.get()
		passed := bc.checker.Pass(got)
		res.checks = append(res.checks, CheckResult{
			Passed: passed,
			Reason: r.explainCheck(bc, got, passed),
			label:  bc.label,
		})
		if !passed {
			res.nFailed++
		}
	}
	return res
}

type baseResults struct {
	checks  []CheckResult
	nFailed int
}

func (res baseResults) Checks() []CheckResult {
	return res.checks
}

func (res baseResults) Passed() bool {
	return res.nFailed == 0
}

func (res baseResults) Failed() bool {
	return !res.Passed()
}

func (res baseResults) NPassed() int {
	return res.NChecks() - res.nFailed
}

func (res baseResults) NFailed() int {
	return res.nFailed
}

func (res baseResults) NChecks() int {
	return len(res.checks)
}
