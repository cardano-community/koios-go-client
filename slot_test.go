// SPDX-License-Identifier: Apache-2.0
//
// Copyright © 2026 The Cardano Community Authors

package koios_test

import (
	"encoding/json"
	"errors"
	"strconv"
	"testing"

	koios "github.com/cardano-community/koios-go-client/v5"
)

func TestSlotUnmarshalJSONBounds(t *testing.T) {
	max := ^uint(0)
	overflow := "18446744073709551616"
	if strconv.IntSize == 32 {
		overflow = "4294967296"
	}
	for _, quoted := range []bool{false, true} {
		for _, value := range []uint{0, 42331172, max - 1, max} {
			raw := strconv.FormatUint(uint64(value), 10)
			if quoted {
				raw = strconv.Quote(raw)
			}
			t.Run(raw, func(t *testing.T) {
				var slot koios.Slot
				if err := json.Unmarshal([]byte(raw), &slot); err != nil {
					t.Fatal(err)
				}
				if slot != koios.Slot(value) {
					t.Fatalf("slot = %d, want %d", slot, value)
				}
			})
		}
		raw := overflow
		if quoted {
			raw = strconv.Quote(raw)
		}
		t.Run(raw, func(t *testing.T) {
			slot := koios.Slot(42)
			if err := json.Unmarshal([]byte(raw), &slot); !errors.Is(err, strconv.ErrRange) {
				t.Fatalf("overflow error = %v, want %v", err, strconv.ErrRange)
			}
			if slot != 42 {
				t.Fatalf("overflow changed slot to %d", slot)
			}
		})
	}
}
