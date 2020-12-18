package main

import (
	"archive/tar"
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/nicholastoddsmith/aesrw"
	"github.com/ulikunitz/xz/lzma"
	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/sha3"
)

func decrypt() {
	in, err := os.Open(*target)
	if err != nil {
		panic(err)
	}
	var outName string
	if *outLoc != "" {
		outName = *outLoc
	} else {
		outName = strings.ReplaceAll(*target, ".crypt", "")
	}

	//Get hash, salt and mode
	hash := make([]byte, 64)
	salt := make([]byte, 16)
	mode := make([]byte, 1)

	buf := bufio.NewReader(in)

	buf.Read(hash)
	buf.Read(salt)
	buf.Read(mode)
	fmt.Printf("%x\n%x\n%x\n", hash, salt, mode)

	key := argon2.IDKey([]byte(*pass), salt, 1, 64*1024, 4, 32)
	sha := sha3.Sum512(key)
	newHash := sha[:]
	if comp := bytes.Compare(newHash, hash); comp != 0 {
		log.Println("Hash did not match")
		if !*force {
			os.Exit(1)
		}
	}

	// Spawn AES decrypter
	// Throwing EOF, cannot read IV -> ok
	tmp := bufio.NewReader(in)
	tmp.Discard(len(hash) + len(salt) + len(mode))
	aesr, err := aesrw.NewReader(tmp, key)
	if err != nil {
		panic(err)
	}

	modeEnum := mode[0]
	switch modeEnum {
	case comp:
		lzmar, err := lzma.NewReader2(aesr)
		if err != nil {
			panic(err)
		}
		err = extractTar(lzmar, outName)
		if err != nil {
			panic(err)
		}
	case unComp:
		err = extractTar(aesr, outName)
		if err != nil {
			panic(err)
		}
	default:
		log.Fatalln("Cannot determine compression type")
	}

	log.Println("Done")

}

func validRelPath(p string) bool {
	if p == "" || strings.Contains(p, `\`) || strings.HasPrefix(p, "/") || strings.Contains(p, "../") {
		return false
	}
	return true
}

func extractTar(src io.Reader, dst string) (err error) {
	tarr := tar.NewReader(src)
	for {
		thead, err := tarr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if !validRelPath(thead.Name) {
			return fmt.Errorf("Tar contains an invalid path")
		}

		out := filepath.Join(dst, thead.Name)

		switch thead.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(out, 0750); err != nil {
				return err
			}
		case tar.TypeReg:
			fileToWrite, err := os.OpenFile(out, os.O_CREATE|os.O_RDWR, os.FileMode(thead.Mode))
			if err != nil {
				return err
			}
			if _, err := io.Copy(fileToWrite, tarr); err != nil {
				return err
			}
			fileToWrite.Close()
		}
	}
	return
}
