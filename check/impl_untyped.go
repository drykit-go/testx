package check

import "fmt"

type untypedCheckFactory struct{}

func (untypedCheckFactory) Custom(desc string, f UntypedPassFunc) UntypedChecker {
	return untypedCheck{
		passFunc: f,
		explFunc: func(label string, got interface{}) string {
			return fmt.Sprintf(
				"%s: %s, got %v",
				label, desc, got,
			)
		},
	}
}
