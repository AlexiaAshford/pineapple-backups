package HbookerAPI

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
)

var (
	//IV 偏移量
	IV = []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
)

//SHA256 sha256 编码
func SHA256(data []byte) []byte {
	ret := sha256.Sum256(data)
	return ret[:]
}

//Base64Decode Base64 解码
func Base64Decode(encoded string) ([]byte, error) {
	if decoded, err := base64.StdEncoding.DecodeString(encoded); err == nil {
		return decoded, nil
	} else {
		return nil, err
	}
}

//LoadKey 读取解密密钥
func LoadKey(EncryptKey string) []byte {
	Key := SHA256([]byte(EncryptKey))
	return Key[:32]
}

//AESDecrypt AES 解密
func AESDecrypt(EncryptKey string, ciphertext []byte) ([]byte, error) {
	if EncryptKey == "" {
		EncryptKey = "zG2nSeEfSHfvTCHy5LCcqtBbQehKNLXn"
	}
	if block, err := aes.NewCipher(LoadKey(EncryptKey)); err == nil {
		blockModel, plainText := cipher.NewCBCDecrypter(block, IV), make([]byte, len(ciphertext))
		blockModel.CryptBlocks(plainText, ciphertext)
		return plainText[:(len(plainText) - int(plainText[len(plainText)-1]))], nil
	} else {
		return nil, err
	}
}

//Decode 入口函数
func Decode(content string, EncryptKey string) []byte {
	if decoded, err := Base64Decode(content); err == nil {
		if raw, ok := AESDecrypt(EncryptKey, decoded); ok == nil {
			return raw
		} else {
			fmt.Println("AESDecrypt Error:", ok)
		}
	} else {
		fmt.Println("Base64Decode Error:", err)
	}
	panic("Decrypt Error, Please Check Your Key!")
}
