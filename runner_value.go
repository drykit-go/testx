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
	return r.baseResults()
}

func (r *valueRunner) MustBe(exp interface{}) ValueRunner {
	pass := func(got interface{}) bool { return deq(got, exp) }
	expl := func(label string, got interface{}) string {
		return fmt.Sprintf("%s: expect %v, got %v", label, exp, got)
	}
	r.addCheck(baseCheck{
		"value",
		func() gottype { return r.value },
		check.NewUntypedChecker(pass, expl),
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
		r.addCheck(baseCheck{
			"value",
			func() gottype { return r.value },
			check.NewUntypedChecker(pass, expl),
		})
	}
	return r
}

func (r *valueRunner) MustPass(checkers ...interface{}) ValueRunner {
	r.addChecks("value", func() gottype { return r.value }, checkers, false)
	return r
}

func Value(v interface{}) ValueRunner {
	return &valueRunner{value: v}
}
