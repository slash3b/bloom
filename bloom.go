package bloom

import (
	"crypto/md5"
	"encoding/binary"
)

type Filter struct {
	bits []byte
	size int
}

func New(s int) Filter {
	return Filter{
		bits: make([]byte, s),
		size: s,
	}
}

func (f *Filter) indexes(s string) []int {
	res := md5.Sum([]byte(s))

	a := int(binary.LittleEndian.Uint64(res[0:md5.Size/2]) % uint64(f.size*8))
	a2 := a / f.size
	a3 := a % 8

	b := int(binary.LittleEndian.Uint64(res[md5.Size/2:]) % uint64(f.size*8))
	b2 := b / f.size
	b3 := b % 8

	return []int{a2, a3, b2, b3}
}

// Set adds entry into the Filter
func (f *Filter) Set(s string) {
	idxs := f.indexes(s)

	f.bits[idxs[0]] |= 1 << idxs[1]
	f.bits[idxs[2]] |= 1 << idxs[3]
}

// Get checks if element exists.
func (f *Filter) Get(s string) bool {
	idxs := f.indexes(s)

	for i := 0; i < len(idxs); i = i + 2 {
		if f.bits[idxs[i]]&(1<<idxs[i+1]) == 0 {
			return false
		}
	}

	return true
}
