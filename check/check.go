package check

type (
	// PassFunc is the required method to implement Passer.
	// It returns a boolean that indicates whether the got value
	// passes the current check.
	PassFunc[T any] func(got T) bool

	// ExplainFunc is the required method to implement Explainer.
	// It returns a string explaining why the gotten value failed the check.
	// The label provides some context, such as "response code".
	ExplainFunc func(label string, got any) string
)

type (
	// Passer provides a method Pass that returns a bool that indicates
	// whether the got value passes the current check.
	Passer[T any] interface{ Pass(got T) bool }

	// Explainer provides a method Explain describing the reason of a failed check.
	Explainer interface {
		Explain(label string, got any) string
	}
)

// Checker satisfies both Passer and Explainer interfaces.
type Checker[T any] interface {
	Passer[T]
	Explainer
}

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
