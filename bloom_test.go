package bloom_test

import (
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
