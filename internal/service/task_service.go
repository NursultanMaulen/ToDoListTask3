package service

import (
	"time"

	"myproject/internal/domain"
)

// TaskService сервис для работы с задачами
type TaskService struct {
	repo domain.TaskRepository
}

// NewTaskService создает новый экземпляр TaskService
func NewTaskService(repo domain.TaskRepository) *TaskService {
	return &TaskService{
		repo: repo,
	}
}

// CreateTask создает новую задачу
func (s *TaskService) CreateTask(text string, dueDate *time.Time, priority string) (*domain.Task, error) {
	nextID, err := s.repo.GetNextID()
	if err != nil {
		return nil, err
	}

	task := &domain.Task{
		ID:        nextID,
		Text:      text,
		Completed: false,
		DueDate:   dueDate,
		Priority:  priority,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.repo.Create(task); err != nil {
		return nil, err
	}

	return task, nil
}

// ToggleTask переключает статус выполнения задачи
func (s *TaskService) ToggleTask(id int) error {
	task, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	task.Completed = !task.Completed
	task.UpdatedAt = time.Now()

	return s.repo.Update(task)
}

// DeleteTask удаляет задачу
func (s *TaskService) DeleteTask(id int) error {
	return s.repo.Delete(id)
}

// GetAllTasks возвращает все задачи
func (s *TaskService) GetAllTasks() ([]domain.Task, error) {
	return s.repo.GetAll()
}

// UpdateTask обновляет задачу
func (s *TaskService) UpdateTask(task *domain.Task) error {
	task.UpdatedAt = time.Now()
	return s.repo.Update(task)
}
