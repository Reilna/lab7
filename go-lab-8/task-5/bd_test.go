package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateUser(t *testing.T) {
	connectDB()
	defer closeDB()

	user := User{Name: "Тестовый Пользователь", Age: 25}
	body, _ := json.Marshal(user)
	req, err := http.NewRequest("POST", "/users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(createUser)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("Неверный код ответа: получили %v, ожидали %v", status, http.StatusCreated)
	}

	var createdUser User
	if err := json.NewDecoder(rr.Body).Decode(&createdUser); err != nil {
		t.Errorf("Ошибка при декодировании ответа: %v", err)
	}

	if createdUser.ID == 0 {
		t.Error("Созданный пользователь не имеет ID")
	}
}

func TestGetUsers(t *testing.T) {
	connectDB()
	defer closeDB()

	req, err := http.NewRequest("GET", "/users", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(getUsers)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Неверный код ответа: получили %v, ожидали %v", status, http.StatusOK)
	}

	var users []User
	if err := json.NewDecoder(rr.Body).Decode(&users); err != nil {
		t.Errorf("Ошибка при декодировании ответа: %v", err)
	}

	if len(users) == 0 {
		t.Error("Ожидалось получить хотя бы одного пользователя, но список пуст")
	}
}

// func TestUpdateUser(t *testing.T) {
// 	connectDB()
// 	defer closeDB()

// 	// Создание тестового пользователя
// 	initialUser := User{Name: "Пользователь", Age: 20}
// 	createdUserID := createTestUser(t, initialUser.Name, initialUser.Age)

// 	// Обновление данных пользователя
// 	updatedUser := User{Name: "Обновленный Пользователь", Age: 30}
// 	body, _ := json.Marshal(updatedUser)
// 	req, err := http.NewRequest("PUT", "/users/"+strconv.Itoa(createdUserID), bytes.NewBuffer(body))
// 	req.Header.Set("Content-Type", "application/json")
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	rr := httptest.NewRecorder()
// 	handler := http.HandlerFunc(updateUser)
// 	handler.ServeHTTP(rr, req)

// 	if status := rr.Code; status != http.StatusOK {
// 		t.Errorf("Неверный код ответа: получили %v, ожидали %v", status, http.StatusOK)
// 	}

// 	// Проверка

// 	var user User
// 	err = json.NewDecoder(rr.Body).Decode(&user)
// 	if err != nil {
// 		t.Errorf("Ошибка при декодировании ответа: %v", err)
// 	}

// 	if user.Name != updatedUser.Name || user.Age != updatedUser.Age {
// 		t.Error("Данные пользователя не были обновлены корректно")
// 	}
// }

// func TestDeleteUser(t *testing.T) {
// 	connectDB()
// 	defer closeDB()

// 	// Создание тестового пользователя
// 	testUser := User{Name: "Удаляемый Пользователь", Age: 20}
// 	createdUserID := createTestUser(t, testUser.Name, testUser.Age)

// 	req, err := http.NewRequest("DELETE", "/users/"+strconv.Itoa(createdUserID), nil)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	rr := httptest.NewRecorder()
// 	handler := http.HandlerFunc(deleteUser)
// 	handler.ServeHTTP(rr, req)

// 	if status := rr.Code; status != http.StatusNoContent {
// 		t.Errorf("Неверный код ответа: получили %v, ожидали %v", status, http.StatusNoContent)
// 	}

// 	// Проверка, что пользователь действительно удален
// 	req, err = http.NewRequest("GET", "/users/"+strconv.Itoa(createdUserID), nil)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	rr = httptest.NewRecorder()
// 	handler = http.HandlerFunc(getUserByID)
// 	handler.ServeHTTP(rr, req)

// 	if status := rr.Code; status != http.StatusNotFound {
// 		t.Errorf("Пользователь не был удален, код ответа: %v", status)
// 	}
// }

// Утилита для создания тестового пользователя
func createTestUser(t *testing.T, name string, age int) int {
	user := User{Name: name, Age: age}
	body, _ := json.Marshal(user)
	req, err := http.NewRequest("POST", "/users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(createUser)
	handler.ServeHTTP(rr, req)

	var createdUser User
	if err := json.NewDecoder(rr.Body).Decode(&createdUser); err != nil {
		t.Errorf("Ошибка при декодировании ответа: %v", err)
	}

	return createdUser.ID
}
