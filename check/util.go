package check

import (
	"encoding/json"
	"io"
	"log"
	"reflect"
)

// sameJSONproduced returns true if xdata and ydata result in the same JSON value,
// ignoring the keys order.
// xptr and yptr must be pointers to interface{}, as their values are filled
// with marshaled+unmarshaled json from xdata and ydata respectively.
//
// It panics on the first error encountered in the process.
func sameJSONproduced(xdata, ydata, xptr, yptr interface{}) bool {
	bx, err := json.Marshal(xdata)
	if err != nil {
		log.Panicf("failed to marshal value: %#v\n%v", xdata, err)
	}
	by, err := json.Marshal(ydata)
	if err != nil {
		log.Panicf("failed to marshal value: %#v\n%v", ydata, err)
	}
	return sameJSON(bx, by, &xptr, &yptr)
}

// sameJSON returns true if x and y evaluate to the same JSON value,
// ignoring the keys order.
// xptr and yptr must be pointers to interface{}, as their values are filled
// with unmarshaled x and y respectively.
//
// It panics on the first error encountered in the process.
func sameJSON(x, y []byte, xptr, yptr interface{}) bool {
	if err := json.Unmarshal(x, &xptr); err != nil {
		log.Panic(err)
	}
	if err := json.Unmarshal(y, &yptr); err != nil {
		log.Panic(err)
	}
	return reflect.DeepEqual(xptr, yptr)
}

func deq(a, b interface{}) bool {
	return reflect.DeepEqual(a, b)
}

// mustReadIO reads and closes a reader. It panics if any error occurs
// in the process, in which case the provided label is used to specify
// some context in the error message.
func mustReadIO(label string, reader io.ReadCloser) []byte {
	b, err := io.ReadAll(reader)
	if err != nil {
		log.Panicf("error reading %s: %s", label, err)
	}
	return b
}
