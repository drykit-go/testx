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

var _ HTTPHandlerRunner = (*handlerRunner)(nil)

type handlerRunner struct {
	baseRunner

	hf http.HandlerFunc
	mw func(http.HandlerFunc) http.HandlerFunc
	rr *httptest.ResponseRecorder
	rq *http.Request

	gotRequest *http.Request
	response   *http.Response
	duration   time.Duration
}

func (r *handlerRunner) WithRequest(request *http.Request) HTTPHandlerRunner {
	return &handlerRunner{
		baseRunner: r.baseRunner,
		hf:         r.hf,
		mw:         r.mw,
		rr:         httptest.NewRecorder(),
		rq:         request,
	}
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
			get:     func() gottype { return r.gotRequest },
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

func (r *handlerRunner) Run(t *testing.T) {
	t.Helper()
	r.dryRun()
	r.run(t)
}

func (r *handlerRunner) DryRun() HandlerResulter {
	r.dryRun()
	return handlerResults{
		baseResults: r.baseResults(),
		duration:    r.duration,
		response:    r.response,
	}
}

func (r *handlerRunner) dryRun() {
	r.rr = httptest.NewRecorder()
	if r.rq == nil {
		r.rq = r.defaultRequest()
	}

	hf := r.mw(r.interceptRequest(r.hf))
	r.hf = hf

	main := func() {
		hf(r.rr, r.rq)
	}

	r.duration = timeFunc(main)
	r.response = r.rr.Result() //nolint:bodyclose
	r.response.Header = r.rr.Header()
}

func (r *handlerRunner) defaultRequest() *http.Request {
	req, _ := http.NewRequest("GET", "/", nil)
	return req
}

func newHandlerRunner(
	hf http.HandlerFunc,
	middlewares ...func(http.HandlerFunc) http.HandlerFunc,
) HTTPHandlerRunner {
	runner := &handlerRunner{hf: hf}
	runner.setMiddleware(middlewares...)
	return runner
}

func (r *handlerRunner) interceptRequest(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		r.gotRequest = req.Clone(req.Context())
		next(w, req)
	}
}

func (r *handlerRunner) setMiddleware(middlewares ...func(http.HandlerFunc) http.HandlerFunc) {
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
