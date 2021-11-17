package check

import (
	"fmt"
	"reflect"
	"sort"

	"github.com/drykit-go/cond"

	"github.com/drykit-go/testx/internal/reflectutil"
)

// MapCheckerProvider provides checks on kind map.
type MapCheckerProvider[K comparable, V any] struct{ ValueCheckerProvider[map[K]V] }

// Len checks the gotten map passes the given Checker[int].
func (p MapCheckerProvider[K, V]) Len(c Checker[int]) Checker[map[K]V] {
	var gotlen int
	pass := func(got map[K]V) bool {
		reflectutil.MustBeOfKind(got, reflect.Map)
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

// HasKeys checks the gotten map has the given keys set.
func (p MapCheckerProvider[K, V]) HasKeys(keys ...K) Checker[map[K]V] {
	var missing []string
	pass := func(got map[K]V) bool {
		reflectutil.MustBeOfKind(got, reflect.Map)
		for _, expk := range keys {
			if _, found := p.get(got, expk); !found {
				missing = append(missing, fmt.Sprint(expk))
			}
		}
		return len(missing) == 0
	}
	expl := func(label string, got any) string {
		return p.explain(label, "to have keys "+p.formatList(missing), got)
	}
	return NewChecker(pass, expl)
}

// HasNotKeys checks the gotten map has the given keys set.
func (p MapCheckerProvider[K, V]) HasNotKeys(keys ...K) Checker[map[K]V] {
	var badkeys []string
	pass := func(got map[K]V) bool {
		reflectutil.MustBeOfKind(got, reflect.Map)
		for _, expk := range keys {
			if _, found := p.get(got, expk); found {
				badkeys = append(badkeys, fmt.Sprint(expk))
			}
		}
		return len(badkeys) == 0
	}
	expl := func(label string, got any) string {
		return p.explainNot(label, "to have keys "+p.formatList(badkeys), got)
	}
	return NewChecker(pass, expl)
}

// HasValues checks the gotten map has the given values set.
func (p MapCheckerProvider[K, V]) HasValues(values ...V) Checker[map[K]V] {
	var missing []string
	pass := func(got map[K]V) bool {
		reflectutil.MustBeOfKind(got, reflect.Map)
		for _, expv := range values {
			if !p.hasValue(got, expv) {
				missing = append(missing, fmt.Sprint(expv))
			}
		}
		return len(missing) == 0
	}
	expl := func(label string, got any) string {
		return p.explain(label, "to have values "+p.formatList(missing), got)
	}
	return NewChecker(pass, expl)
}

// HasNotValues checks the gotten map has not the given values set.
func (p MapCheckerProvider[K, V]) HasNotValues(values ...V) Checker[map[K]V] {
	var badvalues []string
	pass := func(got map[K]V) bool {
		reflectutil.MustBeOfKind(got, reflect.Map)
		for _, badv := range values {
			if p.hasValue(got, badv) {
				badvalues = append(badvalues, fmt.Sprint(badv))
			}
		}
		return len(badvalues) == 0
	}
	expl := func(label string, got any) string {
		return p.explainNot(label, "to have values "+p.formatList(badvalues), got)
	}
	return NewChecker(pass, expl)
}

// CheckValues checks the gotten map's values corresponding to the given keys
// pass the given checker. A key not found is considered a fail.
// If len(keys) == 0, the check is made on all map values.
func (p MapCheckerProvider[K, V]) CheckValues(c Checker[V], keys ...K) Checker[map[K]V] { //nolint: gocognit // TODO: refactor
	var badentries []string
	allKeys := len(keys) == 0
	pass := func(got map[K]V) bool {
		reflectutil.MustBeOfKind(got, reflect.Map)
		if allKeys {
			p.walk(got, func(gotk K, gotv V) {
				if !c.Pass(gotv) {
					badentries = append(badentries, fmt.Sprintf("%v:%v", gotk, gotv))
				}
			})
			sort.Strings(badentries)
		} else {
			for _, expk := range keys {
				gotv, ok := p.get(got, expk)
				if !ok || !c.Pass(gotv) {
					badentries = append(badentries, fmt.Sprintf("%v:%v", expk, gotv))
				}
			}
		}
		return len(badentries) == 0
	}
	expl := func(label string, _ any) string {
		checkedKeys := cond.String("all keys", fmt.Sprintf("keys %v", keys), allKeys)
		return p.explainCheck(label,
			fmt.Sprintf("values for %s to pass Checker[V]", checkedKeys),
			c.Explain("values", p.formatList(badentries)),
		)
	}
	return NewChecker(pass, expl)
}

// get returns gotmap[key] and a bool representing whether a match is found.
func (p MapCheckerProvider[K, V]) get(gotmap map[K]V, key K) (V, bool) {
	iter := reflect.ValueOf(gotmap).MapRange()
	for iter.Next() {
		if k := iter.Key().Interface(); p.deq(k, key) {
			v, ok := iter.Value().Interface().(V)
			if !ok {
				var vnil V
				return vnil, true
			}
			return v, true
		}
	}
	var vnil V
	return vnil, false
}

// hasValue returns true if gotmap matches the specified value.
func (p MapCheckerProvider[K, V]) hasValue(gotmap map[K]V, value V) bool {
	iter := reflect.ValueOf(gotmap).MapRange()
	for iter.Next() {
		if gotv := iter.Value().Interface(); p.deq(gotv, value) {
			return true
		}
	}
	return false
}

func (MapCheckerProvider[K, V]) walk(gotmap map[K]V, f func(k K, v V)) {
	vmap := reflect.ValueOf(gotmap)
	iter := vmap.MapRange()
	for iter.Next() {
		var knil K
		k, ok := iter.Key().Interface().(K)
		if !ok {
			k = knil
		}
		var vnil V
		v, ok := iter.Value().Interface().(V)
		if !ok {
			v = vnil
		}
		f(k, v)
	}
}
