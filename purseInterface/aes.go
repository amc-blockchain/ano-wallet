package purseInterface

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
)

const keyStr = ""

func AesOperate(operateStr string, operateDir bool) (reStr string, reErr error) {

	key := []byte(keyStr)
	if operateDir {

		crypted, err := aesEncrypt([]byte(operateStr), key)
		if err != nil {
			reErr = err
		}
		reStr = base64.StdEncoding.EncodeToString(crypted)
	} else {

		oBytes, err := base64.StdEncoding.DecodeString(operateStr)
		if err == nil {
			origData, err := aesDecrypt(oBytes, key)
			if err != nil {
				reErr = err
			}
			reStr = string(origData)
		}
	}
	return reStr, reErr
}

/*
 *
 * origData - ([]byte)
 * key -
 */
func aesEncrypt(origData, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	origData = pkcs5Padding(origData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

/*
 *
 * crypted - ([]byte)
 * key -
 */
func aesDecrypt(crypted, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData = pkcs5UnPadding(origData)
	return origData, nil
}

func pkcs5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func pkcs5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}
