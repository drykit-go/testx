package check

import (
	"fmt"
	"reflect"

	"github.com/drykit-go/testx/internal/reflectutil"
)

// sliceCheckerProvider provides checks on kind slice.
type sliceCheckerProvider struct{ valueCheckerProvider }

// Len checks the length of the gotten slice passes the given IntChecker.
func (p sliceCheckerProvider) Len(c IntChecker) ValueChecker {
	var gotlen int
	pass := func(got interface{}) bool {
		reflectutil.MustBeOfKind(got, reflect.Slice)
		gotlen = reflect.ValueOf(got).Len()
		return c.Pass(gotlen)
	}
	expl := func(label string, got interface{}) string {
		return p.explainCheck(label,
			"length to pass IntChecker",
			c.Explain("length", gotlen),
		)
	}
	return NewValueChecker(pass, expl)
}

// Cap checks the capacity of the gotten slice passes the given IntChecker.
func (p sliceCheckerProvider) Cap(c IntChecker) ValueChecker {
	var gotcap int
	pass := func(got interface{}) bool {
		reflectutil.MustBeOfKind(got, reflect.Slice)
		gotcap = reflect.ValueOf(got).Cap()
		return c.Pass(gotcap)
	}
	expl := func(label string, got interface{}) string {
		return p.explainCheck(label,
			"capacity to pass IntChecker",
			c.Explain("capacity", gotcap),
		)
	}
	return NewValueChecker(pass, expl)
}

// HasValues checks the gotten slice has the given values set.
func (p sliceCheckerProvider) HasValues(values ...interface{}) ValueChecker {
	var missing []string
	pass := func(got interface{}) bool {
		reflectutil.MustBeOfKind(got, reflect.Slice)
		for _, expv := range values {
			if !p.hasValue(got, expv) {
				missing = append(missing, fmt.Sprint(expv))
			}
		}
		return len(missing) == 0
	}
	expl := func(label string, got interface{}) string {
		return p.explain(label,
			"to have values "+p.formatList(missing),
			got,
		)
	}
	return NewValueChecker(pass, expl)
}

// HasNotValues checks the gotten slice has not the given values set.
func (p sliceCheckerProvider) HasNotValues(values ...interface{}) ValueChecker {
	var badvalues []string
	pass := func(got interface{}) bool {
		reflectutil.MustBeOfKind(got, reflect.Slice)
		for _, badv := range values {
			if p.hasValue(got, badv) {
				badvalues = append(badvalues, fmt.Sprint(badv))
			}
		}
		return len(badvalues) == 0
	}
	expl := func(label string, got interface{}) string {
		return p.explainNot(label,
			"to have values "+p.formatList(badvalues),
			got,
		)
	}
	return NewValueChecker(pass, expl)
}

// CheckValues checks the values of the gotten slice pass the given ValueChecker.
// If a filterFunc is provided, the values not passing it are ignored.
func (p sliceCheckerProvider) CheckValues(c ValueChecker, filters ...func(i int, v interface{}) bool) ValueChecker {
	var badvalues []string
	pass := func(got interface{}) bool {
		reflectutil.MustBeOfKind(got, reflect.Slice)
		p.walk(got, filters, func(i int, v interface{}) {
			if !c.Pass(v) {
				badvalues = append(badvalues, fmt.Sprintf("%d:%v", i, v))
			}
		})
		return len(badvalues) == 0
	}
	expl := func(label string, _ interface{}) string {
		return p.explainCheck(label,
			"values to pass ValueChecker",
			c.Explain("values", p.formatList(badvalues)),
		)
	}
	return NewValueChecker(pass, expl)
}

// Helpers

// hasValue returns true if slice has a value equal to expv.
func (p sliceCheckerProvider) hasValue(slice, expv interface{}) bool {
	return p.walkUntil(slice, nil, func(_ int, v interface{}) bool {
		return p.deq(v, expv)
	})
}

// walk iterates over a slice until the end is reached.
// It calls f(i, v) each iteration if (i, v) pass the given filters.
func (p sliceCheckerProvider) walk(
	slice interface{},
	filters []func(int, interface{}) bool,
	f func(i int, v interface{}),
) {
	p.walkUntil(slice, filters, func(i int, v interface{}) bool {
		f(i, v)
		return false
	})
}

// walksUntil behaves like walk excepts it returns early if the stop func
// returns true for the current iteration. In returns true if it was stopped
// early, false otherwise.
func (p sliceCheckerProvider) walkUntil(
	slice interface{},
	filters []func(int, interface{}) bool,
	stop func(int, interface{}) bool,
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
func (p sliceCheckerProvider) mergeFilters(
	filters ...func(int, interface{}) bool,
) func(int, interface{}) bool {
	if len(filters) == 0 {
		return func(int, interface{}) bool { return true }
	}
	return func(i int, v interface{}) bool {
		curr := filters[0]
		next := p.mergeFilters(filters[1:]...)
		return curr(i, v) && next(i, v)
	}
}
