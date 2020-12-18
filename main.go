package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

var (
	target      = flag.String("target", "", "the target to operate, can be either a directory or a file")
	pass        = flag.String("password", "", "the password")
	outLoc      = flag.String("o", "", "the output location")
	decr        = flag.Bool("decrypt", false, "use it to decrypt a file already encrypted")
	compression = flag.Bool("compression", false, "enable compression (lzma2)")
	force       = flag.Bool("force", false, "when decrypting ingores the hash and extracts everything using the given password")
)

const (
	comp byte = iota
	unComp
)

func main() {
	flag.Parse()
	if *target == "" {
		log.Fatalln("Dir cannot be empty")
	} else {
		info, err := os.Stat(*target)
		if err != nil {
			panic(err)
		}
		if reg := info.Mode().IsRegular() || info.IsDir(); !reg {
			fmt.Println("The file is not regular, is it a device?")
			os.Exit(1)
		}
		if info.IsDir() && *decr {
			fmt.Println("Cannot use a directory as the decompression target")
			os.Exit(1)
		}
	}
	if *pass == "" {
		fmt.Println("Password cannot be empty")
		os.Exit(1)
	}
	if *decr {
		decrypt()
	} else {
		encrypt()
	}
}
