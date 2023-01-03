package unstruct

import (
	"os"
	"reflect"
	"strings"
	"unicode"

	"github.com/aereal/unstruct/internal"
)

type environmentSourceConfig struct {
	sliceDelimiter string
	prefix         string
}

type EnvironmentSourceOption func(*environmentSourceConfig)

func WithSliceDelimiter(delimiter string) EnvironmentSourceOption {
	return func(cfg *environmentSourceConfig) {
		cfg.sliceDelimiter = delimiter
	}
}

// WithPrefix indicates the EnvironmentSource fetches environment variables with given prefix.
func WithPrefix(prefix string) EnvironmentSourceOption {
	return func(cfg *environmentSourceConfig) {
		cfg.prefix = prefix
	}
}

func NewEnvironmentSource(opts ...EnvironmentSourceOption) *EnvironmentSource {
	src := &EnvironmentSource{}
	for _, o := range opts {
		o(&src.environmentSourceConfig)
	}
	return src
}

// EnvironmentSource is a Source that reads values from environment variables.
type EnvironmentSource struct {
	environmentSourceConfig
}

var _ Source = &EnvironmentSource{}

const defaultSliceDelimiter = ","

func (s *EnvironmentSource) FillValue(path Path, target reflect.Value) error {
	var parts []string
	for _, sf := range path {
		parts = append(parts, strings.ToUpper(toSnakeCase(sf.Name)))
	}
	name := strings.Join(parts, "_")
	if s.prefix != "" {
		name = s.prefix + name
	}
	val := os.Getenv(name)
	switch kind := target.Type().Kind(); kind {
	case reflect.Slice:
		delimiter := s.sliceDelimiter
		if delimiter == "" {
			delimiter = defaultSliceDelimiter
		}
		parts := strings.Split(val, delimiter)

		// grow slice size
		if size := len(parts); size > target.Cap() {
			grown := reflect.MakeSlice(target.Type(), target.Len(), size)
			reflect.Copy(grown, target)
			target.Set(grown)
			target.SetLen(size)
		}

		for i, el := range parts {
			if err := internal.DecodeStringToScalarType(el, target.Index(i)); err != nil {
				return err
			}
		}
	default:
		if err := internal.DecodeStringToScalarType(val, target); err != nil {
			return err
		}
	}
	return nil
}

func toSnakeCase(s string) string {
	b := &strings.Builder{}
	for i, r := range s {
		switch {
		case i == 0:
			b.WriteRune(unicode.ToLower(r))
		case unicode.IsUpper(r):
			b.WriteRune('_')
			b.WriteRune(unicode.ToLower(r))
		default:
			b.WriteRune(r)
		}
	}
	return b.String()
}
