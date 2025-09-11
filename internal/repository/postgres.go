package repository

import (
	"database/sql"
	"time"

	"myproject/internal/domain"

	_ "github.com/lib/pq"
)

// PostgresRepository реализация TaskRepository для PostgreSQL
type PostgresRepository struct {
	db *sql.DB
}

// NewPostgresRepository создает новый экземпляр PostgresRepository
func NewPostgresRepository(connectionString string) (*PostgresRepository, error) {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}

	// Устанавливаем максимальное время ожидания для проверки соединения
	db.SetConnMaxLifetime(time.Second * 5)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)

	// Проверяем соединение несколько раз
	for i := 0; i < 3; i++ {
		if err = db.Ping(); err == nil {
			break
		}
		time.Sleep(time.Second)
	}
	if err != nil {
		return nil, err
	}

	return &PostgresRepository{db: db}, nil
}

// InitDB инициализирует структуру базы данных
func (r *PostgresRepository) InitDB() error {
	// Начинаем транзакцию
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Создаем таблицу tasks
	_, err = tx.Exec(`
		CREATE TABLE IF NOT EXISTS tasks (
			id SERIAL PRIMARY KEY,
			text TEXT NOT NULL,
			completed BOOLEAN DEFAULT FALSE,
			due_date TIMESTAMP,
			priority TEXT,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return err
	}

	// Создаем таблицу settings
	_, err = tx.Exec(`
		CREATE TABLE IF NOT EXISTS settings (
			id SERIAL PRIMARY KEY,
			dark_mode BOOLEAN DEFAULT FALSE
		)
	`)
	if err != nil {
		return err
	}

	// Проверяем существование начальных настроек
	var count int
	err = tx.QueryRow("SELECT COUNT(*) FROM settings").Scan(&count)
	if err != nil {
		return err
	}

	// Если настроек нет, создаем дефолтные
	if count == 0 {
		_, err = tx.Exec(`
			INSERT INTO settings (dark_mode) 
			VALUES (false)
		`)
		if err != nil {
			return err
		}
	}

	// Фиксируем транзакцию
	return tx.Commit()
}

// Create добавляет новую задачу
func (r *PostgresRepository) Create(task *domain.Task) error {
	query := `
		INSERT INTO tasks (text, completed, due_date, priority, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id`

	now := time.Now()
	task.CreatedAt = now
	task.UpdatedAt = now

	return r.db.QueryRow(
		query,
		task.Text,
		task.Completed,
		task.DueDate,
		task.Priority,
		task.CreatedAt,
		task.UpdatedAt,
	).Scan(&task.ID)
}

// Update обновляет существующую задачу
func (r *PostgresRepository) Update(task *domain.Task) error {
	query := `
		UPDATE tasks
		SET text = $1, completed = $2, due_date = $3, priority = $4, updated_at = $5
		WHERE id = $6`

	task.UpdatedAt = time.Now()

	result, err := r.db.Exec(
		query,
		task.Text,
		task.Completed,
		task.DueDate,
		task.Priority,
		task.UpdatedAt,
		task.ID,
	)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// Delete удаляет задачу по ID
func (r *PostgresRepository) Delete(id int) error {
	query := `DELETE FROM tasks WHERE id = $1`
	
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// GetByID получает задачу по ID
func (r *PostgresRepository) GetByID(id int) (*domain.Task, error) {
	task := &domain.Task{}
	
	query := `
		SELECT id, text, completed, due_date, priority, created_at, updated_at
		FROM tasks
		WHERE id = $1`

	err := r.db.QueryRow(query, id).Scan(
		&task.ID,
		&task.Text,
		&task.Completed,
		&task.DueDate,
		&task.Priority,
		&task.CreatedAt,
		&task.UpdatedAt,
	)
	
	if err != nil {
		return nil, err
	}

	return task, nil
}

// GetAll получает все задачи
func (r *PostgresRepository) GetAll() ([]domain.Task, error) {
	query := `
		SELECT id, text, completed, due_date, priority, created_at, updated_at
		FROM tasks
		ORDER BY created_at DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []domain.Task
	for rows.Next() {
		var task domain.Task
		err := rows.Scan(
			&task.ID,
			&task.Text,
			&task.Completed,
			&task.DueDate,
			&task.Priority,
			&task.CreatedAt,
			&task.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

// GetNextID получает следующий доступный ID
func (r *PostgresRepository) GetNextID() (int, error) {
	var nextID int
	err := r.db.QueryRow("SELECT COALESCE(MAX(id) + 1, 1) FROM tasks").Scan(&nextID)
	if err != nil {
		return 0, err
	}
	return nextID, nil
}

// GetSettings получает настройки приложения
func (r *PostgresRepository) GetSettings() (*domain.Settings, error) {
	settings := &domain.Settings{}
	
	err := r.db.QueryRow(`
		SELECT id, dark_mode 
		FROM settings 
		ORDER BY id 
		LIMIT 1
	`).Scan(&settings.ID, &settings.DarkMode)

	if err == sql.ErrNoRows {
		// Если настроек нет, создаем дефолтные
		err = r.db.QueryRow(`
			INSERT INTO settings (dark_mode) 
			VALUES (false) 
			RETURNING id, dark_mode
		`).Scan(&settings.ID, &settings.DarkMode)
	}

	if err != nil {
		return nil, err
	}

	return settings, nil
}

// UpdateSettings обновляет настройки приложения
func (r *PostgresRepository) UpdateSettings(settings *domain.Settings) error {
	_, err := r.db.Exec(`
		UPDATE settings 
		SET dark_mode = $1 
		WHERE id = $2
	`, settings.DarkMode, settings.ID)
	return err
}
