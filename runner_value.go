package testx

import (
	"testing"

	"github.com/drykit-go/testx/check"
)

var _ ValueRunner = (*valueRunner)(nil)

type valueRunner struct {
	baseRunner
	value interface{}
}

func (r *valueRunner) Run(t *testing.T) {
	t.Helper()
	r.run(t)
}

func (r *valueRunner) DryRun() Resulter {
	return r.baseResults()
}

func (r *valueRunner) Exp(exp interface{}) ValueRunner {
	r.addCheck(baseCheck{
		label:   "value",
		get:     func() gottype { return r.value },
		checker: check.Value.Is(exp),
	})
	return r
}

func (r *valueRunner) Not(values ...interface{}) ValueRunner {
	r.addCheck(baseCheck{
		label:   "value",
		get:     func() gottype { return r.value },
		checker: check.Value.Not(values...),
	})
	return r
}

func (r *valueRunner) Pass(checkers ...check.ValueChecker) ValueRunner {
	r.addChecks("value", func() gottype { return r.value }, checkers, false)
	return r
}

func newValueRunner(v interface{}) ValueRunner {
	return &valueRunner{value: v}
}
