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

package internal

import (
	"net/http"
	"net/url"
)

// APITestRequestSpec used to define requests for testdata.
type APITestRequestSpec struct {
	Method string      `json:"method"`
	Query  url.Values  `json:"query,omitempty"`
	Body   []byte      `json:"body,omitempty"`
	Header http.Header `json:"header,omitempty"`
}

// APITestRequestSpec used to define responses for testdata.
type APITestResponseSpec struct {
	Code   int         `json:"code"`
	Header http.Header `json:"headers"`
	Body   any         `json:"body"`
}

// APITestSpec is testdata spec.
type APITestSpec struct {
	Network  string              `json:"network"`
	Filename string              `json:"filename"`
	Endpoint string              `json:"endpoint"`
	Request  APITestRequestSpec  `json:"request"`
	Response APITestResponseSpec `json:"response"`
}
