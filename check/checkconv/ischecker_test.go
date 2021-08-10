package checkconv_test

import (
	"testing"

	"github.com/drykit-go/testx/check/checkconv"
)

type onlyPasser struct{}

func (onlyPasser) Pass(int) bool { return true }

type onlyExplainer struct{}

func (onlyExplainer) PassX(int) bool                     { return true }
func (onlyExplainer) Explain(string, interface{}) string { return "" }

type badPasser struct{}

func (badPasser) Pass(int) int                       { return 0 }
func (badPasser) Explain(string, interface{}) string { return "" }

type badExplainerIn struct{}

func (badExplainerIn) Pass(int) bool                           { return true }
func (badExplainerIn) Explain(interface{}, interface{}) string { return "" }

type badExplainerOut struct{}

func (badExplainerOut) Pass(int) bool                           { return true }
func (badExplainerOut) Explain(string, interface{}) interface{} { return "" }

type checkerAsFields struct {
	Pass    func(int) bool
	Explain func(string, interface{}) string
}

type validCheckerInt struct{}

func (validCheckerInt) Pass(int) bool                      { return true }
func (validCheckerInt) Explain(string, interface{}) string { return "" }

type validCheckerInterface struct{}

func (validCheckerInterface) Pass(interface{}) bool              { return true }
func (validCheckerInterface) Explain(string, interface{}) string { return "" }

func TestIsChecker(t *testing.T) { //nolint: gocognit
	t.Run("should not be considered a checker", func(t *testing.T) {
		values := []interface{}{
			"a string",
			42,
			func(int) bool { return true },
			onlyPasser{},
			onlyExplainer{},
			badPasser{},
			badExplainerIn{},
			badExplainerOut{},
			checkerAsFields{
				Pass:    func(int) bool { return true },
				Explain: func(string, interface{}) string { return "" },
			},
		}

		for _, v := range values {
			if checkconv.IsChecker(v) {
				t.Errorf("value %v was wrongly considered a checker", v)
			}
		}
	})

	t.Run("should be considered a checker", func(t *testing.T) {
		values := []interface{}{
			validCheckerInt{},
			validCheckerInterface{},
		}

		for _, v := range values {
			if !checkconv.IsChecker(v) {
				t.Errorf("checker %v was wrongly considered not a checker", v)
			}
		}
	})
}
