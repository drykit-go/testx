package check

import "fmt"

type untypedCheckerFactory struct{}

func (untypedCheckerFactory) Custom(desc string, f UntypedPassFunc) UntypedChecker {
	return untypedChecker{
		passFunc: f,
		explFunc: func(label string, got interface{}) string {
			return fmt.Sprintf(
				"%s: %s, got %v",
				label, desc, got,
			)
		},
	}
}
