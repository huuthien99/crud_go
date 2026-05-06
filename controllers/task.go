package controllers

import (
	"auth_crud/config"
	"auth_crud/models"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func CreateTask(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	var input models.TaskInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task := models.Task{
		Title:       input.Title,
		Description: input.Description,
		UserID:      userID,
		Status:      "pending",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	const taskFields = `
		id, title, description, status, user_id, created_at, updated_at
	`

	query := `
	INSERT INTO tasks(title, description, status, user_id, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING ` + taskFields

	err := config.DB.QueryRow(
		query,
		task.Title,
		task.Description,
		task.Status,
		task.UserID,
		task.CreatedAt,
		task.UpdatedAt,
	).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.Status,
		&task.UserID,
		&task.CreatedAt,
		&task.UpdatedAt,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Created successfully!", "task": task})

}

func GetTaskList(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	var query models.TaskQuery

	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	page := 1
	limit := 10

	if query.Page != nil && *query.Page > 1 {
		page = *query.Page
	}

	if query.Limit != nil && *query.Limit > 1 {
		limit = *query.Limit
	}

	offset := (page - 1) * limit

	where := "WHERE user_id = $1"

	args := []any{userID}

	argId := 2

	if query.Search != "" {
		where += " AND title ILIKE $" + strconv.Itoa(argId)
		args = append(args, "%"+query.Search+"%")
		argId++
	}

	if query.Status != "" {
		where += " AND status = $" + strconv.Itoa(argId)
		args = append(args, query.Status)
		argId++
	}

	countQuery := "SELECT COUNT(*) FROM tasks " + where

	var total int64
	err := config.DB.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count tasks"})
		return
	}

	dataQuery := `
		SELECT id, title, description, status, user_id, created_at, updated_at
		FROM tasks
		` + where + `
		ORDER BY created_at DESC
		LIMIT $` + strconv.Itoa(argId) + `
		OFFSET $` + strconv.Itoa(argId+1)

	args = append(args, limit, offset)

	rows, err := config.DB.Query(dataQuery, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get tasks"})
		return
	}

	defer rows.Close()

	var tasks []models.Task

	for rows.Next() {
		var task models.Task

		err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&task.Status,
			&task.UserID,
			&task.CreatedAt,
			&task.UpdatedAt,
		)
		if err != nil {
			log.Println("Scan stack fail", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error!"})
			return
		}

		tasks = append(tasks, task)
	}
	if err := rows.Err(); err != nil {
		log.Println("Row iteration error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": tasks,
		"pagination": gin.H{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": (total + int64(limit) - 1) / int64(limit),
		},
	})
}

func GetTaskHelper(userID uint, taskID uint) (models.Task, error) {

	var task models.Task

	query := `
		SELECT id, title, description, status, user_id, created_at, updated_at
		FROM tasks
		WHERE id = $1 AND user_id = $2
		`
	err := config.DB.QueryRow(
		query,
		taskID,
		userID,
	).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.Status,
		&task.UserID,
		&task.CreatedAt,
		&task.UpdatedAt,
	)

	if err != nil {
		return models.Task{}, err
	}
	return task, nil
}

func GetTaskById(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	task, err := GetTaskHelper(userID, uint(id))

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": task,
	})
}

func DeleteTask(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}
	taskId := uint(id)

	query := "DELETE FROM tasks WHERE id = $1 AND user_id = $2"

	result, err := config.DB.Exec(query, taskId, userID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task deleted"})
}

func UpdateTask(c *gin.Context) {

	userID := c.MustGet("user_id").(uint)

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	taskId := uint(id)

	task, err := GetTaskHelper(userID, taskId)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	var input models.TaskInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := `
		UPDATE tasks
		SET title = $1,
		    description = $2,
		    updated_at = NOW()
		WHERE id = $3 AND user_id = $4
		RETURNING id, title, description, status, user_id, created_at, updated_at
	`

	err = config.DB.QueryRow(
		query,
		input.Title,
		input.Description,
		taskId,
		userID,
	).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.Status,
		&task.UserID,
		&task.CreatedAt,
		&task.UpdatedAt,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update task"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "OK",
		"data":    task,
	})
}

func ChangeStatus(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	taskId := uint(id)

	task, err := GetTaskHelper(userID, taskId)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}
	
	var input models.TaskInputChangeStatus
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := `
		UPDATE tasks
		SET status = $1,
		    updated_at = NOW()
		WHERE id = $2 AND user_id = $3
		RETURNING id, title, description, status, user_id, created_at, updated_at
	`
	err = config.DB.QueryRow(
		query,
		input.Status,
		taskId,
		userID,
	).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.Status,
		&task.UserID,
		&task.CreatedAt,
		&task.UpdatedAt,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update task"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "OK",
		"data":    task,
	})
}
