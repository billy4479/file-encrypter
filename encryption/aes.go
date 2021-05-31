package encryption

import (
	"crypto/rand"

	g "github.com/billy4479/file-encrypter/global"
	"github.com/nicholastoddsmith/aesrw"
)

func makeAES(data *encryptionData) (err error) {
	data.salt = make([]byte, g.SaltLength)
	_, err = rand.Reader.Read(data.salt)
	if err != nil {
		return
	}
	data.hash = g.MakeHash([]byte(data.password), data.salt)

	g.EPrintf("Salt: %x\n", data.salt)
	g.EPrintf("Key: %x\n", data.hash)

	_, err = data.outputBuffer.Write(data.salt)
	if err != nil {
		return
	}

	data.aesWriter, err = aesrw.NewWriter(data.outputBuffer, data.hash)
	return
}
