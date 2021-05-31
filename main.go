package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/billy4479/file-encrypter/decryption"
	"github.com/billy4479/file-encrypter/encryption"
)

func main() {
	decrypt := flag.Bool("d", false, "Decrypt switch. Default: false")
	outputPath := flag.String("o", "", "Output path. Default: the same as input but with '.crypt' at the end")
	password := flag.String("p", "", "Password. Required")
	inputFile := flag.String("i", "", "Input file. Required")
	useCompression := flag.Bool("c", false, "Use lzma2 compression. Default: false")

	flag.Parse()

	if *inputFile == "" {
		fmt.Println("Please enter an input file.")
		os.Exit(1)
	}

	if *password == "" {
		fmt.Println("Please enter a password.")
		os.Exit(1)
	}

	var err error
	if *decrypt {
		err = decryption.Decrypt(*inputFile, *outputPath)
	} else {
		err = encryption.Encrypt(*inputFile, *outputPath, *useCompression)
	}

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
