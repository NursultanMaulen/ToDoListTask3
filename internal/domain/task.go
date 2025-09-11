package domain

import "time"

// Task представляет собой доменную модель задачи
type Task struct {
	ID        int        `json:"id"`
	Text      string     `json:"text"`
	Completed bool       `json:"completed"`
	DueDate   *time.Time `json:"dueDate"`
	Priority  string     `json:"priority"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
}

// TaskRepository интерфейс для работы с хранилищем задач
type TaskRepository interface {
	Create(task *Task) error
	Update(task *Task) error
	Delete(id int) error
	GetByID(id int) (*Task, error)
	GetAll() ([]Task, error)
	GetNextID() (int, error)
}

// Settings представляет настройки приложения
type Settings struct {
	ID       int  `json:"id"`
	DarkMode bool `json:"darkMode"`
}
