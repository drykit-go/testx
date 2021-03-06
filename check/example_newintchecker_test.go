package check_test

import (
	"fmt"

	"github.com/drykit-go/testx"
	"github.com/drykit-go/testx/check"
	"github.com/drykit-go/testx/checkconv"
)

/*
	Example: implementation of a custom checker of a type
	defined by package check
*/

func Example_newIntChecker() {
	checkIsEven := check.NewIntChecker(
		func(got int) bool { return got&1 == 0 },
		func(label string, got interface{}) string {
			return fmt.Sprintf("%s: expect even int, got %v", label, got)
		},
	)

	resultPass := testx.Value(42).Pass(checkconv.FromInt(checkIsEven)).DryRun()
	resultFail := testx.Value(43).Pass(checkconv.FromInt(checkIsEven)).DryRun()

	fmt.Println(resultPass.Passed())
	fmt.Println(resultFail.Passed())

	// Output:
	// true
	// false
}
