package util

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
)

var (
	aes_key = "#yadda@deepglint"
	aes_iv  = "^On1shiuva4$"
)

func Encrypt(origData []byte) ([]byte, error) {
	return aesEncrypt(origData, []byte(aes_key))
}

func Decrypt(crypted []byte) ([]byte, error) {
	return aesDecrypt(crypted, []byte(aes_key))
}

// 3DES加密
func aesEncrypt(origData, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	origData = pKCS5Padding(origData, block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, []byte(aes_iv)[:blockSize])
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

// 3DES解密
func aesDecrypt(crypted, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, []byte(aes_iv)[:blockSize])
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData = pKCS5UnPadding(origData)
	return origData, nil
}

func pKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func pKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}
