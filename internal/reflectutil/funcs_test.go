package reflectutil_test

import (
	"errors"
	"testing"

	"github.com/drykit-go/testx/internal/reflectutil"
	"github.com/drykit-go/testx/internal/testutil"
)

func validFunc(in int) int {
	return in
}

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
		got := reflectutil.FuncName(validFunc)
		exp := "reflectutil_test.validFunc"
		if got != exp {
			t.Errorf("exp %s, got %s", exp, got)
		}
	})
}

func TestNewFunc(t *testing.T) {
	t.Run("non-func returns error", func(t *testing.T) {
		nofunc := struct{}{}
		f, err := reflectutil.NewFunc(nofunc)
		if err == nil {
			t.Error("exp non-nil error, got nil")
		}
		if !errors.Is(err, reflectutil.ErrNotAFunc) {
			t.Errorf("got unexpected error: %s", err)
		}
		if f != nil {
			t.Errorf("exp nil *Func, got %v", f)
		}
	})

	t.Run("valid func returns *Func", func(t *testing.T) {
		f, err := reflectutil.NewFunc(validFunc)
		if err != nil {
			t.Errorf("got unexpected error: %s", err)
		}
		if got, exp := f.Name, reflectutil.FuncName(validFunc); got != exp {
			t.Errorf("f.Name: exp %s, got %s", exp, got)
		}
	})
}
