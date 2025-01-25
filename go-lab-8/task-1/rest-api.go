package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

// var users []User

var users = []User{
	{ID: 1, Name: "Andrew", Age: 24},
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	log.Println("GET /users - Получение данных обо всех пользователях")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])
	log.Printf("GET /users/%s - Получение данных о пользователe с ID: %d\n", params["id"], id)

	for _, user := range users {
		if user.ID == id {
			json.NewEncoder(w).Encode(user)
			return
		}
	}
	log.Printf("Пользователь с ID %d не найден\n", id)
	http.Error(w, "User not found", http.StatusNotFound)
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	log.Println("POST /users - Создание нового пользователя")
	var newUser User
	_ = json.NewDecoder(r.Body).Decode(&newUser)
	newUser.ID = len(users) + 1
	users = append(users, newUser)
	log.Printf("Создан новый пользователь: %+v\n", newUser)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newUser)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])
	log.Printf("PUT /users/%s - Обновление информации о пользователе с ID: %d\n", params["id"], id)

	for i, user := range users {
		if user.ID == id {
			_ = json.NewDecoder(r.Body).Decode(&user)
			users[i] = user
			log.Printf("Информация о пользователе обновлена: %+v\n", user)
			json.NewEncoder(w).Encode(user)
			return
		}
	}
	log.Printf("Пользователь с ID %d не найден\n", id)
	http.Error(w, "Пользователь не найден", http.StatusNotFound)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])
	log.Printf("DELETE /users/%s - Удаление пользователя с ID: %d\n", params["id"], id)

	for i, user := range users {
		if user.ID == id {
			users = append(users[:i], users[i+1:]...)
			log.Printf("Пользователь с ID %d удален\n", id)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
	log.Printf("Пользователь с ID %d не найден\n", id)
	http.Error(w, "Пользователь не найден", http.StatusNotFound)
}

// func loggingMiddleware(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		log.Printf("Request: %s %s", r.Method, r.RequestURI)
// 		next.ServeHTTP(w, r)
// 	})
// }

func main() {
	r := mux.NewRouter()
	// r.Use(loggingMiddleware)
	r.HandleFunc("/users", GetUsers).Methods("GET")
	r.HandleFunc("/users/{id}", GetUser).Methods("GET")
	r.HandleFunc("/users", CreateUser).Methods("POST")
	r.HandleFunc("/users/{id}", UpdateUser).Methods("PUT")
	r.HandleFunc("/users/{id}", DeleteUser).Methods("DELETE")

	log.Println("Сервер запущен на порту localhost:8080")
	http.ListenAndServe(":8080", r)
}
