package unstruct

import (
	"reflect"
	"strings"
)

func NewMapSource(m map[string]any) *MapSource {
	src := &MapSource{m: map[string]any{}}
	if m != nil {
		src.m = m
	}
	return src
}

type MapSource struct {
	m map[string]any
}

var _ Source = &MapSource{}

func (s *MapSource) FillValue(path Path, target reflect.Value) error {
	paths := make([]string, 0, len(path))
	for _, sf := range path {
		paths = append(paths, sf.Name)
	}
	name := strings.Join(paths, "/")
	v, ok := s.m[name]
	if !ok {
		return ErrValueNotFound
	}
	vt := reflect.TypeOf(v)
	switch kind := target.Type().Kind(); kind {
	case reflect.Slice:
		if !vt.AssignableTo(target.Type()) {
			return &UnsupportedTypeError{Kind: target.Kind().String()}
		}
		rv := reflect.ValueOf(v)
		if size := rv.Len(); size > target.Cap() {
			grown := reflect.MakeSlice(target.Type(), target.Len(), size)
			_ = reflect.Copy(grown, target)
			target.Set(grown)
			target.SetLen(size)
		}
		for i := 0; i < rv.Len(); i++ {
			target.Index(i).Set(rv.Index(i))
		}
	default:
		if !vt.AssignableTo(target.Type()) {
			return &UnsupportedTypeError{vt.Kind().String()}
		}
		target.Set(reflect.ValueOf(v))
	}
	return nil
}
