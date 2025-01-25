package main

import (
	"crypto/md5"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"hash"
)

func hashString(hashFunc string, input string) string {
	var h hash.Hash
	switch hashFunc {
	case "md5":
		h = md5.New()
	case "sha256":
		h = sha256.New()
	case "sha512":
		h = sha512.New()
	default:
		fmt.Println("Неизвестная функция хэширования")
		return ""
	}
	h.Write([]byte(input))
	return hex.EncodeToString(h.Sum(nil))
}

func main() {
	fmt.Println("Введите строку для хэширования:")
	var input string
	fmt.Scanln(&input)

	fmt.Println("Выберите хэш-функцию:")
	fmt.Println("1. MD5")
	fmt.Println("2. SHA-256")
	fmt.Println("3. SHA-512")
	var choice int
	fmt.Scanln(&choice)

	var hashFunc string
	switch choice {
	case 1:
		hashFunc = "md5"
	case 2:
		hashFunc = "sha256"
	case 3:
		hashFunc = "sha512"
	default:
		fmt.Println("Некорректный выбор")
		return
	}

	hashed := hashString(hashFunc, input)
	fmt.Printf("Хэш: %s\n", hashed)

	fmt.Print("Введите строку для проверки целостности: ")
	var checkInput string
	fmt.Scanln(&checkInput)
	fmt.Print("Введите хэш для проверки: ")
	var checkHash string
	fmt.Scanln(&checkHash)

	if hashString(hashFunc, checkInput) == checkHash {
		fmt.Println("Целостность данных подтверждена.")
	} else {
		fmt.Println("Целостность данных опровергнута.")
	}
}
