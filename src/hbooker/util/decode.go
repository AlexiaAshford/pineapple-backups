package util

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"
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
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, err
	}
	return decoded, nil
}

//LoadKey 读取解密密钥
func LoadKey(EncryptKey string) []byte {
	Key := SHA256([]byte(EncryptKey))
	return Key[:32]
}

//AESDecrypt AES 解密
func AESDecrypt(EncryptKey string, ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(LoadKey(EncryptKey))
	if err != nil {
		return nil, err
	}
	blockModel, plainText := cipher.NewCBCDecrypter(block, IV), make([]byte, len(ciphertext))
	blockModel.CryptBlocks(plainText, ciphertext)
	return plainText[:(len(plainText) - int(plainText[len(plainText)-1]))], nil
}

//Decode 入口函数
func Decode(str string, EncryptKey string) string {
	var err error
	var decoded, raw []byte
	decoded, err = Base64Decode(str)
	if err != nil {
		panic(err)
	}
	raw, err = AESDecrypt(EncryptKey, decoded)
	if err != nil {
		panic(err)
	}
	return string(raw)
}
