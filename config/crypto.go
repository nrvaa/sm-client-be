package config

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"errors"
)

// pkcs7Unpad menghapus padding PKCS7 dari blok data
func pkcs7Unpad(data []byte) ([]byte, error) {
	length := len(data)
	if length == 0 {
		return nil, errors.New("empty data")
	}
	paddingLen := int(data[length-1])
	if paddingLen == 0 || paddingLen > length {
		return nil, errors.New("invalid padding")
	}
	return data[:length-paddingLen], nil
}

// DecryptAES mendekripsi ciphertext HEX menggunakan HEX_KEY dan IV_KEY
func DecryptAES(hexCt string, hexKey string, hexIV string) (string, error) {
	if hexCt == "" {
		return "", nil
	}

	ct, err := hex.DecodeString(hexCt)
	if err != nil {
		return "", err
	}

	key, err := hex.DecodeString(hexKey)
	if err != nil {
		return "", err
	}

	iv, err := hex.DecodeString(hexIV)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	if len(ct)%aes.BlockSize != 0 {
		return "", errors.New("ciphertext is not a multiple of the block size")
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(ct, ct)

	unpadded, err := pkcs7Unpad(ct)
	if err != nil {
		return "", err
	}

	return string(unpadded), nil
}
