package global

import "errors"

const (
	HashLength       = 32
	SaltLength       = 16
	Argon2Time       = 1
	Argon2Threads    = 4
	Argon2Memory     = 64 * 1024
	MagicBytesLength = 2

	OutFileSuffix = ".crypt"
)

var (
	MagicBytes           = [MagicBytesLength]byte{0xDC, 0x00}
	MagicBytesCompressed = [MagicBytesLength]byte{0xDC, 0x01}

	ErrNoMagicBytes  = errors.New("Magic bytes are wrong or missing")
	ErrWrongPassword = errors.New("Wrong password")
)
