package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

var baseURL = "http://localhost:8080"

type User struct {
	Name string `json:"name"`
}

func login(username string) (string, error) {
	user := User{Name: username}
	data, _ := json.Marshal(user)
	resp, err := http.Post(baseURL+"/login", "application/json", bytes.NewBuffer(data))
	if err != nil {
		return "", fmt.Errorf("ошибка при запросе: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("не удалось войти: %s", resp.Status)
	}

	var result map[string]string
	json.NewDecoder(resp.Body).Decode(&result)
	return result["token"], nil
}

func getUsers(token string) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", baseURL+"/users", nil)
	req.Header.Set("Authorization", token)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Ошибка при получении пользователей: %v", err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("Список пользователей:", string(body))
}

func createUser(token, name string, age int, role string) {
	client := &http.Client{}
	newUser := map[string]interface{}{"name": name, "age": age, "role": role}
	data, _ := json.Marshal(newUser)
	req, _ := http.NewRequest("POST", baseURL+"/users", bytes.NewBuffer(data))
	req.Header.Set("Authorization", token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Ошибка при добавлении пользователя: %v", err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("Результат добавления пользователя:", string(body))
}

func main() {
	fmt.Println()
	// token, err := login("Андрейка")
	token, err := login("Admin")

	if err != nil {
		log.Fatalf("Ошибка аутентификации: %v", err)
	}
	fmt.Println("Токен аутентификации:", token)
	fmt.Println()

	getUsers(token)
	createUser(token, "Типчик", 22, "user")
	getUsers(token)
}
