package testx_test

import (
	"errors"
	"log"
	"reflect"
	"strings"
	"testing"

	"github.com/drykit-go/testx"
	"github.com/drykit-go/testx/check"
	"github.com/drykit-go/testx/checkconv"
)

// Tests

var expFixedArgs = map[string]any{
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
		testx.Table(evenSingle).Cases(cases).Run(t)
	})

	t.Run("single in multiple out", func(t *testing.T) {
		testx.Table(evenMultipleOut).Config(testx.TableConfig{
			OutPos: outPos,
		}).
			Cases(cases).
			Run(t)
	})

	t.Run("multiple in single out", func(t *testing.T) {
		testx.Table(evenMultipleIn).Config(testx.TableConfig{
			InPos:     inPos,
			FixedArgs: []any{a0, a2}, // len(FixedArgs) == nparams-1
		}).
			Cases(cases).
			Run(t)
	})

	t.Run("multiple in multiple out", func(t *testing.T) {
		testx.Table(evenMultipleInOut).Config(testx.TableConfig{
			InPos:     inPos,
			OutPos:    outPos,
			FixedArgs: []any{0: a0, 2: a2}, // len(FixedArgs) == nparams
		}).
			Cases(cases).
			Run(t)
	})

	t.Run("using check.IntChecker", func(t *testing.T) {
		testx.Table(double).
			Cases([]testx.Case{
				{In: 21, Pass: checkconv.AssertMany(check.Int.Is(42))},
				{In: -4, Pass: checkconv.AssertMany(check.Int.InRange(-10, 0))},
			}).
			Run(t)
	})

	t.Run("expect nil value", func(t *testing.T) {
		runner := testx.Table(func(wantnil bool) any {
			if wantnil {
				return nil
			}
			return 0
		}).Cases([]testx.Case{
			{In: false, Exp: 0},
			{In: true},                    // Exp == nil, no check added
			{In: true, Exp: testx.ExpNil}, // expect nil value
		})

		runner.Run(t)

		if n := runner.DryRun().NChecks(); n != 2 {
			t.Errorf("exp 2 checks, got %d", n)
		}
	})

	t.Run("Case.Not checks", func(t *testing.T) {
		results := testx.Table(func(n int) int { return n }).
			Cases([]testx.Case{
				{In: 0, Not: []any{-1, 1}}, // pass
				{In: 0, Not: []any{0}},     // fail
			}).
			DryRun()

		if nc := results.NChecks(); nc != 2 {
			t.Errorf("exp 2 checks, got %d", nc)
		}
		if results.FailedAt(0) {
			t.Error("exp Case 0 to pass, got fail")
		}
		if results.PassedAt(1) {
			t.Error("exp Case 1 to fail, got pass")
		}
	})

	t.Run("test case labels", func(t *testing.T) {
		results := testx.Table(divide).Config(testx.TableConfig{
			InPos:     1,
			OutPos:    1,
			FixedArgs: testx.Args{42.0},
		}).Cases([]testx.Case{
			{In: 0.0, Exp: testx.ExpNil, Lab: "zeroth case"}, // fail
			{In: 0.0, Exp: testx.ExpNil, Lab: "first case"},  // fail
		}).DryRun()

		expLabels := []string{
			`Table.Cases[0] "zeroth case" testx_test.divide(42, 0)`,
			`Table.Cases[1] "first case" testx_test.divide(42, 0)`,
		}

		for i, c := range results.Checks() {
			got := c.Reason
			exp := expLabels[i]
			if !strings.HasPrefix(got, exp) {
				t.Errorf("bad label output\nexp %s\ngot %s", got, exp)
			}
		}
	})
}

func TestExpNil(t *testing.T) {
	t.Run("Exp=ExpNil expects nil", func(t *testing.T) {
		f := func(int) any { return nil }
		res := testx.Table(f).Cases([]testx.Case{
			{In: 0, Exp: testx.ExpNil},
		}).DryRun()

		if n := res.NChecks(); n != 1 {
			t.Errorf("exp 1 check, got %d", n)
		}
		if res.Failed() {
			t.Error("nil did not pass Case.Exp == ExpNil")
		}
	})

	t.Run("Exp=ExpNil does not expect 0", func(t *testing.T) {
		f := func(int) int { return 0 }
		res := testx.Table(f).Cases([]testx.Case{
			{In: 0, Exp: testx.ExpNil},
		}).DryRun()

		if n := res.NChecks(); n != 1 {
			t.Errorf("exp 1 check, got %d", n)
		}
		if res.Passed() {
			t.Error("0 did pass Case.Exp == ExpNil")
		}
	})

	t.Run("Exp=0 does not expect nil", func(t *testing.T) {
		f := func(int) any { return nil }
		res := testx.Table(f).Cases([]testx.Case{
			{In: 0, Exp: 0},
		}).DryRun()

		if n := res.NChecks(); n != 1 {
			t.Errorf("exp 1 check, got %d", n)
		}
		if res.Passed() {
			t.Error("nil did pass Case.Exp == 0")
		}
	})
}

func TestTableRunnerResults(t *testing.T) {
	t.Run("pass", func(t *testing.T) {
		res := testx.
			Table(evenSingle).
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

	t.Run("fail", func(t *testing.T) {
		res := testx.
			Table(evenSingle).
			Cases([]testx.Case{
				{In: 10, Exp: true, Lab: "even number"}, // pass
				{In: -1, Exp: true, Lab: "odd number"},  // fail
				{In: -1, Exp: true},                     // fail
			}).
			DryRun()

		exp := tableResults{
			baseResults: baseResults{
				passed:  false,
				failed:  true,
				nPassed: 1,
				nFailed: 2,
				nChecks: 3,
				checks: []testx.CheckResult{
					{Passed: true, Reason: ""},
					{Passed: false, Reason: "Table.Cases[1] \"odd number\" testx_test.evenSingle(-1):\nexp true\ngot false"},
					{Passed: false, Reason: "Table.Cases[2] testx_test.evenSingle(-1):\nexp true\ngot false"},
				},
			},
			passedAt:    map[int]bool{0: true, 1: false, 2: false},
			failedAt:    map[int]bool{0: false, 1: true, 2: true},
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

func evenMultipleOut(a1 int) (string, any, bool, int) {
	return "", struct{}{}, evenSingle(a1), -1
}

func evenMultipleIn(a0 []byte, a1 int, a2 map[rune][][]float64) bool {
	panicOnUnexpectedArgs(a0, a2)
	return evenSingle(a1)
}

func evenMultipleInOut(a0 []byte, a1 int, a2 map[rune][][]float64) (string, any, bool, int) {
	panicOnUnexpectedArgs(a0, a2)
	return evenMultipleOut(a1)
}

func double(n int) int {
	return 2 * n
}

func divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, errors.New("division by 0")
	}
	return a / b, nil
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
