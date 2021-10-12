package testx

import (
	"time"
)

// timeFunc executes the given func and returns the elapsed time
// during the execution.
func timeFunc(f func()) time.Duration {
	t0 := time.Now()
	f()
	return time.Since(t0)
}
