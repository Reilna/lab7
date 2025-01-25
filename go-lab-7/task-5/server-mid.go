package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// Middleware для логирования запросов
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("Начал %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r) // вызов следующего обработчика
		log.Printf("Закончил %s %s in %v", r.Method, r.URL.Path, time.Since(start))
	})
}

// Обработчик для GET-запроса /hello
func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Привет всем")
}

// Обработчик для POST-запроса /data
func dataHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// Чтение тела запроса
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Ошибка при чтении", http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		// Вывод данных в консоль
		// log.Printf("Начатый %s %s", r.Method, r.URL.Path)
		log.Printf("Полученные данные %s %s: %s\n", r.Method, r.URL.Path, body)

		// Ответ клиенту
		fmt.Fprintf(w, "Данные получены")
	} else {
		http.Error(w, "Этот метод не доступен", http.StatusMethodNotAllowed)
	}
}

func main() {
	mux := http.NewServeMux()

	// Регистрация маршрутов
	mux.HandleFunc("/hello", helloHandler)
	mux.HandleFunc("/data", dataHandler)

	// Обёртывание маршрутов в middleware для логирования
	loggedMux := loggingMiddleware(mux)

	fmt.Println("HTTP-сервер middleware запущен и внимает порту 8080...")
	log.Fatal(http.ListenAndServe("localhost:8080", loggedMux))
}
