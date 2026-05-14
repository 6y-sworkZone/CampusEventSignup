package controllers

import (
	"campus-events/config"
	"campus-events/models"
	"campus-events/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	role := c.Query("role")
	search := c.Query("search")

	offset := (page - 1) * pageSize

	query := config.DB.Model(&models.User{})

	if role != "" {
		query = query.Where("role = ?", role)
	}

	if search != "" {
		query = query.Where("username LIKE ? OR student_id LIKE ? OR employee_id LIKE ?", 
			"%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	var users []models.User
	var total int64

	query.Count(&total)
	query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&users)

	c.JSON(http.StatusOK, utils.SuccessResponse(gin.H{
		"list":  users,
		"total": total,
		"page":  page,
	}))
}

func ToggleUserStatus(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(400, "无效的用户ID"))
		return
	}

	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, utils.ErrorResponse(404, "用户不存在"))
		return
	}

	if user.Status == models.StatusActive {
		user.Status = models.StatusDisabled
	} else {
		user.Status = models.StatusActive
	}

	config.DB.Save(&user)

	c.JSON(http.StatusOK, utils.SuccessResponse(gin.H{
		"user":   user,
		"status": user.Status,
	}))
}
