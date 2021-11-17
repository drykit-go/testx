package testx

import (
	"testing"

	"github.com/drykit-go/testx/check"
)

var _ ValueRunner[any] = (*valueRunner[any])(nil)

type valueRunner[T any] struct {
	baseRunner
	value T
}

func (r *valueRunner[T]) Run(t *testing.T) {
	t.Helper()
	r.run(t)
}

func (r *valueRunner[T]) DryRun() Resulter {
	return r.dryRun()
}

func (r *valueRunner[T]) Exp(value T) ValueRunner[T] {
	r.addValueCheck(check.Value[any]().Is(value))
	return r
}

func (r *valueRunner[T]) Not(values ...T) ValueRunner[T] {
	// TODO: find a way to cast properly
	valuesitf := []any{}
	for _, v := range values {
		valuesitf = append(valuesitf, v)
	}
	r.addValueCheck(check.Value[any]().Not(valuesitf...))
	return r
}

func (r *valueRunner[T]) Pass(checkers ...check.Checker[T]) ValueRunner[T] {
	for _, c := range checkers {
		r.addValueCheck(check.Wrap(c))
	}
	return r
}

func (r *valueRunner[T]) addValueCheck(c check.Checker[any]) {
	r.addCheck(baseCheck{
		label:   "value",
		get:     func() gottype { return r.value },
		checker: c,
	})
}

func newValueRunner[T any](v T) ValueRunner[T] {
	return &valueRunner[T]{value: v}
}
