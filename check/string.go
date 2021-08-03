package check

import (
	"fmt"
	"regexp"
	"strings"
)

type stringCheck struct {
	passFunc StringPassFunc
	explFunc ExplainFunc
}

func (c stringCheck) Pass(got string) bool {
	return c.passFunc(got)
}

func (c stringCheck) Explain(label string, got interface{}) string {
	return c.explFunc(label, got)
}

type stringValue struct{}

func (stringValue) Equal(tar string) StringChecker {
	return stringCheck{
		passFunc: func(got string) bool { return got == tar },
		explFunc: func(label string, got interface{}) string {
			return fmt.Sprintf(
				"expect %s to equal %v, got %v",
				label, tar, got,
			)
		},
	}
}

func (stringValue) Len(c IntChecker) StringChecker {
	return stringCheck{
		passFunc: func(got string) bool { return c.Pass(len(got)) },
		explFunc: func(label string, got interface{}) string {
			return fmt.Sprintf(
				"unexpected %s length: %s",
				label, c.Explain(label, len(got.(string))),
			)
		},
	}
}

func (stringValue) Match(rgx *regexp.Regexp) StringChecker {
	return stringCheck{
		passFunc: func(got string) bool { return rgx.MatchString(got) },
		explFunc: func(label string, got interface{}) string {
			return fmt.Sprintf(
				"expect %s to match regexp %s, got %s",
				label, rgx.String(), got,
			)
		},
	}
}

func (stringValue) NotMatch(rgx *regexp.Regexp) StringChecker {
	return stringCheck{
		passFunc: func(got string) bool { return !rgx.MatchString(got) },
		explFunc: func(label string, got interface{}) string {
			return fmt.Sprintf(
				"expect %s not to match regexp %s, got %s",
				label, rgx.String(), got,
			)
		},
	}
}

func (stringValue) Contains(tar string) StringChecker {
	return stringCheck{
		passFunc: func(got string) bool { return strings.Contains(got, tar) },
		explFunc: func(label string, got interface{}) string {
			return fmt.Sprintf(
				"expect %s to contain substring %s, got %s",
				label, tar, got,
			)
		},
	}
}

func (stringValue) NotContains(tar string) StringChecker {
	return stringCheck{
		passFunc: func(got string) bool { return !strings.Contains(got, tar) },
		explFunc: func(label string, got interface{}) string {
			return fmt.Sprintf(
				"expect %s not to contain substring %s, got %s",
				label, tar, got,
			)
		},
	}
}
