package check

import (
	"fmt"
	"regexp"
	"strings"
)

type stringCheckFactory struct{}

func (stringCheckFactory) Equal(tar string) StringChecker {
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

func (stringCheckFactory) Len(c IntChecker) StringChecker {
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

func (stringCheckFactory) Match(rgx *regexp.Regexp) StringChecker {
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

func (stringCheckFactory) NotMatch(rgx *regexp.Regexp) StringChecker {
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

func (stringCheckFactory) Contains(tar string) StringChecker {
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

func (stringCheckFactory) NotContains(tar string) StringChecker {
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
