package testx

import (
	"testing"

	"github.com/drykit-go/cond"

	"github.com/drykit-go/testx/check"
	"github.com/drykit-go/testx/checkconv"
)

type (
	gottype interface{}
	getfunc func() gottype

	baseCheck struct {
		label        string
		explainLabel string
		get          getfunc
		checker      check.ValueChecker
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
			r.fail(t, current.checker.Explain(current.label, got))
		}
	}
}

func (r *baseRunner) fail(t *testing.T, msg string) {
	t.Helper()
	t.Error(msg)
}

func (r *baseRunner) baseResults() baseResults {
	results := baseResults{}
	for _, bc := range r.checks {
		getExplain := func(bc baseCheck, got interface{}, passed bool) string {
			if passed {
				return ""
			}
			label := cond.DefaultString(bc.explainLabel, bc.label)
			return bc.checker.Explain(label, got)
		}
		got := bc.get()
		passed := bc.checker.Pass(got)
		results.checks = append(results.checks, CheckResult{
			Passed: passed,
			Reason: getExplain(bc, got, passed),
			label:  bc.label,
		})
		if !passed {
			results.nFailed++
		}
	}
	return results
}

type baseResults struct {
	checks  []CheckResult
	nFailed int
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
