package testx

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/drykit-go/testx/check"
	"github.com/drykit-go/testx/checkconv"
	"github.com/drykit-go/testx/internal/httpconv"
	"github.com/drykit-go/testx/internal/ioutil"
)

var _ HTTPHandlerRunner = (*httpHandlerRunner)(nil)

type httpHandlerRunner struct {
	baseRunner

	in  httpHandlerRunnerInput
	got httpHandlerRunnerResults
}

func (r *httpHandlerRunner) WithRequest(request *http.Request) HTTPHandlerRunner {
	return &httpHandlerRunner{
		baseRunner: r.baseRunner,
		in:         r.in.withRequest(request),
	}
}

func (r *httpHandlerRunner) Duration(checkers ...check.DurationChecker) HTTPHandlerRunner {
	for _, c := range checkers {
		r.addCheck(baseCheck{
			label:   "handling duration",
			get:     func() gottype { return r.got.duration },
			checker: checkconv.FromDuration(c),
		})
	}
	return r
}

func (r *httpHandlerRunner) Request(checkers ...check.HTTPRequestChecker) HTTPHandlerRunner {
	for _, c := range checkers {
		r.addCheck(baseCheck{
			label:   "http request",
			get:     func() gottype { return r.got.request },
			checker: checkconv.FromHTTPRequest(c),
		})
	}
	return r
}

func (r *httpHandlerRunner) Response(checkers ...check.HTTPResponseChecker) HTTPHandlerRunner {
	for _, c := range checkers {
		r.addCheck(baseCheck{
			label:   "http response",
			get:     func() gottype { return r.got.response },
			checker: checkconv.FromHTTPResponse(c),
		})
	}
	return r
}

func (r *httpHandlerRunner) Run(t *testing.T) {
	t.Helper()
	r.setResults()
	r.run(t)
}

func (r *httpHandlerRunner) DryRun() HandlerResulter {
	r.setResults()
	results := r.got
	results.baseResults = r.dryRun()
	return results
}

func (r *httpHandlerRunner) setResults() {
	rr := httptest.NewRecorder()
	if r.in.rq == nil {
		r.in.rq = r.defaultRequest()
	}

	handler := r.in.mw(r.interceptRequest(r.in.hf))
	r.got.duration = timeFunc(func() {
		handler(rr, r.in.rq)
	})
	r.got.response = rr.Result() //nolint:bodyclose
	r.got.response.Header = rr.Header()
}

func (r *httpHandlerRunner) defaultRequest() *http.Request {
	req, _ := http.NewRequest("GET", "/", nil)
	return req
}

func newHTTPHandlerRunner(
	hf http.HandlerFunc,
	middlewares ...func(http.HandlerFunc) http.HandlerFunc,
) HTTPHandlerRunner {
	runner := &httpHandlerRunner{in: httpHandlerRunnerInput{hf: hf}}
	runner.setMergedMiddlewares(middlewares...)
	return runner
}

func (r *httpHandlerRunner) interceptRequest(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		r.got.request = req.Clone(req.Context())
		next(w, req)
	}
}

func (r *httpHandlerRunner) setMergedMiddlewares(middlewares ...func(http.HandlerFunc) http.HandlerFunc) {
	r.in.mw = httpconv.Merge(middlewares...)
}

type httpHandlerRunnerInput struct {
	hf http.HandlerFunc
	mw func(http.HandlerFunc) http.HandlerFunc
	rq *http.Request
}

func (in httpHandlerRunnerInput) withRequest(rq *http.Request) httpHandlerRunnerInput {
	return httpHandlerRunnerInput{hf: in.hf, mw: in.mw, rq: rq}
}

type httpHandlerRunnerResults struct {
	baseResults
	request  *http.Request
	response *http.Response
	duration time.Duration
}

var _ HandlerResulter = (*httpHandlerRunnerResults)(nil)

func (res httpHandlerRunnerResults) ResponseHeader() http.Header {
	return res.response.Header
}

func (res httpHandlerRunnerResults) ResponseStatus() string {
	return res.response.Status
}

func (res httpHandlerRunnerResults) ResponseCode() int {
	return res.response.StatusCode
}

func (res httpHandlerRunnerResults) ResponseBody() []byte {
	return ioutil.NopRead(&res.response.Body)
}

func (res httpHandlerRunnerResults) ResponseDuration() time.Duration {
	return res.duration
}
