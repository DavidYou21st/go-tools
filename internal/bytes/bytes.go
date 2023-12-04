//go:build go1.20

package bytes

import (
	"math/rand"
	"unsafe"
)

// RandomBytes generates len random bytes.
func RandomBytes(len int) ([]byte, error) {
	b := make([]byte, len)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// StringToBytes converts string to byte slice without a memory allocation.
func StringToBytes(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

// BytesToString converts byte slice to string without a memory allocation.
func BytesToString(b []byte) string {
	if len(b) == 0 {
		return ""
	}
	return unsafe.String(&b[0], len(b))
}
