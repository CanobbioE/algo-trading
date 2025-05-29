package utilities

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

func ToPointer[T comparable](val T) *T {
	var defaultValue T
	if val == defaultValue {
		return nil
	}
	return &val
}

func SafeConvert[T any](val any) T {
	var defaultValue T
	if v, ok := val.(T); ok {
		return v
	}
	return defaultValue
}

func MustReturn[T any](val T, err error) T {
	Must(err)
	return val
}

func Must(err error) {
	if err != nil {
		panic(err)
	}
}
