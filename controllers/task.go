package controllers

import (
	"auth_crud/config"
	"auth_crud/models"
	"net/http"
	"strconv"

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
	}

	if err := config.DB.Create(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create todo"})
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

	if query.Page < 1 {
		query.Page = 1
	}

	if query.Limit < 1 {
		query.Limit = 10
	}

	offset := (query.Page - 1) * query.Limit

	var total int64

	var tasks []models.Task

	db := config.DB.Model(&models.Task{}).Where("user_id = ?", userID)

	if query.Search != "" {
		db = db.Where("title ILIKE ?", "%"+query.Search+"%")
	}

	if query.Status != "" {
		db = db.Where("status = ?", query.Status)
	}

	db.Count(&total)
	db.Order("created_at desc").Offset(offset).Limit(query.Limit).Find(&tasks)

	c.JSON(http.StatusOK, gin.H{
		"data": tasks,
		"pagination": gin.H{
			"page":        query.Page,
			"limit":       query.Limit,
			"total":       total,
			"total_pages": (total + int64(query.Limit) - 1) / int64(query.Limit),
		},
	})
}

func GetTaskHelper(c *gin.Context) (models.Task, bool) {
	userID := c.MustGet("user_id").(uint)
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return models.Task{}, false
	}

	var task models.Task

	if err := config.DB.Where("id = ? AND user_id = ?", id, userID).First(&task).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return models.Task{}, false
	}
	return task, true
}

func GetTaskById(c *gin.Context) {
	task, exist := GetTaskHelper(c)

	if !exist {
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": task,
	})
}

func DeleteTask(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	userID := c.MustGet("user_id").(uint)
	result := config.DB.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Task{})
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Task deleted"})
}

func UpdateTask(c *gin.Context) {
	task, exist := GetTaskHelper(c)

	if !exist {
		return
	}

	var input models.TaskInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// config.DB.Model(&task).Updates(input)
	// c.JSON(http.StatusOK, task)

	task.Title = input.Title
	task.Description = input.Description
	if input.Status != "" {
		task.Status = input.Status
	}

	if err := config.DB.Save(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update task"})
		return
	}

	c.JSON(http.StatusOK, task)
}

func ChangeStatus(c *gin.Context) {
	task, exist := GetTaskHelper(c)

	if !exist {
		return
	}

	var input models.TaskInputChangeStatus
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task.Status = input.Status

	config.DB.Save(&task)
	c.JSON(http.StatusOK, task)
}
