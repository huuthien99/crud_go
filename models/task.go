package models

import (
	"time"
)

type Task struct {
	ID          uint      `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	UserID      uint      `json:"user_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type TaskInput struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
}

type TaskInputChangeStatus struct {
	Status string `json:"status" binding:"required,oneof=pending processing completed"`
}

type TaskQuery struct {
	Page   *int   `form:"page"`
	Limit  *int   `form:"limit"`
	Search string `form:"search"`
	Status string `form:"status"`
}
