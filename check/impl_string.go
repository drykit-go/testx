package check

import (
	"fmt"
	"regexp"
	"strings"
)

type stringCheckerFactory struct{}

func (stringCheckerFactory) Is(tar string) StringChecker {
	pass := func(got string) bool { return got == tar }
	expl := func(label string, got interface{}) string {
		return fmt.Sprintf(
			"expect %s to equal %v, got %v",
			label, tar, got,
		)
	}
	return NewStringChecker(pass, expl)
}

func (stringCheckerFactory) Len(c IntChecker) StringChecker {
	pass := func(got string) bool { return c.Pass(len(got)) }
	expl := func(label string, got interface{}) string {
		return fmt.Sprintf(
			"unexpected %s length: %s",
			label, c.Explain(label, len(got.(string))),
		)
	}
	return NewStringChecker(pass, expl)
}

func (stringCheckerFactory) Match(rgx *regexp.Regexp) StringChecker {
	pass := func(got string) bool { return rgx.MatchString(got) }
	expl := func(label string, got interface{}) string {
		return fmt.Sprintf(
			"expect %s to match regexp %s, got %s",
			label, rgx.String(), got,
		)
	}
	return NewStringChecker(pass, expl)
}

func (stringCheckerFactory) NotMatch(rgx *regexp.Regexp) StringChecker {
	pass := func(got string) bool { return !rgx.MatchString(got) }
	expl := func(label string, got interface{}) string {
		return fmt.Sprintf(
			"expect %s not to match regexp %s, got %s",
			label, rgx.String(), got,
		)
	}
	return NewStringChecker(pass, expl)
}

func (stringCheckerFactory) Contains(tar string) StringChecker {
	pass := func(got string) bool { return strings.Contains(got, tar) }
	expl := func(label string, got interface{}) string {
		return fmt.Sprintf(
			"expect %s to contain substring %s, got %s",
			label, tar, got,
		)
	}
	return NewStringChecker(pass, expl)
}

func (stringCheckerFactory) NotContains(tar string) StringChecker {
	pass := func(got string) bool { return !strings.Contains(got, tar) }
	expl := func(label string, got interface{}) string {
		return fmt.Sprintf(
			"expect %s not to contain substring %s, got %s",
			label, tar, got,
		)
	}
	return NewStringChecker(pass, expl)
}
