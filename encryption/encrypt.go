package encryption

import (
	"archive/tar"
	"bufio"
	"os"

	"github.com/nicholastoddsmith/aesrw"
	"github.com/ulikunitz/xz/lzma"
)

type encryptionData struct {
	inputPath  string
	outputPath string

	outputFile   *os.File
	outputBuffer *bufio.Writer

	aesWriter  *aesrw.AESWriter
	lzmaWriter *lzma.Writer2
	tarWriter  *tar.Writer

	password string
	salt     []byte
	hash     []byte

	isCompressed bool
}

func (d *encryptionData) Close() (err error) {
	if d.tarWriter != nil {
		err = d.tarWriter.Close()
		if err != nil {
			return
		}
	}

	if d.lzmaWriter != nil {
		err = d.lzmaWriter.Close()
		if err != nil {
			return
		}
	}

	if d.aesWriter != nil {
		err = d.aesWriter.Close()
	}

	d.outputBuffer.Flush()

	return
}

func Encrypt(inputPath string, outputPath string, useCompression bool) (err error) {
	data := &encryptionData{
		inputPath:    inputPath,
		outputPath:   outputPath,
		isCompressed: useCompression,
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

	err = makeAES(data)
	if err != nil {
		return
	}

	err = makeLZMA(data)
	if err != nil {
		return
	}

	err = makeTar(data)
	return
}
