package check

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"reflect"
)

type bytesCheckerFactory struct{}

func (bytesCheckerFactory) Is(tar []byte) BytesChecker {
	pass := func(got []byte) bool { return bytes.Equal(got, tar) }
	expl := func(label string, got interface{}) string {
		return fmt.Sprintf(
			"expect %s to equal %v, got %v",
			label, tar, got,
		)
	}
	return NewBytesChecker(pass, expl)
}

func (bytesCheckerFactory) SameJSON(tar []byte) BytesChecker {
	var decGot, decTar interface{}
	pass := func(got []byte) bool {
		if err := json.Unmarshal(got, &decGot); err != nil {
			log.Panic(err)
		}
		if err := json.Unmarshal(tar, &decTar); err != nil {
			log.Panic(err)
		}
		return reflect.DeepEqual(decGot, decTar)
	}
	expl := func(label string, got interface{}) string {
		return fmt.Sprintf(
			"expected json encoded value to equal %#v, got %#v",
			decTar, decGot,
		)
	}
	return NewBytesChecker(pass, expl)
}

func (bytesCheckerFactory) Len(c IntChecker) BytesChecker {
	pass := func(got []byte) bool { return c.Pass(len(got)) }
	expl := func(label string, got interface{}) string {
		return fmt.Sprintf(
			"unexpected %s length: %s",
			label, c.Explain(label, len(got.([]byte))),
		)
	}
	return NewBytesChecker(pass, expl)
}
