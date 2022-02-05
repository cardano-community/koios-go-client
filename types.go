// Copyright 2022 The Howijd.Network Authors
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

package koios

import "github.com/shopspring/decimal"

type (
	// Address defines type for _address.
	Address string

	// AnyAddress defines type for _any_address.
	AnyAddress string

	// AssetName defines type for _asset_name.
	AssetName string

	// AssetPolicy defines type for _asset_policy.
	AssetPolicy string

	// BlockHash defines type for _block_hash.
	BlockHash string

	// EarnedEpochNo defines type for _earned_epoch_no.
	EarnedEpochNo string

	// EpochNo defines type for _epoch_no.
	EpochNo uint64

	// PoolBech32 defines type for _pool_bech32.
	PoolBech32 string

	// PoolBech32Optional defines type for _pool_bech32_optional.
	PoolBech32Optional string

	// ScriptHash defines type for _script_hash.
	ScriptHash string

	// StakeAddress defines type for _stake_address.
	StakeAddress string

	// Lovelace defines type for ADA lovelaces. This library uses forked snapshot
	// of github.com/shopspring/decimal package to provide. JSON and XML
	// serialization/deserialization and make it ease to work with calculations
	// and deciimal precisions of ADA lovelace and native assets.
	// SEE: https://github.com/howijd/decimal
	// issues and bug reports are welcome to:
	// https://github.com/howijd/decimal/issues.
	Lovelace decimal.Decimal
)
