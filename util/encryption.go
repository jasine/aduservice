package util

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"errors"
)

//"github.com/deepglint/aduservice/util"
//Encrypt(origData []byte) ([]byte, error)
//Decrypt(crypted []byte) ([]byte, error)
var (
	aes_key  = "#yadda@deepglint"
	aes_iv   = "^On1shiuva4$"
	aes_flag = []byte("#muse&libra-t")
)

func Encrypt(origData []byte) ([]byte, error) {
	return aesEncrypt(append(origData, aes_flag...), []byte(aes_key))
}

func Decrypt(crypted []byte) ([]byte, error) {
	data, err := aesDecrypt(crypted, []byte(aes_key))
	if err != nil {
		return data, err
	}
	if aescheck(data) {
		return data[:len(data)-len(aes_flag)], nil
	} else {
		return data, errors.New("not valid")
	}
}
func aescheck(data []byte) bool {
	if len(data) >= len(aes_flag) && bytes.Equal(aes_flag, data[len(data)-len(aes_flag):]) {
		return true
	}
	return false
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
