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

import (
	"context"
	"encoding/json"
)

// GetTip returns the tip info about the latest block seen by chain.
func (c *Client) GetTip(ctx context.Context) (*TipResponse, error) {
	rsp, err := c.GET(ctx, "/tip")
	if err != nil {
		return nil, err
	}
	res := &TipResponse{}
	res.setStatus(rsp)
	body, err := readResponseBody(rsp)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(body, &res.Tip); err != nil {
		return nil, err
	}

	return res, nil
}
