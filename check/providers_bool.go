package check

import (
	"fmt"
)

// boolCheckerProvider provides checks on type bool.
type boolCheckerProvider struct{}

// Is checks the gotten bool is equal to the target.
func (boolCheckerProvider) Is(tar bool) BoolChecker {
	pass := func(got bool) bool { return got == tar }
	expl := func(label string, got interface{}) string {
		return fmt.Sprintf(
			"expect %s to be %v, got %v",
			label, tar, got,
		)
	}
	return NewBoolChecker(pass, expl)
}
