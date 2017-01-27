package model

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
)

var commonIV = []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f}
var key_text = "astaxie12798akljzmknm.ahkjkljl;k"

func Encrypt(param []byte) []byte {
	c, err := initCipher()
	if err != nil {
		fmt.Print(err)
	}

	// Encrypted string
	cfb := cipher.NewCFBEncrypter(c, commonIV)
	ciphertext := make([]byte, len(param))
	cfb.XORKeyStream(ciphertext, param)

	return ciphertext
}

func Decrypt(param []byte) string {
	c, err := initCipher()
	if err != nil {
		fmt.Print(err)
	}

	// Decrypt strings
	cfbdec := cipher.NewCFBDecrypter(c, commonIV)
	plaintextCopy := make([]byte, len(param))
	cfbdec.XORKeyStream(plaintextCopy, param)

	s := string(plaintextCopy[:len(plaintextCopy)])

	return s
}

func initCipher() (cipher.Block, error) {
	// Create the aes encryption algorithm
	c, err := aes.NewCipher([]byte(key_text))
	if err != nil {
		fmt.Print(err)
	}

	return c, err
}
