package bloom_test

import (
	"errors"
	"fmt"
	"io"
	"testing"

	"bloom"
)

var benchmarkResult bool

func TestBloom(t *testing.T) {
	bl := bloom.New(100)

	bl.Set("foo")

	res := bl.Get("foo")
	if !res {
		t.Errorf("foo is missing, should be in the bitset")
	}

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

func BenchmarkSet(b *testing.B) {
	keys := benchmarkKeys(1024)

	for _, size := range []int{100, 10_000, 1_000_000} {
		b.Run(fmt.Sprintf("filter_bytes=%d", size), func(b *testing.B) {
			bl := bloom.New(size)
			b.ReportAllocs()
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				bl.Set(keys[i%len(keys)])
			}
		})
	}
}

func BenchmarkGetHit(b *testing.B) {
	keys := benchmarkKeys(1024)

	for _, size := range []int{100, 10_000, 1_000_000} {
		b.Run(fmt.Sprintf("filter_bytes=%d", size), func(b *testing.B) {
			bl := bloom.New(size)
			for _, key := range keys {
				bl.Set(key)
			}

			b.ReportAllocs()
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				benchmarkResult = bl.Get(keys[i%len(keys)])
			}
		})
	}
}

func BenchmarkGetMiss(b *testing.B) {
	insertedKeys := benchmarkKeys(1024)
	missingKeys := benchmarkKeysFrom("missing", 1024)

	for _, size := range []int{100, 10_000, 1_000_000} {
		b.Run(fmt.Sprintf("filter_bytes=%d", size), func(b *testing.B) {
			bl := bloom.New(size)
			for _, key := range insertedKeys {
				bl.Set(key)
			}

			b.ReportAllocs()
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				benchmarkResult = bl.Get(missingKeys[i%len(missingKeys)])
			}
		})
	}
}

func benchmarkKeys(n int) []string {
	return benchmarkKeysFrom("payload", n)
}

func benchmarkKeysFrom(prefix string, n int) []string {
	keys := make([]string, n)
	for i := 0; i < n; i++ {
		keys[i] = fmt.Sprintf("%s-%08d", prefix, i)
	}

	return keys
}

func TestReader(t *testing.T) {
	bl := bloom.New(100)

	bl.Set("foo")

	p := make([]byte, 3)
	for {
		n, err := bl.Read(p)
		if n == 0 && err == nil {
			fmt.Println("nothing has been read, aborting")

			break
		}

		if n == 0 && errors.Is(err, io.EOF) {
			fmt.Println("all good!")

			break
		}

		fmt.Printf("have read %d bytes\n", n)
	}
}

func TestWriter(t *testing.T) {
	bl := bloom.New(10)

	p := make([]byte, 3)
	for {
		n, err := bl.Write(p)
		if err != nil {
			fmt.Println(err)

			break
		}

		fmt.Printf("have written %d bytes\n", n)
	}
}
