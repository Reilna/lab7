package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4"
)

var db *pgx.Conn

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func connectDB() {
	var err error
	connStr := "postgres://artematrr:@localhost:5432/lb8_users"
	db, err = pgx.Connect(context.Background(), connStr)
	if err != nil {
		log.Fatal("Не удалось установить соединение с базой данных:", err)
	}
}

func closeDB() {
	if db != nil {
		err := db.Close(context.Background())
		if err != nil {
			log.Fatal("Не удалось закрыть соединение:", err)
		}
	}
}

// Middleware для централизованной обработки ошибок
func errorHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				http.Error(w, "Произошла ошибка на сервере", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// Валидация полей User
func validateUser(user User) error {
	if strings.TrimSpace(user.Name) == "" {
		return errors.New("поле 'name' не может быть пустым")
	}
	if user.Age <= 0 {
		return errors.New("поле 'age' должно быть положительным числом")
	}
	return nil
}

// Обработчик для получения списка пользователей с пагинацией и фильтрацией
func getUsers(w http.ResponseWriter, r *http.Request) {
	// Извлекаем параметры пагинации и фильтрации из запроса
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")
	nameFilter := r.URL.Query().Get("name")
	ageFilterStr := r.URL.Query().Get("age")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 3 // Значения по умолчанию для пагинации
	}
	offset, _ := strconv.Atoi(offsetStr) // offset по умолчанию 0

	query := "SELECT id, name, age FROM users WHERE 1=1"
	args := []interface{}{}
	argIdx := 1

	// Добавляем фильтрацию по имени, если указан параметр name
	if nameFilter != "" {
		query += " AND name ILIKE '%' || $" + strconv.Itoa(argIdx) + " || '%'"
		args = append(args, nameFilter)
		argIdx++
	}

	// Добавляем фильтрацию по возрасту, если указан параметр age
	if ageFilterStr != "" {
		ageFilter, err := strconv.Atoi(ageFilterStr)
		if err == nil {
			query += " AND age = $" + strconv.Itoa(argIdx)
			args = append(args, ageFilter)
			argIdx++
		}
	}

	// Добавляем пагинацию
	query += " LIMIT $" + strconv.Itoa(argIdx) + " OFFSET $" + strconv.Itoa(argIdx+1)
	args = append(args, limit, offset)

	// Выполняем запрос к базе данных
	rows, err := db.Query(context.Background(), query, args...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Формируем список пользователей
	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Name, &user.Age); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		users = append(users, user)
	}

	// Возвращаем результат в формате JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// Создание нового пользователя с валидацией
func createUser(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Неверный ввод", http.StatusBadRequest)
		return
	}

	// Валидация данных пользователя
	if err := validateUser(user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	query := "INSERT INTO users (name, age) VALUES ($1, $2) RETURNING id"
	err := db.QueryRow(context.Background(), query, user.Name, user.Age).Scan(&user.ID)
	if err != nil {
		http.Error(w, "Ошибка при добавлении пользователя", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func main() {
	connectDB()
	defer closeDB()

	r := mux.NewRouter()
	r.Use(errorHandler)

	r.HandleFunc("/users", getUsers).Methods("GET")
	r.HandleFunc("/users", createUser).Methods("POST")

	fmt.Println("Сервер запущен на порту localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
