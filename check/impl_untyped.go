package check

import "fmt"

type untypedValue struct{}

func (untypedValue) Custom(desc string, f UntypedPassFunc) UntypedChecker {
	return untypedCheck{
		passFunc: f,
		explFunc: func(label string, got interface{}) string {
			return fmt.Sprintf(
				"expect %s to %s, got %v",
				label, desc, got,
			)
		},
	}
}
