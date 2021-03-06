package check_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/drykit-go/testx"
	"github.com/drykit-go/testx/check"
)

/*
	Example: implementation of a custom checker
*/

// StatusOKChecker is a custom checker that implements IntChecker interface.
// In consequence in can be used in any function that accepts an IntChecker.
type StatusOKChecker struct{}

// Pass satisfies IntPasser interface.
func (c StatusOKChecker) Pass(got int) bool {
	return (got >= 200 && got < 300) || got == 304
}

// Explain satisfies Explainer interface.
func (c StatusOKChecker) Explain(label string, got interface{}) string {
	return fmt.Sprintf("%s: got bad http code: %v", label, got)
}

// HandleNotFound is a http.HandlerFunc that always responds 404 NotFound.
func HandleNotFound(w http.ResponseWriter, _ *http.Request) {
	http.Error(w, "Not Found", 404)
}

func Example_customChecker() {
	request := httptest.NewRequest("GET", "/", nil)

	results := testx.HTTPHandlerFunc(HandleNotFound).WithRequest(request).
		// check.HTTPResponse.StatusCode accepts an IntChecker,
		// StatusOKChecker satisfies IntChecker interface.
		Response(check.HTTPResponse.StatusCode(StatusOKChecker{})).
		DryRun()

	fmt.Println(results.Passed())
	fmt.Println(results.ResponseCode())
	fmt.Println(results.Checks()[0].Reason)

	// Output:
	// false
	// 404
	// http response:
	// exp status code to pass IntChecker
	// got explanation: status code: got bad http code: 404
}
