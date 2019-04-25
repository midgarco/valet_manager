package crypt

import (
	"bytes"
	"crypto/cipher"
	"crypto/des"
	"crypto/md5"
	"encoding/base64"
	"errors"

	"github.com/midgarco/valet_manager/config"
)

func getCipher(key string) (cipher.Block, error) {
	key16Buffer := []byte(key)
	md5Sum := md5.Sum(key16Buffer)

	var key24Buffer []byte
	key24Buffer = append(key24Buffer, md5Sum[:16]...)
	key24Buffer = append(key24Buffer, md5Sum[:8]...)

	return des.NewTripleDESCipher(key24Buffer)
}

func padBlocks(input []byte, blockSize int) []byte {
	length := blockSize - (len(input) % blockSize)
	padding := bytes.Repeat([]byte{byte(length)}, length)
	return append(input, padding...)
}

func unpadBlocks(input []byte) ([]byte, error) {
	length := len(input)
	padding := int(input[length-1])
	if (length - padding) < 0 {
		return input, errors.New("slice bounds out of range")
	}
	return input[:(length - padding)], nil
}

// Encrypt ...
func Encrypt(input string) (string, error) {
	key := config.Get("APP_KEY")
	if key == "" {
		return "", errors.New("application key not configured")
	}

	cipher, err := getCipher(key)
	if err != nil {
		return "", err
	}

	blockSize := cipher.BlockSize()
	b := padBlocks([]byte(input), blockSize)
	output := make([]byte, len(b))
	ob := output

	for len(b) > 0 {
		cipher.Encrypt(ob[:blockSize], b[:blockSize])
		b = b[blockSize:]
		ob = ob[blockSize:]
	}

	return base64.StdEncoding.EncodeToString(output), nil
}

// Decrypt ...
func Decrypt(input string) (string, error) {
	key := config.Get("APP_KEY")
	if key == "" {
		return "", errors.New("application key not configured")
	}

	cipher, err := getCipher(key)
	if err != nil {
		return "", err
	}

	blockSize := cipher.BlockSize()
	b, err := base64.StdEncoding.DecodeString(input)
	if err != nil {
		return "", err
	}
	output := make([]byte, len(b))
	ob := output

	for len(b) > 0 {
		cipher.Decrypt(ob[:blockSize], b[:blockSize])
		b = b[blockSize:]
		ob = ob[blockSize:]
	}

	output, err = unpadBlocks(output)
	if err != nil {
		return "", err
	}

	return string(output[:]), nil
}
