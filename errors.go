package unstruct

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/aereal/unstruct/internal"
)

var (
	ErrInvalidTarget = errors.New("target must be a pointer to the struct")
	ErrValueNotFound = errors.New("no value found")
)

// UnsupportedTypeError is an error type represents target value type is not supported by Source.
type UnsupportedTypeError = internal.UnsupportedTypeError

// DecodeFieldError is an error type represents the decode failure on the field.
type DecodeFieldError struct {
	// Field is a struct field that error occurs.
	Field reflect.StructField

	// Err is an error that occurs.
	Err error
}

var _ error = &DecodeFieldError{}

func (e *DecodeFieldError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field.Name, e.Err.Error())
}

func (e *DecodeFieldError) Unwrap() error {
	return e.Err
}

func (e *DecodeFieldError) Is(other error) bool {
	if e == nil {
		return false
	}
	var oe *DecodeFieldError
	if !errors.As(other, &oe) {
		return false
	}
	return e.Field.Name == oe.Field.Name && e.Field.PkgPath == oe.Field.PkgPath && errors.Is(e.Err, oe.Err)
}
