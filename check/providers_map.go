package check

import (
	"fmt"
	"reflect"
	"strings"
)

// mapCheckerProvider provides checks on kind map.
type mapCheckerProvider struct{}

// SameJSON checks the gotten map and the target value
// result in the same JSON.
// It panics if any error occurs in the marshaling process.
func (mapCheckerProvider) SameJSON(tar interface{}) ValueChecker {
	return Value.SameJSON(tar)
}

// Len checks the gotten map passes the given IntChecker.
func (mapCheckerProvider) Len(c IntChecker) ValueChecker {
	pass := func(got interface{}) bool {
		panicOnUnexpectedKind(got, reflect.Map)
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

// HasKey checks the gotten map has the given keys set.
func (p mapCheckerProvider) HasKeys(keys ...interface{}) ValueChecker {
	var missing []string
	pass := func(got interface{}) bool {
		panicOnUnexpectedKind(got, reflect.Map)
		for _, expk := range keys {
			if _, found := p.get(got, expk); !found {
				missing = append(missing, fmt.Sprint(expk))
			}
		}
		return len(missing) == 0
	}
	expl := func(label string, _ interface{}) string {
		return fmt.Sprintf(
			"%s misses expected keys: %s",
			label, strings.Join(missing, ","),
		)
	}
	return NewValueChecker(pass, expl)
}

// HasNotKey checks the gotten map has the given keys set.
func (p mapCheckerProvider) HasNotKeys(keys ...interface{}) ValueChecker {
	var badKeys []string
	pass := func(got interface{}) bool {
		panicOnUnexpectedKind(got, reflect.Map)
		for _, expk := range keys {
			if _, found := p.get(got, expk); found {
				badKeys = append(badKeys, fmt.Sprint(expk))
			}
		}
		return len(badKeys) == 0
	}
	expl := func(label string, _ interface{}) string {
		return fmt.Sprintf(
			"%s has unexpected keys: %s",
			label, strings.Join(badKeys, ","),
		)
	}
	return NewValueChecker(pass, expl)
}

// HasValues checks the gotten map has the given values set.
func (p mapCheckerProvider) HasValues(values ...interface{}) ValueChecker {
	var missing []string
	pass := func(got interface{}) bool {
		panicOnUnexpectedKind(got, reflect.Map)
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

// HasNotValues checks the gotten map has not the given values set.
func (p mapCheckerProvider) HasNotValues(values ...interface{}) ValueChecker {
	var badValues []string
	pass := func(got interface{}) bool {
		panicOnUnexpectedKind(got, reflect.Map)
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

func (p mapCheckerProvider) CheckValues(c ValueChecker, keys []interface{}) ValueChecker {
	var badEntries []string
	pass := func(got interface{}) bool {
		panicOnUnexpectedKind(got, reflect.Map)
		for _, expk := range keys {
			gotv, ok := p.get(got, expk)
			if !ok || !c.Pass(gotv) {
				badEntries = append(badEntries, fmt.Sprint(expk))
			}
		}
		return len(badEntries) == 0
	}
	expl := func(label string, _ interface{}) string {
		return fmt.Sprintf(
			"%s has missing or unexpected values for keys: %s",
			label, strings.Join(badEntries, ","),
		)
	}
	return NewValueChecker(pass, expl)
}

// get returns gotmap[key] and a bool representing whether a match is found.
func (mapCheckerProvider) get(gotmap, key interface{}) (interface{}, bool) {
	iter := reflect.ValueOf(gotmap).MapRange()
	for iter.Next() {
		if k := iter.Key().Interface(); deq(k, key) {
			return iter.Value().Interface(), true
		}
	}
	return nil, false
}

// hasValue returns true if gotmap matches the specified value.
func (mapCheckerProvider) hasValue(gotmap, value interface{}) bool {
	iter := reflect.ValueOf(gotmap).MapRange()
	for iter.Next() {
		if gotv := iter.Value().Interface(); deq(gotv, value) {
			return true
		}
	}
	return false
}
