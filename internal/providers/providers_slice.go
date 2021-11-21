package providers

import (
	"fmt"
	"reflect"

	check "github.com/drykit-go/testx/internal/checktypes"
	"github.com/drykit-go/testx/internal/reflectutil"
)

// SliceCheckerProvider provides checks on kind slice.
type SliceCheckerProvider[Elem any] struct{ ValueCheckerProvider[[]Elem] }

// Len checks the length of the gotten slice passes the given Checker[int].
func (p SliceCheckerProvider[Elem]) Len(c check.Checker[int]) check.Checker[[]Elem] {
	var gotlen int
	pass := func(got []Elem) bool {
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
	return check.NewChecker(pass, expl)
}

// Cap checks the capacity of the gotten slice passes the given Checker[int].
func (p SliceCheckerProvider[Elem]) Cap(c check.Checker[int]) check.Checker[[]Elem] {
	var gotcap int
	pass := func(got []Elem) bool {
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
	return check.NewChecker(pass, expl)
}

// HasValues checks the gotten slice has the given values set.
func (p SliceCheckerProvider[Elem]) HasValues(values ...Elem) check.Checker[[]Elem] {
	var missing []string
	pass := func(got []Elem) bool {
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
	return check.NewChecker(pass, expl)
}

// HasNotValues checks the gotten slice has not the given values set.
func (p SliceCheckerProvider[Elem]) HasNotValues(values ...Elem) check.Checker[[]Elem] {
	var badvalues []string
	pass := func(got []Elem) bool {
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
	return check.NewChecker(pass, expl)
}

// CheckValues checks the values of the gotten slice passes
// the given Checker[Elem].
// If a filterFunc is provided, the values not passing it are ignored.
func (p SliceCheckerProvider[Elem]) CheckValues(
	c check.Checker[Elem],
	filters ...func(i int, v Elem) bool,
) check.Checker[[]Elem] {
	var badvalues []string
	pass := func(got []Elem) bool {
		reflectutil.MustBeOfKind(got, reflect.Slice)
		p.walk(got, filters, func(i int, v Elem) {
			if !c.Pass(v) {
				badvalues = append(badvalues, fmt.Sprintf("%d:%v", i, v))
			}
		})
		return len(badvalues) == 0
	}
	expl := func(label string, _ any) string {
		return p.explainCheck(label,
			"values to pass Checker[Elem]",
			c.Explain("values", p.formatList(badvalues)),
		)
	}
	return check.NewChecker(pass, expl)
}

// Helpers

// hasValue returns true if slice has a value equal to expv.
func (p SliceCheckerProvider[Elem]) hasValue(slice []Elem, expv Elem) bool {
	return p.walkUntil(slice, nil, func(_ int, v Elem) bool {
		return p.deq(v, expv)
	})
}

// walk iterates over a slice until the end is reached.
// It calls f(i, v) each iteration if (i, v) pass the given filters.
func (p SliceCheckerProvider[Elem]) walk(
	slice []Elem,
	filters []func(int, Elem) bool,
	f func(i int, v Elem),
) {
	p.walkUntil(slice, filters, func(i int, v Elem) bool {
		f(i, v)
		return false
	})
}

// walksUntil behaves like walk excepts it returns early if the stop func
// returns true for the current iteration. In returns true if it was stopped
// early, false otherwise.
func (p SliceCheckerProvider[Elem]) walkUntil(
	slice []Elem,
	filters []func(int, Elem) bool,
	stop func(int, Elem) bool,
) bool {
	vslice := reflect.ValueOf(slice)
	l := vslice.Len()
	for i := 0; i < l; i++ {
		v := vslice.Index(i).Interface().(Elem)
		filter := p.mergeFilters(filters...)
		passed := filter(i, v)
		if passed && stop(i, v) {
			return true
		}
	}
	return false
}

// mergeFilters combinates several filtering funcs into one.
func (p SliceCheckerProvider[Elem]) mergeFilters(
	filters ...func(int, Elem) bool,
) func(int, Elem) bool {
	if len(filters) == 0 {
		return func(int, Elem) bool { return true }
	}
	return func(i int, v Elem) bool {
		curr := filters[0]
		next := p.mergeFilters(filters[1:]...)
		return curr(i, v) && next(i, v)
	}
}
