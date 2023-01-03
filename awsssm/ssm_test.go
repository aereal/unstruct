package awsssm_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/aereal/unstruct"
	"github.com/aereal/unstruct/awsssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"github.com/google/go-cmp/cmp"
)

var (
	wantStr         = "str"
	wantInt         = 1
	wantStringSlice = []string{"a", "b", "c"}
)

type namedTarget struct {
	StringField      string
	IntField         int
	StringSliceField []string
}

func newTestSSMClient(handler http.Handler) (awsssm.SSMClient, func()) {
	srv := httptest.NewServer(handler)
	options := ssm.Options{
		EndpointResolver: ssm.EndpointResolverFromURL(srv.URL),
	}
	return ssm.New(options), func() { srv.Close() }
}

func TestSource_FillValue(t *testing.T) {
	pathPrefix := "/my"
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		out := &ssm.GetParametersByPathOutput{}
		out.Parameters = []types.Parameter{
			{Name: ref(fmt.Sprintf("%s/StringField", pathPrefix)), Value: ref(wantStr), Type: types.ParameterTypeString},
			{Name: ref(fmt.Sprintf("%s/IntField", pathPrefix)), Value: ref(strconv.Itoa(wantInt)), Type: types.ParameterTypeString},
			{Name: ref(fmt.Sprintf("%s/StringSliceField", pathPrefix)), Value: ref(strings.Join(wantStringSlice, ",")), Type: types.ParameterTypeStringList},
		}
		_ = json.NewEncoder(w).Encode(out)
	})
	client, clean := newTestSSMClient(handler)
	t.Cleanup(clean)
	srv, err := awsssm.New(awsssm.WithClient(client), awsssm.WithPathPrefix(pathPrefix))
	if err != nil {
		t.Fatal(err)
	}
	dec := unstruct.NewDecoder(srv)
	var target namedTarget
	if err := dec.Decode(&target); err != nil {
		t.Fatal(err)
	}
	want := namedTarget{
		StringField:      wantStr,
		IntField:         wantInt,
		StringSliceField: wantStringSlice,
	}
	if diff := cmp.Diff(want, target); diff != "" {
		t.Errorf("-want, +got:\n%s", diff)
	}
}

func ref[T any](v T) *T {
	return &v
}
