package check

import (
	"fmt"
	"time"
)

type durationCheckerFactory struct{}

func (f durationCheckerFactory) Over(tar time.Duration) DurationChecker {
	pass := func(got time.Duration) bool { return f.ms(got) > f.ms(tar) }
	expl := func(label string, got interface{}) string {
		return fmt.Sprintf(
			"expect %s to be over %vms, got %vms",
			label, f.ms(tar), f.ms(got.(time.Duration)),
		)
	}
	return NewDurationChecker(pass, expl)
}

func (f durationCheckerFactory) Under(tar time.Duration) DurationChecker {
	pass := func(got time.Duration) bool { return f.ms(got) < f.ms(tar) }
	expl := func(label string, got interface{}) string {
		return fmt.Sprintf(
			"expect %s to under %vms, got %vms",
			label, f.ms(tar), f.ms(got.(time.Duration)),
		)
	}
	return NewDurationChecker(pass, expl)
}

// Helpers

func (f durationCheckerFactory) ms(d time.Duration) int64 {
	return d.Milliseconds()
}
