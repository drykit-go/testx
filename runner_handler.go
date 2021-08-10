package testx

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/drykit-go/testx/check"
)

var _ HandlerRunner = (*handlerRunner)(nil)

type handlerRunner struct {
	baseRunner

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

func (r *handlerRunner) Run(t *testing.T) {
	main := func() { r.hf(r.rr, r.rq) }
	if r.hasDurationCheck {
		r.handlingDuration = timeFunc(main)
	} else {
		main()
	}

	r.setResponse(r.rr)
	r.run(t)
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
		func() gotType { return r.handlingDuration },
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
