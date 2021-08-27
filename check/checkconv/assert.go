package checkconv

import (
	"reflect"

	"github.com/drykit-go/testx/check"
)

// Assert returns a check.UntypedChecker built upon the given checker
// and a bool indicating whether it was successful.
//
// Contrary to Cast, it can perform conversions from checker types
// unknown by package check. That means it can work with any custom
// implementation provided it is valid (see IsChecker for details).
//
// However, Cast should be the first choice for a known checker type
// as Assert is about 10 times slower.
func Assert(anyChecker interface{}) (c check.UntypedChecker, ok bool) {
	if !IsChecker(anyChecker) {
		return
	}

	v := reflect.ValueOf(anyChecker)
	c = check.NewUntypedChecker(
		func(got interface{}) bool {
			gotv := reflect.ValueOf(got)
			return v.MethodByName(signaturePass.name).
				Call([]reflect.Value{gotv})[0].
				Bool()
		},
		func(label string, got interface{}) string {
			labv := reflect.ValueOf(label)
			gotv := reflect.ValueOf(got)
			return v.MethodByName(signatureExpl.name).
				Call([]reflect.Value{labv, gotv})[0].
				String()
		},
	)
	ok = true
	return
}
