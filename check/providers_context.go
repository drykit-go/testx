package check

import (
	"context"
	"fmt"
	"strings"

	"github.com/drykit-go/cond"
)

// contextCheckerProvider provides checks on type context.Context.
type contextCheckerProvider struct{ baseCheckerProvider }

// Done checks the gotten context is done.
func (p contextCheckerProvider) Done(expectDone bool) ContextChecker {
	var err error
	done := func() bool { return err != nil }
	pass := func(got context.Context) bool {
		err = got.Err()
		return done() == expectDone
	}
	expl := func(label string, _ interface{}) string {
		notString := cond.String("", "not ", expectDone)
		expString := fmt.Sprintf("context %sto be done", notString)
		gotString := cond.String(fmt.Sprint(err), "context not done", done())
		return p.explain(label, expString, gotString)
	}
	return NewContextChecker(pass, expl)
}

// HasKeys checks the gotten context has the given keys set.
func (p contextCheckerProvider) HasKeys(keys ...interface{}) ContextChecker {
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
	expl := func(label string, got interface{}) string {
		return p.explain(label,
			"to have keys "+strings.Join(missing, ","),
			"keys not set",
		)
	}
	return NewContextChecker(pass, expl)
}

// Value checks the gotten context's value for the given key passes
// the given ValueChecker. It fails if value is nil.
//
// Examples:
// 	Context.Value("userID", Value.Is("abcde"))
// 	Context.Value("userID", checkconv.Assert(String.Contains("abc")))
func (p contextCheckerProvider) Value(key interface{}, c ValueChecker) ContextChecker {
	var v interface{}
	pass := func(got context.Context) bool {
		v = got.Value(key)
		return v != nil && c.Pass(v)
	}
	expl := func(label string, got interface{}) string {
		return p.explainCheck(label,
			fmt.Sprintf("value for key %v to pass ValueChecker", key),
			c.Explain("value", v),
		)
	}
	return NewContextChecker(pass, expl)
}
