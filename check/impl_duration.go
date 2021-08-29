package check

import (
	"fmt"
	"time"
)

type durationCheckerFactory struct{}

func (durationCheckerFactory) Over(tar time.Duration) DurationChecker {
	pass := func(got time.Duration) bool { return ms(got) > ms(tar) }
	expl := func(label string, got interface{}) string {
		return fmt.Sprintf(
			"expect %s to be over %vms, got %vms",
			label, ms(tar), ms(got.(time.Duration)),
		)
	}
	return NewDurationChecker(pass, expl)
}

func (durationCheckerFactory) Under(tar time.Duration) DurationChecker {
	pass := func(got time.Duration) bool { return ms(got) < ms(tar) }
	expl := func(label string, got interface{}) string {
		return fmt.Sprintf(
			"expect %s to under %vms, got %vms",
			label, ms(tar), ms(got.(time.Duration)),
		)
	}
	return NewDurationChecker(pass, expl)
}

func ms(d time.Duration) int64 { return d.Milliseconds() }
