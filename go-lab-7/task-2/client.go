package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	// Подключаемся к серверу
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Ошибка подключения к серверу:", err)
		return
	}
	defer conn.Close()

	// Вводим сообщение
	fmt.Print("Введите сообщение: ")
	message, _ := bufio.NewReader(os.Stdin).ReadString('\n')

	// Отправляем сообщение серверу
	_, err = conn.Write([]byte(message))
	if err != nil {
		fmt.Println("Ошибка при отправке:", err)
		return
	}

	// Ожидаем ответ от сервера
	reply, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Println("Ошибка при получении ответа:", err)
		return
	}

	// Выводим ответ сервера
	fmt.Println("Ответ сервера:", reply)
}
