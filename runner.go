package testx

import (
	"testing"

	"github.com/drykit-go/testx/check"
	"github.com/drykit-go/testx/check/checkconv"
)

type Runner interface {
	Run(t *testing.T)
}

type baseRunner struct {
	// t      *testing.T
	checks []testCheck
}

func (r *baseRunner) addCheck(c testCheck) {
	r.checks = append(r.checks, c)
}

// addChecks is unused for now. It could be used to avoid an addXxxChecks method
// for each typed Checker, but it implies to make the conversion from []TChecker
// to []interface{} upstream (which requires to iterate on the slice just for
// the conversion).
func (r *baseRunner) addChecks(label string, get getFunc, checks []interface{}) {
	for _, c := range checks {
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

func (r *baseRunner) run(t *testing.T) {
	for _, current := range r.checks {
		got := current.get()
		if !current.check.Pass(got) {
			fail(t, current.check.Explain(current.label, got))
		}
	}
}
