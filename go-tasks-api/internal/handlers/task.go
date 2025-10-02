package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/Kathecam/go-tasks-api/internal/dto"
	"github.com/Kathecam/go-tasks-api/internal/errors"
	"github.com/Kathecam/go-tasks-api/internal/validator"
)

// TaskHandler maneja las operaciones de tareas
type TaskHandler struct {
	validator *validator.Validator
}

// NewTaskHandler crea una nueva instancia del handler
func NewTaskHandler() *TaskHandler {
	return &TaskHandler{
		validator: validator.New(),
	}
}

// GetTasks maneja GET /tasks
func (h *TaskHandler) GetTasks(c *fiber.Ctx) error {
	// Por ahora retornamos datos mock
	tasks := []dto.TaskResponse{
		{
			ID:          "1",
			Title:       "Learn Go",
			Description: "Study Go programming language",
			Priority:    "high",
			Completed:   false,
			CreatedAt:   time.Now().Add(-24 * time.Hour),
			UpdatedAt:   time.Now().Add(-24 * time.Hour),
		},
		{
			ID:          "2",
			Title:       "Build API",
			Description: "Create REST API with Fiber",
			Priority:    "medium",
			Completed:   true,
			CreatedAt:   time.Now().Add(-12 * time.Hour),
			UpdatedAt:   time.Now().Add(-6 * time.Hour),
		},
	}

	response := dto.TaskListResponse{
		Tasks: tasks,
		Total: len(tasks),
		Page:  1,
		Limit: 10,
	}

	return c.JSON(response)
}

// GetTaskByID maneja GET /tasks/:id
func (h *TaskHandler) GetTaskByID(c *fiber.Ctx) error {
	id := c.Params("id")

	// Validar que el ID no esté vacío
	if id == "" {
		return errors.ErrInvalidInput.WithDetails(map[string]string{
			"id": "Task ID is required",
		})
	}

	// Validar formato UUID (opcional pero buena práctica)
	if _, err := uuid.Parse(id); err != nil {
		return errors.ErrInvalidInput.WithDetails(map[string]string{
			"id": "Task ID must be a valid UUID",
		})
	}

	// Simular que no encontramos la tarea algunas veces
	if id == "non-existent" {
		return errors.ErrTaskNotFound
	}

	// Retornar tarea mock
	task := dto.TaskResponse{
		ID:          id,
		Title:       "Learn Go",
		Description: "Study Go programming language",
		Priority:    "high",
		Completed:   false,
		CreatedAt:   time.Now().Add(-24 * time.Hour),
		UpdatedAt:   time.Now().Add(-24 * time.Hour),
	}

	return c.JSON(dto.SuccessResponse{
		Success: true,
		Message: "Task retrieved successfully",
		Data:    task,
	})
}

