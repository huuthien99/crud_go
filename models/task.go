package models

import (
	"time"
)

type Task struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	Title       string    `gorm:"size:255" json:"title"`
	Description string    `json:"description"`
	Status      string    `gorm:"default:pending;type:varchar(20)" json:"status"`
	UserID      uint      `json:"user_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type TaskInput struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	Status      string `json:"status" binding:"omitempty,oneof=pending processing completed"`
}

type TaskInputChangeStatus struct {
	Status string `json:"status" binding:"required,omitempty,oneof=pending processing completed"`
}

type TaskQuery struct {
	Page   int    `form:"page"`
	Limit  int    `form:"limit"`
	Search string `form:"search"`
	Status string `form:"status"`
}
