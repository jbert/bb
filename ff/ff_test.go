package ff

import "testing"

func TestAdd(t *testing.T) {
	testCases := []struct {
		p        Poly
		q        Poly
		expected Poly
	}{
		{0x00, 0x00, 0x00},
		{0x01, 0x00, 0x01},
		{0x01, 0x01, 0x00},
		{0xff, 0x01, 0xfe},
	}
	for _, tc := range testCases {
		got := tc.p.Add(tc.q)
		if got != tc.expected {
			t.Errorf("p %s q %s: got %s expected %s", tc.p, tc.q, got, tc.expected)
		}
	}
}

func TestMulWithHighBits(t *testing.T) {
	testCases := []struct {
		p        Poly
		q        Poly
		expected uint16
	}{
		{0b00000000, 0b00000000, 0b00000000},
		{0b00000001, 0b00000000, 0b00000000},
		{0b00000000, 0b00000001, 0b00000000},
		{0b00000001, 0b00000001, 0b00000001},
		{0b00000110, 0b00010001, 0b01100110},
		{0b10000000, 0b00000010, 0b100000000},
		{0b10000000, 0b00000011, 0b110000000},
		{0b00000010, 0b10000001, 0b100000010},
	}
	for _, tc := range testCases {
		got := tc.p.mulWithHighBits(tc.q)
		if got != tc.expected {
			t.Errorf("p %s q %s: got %016b expected %016b", tc.p, tc.q, got, tc.expected)
		}
	}
}

func TestNewFromUint16(t *testing.T) {
	testCases := []struct {
		n        uint16
		expected Poly
	}{
		{0b0000000000000000, 0b00000000},
		{0b0000000000000001, 0b00000001},
		{0b0000000000000010, 0b00000010},
		{0b0000000010000000, 0b10000000},
		{0b0000000100000000, 0b00011011},
		{0b0000001000000000, 0b00110110},
		{0b0000010000000000, 0b01101100},
		{0b0000100000000000, 0b11011000},
		// X12
		// = (X4 + X3 + X + 1) X4
		// = X8 + X7 + X5 + X4
		// = (X4 + X3 + X + 1) + X7 + X5 + X4
		// = X7 + X5 + X3 + X + 1
		//
		//   00000001 10110000
		// = 00000000 10110000 + 00011011
		// = 00000000 10101011
		// = 00000000 10101011
		{0b0001000000000000, 0b10101011},
	}
	for _, tc := range testCases {
		t.Logf("n %016b", tc.n)
		got := NewFromUint16(tc.n)
		if got != tc.expected {
			t.Errorf("n %016b: got %08b expected %08b", tc.n, got, tc.expected)
		}
	}
}

func TestMul(t *testing.T) {
	testCases := []struct {
		p        Poly
		q        Poly
		expected Poly
	}{
		{0b00000000, 0b00000000, 0b00000000},
		{0b00000001, 0b00000000, 0b00000000},
		{0b00000000, 0b00000001, 0b00000000},
		{0b00000001, 0b00000001, 0b00000001},
		{0b00000110, 0b00010001, 0b01100110},
		{0b10000000, 0b00000010, 0b00011011},
		{0x02, 0x81, 0x19},
		{0x03, 0x8f, 0x8a},
		{0x03, 0xf0, 0x0b},
	}
	for _, tc := range testCases {
		got := tc.p.Mul(tc.q)
		if got != tc.expected {
			t.Errorf("p %s q %s: got %08b expected %08b", tc.p, tc.q, got, tc.expected)
		}
	}
}
