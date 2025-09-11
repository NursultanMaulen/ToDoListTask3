# Todo List Desktop Application

Modern desktop todo list application built with Wails, React, and PostgreSQL. Features a beautiful dark/light theme and responsive design.

## Demo

Watch the application in action:

![Application Demo](./demo/demo.gif)

## Screenshots

### Light and Dark Themes

| Light Mode                                  | Dark Mode                                 |
| ------------------------------------------- | ----------------------------------------- |
| ![Light Mode](./screenshots/light-mode.png) | ![Dark Mode](./screenshots/dark-mode.png) |

### Task Filtering

![Filtering Interface](./screenshots/filtering.png)

### Database View

![Tasks Database](./screenshots/tasksdb.png)

### Console Output

![Console View](./screenshots/console.png)

## Checklist

### 1. Интерфейс пользователя

Основная часть:

- [x] Создан интерфейс с текстовым полем для ввода новой задачи
- [x] Добавлена кнопка для добавления задачи в список
- [x] Реализовано отображение списка всех задач
- [x] Использован CSS для стилизации интерфейса
- [x] Добавлены значки и цвета для обозначения статуса задач

Бонусная часть:

- [x] Реализована адаптивная верстка
- [x] Добавлено переключение светлой/тёмной темы

### 2. Добавление задач

Основная часть:

- [x] Реализован функционал добавления новой задачи
- [x] Добавлена валидация ввода

Бонусная часть:

- [x] Реализовано добавление задач с датой выполнения
- [x] Добавлена установка приоритета (низкий, средний, высокий)

### 3. Удаление задач

Основная часть:

- [x] Реализовано удаление задач из списка

Бонусная часть:

- [x] Добавлено модальное окно подтверждения удаления

### 4. Управление выполнением задач

Основная часть:

- [x] Реализована отметка задачи как выполненной
- [x] Добавлено зачеркивание текста выполненных задач

Бонусная часть:

- [x] Реализовано перемещение в раздел "Выполненные задачи"
- [x] Добавлена возможность отмены выполнения задачи

### 5. Сохранение состояния

Основная часть:

- [x] Реализовано сохранение состояния при закрытии
- [x] Добавлена загрузка состояния при запуске

Бонусная часть:

- [x] Интегрирована PostgreSQL база данных
- [x] Реализована архитектура repo → service → usecase

### 6. Фильтрация и сортировка задач

Основная часть:

- [x] Добавлена фильтрация по статусу
- [x] Реализована сортировка по дате добавления

Бонусная часть:

- [x] Добавлена сортировка по приоритету
- [x] Реализована фильтрация по дате

## Technologies

- [Wails](https://wails.io/) - Go framework for desktop applications
- [React](https://reactjs.org/) - Frontend UI library
- [PostgreSQL](https://www.postgresql.org/) - Database
- Modern CSS with variables for theming

## Prerequisites

Before you begin, ensure you have the following installed:

- Go 1.21 or later
- Node.js 16 or later
- PostgreSQL
- Wails CLI

## Getting Started

### Installation

1. Clone the repository:

```bash
git clone <repository-url>
cd myproject
```

2. Install frontend dependencies:

```bash
cd frontend
npm install
cd ..
```

3. Set up the PostgreSQL database:

```sql
CREATE DATABASE todolist;
```

4. Configure the database connection:
   - Copy `config.example.json` to `config.json`
   - Update the database settings in `config.json` with your PostgreSQL credentials

### Development

- Run the application in development mode:

```bash
wails dev
```

## Project Structure

```
myproject/
├── app.go              # Main application logic
├── main.go            # Entry point
├── frontend/          # React frontend
│   ├── src/
│   │   ├── App.jsx   # Main React component
│   │   └── App.css   # Styles with theme support
├── internal/          # Backend Go code
│   ├── domain/       # Business logic
│   ├── repository/   # Database interactions
│   └── service/      # Business services
└── build/            # Build configurations
```

## Features in Detail

### Task Management

- Create new tasks with title, priority, and due date
- Mark tasks as complete/incomplete
- Delete tasks with confirmation
- Edit existing tasks

### User Interface

- Responsive design that works on all screen sizes
- Smooth animations and transitions
- Intuitive controls and layout
- Dark/light theme toggle with CSS variables

### Data Persistence

- All tasks are stored in PostgreSQL
- Reliable data storage and retrieval
- Efficient database queries

## License

This project is licensed under the MIT License - see the LICENSE file for details.
