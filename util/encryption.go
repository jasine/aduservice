package util

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
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
	return aesEncrypt(origData, []byte(aes_key))
}

func Decrypt(crypted []byte) ([]byte, error) {
	return aesDecrypt(crypted, []byte(aes_key))
}

// 加入md5校验合法性
func computeMd5(data []byte) []byte {
	md5Ctx := md5.New()
	md5Ctx.Write(data)
	md5bs := md5Ctx.Sum(nil)
	return md5bs
}

func aescheck(data []byte) bool {
	if len(data) >= len(aes_flag) && bytes.Equal(aes_flag, data[len(data)-len(aes_flag):]) {
		return true
	}
	return false
}

// AES加密
func aesEncrypt(origData, key []byte) ([]byte, error) {
	origData = append(origData, computeMd5(origData)...)
	mmmy(origData)
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

// AES解密
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

	if len(origData) < 16 || !bytes.Equal(origData[len(origData)-16:], computeMd5(origData[:len(origData)-16])) {
		return origData, errors.New("not valid")
	}
	return origData[:len(origData)-16], nil
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
