package global

import (
	"fmt"
	"os"

	"golang.org/x/crypto/argon2"
)

func MakeHash(plainText []byte, salt []byte) []byte {
	return argon2.IDKey([]byte(plainText), []byte(salt), Argon2Time, Argon2Memory, Argon2Threads, HashLength)
}

func EPrintf(format string, a ...interface{}) (int, error) {
	return fmt.Fprintf(os.Stderr, format, a...)
}
