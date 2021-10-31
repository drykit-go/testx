package check

// TODO: rename file check.go when the other is removed

type PassFunc[T any] func(got T) bool

type (
	Passer[T any]  interface{ Pass(got T) bool }
	Checker[T any] interface {
		Passer[T]
		Explainer
	}
)

type checker[T any] struct {
	pass PassFunc[T]
	expl ExplainFunc
}

func (c checker[T]) Pass(got T) bool {
	return c.pass(got)
}

func (c checker[T]) Explain(label string, got any) string {
	return c.expl(label, got)
}

func NewChecker[T any](
	passFunc PassFunc[T],
	explainFunc ExplainFunc,
) Checker[T] {
	return checker[T]{
		pass: passFunc,
		expl: explainFunc,
	}
}
