package check

import "fmt"

type intCheckerFactory struct{}

func (intCheckerFactory) Is(tar int) IntChecker {
	pass := func(got int) bool { return got == tar }
	expl := func(label string, got interface{}) string {
		return fmt.Sprintf(
			"expect %s == %d, got %d",
			label, tar, got,
		)
	}
	return NewIntChecker(pass, expl)
}

func (intCheckerFactory) Not(tar int) IntChecker {
	pass := func(got int) bool { return got != tar }
	expl := func(label string, got interface{}) string {
		return fmt.Sprintf(
			"expect %s != %d, got %d",
			label, tar, got,
		)
	}
	return NewIntChecker(pass, expl)
}

func (f intCheckerFactory) InRange(lo, hi int) IntChecker {
	pass := func(got int) bool { return f.inrange(got, lo, hi) }
	expl := func(label string, got interface{}) string {
		return fmt.Sprintf(
			"expect %s in range [%d:%d], got %d",
			label, lo, hi, got,
		)
	}
	return NewIntChecker(pass, expl)
}

func (f intCheckerFactory) OutRange(lo, hi int) IntChecker {
	pass := func(got int) bool { return !f.inrange(got, lo, hi) }
	expl := func(label string, got interface{}) string {
		return fmt.Sprintf(
			"expect %s not in range [%d:%d], got %d",
			label, lo, hi, got,
		)
	}
	return NewIntChecker(pass, expl)
}

func (f intCheckerFactory) GT(tar int) IntChecker {
	pass := func(got int) bool { return !f.lte(got, tar) }
	expl := func(label string, got interface{}) string {
		return fmt.Sprintf(
			"expect %s > %d, got %d",
			label, tar, got,
		)
	}
	return NewIntChecker(pass, expl)
}

func (f intCheckerFactory) GTE(tar int) IntChecker {
	pass := func(got int) bool { return !f.lt(got, tar) }
	expl := func(label string, got interface{}) string {
		return fmt.Sprintf(
			"expect %s >= to %d, got %d",
			label, tar, got,
		)
	}
	return NewIntChecker(pass, expl)
}

func (f intCheckerFactory) LT(tar int) IntChecker {
	pass := func(got int) bool { return f.lt(got, tar) }
	expl := func(label string, got interface{}) string {
		return fmt.Sprintf(
			"expect %s < %d, got %d",
			label, tar, got,
		)
	}
	return NewIntChecker(pass, expl)
}

func (f intCheckerFactory) LTE(tar int) IntChecker {
	pass := func(got int) bool { return f.lte(got, tar) }
	expl := func(label string, got interface{}) string {
		return fmt.Sprintf(
			"expect %s <= to %d, got %d",
			label, tar, got,
		)
	}
	return NewIntChecker(pass, expl)
}

// Helpers

func (f intCheckerFactory) lt(a, b int) bool  { return a < b }
func (f intCheckerFactory) lte(a, b int) bool { return a <= b }
func (f intCheckerFactory) inrange(n, lo, hi int) bool {
	return !f.lt(n, lo) && f.lte(n, hi)
}
