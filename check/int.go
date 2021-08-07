package check

import "fmt"

type intValue struct{}

func (intValue) InRange(lo, hi int) IntChecker {
	return intCheck{
		passFunc: func(got int) bool {
			return got >= lo && got <= hi
		},
		explFunc: func(label string, got interface{}) string {
			return fmt.Sprintf(
				"expect %s to be in range [%d:%d], got %d",
				label, lo, hi, got,
			)
		},
	}
}

func (intValue) NotInRange(lo, hi int) IntChecker {
	return intCheck{
		passFunc: func(got int) bool {
			return got < lo || got > hi
		},
		explFunc: func(label string, got interface{}) string {
			return fmt.Sprintf(
				"expect %s not to be in range [%d:%d], got %d",
				label, lo, hi, got,
			)
		},
	}
}

func (intValue) Equal(tar int) IntChecker {
	return intCheck{
		passFunc: func(got int) bool {
			return got == tar
		},
		explFunc: func(label string, got interface{}) string {
			return fmt.Sprintf(
				"expect %s to equal %d, got %d",
				label, tar, got,
			)
		},
	}
}

func (intValue) NotEqual(tar int) IntChecker {
	return intCheck{
		passFunc: func(got int) bool {
			return got != tar
		},
		explFunc: func(label string, got interface{}) string {
			return fmt.Sprintf(
				"expect %s not to equal %d, got %d",
				label, tar, got,
			)
		},
	}
}

func (intValue) GreaterThan(tar int) IntChecker {
	return intCheck{
		passFunc: func(got int) bool {
			return got > tar
		},
		explFunc: func(label string, got interface{}) string {
			return fmt.Sprintf(
				"expect %s to be greater than %d, got %d",
				label, tar, got,
			)
		},
	}
}

func (intValue) GreaterOrEqual(tar int) IntChecker {
	return intCheck{
		passFunc: func(got int) bool {
			return got >= tar
		},
		explFunc: func(label string, got interface{}) string {
			return fmt.Sprintf(
				"expect %s to be greater or equal to %d, got %d",
				label, tar, got,
			)
		},
	}
}

func (intValue) LesserThan(tar int) IntChecker {
	return intCheck{
		passFunc: func(got int) bool {
			return got < tar
		},
		explFunc: func(label string, got interface{}) string {
			return fmt.Sprintf(
				"expect %s to be lesser than %d, got %d",
				label, tar, got,
			)
		},
	}
}

func (intValue) LesserOrEqual(tar int) IntChecker {
	return intCheck{
		passFunc: func(got int) bool {
			return got <= tar
		},
		explFunc: func(label string, got interface{}) string {
			return fmt.Sprintf(
				"expect %s to be lesser or equal to %d, got %d",
				label, tar, got,
			)
		},
	}
}
