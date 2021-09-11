package reflectutil

import "errors"

// ErrNotAFunc is returned when a func expects a reflect.Func kind
// and receives a different one.
var ErrNotAFunc = errors.New("expect a func input")
