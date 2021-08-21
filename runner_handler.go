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

	results handlerResults

	hasDurationCheck bool
}

func (r *handlerRunner) Run(t *testing.T) {
	main := func() { r.hf(r.rr, r.rq) }
	if r.hasDurationCheck {
		r.response.duration = timeFunc(main)
	} else {
		main()
	}

	r.setResponse(r.rr)
	r.run(t)
}

func (r *handlerRunner) DryRun() HandlerResulter {
	main := func() { r.hf(r.rr, r.rq) }
	if r.hasDurationCheck {
		r.response.duration = timeFunc(main)
	} else {
		main()
	}
	r.setResponse(r.rr)

	for _, c := range r.checks {
		r.updateBaseResults(c)
	}
	r.updateSpecificResults()
	return r.results
}

func (r *handlerRunner) updateSpecificResults() {
	r.results = handlerResults{
		baseResults: r.baseResults,
		response:    r.response,
	}
}

func (r *handlerRunner) setResponse(rr *httptest.ResponseRecorder) {
	result := rr.Result()
	defer result.Body.Close()
	r.response.header = rr.Header()
	r.response.status = result.Status
	r.response.code = result.StatusCode
	r.response.body = mustReadIO("response body", result.Body)
}

func (r *handlerRunner) ResponseStatus(checks ...check.StringChecker) HandlerRunner {
	r.addStringChecks(
		"response status",
		func() gotType { return r.response.status },
		checks,
	)
	return r
}

func (r *handlerRunner) ResponseCode(checks ...check.IntChecker) HandlerRunner {
	r.addIntChecks(
		"response code",
		func() gotType { return r.response.code },
		checks,
	)
	return r
}

func (r *handlerRunner) ResponseBody(checks ...check.BytesChecker) HandlerRunner {
	r.addBytesChecks(
		"response body",
		func() gotType { return r.response.body },
		checks,
	)
	return r
}

func (r *handlerRunner) Duration(checks ...check.DurationChecker) HandlerRunner {
	r.hasDurationCheck = true
	r.addDurationChecks(
		"handling duration",
		func() gotType { return r.response.duration },
		checks,
	)
	return r
}

func (r *handlerRunner) ResponseHeader(checks ...check.HTTPHeaderChecker) HandlerRunner {
	r.addHTTPHeaderChecks(
		"response header",
		func() gotType { return r.response.header },
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
