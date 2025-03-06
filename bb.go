package bb

import (
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
)

// 00010203 -> 01020300
func RotWord(n uint32) uint32 {
	b := n & 0xff000000
	b >>= 24
	n <<= 8
	n |= b
	return n
}

var sbox_en = [256]byte{
	0x63, 0x7c, 0x77, 0x7b, 0xf2, 0x6b, 0x6f, 0xc5, 0x30, 0x01, 0x67, 0x2b, 0xfe, 0xd7, 0xab, 0x76,
	0xca, 0x82, 0xc9, 0x7d, 0xfa, 0x59, 0x47, 0xf0, 0xad, 0xd4, 0xa2, 0xaf, 0x9c, 0xa4, 0x72, 0xc0,
	0xb7, 0xfd, 0x93, 0x26, 0x36, 0x3f, 0xf7, 0xcc, 0x34, 0xa5, 0xe5, 0xf1, 0x71, 0xd8, 0x31, 0x15,
	0x04, 0xc7, 0x23, 0xc3, 0x18, 0x96, 0x05, 0x9a, 0x07, 0x12, 0x80, 0xe2, 0xeb, 0x27, 0xb2, 0x75,
	0x09, 0x83, 0x2c, 0x1a, 0x1b, 0x6e, 0x5a, 0xa0, 0x52, 0x3b, 0xd6, 0xb3, 0x29, 0xe3, 0x2f, 0x84,
	0x53, 0xd1, 0x00, 0xed, 0x20, 0xfc, 0xb1, 0x5b, 0x6a, 0xcb, 0xbe, 0x39, 0x4a, 0x4c, 0x58, 0xcf,
	0xd0, 0xef, 0xaa, 0xfb, 0x43, 0x4d, 0x33, 0x85, 0x45, 0xf9, 0x02, 0x7f, 0x50, 0x3c, 0x9f, 0xa8,
	0x51, 0xa3, 0x40, 0x8f, 0x92, 0x9d, 0x38, 0xf5, 0xbc, 0xb6, 0xda, 0x21, 0x10, 0xff, 0xf3, 0xd2,
	0xcd, 0x0c, 0x13, 0xec, 0x5f, 0x97, 0x44, 0x17, 0xc4, 0xa7, 0x7e, 0x3d, 0x64, 0x5d, 0x19, 0x73,
	0x60, 0x81, 0x4f, 0xdc, 0x22, 0x2a, 0x90, 0x88, 0x46, 0xee, 0xb8, 0x14, 0xde, 0x5e, 0x0b, 0xdb,
	0xe0, 0x32, 0x3a, 0x0a, 0x49, 0x06, 0x24, 0x5c, 0xc2, 0xd3, 0xac, 0x62, 0x91, 0x95, 0xe4, 0x79,
	0xe7, 0xc8, 0x37, 0x6d, 0x8d, 0xd5, 0x4e, 0xa9, 0x6c, 0x56, 0xf4, 0xea, 0x65, 0x7a, 0xae, 0x08,
	0xba, 0x78, 0x25, 0x2e, 0x1c, 0xa6, 0xb4, 0xc6, 0xe8, 0xdd, 0x74, 0x1f, 0x4b, 0xbd, 0x8b, 0x8a,
	0x70, 0x3e, 0xb5, 0x66, 0x48, 0x03, 0xf6, 0x0e, 0x61, 0x35, 0x57, 0xb9, 0x86, 0xc1, 0x1d, 0x9e,
	0xe1, 0xf8, 0x98, 0x11, 0x69, 0xd9, 0x8e, 0x94, 0x9b, 0x1e, 0x87, 0xe9, 0xce, 0x55, 0x28, 0xdf,
	0x8c, 0xa1, 0x89, 0x0d, 0xbf, 0xe6, 0x42, 0x68, 0x41, 0x99, 0x2d, 0x0f, 0xb0, 0x54, 0xbb, 0x16,
}

func BytesToWord(bs [4]byte) uint32 {
	return uint32(bs[0])<<24 + uint32(bs[1])<<16 + uint32(bs[2])<<8 + uint32(bs[3])<<0
}

func WordToBytes(n uint32) [4]byte {
	return [4]byte{
		byte((n & 0xff000000) >> 24),
		byte((n & 0x00ff0000) >> 16),
		byte((n & 0x0000ff00) >> 8),
		byte((n & 0x000000ff) >> 0),
	}
}

func SubWord(n uint32) uint32 {
	bs := WordToBytes(n)
	for i := range bs {
		bs[i] = sbox_en[bs[i]]
	}
	return BytesToWord(bs)
}

