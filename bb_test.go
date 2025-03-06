package bb

import (
	"testing"
)

func TestRotWord(t *testing.T) {
	testCases := []struct {
		in       uint32
		expected uint32
	}{
		{0x00010203, 0x01020300},
		{0x00000000, 0x00000000},
		{0x01000000, 0x00000001},
		{0x00010000, 0x01000000},
		{0x00000100, 0x00010000},
		{0xfffefdfc, 0xfefdfcff},
	}

	for _, tc := range testCases {
		got := RotWord(tc.in)
		if got != tc.expected {
			t.Errorf("got %08x expected %08x", got, tc.expected)
		}
	}
}

func TestWordToBytes(t *testing.T) {
	testCases := []struct {
		in       uint32
		expected [4]byte
	}{
		{0x00010203, [4]byte{0x00, 0x01, 0x02, 0x03}},
		{0xfffefdfc, [4]byte{0xff, 0xfe, 0xfd, 0xfc}},
	}

	for _, tc := range testCases {
		got := WordToBytes(tc.in)
		if got != tc.expected {
			t.Errorf("got %08x expected %v", got, tc.expected)
		}
	}
}

func TestRcon(t *testing.T) {
	testCases := []struct {
		in       byte
		expected uint32
	}{
		{1, 0x01000000},
		{2, 0x02000000},
		{3, 0x04000000},
		{4, 0x08000000},
	}

	for _, tc := range testCases {
		got := Rcon(tc.in)
		if got != tc.expected {
			t.Errorf("got %08x expected %v", got, tc.expected)
		}
	}
}

func TestExpandKey(t *testing.T) {
	keyCols := [4]uint32{0x2b7e1516, 0x28aed2a6, 0xabf71588, 0x09cf4f3c}
	key := ColsToKey(keyCols)
	expected := [11]Key{
		ColsToKey([4]uint32{0x2b7e1516, 0x28aed2a6, 0xabf71588, 0x09cf4f3c}),
		ColsToKey([4]uint32{0xa0fafe17, 0x88542cb1, 0x23a33939, 0x2a6c7605}),
		ColsToKey([4]uint32{0xf2c295f2, 0x7a96b943, 0x5935807a, 0x7359f67f}),
		ColsToKey([4]uint32{0x3d80477d, 0x4716fe3e, 0x1e237e44, 0x6d7a883b}),
		ColsToKey([4]uint32{0xef44a541, 0xa8525b7f, 0xb671253b, 0xdb0bad00}),
		ColsToKey([4]uint32{0xd4d1c6f8, 0x7c839d87, 0xcaf2b8bc, 0x11f915bc}),
		ColsToKey([4]uint32{0x6d88a37a, 0x110b3efd, 0xdbf98641, 0xca0093fd}),
		ColsToKey([4]uint32{0x4e54f70e, 0x5f5fc9f3, 0x84a64fb2, 0x4ea6dc4f}),
		ColsToKey([4]uint32{0xead27321, 0xb58dbad2, 0x312bf560, 0x7f8d292f}),
		ColsToKey([4]uint32{0xac7766f3, 0x19fadc21, 0x28d12941, 0x575c006e}),
		ColsToKey([4]uint32{0xd014f9a8, 0xc9ee2589, 0xe13f0cc8, 0xb6630ca6}),
	}

	got := KeyExpansion(key)
	if got != expected {
		t.Errorf("got %v expected %v", got, expected)
	}
}

func TestSubBytes(t *testing.T) {
	in := State([]byte{
		0x00, 0x04, 0x08, 0x0c,
		0x01, 0x05, 0x09, 0x0d,
		0x02, 0x06, 0x0a, 0x0e,
		0x03, 0x07, 0x0b, 0x0f,
	})
	got := in.SubBytes()
	expected := State([]byte{
		0x63, 0xf2, 0x30, 0xfe,
		0x7c, 0x6b, 0x01, 0xd7,
		0x77, 0x6f, 0x67, 0xab,
		0x7b, 0xc5, 0x2b, 0x76,
	})

	if got != expected {
		t.Errorf("Subbytes failed")
	}
}

func TestShiftRows(t *testing.T) {
	in := State([]byte{
		0x63, 0xf2, 0x30, 0xfe,
		0x7c, 0x6b, 0x01, 0xd7,
		0x77, 0x6f, 0x67, 0xab,
		0x7b, 0xc5, 0x2b, 0x76,
	})
	got := in.ShiftRows()
	expected := State([]byte{
		0x63, 0xf2, 0x30, 0xfe,
		0x6b, 0x01, 0xd7, 0x7c,
		0x67, 0xab, 0x77, 0x6f,
		0x76, 0x7b, 0xc5, 0x2b,
	})

	if got != expected {
		t.Errorf("ShiftRows failed got:\n%s\nexpected\n%s\n", got, expected)
	}
}

func TestMixColumns(t *testing.T) {
	in := State([]byte{
		0x63, 0xf2, 0x30, 0xfe,
		0x6b, 0x01, 0xd7, 0x7c,
		0x67, 0xab, 0x77, 0x6f,
		0x76, 0x7b, 0xc5, 0x2b,
	})
	expected := State([]byte{
		0x6a, 0x2c, 0xb0, 0x27,
		0x6a, 0x6d, 0xd9, 0x9c,
		0x5c, 0x33, 0x5d, 0x21,
		0x45, 0x51, 0x61, 0x5c,
	})
	got := in.MixColumns()

	if got != expected {
		t.Errorf("MixColumns failed got:\n%s\nexpected\n%s\n", got, expected)
	}
}
