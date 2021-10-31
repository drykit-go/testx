package testx_test

import (
	"testing"

	"github.com/drykit-go/testx"
	"github.com/drykit-go/testx/check"
)

func TestValueRunner(t *testing.T) {
	t.Run("should pass", func(t *testing.T) {
		res := testx.Value(42).
			Exp(42).
			Not(3, -1).
			Pass(check.Int.InRange(41, 43), check.Int.GT(10)).
			DryRun()

		exp := baseResults{
			passed:  true,
			failed:  false,
			nPassed: 4,
			nFailed: 0,
			nChecks: 4,
			checks: []testx.CheckResult{
				{Passed: true, Reason: ""},
				{Passed: true, Reason: ""},
				{Passed: true, Reason: ""},
				{Passed: true, Reason: ""},
			},
		}

		assertEqualBaseResults(t, res, exp)
	})

	t.Run("should fail", func(t *testing.T) {
		res := testx.Value(42).
			Exp(99).
			Not(99).
			Not(99, 42).
			Pass(check.Int.LT(10)).
			DryRun()

		exp := baseResults{
			passed:  false,
			failed:  true,
			nPassed: 1,
			nFailed: 3,
			nChecks: 4,
			checks: []testx.CheckResult{
				{Passed: false, Reason: "value:\nexp 99\ngot 42"},
				{Passed: true, Reason: ""},
				{Passed: false, Reason: "value:\nexp not 42\ngot 42"},
				{Passed: false, Reason: "value:\nexp < 10\ngot 42"},
			},
		}

		assertEqualBaseResults(t, res, exp)
	})
}
