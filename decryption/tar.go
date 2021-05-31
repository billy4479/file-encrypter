package decryption

import (
	"archive/tar"
	"io"
	"os"
	"path/filepath"

	g "github.com/billy4479/file-encrypter/global"
)

func unTar(data *decryptionData) error {
	if data.isCompressed {
		data.tarReader = tar.NewReader(data.lzmaReader)
	} else {
		data.tarReader = tar.NewReader(data.aesReader)
	}

	for {
		header, err := data.tarReader.Next()

		switch {
		case err == io.EOF:
			return nil

		case err != nil:
			return err

		case header == nil:
			continue
		}

		target := filepath.Join(data.outputPath, header.Name)
		g.EPrintf("Decrypting %s...\n", target)

		info := header.FileInfo()

		if info.IsDir() {
			if err := os.MkdirAll(target, 0755); err != nil {
				return err
			}
			continue
		}

		f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
		if err != nil {
			return err
		}

		if _, err := io.Copy(f, data.tarReader); err != nil {
			f.Close()
			return err
		}

		f.Close()
	}
}
