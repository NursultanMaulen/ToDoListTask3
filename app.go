package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"myproject/internal/domain"
	"myproject/internal/repository"
	"myproject/internal/service"
)

// Config структура для конфигурации
type Config struct {
	Database struct {
		Host     string `json:"host"`
		Port     int    `json:"port"`
		User     string `json:"user"`
		Password string `json:"password"`
		DBName   string `json:"dbname"`
	} `json:"database"`
}

// App struct
type App struct {
	ctx         context.Context
	taskService *service.TaskService
	repo        *repository.PostgresRepository
	initOnce    sync.Once
	config      Config
}

// loadConfig загружает конфигурацию из файла
func loadConfig() (Config, error) {
	var config Config
	
	// Пытаемся загрузить config.json
	file, err := os.ReadFile("config.json")
	if err != nil {
		// Если config.json не найден, пробуем загрузить пример конфигурации
		file, err = os.ReadFile("config.example.json")
		if err != nil {
			return config, fmt.Errorf("error reading config file: %v", err)
		}
	}

	if err := json.Unmarshal(file, &config); err != nil {
		return config, fmt.Errorf("error parsing config file: %v", err)
	}

	return config, nil
}

// NewApp создает новый экземпляр приложения
func NewApp() *App {
	config, err := loadConfig()
	if err != nil {
		panic(fmt.Sprintf("Failed to load configuration: %v", err))
	}

	// Формируем строку подключения к PostgreSQL
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Database.Host, config.Database.Port, config.Database.User, 
		config.Database.Password, config.Database.DBName)

	println("Connecting to database:", connStr)

	// Инициализируем репозиторий
	repo, err := repository.NewPostgresRepository(connStr)
	if err != nil {
		println("Error connecting to database:", err.Error())
		panic(fmt.Sprintf("Failed to connect to database: %v", err))
	}

	println("Successfully connected to database")

	// Инициализируем структуру базы данных
	if err := repo.InitDB(); err != nil {
		println("Error initializing database:", err.Error())
		panic(fmt.Sprintf("Failed to initialize database: %v", err))
	}

	println("Database initialized successfully")

	// Создаем сервис для работы с задачами
	taskService := service.NewTaskService(repo)

	return &App{
		ctx:         context.Background(),
		taskService: taskService,
		repo:        repo,
		config:      config,
	}
}

// ensureInitialized гарантирует, что все зависимости инициализированы
func (a *App) ensureInitialized() {
	a.initOnce.Do(func() {
		if a.ctx == nil {
			a.ctx = context.Background()
		}

		// Если конфигурация не загружена, загружаем её
		if (Config{}) == a.config {
			config, err := loadConfig()
			if err != nil {
				println("Error loading config:", err.Error())
				return
			}
			a.config = config
		}

		// Если репозиторий не инициализирован, инициализируем его
		if a.repo == nil {
			connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
				a.config.Database.Host, a.config.Database.Port, a.config.Database.User, 
				a.config.Database.Password, a.config.Database.DBName)
			println("ensureInitialized: connecting to database:", connStr)
			repo, err := repository.NewPostgresRepository(connStr)
			if err != nil {
				println("ensureInitialized: repo connect error:", err.Error())
				return
			}
			if err := repo.InitDB(); err != nil {
				println("ensureInitialized: initdb error:", err.Error())
				return
			}
			a.repo = repo
		}

		if a.taskService == nil && a.repo != nil {
			a.taskService = service.NewTaskService(a.repo)
		}
	})
}

// Startup вызывается при запуске приложения
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
}

// AddTask добавляет новую задачу
func (a *App) AddTask(text string, dueDate *time.Time, priority string) domain.Task {
	println("AddTask called with context:", a.ctx != nil)
	a.ensureInitialized()
	if a.taskService == nil {
		println("taskService is nil in AddTask")
		return domain.Task{}
	}
	task, err := a.taskService.CreateTask(text, dueDate, priority)
	if err != nil {
		println("Error adding task:", err.Error())
		return domain.Task{}
	}
	return *task
}

// GetTasks возвращает все задачи
func (a *App) GetTasks() []domain.Task {
	println("GetTasks called")
	a.ensureInitialized()
	if a.taskService == nil {
		println("taskService is nil!")
		return []domain.Task{}
	}
	tasks, err := a.taskService.GetAllTasks()
	if err != nil {
		println("Error getting tasks:", err.Error())
		return []domain.Task{}
	}
	println("Successfully got", len(tasks), "tasks")
	return tasks
}

// ToggleTask переключает статус задачи
func (a *App) ToggleTask(id int) []domain.Task {
	a.ensureInitialized()
	if a.taskService == nil {
		println("Context or taskService is nil in ToggleTask")
		return []domain.Task{}
	}
	err := a.taskService.ToggleTask(id)
	if err != nil {
		println("Error toggling task:", err.Error())
	}
	return a.GetTasks()
}

// DeleteTask удаляет задачу
func (a *App) DeleteTask(id int) []domain.Task {
	a.ensureInitialized()
	if a.taskService == nil {
		println("Context or taskService is nil in DeleteTask")
		return []domain.Task{}
	}
	err := a.taskService.DeleteTask(id)
	if err != nil {
		println("Error deleting task:", err.Error())
	}
	return a.GetTasks()
}

// GetDarkMode возвращает текущую тему
func (a *App) GetDarkMode() bool {
	a.ensureInitialized()
	if a.repo == nil {
		println("Context or repo is nil in GetDarkMode")
		return false
	}
	settings, err := a.repo.GetSettings()
	if err != nil {
		println("Error getting dark mode:", err.Error())
		return false
	}
	return settings.DarkMode
}

// SetDarkMode устанавливает тему
func (a *App) SetDarkMode(isDark bool) bool {
	a.ensureInitialized()
	if a.repo == nil {
		println("Context or repo is nil in SetDarkMode")
		return isDark
	}
	settings, err := a.repo.GetSettings()
	if err != nil {
		println("Error getting settings:", err.Error())
		return isDark
	}
	
	settings.DarkMode = isDark
	err = a.repo.UpdateSettings(settings)
	if err != nil {
		println("Error updating dark mode:", err.Error())
		return isDark
	}
	return settings.DarkMode
}
