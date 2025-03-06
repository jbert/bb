package bb

import (
	"encoding/hex"
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
			t.Errorf("got %08x expected %08x", got, tc.expected)
		}
	}
}

func TestExpandKey(t *testing.T) {
	key := MustKeyFromHex("2b7e151628aed2a6abf7158809cf4f3c")
	expected := [11]Key{
		MustKeyFromHex("2b7e151628aed2a6abf7158809cf4f3c"),
		MustKeyFromHex("a0fafe1788542cb123a339392a6c7605"),
		MustKeyFromHex("f2c295f27a96b9435935807a7359f67f"),
		MustKeyFromHex("3d80477d4716fe3e1e237e446d7a883b"),
		MustKeyFromHex("ef44a541a8525b7fb671253bdb0bad00"),
		MustKeyFromHex("d4d1c6f87c839d87caf2b8bc11f915bc"),
		MustKeyFromHex("6d88a37a110b3efddbf98641ca0093fd"),
		MustKeyFromHex("4e54f70e5f5fc9f384a64fb24ea6dc4f"),
		MustKeyFromHex("ead27321b58dbad2312bf5607f8d292f"),
		MustKeyFromHex("ac7766f319fadc2128d12941575c006e"),
		MustKeyFromHex("d014f9a8c9ee2589e13f0cc8b6630ca6"),
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

func TestAddRoundKey(t *testing.T) {
	in := State([]byte{
		0x6a, 0x2c, 0xb0, 0x27,
		0x6a, 0x6d, 0xd9, 0x9c,
		0x5c, 0x33, 0x5d, 0x21,
		0x45, 0x51, 0x61, 0x5c,
	})
	expected := State([]byte{
		0xbc, 0xfe, 0x6a, 0xf1,
		0xc0, 0xc2, 0x7f, 0x37,
		0x28, 0x41, 0x25, 0x57,
		0xb8, 0xab, 0x90, 0xa2,
	})
	k := MustKeyFromHex("d6aa74fdd2af72fadaa678f1d6ab76fe")

	got := in.AddRoundKey(k)

	if got != expected {
		t.Errorf("AddRoundKey failed got:\n%s\nexpected\n%s\n", got, expected)
	}
}

func TestComposite(t *testing.T) {
	in := State([]byte{
		0x00, 0x04, 0x08, 0x0c,
		0x01, 0x05, 0x09, 0x0d,
		0x02, 0x06, 0x0a, 0x0e,
		0x03, 0x07, 0x0b, 0x0f,
	})
	expected := State([]byte{
		0xbc, 0xfe, 0x6a, 0xf1,
		0xc0, 0xc2, 0x7f, 0x37,
		0x28, 0x41, 0x25, 0x57,
		0xb8, 0xab, 0x90, 0xa2,
	})
	k := MustKeyFromHex("d6aa74fdd2af72fadaa678f1d6ab76fe")

	got := in.SubBytes().ShiftRows().MixColumns().AddRoundKey(k)
	if got != expected {
		t.Errorf("Composite failed got:\n%s\nexpected\n%s\n", got, expected)
	}
}

func TestEncrypt(t *testing.T) {
	ptxt := "theblockbreakers"
	key := MustKeyFromHex("2b7e151628aed2a6abf7158809cf4f3c")
	got, err := EncryptBlock(ptxt, key)
	if err != nil {
		t.Fatalf("Can't encrypt block: %s", err)
	}

	expected := State([]byte{
		0xc6, 0x02, 0x23, 0x2f,
		0x9f, 0x5a, 0x93, 0x05,
		0x25, 0x9e, 0xf6, 0xb7,
		0xd0, 0xf3, 0x3e, 0x47,
	})
	if got != expected {
		t.Errorf("Encrypt failed got:\n%s\nexpected\n%s\n", got, expected)
	}

}

func TestAESVector(t *testing.T) {
	key, err := KeyFromHex("000102030405060708090a0b0c0d0e0f")
	if err != nil {
		t.Fatalf("Can't create key: %s", err)
	}
	ptxtBytes, err := hex.DecodeString("00112233445566778899aabbccddeeff")
	if err != nil {
		t.Fatalf("Can't decode plaintext hex: %s", err)
	}
	gotState, err := EncryptBlockBytes(ptxtBytes, key)
	if err != nil {
		t.Fatalf("Can't encrypt block: %s", err)
	}
	got := StateToBlock(gotState)

	expected := MustKeyFromHex("69c4e0d86a7b0430d8cdb78070b4c55a")
	if err != nil {
		t.Fatalf("Can't hex decode expected: %s", err)
	}
	if got != expected {
		t.Errorf("Encrypt failed got:\n%s\nexpected\n%s\n", got, expected)
	}

}
