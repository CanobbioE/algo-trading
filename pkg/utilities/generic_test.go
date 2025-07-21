package utilities_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/CanobbioE/algo-trading/pkg/utilities"
)

func TestDefaultPointer(t *testing.T) {
	t.Run("returns default value when input is nil", func(t *testing.T) {
		want := ""
		if got := utilities.DefaultPointer[string](nil); !reflect.DeepEqual(got, &want) {
			t.Errorf("DefaultPointer() = %+v, want \"\"", got)
		}
	})

	t.Run("returns a pointer to the input when non-nil", func(t *testing.T) {
		want := "test"
		if got := utilities.DefaultPointer(&want); !reflect.DeepEqual(got, &want) {
			t.Errorf("DefaultPointer() = %v, want %v", got, want)
		}
	})
}

func TestToPointer(t *testing.T) {
	t.Run("returns a pointer to a value", func(t *testing.T) {
		want := "test"
		if got := utilities.ToPointer("test"); !reflect.DeepEqual(got, &want) {
			t.Errorf("DefaultPointer(STRING) = %v, want %v", got, want)
		}
	})
}

func TestMustReturn(t *testing.T) {
	mockCall := func(shouldErr bool) (string, error) {
		if shouldErr {
			return "", errors.New("error")
		}
		return "test", nil
	}
	t.Run("returns a value when no error occurs", func(t *testing.T) {
		if got := utilities.MustReturn(mockCall(false)); !reflect.DeepEqual(got, "test") {
			t.Errorf("DefaultPointer(STRING) = %v, want \"test\"", got)
		}
	})
	t.Run("panics when an error occurs", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Log("Recovered from expected panic", r)
			}
		}()
		utilities.MustReturn(mockCall(true))
		t.Errorf("MustReturn should panic when an error occurs")
	})
}

func TestMust(t *testing.T) {
	mockCall := func(shouldErr bool) error {
		if shouldErr {
			return errors.New("error")
		}
		return nil
	}
	t.Run("succeeds when no error occurs", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Fatal("Unexpected panic", r)
			}
		}()
		utilities.Must(mockCall(false))
	})
	t.Run("panics when an error occurs", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Log("Recovered from expected panic", r)
			}
		}()
		utilities.Must(mockCall(true))
		t.Errorf("Must should panic when an error occurs")
	})
}
