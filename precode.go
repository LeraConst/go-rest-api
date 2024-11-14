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
	ID          string   `json:"id"`          // ID задачи
	Description string   `json:"description"` // Заголовок
	Note        string   `json:"note"`        // Описание задачи
	Application []string `json:"application"` // Приложения, которыми будете пользоваться
}

var tasks = map[string]Task{
	"1": {
		ID:          "1",
		Description: "Сделать финальное задание темы REST API",
		Note:        "Если сегодня сделаю, то завтра будет свободный день. Ура!",
		Application: []string{
			"VS Code",
			"Terminal",
			"git",
		},
	},
	"2": {
		ID:          "2",
		Description: "Протестировать финальное задание с помощью Postmen",
		Note:        "Лучше это делать в процессе разработки, каждый раз, когда запускаешь сервер и проверяешь хендлер",
		Application: []string{
			"VS Code",
			"Terminal",
			"git",
			"Postman",
		},
	},
}

func getTasksHandler(w http.ResponseWriter, r *http.Request) {
	// сериализуем данные из слайса tasks
	resp, err := json.Marshal(tasks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Устанавливаем заголовки и статус успешного ответа
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, resp)
}

// обрабатывает запрос POST /tasks, принимает и сохраняет задачу в мапе
func postTasksHandler(w http.ResponseWriter, r *http.Request) {
	var task Task
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tasks[task.ID] = task

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

// обрабатывает запрос DELETE /tasks/{id} и удаляет задачу по ID
func deleteTasksHandler(w http.ResponseWriter, r *http.Request) {
	// Извлекаем ID задачи из URL
	taskID := chi.URLParam(r, "id")

	// Проверяем, существует ли задача с данным ID
	if _, exists := tasks[taskID]; !exists {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	// Удаляем задачу из мапы
	delete(tasks, taskID)

	// Устанавливаем заголовки и статус успешного ответа
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

// обрабатывает запрос GET /tasks/{id} и возвращает задачу по ID
func getTaskByIdHandler(w http.ResponseWriter, r *http.Request) {
	// Извлекаем ID задачи из URL
	taskID := chi.URLParam(r, "id")

	// Проверяем, существует ли задача с данным ID
	task, exists := tasks[taskID]
	if !exists {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	// Сериализуем задачу в JSON
	resp, err := json.Marshal(task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Устанавливаем заголовки и статус успешного ответа
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, resp)
}

func main() {
	r := chi.NewRouter()

	// регистрируем в роутере эндпоинт `/tasks` с методом GET, для которого используется обработчик `getTasksHandler`
	r.Get("/tasks", getTasksHandler)
	// регистрируем в роутере эндпоинт `/tasks` с методом POST, для которого используется обработчик `postTasksHandler`
	r.Post("/tasks", postTasksHandler)
	// регистрируем в роутере эндпоинт `/tasks/{id}` с методом DELETE, для которого используется обработчик `deleteTasksHandler`
	r.Delete("/tasks/{id}", deleteTasksHandler)
	// регистрируем в роутере эндпоинт `/tasks/{id}` с методом GET, для которого используется обработчик `getTaskByIdHandler`
	r.Get("/tasks/{id}", getTaskByIdHandler)

	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
