// client.go
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

var sessionToken string

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Age      int    `json:"age"`
	Password string `json:"password"`
}

func login(username, password string) error {
	user := User{Name: username, Password: password}
	body, _ := json.Marshal(user)
	resp, err := http.Post("http://localhost:8080/login", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("ошибка подключения: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		sessionToken = resp.Header.Get("Authorization")
		fmt.Println("Авторизация успешна")
		return nil
	} else {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("ошибка авторизации: %s", string(body))
	}
}

func authorizedRequest(method, url string, data interface{}) (*http.Response, error) {
	var req *http.Request
	var err error

	if data != nil {
		body, _ := json.Marshal(data)
		req, err = http.NewRequest(method, url, bytes.NewBuffer(body))
	} else {
		req, err = http.NewRequest(method, url, nil)
	}
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", sessionToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	return client.Do(req)
}

func getUsers() {
	resp, err := authorizedRequest("GET", "http://localhost:8080/users", nil)
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var users []User
		json.NewDecoder(resp.Body).Decode(&users)
		fmt.Println("Пользователи:")
		for _, user := range users {
			fmt.Printf("ID: %d, Имя: %s, Возраст: %d\n", user.ID, user.Name, user.Age)
		}
	} else {
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println("Ошибка:", string(body))
	}
}

func getUserByID(id int) {
	resp, err := authorizedRequest("GET", fmt.Sprintf("http://localhost:8080/users/%d", id), nil)
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var user User
		json.NewDecoder(resp.Body).Decode(&user)
		fmt.Printf("Пользователь: ID: %d, Имя: %s, Возраст: %d\n", user.ID, user.Name, user.Age)
	} else {
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println("Ошибка:", string(body))
	}
}

func createUser(name string, age int, password string) {
	user := User{Name: name, Age: age, Password: password}
	resp, err := authorizedRequest("POST", "http://localhost:8080/users", user)
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusCreated {
		fmt.Println("Пользователь успешно создан")
	} else {
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println("Ошибка:", string(body))
	}
}

func updateUser(id int, name string, age int, password string) {
	user := User{ID: id, Name: name, Age: age, Password: password}
	resp, err := authorizedRequest("PUT", fmt.Sprintf("http://localhost:8080/users/%d", id), user)
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		fmt.Println("Пользователь успешно обновлен")
	} else {
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println("Ошибка:", string(body))
	}
}

func deleteUser(id int) {
	resp, err := authorizedRequest("DELETE", fmt.Sprintf("http://localhost:8080/users/%d", id), nil)
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		fmt.Println("Пользователь успешно удален")
	} else {
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println("Ошибка:", string(body))
	}
}

func main() {
	fmt.Println("Введите имя пользователя для входа:")
	var username string
	fmt.Scan(&username)

	fmt.Println("Введите пароль:")
	var password string
	fmt.Scan(&password)

	if err := login(username, password); err != nil {
		fmt.Println(err)
		return
	}

	for {
		fmt.Println("\nВыберите:")
		fmt.Println("\t1. Показать всех пользователей")
		fmt.Println("\t2. Найти пользователя по ID")
		fmt.Println("\t3. Создать нового пользователя")
		fmt.Println("\t4. Обновить данные пользователя")
		fmt.Println("\t5. Удалить пользователя")
		fmt.Println("\t0. Выход")
		fmt.Println("Действие:")

		var choice int
		fmt.Scan(&choice)

		switch choice {
		case 1:
			getUsers()
		case 2:
			fmt.Println("Введите ID пользователя:")
			var id int
			fmt.Scan(&id)
			getUserByID(id)
		case 3:
			fmt.Println("Введите имя:")
			var name string
			fmt.Scan(&name)
			fmt.Println("Введите возраст:")
			var age int
			fmt.Scan(&age)
			fmt.Println("Введите пароль:")
			var password string
			fmt.Scan(&password)
			createUser(name, age, password)
		case 4:
			fmt.Println("Введите ID пользователя:")
			var id int
			fmt.Scan(&id)
			fmt.Println("Введите новое имя:")
			var name string
			fmt.Scan(&name)
			fmt.Println("Введите новый возраст:")
			var age int
			fmt.Scan(&age)
			fmt.Println("Введите новый пароль:")
			var password string
			fmt.Scan(&password)
			updateUser(id, name, age, password)
		case 5:
			fmt.Println("Введите ID пользователя:")
			var id int
			fmt.Scan(&id)
			deleteUser(id)
		case 0:
			fmt.Println("Выход...")
			return
		default:
			fmt.Println("Неверный выбор, попробуйте снова.")
		}
	}
}
