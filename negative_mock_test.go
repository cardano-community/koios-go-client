// SPDX-License-Identifier: Apache-2.0
//
// Copyright © 2022 The Cardano Community Authors

package koios_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"

	koios "github.com/cardano-community/koios-go-client/v5"
)

func mockStatusClient(t *testing.T, status int, body string) *koios.Client {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
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
	client, err := koios.New(koios.Host(u.Hostname()), koios.Port(uint16(port)), koios.Scheme("http"))
	if err != nil {
		t.Fatal(err)
	}
	return client
}

func TestMockHTTPErrorKeepsKoiosError(t *testing.T) {
	client := mockStatusClient(t, http.StatusBadRequest, `{"code":"PGRST100","message":"bad query"}`)
	res, err := client.GetDRepList(context.Background(), nil)
	if !errors.Is(err, koios.ErrResponse) {
		t.Fatalf("want ErrResponse, got %v", err)
	}
	if res.Error == nil || !strings.Contains(res.Error.Message, "bad query") {
		t.Fatalf("Koios error body not surfaced: %+v", res.Error)
	}
	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("status code = %d", res.StatusCode)
	}
}

func TestMockMalformedJSON(t *testing.T) {
	client, _ := mockClientRaw(t, `{"not":"a list"}`)
	if _, err := client.GetDRepList(context.Background(), nil); err == nil {
		t.Fatal("expected a decode error")
	}
}

func TestMockRequiredParams(t *testing.T) {
	client, capture := mockClient(t, "asset_list")
	if _, err := client.GetAssets(context.Background(), "", "", nil); !errors.Is(err, koios.ErrNoPolicyID) {
		t.Fatalf("GetAssets: want ErrNoPolicyID, got %v", err)
	}
	if _, err := client.GetPoolInvalidDelegators(context.Background(), "", nil, nil); !errors.Is(err, koios.ErrNoPoolID) {
		t.Fatalf("GetPoolInvalidDelegators: want ErrNoPoolID, got %v", err)
	}
	if capture.Method != "" {
		t.Fatalf("no request should be sent, got %s %s", capture.Method, capture.Path)
	}
}

func TestMockCommitteeInfoEmpty(t *testing.T) {
	client, _ := mockClientRaw(t, `[]`)
	if _, err := client.GetCommitteeInfo(context.Background(), nil); !errors.Is(err, koios.ErrNoData) {
		t.Fatalf("want ErrNoData, got %v", err)
	}
}

func TestMockCollateralOutputAcceptsArray(t *testing.T) {
	var out koios.TxCollateralOutput
	in := `[{"tx_hash":"aa","tx_index":3,"value":"5","asset_list":[]}]`
	if err := json.Unmarshal([]byte(in), &out); err != nil {
		t.Fatal(err)
	}
	mockEq(t, "tx_index", out.TxIndex, 3)
	mockEq(t, "tx_hash", out.TxHash, koios.TxHash("aa"))
}

// A request that fails before it is sent must not leave goroutines behind.
func TestMockCancelledRequestsDoNotLeak(t *testing.T) {
	client, _ := mockClientRaw(t, `[]`)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	before := runtime.NumGoroutine()
	for range 100 {
		_, _ = client.GetBlocksTxs(ctx, []koios.BlockHash{"aa"}, nil)
		_, _ = client.GetPoolOwnerHistory(ctx, []koios.PoolID{"pool1"}, nil)
	}
	time.Sleep(100 * time.Millisecond)
	if grown := runtime.NumGoroutine() - before; grown > 20 {
		t.Fatalf("goroutines grew by %d", grown)
	}
}
