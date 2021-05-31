package encryption

import (
	"bufio"
	"os"
	"path/filepath"

	g "github.com/billy4479/file-encrypter/global"
)

func openOutput(data *encryptionData) (err error) {

	if data.outputPath == "" {
		data.outputPath = data.inputPath + g.OutFileSuffix
	}

	var stat os.FileInfo
	stat, err = os.Stat(data.outputPath)
	if err == nil {
		if stat.IsDir() {
			data.outputPath = filepath.Join(data.outputPath, filepath.Base(data.inputPath)+g.OutFileSuffix)
		}
	}

	data.outputFile, err = os.Create(data.outputPath)
	if err != nil {
		return
	}

	data.outputBuffer = bufio.NewWriter(data.outputFile)

	if data.isCompressed {
		_, err = data.outputBuffer.Write(g.MagicBytesCompressed[:])
	} else {
		_, err = data.outputBuffer.Write(g.MagicBytes[:])
	}
	if err != nil {
		return
	}

	return
}

func openInput(data *encryptionData) (err error) {
	_, err = os.Stat(data.inputPath)
	return
}
