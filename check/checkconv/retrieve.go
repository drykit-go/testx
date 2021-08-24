package checkconv

import (
	"reflect"

	"github.com/drykit-go/testx/check"
)

// Retrieve rebuilds a check.UntypedChecker that is masked by interface{} type.
// Its main purpose is for runners that accept several checker types, hence
// are forced to use an interface{} type.
// Contrary to UntypedChecker(), it also work on custom checkers unknown by
// package check.
func Retrieve(in interface{}) (c check.UntypedChecker, ok bool) {
	if !IsChecker(in) {
		return nil, false
	}
	inval := reflect.ValueOf(in)
	return check.NewUntypedCheck(
		func(got interface{}) bool {
			gotv := reflect.ValueOf(got)
			return inval.MethodByName("Pass").
				Call([]reflect.Value{gotv})[0].
				Bool()
		},
		func(label string, got interface{}) string {
			labv := reflect.ValueOf(label)
			gotv := reflect.ValueOf(got)
			return inval.MethodByName("Explain").
				Call([]reflect.Value{labv, gotv})[0].
				String()
		},
	), true
}
