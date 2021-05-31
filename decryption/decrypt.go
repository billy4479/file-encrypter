package decryption

import (
	"archive/tar"
	"bufio"
	"os"

	"github.com/nicholastoddsmith/aesrw"
	"github.com/ulikunitz/xz/lzma"
)

type decryptionData struct {
	inputPath  string
	outputPath string

	inputFile       *os.File
	inputFileBuffer *bufio.Reader

	aesReader  *aesrw.AESReader
	lzmaReader *lzma.Reader2
	tarReader  *tar.Reader

	password string
	salt     []byte

	isCompressed bool
}

func (d *decryptionData) Close() (err error) {
	if d.inputFile != nil {
		err = d.inputFile.Close()
	}
	return
}

func Decrypt(inputPath string, outputPath string) (err error) {
	data := &decryptionData{
		inputPath:  inputPath,
		outputPath: outputPath,
	}
	defer data.Close()

	err = openInput(data)
	if err != nil {
		return
	}

	err = openOutput(data)
	if err != nil {
		return
	}

	err = deAES(data)
	if err != nil {
		return
	}

	err = unLZMA(data)
	if err != nil {
		return
	}

	err = unTar(data)
	return
}