// CreateTask maneja POST /tasks
func (h *TaskHandler) CreateTask(c *fiber.Ctx) error {
	var req dto.CreateTaskRequest

	// Parse JSON body
	if err := c.BodyParser(&req); err != nil {
		return errors.ErrInvalidJSON
	}

	// Validar estructura
	if err := h.validator.ValidateStruct(&req); err != nil {
		return err // Ya es un AppError con detalles
	}

	// Crear tarea (mock)
	task := dto.TaskResponse{
		ID:          uuid.New().String(),
		Title:       req.Title,
		Description: req.Description,
		Priority:    req.Priority,
		Completed:   false,
		DueDate:     &req.DueDate,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Si priority está vacío, usar default
	if task.Priority == "" {
		task.Priority = "medium"
	}

	return c.Status(201).JSON(dto.SuccessResponse{
		Success: true,
		Message: "Task created successfully",
		Data:    task,
	})
}

// UpdateTask maneja PUT /tasks
func (h *TaskHandler) UpdateTask(c *fiber.Ctx) error {
	var req dto.UpdateTaskRequest
	id := c.Params("id")
	// Parse JSON body
	if err := c.BodyParser(&req); err != nil {
		return errors.ErrInvalidJSON
	}
	// Validar estructura
	if err := h.validator.ValidateStruct(&req); err != nil {
		return err // Ya es un AppError con detalles
	}
	// Validar que el ID no esté vacío
	if id == "" {
		return errors.ErrInvalidInput.WithDetails(map[string]string{
			"id": "Task ID is required",
		})
	}
	//validar que almenos un campo venga en el request
	if req.Title == nil && req.Description == nil && req.Priority == nil && req.DueDate == nil && req.Completed == nil {
		return errors.ErrInvalidInput.WithDetails(map[string]string{
			"body": "nothing to update",
		})
	}
	// Actualizar tarea (mock)
	task := dto.TaskResponse{
		ID:          id,
		Title:       *req.Title,
		Description: *req.Description,
		Priority:    *req.Description,
		Completed:   *req.Completed,
		CreatedAt:   time.Now().Add(-24 * time.Hour),
		UpdatedAt:   time.Now(),
	}
	return c.JSON(dto.SuccessResponse{
		Success: true,
		Message: "Task updated successfully",
		Data:    task,
	})
}

// DeleteTask maneja DELETE /tasks/:id
func (h *TaskHandler) DeleteTask(c *fiber.Ctx) error {
	id := c.Params("id")

	// Validar que el ID no esté vacío
	if id == "" {
		return errors.ErrInvalidInput.WithDetails(map[string]string{
			"id": "Task ID is required",
		})
	}

	// Validar formato UUID (opcional pero buena práctica)
	if _, err := uuid.Parse(id); err != nil {
		return errors.ErrInvalidInput.WithDetails(map[string]string{
			"id": "Task ID must be a valid UUID",
		})
	}

	// Simular que no encontramos la tarea algunas veces
	if id == "non-existent" {
		return errors.ErrTaskNotFound
	}

	return c.JSON(dto.SuccessResponse{
		Success: true,
		Message: "Task deleted successfully",
		Data: fiber.Map{
			"deleted_id": id,
		},
	})
}

// GetTask maneja GET /tasks/:id
func (h *TaskHandler) GetTask(c *fiber.Ctx) error {
	id := c.Params("id")

	// Validar UUID
	if _, err := uuid.Parse(id); err != nil {
		return errors.ErrInvalidInput.WithDetails(map[string]string{
			"id": "Task ID must be a valid UUID",
		})
	}

	// Simular búsqueda - algunas IDs retornan 404
	if id == "00000000-0000-0000-0000-000000000000" {
		return errors.ErrTaskNotFound
	}

	// Mock task response
	task := dto.TaskResponse{
		ID:          id,
		Title:       "Individual Task",
		Description: "Retrieved by ID",
		Priority:    "medium",
		Completed:   false,
		CreatedAt:   time.Now().Add(-2 * time.Hour),
		UpdatedAt:   time.Now().Add(-1 * time.Hour),
	}

	return c.JSON(task)
}

// UpdateTask maneja PUT /tasks/:id
func (h *TaskHandler) UpdateTask(c *fiber.Ctx) error {
	id := c.Params("id")

	// Validar UUID
	if _, err := uuid.Parse(id); err != nil {
		return errors.ErrInvalidInput.WithDetails(map[string]string{
			"id": "Task ID must be a valid UUID",
		})
	}

	var req dto.UpdateTaskRequest

	// Parse JSON body
	if err := c.BodyParser(&req); err != nil {
		return errors.ErrInvalidJSON
	}

	// Validar estructura
	if err := h.validator.ValidateStruct(&req); err != nil {
		return err
	}

	// Verificar que al menos un campo está presente
	if req.Title == nil && req.Description == nil && req.Priority == nil &&
		req.DueDate == nil && req.Completed == nil {
		return errors.ErrInvalidInput.WithDetails(map[string]string{
			"update": "At least one field must be provided to update",
		})
	}

	// Simular que algunas tareas no existen
	if id == "00000000-0000-0000-0000-000000000000" {
		return errors.ErrTaskNotFound
	}

	// Mock updated task
	task := dto.TaskResponse{
		ID:          id,
		Title:       "Updated Task",
		Description: "Task has been updated",
		Priority:    "high",
		Completed:   true,
		CreatedAt:   time.Now().Add(-24 * time.Hour),
		UpdatedAt:   time.Now(),
	}

	// Si se proporcionaron valores, usarlos
	if req.Title != nil {
		task.Title = *req.Title
	}
	if req.Description != nil {
		task.Description = *req.Description
	}
	if req.Priority != nil {
		task.Priority = *req.Priority
	}
	if req.Completed != nil {
		task.Completed = *req.Completed
	}

	return c.JSON(dto.SuccessResponse{
		Success: true,
		Message: "Task updated successfully",
		Data:    task,
	})
}
