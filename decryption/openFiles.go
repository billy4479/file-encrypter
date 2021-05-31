package decryption

import (
	"bufio"
	"bytes"
	"os"
	"path/filepath"
	"strings"

	g "github.com/billy4479/file-encrypter/global"
)

func openInput(data *decryptionData) (err error) {
	data.inputFile, err = os.Open(data.inputPath)
	if err != nil {
		return
	}

	data.inputFileBuffer = bufio.NewReader(data.inputFile)

	magicAndSalt := make([]byte, g.MagicBytesLength+g.SaltLength)
	_, err = data.inputFileBuffer.Read(magicAndSalt)
	if err != nil {
		return
	}

	magic := magicAndSalt[:g.MagicBytesLength]
	g.EPrintf("Magic bytes: 0x%x\n", magic)
	if bytes.Equal(magic, g.MagicBytesCompressed[:]) {
		data.isCompressed = true
	} else if !bytes.Equal(magic, g.MagicBytes[:]) {
		err = g.ErrNoMagicBytes
		return
	}

	data.salt = magicAndSalt[g.MagicBytesLength : g.MagicBytesLength+g.SaltLength]
	g.EPrintf("Salt: %x\n", data.salt)

	return
}

func openOutput(data *decryptionData) (err error) {
	if data.outputPath == "" {
		if strings.HasSuffix(data.inputFile.Name(), g.OutFileSuffix) {
			data.outputPath = filepath.Dir(data.inputPath)
			// data.outputPath = data.inputPath[:len(data.inputFile.Name())-len(g.OutFileSuffix)]
		}
	} else {
		err = os.MkdirAll(data.outputPath, 0755)
		if err != nil {
			return
		}
	}

	g.EPrintf("Output location: %s\n", data.outputPath)
	return
}
