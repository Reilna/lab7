package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
)

func encrypt(plainText string, key string) (string, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	plainTextBytes := pad([]byte(plainText), block.BlockSize())
	cipherText := make([]byte, aes.BlockSize+len(plainTextBytes))

	iv := cipherText[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	cipher.NewCBCEncrypter(block, iv).CryptBlocks(cipherText[aes.BlockSize:], plainTextBytes)
	return base64.StdEncoding.EncodeToString(cipherText), nil
}

func decrypt(cipherText string, key string) (string, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	cipherTextBytes, _ := base64.StdEncoding.DecodeString(cipherText)
	if len(cipherTextBytes) < aes.BlockSize {
		return "", fmt.Errorf("ciphertext too short")
	}
	iv := cipherTextBytes[:aes.BlockSize]
	cipherTextBytes = cipherTextBytes[aes.BlockSize:]

	cbc := cipher.NewCBCDecrypter(block, iv)
	cbc.CryptBlocks(cipherTextBytes, cipherTextBytes)

	plainText := unpad(cipherTextBytes)
	return string(plainText), nil
}

func pad(src []byte, blockSize int) []byte {
	padding := blockSize - len(src)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(src, padtext...)
}

func unpad(src []byte) []byte {
	length := len(src)
	unpadding := int(src[length-1])
	return src[:(length - unpadding)]
}

func main() {
	var input string
	var key string

	fmt.Print("Введите строку для шифрования: ")
	fmt.Scanln(&input)
	fmt.Println("Введите секретный ключ:")
	fmt.Scanln(&key)

	encrypted, err := encrypt(input, key)
	if err != nil {
		fmt.Println("Ошибка шифрования:", err)
		return
	}
	fmt.Printf("Зашифрованный текст: %s\n", encrypted)

	fmt.Print("Введите зашифрованный текст для расшифрования: ")
	var cipherText string
	fmt.Scanln(&cipherText)
	decrypted, err := decrypt(cipherText, key)
	if err != nil {
		fmt.Println("Ошибка расшифрования:", err)
		return
	}
	fmt.Printf("Расшифрованный текст: %s\n", decrypted)
}
