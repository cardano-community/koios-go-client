// SPDX-License-Identifier: Apache-2.0
//
// Copyright © 2022 The Cardano Community Authors

package koios_test

import (
	"context"
	"errors"
	"testing"

	koios "github.com/cardano-community/koios-go-client/v5"
)

func TestMockOgmios(t *testing.T) {
	client, capture := mockClient(t, "ogmios")
	res, err := client.Ogmios(context.Background(), koios.OgmiosRequest{Method: koios.OgmiosQueryNetworkTip}, nil)
	mockNoErr(t, err)
	mockExpect(t, capture, "POST", "/ogmios", nil, `{"jsonrpc": "2.0", "method": "queryNetwork/tip"}`)
	mockJSONEq(t, "data", res.Data, `{"jsonrpc": "2.0", "method": "queryNetwork/tip",
		"result": {"slot": 59886800, "id": "df5678c9774b7bc8c60a4c83b63c3676e618640ad05f7d1ee775b68939cf77d1"}}`)
	mockEq(t, "status_code", res.StatusCode, 200)
}

func TestMockOgmiosParams(t *testing.T) {
	client, capture := mockClientRaw(t, `{"jsonrpc": "2.0", "method": "queryLedgerState/liveStakeDistribution", "result": {}}`)
	_, err := client.Ogmios(context.Background(), koios.OgmiosRequest{
		JSONRPC: "2.0",
		Method:  koios.OgmiosQueryLedgerStateLiveStakeDistribution,
		Params:  map[string]any{"pools": []string{"pool155efqn9xpcf73pphkk88cmlkdwx4ulkg606tne970qswczg3asc"}},
	}, nil)
	mockNoErr(t, err)
	mockExpect(t, capture, "POST", "/ogmios", nil, `{"jsonrpc": "2.0", "method": "queryLedgerState/liveStakeDistribution",
		"params": {"pools": ["pool155efqn9xpcf73pphkk88cmlkdwx4ulkg606tne970qswczg3asc"]}}`)
}

func TestMockOgmiosNoMethod(t *testing.T) {
	client, _ := mockClientRaw(t, `{}`)
	_, err := client.Ogmios(context.Background(), koios.OgmiosRequest{}, nil)
	if !errors.Is(err, koios.ErrNoOgmiosMethod) {
		t.Errorf("err = %v, want ErrNoOgmiosMethod", err)
	}
}
