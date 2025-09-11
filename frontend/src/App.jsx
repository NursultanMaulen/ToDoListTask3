/**
 * Todo List Application
 *
 * Функционал:
 * 1. Управление задачами:
 *    - Создание новых задач
 *    - Установка приоритета (низкий, средний, высокий)
 *    - Установка дедлайна (опционально)
 *    - Отметка о выполнении
 *    - Удаление задач с подтверждением
 *
 * 2. Интерфейс:
 *    - Адаптивный дизайн
 *    - Поддержка светлой/темной темы
 *    - Модальные окна для подтверждений
 *    - Визуальная индикация приоритетов
 */

import { useState, useEffect } from "react";
import "./App.css";
import {
  AddTask,
  GetTasks,
  ToggleTask,
  DeleteTask,
  GetDarkMode,
  SetDarkMode,
} from "../wailsjs/go/main/App";

function App() {
  // Состояния для управления задачами
  const [tasks, setTasks] = useState([]);
  const [newTask, setNewTask] = useState("");
  const [dueDate, setDueDate] = useState("");
  const [priority, setPriority] = useState("medium");

  // Состояния интерфейса
  const [darkMode, setDarkMode] = useState(false);
  const [deleteModal, setDeleteModal] = useState(false);
  const [taskToDelete, setTaskToDelete] = useState(null);
  // Фильтрация и сортировка
  const [statusFilter, setStatusFilter] = useState("all"); // all | active | completed
  const [sortBy, setSortBy] = useState("date_desc"); // date_desc | date_asc | priority
  const [dateFilter, setDateFilter] = useState("all"); // all | today | week | overdue

  useEffect(() => {
    const initApp = async () => {
      // Загружаем задачи
      const tasksData = await GetTasks();
      setTasks(tasksData || []);

      // Загружаем сохраненную тему
      const savedDarkMode = await GetDarkMode();
      setDarkMode(savedDarkMode);
    };

    initApp();
  }, []);

  const loadTasks = async () => {
    const tasksData = await GetTasks();
    setTasks(tasksData || []);
  };

  const handleAddTask = async (e) => {
    e.preventDefault();
    if (newTask.trim()) {
      const taskDueDate = dueDate ? new Date(dueDate) : null;
      await AddTask(newTask, taskDueDate, priority);
      setNewTask("");
      setDueDate("");
      setPriority("medium");
      loadTasks();
    }
  };

  const handleToggleTask = async (id) => {
    await ToggleTask(id);
    loadTasks();
  };

  const toggleTheme = async () => {
    const newDarkMode = !darkMode;
    const updated = await SetDarkMode(newDarkMode);
    setDarkMode(updated);
  };

  // Открывает модальное окно подтверждения удаления
  const handleDeleteClick = (e, task) => {
    e.stopPropagation(); // Предотвращаем срабатывание onClick родительского элемента
    setTaskToDelete(task);
    setDeleteModal(true);
  };

  // Подтверждение удаления задачи
  const confirmDelete = async () => {
    if (taskToDelete) {
      await DeleteTask(taskToDelete.id);
      setDeleteModal(false);
      setTaskToDelete(null);
      loadTasks();
    }
  };

  // Отмена удаления
  const cancelDelete = () => {
    setDeleteModal(false);
    setTaskToDelete(null);
  };

  // Применяем фильтрацию и сортировку на клиенте
  const filteredAndSorted = (() => {
    const now = new Date();

    // Клонируем список
    let list = (tasks || []).slice();

    // Фильтр по статусу
    if (statusFilter === "active") {
      list = list.filter((t) => !t.completed);
    } else if (statusFilter === "completed") {
      list = list.filter((t) => t.completed);
    }

    // Фильтр по дате
    if (dateFilter === "today") {
      list = list.filter((t) => {
        if (!t.dueDate) return false;
        const d = new Date(t.dueDate);
        return d.toDateString() === now.toDateString();
      });
    } else if (dateFilter === "week") {
      const weekAhead = new Date(now);
      weekAhead.setDate(now.getDate() + 7);
      list = list.filter((t) => {
        if (!t.dueDate) return false;
        const d = new Date(t.dueDate);
        return d >= now && d <= weekAhead;
      });
    } else if (dateFilter === "overdue") {
      list = list.filter((t) => {
        if (!t.dueDate) return false;
        const d = new Date(t.dueDate);
        return d < now && !t.completed;
      });
    }

    // Сортировка
    const priorityRank = (p) => {
      if (!p) return 1;
      if (p === "high") return 3;
      if (p === "medium") return 2;
      return 1;
    };

    list.sort((a, b) => {
      if (sortBy === "date_asc" || sortBy === "date_desc") {
        const da = a.createdAt ? new Date(a.createdAt) : new Date(0);
        const db = b.createdAt ? new Date(b.createdAt) : new Date(0);
        return sortBy === "date_asc" ? da - db : db - da;
      }

      if (sortBy === "priority") {
        return priorityRank(b.priority) - priorityRank(a.priority);
      }

      return 0;
    });

    return list;
  })();

  return (
    <div id="App" className={darkMode ? "dark-mode" : "light-mode"}>
      <div className="container">
        <header>
          <h1>Todo List</h1>
          <button className="theme-toggle" onClick={toggleTheme}>
            {darkMode ? "☀️" : "🌙"}
          </button>
        </header>

        <form onSubmit={handleAddTask} className="task-form">
          <div className="input-group">
            <input
              type="text"
              value={newTask}
              onChange={(e) => setNewTask(e.target.value)}
              placeholder="Add a new task..."
              className="task-input"
              required
            />
            <input
              type="datetime-local"
              value={dueDate}
              onChange={(e) => setDueDate(e.target.value)}
              className="date-input"
            />
            <select
              value={priority}
              onChange={(e) => setPriority(e.target.value)}
              className="priority-select"
            >
              <option value="low">Low Priority</option>
              <option value="medium">Medium Priority</option>
              <option value="high">High Priority</option>
            </select>
          </div>
          <button type="submit" className="add-btn">
            Add Task
          </button>
        </form>

        {/* Контролы фильтрации и сортировки */}
        <div className="controls">
          <div className="filter-group">
            <label>Filter:</label>
            <select
              value={statusFilter}
              onChange={(e) => setStatusFilter(e.target.value)}
            >
              <option value="all">All</option>
              <option value="active">Active</option>
              <option value="completed">Completed</option>
            </select>
          </div>

          <div className="filter-group">
            <label>Date:</label>
            <select
              value={dateFilter}
              onChange={(e) => setDateFilter(e.target.value)}
            >
              <option value="all">Any</option>
              <option value="today">Today</option>
              <option value="week">This week</option>
              <option value="overdue">Overdue</option>
            </select>
          </div>

          <div className="filter-group">
            <label>Sort:</label>
            <select value={sortBy} onChange={(e) => setSortBy(e.target.value)}>
              <option value="date_desc">Date (newest)</option>
              <option value="date_asc">Date (oldest)</option>
              <option value="priority">Priority</option>
            </select>
          </div>
        </div>

        <div className="tasks-list">
          {(filteredAndSorted || []).map((task) => (
            <div
              key={task.id}
              className={`task-item ${task.completed ? "completed" : ""}`}
              onClick={() => handleToggleTask(task.id)}
            >
              <span className="task-checkbox">
                {task.completed ? "✓" : "○"}
              </span>
              <div className="task-content">
                <span className="task-text">{task.text}</span>
                <div className="task-details">
                  <span className={`priority-badge ${task.priority}`}>
                    {task.priority}
                  </span>
                  {task.dueDate && (
                    <span className="due-date">
                      Deadline: {new Date(task.dueDate).toLocaleString()}
                    </span>
                  )}
                </div>
              </div>
              <button
                className="delete-btn"
                onClick={(e) => handleDeleteClick(e, task)}
              >
                ×
              </button>
            </div>
          ))}
        </div>

        {/* Модальное окно подтверждения удаления */}
        {deleteModal && (
          <div className="modal-overlay">
            <div className="modal">
              <h2>Подтверждение удаления</h2>
              <p>
                Вы уверены, что хотите удалить задачу "{taskToDelete?.text}"?
              </p>
              <div className="modal-buttons">
                <button className="modal-btn cancel" onClick={cancelDelete}>
                  Отмена
                </button>
                <button className="modal-btn confirm" onClick={confirmDelete}>
                  Удалить
                </button>
              </div>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}

export default App;
