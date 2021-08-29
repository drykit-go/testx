package testx_test

import (
	"log"
	"reflect"
	"testing"

	"github.com/drykit-go/testx"
	"github.com/drykit-go/testx/check"
)

// Tests

var expFixedArgs = map[string]interface{}{
	"a0": []byte("arg0"),
	"a2": map[rune][][]float64{'π': {[]float64{3.14}}},
}

// TestTableRunner ensures testx.Table behaves correctly, in particular
// when dealing with functions with multiple inputs and outputs.
func TestTableRunner(t *testing.T) {
	cases := []testx.Case{
		{In: 42, Exp: true},
		{In: 99, Exp: false, Lab: "odd number"},
	}

	const (
		inPos  = 1
		outPos = 2
	)

	a0, a2 := expFixedArgs["a0"], expFixedArgs["a2"]

	t.Run("single in single out", func(t *testing.T) {
		testx.Table(evenSingle, nil).Cases(cases).Run(t)
	})

	t.Run("single in multiple out", func(t *testing.T) {
		testx.Table(evenMultipleOut, &testx.TableConfig{
			OutPos: outPos,
		}).
			Cases(cases).
			Run(t)
	})

	t.Run("multiple in single out", func(t *testing.T) {
		testx.Table(evenMultipleIn, &testx.TableConfig{
			InPos:     inPos,
			FixedArgs: []interface{}{a0, a2}, // len(FixedArgs) == nparams-1
		}).
			Cases(cases).
			Run(t)
	})

	t.Run("multiple in multiple out", func(t *testing.T) {
		testx.Table(evenMultipleInOut, &testx.TableConfig{
			InPos:     inPos,
			OutPos:    outPos,
			FixedArgs: []interface{}{0: a0, 2: a2}, // len(FixedArgs) == nparams
		}).
			Cases(cases).
			Run(t)
	})

	t.Run("using check.IntChecker", func(t *testing.T) {
		testx.Table(double, nil).
			Cases([]testx.Case{
				{In: 21, Exp: check.Int.Is(42)},
				{In: -4, Exp: check.Int.InRange(-10, 0)},
			}).
			Run(t)
	})
}

func TestTableRunnerResults(t *testing.T) {
	t.Run("should pass", func(t *testing.T) {
		res := testx.
			Table(evenSingle, nil).
			Cases([]testx.Case{
				{In: 10, Exp: true, Lab: "even number"},
				{In: 11, Exp: false, Lab: "odd number"},
			}).
			DryRun()

		exp := tableResults{
			baseResults: baseResults{
				passed:  true,
				failed:  false,
				nPassed: 2,
				nFailed: 0,
				nChecks: 2,
				checks: []testx.CheckResult{
					{Passed: true, Reason: ""},
					{Passed: true, Reason: ""},
				},
			},
			passedAt:    map[int]bool{0: true, 1: true},
			failedAt:    map[int]bool{0: false, 1: false},
			passedLabel: map[string]bool{"even number": true, "odd number": true},
			failedLabel: map[string]bool{"even number": false, "odd number": false},
		}

		assertEqualTableResults(t, res, exp)
	})

	t.Run("should fail", func(t *testing.T) {
		res := testx.
			Table(evenSingle, nil).
			Cases([]testx.Case{
				{In: 10, Exp: true, Lab: "even number"},
				{In: 11, Exp: true, Lab: "odd number"},
			}).
			DryRun()

		exp := tableResults{
			baseResults: baseResults{
				passed:  false,
				failed:  true,
				nPassed: 1,
				nFailed: 1,
				nChecks: 2,
				checks: []testx.CheckResult{
					{Passed: true, Reason: ""},
					{Passed: false, Reason: "odd number\ntestx_test.evenSingle(11) -> expect true, got false"},
				},
			},
			passedAt:    map[int]bool{0: true, 1: false},
			failedAt:    map[int]bool{0: false, 1: true},
			passedLabel: map[string]bool{"even number": true, "odd number": false},
			failedLabel: map[string]bool{"even number": false, "odd number": true},
		}

		assertEqualTableResults(t, res, exp)
	})
}

// Tested funcs

func evenSingle(a1 int) bool {
	return a1&1 == 0
}

func evenMultipleOut(a1 int) (string, interface{}, bool, int) {
	return "", struct{}{}, evenSingle(a1), -1
}

func evenMultipleIn(a0 []byte, a1 int, a2 map[rune][][]float64) bool {
	panicOnUnexpectedArgs(a0, a2)
	return evenSingle(a1)
}

func evenMultipleInOut(a0 []byte, a1 int, a2 map[rune][][]float64) (string, interface{}, bool, int) {
	panicOnUnexpectedArgs(a0, a2)
	return evenMultipleOut(a1)
}

func double(n int) int {
	return 2 * n
}

// Helpers

func panicOnUnexpectedArgs(a0 []byte, a2 map[rune][][]float64) {
	deq := reflect.DeepEqual
	if !deq(a0, expFixedArgs["a0"]) || !deq(a2, expFixedArgs["a2"]) {
		log.Panicf(
			"received unexpected args:\na0: %#v\nexp0: %#v\na2: %#v\nexp2: %#v",
			a0, expFixedArgs["a0"], a2, expFixedArgs["a2"],
		)
	}
}

type tableResults struct {
	baseResults
	passedAt    map[int]bool
	failedAt    map[int]bool
	passedLabel map[string]bool
	failedLabel map[string]bool
}

func assertEqualTableResults(t *testing.T, res testx.TableResulter, exp tableResults) {
	t.Helper()
	assertEqualBaseResults(t, res, exp.baseResults)
	for i, v := range exp.passedAt {
		if got := res.PassedAt(i); got != v {
			failBadResults(t, "PassedAt", got, v)
		}
	}
	for i, v := range exp.failedAt {
		if got := res.FailedAt(i); got != v {
			failBadResults(t, "FailedAt", got, v)
		}
	}
	for k, v := range exp.passedLabel {
		if got := res.PassedLabel(k); got != v {
			failBadResults(t, "PassedLabel", got, v)
		}
	}
	for k, v := range exp.failedLabel {
		if got := res.FailedLabel(k); got != v {
			failBadResults(t, "FailedLabel", got, v)
		}
	}
}
