package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const baseURL = "http://localhost:8080/users"

type User struct {
	Id   primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name string             `bson:"name" json:"name"`
	Age  int                `bson:"age" json:"age"`
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
			getUserByID()
		case 3:
			createUser()
		case 4:
			updateUser()
		case 5:
			deleteUser()
		case 6:
			fmt.Println("Выход из программы...")
			os.Exit(0)
		default:
			fmt.Println("Неверный выбор. Попробуйте еще раз.")
		}
	}
}

func getAllUsers() {
	resp, err := http.Get(baseURL)
	if err != nil {
		fmt.Println("Ошибка при получении пользователей:", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Ошибка при чтении ответа:", err)
		return
	}

	if resp.StatusCode == http.StatusOK {
		var users []User
		json.Unmarshal(body, &users)
		for _, user := range users {
			fmt.Printf("ID: %s, Name: %s, Age: %d\n", user.Id.Hex(), user.Name, user.Age)
		}
	} else {
		fmt.Println("Ошибка:", string(body))
	}
}

func getUserByID() {
	fmt.Print("Введите ID пользователя: ")
	var id string
	fmt.Scan(&id)

	resp, err := http.Get(baseURL + "/" + id)
	if err != nil {
		fmt.Println("Ошибка при получении пользователя:", err)
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
		json.Unmarshal(body, &user)
		fmt.Printf("ID: %s, Name: %s, Age: %d\n", user.Id.Hex(), user.Name, user.Age)
	} else {
		fmt.Println("Ошибка:", string(body))
	}
}

func createUser() {
	var user User
	fmt.Print("Введите имя пользователя: ")
	fmt.Scan(&user.Name)
	fmt.Print("Введите возраст пользователя: ")
	fmt.Scan(&user.Age)

	body, _ := json.Marshal(user)
	resp, err := http.Post(baseURL, "application/json", bytes.NewBuffer(body))
	if err != nil {
		fmt.Println("Ошибка при создании пользователя:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusCreated {
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println("Пользователь создан:", string(body))
	} else {
		fmt.Println("Ошибка:", resp.Status)
	}
}

func updateUser() {
	var user User
	fmt.Print("Введите ID пользователя для обновления: ")
	fmt.Scan(&user.Id)
	fmt.Print("Введите новое имя пользователя: ")
	fmt.Scan(&user.Name)
	fmt.Print("Введите новый возраст пользователя: ")
	fmt.Scan(&user.Age)

	body, _ := json.Marshal(user)
	req, err := http.NewRequest("PUT", baseURL+"/"+user.Id.Hex(), bytes.NewBuffer(body))
	if err != nil {
		fmt.Println("Ошибка при обновлении пользователя:", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Ошибка при обновлении пользователя:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println("Пользователь обновлён:", string(body))
	} else {
		fmt.Println("Ошибка:", resp.Status)
	}
}

func deleteUser() {
	fmt.Print("Введите ID пользователя для удаления: ")
	var id string
	fmt.Scan(&id)

	req, err := http.NewRequest("DELETE", baseURL+"/"+id, nil)
	if err != nil {
		fmt.Println("Ошибка при удалении пользователя:", err)
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Ошибка при удалении пользователя:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent {
		fmt.Println("Пользователь удалён")
	} else {
		fmt.Println("Ошибка:", resp.Status)
	}
}
