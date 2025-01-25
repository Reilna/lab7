package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// Обработчик для GET-запроса /hello
func helloHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		fmt.Fprintf(w, "Всем привет")
	} else {
		http.Error(w, "Этот метод не доступен", http.StatusMethodNotAllowed)
	}
}

// Обработчик для POST-запроса /data
func dataHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// Чтение тела запроса
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Ошиюка при чтении", http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		// Парсинг JSON
		var data map[string]interface{}
		if err := json.Unmarshal(body, &data); err != nil {
			http.Error(w, "Плохой JSON", http.StatusBadRequest)
			return
		}

		// Вывод данных в консоль
		fmt.Printf("Полученные данные: %v\n", data)

		// Ответ клиенту
		fmt.Fprintf(w, "Данные получены")
	} else {
		http.Error(w, "Этот метод не доступен", http.StatusMethodNotAllowed)
	}
}

func main() {
	http.HandleFunc("/hello", helloHandler)
	http.HandleFunc("/data", dataHandler)

	fmt.Println("HTTP сервер запущен и внимает порту 8080...")
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}
