package testx

import (
	"testing"

	"github.com/drykit-go/testx/check"
	"github.com/drykit-go/testx/check/checkconv"
)

type baseTest struct {
	// t      *testing.T
	checks []testCheck
}

func (test *baseTest) addCheck(c testCheck) {
	test.checks = append(test.checks, c)
}

// addChecks is unused for now. It could be used to avoid an addXxxChecks method
// for each typed Checker, but it implies to make the conversion from []TChecker
// to []interface{} upstream (which requires to iterate on the slice just for
// the conversion).
func (test *baseTest) addChecks(label string, get getFunc, checks []interface{}) {
	for _, c := range checks {
		test.addCheck(testCheck{label: label, get: get, check: checkconv.UntypedChecker(c)})
	}
}

func (test *baseTest) addIntChecks(label string, get getFunc, checks []check.IntChecker) {
	for _, c := range checks {
		test.addCheck(testCheck{label: label, get: get, check: checkconv.FromInt(c)})
	}
}

func (test *baseTest) addBytesChecks(label string, get getFunc, checks []check.BytesChecker) {
	for _, c := range checks {
		test.addCheck(testCheck{label: label, get: get, check: checkconv.FromBytes(c)})
	}
}

func (test *baseTest) addStringChecks(label string, get getFunc, checks []check.StringChecker) {
	for _, c := range checks {
		test.addCheck(testCheck{label: label, get: get, check: checkconv.FromString(c)})
	}
}

func (test *baseTest) addDurationChecks(label string, get getFunc, checks []check.DurationChecker) {
	for _, c := range checks {
		test.addCheck(testCheck{label: label, get: get, check: checkconv.FromDuration(c)})
	}
}

func (test *baseTest) addHTTPHeaderChecks(label string, get getFunc, checks []check.HTTPHeaderChecker) {
	for _, c := range checks {
		test.addCheck(testCheck{label: label, get: get, check: checkconv.FromHTTPHeader(c)})
	}
}

func (test *baseTest) run(t *testing.T) {
	for _, current := range test.checks {
		got := current.get()
		if !current.check.Pass(got) {
			fail(t, current.check.Explain(current.label, got))
		}
	}
}
