package testx

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/drykit-go/testx/check"
	"github.com/drykit-go/testx/checkconv"
	"github.com/drykit-go/testx/internal/ioutil"
)

var _ HTTPHandlerRunner = (*handlerRunner)(nil)

type handlerRunner struct {
	baseRunner

	hf http.HandlerFunc
	rr *httptest.ResponseRecorder
	rq *http.Request

	response *http.Response

	duration time.Duration
}

func (r *handlerRunner) Run(t *testing.T) {
	t.Helper()
	r.dryRun()
	r.run(t)
}

func (r *handlerRunner) DryRun() HandlerResulter {
	r.dryRun()
	return handlerResults{
		baseResults:  r.baseResults(),
		duration:     r.duration,
		response:     r.response,
		responseBody: ioutil.NopRead(&r.response.Body),
	}
}

func (r *handlerRunner) dryRun() {
	main := func() { r.hf(r.rr, r.rq) }
	r.duration = timeFunc(main)
	r.response = r.rr.Result() //nolint:bodyclose
	r.response.Header = r.rr.Header()
}

func (r *handlerRunner) Duration(checks ...check.DurationChecker) HTTPHandlerRunner {
	r.addDurationChecks(
		"handling duration",
		func() gottype { return r.duration },
		checks,
	)
	return r
}

func (r *handlerRunner) Request(checkers ...check.HTTPRequestChecker) HTTPHandlerRunner {
	for _, c := range checkers {
		r.addCheck(baseCheck{
			label:   "http request",
			get:     func() gottype { return r.rq },
			checker: checkconv.FromHTTPRequest(c),
		})
	}
	return r
}

func (r *handlerRunner) Response(checkers ...check.HTTPResponseChecker) HTTPHandlerRunner {
	for _, c := range checkers {
		r.addCheck(baseCheck{
			label:   "http response",
			get:     func() gottype { return r.response },
			checker: checkconv.FromHTTPResponse(c),
		})
	}
	return r
}

func newHandlerRunner(hf http.HandlerFunc, r *http.Request) HTTPHandlerRunner {
	return &handlerRunner{
		hf: hf,
		rr: httptest.NewRecorder(),
		rq: r,
	}
}

type handlerResults struct {
	baseResults
	duration     time.Duration
	response     *http.Response
	responseBody []byte
}

var _ HandlerResulter = (*handlerResults)(nil)

func (r handlerResults) ResponseHeader() http.Header {
	return r.response.Header
}

func (r handlerResults) ResponseStatus() string {
	return r.response.Status
}

func (r handlerResults) ResponseCode() int {
	return r.response.StatusCode
}

func (r handlerResults) ResponseBody() []byte {
	return r.responseBody
}

func (r handlerResults) ResponseDuration() time.Duration {
	return r.duration
}
