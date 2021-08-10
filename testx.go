package testx

import (
	"testing"

	"github.com/drykit-go/testx/check"
)

type Runner interface {
	Run(t *testing.T)
}

type ValueRunner interface {
	Runner
	MustBe(v interface{}) ValueRunner
	MustNotBe(v ...interface{}) ValueRunner
	MustPass(checker ...interface{}) ValueRunner
}

type TableRunner interface {
	Runner
	Cases(cases []Case) TableRunner
}

type HandlerRunner interface {
	Runner
	ResponseHeader(...check.HTTPHeaderChecker) HandlerRunner
	ResponseStatus(...check.StringChecker) HandlerRunner
	ResponseCode(...check.IntChecker) HandlerRunner
	ResponseBody(...check.BytesChecker) HandlerRunner
	Duration(...check.DurationChecker) HandlerRunner
}
