package testx

import (
	"fmt"
	"testing"

	"github.com/drykit-go/testx/check"
)

var _ ValueRunner = (*valueRunner)(nil)

type valueRunner struct {
	baseRunner
	value interface{}
}

func (r *valueRunner) Run(t *testing.T) {
	r.run(t)
}

func (r *valueRunner) MustBe(exp interface{}) ValueRunner {
	pass := func(got interface{}) bool { return deq(got, exp) }
	expl := func(label string, got interface{}) string {
		return fmt.Sprintf("%s: expect to be %v, is %v", label, exp, got)
	}
	r.addCheck(testCheck{
		"value check",
		func() gotType { return r.value },
		check.NewUntypedCheck(pass, expl),
	})
	return r
}

func (r *valueRunner) MustNotBe(values ...interface{}) ValueRunner {
	for _, nexp := range values {
		pass := func(got interface{}) bool { return !deq(got, nexp) }
		expl := func(label string, got interface{}) string {
			return fmt.Sprintf("%s: expect not to be %v, is %v", label, nexp, got)
		}
		r.addCheck(testCheck{
			"value check",
			func() gotType { return r.value },
			check.NewUntypedCheck(pass, expl),
		})
	}
	return r
}

func (r *valueRunner) MustPass(checks ...interface{}) ValueRunner {
	r.addChecks("value check", func() gotType { return r.value }, checks)
	return r
}

func Value(v interface{}) ValueRunner {
	return &valueRunner{value: v}
}
