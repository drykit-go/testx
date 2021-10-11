package check_test

import (
	"context"
	"testing"

	"github.com/drykit-go/testx/check"
	"github.com/drykit-go/testx/checkconv"
)

func TestContextCheckerProvider(t *testing.T) {
	ctxNotDone := func() (context.Context, context.CancelFunc) {
		return context.WithCancel(context.Background())
	}
	ctxDone := func() context.Context {
		ctx, cancel := ctxNotDone()
		cancel()
		return ctx
	}
	ctxVal := func(key, val interface{}) context.Context {
		return context.WithValue(context.Background(), key, val)
	}

	t.Run("Done pass", func(t *testing.T) {
		checkDone := check.Context.Done(true)
		ctxDone := ctxDone()
		assertPassContextChecker(t, "Done", checkDone, ctxDone)

		checkNotDone := check.Context.Done(false)
		ctxNotDone, cancel := ctxNotDone()
		assertPassContextChecker(t, "Done", checkNotDone, ctxNotDone)
		cancel()
	})

	t.Run("Done fail", func(t *testing.T) {
		checkDone := check.Context.Done(true)
		ctxNotDone, cancel := ctxNotDone()
		assertFailContextChecker(t, "Done", checkDone, ctxNotDone)
		cancel()

		checkNotDone := check.Context.Done(false)
		ctxDone := ctxDone()
		assertFailContextChecker(t, "Done", checkNotDone, ctxDone)
	})

	t.Run("HasKeys pass", func(t *testing.T) {
		c := check.Context.HasKeys("user")
		ctx := ctxVal("user", struct{}{})
		assertPassContextChecker(t, "HasKeys", c, ctx)
	})

	t.Run("HasKeys fail", func(t *testing.T) {
		c := check.Context.HasKeys("user", "token")
		ctx := ctxVal("user", struct{}{})
		assertFailContextChecker(t, "HasKeys", c, ctx)
	})

	t.Run("Value pass", func(t *testing.T) {
		c := check.Context.Value("userID", checkconv.FromInt(check.Int.GT(0)))
		ctx := ctxVal("userID", 42)
		assertPassContextChecker(t, "Value", c, ctx)
	})

	t.Run("Value fail", func(t *testing.T) {
		c := check.Context.Value("userID", check.Value.Is(0))

		ctxMissingKey := context.Background()
		assertFailContextChecker(t, "Value", c, ctxMissingKey)

		ctxBadValue := ctxVal("userID", -1)
		assertFailContextChecker(t, "Value", c, ctxBadValue)
	})
}

// Helpers

//nolint: revive // context-as-argument rule not relevant here
func assertPassContextChecker(t *testing.T, method string, c check.ContextChecker, ctx context.Context) {
	t.Helper()
	if !c.Pass(ctx) {
		failContextCheckerTest(t, true, method, ctx, c.Explain)
	}
}

//nolint: revive // context-as-argument rule not relevant here
func assertFailContextChecker(t *testing.T, method string, c check.ContextChecker, ctx context.Context) {
	t.Helper()
	if c.Pass(ctx) {
		failContextCheckerTest(t, false, method, ctx, c.Explain)
	}
}

//nolint: revive // context-as-argument rule not relevant here
func failContextCheckerTest(t *testing.T, expPass bool, method string, ctx context.Context, explain check.ExplainFunc) {
	t.Helper()
	failCheckerTest(t, expPass, "Context."+method, explain("Context value", ctx))
}
