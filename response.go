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
package koios

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// ReadResponseBody is reading http.Response aand closing it after read.
func ReadResponseBody(rsp *http.Response) (body []byte, err error) {
	if rsp == nil {
		return nil, nil
	}

	defer func() { _ = rsp.Body.Close() }()

	return io.ReadAll(rsp.Body)
}

// ReadAndUnmarshalResponse is helper to unmarchal json responses.
func ReadAndUnmarshalResponse(rsp *http.Response, res *Response, dest interface{}) error {
	body, err := ReadResponseBody(rsp)
	if !strings.Contains(rsp.Header.Get("Content-Type"), "json") {
		return fmt.Errorf("%w: %s", ErrResponseIsNotJSON, string(body))
	}

	res.applyError(body, err)
	if len(body) == 0 || err != nil {
		return err
	}

	defer res.ready()
	err = json.Unmarshal(body, dest)
	res.applyError(body, err)
	return err
}

func (r *Response) applyError(body []byte, err error) {
	if err == nil {
		return
	}

	if r.Error == nil {
		r.Error = &ResponseError{}
	}
	if len(body) != 0 {
		berr := json.Unmarshal(body, r.Error)
		if berr != nil {
			r.Error.Message = berr.Error()
		}
	}
	defer r.ready()

	if len(r.Error.Message) == 0 {
		r.Error.Message = err.Error()
	} else {
		r.Error.Message = fmt.Sprintf("%s: %s", err.Error(), r.Error.Message)
	}
}

func (r *Response) ready() {
	if r.Stats == nil {
		return
	}
	r.Stats.ReqDur = time.Since(r.Stats.ReqStartedAt)
	r.Stats.ReqDurStr = fmt.Sprint(r.Stats.ReqDur)
}

func (r *Response) applyRsp(rsp *http.Response) {
	r.StatusCode = rsp.StatusCode
	r.RequestMethod = rsp.Request.Method
	r.Status = rsp.Status
	r.Date = rsp.Header.Get("date")
	r.ContentRange = rsp.Header.Get("content-range")
	r.ContentLocation = rsp.Header.Get("content-location")
}
