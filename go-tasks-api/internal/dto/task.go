package dto

import (
	"time"
)

// CreateTaskRequest define la estructura para crear una tarea
type CreateTaskRequest struct {
	Title       string `json:"title" validate:"required,min=1,max=100"`
	Description string `json:"description" validate:"max=500"`
	Priority    string `json:"priority" validate:"omitempty,oneof=low medium high"`
	DueDate     string `json:"due_date" validate:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
}

// UpdateTaskRequest define la estructura para actualizar una tarea
type UpdateTaskRequest struct {
	Title       *string `json:"title" validate:"omitempty,min=1,max=100"`
	Description *string `json:"description" validate:"omitempty,max=500"`
	Priority    *string `json:"priority" validate:"omitempty,oneof=low medium high"`
	DueDate     *string `json:"due_date" validate:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
	Completed   *bool   `json:"completed"`
}

// TaskResponse define la estructura de respuesta de una tarea
type TaskResponse struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Priority    string    `json:"priority"`
	Completed   bool      `json:"completed"`
	DueDate     *string   `json:"due_date,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TaskListResponse define la respuesta para listar tareas
type TaskListResponse struct {
	Tasks []TaskResponse `json:"tasks"`
	Total int            `json:"total"`
	Page  int            `json:"page"`
	Limit int            `json:"limit"`
}

// ErrorResponse define la estructura estándar de errores
type ErrorResponse struct {
	Error   string            `json:"error"`
	Message string            `json:"message"`
	Code    int               `json:"code"`
	Details map[string]string `json:"details,omitempty"`
}

// SuccessResponse define respuestas exitosas simples
type SuccessResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}
