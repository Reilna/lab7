package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
)

var jwtKey = []byte("privet")

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
	Role string `json:"role"`
}

var users = []User{
	{ID: 1, Name: "Андрейка", Age: 24, Role: "user"},
	{ID: 2, Name: "Admin", Age: 30, Role: "admin"},
}

func generateJWT(user User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})
	return token.SignedString(jwtKey)
}

func authMiddleware(next http.Handler, allowedRoles ...string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "Токен отсутствует", http.StatusUnauthorized)
			return
		}

		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "Неверный токен", http.StatusUnauthorized)
			return
		}

		role, ok := claims["role"].(string)
		if !ok {
			http.Error(w, "Ошибка при чтении роли", http.StatusUnauthorized)
			return
		}

		// Проверяем, разрешена ли роль пользователя
		for _, allowedRole := range allowedRoles {
			if role == allowedRole {
				next.ServeHTTP(w, r)
				return
			}
		}
		http.Error(w, "Недостаточно прав", http.StatusForbidden)
	})
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var newUser User
	_ = json.NewDecoder(r.Body).Decode(&newUser)
	newUser.ID = len(users) + 1
	users = append(users, newUser)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newUser)
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		var u User
		json.NewDecoder(r.Body).Decode(&u)
		for _, user := range users {
			if user.Name == u.Name {
				token, err := generateJWT(user)
				if err != nil {
					http.Error(w, "Не удалось создать токен", http.StatusInternalServerError)
					return
				}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(map[string]string{"token": token})
				return
			}
		}
		http.Error(w, "Пользователь не найден", http.StatusUnauthorized)
	}).Methods("POST")

	r.Handle("/users", authMiddleware(http.HandlerFunc(GetUsers), "user", "admin")).Methods("GET")
	r.Handle("/users", authMiddleware(http.HandlerFunc(CreateUser), "admin")).Methods("POST")

	log.Println("Сервер запущен на порту localhost:8080")
	http.ListenAndServe(":8080", r)
}
