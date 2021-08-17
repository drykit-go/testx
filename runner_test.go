package testx_test

import "time"

type baseResults struct {
	Passed, Failed            bool
	NPassed, NFailed, NChecks int
	ExecTime                  time.Duration
}
