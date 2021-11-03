<h1 align="center">testx/check</h1>

<table align="center">
	<tr>
		<td><a href="../README.md">testx</a></td>
		<td>testx/check</td>
	</tr>
</table>

Package `check` provides _checkers_ interfaces and implementations
for `testx` test runners.

## Table of contents

- [Checker interfaces](#checker-interfaces)
- [Provided checkers](#provided-checkers)
- [Custom checkers](#custom-checkers)
  - [Extend provided checkers](#extend-provided-checkers)
    - [Using a `New` function](#using-a-new-function)
    - [Using `check.Value.Custom`](#using-checkvaluecustom)
  - [Custom checkers of any type](#custom-checkers-of-any-type)
  - [Contribute!](#contribute)

## Checker interfaces

A _checker_ is a generic interface composed of 2 methods `Pass` and `Explain`:

```go
// Checker provides methods to run checks on any typed value.
type Checker[T any] interface {
    Pass(got T) bool
    Explain(label string, got interface{}) string
}
```

## Provided checkers

Package `check` provides a collection of basic checkers for common types and kinds:

<details>
  <summary>Covered types</summary>
  <table>
    <thead>
      <tr>
        <th>Go type</th>
        <th>Checker provider</th>
        <th>Interface</th>
      </tr>
    </thead>
    <tbody>
      <tr>
        <td><code>int</code>, <code>intN</code> (8, 16, 32, 64)</td>
        <td><code>check.Int</code>, <code>check.IntN</code></td>
        <td>
          <a href="https://pkg.go.dev/github.com/drykit-go/testx/check#NumberCheckerProvider">
            <code>NumberCheckerProvider[int]</code>, <code>NumberCheckerProvider[intN]</code>
          </a>
        </td>
      </tr>
      <tr>
        <td><code>uint</code>, <code>uintN</code> (8, 16, 32, 64)</td>
        <td><code>check.Uint</code>, <code>check.UintN</code></td>
        <td>
          <a href="https://pkg.go.dev/github.com/drykit-go/testx/check#NumberCheckerProvider">
            <code>NumberCheckerProvider[uint]</code>, <code>NumberCheckerProvider[uintN]</code>
          </a>
        </td>
      </tr>
      <tr>
        <td><code>floatN</code> (32, 64)</td>
        <td><code>check.FloatN</code></td>
        <td>
          <a href="https://pkg.go.dev/github.com/drykit-go/testx/check#NumberCheckerProvider">
            <code>NumberCheckerProvider[floatN]</code>
          </a>
        </td>
      </tr>
      <tr>
        <td><code>bool</code></td>
        <td><code>check.Bool</code></td>
        <td>
          <a href="https://pkg.go.dev/github.com/drykit-go/testx/check#BoolCheckerProvider">
            <code>BoolCheckerProvider</code>
          </a>
        </td>
      </tr>
      <tr>
        <td><code>string</code></td>
        <td><code>check.String</code></td>
        <td>
          <a href="https://pkg.go.dev/github.com/drykit-go/testx/check#StringCheckerProvider">
            <code>StringCheckerProvider</code>
          </a>
        </td>
      </tr>
      <tr>
        <td><code>[]byte</code></td>
        <td><code>check.Bytes</code></td>
        <td>
          <a href="https://pkg.go.dev/github.com/drykit-go/testx/check#BytesCheckerProvider">
            <code>BytesCheckerProvider</code>
          </a>
        </td>
      </tr>
      <tr>
        <td><code>time.Duration</code></td>
        <td><code>check.Duration</code></td>
        <td>
          <a href="https://pkg.go.dev/github.com/drykit-go/testx/check#DurationCheckerProvider">
            <code>DurationCheckerProvider</code>
          </a>
        </td>
      </tr>
      <tr>
        <td><code>context.Context</code></td>
        <td><code>check.Context</code></td>
        <td>
          <a href="https://pkg.go.dev/github.com/drykit-go/testx/check#ContextCheckerProvider">
            <code>ContextCheckerProvider</code>
          </a>
        </td>
      </tr>
      <tr>
        <td><code>http.Header</code></td>
        <td><code>check.HTTPHeader</code></td>
        <td>
          <a href="https://pkg.go.dev/github.com/drykit-go/testx/check#HTTPHeaderCheckerProvider">
            <code>HTTPHeaderCheckerProvider</code>
          </a>
        </td>
      </tr>
      <tr>
        <td><code>*http.Request</code></td>
        <td><code>check.HTTPRequest</code></td>
        <td>
          <a href="https://pkg.go.dev/github.com/drykit-go/testx/check#HTTPRequestCheckerProvider">
            <code>HTTPRequestCheckerProvider</code>
          </a>
        </td>
      </tr>
      <tr>
        <td><code>*http.Response</code></td>
        <td><code>check.HTTPResponse</code></td>
        <td>
          <a href="https://pkg.go.dev/github.com/drykit-go/testx/check#HTTPResponseCheckerProvider">
            <code>HTTPResponseCheckerProvider</code>
          </a>
        </td>
      </tr>
      <tr>
        <td><code>any</code></td>
        <td><code>check.Value</code></td>
        <td>
          <a href="https://pkg.go.dev/github.com/drykit-go/testx/check#ValueCheckerProvider">
            <code>ValueCheckerProvider</code>
          </a>
        </td>
      </tr>
    </tbody>
    <thead>
      <tr>
        <th>Go kind</th>
        <th>Checker provider</th>
      </tr>
    </thead>
    <tbody>
      <tr>
        <td><code>slice</code></td>
        <td><code>check.Slice</code></td>
        <td>
          <a href="https://pkg.go.dev/github.com/drykit-go/testx/check#SliceCheckerProvider">
            <code>SliceCheckerProvider</code>
          </a>
        </td>
      </tr>
      <tr>
        <td><code>map</code></td>
        <td><code>check.Map</code></td>
        <td>
          <a href="https://pkg.go.dev/github.com/drykit-go/testx/check#MapCheckerProvider">
            <code>MapCheckerProvider</code>
          </a>
        </td>
      </tr>
      <tr>
        <td><code>struct</code></td>
        <td><code>check.Struct</code></td>
        <td>
          <a href="https://pkg.go.dev/github.com/drykit-go/testx/check#StructCheckerProvider">
            <code>StructCheckerProvider</code>
          </a>
        </td>
      </tr>
    </tbody>
  </table>
</details>

## Custom checkers

This section is about creating custom checkers, either extending the provided ones or full-fledged customs.

### Extend provided checkers

There are 2 ways to extend a provided checker.

#### Using a `New` function

Package `check` provides a `NewChecker` function to create a custom checker.

It takes a `PassFunc` that determinates if a got value passes the check
and an `ExplainFunc` that outputs the reason of a failed check,
then returns a checker for that type.

```go
func Test42IsEven(t *testing.T) {
    checkIsEven := check.NewChecker(
        func(got int) bool { return got&1 == 0 },
        func(label string, got any) string {
            return fmt.Sprintf("%s: expect even int, got %v", label, got)
        },
      )

    testx.Value(42).Pass(checkIsEven).Run(t)
}
```

Related examples:

- [NewChecker](https://pkg.go.dev/github.com/drykit-go/testx/check#example-package-NewChecker)

#### Using `check.Value.Custom`

```go
func Test42IsEven(t *testing.T) {
    checkIsEven := check.Value.Custom(
        "even int",
        func(got any) bool {
          gotInt, ok := got.(int)
          return ok && got&1 == 0
        },
      )

    testx.Value(42).Pass(checkIsEven).Run(t)
}
```

### Custom checker implementation

Any type that implements `check.Checker[T any]` interface is a valid checker
and can be used for `testx` runners.

For runners that require a `check.Checker[interface{}]` specifically
(such as `testx.Case.Pass`), use `check.Wrap` to perform the conversion.

Related examples:

- [CustomChecker](https://pkg.go.dev/github.com/drykit-go/testx/check#example-package-CustomChecker):
  implementation of a custom checker.
- [CustomCheckerClosures](https://pkg.go.dev/github.com/drykit-go/testx/check#example-package-CustomCheckerClosures):
  advanced implementation of a parameterized custom checker
  for type `complex128`, using closures.

### Contribute!

If you believe a particuliar checker should be natively implemented,
feel free to contribute!

- [📋 File an issue](https://github.com/drykit-go/testx/issues/new)
- [🔀 Open a Pull request](https://github.com/drykit-go/testx/pulls)
- [💬 Open a discussion](https://github.com/drykit-go/testx/discussions/new)
