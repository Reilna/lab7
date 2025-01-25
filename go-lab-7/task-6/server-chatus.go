package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

// Обновление соединений веб-сокетов
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// Хранилище подключений
var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan string)

// Обработчик веб-сокет соединений
func handleConnections(w http.ResponseWriter, r *http.Request) {
	// Обновление соединения до веб-сокета
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Ошибка при обновлении до веб-сокета:", err)
		return
	}
	defer ws.Close()

	// Добавляем клиента в карту
	clients[ws] = true

	for {
		// Чтение сообщения от клиента (ожидаем текст)
		_, msg, err := ws.ReadMessage()
		if err != nil {
			fmt.Println("Ошибка при чтении сообщения:", err)
			delete(clients, ws)
			break
		}
		// Отправляем сообщение в канал broadcast
		broadcast <- string(msg)
	}
}

// Рассылка сообщений всем подключённым клиентам
func handleMessages() {
	for {
		// Получаем сообщение из канала
		msg := <-broadcast
		fmt.Println(msg)

		// Отправляем сообщение всем клиентам
		for client := range clients {
			err := client.WriteMessage(websocket.TextMessage, []byte(msg))
			if err != nil {
				fmt.Println("Ошибка при отправке сообщения:", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}

func main() {
	// Запуск обработчика сообщений
	go handleMessages()

	// Запуск HTTP-сервера с веб-сокетами
	http.HandleFunc("/ws", handleConnections)

	fmt.Println("Сервер веб-сокетов запущен на :8080")
	err := http.ListenAndServe("localhost:8080", nil)
	if err != nil {
		fmt.Println("Ошибка при запуске сервера:", err)
	}
}
