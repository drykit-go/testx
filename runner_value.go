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

func (r *valueRunner) Exp(value interface{}) ValueRunner {
	r.addValueCheck(check.Value.Is(value))
	return r
}

func (r *valueRunner) Not(values ...interface{}) ValueRunner {
	r.addValueCheck(check.Value.Not(values...))
	return r
}

func (r *valueRunner) Pass(checkers ...check.ValueChecker) ValueRunner {
	for _, c := range checkers {
		r.addValueCheck(c)
	}
	return r
}

func (r *valueRunner) addValueCheck(c check.ValueChecker) {
	r.addCheck(baseCheck{
		label:   "value",
		get:     func() gottype { return r.value },
		checker: c,
	})
}

func newValueRunner(v interface{}) ValueRunner {
	return &valueRunner{value: v}
}
