package check

import (
	"fmt"
	"reflect"
	"strings"
)

// mapCheckerProvider provides checks on kind map.
type mapCheckerProvider struct{ baseCheckerProvider }

// SameJSON checks the gotten map and the target value
// result in the same JSON.
// It panics if any error occurs in the marshaling process.
func (mapCheckerProvider) SameJSON(tar interface{}) ValueChecker {
	return Value.SameJSON(tar)
}

// Len checks the gotten map passes the given IntChecker.
func (p mapCheckerProvider) Len(c IntChecker) ValueChecker {
	var gotlen int
	pass := func(got interface{}) bool {
		panicOnUnexpectedKind(got, reflect.Map)
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
	expl := func(label string, got interface{}) string {
		return p.explain(label, "to have keys "+strings.Join(missing, ","), got)
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
	expl := func(label string, got interface{}) string {
		return p.explainNot(label, "to have keys "+strings.Join(badKeys, ","), got)
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
	expl := func(label string, got interface{}) string {
		return p.explain(label, "to have values "+strings.Join(missing, ","), got)
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
	expl := func(label string, got interface{}) string {
		return p.explainNot(label, "to have values "+strings.Join(badValues, ","), got)
	}
	return NewValueChecker(pass, expl)
}

// CheckValues checks the gotten map's values corresponding to the given keys
// pass the given checker. A key not found is considered a fail.
// If len(keys) == 0, the check is made on all map values.
func (p mapCheckerProvider) CheckValues(c ValueChecker, keys ...interface{}) ValueChecker { //nolint: gocognit // TODO: refactor
	var badEntries []string
	pass := func(got interface{}) bool {
		panicOnUnexpectedKind(got, reflect.Map)
		if len(keys) == 0 {
			p.walk(got, func(gotk, gotv interface{}) {
				if !c.Pass(gotv) {
					badEntries = append(badEntries, fmt.Sprint(gotk))
				}
			})
		} else {
			for _, expk := range keys {
				gotv, ok := p.get(got, expk)
				if !ok || !c.Pass(gotv) {
					badEntries = append(badEntries, fmt.Sprint(expk))
				}
			}
		}
		return len(badEntries) == 0
	}
	expl := func(label string, _ interface{}) string {
		return p.explainCheck(label,
			fmt.Sprintf("values for keys %v to pass ValueChecker", keys),
			c.Explain("values", "fail"),
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

func (mapCheckerProvider) walk(gotmap interface{}, f func(k, v interface{})) {
	vmap := reflect.ValueOf(gotmap)
	iter := vmap.MapRange()
	for iter.Next() {
		k := iter.Key().Interface()
		v := iter.Value().Interface()
		f(k, v)
	}
}
