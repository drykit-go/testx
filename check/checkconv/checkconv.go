// Package checkconv provides functions to convert typed checks
// into generic ones.
package checkconv

import (
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
