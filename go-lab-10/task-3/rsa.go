package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"
)

func generateKeys() {
	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		fmt.Println("Ошибка генерации ключа:", err)
		return
	}

	privFile, err := os.Create("private_key.pem")
	if err != nil {
		fmt.Println("Ошибка создания файла:", err)
		return
	}
	defer privFile.Close()

	privBytes := x509.MarshalPKCS1PrivateKey(privKey)
	privPem := &pem.Block{Type: "PRIVATE KEY", Bytes: privBytes}
	if err := pem.Encode(privFile, privPem); err != nil {
		fmt.Println("Ошибка записи ключа в файл:", err)
		return
	}

	pubKey := &privKey.PublicKey
	pubFile, err := os.Create("public_key.pem")
	if err != nil {
		fmt.Println("Ошибка создания файла:", err)
		return
	}
	defer pubFile.Close()

	pubBytes := x509.MarshalPKCS1PublicKey(pubKey)
	pubPem := &pem.Block{Type: "PUBLIC KEY", Bytes: pubBytes}
	if err := pem.Encode(pubFile, pubPem); err != nil {
		fmt.Println("Ошибка записи ключа в файл:", err)
		return
	}

	fmt.Println("Ключи успешно сгенерированы и сохранены в файлы.")
}

func signMessage(message string) {
	privKeyBytes, err := ioutil.ReadFile("private_key.pem")
	if err != nil {
		fmt.Println("Ошибка чтения приватного ключа:", err)
		return
	}

	privPem, _ := pem.Decode(privKeyBytes)
	privKey, err := x509.ParsePKCS1PrivateKey(privPem.Bytes)
	if err != nil {
		fmt.Println("Ошибка парсинга приватного ключа:", err)
		return
	}

	signature, err := rsa.SignPKCS1v15(rand.Reader, privKey, 0, []byte(message))
	if err != nil {
		fmt.Println("Ошибка подписывания сообщения:", err)
		return
	}

	err = ioutil.WriteFile("signature.sig", signature, 0644)
	if err != nil {
		fmt.Println("Ошибка записи подписи в файл:", err)
		return
	}

	fmt.Println("Сообщение успешно подписано и подпись сохранена в файл signature.sig.")
}

func verifySignature(message string) {
	signature, err := ioutil.ReadFile("signature.sig")
	if err != nil {
		fmt.Println("Ошибка чтения подписи:", err)
		return
	}

	pubKeyBytes, err := ioutil.ReadFile("public_key.pem")
	if err != nil {
		fmt.Println("Ошибка чтения публичного ключа:", err)
		return
	}

	pubPem, _ := pem.Decode(pubKeyBytes)
	pubKey, err := x509.ParsePKCS1PublicKey(pubPem.Bytes)
	if err != nil {
		fmt.Println("Ошибка парсинга публичного ключа:", err)
		return
	}

	err = rsa.VerifyPKCS1v15(pubKey, 0, []byte(message), signature)
	if err != nil {
		fmt.Println("Подпись недействительна.")
	} else {
		fmt.Println("Подпись действительна.")
	}
}

func main() {
	fmt.Println("Выберите действие:")
	fmt.Println("1. Сгенерировать ключи")
	fmt.Println("2. Подписать сообщение")
	fmt.Println("3. Проверить подпись")
	var choice int
	fmt.Scan(&choice)

	switch choice {
	case 1:
		generateKeys()
	case 2:
		fmt.Println("Введите сообщение для подписи:")
		var message string
		fmt.Scan(&message)
		signMessage(message)
	case 3:
		fmt.Println("Введите сообщение для проверки подписи:")
		var message string
		fmt.Scan(&message)
		verifySignature(message)
	default:
		fmt.Println("Некорректный выбор.")
	}
}
