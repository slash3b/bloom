package bloom

import (
	"io"

	"github.com/OneOfOne/xxhash"
)

type Filter struct {
	bits    []byte
	size    int
	bitsize uint64

	ioReaderOffset int
}

func New(s int) Filter {
	return Filter{
		bits:    make([]byte, s),
		size:    s,
		bitsize: uint64(s * 8),
	}
}

func (f *Filter) Read(p []byte) (int, error) {
	n := copy(p, f.bits[f.ioReaderOffset:])

	if f.ioReaderOffset >= f.size {
		return 0, io.EOF
	}

	if n == 0 {
		return 0, nil
	}

	f.ioReaderOffset += n

	return n, nil
}

// Set adds entry into the Filter
func (f *Filter) Set(s string) {
	h1 := xxhash.ChecksumString64S(s, 0)
	h2 := xxhash.ChecksumString64S(s, 42)

	a := int(h1 % f.bitsize)
	b := int(h2 % f.bitsize)

	f.bits[a/8] |= 1 << (a % 8)
	f.bits[b/8] |= 1 << (b % 8)
}

// Get checks if element exists.
func (f *Filter) Get(s string) bool {
	h1 := xxhash.ChecksumString64S(s, 0)
	h2 := xxhash.ChecksumString64S(s, 42)

	a := int(h1 % f.bitsize)
	b := int(h2 % f.bitsize)

	if f.bits[a/8]&(1<<(a%8)) == 0 {
		return false
	}

	if f.bits[b/8]&(1<<(b%8)) == 0 {
		return false
	}

	return true
}
