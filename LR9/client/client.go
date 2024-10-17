package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
)

const apiBaseURL = "http://localhost:8080/users"

type User struct {
	Id   string `json:"id,omitempty"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

var sessionToken string
var wg sync.WaitGroup

func main() {
	for {
		if sessionToken == "" {
			fmt.Println("================================")
			fmt.Println("Выберите операцию:")
			fmt.Println("1. Авторизоваться")
			fmt.Println("2. Зарегистрироваться")
			fmt.Println("3. Выход")
			fmt.Println("================================")

			var choice int
			fmt.Print("Введите номер операции: ")
			fmt.Scan(&choice)

			wg.Add(1)
			switch choice {
			case 1:
				go login()
			case 2:
				go register()
			case 3:
				fmt.Println("Выход из программы.")
				wg.Done()
				os.Exit(0)
			default:
				fmt.Println("Неверный выбор, попробуйте снова.")
				wg.Done()
			}
			wg.Wait()
		} else {
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
			fmt.Print("Введите номер операции: ")
			fmt.Scan(&choice)

			wg.Add(1)
			switch choice {
			case 1:
				go getAllUsers()
			case 2:
				go getUser()
			case 3:
				go createUser()
			case 4:
				go updateUser()
			case 5:
				go deleteUser()
			case 6:
				fmt.Println("Выход из программы.")
				wg.Done()
				os.Exit(0)
			default:
				fmt.Println("Неверный выбор, попробуйте снова.")
				wg.Done()
			}
			wg.Wait()
		}
	}
}

func login() {
	defer wg.Done()

	for {
		var login, password string
		fmt.Print("Введите логин: ")
		fmt.Scan(&login)
		fmt.Print("Введите пароль: ")
		fmt.Scan(&password)

		credentials := map[string]string{"login": login, "password": password}
		body, _ := json.Marshal(credentials)

		resp, err := http.Post(apiBaseURL+"/login", "application/json", bytes.NewBuffer(body))
		if err != nil {
			fmt.Printf("Ошибка при авторизации: %v\n", err)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			var result map[string]interface{}
			body, _ := ioutil.ReadAll(resp.Body)
			json.Unmarshal(body, &result)

			sessionToken = result["token"].(string)
			fmt.Println("Авторизация успешна, токен сохранён.")
			break
		} else {
			handleErrorResponse(resp)
			fmt.Println("Попробуйте снова.")
		}
	}
}

func register() {
	defer wg.Done()

	var login, password string
	fmt.Print("Введите логин: ")
	fmt.Scan(&login)
	fmt.Print("Введите пароль: ")
	fmt.Scan(&password)

	credentials := map[string]string{"login": login, "password": password}
	body, _ := json.Marshal(credentials)

	resp, err := http.Post(apiBaseURL+"/register", "application/json", bytes.NewBuffer(body))
	if err != nil {
		fmt.Printf("Ошибка при регистрации: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusCreated {
		fmt.Println("Регистрация прошла успешно.")
	} else {
		handleErrorResponse(resp)
	}
}

func getAllUsers() {
	defer wg.Done()
	req, err := http.NewRequest("GET", apiBaseURL, nil)
	if err != nil {
		fmt.Printf("Ошибка при создании запроса: %v\n", err)
		return
	}
	req.Header.Set("Authorization", "Bearer "+sessionToken)

	client := &http.Client{}
	resp, err := client.Do(req)
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
		handleErrorResponse(resp)
	}
}

func createUser() {
	defer wg.Done()
	var user User
	fmt.Print("Введите имя пользователя: ")
	fmt.Scan(&user.Name)

	for {
		fmt.Print("Введите возраст пользователя: ")
		_, err := fmt.Scan(&user.Age)
		if err != nil {
			fmt.Println("Ошибка: введите корректное целое число для возраста.")
			fmt.Scanln()
			continue
		}

		if user.Age <= 0 {
			fmt.Println("Возраст должен быть положительным числом.")
		} else {
			break
		}
	}

	body, _ := json.Marshal(user)
	req, err := http.NewRequest("POST", apiBaseURL, bytes.NewBuffer(body))
	if err != nil {
		fmt.Printf("Ошибка при создании пользователя: %v\n", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+sessionToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Ошибка при создании пользователя: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusCreated {
		var result map[string]interface{}
		body, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal(body, &result)

		createdUser := result["user"].(map[string]interface{})
		fmt.Println("Пользователь успешно создан:")
		fmt.Printf("ID: %s\nИмя: %s\nВозраст: %d\n", createdUser["id"], createdUser["name"], int(createdUser["age"].(float64)))
	} else {
		handleErrorResponse(resp)
	}
}

func updateUser() {
	defer wg.Done()
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
	req.Header.Set("Authorization", "Bearer "+sessionToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Ошибка при обновлении пользователя: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var result map[string]interface{}
		body, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal(body, &result)

		updatedUser := result["user"].(map[string]interface{})
		fmt.Println("Пользователь успешно обновлен:")
		fmt.Printf("ID: %s\nИмя: %s\nВозраст: %d\n", updatedUser["id"], updatedUser["name"], int(updatedUser["age"].(float64)))
	} else {
		handleErrorResponse(resp)
	}
}

func getUser() {
	defer wg.Done()
	var id string
	fmt.Print("Введите ID пользователя: ")
	fmt.Scan(&id)

	req, err := http.NewRequest("GET", apiBaseURL+"/"+id, nil)
	if err != nil {
		fmt.Printf("Ошибка при создании запроса: %v\n", err)
		return
	}
	req.Header.Set("Authorization", "Bearer "+sessionToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Ошибка при получении пользователя: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var user User
		body, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal(body, &user)
		fmt.Printf("Данные пользователя: ID: %s, Имя: %s, Возраст: %d\n", user.Id, user.Name, user.Age)
	} else {
		handleErrorResponse(resp)
	}
}

func deleteUser() {
	defer wg.Done()
	var id string
	fmt.Print("Введите ID пользователя для удаления: ")
	fmt.Scan(&id)

	req, err := http.NewRequest("DELETE", apiBaseURL+"/"+id, nil)
	if err != nil {
		fmt.Printf("Ошибка при удалении пользователя: %v\n", err)
		return
	}
	req.Header.Set("Authorization", "Bearer "+sessionToken)

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
		handleErrorResponse(resp)
	}
}

func handleErrorResponse(resp *http.Response) {
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Printf("Ошибка: %s. Подробности: %s\n", resp.Status, string(body))
}
