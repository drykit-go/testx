package check

import (
	"context"
	"fmt"

	"github.com/drykit-go/cond"
)

// ContextCheckerProvider provides checks on type context.Context.
type ContextCheckerProvider struct{ baseCheckerProvider }

// Done checks the gotten context is done.
func (p ContextCheckerProvider) Done(expectDone bool) Checker[context.Context] {
	var err error
	done := func() bool { return err != nil }
	pass := func(got context.Context) bool {
		err = got.Err()
		return done() == expectDone
	}
	expl := func(label string, _ any) string {
		notString := cond.String("", "not ", expectDone)
		expString := fmt.Sprintf("context %sto be done", notString)
		gotString := cond.String(fmt.Sprint(err), "context not done", done())
		return p.explain(label, expString, gotString)
	}
	return NewChecker(pass, expl)
}

// HasKeys checks the gotten context has the given keys set.
func (p ContextCheckerProvider) HasKeys(keys ...any) Checker[context.Context] {
	var missing []string
	pass := func(got context.Context) bool {
		for _, expk := range keys {
			v := got.Value(expk)
			if v == nil {
				missing = append(missing, fmt.Sprint(expk))
			}
		}
		return len(missing) == 0
	}
	expl := func(label string, got any) string {
		return p.explain(label,
			"to have keys "+p.formatList(missing),
			"keys not set",
		)
	}
	return NewChecker(pass, expl)
}

// Value checks the gotten context's value for the given key passes
// the given ValueChecker. It fails if value is nil.
//
// Examples:
// 	Context.Value("userID", Value.Is("abcde"))
// 	Context.Value("userID", Wrap(String.Contains("abc")))
func (p ContextCheckerProvider) Value(key any, c Checker[any]) Checker[context.Context] {
	var v any
	pass := func(got context.Context) bool {
		v = got.Value(key)
		return v != nil && c.Pass(v)
	}
	expl := func(label string, got any) string {
		return p.explainCheck(label,
			fmt.Sprintf("value for key %v to pass Checker[any]", key),
			c.Explain("value", v),
		)
	}
	return NewChecker(pass, expl)
}
