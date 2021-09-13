package testx_test

import (
	"fmt"
	"testing"

	"github.com/drykit-go/testx"
	"github.com/drykit-go/testx/check"
	"github.com/drykit-go/testx/checkconv"
)

// Example

func ExampleValueRunner() {
	results := testx.Value(42).
		Exp(42).
		ExpNot(3, "hello").
		Pass(checkconv.FromInt(check.Int.InRange(41, 43))).
		// Run(t) // can be used in a test func
		DryRun()

	fmt.Println(results.Passed())
	// Output: true
}

// Tests

func TestValueRunner(t *testing.T) {
	t.Run("should pass", func(t *testing.T) {
		res := testx.Value(42).
			Exp(42).
			ExpNot(3, "hello").
			Pass(checkconv.FromInt(check.Int.InRange(41, 43))).
			DryRun()

		exp := baseResults{
			passed:  true,
			failed:  false,
			nPassed: 3,
			nFailed: 0,
			nChecks: 3,
			checks: []testx.CheckResult{
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
			ExpNot(99).
			ExpNot(99, 42).
			Pass(checkconv.FromInt(check.Int.LT(10))).
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
