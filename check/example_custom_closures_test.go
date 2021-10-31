package check_test

import (
	"fmt"

	"github.com/drykit-go/testx"
)

/*
	Example: implementation of a custom checker using closures.
	This example introduces a pattern that has 2 benefits:
		- parameterized checkers (e.g. check.Int.InRange(0, 10))
		- communication between Pass and Explain funcs

	This pattern is used extensively for all providers implementations
	in package check.
*/

// Complex128Checker is a custom checker. Its behavior can be set dynamically
// because its methods Pass and Explain return the result of settable
// func fields pass and explain respectively.
type Complex128Checker struct {
	pass    func(got complex128) bool
	explain func(label string, got interface{}) string
}

// Pass do not satisfy any interface declared by check, but has a valid
// signature Pass(got T) bool, allowing Complex128Checker to be casted
// as a Checker[any] using checkconv.Cast.
func (c Complex128Checker) Pass(got complex128) bool {
	return c.pass(got)
}

// Explain satisfies check.Explainer interface.
func (c Complex128Checker) Explain(label string, got any) string {
	return c.explain(label, got)
}

// checkComplex128ImagIsInRange returns a new Complex128Checker that checks
// the imaginary part of a complex128 is in range [lo:hi].
func checkComplex128ImagIsInRange(lo, hi float64) Complex128Checker {
	var gotImag float64 // gotImag is accessible by both pass and expl
	pass := func(got complex128) bool {
		gotImag = imag(got)
		return lo <= gotImag && gotImag <= hi
	}
	expl := func(label string, got interface{}) string {
		return fmt.Sprintf(
			"%s: exp imag(%v) in range [%.0f:%.0f], got %.0f",
			label, got, lo, hi, gotImag,
		)
	}
	return Complex128Checker{pass: pass, explain: expl}
}

func Example_customCheckerClosures() {
	const testedValue = 42i

	results := testx.Value(testedValue).Pass(
		checkComplex128ImagIsInRange(41, 43), // pass
		checkComplex128ImagIsInRange(0, 1),   // fail
	).DryRun()

	fmt.Println(results.Checks())

	// Output:
	// [{passed} {failed value: exp imag((0+42i)) in range [0:1], got 42}]
}
