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
- [Package `check`](#package-check)
- [Package `checkconv`](#package-checkconv)

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
        Exp(42).
        Not(3, "hello").
        Pass(checkconv.FromInt(check.Int.InRange(41, 43))).
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
    testx.Table(isEven, nil).Cases([]testx.Case{
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
    // Check returns a slice of CheckResults listing the runned checks
    Checks() []CheckResult
    // Passed returns true if all checks passed.
    Passed() bool
    // Failed returns true if one check or more failed.
    Failed() bool
    // NChecks returns the number of checks.
    NChecks() int
    // NPassed returns the number of checks that passed.
    NPassed() int
    // NPassed returns the number of checks that failed.
    NFailed() int
}
```

See [ExampleValueRunner_dryRun](./example_value_test.go) example.

## Package `check`

As you may have noticed from the previous examples, `testx` runners heavily
rely on _checkers_ defined in package `check`.

> Package `check` provides extensible and customizable checkers to perform checks on typed values.

📚 [Package `check` README](./check/README.md)

## Package `checkconv`

> Package `checkconv` provides conversion utilities to convert any typed checker to a `check.ValueChecker`

📚 [Package `checkconv` README](./checkconv/README.md)
