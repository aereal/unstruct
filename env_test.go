package unstruct_test

import (
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/aereal/unstruct"
	"github.com/google/go-cmp/cmp"
)

type data struct {
	FieldStr      string
	FieldInt      int
	FieldFloat    float64
	FieldBool     bool
	FieldStrSlice []string
}

func TestEnvironmentSource_FillValue(t *testing.T) {
	type testCase struct {
		name    string
		env     map[string]string
		opts    []unstruct.EnvironmentSourceOption
		want    any
		wantErr error
	}
	testCases := []testCase{
		{
			name: "ok",
			want: data{FieldStr: "str", FieldInt: 12, FieldFloat: 3.14, FieldBool: true, FieldStrSlice: []string{"a", "b", "c"}},
			env:  map[string]string{"FIELD_STR": "str", "FIELD_INT": "12", "FIELD_FLOAT": "3.14", "FIELD_BOOL": "true", "FIELD_STR_SLICE": "a,b,c"},
		},
		{
			name: "ok/with prefix",
			opts: []unstruct.EnvironmentSourceOption{unstruct.WithPrefix("X_")},
			want: data{FieldStr: "str", FieldInt: 12, FieldFloat: 3.14, FieldBool: true, FieldStrSlice: []string{"a", "b", "c"}},
			env:  map[string]string{"X_FIELD_STR": "str", "X_FIELD_INT": "12", "X_FIELD_FLOAT": "3.14", "X_FIELD_BOOL": "true", "X_FIELD_STR_SLICE": "a,b,c"},
		},
		{
			name: "ok/with custom delimiter",
			want: data{FieldStr: "str", FieldInt: 12, FieldFloat: 3.14, FieldBool: true, FieldStrSlice: []string{"a", "b", "c"}},
			opts: []unstruct.EnvironmentSourceOption{unstruct.WithSliceDelimiter(";")},
			env:  map[string]string{"FIELD_STR": "str", "FIELD_INT": "12", "FIELD_FLOAT": "3.14", "FIELD_BOOL": "true", "FIELD_STR_SLICE": "a;b;c"},
		},
		{
			name:    "ng/lacked required field",
			env:     map[string]string{},
			wantErr: &unstruct.DecodeFieldError{Err: unstruct.ErrValueNotFound, Field: reflect.StructField{Name: "FieldInt"}},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			for k, v := range tc.env {
				t.Setenv(k, v)
			}
			codec := unstruct.NewDecoder(unstruct.NewEnvironmentSource(tc.opts...))
			var got data
			err := codec.Decode(&got)
			if (err != nil) != (tc.wantErr != nil) {
				t.Fatalf("error: want=%v got=%v", tc.wantErr, err)
			}
			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("error: want=%v got=%v", tc.wantErr, err)
			}
			if err != nil {
				return
			}
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("decoded value (-want, +got):\n%s", diff)
			}
		})
	}
}

type parent struct {
	Data data
}

func TestEnvironmentSource_FillValue_nested(t *testing.T) {
	env := map[string]string{"DATA_FIELD_STR": "str", "DATA_FIELD_INT": "12", "DATA_FIELD_FLOAT": "3.14", "DATA_FIELD_BOOL": "true", "DATA_FIELD_STR_SLICE": "a,b,c"}
	want := parent{Data: data{FieldStr: "str", FieldInt: 12, FieldFloat: 3.14, FieldBool: true, FieldStrSlice: []string{"a", "b", "c"}}}
	for k, v := range env {
		t.Setenv(k, v)
	}
	codec := unstruct.NewDecoder(unstruct.NewEnvironmentSource())
	var got parent
	err := codec.Decode(&got)
	if err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("decoded value (-want, +got):\n%s", diff)
	}
}

func TestEnvironmentSource_FillValue_textmarshaler(t *testing.T) {
	type marshalerData struct {
		Timestamp time.Time
	}
	want := marshalerData{
		Timestamp: time.Unix(1672760375, 0),
	}
	t.Setenv("TIMESTAMP", "2023-01-04T00:39:35+09:00")
	codec := unstruct.NewDecoder(unstruct.NewEnvironmentSource())
	var got marshalerData
	if err := codec.Decode(&got); err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("-want, +got:\n%s", diff)
	}
}
