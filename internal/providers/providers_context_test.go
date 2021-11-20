package providers_test

import (
	"context"
	"testing"

	"github.com/drykit-go/testx/check"
	"github.com/drykit-go/testx/internal/providers"
)

func TestContextCheckerProvider(t *testing.T) {
	checkContext := providers.ContextCheckerProvider{}

	ctxNotDone := func() (context.Context, context.CancelFunc) {
		return context.WithCancel(context.Background())
	}
	ctxDone := func() context.Context {
		ctx, cancel := ctxNotDone()
		cancel()
		return ctx
	}
	ctxVal := func(key, val any) context.Context {
		return context.WithValue(context.Background(), key, val)
	}

	t.Run("Done pass", func(t *testing.T) {
		checkDone := checkContext.Done(true)
		ctxDone := ctxDone()
		assertPassChecker(t, "Context.Done", checkDone, ctxDone)

		checkNotDone := checkContext.Done(false)
		ctxNotDone, cancel := ctxNotDone()
		assertPassChecker(t, "Context.Done", checkNotDone, ctxNotDone)
		cancel()
	})

	t.Run("Done fail", func(t *testing.T) {
		checkDone := checkContext.Done(true)
		ctxNotDone, cancel := ctxNotDone()
		assertFailChecker(t, "Context.Done", checkDone, ctxNotDone, makeExpl(
			"context to be done",
			"context not done",
		))
		cancel()

		checkNotDone := checkContext.Done(false)
		ctxDone := ctxDone()
		assertFailChecker(t, "Context.Done", checkNotDone, ctxDone, makeExpl(
			"context not to be done",
			"context canceled",
		))
	})

	t.Run("HasKeys pass", func(t *testing.T) {
		c := checkContext.HasKeys("user")
		ctx := ctxVal("user", struct{}{})
		assertPassChecker(t, "Context.HasKeys", c, ctx)
	})

	t.Run("HasKeys fail", func(t *testing.T) {
		c := checkContext.HasKeys("secret", "user", "token")
		ctx := ctxVal("user", struct{}{})
		assertFailChecker(t, "Context.HasKeys", c, ctx, makeExpl(
			"to have keys [secret, token]",
			"keys not set",
		))
	})

	t.Run("Value pass", func(t *testing.T) {
		c := checkContext.Value("userID", check.Wrap(check.Int.GT(0)))
		ctx := ctxVal("userID", 42)
		assertPassChecker(t, "Context.Value", c, ctx)
	})

	t.Run("Value fail", func(t *testing.T) {
		c := checkContext.Value("userID", check.Value[any]().Is(0))

		ctxMissingKey := context.Background()
		assertFailChecker(t, "Context.Value", c, ctxMissingKey, makeExpl(
			"value for key userID to pass Checker[any]",
			"explanation: value:\n"+makeExpl("0", "<nil>"),
		))

		ctxBadValue := ctxVal("userID", -1)
		assertFailChecker(t, "Context.Value", c, ctxBadValue, makeExpl(
			"value for key userID to pass Checker[any]",
			"explanation: value:\n"+makeExpl("0", "-1"),
		))
	})
}
