package encryption

import (
	"archive/tar"
	"io"
	"os"
	"path/filepath"
	"strings"

	g "github.com/billy4479/file-encrypter/global"
)

func makeTar(data *encryptionData) error {
	// data.tarWriter = tar.NewWriter(io.MultiWriter(os.Stdout, data.aesWriter))

	if data.isCompressed {
		data.tarWriter = tar.NewWriter(data.lzmaWriter)
	} else {
		data.tarWriter = tar.NewWriter(data.aesWriter)
	}

	err := filepath.Walk(data.inputPath, func(file string, fi os.FileInfo, err error) error {

		if err != nil {
			return err
		}

		if (!fi.Mode().IsRegular() && !fi.IsDir()) || file == data.outputPath {
			return nil
		}

		header, err := tar.FileInfoHeader(fi, fi.Name())
		if err != nil {
			return err
		}
		header.Name = strings.TrimPrefix(strings.Replace(file, data.outputPath, "", -1), string(filepath.Separator))

		g.EPrintf("Encrypting %s...\n", header.Name)

		if err := data.tarWriter.WriteHeader(header); err != nil {
			return err
		}

		if fi.IsDir() {
			return nil
		}

		f, err := os.Open(file)
		if err != nil {
			return err
		}

		if _, err := io.Copy(data.tarWriter, f); err != nil {
			f.Close()
			return err
		}

		f.Close()

		return nil
	})

	if err != nil {
		return err
	}

	return data.tarWriter.Flush()
}
