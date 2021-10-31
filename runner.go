package testx

import (
	"testing"

	"github.com/drykit-go/testx/check"
)

type (
	gottype any
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

func (r *baseRunner) addCheck(bc baseCheck) {
	r.checks = append(r.checks, bc)
}

func (r *baseRunner) addChecks(label string, get getfunc, checkers []check.Checker[any]) {
	for _, c := range checkers {
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

func (r *baseRunner) dryRun() baseResults {
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

func (r *baseRunner) explainCheck(bc baseCheck, got any, passed bool) string {
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

func (r *baseRunner) fail(t *testing.T, msg string) {
	t.Helper()
	t.Error(msg)
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
