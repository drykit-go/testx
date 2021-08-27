package check

import "fmt"

type valueCheckerFactory struct{}

func (valueCheckerFactory) Custom(desc string, f ValuePassFunc) ValueChecker {
	return valueChecker{
		passFunc: f,
		explFunc: func(label string, got interface{}) string {
			return fmt.Sprintf(
				"%s: %s, got %v",
				label, desc, got,
			)
		},
	}
}
