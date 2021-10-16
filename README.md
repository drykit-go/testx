<p align="center">    
  <a href="https://pkg.go.dev/github.com/drykit-go/testx#section-documentation">
    <img alt="Go Reference" src="https://pkg.go.dev/badge/github.com/drykit-go/testx.svg" />
  </a>
  <a href="https://goreportcard.com/report/github.com/drykit-go/testx">
    <img alt="Go Report Card" src="https://goreportcard.com/badge/github.com/drykit-go/testx" />
  </a>
  <a href="https://github.com/drykit-go/testx/releases">
    <img alt="Latest version" src="https://img.shields.io/github/v/tag/drykit-go/testx?label=release">
  </a>
</p>
<p align="center">
  <a href="https://circleci.com/gh/circleci/circleci-docs">
    <img alt="CircleCI" src="https://circleci.com/gh/circleci/circleci-docs.svg?style=shield" />
  </a>
  <a href="https://codecov.io/gh/drykit-go/testx">
    <img alt="Codecov" src="https://codecov.io/gh/drykit-go/testx/branch/main/graph/badge.svg?token=XZRUXJDFJE"/>
  </a>
</p>

<table align="center">
	<tr>
		<td>testx</td>
		<td><a href="./check/README.md">testx/check</a></td>
		<td><a href="./checkconv/README.md">testx/checkconv</a></td>
	</tr>
</table>

# testx

Package `testx` provides test runners to accelerate the writing
of unit tests and reduce boilerplate.

## Table of contents

- [Runners](#runners)
  - [`ValueRunner`](#valuerunner)
  - [`HTTPHandlerRunner`](#httphandlerrunner)
  - [`TableRunner`](#tablerunner)
- [Running tests](#running-tests)
  - [Method `Run`](#method-run)
  - [Method `DryRun`](#method-dryrun)
- [Further documentation](#further-documentation)

## Runners

`testx` provides 3 types of runners:

- `ValueRunner` runs tests on a single value.
- `HTTPHandlerRunner` runs tests on types `http.Handler` and `http.HandlerFunc`.
- `TableRunner` runs tests on a single function with a series of test cases.

### `ValueRunner`

`ValueRunner` runs tests on a single value.

```go
func TestGet42(t *testing.T) {
    testx.Value(Get42()).
        Exp(42).                       // expect 42
        Not(3, "hello").               // expect not 3 nor "hello"
        Pass(checkconv.AssertMany(     // expect to pass input checkers:
            check.Int.InRange(41, 43), //     expect in range [41:43]
            // ...
        )...).
        Run(t)
}
```

More examples in file [example_value_test.go](./example_value_test.go).

### `HTTPHandlerRunner`

`HTTPHandlerRunner` runs tests on a `http.Handler` or `http.HandlerFunc`.
It provides methods to perform checks:
- on the input request (e.g. to ensure it has been attached an expected context
  by some middleware)
- on the written response (status code, body, header...)
- on the execution time.

```go
func TestHandleGetMovieByID(t *testing.T) {
    r, _ := http.NewRequest("GET", "/movies/42", nil)
    // Note: WithRequest can be omitted if the input request is not relevant.
    // In that case it defaults to http.NewRequest("GET", "/", nil).
    testx.HTTPHandlerFunc(HandleGetMovieByID).WithRequest(r).
        Response(
            check.HTTPResponse.StatusCode(check.Int.InRange(200, 299)),
            check.HTTPResponse.Body(check.Bytes.Contains([]byte(`"id":42`))),
        ).
        Duration(check.Duration.Under(10 * time.Millisecond)).
        Run(t)
}
```

More examples in file [example_httphandler_test.go](./example_httphandler_test.go).

### `TableRunner`

`TableRunner` runs a series of test cases on a single function.

For unadic functions (1 parameter, 1 return value), its usage is straightforward:

```go
func isEven(n int) { return n&1 == 0 }

func TestIsEven(t *testing.T) {
    testx.Table(isEven).Cases([]testx.Case{
        {In: 0, Exp: true},
        {In: 1, Exp: false},
        {In: -1, Exp: false},
        {In: -2, Exp: true},
    }).Run(t)
}
```

Note that `TableRunner` supports any function type (any parameters number,
any return values numbers). If the tested function is non-unadic, it requires
an additional configuration to know where to inject `Case.In` and which
return value to compare `Case.Exp` with.

See file [example_table_test.go](./example_table_test.go) for detailed examples.

## Running tests

All runners expose two methods to run the tests: `Run` and `DryRun`.

### Method `Run`

`Run(t *testing.T)` runs the tests, fails `t` if any check fails,
and outputs the results like standard tests:

```
--- FAIL: TestMyHandler (0.00s)
  /my-repo/myhandler_test.go:64: response code:
      exp 401
      got 200
FAIL
FAIL	my-repo	0.247s
FAIL
```

### Method `DryRun`

`DryRun()` runs the tests, store the results and returns a `Resulter` interface
to access the stored results:

```go
// Resulter provides methods to read test results after a dry run.
type Resulter interface {
    // Checks returns a slice of CheckResults listing the runned checks
    Checks() []CheckResult
    // Passed returns true if all checks passed.
    Passed() bool
    // Failed returns true if one check or more failed.
    Failed() bool
    // NChecks returns the number of checks.
    NChecks() int
    // NPassed returns the number of checks that passed.
    NPassed() int
    // NFailed returns the number of checks that failed.
    NFailed() int
}
```

See [ExampleValueRunner_dryRun](./example_value_test.go) example.


## Further documentation

- [Go package documentation](https://pkg.go.dev/github.com/drykit-go/testx#section-documentation)

- Package `check` 📄 [Readme](./check/README.md)
  > Package `check` provides extensible and customizable checkers to perform checks on typed values.

- Package `check` 📄 [Readme](./checkconv/README.md)
  > Package `checkconv` provides conversion utilities to convert any typed checker to a `check.ValueChecker`
  
