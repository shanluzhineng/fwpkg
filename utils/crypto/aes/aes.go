// Package aes provides aes encryption/decryption utilities
package aes

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"

	"github.com/shanluzhineng/fwpkg/utils/crypto"
)

// Encrypt string to base64 crypto using AES
func Encrypt(key []byte, text string) (retVal string, err error) {
	plaintext := []byte(text)
	var block cipher.Block
	block, err = aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the cipherText.
	cipherText := make([]byte, aes.BlockSize+len(plaintext))
	iv := cipherText[:aes.BlockSize]
	if _, err = io.ReadFull(rand.Reader, iv); err == nil {
		stream := cipher.NewCFBEncrypter(block, iv)
		stream.XORKeyStream(cipherText[aes.BlockSize:], plaintext)

		// convert to base64
		retVal = base64.URLEncoding.EncodeToString(cipherText)
	}
	return
}

// Decrypt from base64 to decrypted string
func Decrypt(key []byte, cryptoText string) (retVal string, err error) {
	cipherText, _ := base64.URLEncoding.DecodeString(cryptoText)

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the cipherText.
	if len(cipherText) < aes.BlockSize {
		return "", crypto.ErrCipherTooShort
	}

	var block cipher.Block
	block, err = aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)

	// XORKeyStream can work in-place if the two arguments are the same.
	stream.XORKeyStream(cipherText, cipherText)

	return fmt.Sprintf("%s", cipherText), nil
}
