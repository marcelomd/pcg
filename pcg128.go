package pcg

import (
	"math/bits"
)

// 128 state
type Pcg128 struct {
	mulHigh   uint64
	mulLow    uint64
	incHigh   uint64
	incLow    uint64
	stateHigh uint64
	stateLow  uint64
	step      func()
	output64  func() uint64
	output128 func() (uint64, uint64)
}

const (
	mulHigh128   uint64 = 0x2360ed051fc65da4 // 2549297995355413924
	mulLow128    uint64 = 0x4385df649fccf645 // 4865540595714422341
	incHigh128   uint64 = 0x5851f42d4c957f2d // 6364136223846793005
	incLow128    uint64 = 0x14057b7ef767814f // 1442695040888963407
	stateHigh128 uint64 = 1
	stateLow128  uint64 = 1
)

func NewPcg128(mulHigh, mulLow, incHigh, incLow, stateHigh, stateLow uint64) *Pcg128 {
	return &Pcg128{
		mulHigh:   mulHigh,
		mulLow:    mulLow,
		incHigh:   incHigh,
		incLow:    incLow,
		stateHigh: stateHigh,
		stateLow:  stateLow,
	}
}

func NewDefaultPcg128() *Pcg128 {
	return &Pcg128{
		mulHigh:   mulHigh128,
		mulLow:    mulLow128,
		incHigh:   incHigh128,
		incLow:    incLow128,
		stateHigh: stateHigh128,
		stateLow:  stateLow128,
	}
}

func (p *Pcg128) Copy() *Pcg128 {
	return &Pcg128{
		mulHigh:   p.mulHigh,
		mulLow:    p.mulLow,
		incHigh:   p.incHigh,
		incLow:    p.incLow,
		stateHigh: p.stateHigh,
		stateLow:  p.stateLow,
	}
}

func (p *Pcg128) mul() {
	hi, lo := bits.Mul64(p.stateLow, p.mulLow)
	hi += p.stateHigh * p.mulLow
	hi += p.stateLow * p.mulHigh
	p.stateLow = lo
	p.stateHigh = hi
}

func (p *Pcg128) inc() {
	var carry uint64
	p.stateLow, carry = bits.Add64(p.stateLow, p.incLow, 0)
	p.stateHigh, _ = bits.Add64(p.stateHigh, p.incHigh, carry)
}

//----- MCG
func (p *Pcg128) mcg_128_step() {
	p.mul()
}

func (p *Pcg128) mcg_128_seed(stateHigh, stateLow uint64) {
	p.stateHigh = stateHigh
	p.stateLow = stateLow | 1
}

//----- LCG
func (p *Pcg128) lcg_128_step() {
	p.mul()
	p.inc()
}

func (p *Pcg128) lcg_128_seed(stateHigh, stateLow uint64) {
	p.stateHigh = 0
	p.stateLow = 0
	p.lcg_128_step()
	var carry uint64
	p.stateLow, carry = bits.Add64(p.stateLow, stateLow, 0)
	p.stateHigh, _ = bits.Add64(p.stateHigh, stateHigh, carry)
	p.lcg_128_step()
}

func (p *Pcg128) lcg_128_setseq_seed(stateHigh, stateLow, seqHigh, seqLow uint64) {
	p.stateHigh = 0
	p.stateLow = 0
	p.incHigh = seqHigh
	p.incLow = (seqLow << 1) | 1
	p.lcg_128_step()
	var carry uint64
	p.stateLow, carry = bits.Add64(p.stateLow, stateLow, 0)
	p.stateHigh, _ = bits.Add64(p.stateHigh, stateHigh, carry)
	p.lcg_128_step()
}

//----- PCG
func (p *Pcg128) output_xsl_rr_128_64() uint64 {
	return bits.RotateLeft64(p.stateHigh^p.stateLow, -int(p.stateHigh>>58))
}

func (p *Pcg128) output_xsl_rr_rr_128_128() (uint64, uint64) {
	low := bits.RotateLeft64(p.stateHigh^p.stateLow, -int(p.stateHigh>>58))
	high := bits.RotateLeft64(p.stateHigh, -int(low&63))
	return high, low
}

//----- Default lcg, oneseq, XSL RR
func (p *Pcg128) Seed(stateHigh, stateLow uint64) {
	p.lcg_128_seed(stateHigh, stateLow)
}

func (p *Pcg128) Uint64() uint64 {
	p.lcg_128_step()
	return p.output_xsl_rr_128_64()
}

func (p *Pcg128) Uint128() (uint64, uint64) {
	p.lcg_128_step()
	return p.output_xsl_rr_rr_128_128()
}

var defaultPcg128 = NewDefaultPcg128()

func Seed128(stateHigh, stateLow uint64) {
	defaultPcg128.Seed(stateHigh, stateLow)
}

func Uint64() uint64 {
	return defaultPcg128.Uint64()
}

func Uint128() (uint64, uint64) {
	return defaultPcg128.Uint128()
}
