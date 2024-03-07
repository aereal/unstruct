package unstruct

import "reflect"

func NewCascadeSource(srcs ...Source) *CascadeSource {
	cs := &CascadeSource{srcs: srcs}
	return cs
}

type CascadeSource struct {
	srcs []Source
}

var _ Source = &CascadeSource{}

func (s *CascadeSource) FillValue(path Path, target reflect.Value) error {
	for _, src := range s.srcs {
		if err := src.FillValue(path, target); err == nil {
			return nil
		}
	}
	return ErrValueNotFound
}
