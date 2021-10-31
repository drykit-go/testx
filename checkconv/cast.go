package checkconv

import (
	"reflect"

	"github.com/drykit-go/testx/check"
)

// Cast returns a check.Checker[any] built upon the given checker
// and a bool indicating whether it was successful.
//
// Contrary to Assert, it can perform conversions from checker types
// unknown by package check. That means it can work with any custom
// implementation provided it is valid (see IsChecker for details).
//
// However, Assert should be the first choice for a known checker type
// as Cast is about 10 times slower.
func Cast(anyChecker interface{}) (c check.Checker[any], ok bool) {
	if !IsChecker(anyChecker) {
		return
	}

	v := reflect.ValueOf(anyChecker)
	c = check.NewChecker(
		func(got interface{}) bool {
			gotv := reflect.ValueOf(got)
			return v.MethodByName(signaturePass.Name).
				Call([]reflect.Value{gotv})[0].
				Bool()
		},
		func(label string, got interface{}) string {
			labv := reflect.ValueOf(label)
			gotv := reflect.ValueOf(got)
			return v.MethodByName(signatureExpl.Name).
				Call([]reflect.Value{labv, gotv})[0].
				String()
		},
	)
	ok = true
	return
}

// CastMany converts the given checkers as described by Cast,
// and returns them as a slice of check.ValueChecker and a bool
// indicating whether all conversions were successful.
//
// An invalid checker in the args list is silently dismissed,
// this the resulting checkers length can be inferior to the number of args
// if ok === false.
func CastMany(anyCheckers ...interface{}) (checkers []check.Checker[any], ok bool) {
	ok = true
	for _, in := range anyCheckers {
		c, valid := Cast(in)
		if valid {
			checkers = append(checkers, c)
		} else {
			ok = false
		}
	}
	return
}
