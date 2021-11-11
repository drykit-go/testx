package check

import (
	"fmt"
	"reflect"

	"github.com/drykit-go/testx/internal/reflectutil"
)

// SliceCheckerProvider provides checks on kind slice.
type SliceCheckerProvider struct{ ValueCheckerProvider }

// Len checks the length of the gotten slice passes the given Checker[int].
func (p SliceCheckerProvider) Len(c Checker[int]) Checker[any] {
	var gotlen int
	pass := func(got any) bool {
		reflectutil.MustBeOfKind(got, reflect.Slice)
		gotlen = reflect.ValueOf(got).Len()
		return c.Pass(gotlen)
	}
	expl := func(label string, got any) string {
		return p.explainCheck(label,
			"length to pass Checker[int]",
			c.Explain("length", gotlen),
		)
	}
	return NewChecker(pass, expl)
}

// Cap checks the capacity of the gotten slice passes the given Checker[int].
func (p SliceCheckerProvider) Cap(c Checker[int]) Checker[any] {
	var gotcap int
	pass := func(got any) bool {
		reflectutil.MustBeOfKind(got, reflect.Slice)
		gotcap = reflect.ValueOf(got).Cap()
		return c.Pass(gotcap)
	}
	expl := func(label string, got any) string {
		return p.explainCheck(label,
			"capacity to pass Checker[int]",
			c.Explain("capacity", gotcap),
		)
	}
	return NewChecker(pass, expl)
}

// HasValues checks the gotten slice has the given values set.
func (p SliceCheckerProvider) HasValues(values ...any) Checker[any] {
	var missing []string
	pass := func(got any) bool {
		reflectutil.MustBeOfKind(got, reflect.Slice)
		for _, expv := range values {
			if !p.hasValue(got, expv) {
				missing = append(missing, fmt.Sprint(expv))
			}
		}
		return len(missing) == 0
	}
	expl := func(label string, got any) string {
		return p.explain(label,
			"to have values "+p.formatList(missing),
			got,
		)
	}
	return NewChecker(pass, expl)
}

// HasNotValues checks the gotten slice has not the given values set.
func (p SliceCheckerProvider) HasNotValues(values ...any) Checker[any] {
	var badvalues []string
	pass := func(got any) bool {
		reflectutil.MustBeOfKind(got, reflect.Slice)
		for _, badv := range values {
			if p.hasValue(got, badv) {
				badvalues = append(badvalues, fmt.Sprint(badv))
			}
		}
		return len(badvalues) == 0
	}
	expl := func(label string, got any) string {
		return p.explainNot(label,
			"to have values "+p.formatList(badvalues),
			got,
		)
	}
	return NewChecker(pass, expl)
}

// CheckValues checks the values of the gotten slice passes
// the given Checker[any].
// If a filterFunc is provided, the values not passing it are ignored.
func (p SliceCheckerProvider) CheckValues(
	c Checker[any],
	filters ...func(i int, v any) bool,
) Checker[any] {
	var badvalues []string
	pass := func(got any) bool {
		reflectutil.MustBeOfKind(got, reflect.Slice)
		p.walk(got, filters, func(i int, v any) {
			if !c.Pass(v) {
				badvalues = append(badvalues, fmt.Sprintf("%d:%v", i, v))
			}
		})
		return len(badvalues) == 0
	}
	expl := func(label string, _ any) string {
		return p.explainCheck(label,
			"values to pass Checker[any]",
			c.Explain("values", p.formatList(badvalues)),
		)
	}
	return NewChecker(pass, expl)
}

// Helpers

// hasValue returns true if slice has a value equal to expv.
func (p SliceCheckerProvider) hasValue(slice, expv any) bool {
	return p.walkUntil(slice, nil, func(_ int, v any) bool {
		return p.deq(v, expv)
	})
}

// walk iterates over a slice until the end is reached.
// It calls f(i, v) each iteration if (i, v) pass the given filters.
func (p SliceCheckerProvider) walk(
	slice any,
	filters []func(int, any) bool,
	f func(i int, v any),
) {
	p.walkUntil(slice, filters, func(i int, v any) bool {
		f(i, v)
		return false
	})
}

// walksUntil behaves like walk excepts it returns early if the stop func
// returns true for the current iteration. In returns true if it was stopped
// early, false otherwise.
func (p SliceCheckerProvider) walkUntil(
	slice any,
	filters []func(int, any) bool,
	stop func(int, any) bool,
) bool {
	vslice := reflect.ValueOf(slice)
	l := vslice.Len()
	for i := 0; i < l; i++ {
		v := vslice.Index(i).Interface()
		filter := p.mergeFilters(filters...)
		passed := filter(i, v)
		if passed && stop(i, v) {
			return true
		}
	}
	return false
}

// mergeFilters combinates several filtering funcs into one.
func (p SliceCheckerProvider) mergeFilters(
	filters ...func(int, any) bool,
) func(int, any) bool {
	if len(filters) == 0 {
		return func(int, any) bool { return true }
	}
	return func(i int, v any) bool {
		curr := filters[0]
		next := p.mergeFilters(filters[1:]...)
		return curr(i, v) && next(i, v)
	}
}
