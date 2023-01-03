package awsssm

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"sync"

	"github.com/aereal/unstruct"
	"github.com/aereal/unstruct/internal"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
)

var (
	ErrClientRequired     = errors.New("ssm client is required")
	ErrPathPrefixRequired = errors.New("path prefix is required")
)

// Option is a function configures a Source's behavior.
type Option func(*config)

// WithContext indicates the Source to use given context.
func WithContext(ctx context.Context) Option {
	return func(c *config) {
		c.ctx = ctx
	}
}

// WithClient indicates the Source to use given SSM client.
func WithClient(client SSMClient) Option {
	return func(c *config) {
		c.client = client
	}
}

// WithPathPrefix indicates the Source to get parameters with given path prefix.
func WithPathPrefix(pathPrefix string) Option {
	return func(c *config) {
		c.pathPrefix = pathPrefix
	}
}

type SSMClient interface {
	GetParametersByPath(ctx context.Context, params *ssm.GetParametersByPathInput, optFns ...func(*ssm.Options)) (*ssm.GetParametersByPathOutput, error)
}

type config struct {
	client     SSMClient
	pathPrefix string
	ctx        context.Context
}

func New(opts ...Option) (*Source, error) {
	src := &Source{}
	for _, o := range opts {
		o(&src.config)
	}
	if src.client == nil {
		return nil, ErrClientRequired
	}
	if src.pathPrefix == "" {
		return nil, ErrPathPrefixRequired
	}
	return src, nil
}

type Source struct {
	config

	init         sync.Once
	err          error
	paramsByName map[string]types.Parameter
}

var _ unstruct.Source = &Source{}

func (s *Source) doInit() {
	ctx := s.ctx
	if ctx == nil {
		ctx = context.Background()
	}
	input := &ssm.GetParametersByPathInput{
		Recursive: ref(true),
		Path:      ref(s.pathPrefix),
	}
	out, err := s.client.GetParametersByPath(ctx, input)
	if err != nil {
		s.err = err
		return
	}
	s.paramsByName = map[string]types.Parameter{}
	for _, p := range out.Parameters {
		if p.Name == nil {
			continue
		}
		s.paramsByName[*p.Name] = p
	}
}

const (
	pathDelimiter    = "/"
	strListDelimiter = ","
)

func (s *Source) FillValue(path unstruct.Path, target reflect.Value) error {
	s.init.Do(func() { s.doInit() })
	if s.err != nil {
		return s.err
	}
	hierarchy := make([]string, 0, len(path)+1)
	hierarchy = append(hierarchy, s.pathPrefix)
	for _, f := range path {
		hierarchy = append(hierarchy, f.Name)
	}
	name := strings.Join(hierarchy, pathDelimiter)
	param, found := s.paramsByName[name]
	if !found || param.Value == nil {
		return unstruct.ErrValueNotFound
	}
	val := *param.Value
	switch kind := target.Type().Kind(); kind {
	case reflect.Slice:
		parts := strings.Split(val, strListDelimiter)

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

func ref[T any](v T) *T {
	return &v
}
