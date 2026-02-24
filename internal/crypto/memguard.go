package crypto

import "github.com/awnumar/memguard"

func InitMemguard() {
	memguard.CatchInterrupt()
}

func DestroyMemguard() {
	memguard.Purge()
}

func NewLockedKey(data []byte) (*memguard.LockedBuffer, error) {
	buf := memguard.NewBufferFromBytes(data)
	ZeroBytes(data)
	return buf, nil
}

func SealKey(buf *memguard.LockedBuffer) *memguard.Enclave {
	return buf.Seal()
}

func UnsealKey(enc *memguard.Enclave) (*memguard.LockedBuffer, error) {
	return enc.Open()
}

func ZeroBytes(b []byte) {
	memguard.WipeBytes(b)
}
