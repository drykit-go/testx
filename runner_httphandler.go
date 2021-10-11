package testx

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/drykit-go/testx/check"
	"github.com/drykit-go/testx/checkconv"
	"github.com/drykit-go/testx/internal/httputil/middleware"
	"github.com/drykit-go/testx/internal/ioutil"
)

var _ HTTPHandlerRunner = (*httpHandlerRunner)(nil)

type httpHandlerRunner struct {
	baseRunner

	hf http.HandlerFunc
	mw func(http.HandlerFunc) http.HandlerFunc
	rr *httptest.ResponseRecorder
	rq *http.Request

	gotRequest  *http.Request
	gotResponse *http.Response
	gotDuration time.Duration
}

func (r *httpHandlerRunner) WithRequest(request *http.Request) HTTPHandlerRunner {
	return &httpHandlerRunner{
		baseRunner: r.baseRunner,
		hf:         r.hf,
		mw:         r.mw,
		rr:         httptest.NewRecorder(),
		rq:         request,
	}
}

func (r *httpHandlerRunner) Duration(checks ...check.DurationChecker) HTTPHandlerRunner {
	r.addDurationChecks(
		"handling duration",
		func() gottype { return r.gotDuration },
		checks,
	)
	return r
}

func (r *httpHandlerRunner) Request(checkers ...check.HTTPRequestChecker) HTTPHandlerRunner {
	for _, c := range checkers {
		r.addCheck(baseCheck{
			label:   "http request",
			get:     func() gottype { return r.gotRequest },
			checker: checkconv.FromHTTPRequest(c),
		})
	}
	return r
}

func (r *httpHandlerRunner) Response(checkers ...check.HTTPResponseChecker) HTTPHandlerRunner {
	for _, c := range checkers {
		r.addCheck(baseCheck{
			label:   "http response",
			get:     func() gottype { return r.gotResponse },
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
	return handlerResults{
		baseResults: r.baseResults(),
		duration:    r.gotDuration,
		response:    r.gotResponse,
	}
}

func (r *httpHandlerRunner) setResults() {
	r.rr = httptest.NewRecorder()
	if r.rq == nil {
		r.rq = r.defaultRequest()
	}

	handler := r.mw(r.interceptRequest(r.hf))
	r.gotDuration = timeFunc(func() {
		handler(r.rr, r.rq)
	})
	r.gotResponse = r.rr.Result() //nolint:bodyclose
	r.gotResponse.Header = r.rr.Header()
}

func (r *httpHandlerRunner) defaultRequest() *http.Request {
	req, _ := http.NewRequest("GET", "/", nil)
	return req
}

func newHTTPHandlerRunner(
	hf http.HandlerFunc,
	middlewares ...func(http.HandlerFunc) http.HandlerFunc,
) HTTPHandlerRunner {
	runner := &httpHandlerRunner{hf: hf}
	runner.setMergedMiddlewares(middlewares...)
	return runner
}

func (r *httpHandlerRunner) interceptRequest(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		r.gotRequest = req.Clone(req.Context())
		next(w, req)
	}
}

func (r *httpHandlerRunner) setMergedMiddlewares(middlewares ...func(http.HandlerFunc) http.HandlerFunc) {
	r.mw = middleware.MergeRight(middlewares...)
}

func Adapt(hf http.HandlerFunc, adapters ...func(http.HandlerFunc) http.HandlerFunc) http.HandlerFunc {
	for _, adapter := range adapters {
		hf = adapter(hf)
	}
	return hf
}

type handlerResults struct {
	baseResults
	duration time.Duration
	response *http.Response
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
	return ioutil.NopRead(&r.response.Body)
}

func (r handlerResults) ResponseDuration() time.Duration {
	return r.duration
}
