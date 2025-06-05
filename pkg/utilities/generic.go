package utilities

// OptionalWithFallback if optional is nil or the zero-value for T, fallback is returned.
func OptionalWithFallback[T comparable](optional *T, fallback T) T {
	if optional == nil {
		return fallback
	}
	var defaultValue T
	if defaultValue == *optional {
		return fallback
	}
	return *optional
}

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
	var defaultValue T
	if val == defaultValue {
		return nil
	}
	return &val
}

// SafeConvert cast val to the T.
// If val is not T, the default value for T is returned.
func SafeConvert[T any](val any) T {
	var defaultValue T
	if v, ok := val.(T); ok {
		return v
	}
	return defaultValue
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
