package main

import (
	"archive/tar"
	"bufio"
	"crypto/rand"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/nicholastoddsmith/aesrw"
	"github.com/ulikunitz/xz/lzma"
	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/sha3"
)

func encrypt() {
	var outName string
	if *outLoc != "" {
		outName = *outLoc
	} else {
		outName = *target + ".crypt"
	}

	out, err := os.Create(outName)
	if err != nil {
		panic(err)
	}
	defer out.Close()
	buf := bufio.NewWriter(out)

	// Create hash and sale
	salt := make([]byte, 16)
	_, err = rand.Read(salt)
	if err != nil {
		panic(err)
	}
	key := argon2.IDKey([]byte(*pass), salt, 1, 64*1024, 4, 32)
	hash := sha3.Sum512(key)

	// Append hash and salt

	buf.Write(hash[:])
	buf.Write(salt)

	fmt.Printf("%x\n%x\n", hash, salt)

	if *compression {
		buf.Write([]byte{comp})
	} else {
		buf.Write([]byte{unComp})
	}

	// Encrypting
	aesw, err := aesrw.NewWriter(buf, key)
	if err != nil {
		panic(err)
	}
	defer aesw.Close()

	// Compressing
	if *compression {
		lw, err := lzma.NewWriter2(aesw)
		if err != nil {
			panic(err)
		}
		err = createTar(*target, lw)
		if err != nil {
			panic(err)
		}
		lw.Close()

	} else {
		err = createTar(*target, aesw)
		if err != nil {
			panic(err)
		}
	}

	log.Println("Done.")
}

func createTar(src string, buf io.Writer) (err error) {
	tarw := tar.NewWriter(buf)
	err = filepath.Walk(*target, func(path string, info os.FileInfo, inErr error) (err error) {
		if inErr != nil {
			return
		}
		thead, err := tar.FileInfoHeader(info, path)
		if err != nil {
			return
		}
		thead.Name = filepath.ToSlash(thead.Name)
		err = tarw.WriteHeader(thead)
		if err != nil {
			return
		}

		if info.IsDir() {
			return
		}

		fileData, err := os.Open(path)
		if err != nil {
			return
		}

		_, err = io.Copy(tarw, fileData)
		if err != nil {
			return
		}
		fileData.Close()

		return nil
	})
	if err != nil {
		return
	}
	err = tarw.Close()
	return
}
