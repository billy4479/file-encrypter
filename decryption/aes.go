package decryption

import (
	g "github.com/billy4479/file-encrypter/global"
	"github.com/nicholastoddsmith/aesrw"
)

func deAES(data *decryptionData) (err error) {
	hash := g.MakeHash([]byte(data.password), data.salt)
	g.EPrintf("Key: %x\n", hash)
	data.aesReader, err = aesrw.NewReader(data.inputFileBuffer, hash)
	return
}
