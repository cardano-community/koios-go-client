// Copyright 2022 The Cardano Community Authors
// SPDX-License-Identifier: Apache-2.0
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
//
//   http://www.apache.org/licenses/LICENSE-2.0
//   or LICENSE file in repository root.
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package koios_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/cardano-community/koios-go-client/v2"
	"github.com/stretchr/testify/assert"
)

func TestBlocks(t *testing.T) {
	client, err := getClient()
	if !assert.NoError(t, err) {
		return
	}
	blocksTest(t, client)
}

func blocksTest(t TestingT, client *koios.Client) {

	opts := client.NewRequestOptions()
	opts.SetPageSize(10)
	res, err := client.GetBlocks(context.Background(), opts)
	if !assert.NoError(t, err) {
		return
	}
	assertEqual(t, 10, len(res.Data), "total blocks returned")

	assertNotEmpty(t, res.Data[0].Hash, "hash")
	assertGreater(t, res.Data[0].EpochNo, koios.EpochNo(0), "epoch_no")
	assertGreater(t, res.Data[0].AbsSlot, koios.Slot(0), "abs_slot")
	assertGreater(t, res.Data[0].EpochSlot, koios.Slot(0), "epoch_slot")
	assertGreater(t, res.Data[0].Height, 0, "block_height")
	assertGreater(t, res.Data[0].Size, 0, "block_size")
	assertTimeNotZero(t, res.Data[0].Time, "block_time")
	if res.Data[0].TxCount == 0 {
		githubActionWarning("/blocks", fmt.Sprintf("block(%s) tx count is 0", res.Data[0].Hash))
	}
	assertNotEmpty(t, res.Data[0].VrfKey, "vrf_key")
	assertNotEmpty(t, res.Data[0].Pool, "pool")
	assertGreater(t, res.Data[0].OpCertCounter, 0, "op_cert_counter")
	assertGreater(t, res.Data[0].ProtoMajor, 0, "proto_major")
	// assertGreater(t, res.Data[0].ProtoMinor, 0, "proto_minor")
}

func TestBlockInfo(t *testing.T) {
	client, err := getClient()
	if !assert.NoError(t, err) {
		return
	}
	blockInfoTest(t, networkBlockHash(), client)
}

func blockInfoTest(t TestingT, hash koios.BlockHash, client *koios.Client) {
	res, err := client.GetBlockInfo(context.Background(), hash, nil)
	if !assert.NoError(t, err) {
		return
	}

	assertNotEmpty(t, res.Data.Hash, "hash")
	assertGreater(t, res.Data.EpochNo, koios.EpochNo(0), "epoch_no")
	assertGreater(t, res.Data.AbsSlot, koios.Slot(0), "abs_slot")
	assertGreater(t, res.Data.EpochSlot, koios.Slot(0), "epoch_slot")
	assertGreater(t, res.Data.Height, 0, "block_height")
	assertGreater(t, res.Data.Size, 0, "block_size")
	assertTimeNotZero(t, res.Data.Time, "block_time")

	if res.Data.TxCount == 0 {
		githubActionWarning("/block_info", fmt.Sprintf("block(%s) tx count is 0", res.Data.Hash))
	} else {
		assertIsPositive(t, res.Data.TotalOutput, "total_output")
		assertIsPositive(t, res.Data.TotalFees, "total_fees")
	}

	assertNotEmpty(t, res.Data.VrfKey, "vrf_key")
	assertNotEmpty(t, res.Data.OpCert, "op_cert")
	assertGreater(t, res.Data.OpCertCounter, 0, "op_cert_counter")
	assertNotEmpty(t, res.Data.Pool, "pool")
	assertGreater(t, res.Data.ProtoMajor, 0, "proto_major")
	// assertGreater(t, res.Data.ProtoMinor, 0, "proto_minor")

	assertGreater(t, res.Data.Confirmations, 0, "num_confirmations")
	assertNotEmpty(t, res.Data.ParentHash, "parent_hash")
	assertNotEmpty(t, res.Data.ChildHash, "child_hash")
}

func TestBlockTxs(t *testing.T) {
	client, err := getClient()
	if !assert.NoError(t, err) {
		return
	}
	blockTxsTest(t, networkBlockHash(), client)
}

func blockTxsTest(t TestingT, hash koios.BlockHash, client *koios.Client) {
	res, err := client.GetBlockTxs(context.Background(), hash, nil)
	if err != nil {
		if assert.ErrorIs(t, err, koios.ErrNoData) {
			githubActionWarning("BlockTxs", err.Error())
			return
		}
		assert.NoError(t, err)
		return
	}
	assertEqual(t, hash, res.Data.BlockHash, "req/res block hashes do not match")
	assertGreater(t, len(res.Data.TxHashes), 0, "tx_hashes")
}
