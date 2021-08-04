package testix

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/drykit-go/testix/check"
)

type HandlerTest struct {
	*HandlerFuncTest
}

type HandlerFuncTest struct {
	baseTest
	hf http.HandlerFunc
	rr *httptest.ResponseRecorder
	rq *http.Request

	hasDurationCheck bool
	handlingDuration time.Duration

	response struct {
		header http.Header
		status string
		code   int
		body   []byte
	}
}

type Runner interface {
	Run(t *testing.T)
}

type HandlerTester interface {
	Runner
	// ContentLength(check.IntChecker) HandlerTester
	ResponseHeader(...check.HTTPHeaderChecker) HandlerTester
	ResponseStatus(...check.StringChecker) HandlerTester
	ResponseCode(...check.IntChecker) HandlerTester
	ResponseBody(...check.BytesChecker) HandlerTester
	Duration(...check.DurationChecker) HandlerTester
}

func (test *HandlerFuncTest) Run(t *testing.T) {
	main := func() { test.hf(test.rr, test.rq) }
	if test.hasDurationCheck {
		test.handlingDuration = timeFunc(main)
	} else {
		main()
	}

	test.setResponse(test.rr)
	test.run(t)
}

func (test *HandlerFuncTest) setResponse(rr *httptest.ResponseRecorder) {
	result := rr.Result()
	test.response.header = rr.Header()
	test.response.status = result.Status
	test.response.code = result.StatusCode
	test.response.body = mustReadIO("response body", result.Body)
}

func (test *HandlerFuncTest) ResponseStatus(checks ...check.StringChecker) HandlerTester {
	test.addStringChecks(
		"response status",
		func() gotType { return test.response.status },
		checks,
	)
	return test
}

func (test *HandlerFuncTest) ResponseCode(checks ...check.IntChecker) HandlerTester {
	test.addIntChecks(
		"response code",
		func() gotType { return test.response.code },
		checks,
	)
	return test
}

func (test *HandlerFuncTest) ResponseBody(checks ...check.BytesChecker) HandlerTester {
	test.addBytesChecks(
		"response body",
		func() gotType { return test.response.body },
		checks,
	)
	return test
}

func (test *HandlerFuncTest) Duration(checks ...check.DurationChecker) HandlerTester {
	test.hasDurationCheck = true
	test.addDurationChecks(
		"handling duration",
		func() gotType { return test.handlingDuration },
		checks,
	)
	return test
}

func (test *HandlerFuncTest) ResponseHeader(checks ...check.HTTPHeaderChecker) HandlerTester {
	test.addHTTPHeaderChecks(
		"response header",
		func() gotType { return test.response.header },
		checks,
	)
	return test
}

func HandlerFunc(hf http.HandlerFunc, r *http.Request) HandlerTester {
	return &HandlerFuncTest{
		hf: hf,
		rr: httptest.NewRecorder(),
		rq: r,
	}
}

func Handler(h http.Handler, r *http.Request) HandlerTester {
	return HandlerFunc(h.ServeHTTP, r)
}
