package check

import (
	"fmt"
	"reflect"
	"strings"
)

// sliceCheckerProvider provides checks on kind slice.
type sliceCheckerProvider struct{}

// SameJSON checks the gotten slice and the target value
// produce the same JSON.
// It panics if any error occurs in the marshaling process.
func (sliceCheckerProvider) SameJSON(tar interface{}) ValueChecker {
	return Value.SameJSON(tar)
}

// Len checks the length of the gotten slice passes the given IntChecker.
func (sliceCheckerProvider) Len(c IntChecker) ValueChecker {
	pass := func(got interface{}) bool {
		panicOnUnexpectedKind(got, reflect.Slice)
		return c.Pass(reflect.ValueOf(got).Len())
	}
	expl := func(label string, got interface{}) string {
		return fmt.Sprintf(
			"exp %s length to pass IntChecker\ngot %s",
			label, c.Explain(label, got),
		)
	}
	return NewValueChecker(pass, expl)
}

// Cap checks the capacity of the gotten slice passes the given IntChecker.
func (sliceCheckerProvider) Cap(c IntChecker) ValueChecker {
	pass := func(got interface{}) bool {
		panicOnUnexpectedKind(got, reflect.Slice)
		return c.Pass(reflect.ValueOf(got).Cap())
	}
	expl := func(label string, got interface{}) string {
		return fmt.Sprintf(
			"exp %s capacity to pass IntChecker\ngot %s",
			label, c.Explain(label, got),
		)
	}
	return NewValueChecker(pass, expl)
}

// HasValues checks the gotten slice has the given values set.
func (p sliceCheckerProvider) HasValues(values ...interface{}) ValueChecker { //nolint: dupl
	var missing []string
	pass := func(got interface{}) bool {
		panicOnUnexpectedKind(got, reflect.Slice)
		for _, expv := range values {
			if !p.hasValue(got, expv) {
				missing = append(missing, fmt.Sprint(expv))
			}
		}
		return len(missing) == 0
	}
	expl := func(label string, _ interface{}) string {
		return fmt.Sprintf(
			"%s misses expected values: %s",
			label, strings.Join(missing, ","),
		)
	}
	return NewValueChecker(pass, expl)
}

// HasNotValues checks the gotten slice has not the given values set.
func (p sliceCheckerProvider) HasNotValues(values ...interface{}) ValueChecker {
	var badValues []string
	pass := func(got interface{}) bool {
		panicOnUnexpectedKind(got, reflect.Slice)
		for _, badv := range values {
			if p.hasValue(got, badv) {
				badValues = append(badValues, fmt.Sprint(badv))
			}
		}
		return len(badValues) == 0
	}
	expl := func(label string, _ interface{}) string {
		return fmt.Sprintf(
			"%s has unexpected values: %s",
			label, strings.Join(badValues, ","),
		)
	}
	return NewValueChecker(pass, expl)
}

// CheckValues checks the values of the gotten slice pass the given ValueChecker.
// If a filterFunc is provided, the values not passing it are ignored.
func (p sliceCheckerProvider) CheckValues(c ValueChecker, filters ...func(i int, v interface{}) bool) ValueChecker {
	var badValues []string
	pass := func(got interface{}) bool {
		panicOnUnexpectedKind(got, reflect.Slice)
		p.walk(got, filters, func(i int, v interface{}) {
			if !c.Pass(v) {
				badValues = append(badValues, fmt.Sprintf("%d:%v", i, v))
			}
		})
		return len(badValues) == 0
	}
	expl := func(label string, _ interface{}) string {
		return fmt.Sprintf(
			"%s has unexpected values: %s",
			label, strings.Join(badValues, ","),
		)
	}
	return NewValueChecker(pass, expl)
}

// Helpers

// hasValue returns true if slice has a value equal to expv.
func (p sliceCheckerProvider) hasValue(slice, expv interface{}) bool {
	return p.walkUntil(slice, nil, func(_ int, v interface{}) bool {
		return deq(v, expv)
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
		v := vslice.Index(i)
		filter := p.applyFilters(filters...)
		ignored := filter(i, v)
		if !ignored && stop(i, v) {
			return true
		}
	}
	return false
}

func (p sliceCheckerProvider) applyFilters(
	filters ...func(int, interface{}) bool,
) func(int, interface{}) bool {
	if len(filters) == 0 {
		return func(int, interface{}) bool { return true }
	}
	return func(i int, v interface{}) bool {
		curr := filters[0]
		next := p.applyFilters(filters[1:]...)
		return curr(i, v) && next(i, v)
	}
}
