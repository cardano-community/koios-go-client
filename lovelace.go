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
	"github.com/shopspring/decimal"
)

func NewLovelace(value int64, exp int32) Lovelace {
	return Lovelace{Decimal: decimal.New(value, exp)}
}

func NewLovelaceFromString(value string) (Lovelace, error) {
	val, err := decimal.NewFromString(value)
	return Lovelace{Decimal: val}, err
}

// func (l *Lovelace) UnmarshalJSON(decimalBytes []byte) error {
// 	dec := &decimal.Decimal{}
// 	err := dec.UnmarshalJSON(decimalBytes)
// 	if err != nil {
// 		return err
// 	}
// 	*l, err = NewFromString(dec.String())
// 	return err
// }

// func NewLovelaceFromBigInt(value *big.Int, exp int32) Lovelace
// func NewLovelaceFromFloat(value float64) Lovelace
// func NewLovelaceFromFloat32(value float32) Lovelace
// func NewLovelaceFromFloatWithExponent(value float64, exp int32) Lovelace
// func NewLovelaceFromFormattedString(value string, replRegexp *regexp.Regexp) (Lovelace, error)
// func NewLovelaceFromInt(value int64) Lovelace
// func NewLovelaceFromInt32(value int32) Lovelace
// func NewLovelaceFromString(value string) (Lovelace, error)
