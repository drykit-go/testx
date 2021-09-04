package check

import (
	"fmt"
	"time"
)

// durationCheckerProvider provides checks on type time.Duration.
type durationCheckerProvider struct{}

// Over checks the gotten time.Duration is over the target duration.
func (p durationCheckerProvider) Over(tar time.Duration) DurationChecker {
	pass := func(got time.Duration) bool { return p.ns(got) > p.ns(tar) }
	expl := func(label string, got interface{}) string {
		return fmt.Sprintf(
			"expect %s to be over %vms, got %vms",
			label, p.ns(tar), p.ns(got.(time.Duration)),
		)
	}
	return NewDurationChecker(pass, expl)
}

// Under checks the gotten time.Duration is under the target duration.
func (p durationCheckerProvider) Under(tar time.Duration) DurationChecker {
	pass := func(got time.Duration) bool { return p.ns(got) < p.ns(tar) }
	expl := func(label string, got interface{}) string {
		return fmt.Sprintf(
			"expect %s to under %vms, got %vms",
			label, p.ns(tar), p.ns(got.(time.Duration)),
		)
	}
	return NewDurationChecker(pass, expl)
}

// InRange checks the gotten time.Duration is in range [lo:hi]
func (p durationCheckerProvider) InRange(lo, hi time.Duration) DurationChecker {
	pass := func(got time.Duration) bool { return p.inrange(got, lo, hi) }
	expl := func(label string, got interface{}) string {
		return fmt.Sprintf(
			"%s:\nexp in range [%vms:%vms]\ngot %vms",
			label, p.ms(lo), p.ms(hi), p.ms(got.(time.Duration)),
		)
	}
	return NewDurationChecker(pass, expl)
}

// OutRange checks the gotten time.Duration is not in range [lo:hi]
func (p durationCheckerProvider) OutRange(lo, hi time.Duration) DurationChecker {
	pass := func(got time.Duration) bool { return !p.inrange(got, lo, hi) }
	expl := func(label string, got interface{}) string {
		return fmt.Sprintf(
			"%s:\nexp not in range [%vms:%vms]\ngot %vms",
			label, p.ms(lo), p.ms(hi), p.ms(got.(time.Duration)),
		)
	}
	return NewDurationChecker(pass, expl)
}

// Helpers

func (p durationCheckerProvider) inrange(d, lo, hi time.Duration) bool {
	nsd, nslo, nshi := p.ns(d), p.ns(lo), p.ns(hi)
	return nslo < nsd && nsd <= nshi
}

func (durationCheckerProvider) ms(d time.Duration) int64 {
	return d.Milliseconds()
}

func (durationCheckerProvider) ns(d time.Duration) int64 {
	return d.Nanoseconds()
}
