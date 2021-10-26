package check

import "fmt"

// intCheckerProvider provides checks on type int.
type intCheckerProvider struct{ baseCheckerProvider }

// Is checks the gotten int is equal to the target.
func (p intCheckerProvider) Is(tar int) IntChecker {
	pass := func(got int) bool { return got == tar }
	expl := func(label string, got interface{}) string {
		return p.explain(label, tar, got)
	}
	return NewIntChecker(pass, expl)
}

// Not checks the gotten int is not equal to the target.
func (p intCheckerProvider) Not(values ...int) IntChecker {
	var match int
	pass := func(got int) bool {
		for _, v := range values {
			if got == v {
				match = v
				return false
			}
		}
		return true
	}
	expl := func(label string, got interface{}) string {
		return p.explainNot(label, match, got)
	}
	return NewIntChecker(pass, expl)
}

// InRange checks the gotten int is in the closed interval [lo:hi].
func (p intCheckerProvider) InRange(lo, hi int) IntChecker {
	pass := func(got int) bool { return p.inrange(got, lo, hi) }
	expl := func(label string, got interface{}) string {
		return p.explain(label, fmt.Sprintf("in range [%v:%v]", lo, hi), got)
	}
	return NewIntChecker(pass, expl)
}

// OutRange checks the gotten int is not in the closed interval [lo:hi].
func (p intCheckerProvider) OutRange(lo, hi int) IntChecker {
	pass := func(got int) bool { return !p.inrange(got, lo, hi) }
	expl := func(label string, got interface{}) string {
		return p.explainNot(label, fmt.Sprintf("in range [%v:%v]", lo, hi), got)
	}
	return NewIntChecker(pass, expl)
}

// GT checks the gotten int is greater than the target.
func (p intCheckerProvider) GT(tar int) IntChecker {
	pass := func(got int) bool { return !p.lte(got, tar) }
	expl := func(label string, got interface{}) string {
		return p.explain(label, fmt.Sprintf("> %v", tar), got)
	}
	return NewIntChecker(pass, expl)
}

// GTE checks the gotten int is greater or equal to the target.
func (p intCheckerProvider) GTE(tar int) IntChecker {
	pass := func(got int) bool { return !p.lt(got, tar) }
	expl := func(label string, got interface{}) string {
		return p.explain(label, fmt.Sprintf(">= %v", tar), got)
	}
	return NewIntChecker(pass, expl)
}

// LT checks the gotten int is lesser than the target.
func (p intCheckerProvider) LT(tar int) IntChecker {
	pass := func(got int) bool { return p.lt(got, tar) }
	expl := func(label string, got interface{}) string {
		return p.explain(label, fmt.Sprintf("< %v", tar), got)
	}
	return NewIntChecker(pass, expl)
}

// LTE checks the gotten int is lesser or equal to the target.
func (p intCheckerProvider) LTE(tar int) IntChecker {
	pass := func(got int) bool { return p.lte(got, tar) }
	expl := func(label string, got interface{}) string {
		return p.explain(label, fmt.Sprintf("<= %v", tar), got)
	}
	return NewIntChecker(pass, expl)
}

// Helpers

func (intCheckerProvider) lt(a, b int) bool  { return a < b }
func (intCheckerProvider) lte(a, b int) bool { return a <= b }
func (p intCheckerProvider) inrange(n, lo, hi int) bool {
	return !p.lt(n, lo) && p.lte(n, hi)
}
