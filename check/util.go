package check

import (
	"encoding/json"
	"fmt"
	"reflect"
)

// sameJSONproduced returns true if xdata and ydata result in the same JSON value,
// ignoring the keys order.
// xptr and yptr must be pointers, as their values are filled
// with marshaled+unmarshaled json from xdata and ydata respectively.
//
// It panics on the first error encountered in the process.
func sameJSONproduced(xdata, ydata, xptr, yptr interface{}) bool {
	mustMarshal := func(in interface{}) []byte {
		b, err := json.Marshal(in)
		if err != nil {
			panic(fmt.Sprintf("failed to marshal value: %#v:\n%v", b, err))
		}
		return b
	}
	bx, by := mustMarshal(xdata), mustMarshal(ydata)
	return sameJSON(bx, by, &xptr, &yptr)
}

// sameJSON returns true if x and y evaluate to the same JSON value,
// ignoring the keys order.
// xptr and yptr must be pointers, as their values are filled
// with unmarshaled x and y respectively.
//
// It panics on the first error encountered in the process.
func sameJSON(x, y []byte, xptr, yptr interface{}) bool {
	mustUnmarshal := func(b []byte, ptr interface{}) {
		if err := json.Unmarshal(b, &ptr); err != nil {
			panic(err)
		}
	}
	mustUnmarshal(x, xptr)
	mustUnmarshal(y, yptr)
	return reflect.DeepEqual(xptr, yptr)
}

func deq(a, b interface{}) bool {
	return reflect.DeepEqual(a, b)
}
