// Code generated by go generate ./...; DO NOT EDIT
// Last generated on 27 Aug 21 01:04 UTC

// Package checkconv provides functions to convert typed checks
// into generic untyped ones.
package checkconv

import (
	"log"
	"net/http"
	"time"

	"github.com/drykit-go/testx/check"
)

// FromBytes returns a new check.UntypedChecker from the given Bytes checker.
// It can be used to facilitate checkers usage by test runners.
func FromBytes(c check.BytesChecker) check.UntypedChecker {
	return check.NewUntypedCheck(
		func(got interface{}) bool { return c.Pass(got.([]byte)) },
		c.Explain,
	)
}

// FromString returns a new check.UntypedChecker from the given String checker.
// It can be used to facilitate checkers usage by test runners.
func FromString(c check.StringChecker) check.UntypedChecker {
	return check.NewUntypedCheck(
		func(got interface{}) bool { return c.Pass(got.(string)) },
		c.Explain,
	)
}

// FromInt returns a new check.UntypedChecker from the given Int checker.
// It can be used to facilitate checkers usage by test runners.
func FromInt(c check.IntChecker) check.UntypedChecker {
	return check.NewUntypedCheck(
		func(got interface{}) bool { return c.Pass(got.(int)) },
		c.Explain,
	)
}

// FromDuration returns a new check.UntypedChecker from the given Duration checker.
// It can be used to facilitate checkers usage by test runners.
func FromDuration(c check.DurationChecker) check.UntypedChecker {
	return check.NewUntypedCheck(
		func(got interface{}) bool { return c.Pass(got.(time.Duration)) },
		c.Explain,
	)
}

// FromHTTPHeader returns a new check.UntypedChecker from the given HTTPHeader checker.
// It can be used to facilitate checkers usage by test runners.
func FromHTTPHeader(c check.HTTPHeaderChecker) check.UntypedChecker {
	return check.NewUntypedCheck(
		func(got interface{}) bool { return c.Pass(got.(http.Header)) },
		c.Explain,
	)
}

// Cast takes a known typed checker (such as check.IntChecker)
// and returns its as a check.UntypedChecker.
// It can be used to facilitate checkers usage by test runners.
// It panics if checker is not a known checker type. For instance,
// a custom checker that implements check.IntChecker will be successfully
// converted, while a valid implementation of an unknown interface,
// such as Float64Checker, will panic.
// For that matter, Assert can be used instead.
func Cast(knownChecker interface{}) check.UntypedChecker {
	switch c := knownChecker.(type) {
	case check.BytesChecker:
		return FromBytes(c)
	case check.StringChecker:
		return FromString(c)
	case check.IntChecker:
		return FromInt(c)
	case check.DurationChecker:
		return FromDuration(c)
	case check.HTTPHeaderChecker:
		return FromHTTPHeader(c)
	case check.UntypedChecker:
		return c
	default:
		log.Panic("attempt to convert unknown checker type")
		return nil
	}
}