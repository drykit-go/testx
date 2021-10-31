package testx_test

import (
	"fmt"
	"testing"

	"github.com/drykit-go/testx"
	"github.com/drykit-go/testx/check"
)

/*
	ValueRunner
*/

func ExampleValueRunner() {
	t := &testing.T{}                 // ignore: emulating a testing context
	get42 := func() int { return 42 } // Some dummy func

	// Run Value test via Run(t *testing.T)
	testx.Value(get42()).
		Exp(42).                         // pass
		Not(3, -1).                      // pass
		Pass(check.Int.InRange(41, 43)). // pass
		Run(t)
}

func ExampleValueRunner_dryRun() {
	get42 := func() int { return 42 } // Some dummy func

	results := testx.Value(get42()).
		Exp(42). // pass
		Pass(
			check.Int.GTE(41),        // pass
			check.Int.InRange(-1, 1), // fail
		).
		DryRun()

	fmt.Println(results.Passed())
	fmt.Println(results.Failed())
	fmt.Println(results.Checks())
	fmt.Println(results.NPassed())
	fmt.Println(results.NFailed())
	fmt.Println(results.NChecks())

	// Output:
	// false
	// true
	// [{passed} {passed} {failed value:
	// exp in range [-1:1]
	// got 42}]
	// 2
	// 1
	// 3
}
