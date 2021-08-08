package testx

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/drykit-go/testx/check"
)

var _ HandlerTestRunner = (*handlerTestRunner)(nil)

type handlerTestRunner struct {
	baseTest

	hf http.HandlerFunc
	rr *httptest.ResponseRecorder
	rq *http.Request

	response struct {
		header http.Header
		status string
		code   int
		body   []byte
	}

	hasDurationCheck bool
	handlingDuration time.Duration
}

type Runner interface {
	Run(t *testing.T)
}

type HandlerTestRunner interface {
	Runner
	ResponseHeader(...check.HTTPHeaderChecker) HandlerTestRunner
	ResponseStatus(...check.StringChecker) HandlerTestRunner
	ResponseCode(...check.IntChecker) HandlerTestRunner
	ResponseBody(...check.BytesChecker) HandlerTestRunner
	Duration(...check.DurationChecker) HandlerTestRunner
}

func (test *handlerTestRunner) Run(t *testing.T) {
	main := func() { test.hf(test.rr, test.rq) }
	if test.hasDurationCheck {
		test.handlingDuration = timeFunc(main)
	} else {
		main()
	}

	test.setResponse(test.rr)
	test.run(t)
}

func (test *handlerTestRunner) setResponse(rr *httptest.ResponseRecorder) {
	result := rr.Result()
	test.response.header = rr.Header()
	test.response.status = result.Status
	test.response.code = result.StatusCode
	test.response.body = mustReadIO("response body", result.Body)
}

func (test *handlerTestRunner) ResponseStatus(checks ...check.StringChecker) HandlerTestRunner {
	test.addStringChecks(
		"response status",
		func() gotType { return test.response.status },
		checks,
	)
	return test
}

func (test *handlerTestRunner) ResponseCode(checks ...check.IntChecker) HandlerTestRunner {
	test.addIntChecks(
		"response code",
		func() gotType { return test.response.code },
		checks,
	)
	return test
}

func (test *handlerTestRunner) ResponseBody(checks ...check.BytesChecker) HandlerTestRunner {
	test.addBytesChecks(
		"response body",
		func() gotType { return test.response.body },
		checks,
	)
	return test
}

func (test *handlerTestRunner) Duration(checks ...check.DurationChecker) HandlerTestRunner {
	test.hasDurationCheck = true
	test.addDurationChecks(
		"handling duration",
		func() gotType { return test.handlingDuration },
		checks,
	)
	return test
}

func (test *handlerTestRunner) ResponseHeader(checks ...check.HTTPHeaderChecker) HandlerTestRunner {
	test.addHTTPHeaderChecks(
		"response header",
		func() gotType { return test.response.header },
		checks,
	)
	return test
}

func HandlerFunc(hf http.HandlerFunc, r *http.Request) HandlerTestRunner {
	return &handlerTestRunner{
		hf: hf,
		rr: httptest.NewRecorder(),
		rq: r,
	}
}

func Handler(h http.Handler, r *http.Request) HandlerTestRunner {
	return HandlerFunc(h.ServeHTTP, r)
}
