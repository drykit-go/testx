package check

import "fmt"

type intCheckFactory struct{}

func (intCheckFactory) InRange(lo, hi int) IntChecker {
	return intCheck{
		passFunc: func(got int) bool {
			return got >= lo && got <= hi
		},
		explFunc: func(label string, got interface{}) string {
			return fmt.Sprintf(
				"expect %s in range [%d:%d], got %d",
				label, lo, hi, got,
			)
		},
	}
}

func (intCheckFactory) NotInRange(lo, hi int) IntChecker {
	return intCheck{
		passFunc: func(got int) bool {
			return got < lo || got > hi
		},
		explFunc: func(label string, got interface{}) string {
			return fmt.Sprintf(
				"expect %s not in range [%d:%d], got %d",
				label, lo, hi, got,
			)
		},
	}
}

func (intCheckFactory) Equal(tar int) IntChecker {
	return intCheck{
		passFunc: func(got int) bool {
			return got == tar
		},
		explFunc: func(label string, got interface{}) string {
			return fmt.Sprintf(
				"expect %s == %d, got %d",
				label, tar, got,
			)
		},
	}
}

func (intCheckFactory) NotEqual(tar int) IntChecker {
	return intCheck{
		passFunc: func(got int) bool {
			return got != tar
		},
		explFunc: func(label string, got interface{}) string {
			return fmt.Sprintf(
				"expect %s != %d, got %d",
				label, tar, got,
			)
		},
	}
}

func (intCheckFactory) GreaterThan(tar int) IntChecker {
	return intCheck{
		passFunc: func(got int) bool {
			return got > tar
		},
		explFunc: func(label string, got interface{}) string {
			return fmt.Sprintf(
				"expect %s > %d, got %d",
				label, tar, got,
			)
		},
	}
}

func (intCheckFactory) GreaterOrEqual(tar int) IntChecker {
	return intCheck{
		passFunc: func(got int) bool {
			return got >= tar
		},
		explFunc: func(label string, got interface{}) string {
			return fmt.Sprintf(
				"expect %s >= to %d, got %d",
				label, tar, got,
			)
		},
	}
}

func (intCheckFactory) LesserThan(tar int) IntChecker {
	return intCheck{
		passFunc: func(got int) bool {
			return got < tar
		},
		explFunc: func(label string, got interface{}) string {
			return fmt.Sprintf(
				"expect %s < %d, got %d",
				label, tar, got,
			)
		},
	}
}

func (intCheckFactory) LesserOrEqual(tar int) IntChecker {
	return intCheck{
		passFunc: func(got int) bool {
			return got <= tar
		},
		explFunc: func(label string, got interface{}) string {
			return fmt.Sprintf(
				"expect %s <= to %d, got %d",
				label, tar, got,
			)
		},
	}
}
