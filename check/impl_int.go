package check

import "fmt"

type intCheckerFactory struct{}

func (intCheckerFactory) InRange(lo, hi int) IntChecker {
	return intChecker{
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

func (intCheckerFactory) OutRange(lo, hi int) IntChecker {
	return intChecker{
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

func (intCheckerFactory) Is(tar int) IntChecker {
	return intChecker{
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

func (intCheckerFactory) Not(tar int) IntChecker {
	return intChecker{
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

func (intCheckerFactory) GT(tar int) IntChecker {
	return intChecker{
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

func (intCheckerFactory) GTE(tar int) IntChecker {
	return intChecker{
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

func (intCheckerFactory) LT(tar int) IntChecker {
	return intChecker{
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

func (intCheckerFactory) LTE(tar int) IntChecker {
	return intChecker{
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
