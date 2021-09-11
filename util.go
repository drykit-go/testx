package testx

import (
	"io"
	"log"
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

// mustReadIO reads and closes a reader. It panics if any error occurs
// in the process, in which case the provided origin allows to specify
// some context in the error message.
func mustReadIO(origin string, reader io.ReadCloser) []byte {
	b, err := io.ReadAll(reader)
	if err != nil {
		log.Panicf("error reading %s: %s", origin, err)
	}
	return b
}
