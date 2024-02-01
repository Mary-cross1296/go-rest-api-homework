package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Task ...
type Task struct {
	ID           string   `json:"id"`
	Description  string   `json:"description"`
	Note         string   `json:"note"`
	Applications []string `json:"applications"`
}

var tasks = map[string]Task{
	"1": {
		ID:          "1",
		Description: "Сделать финальное задание темы REST API",
		Note:        "Если сегодня сделаю, то завтра будет свободный день. Ура!",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
		},
	},
	"2": {
		ID:          "2",
		Description: "Протестировать финальное задание с помощью Postmen",
		Note:        "Лучше это делать в процессе разработки, каждый раз, когда запускаешь сервер и проверяешь хендлер",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
			"Postman",
		},
	},
}

// Ниже напишите обработчики для каждого эндпоинта
// Обработчик для получения всех задач
func getTasks(w http.ResponseWriter, r *http.Request) {
	//сериализуем данные из слайса tasks
	resp, err := json.Marshal(tasks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// в заголовок записываем тип контента, у нас это данные в формате JSON
	w.Header().Set("Content-Type", "application/json")
	//устанавливаем статус ответа 200 ок
	w.WriteHeader(http.StatusOK)
	//записывает данные в формате JSON в тело ответа
	w.Write(resp)
}

// Обработчик для отправки задачи на сервер
func postTask(w http.ResponseWriter, r *http.Request) {
	var task Task
	var taskByte bytes.Buffer

	//Читаем тело запроса
	_, err := taskByte.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//Десериализуем данные из JSON тела запроса и записываем их в переменную
	err = json.Unmarshal(taskByte.Bytes(), &task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//Проверяем, если в мапе task задада с ключом, который совпадает с id новой задачи
	_, ok := tasks[task.ID]
	if !ok {
		//Записываем новую задачу в мапу, id указанный в задаче будет соответсвовать ключу,
		//под которым эта задача запишется в мапу
		tasks[task.ID] = task
		//Добавляем заголовок
		w.Header().Set("Content-Type", "application/json")
		//Возвращаем статус
		w.WriteHeader(http.StatusCreated)
	} else {
		//Возвращаем статус конфликта, т.к. в мапе есть задача с таким id
		w.WriteHeader(http.StatusBadRequest)
	}
}

// Обработчик для получения задачи по ID
func getTaskByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	task, ok := tasks[id]
	if !ok {
		http.Error(w, "Задача не найден", http.StatusBadRequest)
	}

	resp, err := json.Marshal(task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

// Обработчик удаления задачи по ID
func deleteTaskById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	_, ok := tasks[id]
	if !ok {
		http.Error(w, "Задача не найден", http.StatusBadRequest)
	}

	delete(tasks, id)

	resp, err := json.Marshal(tasks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func main() {
	r := chi.NewRouter()

	// здесь регистрируйте ваши обработчики
	r.Get("/tasks", getTasks)
	r.Post("/tasks", postTask)
	r.Get("/tasks/{id}", getTaskByID)
	r.Delete("/tasks/{id}", deleteTaskById)

	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
