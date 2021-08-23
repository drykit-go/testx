package testx

import (
	"net/http"
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
		DryRun() TableResulter
	}

	HandlerRunner interface {
		Runner
		ResponseHeader(...check.HTTPHeaderChecker) HandlerRunner
		ResponseStatus(...check.StringChecker) HandlerRunner
		ResponseCode(...check.IntChecker) HandlerRunner
		ResponseBody(...check.BytesChecker) HandlerRunner
		Duration(...check.DurationChecker) HandlerRunner
		DryRun() HandlerResulter
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

type HandlerResulter interface {
	Resulter
	ResponseHeader() http.Header
	ResponseStatus() string
	ResponseCode() int
	ResponseBody() []byte
	ResponseDuration() time.Duration
}

type TableResulter interface {
	Resulter
	// PassedAt returns true if the ith test case passed.
	PassedAt(index int) bool
	// PassedAt returns true if the ith test case failed.
	FailedAt(index int) bool
	// PassedAt returns true if the test case with matching label passed.
	PassedLabel(label string) bool
	// PassedAt returns true if the test case with matching label failed.
	FailedLabel(label string) bool
}

type CheckResult struct {
	Passed bool
	Reason string
	label  string
}
