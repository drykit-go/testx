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

func (intCheckFactory) OutRange(lo, hi int) IntChecker {
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

func (intCheckFactory) Is(tar int) IntChecker {
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

func (intCheckFactory) Not(tar int) IntChecker {
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

func (intCheckFactory) GT(tar int) IntChecker {
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

func (intCheckFactory) GTE(tar int) IntChecker {
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

func (intCheckFactory) LT(tar int) IntChecker {
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

func (intCheckFactory) LTE(tar int) IntChecker {
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
