package check_test

import (
	"testing"

	"github.com/drykit-go/testx/check"
	"github.com/drykit-go/testx/checkconv"
)

type structTest struct {
	Name string
	Age  int
}

func TestStructCheckerProvider(t *testing.T) {
	s := structTest{Name: "Marcel Patulacci", Age: 42}

	t.Run("SameJSON pass", func(t *testing.T) {
		m := map[string]interface{}{"Name": "Marcel Patulacci", "Age": 42}
		c := check.Struct.SameJSON(m)
		assertPassStructChecker(t, "SameJSON", c, s)
	})

	t.Run("SameJSON fail", func(t *testing.T) {
		m := map[string]interface{}{"Name": "Robert Robichet", "Age": 42}
		c := check.Struct.SameJSON(m)
		assertFailStructChecker(t, "SameJSON", c, s)
	})

	t.Run("IsZero pass", func(t *testing.T) {
		s := structTest{}
		c := check.Struct.IsZero()
		assertPassStructChecker(t, "IsZero", c, s)
	})

	t.Run("IsZero fail", func(t *testing.T) {
		s := structTest{Age: -1}
		c := check.Struct.IsZero()
		assertFailStructChecker(t, "IsZero", c, s)
	})

	t.Run("NotZero pass", func(t *testing.T) {
		s := structTest{Age: -1}
		c := check.Struct.NotZero()
		assertPassStructChecker(t, "NotZero", c, s)
	})

	t.Run("NotZero fail", func(t *testing.T) {
		s := structTest{}
		c := check.Struct.NotZero()
		assertFailStructChecker(t, "NotZero", c, s)
	})

	t.Run("FieldsEqual pass", func(t *testing.T) {
		c := check.Struct.FieldsEqual("Marcel Patulacci", []string{"Name"})
		assertPassStructChecker(t, "FieldsEqual", c, s)
	})

	t.Run("FieldsEqual fail", func(t *testing.T) {
		c := check.Struct.FieldsEqual("Jean-Pierre Avidol", []string{"Name"})
		assertFailStructChecker(t, "FieldsEqual", c, s)
	})

	t.Run("CheckFields pass", func(t *testing.T) {
		c := check.Struct.CheckFields(
			checkconv.FromInt(check.Int.InRange(41, 43)),
			[]string{"Age"},
		)
		assertPassStructChecker(t, "CheckFields", c, s)
	})

	t.Run("CheckFields fail", func(t *testing.T) {
		c := check.Struct.CheckFields(
			checkconv.FromInt(check.Int.OutRange(41, 43)),
			[]string{"Age"},
		)
		assertFailStructChecker(t, "CheckFields", c, s)
	})
}

// Helpers

func assertPassStructChecker(t *testing.T, method string, c check.ValueChecker, s structTest) {
	t.Helper()
	if !c.Pass(s) {
		failStructCheckerTest(t, true, method, s, c.Explain)
	}
}

func assertFailStructChecker(t *testing.T, method string, c check.ValueChecker, s structTest) {
	t.Helper()
	if c.Pass(s) {
		failStructCheckerTest(t, false, method, s, c.Explain)
	}
}

func failStructCheckerTest(t *testing.T, expPass bool, method string, s structTest, explain check.ExplainFunc) {
	t.Helper()
	failCheckerTest(t, expPass, "Struct."+method, explain("struct value", s))
}
