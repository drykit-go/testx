package check

import "fmt"

// valueCheckerProvider provides checks on type interface{}.
type valueCheckerProvider struct{}

// Custom checks the gotten value passes the given ValuePassFunc.
// The description should typically begin with keywords like "expect"
// or "should" for intelligible output.
// For instance, "expect odd number" would output:
// 	> "expect odd number, got 42"
func (valueCheckerProvider) Custom(desc string, f ValuePassFunc) ValueChecker {
	expl := func(label string, got interface{}) string {
		return fmt.Sprintf(
			"%s: %s, got %v",
			label, desc, got,
		)
	}
	return NewValueChecker(f, expl)
}
