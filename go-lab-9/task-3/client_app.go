package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func displayMenu() {
	fmt.Println("\nВыберите:")
	fmt.Println("\t1 - Показать всех пользователей")
	fmt.Println("\t2 - Показать пользователя по ID")
	fmt.Println("\t3 - Добавить нового пользователя")
	fmt.Println("\t4 - Обновить информацию о пользователе")
	fmt.Println("\t5 - Удалить пользователя")
	fmt.Println("\t0 - Выйти")
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

	fmt.Println("Список пользователей:")
	for _, user := range users {
		fmt.Printf("ID: %d, Имя: %s, Возраст: %d\n", user.ID, user.Name, user.Age)
	}
}

func GetUserByID() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Введите ID пользователя: ")
	idStr, _ := reader.ReadString('\n')
	idStr = strings.TrimSpace(idStr)
	id, _ := strconv.Atoi(idStr)

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
		fmt.Printf("ID: %d, Имя: %s, Возраст: %d\n", user.ID, user.Name, user.Age)
	} else {
		fmt.Println("Пользователь не найден:", string(body))
	}
}

func CreateUser() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Введите имя пользователя: ")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)

	fmt.Print("Введите возраст пользователя: ")
	ageStr, _ := reader.ReadString('\n')
	ageStr = strings.TrimSpace(ageStr)
	age, _ := strconv.Atoi(ageStr)

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

func UpdateUser() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Введите ID пользователя для обновления: ")
	idStr, _ := reader.ReadString('\n')
	idStr = strings.TrimSpace(idStr)
	id, _ := strconv.Atoi(idStr)

	fmt.Print("Введите новое имя пользователя: ")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)

	fmt.Print("Введите новый возраст пользователя: ")
	ageStr, _ := reader.ReadString('\n')
	ageStr = strings.TrimSpace(ageStr)
	age, _ := strconv.Atoi(ageStr)

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

func DeleteUser() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Введите ID пользователя для удаления: ")
	idStr, _ := reader.ReadString('\n')
	idStr = strings.TrimSpace(idStr)
	id, _ := strconv.Atoi(idStr)

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
	for {
		displayMenu()
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Действие: ")
		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)

		switch choice {
		case "1":
			GetUsers()
		case "2":
			GetUserByID()
		case "3":
			CreateUser()
		case "4":
			UpdateUser()
		case "5":
			DeleteUser()
		case "0":
			fmt.Println("Завершение работы.")
			return
		default:
			fmt.Println("Неверный выбор, попробуйте еще раз.")
		}
	}
}
