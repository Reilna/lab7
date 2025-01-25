package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var wg sync.WaitGroup
var stopListening = make(chan struct{})

// Список активных соединений
var connections = make(map[net.Conn]struct{})
var connMutex sync.Mutex

func main() {
	ln, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Ошибка при запуске сервера:", err)
		os.Exit(1)
	}
	defer ln.Close()

	// Канал для обработки системных сигналов (для graceful shutdown)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	fmt.Println("TCP-сервер запущен и внимает порту 8080...")

	// Горутинa для graceful shutdown
	go func() {
		<-sigChan
		fmt.Println("\nПолучен сигнал завершения, сервер закрывается")
		close(stopListening)
		ln.Close()       // Закрываем сокет для новых соединений
		connMutex.Lock() // Закрываем все активные соединения
		for conn := range connections {
			conn.Close() // Закрываем каждое соединение
		}
		connMutex.Unlock()
		wg.Wait() // Ожидаем завершения всех горутин
		fmt.Println("Сервер безопасно завершён.")
		os.Exit(0)
	}()

	for {
		select {
		case <-stopListening:
			return
		default:
			// Принимаем соединения
			conn, err := ln.Accept()
			if err != nil {
				select {
				case <-stopListening:
					// Если ошибка вызвана закрытием сервера
					return
				default:
					fmt.Println("Ошибка при принятии соединения:", err)
					continue
				}
			}

			// Добавляем соединение в список активных
			connMutex.Lock()
			connections[conn] = struct{}{}
			connMutex.Unlock()

			// Добавляем новое соединение в группу ожидания
			wg.Add(1)
			go handleConnection(conn)
		}
	}
}

func handleConnection(conn net.Conn) {
	defer wg.Done()

	// Удаляем соединение из списка при завершении
	defer func() {
		connMutex.Lock()
		delete(connections, conn)
		connMutex.Unlock()
		conn.Close()
	}()

	// Читаем сообщение от клиента
	message, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Println("Ошибка при чтении:", err)
		return
	}

	fmt.Println("Сообщение от клиента:", message)

	// Имитация обработки для красоты
	time.Sleep(1 * time.Second)

	// Отправляем подтверждение клиенту
	_, err = conn.Write([]byte("Сообщение получено\n"))
	if err != nil {
		fmt.Println("Ошибка при отправке подтверждения:", err)
	}
}
