// Package checkconv provides functions to convert typed checks
// into generic ones.
package checkconv

import (
	"log"
	"net/http"
	"reflect"
	"time"

	"github.com/drykit-go/testix/check"
)

func FromBytes(c check.BytesChecker) check.UntypedChecker {
	return check.NewUntypedCheck(
		func(got interface{}) bool { return c.Pass(got.([]byte)) }, // TODO: check assertion
		c.Explain,
	)
}

func FromString(c check.StringChecker) check.UntypedChecker {
	return check.NewUntypedCheck(
		func(got interface{}) bool { return c.Pass(got.(string)) }, // TODO: check assertion
		c.Explain,
	)
}

func FromInt(c check.IntChecker) check.UntypedChecker {
	return check.NewUntypedCheck(
		func(got interface{}) bool { return c.Pass(got.(int)) }, // TODO: check assertion
		c.Explain,
	)
}

func FromDuration(c check.DurationChecker) check.UntypedChecker {
	return check.NewUntypedCheck(
		func(got interface{}) bool { return c.Pass(got.(time.Duration)) }, // TODO: check assertion
		c.Explain,
	)
}

func FromHTTPHeader(c check.HTTPHeaderChecker) check.UntypedChecker {
	return check.NewUntypedCheck(
		func(got interface{}) bool { return c.Pass(got.(http.Header)) }, // TODO: check assertion
		c.Explain,
	)
}

// func UntypedChecker(checker interface{}) check.UntypedChecker {
// 	pf := untypedPassFunc(checker)
// }

func untypedPassFunc(c interface{}) check.UntypedPassFunc {
	switch c.(type) {
	case check.BytesPassFunc:
		return untypedBytesPassFunc(c.(check.BytesChecker))

	case check.StringPassFunc:
		return untypedStringPassFunc(c.(check.StringChecker))

	case check.IntPassFunc:
		return untypedIntPassFunc(c.(check.IntChecker))

	case check.DurationPassFunc:
		return untypedDurationPassFunc(c.(check.DurationChecker))

	default:
		log.Fatal("bad conversion")
		return nil
	}
}

func untypedBytesPassFunc(c check.BytesChecker) check.UntypedPassFunc {
	return func(got interface{}) bool { return c.Pass(got.([]byte)) }
}

func untypedIntPassFunc(c check.IntChecker) check.UntypedPassFunc {
	return func(got interface{}) bool { return c.Pass(got.(int)) }
}

func untypedStringPassFunc(c check.StringChecker) check.UntypedPassFunc {
	return func(got interface{}) bool { return c.Pass(got.(string)) }
}

func untypedDurationPassFunc(c check.DurationChecker) check.UntypedPassFunc {
	return func(got interface{}) bool { return c.Pass(got.(time.Duration)) }
}

func PassFunc(passFunc interface{}) check.UntypedPassFunc {
	funcType := reflect.TypeOf(passFunc)
	if funcKind := funcType.Kind(); funcKind != reflect.Func {
		log.Fatalf("WrapPassFunc expects func(Type) bool, got %s", funcKind.String())
	}
	// gotType := funcType.In(0)
	// gotValue := reflect.ValueOf(gotType)

	// untypedPassFuncType := reflect.FuncOf(
	// 	[]reflect.Type{gotType},
	// 	[]reflect.Type{reflect.TypeOf(true)},
	// 	false,
	// )

	// h := http.Header{}
	wrap := func(args []reflect.Value) []reflect.Value {
		// gotRaw := args[0]
		// gotVal := reflect.ValueOf(gotRaw)

		header := http.Header{"Content-Length": []string{""}}
		headerValue := reflect.ValueOf(header)

		// convGot := *(*check.UntypedPassFunc)(unsafe.Pointer(&rawGot))
		return []reflect.Value{
			reflect.
				ValueOf(passFunc).
				Call([]reflect.Value{headerValue})[0],
		}
	}

	// f := reflect.MakeFunc(untypedPassFuncType, wrap)

	makeWrap := func(fptr interface{}) {
		fn := reflect.ValueOf(fptr).Elem()
		v := reflect.MakeFunc(fn.Type(), wrap)
		fn.Set(v)
	}
	var finalPassFunc func(interface{}) bool
	makeWrap(&finalPassFunc)
	// return finalPassFunc

	return func(got interface{}) bool {
		return finalPassFunc(got.(http.Header))
	}
}
