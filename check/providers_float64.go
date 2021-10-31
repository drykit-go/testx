package check

import "fmt"

// float64CheckerProvider provides checks on type float64.
type float64CheckerProvider struct{ baseCheckerProvider }

// Is checks the gotten float64 is equal to the target.
func (p float64CheckerProvider) Is(tar float64) Checker[float64] {
	pass := func(got float64) bool { return got == tar }
	expl := func(label string, got any) string {
		return p.explain(label, tar, got)
	}
	return NewChecker(pass, expl)
}

// Not checks the gotten float64 is not equal to the target.
func (p float64CheckerProvider) Not(values ...float64) Checker[float64] {
	var match float64
	pass := func(got float64) bool {
		for _, v := range values {
			if got == v {
				match = v
				return false
			}
		}
		return true
	}
	expl := func(label string, got any) string {
		return p.explainNot(label, match, got)
	}
	return NewChecker(pass, expl)
}

// InRange checks the gotten float64 is in the closed interval [lo:hi].
func (p float64CheckerProvider) InRange(lo, hi float64) Checker[float64] {
	pass := func(got float64) bool { return p.inrange(got, lo, hi) }
	expl := func(label string, got any) string {
		return p.explain(label, fmt.Sprintf("in range [%v:%v]", lo, hi), got)
	}
	return NewChecker(pass, expl)
}

// OutRange checks the gotten float64 is not in the closed interval [lo:hi].
func (p float64CheckerProvider) OutRange(lo, hi float64) Checker[float64] {
	pass := func(got float64) bool { return !p.inrange(got, lo, hi) }
	expl := func(label string, got any) string {
		return p.explainNot(label, fmt.Sprintf("in range [%v:%v]", lo, hi), got)
	}
	return NewChecker(pass, expl)
}

// GT checks the gotten float64 is greater than the target.
func (p float64CheckerProvider) GT(tar float64) Checker[float64] {
	pass := func(got float64) bool { return !p.lte(got, tar) }
	expl := func(label string, got any) string {
		return p.explain(label, fmt.Sprintf("> %v", tar), got)
	}
	return NewChecker(pass, expl)
}

// GTE checks the gotten float64 is greater or equal to the target.
func (p float64CheckerProvider) GTE(tar float64) Checker[float64] {
	pass := func(got float64) bool { return !p.lt(got, tar) }
	expl := func(label string, got any) string {
		return p.explain(label, fmt.Sprintf(">= %v", tar), got)
	}
	return NewChecker(pass, expl)
}

// LT checks the gotten float64 is lesser than the target.
func (p float64CheckerProvider) LT(tar float64) Checker[float64] {
	pass := func(got float64) bool { return p.lt(got, tar) }
	expl := func(label string, got any) string {
		return p.explain(label, fmt.Sprintf("< %v", tar), got)
	}
	return NewChecker(pass, expl)
}

// LTE checks the gotten float64 is lesser or equal to the target.
func (p float64CheckerProvider) LTE(tar float64) Checker[float64] {
	pass := func(got float64) bool { return p.lte(got, tar) }
	expl := func(label string, got any) string {
		return p.explain(label, fmt.Sprintf("<= %v", tar), got)
	}
	return NewChecker(pass, expl)
}

// Helpers

func (float64CheckerProvider) lt(a, b float64) bool  { return a < b }
func (float64CheckerProvider) lte(a, b float64) bool { return a <= b }
func (p float64CheckerProvider) inrange(n, lo, hi float64) bool {
	return !p.lt(n, lo) && p.lte(n, hi)
}
