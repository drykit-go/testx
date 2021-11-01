package check

// Wrap returns a Checker[any] that wraps the given checker.
func Wrap[T any](checker Checker[T]) Checker[any] {
	return NewChecker(
		func(got interface{}) bool { return checker.Pass(got.(T)) },
		checker.Explain,
	)
}

// Wrap returns a slice of Checker[any] that wrap the given checkers.
func WrapMany[T any](checkers ...Checker[T]) []Checker[any] {
	anyCheckers := make([]Checker[any], len(checkers))
	for i, c := range checkers {
		anyCheckers[i] = Wrap(c)
	}
	return anyCheckers
}
