package crypto

import (
	"bytes"
	"testing"
)

func TestLockedKey_SealUnseal(t *testing.T) {
	original := []byte("32-byte-key-for-testing-1234567!")
	cp := make([]byte, len(original))
	copy(cp, original)

	buf, err := NewLockedKey(cp)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(buf.Bytes(), original) {
		t.Fatal("locked buffer should contain original data")
	}

	enc := SealKey(buf)

	buf2, err := UnsealKey(enc)
	if err != nil {
		t.Fatal(err)
	}
	defer buf2.Destroy()

	if !bytes.Equal(buf2.Bytes(), original) {
		t.Fatal("unsealed key should match original")
	}
}

func TestZeroBytes(t *testing.T) {
	data := []byte("sensitive")
	ZeroBytes(data)
	for _, b := range data {
		if b != 0 {
			t.Fatal("expected all zeros after wipe")
		}
	}
}
