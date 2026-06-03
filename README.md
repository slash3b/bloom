# bloom

A Bloom filter implementation in Go. It uses two xxhash functions with different seeds to map elements into a compact bitset, enabling fast probabilistic membership checks with no false negatives. The filter also implements `io.Reader` and `io.Writer` so its internal state can be serialized and restored.

# run benchmarks

 go test -count=10 -benchtime=2s -bench ./...
