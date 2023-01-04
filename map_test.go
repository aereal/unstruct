package unstruct_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/aereal/unstruct"
	"github.com/google/go-cmp/cmp"
)

func TestMapSource_FillValue(t *testing.T) {
	type testCase struct {
		name    string
		mapVal  map[string]any
		want    any
		wantErr error
	}
	testCases := []testCase{
		{
			name:   "ok",
			want:   data{FieldStr: "str", FieldInt: 12, FieldFloat: 3.14, FieldBool: true, FieldStrSlice: []string{"a", "b", "c"}},
			mapVal: map[string]any{"FieldStr": "str", "FieldInt": 12, "FieldFloat": 3.14, "FieldBool": true, "FieldStrSlice": []string{"a", "b", "c"}},
		},
		{
			name:    "ng/lacked required field",
			mapVal:  map[string]any{"FieldStr": "str", "FieldFloat": 3.14, "FieldBool": true, "FieldStrSlice": []string{"a", "b", "c"}},
			wantErr: &unstruct.DecodeFieldError{Err: unstruct.ErrValueNotFound, Field: reflect.StructField{Name: "FieldInt"}},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			codec := unstruct.NewDecoder(unstruct.NewMapSource(tc.mapVal))
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
