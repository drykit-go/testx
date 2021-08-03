package testix

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/drykit-go/testix/check"
	"github.com/drykit-go/testix/check/checkconv"
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
	ResponseHeader(check.HTTPHeaderChecker) HandlerTester
	ResponseStatus(check.StringChecker) HandlerTester
	ResponseCode(check.IntChecker) HandlerTester
	ResponseBody(check.BytesChecker) HandlerTester
	Duration(check.DurationChecker) HandlerTester
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

func (test *HandlerFuncTest) ResponseStatus(c check.StringChecker) HandlerTester {
	test.addCheck(testCheck{
		label: "response status",
		check: checkconv.FromString(c),
		get:   func() gotType { return test.response.status },
	})
	return test
}

func (test *HandlerFuncTest) ResponseCode(c check.IntChecker) HandlerTester {
	test.addCheck(testCheck{
		label: "response code",
		check: checkconv.FromInt(c),
		get:   func() gotType { return test.response.code },
	})
	return test
}

func (test *HandlerFuncTest) ResponseBody(c check.BytesChecker) HandlerTester {
	test.addCheck(testCheck{
		label: "response body",
		check: checkconv.FromBytes(c),
		get:   func() gotType { return test.response.body },
	})
	return test
}

func (test *HandlerFuncTest) Duration(c check.DurationChecker) HandlerTester {
	test.hasDurationCheck = true
	test.addCheck(testCheck{
		label: "handling duration",
		check: checkconv.FromDuration(c),
		get:   func() gotType { return test.handlingDuration },
	})
	return test
}

// FIXME: test.response.ContentLength / test.rr.Result().ContentLength is always -1
// func (test *HandlerFuncTest) ContentLength(c check.IntChecker) HandlerTester {
// 	test.addCheck(testCheck{
// 		label: "content length",
// 		check: checkconv.FromInt(c),
// 		get:   func() gotType { return test.contentLength },
// 	})
// 	return test
// }

func (test *HandlerFuncTest) ResponseHeader(c check.HTTPHeaderChecker) HandlerTester {
	test.addCheck(testCheck{
		label: "response header",
		check: checkconv.FromHTTPHeader(c),
		get:   func() gotType { return test.response.header },
	})
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
