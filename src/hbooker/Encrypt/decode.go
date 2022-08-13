package Encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
)

// SHA256 sha256 编码
func SHA256(data []byte) []byte {
	ret := sha256.Sum256(data)
	return ret[:]
}

// AESDecrypt AES 解密
func AESDecrypt(EncryptKey string, contentText string) ([]byte, error) {
	if decoded, err := base64.StdEncoding.DecodeString(contentText); err == nil {
		if block, ok := aes.NewCipher(SHA256([]byte(EncryptKey))[:32]); ok == nil {
			blockModel, plainText := cipher.NewCBCDecrypter(block,
				[]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}), make([]byte, len(decoded))
			blockModel.CryptBlocks(plainText, decoded)
			return plainText[:(len(plainText) - int(plainText[len(plainText)-1]))], nil
		} else {
			fmt.Println("AESDecrypt Error:", ok)
			return nil, ok
		}
	} else {
		fmt.Println("Base64Decode Error:", err)
		return nil, err
	}
}

// Decode 入口函数
func Decode(content string, EncryptKey string) []byte {
	if EncryptKey == "" {
		EncryptKey = "zG2nSeEfSHfvTCHy5LCcqtBbQehKNLXn"
	}
	if raw, ok := AESDecrypt(EncryptKey, content); ok == nil {
		return raw
	} else {
		panic("Decrypt Error, Please Check Your Key!")
	}
}
