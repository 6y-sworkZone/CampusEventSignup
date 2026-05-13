package controllers

import (
	"campus-events/config"
	"campus-events/models"
	"campus-events/utils"
	"fmt"
	"net/http"
	"path/filepath"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func UploadAvatar(c *gin.Context) {
	userID := c.GetUint("user_id")

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(400, "获取文件失败"))
		return
	}

	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	uploadPath := fmt.Sprintf("./uploads/avatars/%s", filename)

	if err := c.SaveUploadedFile(file, uploadPath); err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(500, "保存文件失败"))
		return
	}

	var user models.User
	config.DB.First(&user, userID)
	user.Avatar = "/uploads/avatars/" + filename
	config.DB.Save(&user)

	c.JSON(http.StatusOK, utils.SuccessResponse(gin.H{"url": user.Avatar}))
}

func UploadCover(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(400, "获取文件失败"))
		return
	}

	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	uploadPath := fmt.Sprintf("./uploads/covers/%s", filename)

	if err := c.SaveUploadedFile(file, uploadPath); err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(500, "保存文件失败"))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse(gin.H{"url": "/uploads/covers/" + filename}))
}
