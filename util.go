package testx

import (
	"io"
	"log"
	"path/filepath"
	"reflect"
	"runtime"
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

// getFuncName returns the name of the given func prefixed with
// the name of the package it is from.  It panics if f is not a func.
func getFuncName(f interface{}) string {
	fval := reflect.ValueOf(f)
	if kind := fval.Kind(); kind != reflect.Func {
		log.Panicf(
			"calling getFuncName with a non func argument (%s %v)",
			fval.Type().String(), fval.String(),
		)
	}
	fptr := fval.Pointer()
	path := runtime.FuncForPC(fptr).Name()
	return filepath.Base(path)
}

// panicOnErr panics if the given err is not nil.
func panicOnErr(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func condValue(vtrue, vfalse interface{}, istrue bool) interface{} {
	if istrue {
		return vtrue
	}
	return vfalse
}

func condString(vtrue, vfalse string, istrue bool) string {
	return condValue(vtrue, vfalse, istrue).(string)
}

// defaultString returns val if val !== "", def otherwise.
func defaultString(val, def string) string {
	return condString(val, def, val != "")
}
