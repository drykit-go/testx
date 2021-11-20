package providers

import (
	"fmt"
	"time"

	check "github.com/drykit-go/testx/internal/checktypes"
)

// DurationCheckerProvider provides checks on type time.Duration.
type DurationCheckerProvider struct{ baseCheckerProvider }

// Over checks the gotten time.Duration is over the target duration.
func (p DurationCheckerProvider) Over(tar time.Duration) check.Checker[time.Duration] {
	pass := func(got time.Duration) bool { return p.ns(got) > p.ns(tar) }
	expl := func(label string, got any) string {
		return p.explain(label,
			fmt.Sprintf("over %vms", p.ms(tar)),
			fmt.Sprintf("%vms", p.ms(got.(time.Duration))),
		)
	}
	return check.NewChecker(pass, expl)
}

// Under checks the gotten time.Duration is under the target duration.
func (p DurationCheckerProvider) Under(tar time.Duration) check.Checker[time.Duration] {
	pass := func(got time.Duration) bool { return p.ns(got) < p.ns(tar) }
	expl := func(label string, got any) string {
		return p.explain(label,
			fmt.Sprintf("under %vms", p.ms(tar)),
			fmt.Sprintf("%vms", p.ms(got.(time.Duration))),
		)
	}
	return check.NewChecker(pass, expl)
}

// InRange checks the gotten time.Duration is in range [lo:hi]
func (p DurationCheckerProvider) InRange(lo, hi time.Duration) check.Checker[time.Duration] {
	pass := func(got time.Duration) bool { return p.inrange(got, lo, hi) }
	expl := func(label string, got any) string {
		return p.explain(label,
			fmt.Sprintf("in range [%vms:%vms]", p.ms(lo), p.ms(hi)),
			fmt.Sprintf("%vms", p.ms(got.(time.Duration))),
		)
	}
	return check.NewChecker(pass, expl)
}

// OutRange checks the gotten time.Duration is not in range [lo:hi]
func (p DurationCheckerProvider) OutRange(lo, hi time.Duration) check.Checker[time.Duration] {
	pass := func(got time.Duration) bool { return !p.inrange(got, lo, hi) }
	expl := func(label string, got any) string {
		return p.explain(label,
			fmt.Sprintf("not in range [%vms:%vms]", p.ms(lo), p.ms(hi)),
			fmt.Sprintf("%vms", p.ms(got.(time.Duration))),
		)
	}
	return check.NewChecker(pass, expl)
}

// Helpers

func (p DurationCheckerProvider) inrange(d, lo, hi time.Duration) bool {
	nsd, nslo, nshi := p.ns(d), p.ns(lo), p.ns(hi)
	return nslo <= nsd && nsd <= nshi
}

func (DurationCheckerProvider) ms(d time.Duration) int64 {
	return d.Milliseconds()
}

func (DurationCheckerProvider) ns(d time.Duration) int64 {
	return d.Nanoseconds()
}
