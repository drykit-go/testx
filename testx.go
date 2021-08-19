package testx

import (
	"testing"
	"time"

	"github.com/drykit-go/testx/check"
)

type (
	Runner interface {
		Run(t *testing.T)
	}

	ValueRunner interface {
		Runner
		MustBe(v interface{}) ValueRunner
		MustNotBe(v ...interface{}) ValueRunner
		MustPass(checker ...interface{}) ValueRunner
		DryRun() Resulter
	}

	TableRunner interface {
		Runner
		Cases(cases []Case) TableRunner
	}

	HandlerRunner interface {
		Runner
		ResponseHeader(...check.HTTPHeaderChecker) HandlerRunner
		ResponseStatus(...check.StringChecker) HandlerRunner
		ResponseCode(...check.IntChecker) HandlerRunner
		ResponseBody(...check.BytesChecker) HandlerRunner
		Duration(...check.DurationChecker) HandlerRunner
	}
)

// Resulter is the base interface for runners results.
type Resulter interface {
	Checks() []CheckResult
	Passed() bool
	Failed() bool
	NChecks() int
	NPassed() int
	NFailed() int
	ExecTime() time.Duration
}

type CheckResult struct {
	Passed bool
	Reason string
}
