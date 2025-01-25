package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"os"
)

func main() {

	conf := &tls.Config{
		InsecureSkipVerify: true, // проверка сертификата
	}

	conn, err := tls.Dial("tcp", "localhost:8080", conf)
	if err != nil {
		fmt.Println("Ошибка соединения:", err)
		os.Exit(1)
	}
	defer conn.Close()

	fmt.Println("Введите сообщение для отправки:")
	var message string
	fmt.Scanln(&message)

	_, err = conn.Write([]byte(message + "\n"))
	if err != nil {
		fmt.Println("Ошибка при отправке:", err)
		return
	}

	response, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Println("Ошибка при чтении ответа:", err)
		return
	}
	fmt.Println("Ответ сервера:", response)
}
