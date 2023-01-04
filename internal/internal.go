package internal

import (
	"encoding"
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

var (
	ErrOverflowSize     = errors.New("overflow size")
	ErrInvalidBoolValue = errors.New("invalid bool value")

	textUnmarshaler = reflect.TypeOf((*encoding.TextUnmarshaler)(nil)).Elem()
)

func extract(v reflect.Value) (encoding.TextUnmarshaler, bool) {
	if v.Type().Implements(textUnmarshaler) {
		v2 := v
		if v2.IsNil() {
			v2.Set(reflect.New(v2.Type().Elem()))
		}
		if tu, ok := v2.Interface().(encoding.TextUnmarshaler); ok {
			return tu, ok
		}
	}
	return nil, false
}

func ExtractTextUnmarshaler(v reflect.Value) (encoding.TextUnmarshaler, bool) {
	if tu, ok := extract(v); ok {
		return tu, ok
	}
	if v.CanAddr() {
		if tu, ok := extract(v.Addr()); ok {
			return tu, ok
		}
	}
	ptr := reflect.New(v.Type())
	ptr.Elem().Set(v)
	if tu, ok := extract(ptr); ok {
		return tu, ok
	}
	return nil, false
}

func DecodeStringToScalarType(v string, value reflect.Value) error {
	if tu, ok := ExtractTextUnmarshaler(value); ok {
		return tu.UnmarshalText([]byte(v))
	}

	switch kind := value.Kind(); kind {
	case reflect.String:
		value.SetString(v)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		n, err := strconv.ParseInt(v, 10, value.Type().Bits())
		if err != nil {
			return err
		}
		if value.OverflowInt(n) {
			return ErrOverflowSize
		}
		value.SetInt(n)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		n, err := strconv.ParseUint(v, 10, value.Type().Bits())
		if err != nil {
			return err
		}
		if value.OverflowUint(n) {
			return ErrOverflowSize
		}
		value.SetUint(n)
	case reflect.Float32, reflect.Float64:
		n, err := strconv.ParseFloat(v, value.Type().Bits())
		if err != nil {
			return err
		}
		if value.OverflowFloat(n) {
			return ErrOverflowSize
		}
		value.SetFloat(n)
	case reflect.Bool:
		switch v {
		case "true":
			value.SetBool(true)
		case "false":
			value.SetBool(false)
		default:
			return ErrInvalidBoolValue
		}
	default:
		return &UnsupportedTypeError{Kind: kind.String()}
	}
	return nil
}

type UnsupportedTypeError struct {
	Kind string
}

var _ error = &UnsupportedTypeError{}

func (e *UnsupportedTypeError) Error() string {
	return fmt.Sprintf("Unsupported type: %s", e.Kind)
}
