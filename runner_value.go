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

func (r *valueRunner) DryRun() Resulter {
	for _, c := range r.checks {
		r.updateBaseResults(c)
	}
	return r.baseResults
}

func (r *valueRunner) MustBe(exp interface{}) ValueRunner {
	pass := func(got interface{}) bool { return deq(got, exp) }
	expl := func(label string, got interface{}) string {
		return fmt.Sprintf("%s: expect %v, got %v", label, exp, got)
	}
	r.addCheck(testCheck{
		"value",
		func() gotType { return r.value },
		check.NewUntypedCheck(pass, expl),
	})
	return r
}

func (r *valueRunner) MustNotBe(values ...interface{}) ValueRunner {
	for _, nexp := range values {
		nexp := nexp
		pass := func(got interface{}) bool { return !deq(got, nexp) }
		expl := func(label string, got interface{}) string {
			return fmt.Sprintf("%s: expect not %v, got %v", label, nexp, got)
		}
		r.addCheck(testCheck{
			"value",
			func() gotType { return r.value },
			check.NewUntypedCheck(pass, expl),
		})
	}
	return r
}

func (r *valueRunner) MustPass(checkers ...interface{}) ValueRunner {
	r.addChecks("value", func() gotType { return r.value }, checkers)
	return r
}

func Value(v interface{}) ValueRunner {
	return &valueRunner{value: v}
}
