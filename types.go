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

type (
	// Tip defines model for tip.
	Tip []struct {
		// Absolute Slot number (slots not divided into epochs)
		AbsSlot int `json:"abs_slot"`

		// Block Height number on chain
		BlockNo int `json:"block_no"`

		// Timestamp for when the block was created
		BlockTime string `json:"block_time"`

		// Epoch number
		Epoch int `json:"epoch"`

		// Slot number within Epoch
		EpochSlot int `json:"epoch_slot"`

		// Block Hash in hex
		Hash string `json:"hash"`
	}

	// TipResponse wraps response of /tip
	TipResponse struct {
		Response
		Tip Tip `json:"response"`
	}
)
