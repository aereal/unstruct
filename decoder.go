package unstruct

import (
	"reflect"
	"strings"
)

// Source is an interface that fetches the values from external sources and fills into given target.
type Source interface {
	FillValue(path Path, target reflect.Value) error
}

func NewDecoder(srcs ...Source) *Decoder {
	return &Decoder{srcs: srcs}
}

// Decoder runs data sources against given target value and completes target value.
type Decoder struct {
	srcs []Source
}

// Decode decodes the values that comes from Codec's sources into given target.
func (d *Decoder) Decode(target any) error {
	v := reflect.ValueOf(target)
	if v.IsNil() || !(v.Kind() == reflect.Pointer && v.Elem().Kind() == reflect.Struct) {
		return ErrInvalidTarget
	}
	structType := v.Elem()
	typ := structType.Type()
	if err := d.decode(nil, structType, typ); err != nil {
		return err
	}
	return nil
}

func (d *Decoder) decode(path Path, structType reflect.Value, typ reflect.Type) error {
	numFields := typ.NumField()
	for i := 0; i < numFields; i++ {
		field := typ.Field(i)
		fieldValue := structType.Field(i)
		currPath := path.with(field)
		switch kind := field.Type.Kind(); kind {
		case reflect.Map:
			return &DecodeFieldError{Field: field, Err: &UnsupportedTypeError{Kind: kind.String()}}
		case reflect.Struct:
			if err := d.decode(currPath, fieldValue, fieldValue.Type()); err != nil {
				return &DecodeFieldError{Field: field, Err: err}
			}
		default:
			if err := d.fillValue(currPath, field, fieldValue); err != nil {
				return &DecodeFieldError{Field: field, Err: err}
			}
		}
	}
	return nil
}

func (d *Decoder) fillValue(path Path, field reflect.StructField, value reflect.Value) error {
	for _, s := range d.srcs {
		if err := s.FillValue(path, value); err == nil {
			return nil
		}
	}
	opts := parseTag(field.Tag)
	if opts.optional {
		return nil
	}
	return ErrValueNotFound
}

type Path []reflect.StructField

func (p Path) with(v reflect.StructField) Path {
	return append(p, v)
}

type fieldOption struct {
	optional bool
}

const (
	tagKey       = "unstruct"
	tagDelimiter = ","
	optOptional  = "optional"
)

func parseTag(tag reflect.StructTag) *fieldOption {
	opts := &fieldOption{}
	vals := strings.Split(tag.Get(tagKey), tagDelimiter)
	for _, v := range vals {
		switch v {
		case optOptional:
			opts.optional = true
		}
	}
	return opts
}
