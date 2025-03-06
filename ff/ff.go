package ff

import (
	"fmt"
	"strconv"
)

// Helpers for finite field GF(2^8)

// bits 7-0 are X^n
type Poly byte

func NewFromUint16(r uint16) Poly {
	aesX8 := uint16(0b00011011)
	for i := 15; i >= 8; i-- {
		mask := uint16(1 << i)
		if r&mask > 0 {
			r += aesX8 << (i - 8)
		}
		r = r & ^mask
	}
	if r > 255 {
		panic(fmt.Sprintf("logic error: r %04x", r))
	}
	return Poly(r)
}

func (p Poly) String() string {
	return fmt.Sprintf("%08s", strconv.FormatInt(int64(p), 2))
}

func (p Poly) Add(q Poly) Poly {
	return p ^ q
}

func (p Poly) mulWithHighBits(q Poly) uint16 {
	var r uint16

	mask := Poly(1)
	for i := range 8 {
		if p&mask > 0 {
			r ^= uint16(q) << i
		}
		mask <<= 1
	}
	return r
}

// https://davidwong.fr/blockbreakers/aes_3_rcon.html#math
func (p Poly) Mul(q Poly) Poly {

	r := p.mulWithHighBits(q)
	return NewFromUint16(r)
}

func (p Poly) MulByX() Poly {
	r := uint16(p)
	r <<= 1
	return NewFromUint16(r)
}
