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

package testsuite

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNetworkTip(t *testing.T) {
	client, err := getClient()

	if assert.NoError(t, err) {
		tip, err := client.GetTip(context.Background(), nil)
		assert.NoError(t, err)

		assert.Greater(t, tip.Data.AbsSlot, int(100000))
		assert.Greater(t, tip.Data.BlockNo, int(100000))
		assert.Greater(t, tip.Data.EpochNo, int(230))
		assert.Greater(t, tip.Data.EpochSlot, int(1))
		assert.NotEmpty(t, tip.Data.BlockHash)
		assert.NotZero(t, tip.Data.BlockTime)
		assert.False(t, tip.Data.BlockTime.IsZero())
	}
}
