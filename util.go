package testix

import (
	"io"
	"log"
	"time"
)

func timeFunc(f func()) time.Duration {
	t0 := time.Now()
	f()
	return time.Since(t0)
}

func mustReadIO(origin string, reader io.ReadCloser) []byte {
	b, err := io.ReadAll(reader)
	if err != nil {
		log.Fatalf("error reading %s: %s", origin, err)
	}
	return b
}
