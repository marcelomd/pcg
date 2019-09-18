package pcg

import (
	"math/rand"
	"testing"
)

func BenchmarkPcg64_32(b *testing.B) {
	pcg64 := NewDefaultPcg64()
	pcg64.Seed(1)
	for n := 0; n < b.N; n++ {
		pcg64.Uint32()
	}
}

func BenchmarkPcg128_64(b *testing.B) {
	pcg128 := NewDefaultPcg128()
	pcg128.Seed(1, 1)
	for n := 0; n < b.N; n++ {
		pcg128.Uint64()
	}
}

func BenchmarkPcg128_128(b *testing.B) {
	pcg128 := NewDefaultPcg128()
	pcg128.Seed(1, 1)
	for n := 0; n < b.N; n++ {
		pcg128.Uint128()
	}
}

func BenchmarkRandUint32(b *testing.B) {
	r := rand.New(rand.NewSource(1))
	for n := 0; n < b.N; n++ {
		r.Uint32()
	}
}

func BenchmarkRandUint64(b *testing.B) {
	r := rand.New(rand.NewSource(1))
	for n := 0; n < b.N; n++ {
		r.Uint64()
	}
}

func BenchmarkRandRead80(b *testing.B) {
	r := rand.New(rand.NewSource(1))
	for n := 0; n < b.N; n++ {
		b := [10]byte{}
		r.Read(b[:])
	}
}
