package check_test

import (
	"fmt"

	"github.com/drykit-go/testx"
)

/*
	Example: implementation of a custom checker of unknown type
	having no match in interfaces declared by package check.
*/

type MyType struct {
	ID   int
	Name string
}

// MyTypeValidityChecker is a custom checker for type MyType
// that is not handled by package check.
// However it works with any test runner that requires a generic checker
// because it implements Pass(got T) bool and check.Explainer.
type MyTypeValidityChecker struct{}

// Pass do not satisfy any interface declared by check, but has a valid
// signature Pass(got T) bool, thus allowing IsPositiveComplexChecker
// to be recognized as a valid checker.
func (c MyTypeValidityChecker) Pass(got MyType) bool {
	return got.ID >= 0 && len(got.Name) >= 3
}

// Explain satisfies Explainer interface.
func (c MyTypeValidityChecker) Explain(label string, got interface{}) string {
	return fmt.Sprintf("%s: got bad CustomType value: %v", label, got)
}

// identityCustomType returns the input CustomType
func identityCustomType(v MyType) MyType {
	return v
}

func Example_customCheckerUnknownType() {
	checkIsValid := MyTypeValidityChecker{}
	results := testx.Table(identityCustomType, nil).
		Cases([]testx.Case{
			{In: MyType{ID: 0, Name: "yes"}, Exp: checkIsValid}, // pass
			{In: MyType{ID: -1, Name: "no"}, Exp: checkIsValid}, // fail
		}).
		DryRun()

	fmt.Println(results.PassedAt(0))
	fmt.Println(results.PassedAt(1))
	fmt.Println(results.Checks())

	// Output:
	// true
	// false
	// [{passed} {failed : got bad CustomType value: {-1 no}}]
}
