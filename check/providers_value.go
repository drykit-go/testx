package check

import "fmt"

type valueCheckerProvider struct{}

func (valueCheckerProvider) Custom(desc string, f ValuePassFunc) ValueChecker {
	expl := func(label string, got interface{}) string {
		return fmt.Sprintf(
			"%s: %s, got %v",
			label, desc, got,
		)
	}
	return NewValueChecker(f, expl)
}
