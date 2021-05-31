# File Encrypter

This is a simple project to encrypt and decrypt files and folders.

## How it works

First it TARs the files, then, if the option is passed, it uses LZMA2 to compress the TAR archive.

Then everything is encrypted using AES256. The key is derived using Argon2ID.

At the beginning of the encrypted file are stored 2 magic bytes (to determine whether it was compressed or not) and the salt.

## Credits

Thanks to [this](https://medium.com/@skdomino/taring-untaring-files-in-go-6b07cf56bc07) article for tar and un-tar code.

Thanks to [nicholastoddsmith](https://github.com/nicholastoddsmith/) for [aesrw](https://github.com/nicholastoddsmith/aesrw), even if in the future I will implement my own solution.
