package utilities

// DefaultPointer if optional is nil, a pointer to the default value for T is returned.
func DefaultPointer[T comparable](optional *T) *T {
	var defaultValue T
	if optional == nil {
		return &defaultValue
	}
	return optional
}

// ToPointer returns a pointer to val.
func ToPointer[T comparable](val T) *T {
	return &val
}

// MustReturn panics if an error is passed, otherwise it returns the expected value.
func MustReturn[T any](val T, err error) T {
	Must(err)
	return val
}

// Must panics if an error is passed.
func Must(err error) {
	if err != nil {
		panic(err)
	}
}
