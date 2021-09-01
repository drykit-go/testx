package check

import (
	"fmt"
	"time"
)

// durationCheckerProvider provides checks on type time.Duration.
type durationCheckerProvider struct{}

// Over checks the gotten time.Duration is over the target duration.
func (f durationCheckerProvider) Over(tar time.Duration) DurationChecker {
	pass := func(got time.Duration) bool { return f.ms(got) > f.ms(tar) }
	expl := func(label string, got interface{}) string {
		return fmt.Sprintf(
			"expect %s to be over %vms, got %vms",
			label, f.ms(tar), f.ms(got.(time.Duration)),
		)
	}
	return NewDurationChecker(pass, expl)
}

// Under checks the gotten time.Duration is under the target duration.
func (f durationCheckerProvider) Under(tar time.Duration) DurationChecker {
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

func (f durationCheckerProvider) ms(d time.Duration) int64 {
	return d.Milliseconds()
}
