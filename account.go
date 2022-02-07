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
	"net/http"
)

type (

	// AccountListResponse represents response from `/account_list` endpoint.
	AccountListResponse struct {
		Response
		Data []StakeAddress `json:"response"`
	}
)

// GetTip returns the tip info about the latest block seen by chain.
func (c *Client) GetAccountList(ctx context.Context) (res *AccountListResponse, err error) {
	res = &AccountListResponse{}
	rsp, err := c.request(ctx, &res.Response, "GET", nil, "/account_list", nil, nil)
	if err != nil {
		res.applyError(nil, err)
		return
	}

	body, err := readResponseBody(rsp)
	if err != nil {
		res.applyError(body, err)
		return
	}

	accs := []struct {
		ID StakeAddress `json:"id"`
	}{}

	if err = json.Unmarshal(body, &accs); err != nil {
		res.applyError(body, err)
		return
	}

	if rsp.StatusCode != http.StatusOK {
		res.applyError(body, err)
		return
	}
	if len(accs) > 0 {
		for _, a := range accs {
			res.Data = append(res.Data, a.ID)
		}
	}
	return res, nil
}
