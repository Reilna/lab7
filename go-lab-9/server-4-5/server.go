// server.go
package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4"
)

var db *pgx.Conn

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Age      int    `json:"age"`
	Password string `json:"password"`
}

var sessions = make(map[string]int)

func generateSessionToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func authorizeUser(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Неверный ввод", http.StatusBadRequest)
		return
	}

	var storedUser User
	err := db.QueryRow(context.Background(), "SELECT id, password FROM users WHERE name = $1", user.Name).Scan(&storedUser.ID, &storedUser.Password)
	if err != nil || storedUser.Password != user.Password {
		http.Error(w, "Неправильное имя пользователя или пароль", http.StatusUnauthorized)
		return
	}

	token, err := generateSessionToken()
	if err != nil {
		http.Error(w, "Ошибка при создании сессии", http.StatusInternalServerError)
		return
	}
	sessions[token] = storedUser.ID
	w.Header().Set("Authorization", token)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Авторизация успешна"})
}

func requireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if _, exists := sessions[token]; !exists {
			http.Error(w, "Неавторизован", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query(context.Background(), "SELECT id, name, age FROM users")
	if err != nil {
		http.Error(w, "Ошибка получения данных", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Name, &user.Age); err != nil {
			http.Error(w, "Ошибка чтения данных", http.StatusInternalServerError)
			return
		}
		users = append(users, user)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func getUserByID(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	var user User
	err := db.QueryRow(context.Background(), "SELECT id, name, age FROM users WHERE id = $1", id).Scan(&user.ID, &user.Name, &user.Age)
	if err != nil {
		http.Error(w, "Пользователь не найден", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func createUser(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Неверный ввод", http.StatusBadRequest)
		return
	}

	_, err := db.Exec(context.Background(), "INSERT INTO users (name, age, password) VALUES ($1, $2, $3)", user.Name, user.Age, user.Password)
	if err != nil {
		http.Error(w, "Ошибка при создании пользователя", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Пользователь создан"})
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Неверный ввод", http.StatusBadRequest)
		return
	}

	_, err := db.Exec(context.Background(), "UPDATE users SET name=$1, age=$2, password=$3 WHERE id=$4", user.Name, user.Age, user.Password, id)
	if err != nil {
		http.Error(w, "Ошибка при обновлении данных пользователя", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Пользователь обновлен"})
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	_, err := db.Exec(context.Background(), "DELETE FROM users WHERE id = $1", id)
	if err != nil {
		http.Error(w, "Ошибка при удалении пользователя", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Пользователь удален"})
}

func connectDB() {
	var err error
	connStr := "postgres://artematrr:@localhost:5432/lb8_users"
	db, err = pgx.Connect(context.Background(), connStr)
	if err != nil {
		log.Fatal("Не удалось установить соединение с базой данных:", err)
	}

	// Проверка существования администратора
	var exists bool
	err = db.QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM users WHERE id = 1)").Scan(&exists)
	if err != nil {
		log.Fatal("Ошибка при проверке существования администратора:", err)
	}

	if !exists {
		// Добавление администратора, если он не существует
		_, err = db.Exec(context.Background(), "INSERT INTO users (id, name, age, password) VALUES (1, 'admin', 30, 'admin')")
		if err != nil {
			log.Fatal("Ошибка при добавлении пользователя admin:", err)
		}
	}
}
func closeDB() {
	db.Close(context.Background())
}

func main() {
	connectDB()
	defer closeDB()

	r := mux.NewRouter()
	r.HandleFunc("/login", authorizeUser).Methods("POST")

	protected := r.PathPrefix("/").Subrouter()
	protected.Use(requireAuth)
	protected.HandleFunc("/users", getUsers).Methods("GET")
	protected.HandleFunc("/users/{id}", getUserByID).Methods("GET")
	protected.HandleFunc("/users", createUser).Methods("POST")
	protected.HandleFunc("/users/{id}", updateUser).Methods("PUT")
	protected.HandleFunc("/users/{id}", deleteUser).Methods("DELETE")

	fmt.Println("Сервер запущен на порту localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
