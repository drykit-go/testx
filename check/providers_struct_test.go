package check_test

import (
	"testing"

	"github.com/drykit-go/testx/check"
	"github.com/drykit-go/testx/check/checkconv"
)

type (
	personStruct struct {
		Name string
		Age  int
	}
	personMap map[string]interface{}
)

func TestStructCheckerProvider(t *testing.T) {
	p := personStruct{Name: "Marcel Patulacci", Age: 42}

	t.Run("SameJSON pass", func(t *testing.T) {
		m := personMap{"Name": "Marcel Patulacci", "Age": 42}
		c := check.Struct.SameJSON(m)
		assertPassStructChecker(t, "SameJSON", c, p)
	})

	t.Run("SameJSON fail", func(t *testing.T) {
		m := personMap{"Name": "Robert Robichet", "Age": 42}
		c := check.Struct.SameJSON(m)
		assertFailStructChecker(t, "SameJSON", c, p)
	})

	t.Run("IsZero pass", func(t *testing.T) {
		p := personStruct{}
		c := check.Struct.IsZero()
		assertPassStructChecker(t, "IsZero", c, p)
	})

	t.Run("IsZero fail", func(t *testing.T) {
		p := personStruct{Age: -1}
		c := check.Struct.IsZero()
		assertFailStructChecker(t, "IsZero", c, p)
	})

	t.Run("NotZero pass", func(t *testing.T) {
		p := personStruct{Age: -1}
		c := check.Struct.NotZero()
		assertPassStructChecker(t, "NotZero", c, p)
	})

	t.Run("NotZero fail", func(t *testing.T) {
		p := personStruct{}
		c := check.Struct.NotZero()
		assertFailStructChecker(t, "NotZero", c, p)
	})

	t.Run("FieldsEqual pass", func(t *testing.T) {
		c := check.Struct.FieldsEqual("Marcel Patulacci", []string{"Name"})
		assertPassStructChecker(t, "FieldsEqual", c, p)
	})

	t.Run("FieldsEqual fail", func(t *testing.T) {
		c := check.Struct.FieldsEqual("Jean-Pierre Avidol", []string{"Name"})
		assertFailStructChecker(t, "FieldsEqual", c, p)
	})

	t.Run("CheckFields pass", func(t *testing.T) {
		c := check.Struct.CheckFields(
			checkconv.FromInt(check.Int.InRange(41, 43)),
			[]string{"Age"},
		)
		assertPassStructChecker(t, "CheckFields", c, p)
	})

	t.Run("CheckFields fail", func(t *testing.T) {
		c := check.Struct.CheckFields(
			checkconv.FromInt(check.Int.OutRange(41, 43)),
			[]string{"Age"},
		)
		assertFailStructChecker(t, "CheckFields", c, p)
	})
}

// Helpers

func assertPassStructChecker(t *testing.T, method string, c check.ValueChecker, p personStruct) {
	t.Helper()
	if !c.Pass(p) {
		failStructCheckerTest(t, true, method, p, c.Explain)
	}
}

func assertFailStructChecker(t *testing.T, method string, c check.ValueChecker, p personStruct) {
	t.Helper()
	if c.Pass(p) {
		failStructCheckerTest(t, false, method, p, c.Explain)
	}
}

func failStructCheckerTest(t *testing.T, expPass bool, method string, p personStruct, explain check.ExplainFunc) {
	t.Helper()
	failCheckerTest(t, expPass, "Struct."+method, explain("struct value", p))
}
