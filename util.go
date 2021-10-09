package testx

import (
	"reflect"
	"time"
)

var (
	deq = reflect.DeepEqual
	xor = func(a, b bool) bool { return a != b }
)

// timeFunc executes the given func and returns the elapsed time
// during the execution.
func timeFunc(f func()) time.Duration {
	t0 := time.Now()
	f()
	return time.Since(t0)
}
