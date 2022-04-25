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
	"math/big"
	"regexp"

	"github.com/shopspring/decimal"
)

func NewLovelace(value int64, exp int32) Lovelace {
	return Lovelace{Decimal: decimal.New(value, exp)}
}

func NewLovelaceFromString(value string) (Lovelace, error) {
	val, err := decimal.NewFromString(value)
	return Lovelace{Decimal: val}, err
}

func NewLovelaceFromBigInt(value *big.Int, exp int32) Lovelace {
	return Lovelace{Decimal: decimal.NewFromBigInt(value, exp)}
}

func NewLovelaceFromFloat(value float64) Lovelace {
	return Lovelace{Decimal: decimal.NewFromFloat(value)}
}

func NewLovelaceFromFloat32(value float32) Lovelace {
	return Lovelace{Decimal: decimal.NewFromFloat32(value)}
}

func NewLovelaceFromFloatWithExponent(value float64, exp int32) Lovelace {
	return Lovelace{Decimal: decimal.NewFromFloatWithExponent(value, exp)}
}

func NewLovelaceFromFormattedString(value string, replRegexp *regexp.Regexp) (Lovelace, error) {
	val, err := decimal.NewFromFormattedString(value, replRegexp)
	return Lovelace{Decimal: val}, err
}

func NewLovelaceFromInt(value int64) Lovelace {
	return Lovelace{Decimal: decimal.NewFromInt(value)}
}

func NewLovelaceFromInt32(value int32) Lovelace {
	return Lovelace{Decimal: decimal.NewFromInt32(value)}
}
