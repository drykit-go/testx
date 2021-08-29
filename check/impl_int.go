package check

import "fmt"

type intCheckerFactory struct{}

func (intCheckerFactory) InRange(lo, hi int) IntChecker {
	pass := func(got int) bool { return got >= lo && got <= hi }
	expl := func(label string, got interface{}) string {
		return fmt.Sprintf(
			"expect %s in range [%d:%d], got %d",
			label, lo, hi, got,
		)
	}
	return NewIntChecker(pass, expl)
}

func (intCheckerFactory) OutRange(lo, hi int) IntChecker {
	pass := func(got int) bool { return got < lo || got > hi }
	expl := func(label string, got interface{}) string {
		return fmt.Sprintf(
			"expect %s not in range [%d:%d], got %d",
			label, lo, hi, got,
		)
	}
	return NewIntChecker(pass, expl)
}

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

func (intCheckerFactory) GT(tar int) IntChecker {
	pass := func(got int) bool { return got > tar }
	expl := func(label string, got interface{}) string {
		return fmt.Sprintf(
			"expect %s > %d, got %d",
			label, tar, got,
		)
	}
	return NewIntChecker(pass, expl)
}

func (intCheckerFactory) GTE(tar int) IntChecker {
	pass := func(got int) bool { return got >= tar }
	expl := func(label string, got interface{}) string {
		return fmt.Sprintf(
			"expect %s >= to %d, got %d",
			label, tar, got,
		)
	}
	return NewIntChecker(pass, expl)
}

func (intCheckerFactory) LT(tar int) IntChecker {
	pass := func(got int) bool { return got < tar }
	expl := func(label string, got interface{}) string {
		return fmt.Sprintf(
			"expect %s < %d, got %d",
			label, tar, got,
		)
	}
	return NewIntChecker(pass, expl)
}

func (intCheckerFactory) LTE(tar int) IntChecker {
	pass := func(got int) bool { return got <= tar }
	expl := func(label string, got interface{}) string {
		return fmt.Sprintf(
			"expect %s <= to %d, got %d",
			label, tar, got,
		)
	}
	return NewIntChecker(pass, expl)
}
