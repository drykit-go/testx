package testx

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/drykit-go/testx/check"
)

var _ HandlerRunner = (*handlerRunner)(nil)

type httpResponse struct {
	header   http.Header
	status   string
	code     int
	body     []byte
	duration time.Duration
}

type handlerRunner struct {
	baseRunner

	hf http.HandlerFunc
	rr *httptest.ResponseRecorder
	rq *http.Request

	response httpResponse

	hasDurationCheck bool
}

func (r *handlerRunner) Run(t *testing.T) {
	r.dryRun()
	r.run(t)
}

func (r *handlerRunner) DryRun() HandlerResulter {
	r.dryRun()
	return handlerResults{
		baseResults: r.baseResults(),
		response:    r.response,
	}
}

func (r *handlerRunner) dryRun() {
	main := func() { r.hf(r.rr, r.rq) }
	duration := timeFunc(main)
	r.setResponse(r.rr, duration)
}

func (r *handlerRunner) setResponse(rr *httptest.ResponseRecorder, d time.Duration) {
	result := rr.Result()
	defer result.Body.Close()
	r.response.header = rr.Header()
	r.response.status = result.Status
	r.response.code = result.StatusCode
	r.response.body = mustReadIO("response body", result.Body)
	r.response.duration = d
}

func (r *handlerRunner) ResponseStatus(checks ...check.StringChecker) HandlerRunner {
	r.addStringChecks(
		"response status",
		func() gottype { return r.response.status },
		checks,
	)
	return r
}

func (r *handlerRunner) ResponseCode(checks ...check.IntChecker) HandlerRunner {
	r.addIntChecks(
		"response code",
		func() gottype { return r.response.code },
		checks,
	)
	return r
}

func (r *handlerRunner) ResponseBody(checks ...check.BytesChecker) HandlerRunner {
	r.addBytesChecks(
		"response body",
		func() gottype { return r.response.body },
		checks,
	)
	return r
}

func (r *handlerRunner) Duration(checks ...check.DurationChecker) HandlerRunner {
	r.hasDurationCheck = true
	r.addDurationChecks(
		"handling duration",
		func() gottype { return r.response.duration },
		checks,
	)
	return r
}

func (r *handlerRunner) ResponseHeader(checks ...check.HTTPHeaderChecker) HandlerRunner {
	r.addHTTPHeaderChecks(
		"response header",
		func() gottype { return r.response.header },
		checks,
	)
	return r
}

func HandlerFunc(hf http.HandlerFunc, r *http.Request) HandlerRunner {
	return &handlerRunner{
		hf: hf,
		rr: httptest.NewRecorder(),
		rq: r,
	}
}

func Handler(h http.Handler, r *http.Request) HandlerRunner {
	return HandlerFunc(h.ServeHTTP, r)
}

type handlerResults struct {
	baseResults
	response httpResponse
}

var _ HandlerResulter = (*handlerResults)(nil)

func (r handlerResults) ResponseHeader() http.Header {
	return r.response.header
}

func (r handlerResults) ResponseStatus() string {
	return r.response.status
}

func (r handlerResults) ResponseCode() int {
	return r.response.code
}

func (r handlerResults) ResponseBody() []byte {
	return r.response.body
}

func (r handlerResults) ResponseDuration() time.Duration {
	return r.response.duration
}