var rcon = [256]byte{
	0x8d, 0x01, 0x02, 0x04, 0x08, 0x10, 0x20, 0x40, 0x80, 0x1b, 0x36, 0x6c, 0xd8, 0xab, 0x4d, 0x9a,
	0x2f, 0x5e, 0xbc, 0x63, 0xc6, 0x97, 0x35, 0x6a, 0xd4, 0xb3, 0x7d, 0xfa, 0xef, 0xc5, 0x91, 0x39,
	0x72, 0xe4, 0xd3, 0xbd, 0x61, 0xc2, 0x9f, 0x25, 0x4a, 0x94, 0x33, 0x66, 0xcc, 0x83, 0x1d, 0x3a,
	0x74, 0xe8, 0xcb, 0x8d, 0x01, 0x02, 0x04, 0x08, 0x10, 0x20, 0x40, 0x80, 0x1b, 0x36, 0x6c, 0xd8,
	0xab, 0x4d, 0x9a, 0x2f, 0x5e, 0xbc, 0x63, 0xc6, 0x97, 0x35, 0x6a, 0xd4, 0xb3, 0x7d, 0xfa, 0xef,
	0xc5, 0x91, 0x39, 0x72, 0xe4, 0xd3, 0xbd, 0x61, 0xc2, 0x9f, 0x25, 0x4a, 0x94, 0x33, 0x66, 0xcc,
	0x83, 0x1d, 0x3a, 0x74, 0xe8, 0xcb, 0x8d, 0x01, 0x02, 0x04, 0x08, 0x10, 0x20, 0x40, 0x80, 0x1b,
	0x36, 0x6c, 0xd8, 0xab, 0x4d, 0x9a, 0x2f, 0x5e, 0xbc, 0x63, 0xc6, 0x97, 0x35, 0x6a, 0xd4, 0xb3,
	0x7d, 0xfa, 0xef, 0xc5, 0x91, 0x39, 0x72, 0xe4, 0xd3, 0xbd, 0x61, 0xc2, 0x9f, 0x25, 0x4a, 0x94,
	0x33, 0x66, 0xcc, 0x83, 0x1d, 0x3a, 0x74, 0xe8, 0xcb, 0x8d, 0x01, 0x02, 0x04, 0x08, 0x10, 0x20,
	0x40, 0x80, 0x1b, 0x36, 0x6c, 0xd8, 0xab, 0x4d, 0x9a, 0x2f, 0x5e, 0xbc, 0x63, 0xc6, 0x97, 0x35,
	0x6a, 0xd4, 0xb3, 0x7d, 0xfa, 0xef, 0xc5, 0x91, 0x39, 0x72, 0xe4, 0xd3, 0xbd, 0x61, 0xc2, 0x9f,
	0x25, 0x4a, 0x94, 0x33, 0x66, 0xcc, 0x83, 0x1d, 0x3a, 0x74, 0xe8, 0xcb, 0x8d, 0x01, 0x02, 0x04,
	0x08, 0x10, 0x20, 0x40, 0x80, 0x1b, 0x36, 0x6c, 0xd8, 0xab, 0x4d, 0x9a, 0x2f, 0x5e, 0xbc, 0x63,
	0xc6, 0x97, 0x35, 0x6a, 0xd4, 0xb3, 0x7d, 0xfa, 0xef, 0xc5, 0x91, 0x39, 0x72, 0xe4, 0xd3, 0xbd,
	0x61, 0xc2, 0x9f, 0x25, 0x4a, 0x94, 0x33, 0x66, 0xcc, 0x83, 0x1d, 0x3a, 0x74, 0xe8, 0xcb, 0x8d,
}

func Rcon(b byte) uint32 {
	return BytesToWord([4]byte{rcon[b], 0, 0, 0})
}

type Key [16]byte

func (k Key) String() string {
	return hex.EncodeToString(k[:])
}

func ColsToKey(cols [4]uint32) Key {
	var k [16]byte
	for i, col := range cols {
		bs := WordToBytes(col)
		for j, b := range bs {
			k[i*4+j] = b
		}
	}
	return k
}

func KeyToCols(k Key) [4]uint32 {
	var cols [4]uint32
	for i := range cols {
		j := i * 4
		cols[i] = (uint32(k[j]) << 24) + (uint32(k[j+1]) << 16) + (uint32(k[j+2]) << 8) + (uint32(k[j+3]) << 0)
	}
	return cols
}

type ExpandedKey = [11]Key

func KeyExpansion(k Key) ExpandedKey {
	var expKey [11]Key
	for round := range expKey {
		if round == 0 {
			expKey[0] = k
			continue
		}
		prevCols := KeyToCols(expKey[round-1])
		v := RotWord(prevCols[3])
		v = SubWord(v)
		v ^= prevCols[0]
		v ^= Rcon(byte(round))
		var cols [4]uint32
		for j := range cols {
			if j == 0 {
				cols[j] = v
				continue
			}
			cols[j] = cols[j-1] ^ prevCols[j]
		}
		expKey[round] = ColsToKey(cols)
	}
	return expKey
}

type State [16]byte

func PlainBlockToState(ptxt string) (State, error) {
	var s State
	if len(ptxt) != 16 {
		return s, errors.New(fmt.Sprintf("plain text is not one block, length %d not %d", len(ptxt), 16))
	}
	bs := []byte(ptxt)
	for i := range bs {
		col := i / 4
		row := i % 4
		s[row*4+col] = bs[i]
	}
	return s, nil
}

func (s State) String() string {
	sb := strings.Builder{}
	for row := range 4 {
		sb.WriteString(fmt.Sprintf("%02x %02x %02x %02x\n", s[row*4+0], s[row*4+1], s[row*4+2], s[row*4+3]))
	}
	return sb.String()
}

func (state State) SubBytes() State {
	var ret State
	for i, b := range state {
		ret[i] = sbox_en[b]
	}
	return ret
}

func (state State) ShiftRows() State {
	shiftCol := func(step int, col []byte) []byte {
		var ret [4]byte
		for i, b := range col {
			j := (4 + i - step) % 4
			ret[j] = b
		}
		return ret[:]
	}
	var ret State
	for i := range 4 {
		copy(ret[i*4:(i+1)*4], shiftCol(i, state[i*4:(i+1)*4]))
	}
	return ret
}
