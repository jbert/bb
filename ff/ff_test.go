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
	}
	for _, tc := range testCases {
		got := tc.p.mulWithHighBits(tc.q)
		if got != tc.expected {
			t.Errorf("p %s q %s: got %016b expected %016b", tc.p, tc.q, got, tc.expected)
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
	}
	for _, tc := range testCases {
		got := tc.p.Mul(tc.q)
		if got != tc.expected {
			t.Errorf("p %s q %s: got %016b expected %016b", tc.p, tc.q, got, tc.expected)
		}
	}
}
