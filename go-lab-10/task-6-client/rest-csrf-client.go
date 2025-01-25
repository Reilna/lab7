package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

var token string
var csrfToken string

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
	Role string `json:"role"`
}

func login(name string) {
	loginURL := "http://localhost:8080/login"
	user := map[string]string{"name": name}
	jsonData, _ := json.Marshal(user)

	resp, err := http.Post(loginURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Ошибка при входе:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		var response map[string]string
		json.Unmarshal(body, &response)
		token = response["token"]
		csrfToken = response["csrf_token"]
		fmt.Printf("Вход выполнен. Токен: %s\nCSRF-токен: %s\n", token, csrfToken)
	} else {
		fmt.Println("Не удалось войти. Статус код:", resp.StatusCode)
	}
}

func getUsers() {
	userURL := "http://localhost:8080/users"
	req, _ := http.NewRequest("GET", userURL, nil)
	req.Header.Set("Authorization", token)
	req.Header.Set("X-CSRF-Token", csrfToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Ошибка при получении пользователей:", err)
		return
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode == http.StatusOK {
		var users []User
		json.Unmarshal(body, &users)

		fmt.Println("Список пользователей:")
		for _, user := range users {
			fmt.Printf("ID: %d, Имя: %s, Возраст: %d, Роль: %s\n", user.ID, user.Name, user.Age, user.Role)
		}
	} else {
		fmt.Println("Не удалось получить пользователей. Статус код:", resp.StatusCode)
	}
}

func createUser(name string, age int, role string) {
	userURL := "http://localhost:8080/users"
	newUser := map[string]interface{}{
		"name": name,
		"age":  age,
		"role": role,
	}
	jsonData, _ := json.Marshal(newUser)

	req, _ := http.NewRequest("POST", userURL, bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", token)
	req.Header.Set("X-CSRF-Token", csrfToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Ошибка при создании пользователя:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		var user User
		json.Unmarshal(body, &user)
		fmt.Printf("Пользователь создан: {'id': %d, 'name': '%s', 'age': %d, 'role': '%s'}\n", user.ID, user.Name, user.Age, user.Role)
	} else {
		fmt.Println("Не удалось создать пользователя. Нет прав доступа. Статус код:", resp.StatusCode)
	}
}

func main() {
	login("Admin")
	getUsers()
	createUser("Щеколда", 25, "user")
	getUsers()
}
