// Copyright 2022 The Cardano Community Authors
// SPDX-License-Identifier: Apache-2.0
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
//
//	http://www.apache.org/licenses/LICENSE-2.0
//	or LICENSE file in repository root.
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package koios

import (
	"fmt"
	"net/http"
	"net/url"
)

// RequestOptions for the request.
type RequestOptions struct {
	page     uint
	pageSize uint
	locked   bool
	query    url.Values
	headers  http.Header
}

// QuerySet sets the key to value in request query.
// It replaces any existing values.
func (ro *RequestOptions) QuerySet(key, val string) {
	ro.query.Set(key, val)
}

// QueryAdd adds the value to request query by key.
// It appends to any existing values associated with key.
func (ro *RequestOptions) QueryAdd(key, val string) {
	ro.query.Add(key, val)
}

// QueryApply sets all values from provided query.
func (ro *RequestOptions) QueryApply(h url.Values) {
	for k, vv := range h {
		for _, v := range vv {
			ro.QuerySet(k, v)
		}
	}
}

// HeadersSet sets the key to value in request headers.
// It replaces any existing values.
func (ro *RequestOptions) HeaderSet(key, val string) {
	ro.headers.Set(key, val)
}

// HeadersAdd adds the value to request headers by key.
// It appends to any existing values associated with key.
func (ro *RequestOptions) HeaderAdd(key, val string) {
	ro.headers.Add(key, val)
}

// HeadersApply sets all values from provided header.
func (ro *RequestOptions) HeaderApply(h http.Header) {
	for k, vv := range h {
		for _, v := range vv {
			ro.HeaderSet(k, v)
		}
	}
}

// Clone the request options for using it with other request.
func (ro *RequestOptions) Clone() *RequestOptions {
	opts := &RequestOptions{
		headers:  ro.headers.Clone(),
		page:     ro.page,
		pageSize: ro.pageSize,
		query:    ro.query,
		locked:   false,
	}
	q := url.Values{}
	for k, v := range ro.query {
		q[k] = v
	}
	return opts
}

// SetPageSize for request modifies range header to satisfy requested page size.
func (ro *RequestOptions) SetPageSize(size uint) {
	ro.pageSize = size
}

// SetCurrentPage modifies range header of the request to satisfy current page requested.
func (ro *RequestOptions) SetCurrentPage(page uint) {
	ro.page = page
}

// lock the request options.
func (ro *RequestOptions) lock() error {
	if ro.locked {
		return ErrReqOptsAlreadyUsed
	}
	ro.locked = true
	if ro.pageSize != PageSize || ro.page != 1 {
		e := (ro.pageSize * ro.page) - 1
		s := (e + 1) - ro.pageSize
		ro.headers.Set("Range", fmt.Sprintf("%d-%d", s, e))
	}
	return nil
}
