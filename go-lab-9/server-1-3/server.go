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
	connStr := "postgres://artematrr:@localhost:5432/lb8_users_test"
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

func validateUser(user User) error {
	if strings.TrimSpace(user.Name) == "" {
		return errors.New("поле 'name' не может быть пустым")
	}
	if user.Age <= 0 {
		return errors.New("поле 'age' должно быть положительным числом")
	}
	return nil
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")
	nameFilter := r.URL.Query().Get("name")
	ageFilterStr := r.URL.Query().Get("age")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 50 // Значения по умолчанию для пагинации
	}
	offset, _ := strconv.Atoi(offsetStr)

	query := "SELECT id, name, age FROM users WHERE 1=1"
	args := []interface{}{}
	argIdx := 1

	if nameFilter != "" {
		query += " AND name ILIKE '%' || $" + strconv.Itoa(argIdx) + " || '%'"
		args = append(args, nameFilter)
		argIdx++
	}

	if ageFilterStr != "" {
		ageFilter, err := strconv.Atoi(ageFilterStr)
		if err == nil {
			query += " AND age = $" + strconv.Itoa(argIdx)
			args = append(args, ageFilter)
			argIdx++
		}
	}

	query += " LIMIT $" + strconv.Itoa(argIdx) + " OFFSET $" + strconv.Itoa(argIdx+1)
	args = append(args, limit, offset)

	rows, err := db.Query(context.Background(), query, args...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Name, &user.Age); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		users = append(users, user)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func getUserByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	query := "SELECT id, name, age FROM users WHERE id = $1"
	row := db.QueryRow(context.Background(), query, id)

	var user User
	err := row.Scan(&user.ID, &user.Name, &user.Age)
	if err != nil {
		if err == pgx.ErrNoRows {
			http.Error(w, "Пользователь не найден", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
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

func createUserWithID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Неверный ID", http.StatusBadRequest)
		return
	}

	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Неверный ввод", http.StatusBadRequest)
		return
	}

	if err := validateUser(user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Проверяем, существует ли пользователь с таким ID
	var existingUser User
	query := "SELECT id FROM users WHERE id = $1"
	err = db.QueryRow(context.Background(), query, id).Scan(&existingUser.ID)
	if err == nil {
		http.Error(w, "Пользователь с таким ID уже существует", http.StatusConflict)
		return
	}

	// Создаём пользователя с указанным ID
	query = "INSERT INTO users (id, name, age) VALUES ($1, $2, $3)"
	_, err = db.Exec(context.Background(), query, id, user.Name, user.Age)
	if err != nil {
		http.Error(w, "Ошибка при добавлении пользователя", http.StatusInternalServerError)
		return
	}

	user.ID = id
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Неверный ввод", http.StatusBadRequest)
		return
	}

	if user.Name == "" || user.Age <= 0 {
		http.Error(w, "Поля \"name\" и \"age\" не могут быть пустыми", http.StatusBadRequest)
		return
	}

	query := "UPDATE users SET name = $1, age = $2 WHERE id = $3"
	_, err := db.Exec(context.Background(), query, user.Name, user.Age, id)
	if err != nil {
		if err == pgx.ErrNoRows {
			http.Error(w, "Пользователь не найден", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	query := "DELETE FROM users WHERE id = $1"
	_, err := db.Exec(context.Background(), query, id)
	if err != nil {
		if err == pgx.ErrNoRows {
			http.Error(w, "Пользователь не найден", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func main() {
	connectDB()
	defer closeDB()

	r := mux.NewRouter()
	r.Use(errorHandler)

	r.HandleFunc("/users", getUsers).Methods("GET")
	r.HandleFunc("/users/{id}", getUserByID).Methods("GET")
	r.HandleFunc("/users", createUser).Methods("POST")
	r.HandleFunc("/users/{id}", createUserWithID).Methods("POST")
	r.HandleFunc("/users/{id}", updateUser).Methods("PUT")
	r.HandleFunc("/users/{id}", deleteUser).Methods("DELETE")

	fmt.Println("Сервер запущен на порту localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
