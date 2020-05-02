// crypto for calculating password, key or anything else about crypt
package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
)

// prepend to encrypting data for validation check
var Header = "dong"

// CalcId calc the data id from r128 for data storing.
// r128 is the last 128 bytes of user password sha256
func CalcId(r128 string, idx int) string {
	digest := sha256.New()
	digest.Write([]byte(r128))
	digest.Write([]byte(fmt.Sprintf("%d", idx)))

	key := hex.EncodeToString(digest.Sum(nil))
	return key
}

// CalcSecret calc the secret for encrypting
// r128 is the last 128 bytes of user password sha256
func CalcSecret(r128 string, salt string) []byte {
	digest := sha256.New()
	digest.Write([]byte(r128))
	digest.Write([]byte(salt))

	return digest.Sum(nil)[:aes.BlockSize]
}

func padding(buf []byte, blockSize int) []byte {
	nPadding := blockSize - len(buf)%blockSize
	padding := bytes.Repeat([]byte{byte(nPadding)}, nPadding)

	return append(buf, padding...)
}

func unpadding(buf []byte) []byte {
	n := len(buf)
	nPadding := int(buf[n-1])

	return buf[:n-nPadding]
}

// Encrypt encrypt the plain bytes to secret bytes
func Encrypt(plain, key []byte) ([]byte, error) {
	plain = append([]byte(Header), plain...)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	plain = padding(plain, block.BlockSize())

	blockMode := cipher.NewCBCEncrypter(block, key)
	blockMode.CryptBlocks(plain, plain)

	return plain, nil
}

// Decrypt decrypt the secret bytes to plain bytes
func Decrypt(secret, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	plain := make([]byte, len(secret))
	blockMode := cipher.NewCBCDecrypter(block, key)
	blockMode.CryptBlocks(plain, secret)

	if bytes.Compare(plain[:len(Header)], []byte(Header)) != 0 {
		return nil, errors.New("Invalid data header!")
	}

	plain = unpadding(plain)

	return plain[len(Header):], nil
}

// func main() {
// 	plain := "hello"
// 	password := "password"
// 	salt := "saltish"
// 	secret := CalcSecret(password, salt)
// 	encrypted, _ := Encrypt([]byte(plain), secret)
// 	decrypted, _ := Decrypt(encrypted, secret)
// 	println(string(decrypted))
// }
