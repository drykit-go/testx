package testx_test

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/drykit-go/testx"
	"github.com/drykit-go/testx/internal/fmtexpl"
)

var deq = reflect.DeepEqual

type baseResults struct {
	checks                    []testx.CheckResult
	passed, failed            bool
	nPassed, nFailed, nChecks int
}

func assertEqualBaseResults(t *testing.T, res testx.Resulter, exp baseResults) {
	t.Helper()
	if got := toBaseResults(res); !deq(got, exp) {
		var errs []string

		// Validate len(got.checks), early return if invalid
		if explen, gotlen := len(exp.checks), len(got.checks); explen != gotlen {
			failWithErrors(t, "baseResults", fmtexpl.Pretty("len(checks)", explen, gotlen))
			return
		}

		// Validate remaining fields
		for _, fv := range []struct {
			lab string
			got any
			exp any
		}{
			{lab: "passed", got: got.passed, exp: exp.passed},
			{lab: "failed", got: got.failed, exp: exp.failed},
			{lab: "nPassed", got: got.nPassed, exp: exp.nPassed},
			{lab: "nFailed", got: got.nFailed, exp: exp.nFailed},
			{lab: "nChecks", got: got.nChecks, exp: exp.nChecks},
		} {
			if !deq(fv.exp, fv.got) {
				errs = append(errs, fmtexpl.Pretty(fv.lab, fv.exp, fv.got))
			}
		}

		// Validate got.checks
		for i, gotc := range got.checks {
			expc := exp.checks[i]
			if gotc.Passed != expc.Passed {
				errs = append(errs, fmtexpl.Pretty(
					fmt.Sprintf("checks[%d].Passed", i),
					expc.Passed,
					gotc.Passed,
				))
			}
			if gotc.Reason != expc.Reason {
				errs = append(errs, fmtexpl.Pretty(
					fmt.Sprintf("checks[%d].Reason", i),
					expc.Reason,
					gotc.Reason,
				))
			}
		}

		failWithErrors(t, "baseResults", errs...)
	}
}

func failWithErrors(t *testing.T, label string, errs ...string) {
	t.Helper()
	t.Errorf("bad results: %s\n%s", label, strings.Join(errs, "\n"))
}

func failBadResults(t *testing.T, label string, got, exp any) {
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
		checks:  withLabelRemoved(res.Checks()),
		passed:  res.Passed(),
		failed:  res.Failed(),
		nPassed: res.NPassed(),
		nFailed: res.NFailed(),
		nChecks: res.NChecks(),
	}
}
