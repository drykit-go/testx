package reflectutil_test

import (
	"net/http"
	"testing"

	"github.com/drykit-go/testx/internal/reflectutil"
	"github.com/drykit-go/testx/internal/testutil"
)

func TestFuncName(t *testing.T) {
	t.Run("non-func panics", func(t *testing.T) {
		nofunc := struct{}{}
		defer testutil.AssertPanicMessage(t,
			"calling FuncName with a non func argument (struct {})",
		)
		if got, exp := reflectutil.FuncName(nofunc), ""; got != exp {
			t.Errorf("exp %s, got %s", exp, got)
		}
	})

	t.Run("returns func name", func(t *testing.T) {
		if got, exp := reflectutil.FuncName(http.Get), "http.Get"; got != exp {
			t.Errorf("exp %s, got %s", exp, got)
		}
	})
}
