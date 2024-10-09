package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

const apiBaseURL = "http://localhost:8080/users"

type User struct {
	Id   string `json:"id,omitempty"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func main() {
	for {
		fmt.Println("================================")
		fmt.Println("Выберите операцию:")
		fmt.Println("1. Получить всех пользователей")
		fmt.Println("2. Получить пользователя по ID")
		fmt.Println("3. Создать пользователя")
		fmt.Println("4. Обновить пользователя")
		fmt.Println("5. Удалить пользователя")
		fmt.Println("6. Выход")
		fmt.Println("================================")

		var choice int
		fmt.Scan(&choice)

		switch choice {
		case 1:
			getAllUsers()
		case 2:
			getUser()
		case 3:
			createUser()
		case 4:
			updateUser()
		case 5:
			deleteUser()
		case 6:
			fmt.Println("Выход из программы.")
			os.Exit(0)
		default:
			fmt.Println("Неверный выбор, попробуйте снова.")
		}
	}
}

func getAllUsers() {
	resp, err := http.Get(apiBaseURL)
	if err != nil {
		fmt.Printf("Ошибка при получении пользователей: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var users []User
		body, _ := ioutil.ReadAll(resp.Body)
		err := json.Unmarshal(body, &users)
		if err != nil {
			fmt.Printf("Ошибка при разборе данных: %v\n", err)
			return
		}

		fmt.Println("Список пользователей:")
		for _, user := range users {
			fmt.Printf("ID: %s, Имя: %s, Возраст: %d\n", user.Id, user.Name, user.Age)
		}
	} else {
		fmt.Printf("Ошибка: %s\n", resp.Status)
	}
}

func createUser() {
	var user User
	fmt.Print("Введите имя пользователя: ")
	fmt.Scan(&user.Name)
	fmt.Print("Введите возраст пользователя: ")
	fmt.Scan(&user.Age)

	body, _ := json.Marshal(user)
	resp, err := http.Post(apiBaseURL, "application/json", bytes.NewBuffer(body))
	if err != nil {
		fmt.Printf("Ошибка при создании пользователя: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusCreated {
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Printf("Пользователь успешно создан: %s\n", body)
	} else {
		fmt.Printf("Ошибка: %s\n", resp.Status)
	}
}

func getUser() {
	var id string
	fmt.Print("Введите ID пользователя: ")
	fmt.Scan(&id)

	resp, err := http.Get(apiBaseURL + "/" + id)
	if err != nil {
		fmt.Printf("Ошибка при получении пользователя: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Printf("Данные пользователя: %s\n", body)
	} else {
		fmt.Printf("Ошибка: %s\n", resp.Status)
	}
}

func updateUser() {
	var user User
	fmt.Print("Введите ID пользователя: ")
	fmt.Scan(&user.Id)
	fmt.Print("Введите новое имя пользователя: ")
	fmt.Scan(&user.Name)
	fmt.Print("Введите новый возраст пользователя: ")
	fmt.Scan(&user.Age)

	body, _ := json.Marshal(user)
	req, err := http.NewRequest("PUT", apiBaseURL+"/"+user.Id, bytes.NewBuffer(body))
	if err != nil {
		fmt.Printf("Ошибка при обновлении пользователя: %v\n", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Ошибка при обновлении пользователя: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Printf("Пользователь успешно обновлен: %s\n", body)
	} else {
		fmt.Printf("Ошибка: %s\n", resp.Status)
	}
}

func deleteUser() {
	var id string
	fmt.Print("Введите ID пользователя для удаления: ")
	fmt.Scan(&id)

	req, err := http.NewRequest("DELETE", apiBaseURL+"/"+id, nil)
	if err != nil {
		fmt.Printf("Ошибка при удалении пользователя: %v\n", err)
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Ошибка при удалении пользователя: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		fmt.Println("Пользователь успешно удален.")
	} else {
		fmt.Printf("Ошибка: %s\n", resp.Status)
	}
}
