package check

import (
	"fmt"
	"regexp"
	"strings"
)

type stringCheckerFactory struct{}

func (stringCheckerFactory) Is(tar string) StringChecker {
	return stringChecker{
		passFunc: func(got string) bool { return got == tar },
		explFunc: func(label string, got interface{}) string {
			return fmt.Sprintf(
				"expect %s to equal %v, got %v",
				label, tar, got,
			)
		},
	}
}

func (stringCheckerFactory) Len(c IntChecker) StringChecker {
	return stringChecker{
		passFunc: func(got string) bool { return c.Pass(len(got)) },
		explFunc: func(label string, got interface{}) string {
			return fmt.Sprintf(
				"unexpected %s length: %s",
				label, c.Explain(label, len(got.(string))),
			)
		},
	}
}

func (stringCheckerFactory) Match(rgx *regexp.Regexp) StringChecker {
	return stringChecker{
		passFunc: func(got string) bool { return rgx.MatchString(got) },
		explFunc: func(label string, got interface{}) string {
			return fmt.Sprintf(
				"expect %s to match regexp %s, got %s",
				label, rgx.String(), got,
			)
		},
	}
}

func (stringCheckerFactory) NotMatch(rgx *regexp.Regexp) StringChecker {
	return stringChecker{
		passFunc: func(got string) bool { return !rgx.MatchString(got) },
		explFunc: func(label string, got interface{}) string {
			return fmt.Sprintf(
				"expect %s not to match regexp %s, got %s",
				label, rgx.String(), got,
			)
		},
	}
}

func (stringCheckerFactory) Contains(tar string) StringChecker {
	return stringChecker{
		passFunc: func(got string) bool { return strings.Contains(got, tar) },
		explFunc: func(label string, got interface{}) string {
			return fmt.Sprintf(
				"expect %s to contain substring %s, got %s",
				label, tar, got,
			)
		},
	}
}

func (stringCheckerFactory) NotContains(tar string) StringChecker {
	return stringChecker{
		passFunc: func(got string) bool { return !strings.Contains(got, tar) },
		explFunc: func(label string, got interface{}) string {
			return fmt.Sprintf(
				"expect %s not to contain substring %s, got %s",
				label, tar, got,
			)
		},
	}
}
