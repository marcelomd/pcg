package pcg

import (
	"math/bits"
)

// PCG XSH RS 64 state / 32 output
type Pcg64 struct {
	mul   uint64
	inc   uint64
	state uint64
}

const (
	mul64   uint64 = 0x5851f42d4c957f2d // 6364136223846793005
	inc64   uint64 = 0x14057b7ef767814f // 1442695040888963407
	state64 uint64 = 1
)

func NewPcg64(mul, inc, state uint64) *Pcg64 {
	return &Pcg64{mul: mul, inc: inc, state: state}
}

func NewDefaultPcg64() *Pcg64 {
	return &Pcg64{mul: mul64, inc: inc64, state: state64}
}

func (p *Pcg64) Copy() *Pcg64 {
	return &Pcg64{mul: p.mul, inc: p.inc, state: p.state}
}

//----- MCG
func (p *Pcg64) mcg_64_step() {
	p.state = p.state * p.mul
}

func (p *Pcg64) mcg_64_seed(state uint64) {
	p.state = state | 1
}

//----- LCG
func (p *Pcg64) lcg_64_step() {
	p.state = p.state*p.mul + p.inc
}

func (p *Pcg64) lcg_64_seed(state uint64) {
	p.state = 0
	p.lcg_64_step()
	p.state += state
	p.lcg_64_step()
}

func (p *Pcg64) lcg_setseq_64_seed(state, seq uint64) {
	p.state = 0
	p.inc = (seq << 1) | 1
	p.lcg_64_step()
	p.state += state
	p.lcg_64_step()
}

//----- PCG
func (p *Pcg64) output_xsh_rr_64_32(state uint64) uint32 {
	return bits.RotateLeft32(uint32(((state>>18)^state)>>27), -int(state>>59))
}

func (p *Pcg64) output_xsh_rs_64_32(state uint64) uint32 {
	return uint32(((state >> 22) ^ state) >> ((state >> 61) + 22))
}

//----- Default lcg, oneseq, XSH RS
func (p *Pcg64) Seed(state uint64) {
	p.lcg_64_seed(state)
}

func (p *Pcg64) Uint32() uint32 {
	state := p.state
	p.lcg_64_step()
	return p.output_xsh_rs_64_32(state)
}

var defaultPcg64 = NewDefaultPcg64()

func Seed64(state uint64) {
	defaultPcg64.Seed(state)
}

func Uint32() uint32 {
	return defaultPcg64.Uint32()
}
