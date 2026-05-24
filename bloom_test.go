package bloom_test

import (
	"fmt"
	"testing"

	"bloom"
)

func TestBloom(t *testing.T) {
	bl := bloom.New(100)

	bl.Set("foo")

	res := bl.Get("foo")
	if !res {
		t.Errorf("foo is missing, should be in the bitset")
	}
	return

	res = bl.Get("foo2")
	if res {
		t.Errorf("foo2 is present, but should not be in the bitset")
	}
}

// TestFalsePositiveRate inserts a small number of items into a generously
// sized filter, then queries a large set of items that were never inserted.
// A correctly implemented filter of this size should report very few of them
// as present. If the byte-index math collapses every bit into a handful of
// bytes, the filter saturates and reports almost everything as present.
func TestFalsePositiveRate(t *testing.T) {
	// 100 bytes = 800 bits. Plenty of room for 50 entries (~100 bits set).
	bl := bloom.New(100)

	const inserted = 50
	for i := 0; i < inserted; i++ {
		bl.Set(fmt.Sprintf("inserted-%d", i))
	}

	const queries = 1000
	falsePositives := 0
	for i := 0; i < queries; i++ {
		if bl.Get(fmt.Sprintf("absent-%d", i)) {
			falsePositives++
		}
	}

	rate := float64(falsePositives) / float64(queries)
	const maxRate = 0.10 // a correct 800-bit filter sits around ~1.5%
	if rate > maxRate {
		t.Errorf("false-positive rate %.1f%% (%d/%d) exceeds %.0f%% — bits are not spread across all %d bytes",
			rate*100, falsePositives, queries, maxRate*100, 100)
	}
}
