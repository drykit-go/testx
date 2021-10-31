package check

import (
	"testing"

	"github.com/drykit-go/testx/internal/testutil"
)

func TestBaseCheckerProvider_sameJSON(t *testing.T) {
	var (
		borig = []byte(`{"id":42,"name":"Marcel Patulacci"}`)
		bdiff = []byte(`{"id":41,"name":"Robert Robichet"}`)
		bsame = []byte("{\n  \"id\": 42,\n  \"name\": \"Marcel Patulacci\"\n}")
	)

	t.Run("bad dst type", func(t *testing.T) {
		var x map[string]any
		var y int
		defer testutil.AssertPanicMessage(t,
			"json: cannot unmarshal object into Go value of type int",
		)
		if (baseCheckerProvider{}).sameJSON(borig, bsame, &x, &y) {
			t.Error("exp false, got true")
		}
	})

	t.Run("same json", func(t *testing.T) {
		var x, y map[string]any
		if !(baseCheckerProvider{}).sameJSON(borig, bsame, &x, &y) {
			t.Error("exp true, got false", x, y)
		}
	})

	t.Run("diff json", func(t *testing.T) {
		var x, y map[string]any
		if (baseCheckerProvider{}).sameJSON(borig, bdiff, &x, &y) {
			t.Error("exp false, got true")
		}
	})
}

func TestBaseCheckerProvider_sameJSONProduced(t *testing.T) {
	var (
		orig = struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		}{ID: 42, Name: "Marcel Patulacci"}
		same = map[string]any{"id": 42, "name": "Marcel Patulacci"}
		diff = map[string]any{"id": 42, "name": "Robert Robichet"}
	)

	t.Run("same json produced", func(t *testing.T) {
		var x, y map[string]any
		if !(baseCheckerProvider{}).sameJSONProduced(orig, same, &x, &y) {
			t.Error("exp true, got false", x, y)
		}
	})

	t.Run("diff json produced", func(t *testing.T) {
		var x, y map[string]any
		if (baseCheckerProvider{}).sameJSONProduced(orig, diff, &x, &y) {
			t.Error("exp false, got true")
		}
	})

	t.Run("bad input", func(t *testing.T) {
		var x, y map[string]any
		badinput := make(chan int)
		defer testutil.AssertPanicMessage(t,
			"json: unsupported type: chan int",
		)
		if (baseCheckerProvider{}).sameJSONProduced(orig, badinput, &x, &y) {
			t.Error("exp false, got true")
		}
	})
}
