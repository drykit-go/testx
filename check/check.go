package check

import (
	"github.com/drykit-go/testx/internal/checktypes"
)

type Numeric interface{ checktypes.Numeric }

type (
	// PassFunc is the required method to implement Passer.
	// It returns a boolean that indicates whether the got value
	// passes the current check.
	PassFunc[T any] checktypes.PassFunc[T]

	// ExplainFunc is the required method to implement Explainer.
	// It returns a string explaining why the gotten value failed the check.
	// The label provides some context, such as "response code".
	ExplainFunc = checktypes.ExplainFunc
)

type (
	// Passer provides a method Pass that returns a bool that indicates
	// whether the got value passes the current check.
	Passer[T any] interface{ checktypes.Passer[T] }

	// Explainer provides a method Explain describing the reason of a failed check.
	Explainer interface{ checktypes.Explainer }

	// Checker satisfies both Passer and Explainer interfaces.
	Checker[T any] interface{ checktypes.Checker[T] }
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
