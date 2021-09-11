package reflectutil

import (
	"log"
	"reflect"
)

const AnyKind reflect.Kind = 27

// IsZero returns true if v is a zero value.
func IsZero(v interface{}) bool {
	return reflect.ValueOf(v).IsZero()
}

// CallUnwrap calls fn with args and returns the output values
// as []interface{}.
func CallUnwrap(fval reflect.Value, args []interface{}) (output []interface{}) {
	return UnwrapValues(fval.Call(WrapValues(args)))
}

// WrapValues wraps the input values into reflect.Values.
func WrapValues(values []interface{}) (wrapped []reflect.Value) {
	wrapped = make([]reflect.Value, len(values))
	for i, v := range values {
		wrapped[i] = reflect.ValueOf(v)
	}
	return
}

// UnwrapValues unwraps reflect.Values to empty interfaces.
func UnwrapValues(wrapped []reflect.Value) (values []interface{}) {
	values = make([]interface{}, len(wrapped))
	for i, w := range wrapped {
		values[i] = w.Interface()
	}
	return
}

func PanicOnUnexpectedKind(v interface{}, exp reflect.Kind) {
	panicOnUnexpectedKind(reflect.ValueOf(v), exp)
}

func panicOnUnexpectedKind(v reflect.Value, k reflect.Kind) {
	if v.Kind() != k {
		log.Panicf(
			"expect %s kind as input, got %s",
			k.String(), v.String(),
		)
	}
}
