package check_test

import (
	"fmt"
	"net/http"

	"github.com/drykit-go/testx"
	"github.com/drykit-go/testx/check"
)

/*
	Example: implementation of a custom checker
*/

// HTTPCodeChecker is a custom checker that implements IntChecker interface.
// In consequence in can be used in any function that requires an IntChecker
// or a less specific checker.
type HTTPCodeChecker struct{}

// Pass satisfies IntPasser interface.
func (c HTTPCodeChecker) Pass(got int) bool {
	return (got <= 200 && got < 300) || got == 304
}

// Explain satisfies Explainer interface.
func (c HTTPCodeChecker) Explain(label string, got interface{}) string {
	return fmt.Sprintf("%s: got bad http code: %v", label, got)
}

// HandleNotFound is a http.HandlerFunc that always responds 404 NotFound.
func HandleNotFound(w http.ResponseWriter, _ *http.Request) {
	http.Error(w, "Not Found", 404)
}

func Example_customChecker() {
	request, _ := http.NewRequest("GET", "/", nil)

	results := testx.HTTPHandlerFunc(HandleNotFound).WithRequest(request).
		// HTTPResponseCode is a valid IntChecker
		Response(check.HTTPResponse.StatusCode(HTTPCodeChecker{})).
		DryRun()

	fmt.Println(results.Passed())
	fmt.Println(results.ResponseCode())
	fmt.Println(results.Checks())

	// Output:
	// false
	// 404
	// [{failed http response:
	// exp status code to pass IntChecker
	// got 404}]
}
