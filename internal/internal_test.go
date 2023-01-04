package internal_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/aereal/unstruct/internal"
)

func TestExtractTextUnmarshaler(t *testing.T) {
	type testCase struct {
		name  string
		value any
		want  bool
	}
	now := time.Now()
	testCases := []testCase{
		{
			name:  "not implemented",
			value: "",
			want:  false,
		},
		{
			name:  "implemented type",
			value: &now,
			want:  true,
		},
		{
			name:  "implemented underlying type",
			value: now,
			want:  true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(tc.value)
			_, ok := internal.ExtractTextUnmarshaler(v)
			if ok != tc.want {
				t.Errorf("want=%v got=%v", tc.want, ok)
			}
		})
	}
}
