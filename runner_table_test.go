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
		{In: 99, Exp: false},
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
		testx.
			Table(evenMultipleOut, &testx.TableConfig{
				OutPos: outPos,
			}).
			Cases(cases).
			Run(t)
	})

	t.Run("multiple in single out", func(t *testing.T) {
		testx.
			Table(evenMultipleIn, &testx.TableConfig{
				InPos:     inPos,
				FixedArgs: []interface{}{a0, a2}, // len(FixedArgs) == nparams-1
			}).
			Cases(cases).
			Run(t)
	})

	t.Run("multiple in multiple out", func(t *testing.T) {
		testx.
			Table(evenMultipleInOut, &testx.TableConfig{
				InPos:     inPos,
				OutPos:    outPos,
				FixedArgs: []interface{}{0: a0, 2: a2}, // len(FixedArgs) == nparams
			}).
			Cases(cases).
			Run(t)
	})

	t.Run("using check.IntChecker", func(t *testing.T) {
		testx.
			Table(double, nil).
			Cases([]testx.Case{
				{In: 21, Exp: check.Int.Equal(42)},
				{In: -4, Exp: check.Int.InRange(-10, 0)},
			}).
			Run(t)
	})
}

// Tested funcs

func evenSingle(a1 int) bool {
	return a1&1 == 0
}

func evenMultipleOut(a1 int) (string, interface{}, bool, int) { //nolint: gocritic,unamedResult
	return "", struct{}{}, evenSingle(a1), -1
}

func evenMultipleIn(a0 []byte, a1 int, a2 map[rune][][]float64) bool {
	panicOnUnexpectedArgs(a0, a2)
	return evenSingle(a1)
}

func evenMultipleInOut(a0 []byte, a1 int, a2 map[rune][][]float64) (string, interface{}, bool, int) { //nolint: gocritic,unamedResult
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
		log.Fatalf(
			"received unexpected args:\na0: %#v\nexp0: %#v\na2: %#v\nexp2: %#v",
			a0, expFixedArgs["a0"], a2, expFixedArgs["a2"],
		)
	}
}
