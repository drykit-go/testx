package testix_test

import (
	"log"
	"reflect"
	"testing"

	"github.com/drykit-go/testix"
)

var expFixedArgs = map[string]interface{}{
	"a0": []byte("arg0"),
	"a2": map[rune][][]float64{'π': {[]float64{3.14}}},
}

// TestTableRunner ensures testix.Table behaves correctly, in particular
// when dealing with functions with multiple inputs and outputs.
func TestTableRunner(t *testing.T) {
	cases := []testix.Case{
		{In: 42, Exp: true},
		{In: 99, Exp: false},
	}

	const (
		inPos  = 1
		outPos = 2
	)

	a0, a2 := expFixedArgs["a0"], expFixedArgs["a2"]

	t.Run("single in single out", func(t *testing.T) {
		testix.Table(evenSingle, nil).Cases(cases).Run(t)
	})

	t.Run("single in multiple out", func(t *testing.T) {
		testix.Table(evenMultipleOut, &testix.TableConfig{
			OutPos: outPos,
		}).Cases(cases).Run(t)
	})

	t.Run("multiple in single out", func(t *testing.T) {
		testix.Table(evenMultipleIn, &testix.TableConfig{
			InPos: inPos,
			// should accept len(FixedArgs) == nparams-1
			FixedArgs: []interface{}{a0, a2},
		}).Cases(cases).Run(t)
	})

	t.Run("multiple in multiple out", func(t *testing.T) {
		testix.Table(evenMultipleInOut, &testix.TableConfig{
			InPos:  inPos,
			OutPos: outPos,
			// should accept len(FixedArgs) == nparams
			FixedArgs: []interface{}{0: a0, 2: a2},
		}).Cases(cases).Run(t)
	})
}

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

func panicOnUnexpectedArgs(a0 []byte, a2 map[rune][][]float64) {
	if !deq(a0, expFixedArgs["a0"]) || !deq(a2, expFixedArgs["a2"]) {
		log.Fatalf(
			"received unexpected args:\na0: %#v\nexp0: %#v\na2: %#v\nexp2: %#v",
			a0, expFixedArgs["a0"], a2, expFixedArgs["a2"],
		)
	}
}

func deq(a, b interface{}) bool {
	return reflect.DeepEqual(a, b)
}
