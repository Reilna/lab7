package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	// Слушаем порт 8080
	ln, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Ошибка при запуске сервера:", err)
		os.Exit(1)
	}
	defer ln.Close()

	fmt.Println("TCP-сервер запущен и внимает порту 8080...")

	for {
		// Принимаем соединение
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Ошибка при принятии соединения:", err)
			continue
		}

		// Обрабатываем соединение в отдельной функции
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	// Читаем сообщение от клиента
	message, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Println("Ошибка при чтении:", err)
		return
	}

	fmt.Println("Полученное сообщение:", message)

	// Отправляем подтверждение
	_, err = conn.Write([]byte("Сообщение успешно доставлено\n"))
	if err != nil {
		fmt.Println("Ошибка при отправке:", err)
	}
}
