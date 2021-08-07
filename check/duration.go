package check

import (
	"fmt"
	"time"
)

type durationValue struct{}

func (durationValue) Over(tar time.Duration) DurationChecker {
	return durationCheck{
		passFunc: func(got time.Duration) bool {
			return ms(got) > ms(tar)
		},
		explFunc: func(label string, got interface{}) string {
			return fmt.Sprintf(
				"expect %s to be over %vms, got %vms",
				label, ms(tar), ms(got.(time.Duration)),
			)
		},
	}
}

func (durationValue) Under(tar time.Duration) DurationChecker {
	return durationCheck{
		passFunc: func(got time.Duration) bool {
			return ms(got) < ms(tar)
		},
		explFunc: func(label string, got interface{}) string {
			return fmt.Sprintf(
				"expect %s to under %vms, got %vms",
				label, ms(tar), ms(got.(time.Duration)),
			)
		},
	}
}

func ms(d time.Duration) int64 { return d.Milliseconds() }
