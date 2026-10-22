// SPDX-License-Identifier: Apache-2.0
//
// Copyright © 2022 The Cardano Community Authors

package koios_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"testing"

	koios "github.com/cardano-community/koios-go-client/v5"
)

// mockCapture records the last request received by a mock Koios server.
type mockCapture struct {
	Method string
	Path   string
	Query  url.Values
	Body   []byte
}

// mockClient returns a client talking to a mock server that answers with the
// fixtures file of the given name.
func mockClient(t *testing.T, fixture string) (*koios.Client, *mockCapture) {
	t.Helper()
	b, err := os.ReadFile(filepath.Join("fixtures", fixture+".json"))
	if err != nil {
		t.Fatalf("read fixture %s: %v", fixture, err)
	}
	return mockClientRaw(t, string(b))
}

// mockClientRaw returns a client talking to a mock server that answers with body.
func mockClientRaw(t *testing.T, body string) (*koios.Client, *mockCapture) {
	t.Helper()
	capture := &mockCapture{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capture.Method = r.Method
		capture.Path = r.URL.Path
		capture.Query = r.URL.Query()
		capture.Body, _ = io.ReadAll(r.Body)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(body))
	}))
	t.Cleanup(srv.Close)

	u, err := url.Parse(srv.URL)
	if err != nil {
		t.Fatal(err)
	}
	port, err := strconv.ParseUint(u.Port(), 10, 16)
	if err != nil {
		t.Fatal(err)
	}
	client, err := koios.New(
		koios.Host(u.Hostname()),
		koios.Port(uint16(port)),
		koios.Scheme("http"),
	)
	if err != nil {
		t.Fatal(err)
	}
	return client, capture
}

// mockExpect asserts the method, path, query and JSON body of the captured
// request. An empty body asserts that no body was sent. The query must match
// exactly.
func mockExpect(t *testing.T, c *mockCapture, method, path string, query map[string]string, body string) {
	t.Helper()
	mockEq(t, "method", c.Method, method)
	mockEq(t, "path", c.Path, "/api/v1"+path)
	got := map[string]string{}
	for k := range c.Query {
		got[k] = c.Query.Get(k)
	}
	if query == nil {
		query = map[string]string{}
	}
	if !reflect.DeepEqual(got, query) {
		t.Errorf("query = %v, want %v", got, query)
	}
	if body == "" {
		if len(c.Body) != 0 {
			t.Errorf("body = %s, want none", c.Body)
		}
		return
	}
	var gotBody, wantBody any
	if err := json.Unmarshal(c.Body, &gotBody); err != nil {
		t.Fatalf("request body is not JSON: %v: %s", err, c.Body)
	}
	if err := json.Unmarshal([]byte(body), &wantBody); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(gotBody, wantBody) {
		t.Errorf("body = %s, want %s", c.Body, body)
	}
}

func mockEq[T comparable](t *testing.T, name string, got, want T) {
	t.Helper()
	if got != want {
		t.Errorf("%s = %v, want %v", name, got, want)
	}
}

func mockNoErr(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func mockLen[T any](t *testing.T, name string, s []T, want int) {
	t.Helper()
	if len(s) != want {
		t.Fatalf("len(%s) = %d, want %d", name, len(s), want)
	}
}

func mockJSONEq(t *testing.T, name string, got json.RawMessage, want string) {
	t.Helper()
	var g, w any
	if err := json.Unmarshal(got, &g); err != nil {
		t.Fatalf("%s is not JSON: %v: %s", name, err, got)
	}
	if err := json.Unmarshal([]byte(want), &w); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(g, w) {
		t.Errorf("%s = %s, want %s", name, got, want)
	}
}

func mockEpochPtr(n koios.EpochNo) *koios.EpochNo { return &n }
