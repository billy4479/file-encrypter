package encryption

import (
	g "github.com/billy4479/file-encrypter/global"
	"github.com/ulikunitz/xz/lzma"
)

func makeLZMA(data *encryptionData) (err error) {
	if !data.isCompressed {
		return
	}

	g.EPrintf("Enabling lzma compression.")

	data.lzmaWriter, err = lzma.NewWriter2(data.aesWriter)
	return
}
