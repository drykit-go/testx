package testx

import "time"

type baseResults struct {
	checks   []CheckResult
	nFailed  int
	execTime time.Duration
}

func (r baseResults) Checks() []CheckResult {
	return r.checks
}

func (r baseResults) Passed() bool {
	return r.nFailed == 0
}

func (r baseResults) Failed() bool {
	return !r.Passed()
}

func (r baseResults) NPassed() int {
	return r.NChecks() - r.nFailed
}

func (r baseResults) NFailed() int {
	return r.nFailed
}

func (r baseResults) NChecks() int {
	return len(r.checks)
}

func (r baseResults) ExecTime() time.Duration {
	return r.execTime
}
