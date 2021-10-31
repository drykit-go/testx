package check

import (
	"fmt"
	"time"
)

// durationCheckerProvider provides checks on type time.Duration.
type durationCheckerProvider struct{ baseCheckerProvider }

// Over checks the gotten time.Duration is over the target duration.
func (p durationCheckerProvider) Over(tar time.Duration) Checker[time.Duration] {
	pass := func(got time.Duration) bool { return p.ns(got) > p.ns(tar) }
	expl := func(label string, got any) string {
		return p.explain(label,
			fmt.Sprintf("over %vms", p.ms(tar)),
			fmt.Sprintf("%vms", p.ms(got.(time.Duration))),
		)
	}
	return NewChecker(pass, expl)
}

// Under checks the gotten time.Duration is under the target duration.
func (p durationCheckerProvider) Under(tar time.Duration) Checker[time.Duration] {
	pass := func(got time.Duration) bool { return p.ns(got) < p.ns(tar) }
	expl := func(label string, got any) string {
		return p.explain(label,
			fmt.Sprintf("under %vms", p.ms(tar)),
			fmt.Sprintf("%vms", p.ms(got.(time.Duration))),
		)
	}
	return NewChecker(pass, expl)
}

// InRange checks the gotten time.Duration is in range [lo:hi]
func (p durationCheckerProvider) InRange(lo, hi time.Duration) Checker[time.Duration] {
	pass := func(got time.Duration) bool { return p.inrange(got, lo, hi) }
	expl := func(label string, got any) string {
		return p.explain(label,
			fmt.Sprintf("in range [%vms:%vms]", p.ms(lo), p.ms(hi)),
			fmt.Sprintf("%vms", p.ms(got.(time.Duration))),
		)
	}
	return NewChecker(pass, expl)
}

// OutRange checks the gotten time.Duration is not in range [lo:hi]
func (p durationCheckerProvider) OutRange(lo, hi time.Duration) Checker[time.Duration] {
	pass := func(got time.Duration) bool { return !p.inrange(got, lo, hi) }
	expl := func(label string, got any) string {
		return p.explain(label,
			fmt.Sprintf("not in range [%vms:%vms]", p.ms(lo), p.ms(hi)),
			fmt.Sprintf("%vms", p.ms(got.(time.Duration))),
		)
	}
	return NewChecker(pass, expl)
}

// Helpers

func (p durationCheckerProvider) inrange(d, lo, hi time.Duration) bool {
	nsd, nslo, nshi := p.ns(d), p.ns(lo), p.ns(hi)
	return nslo <= nsd && nsd <= nshi
}

func (durationCheckerProvider) ms(d time.Duration) int64 {
	return d.Milliseconds()
}

func (durationCheckerProvider) ns(d time.Duration) int64 {
	return d.Nanoseconds()
}
