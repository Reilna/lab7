package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func GetUsers() {
	resp, err := http.Get("http://localhost:8080/users")
	if err != nil {
		fmt.Println("Ошибка при отправке запроса:", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Ошибка при чтении ответа:", err)
		return
	}

	var users []User
	if err := json.Unmarshal(body, &users); err != nil {
		fmt.Println("Ошибка при обработке данных:", err)
		return
	}

	// Форматированный вывод пользователей
	fmt.Println("Список пользователей:")
	for _, user := range users {
		fmt.Printf("ID: %d, Имя: %s, Возраст: %d\n", user.ID, user.Name, user.Age)
	}
}

func GetUserByID(id int) {
	resp, err := http.Get(fmt.Sprintf("http://localhost:8080/users/%d", id))
	if err != nil {
		fmt.Println("Ошибка при отправке запроса:", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Ошибка при чтении ответа:", err)
		return
	}

	if resp.StatusCode == http.StatusOK {
		var user User
		if err := json.Unmarshal(body, &user); err != nil {
			fmt.Println("Ошибка при обработке данных:", err)
			return
		}
		fmt.Printf("Вывод пользователя - ID: %d, Имя: %s, Возраст: %d\n", user.ID, user.Name, user.Age)
	} else {
		fmt.Println("Ошибка при получении пользователя:", string(body))
	}
}

func CreateUser(name string, age int) {
	user := User{Name: name, Age: age}
	jsonData, _ := json.Marshal(user)

	resp, err := http.Post("http://localhost:8080/users", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Ошибка при отправке запроса:", err)
		return
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode == http.StatusCreated {
		var createdUser User
		json.Unmarshal(body, &createdUser)
		fmt.Printf("Пользователь создан с ID: %d\n", createdUser.ID)
	} else {
		fmt.Println("Ошибка при создании пользователя:", string(body))
	}
}

func CreateUserWithID(id int, name string, age int) {
	user := User{Name: name, Age: age}
	jsonData, _ := json.Marshal(user)

	url := fmt.Sprintf("http://localhost:8080/users/%d", id)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Ошибка при создании пользователя:", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Ошибка при отправке запроса:", err)
		return
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode == http.StatusCreated {
		var createdUser User
		json.Unmarshal(body, &createdUser)
		fmt.Printf("Пользователь создан с ID: %d\n", createdUser.ID)
	} else {
		fmt.Println("Ошибка при создании пользователя:", string(body))
	}
}

func UpdateUser(id int, name string, age int) {
	user := User{Name: name, Age: age}
	jsonData, _ := json.Marshal(user)

	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("http://localhost:8080/users/%d", id), bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Ошибка при создании запроса:", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Ошибка при отправке запроса:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		fmt.Println("Пользователь обновлен успешно.")
	} else {
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println("Ошибка при обновлении пользователя:", string(body))
	}
}

func DeleteUser(id int) {
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("http://localhost:8080/users/%d", id), nil)
	if err != nil {
		fmt.Println("Ошибка при создании запроса:", err)
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Ошибка при отправке запроса:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent {
		fmt.Println("Пользователь удален успешно.")
	} else {
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println("Ошибка при удалении пользователя:", string(body))
	}
}

func main() {
	CreateUser("Алексейчик", -1)
	CreateUserWithID(2, "", 25)
}
