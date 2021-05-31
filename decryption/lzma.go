package decryption

import "github.com/ulikunitz/xz/lzma"

func unLZMA(data *decryptionData) (err error) {
	if !data.isCompressed {
		return
	}

	data.lzmaReader, err = lzma.NewReader2(data.aesReader)
	return
}
