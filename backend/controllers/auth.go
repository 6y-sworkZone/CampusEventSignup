package controllers

import (
	"campus-events/config"
	"campus-events/models"
	"campus-events/utils"
	"net/http"
	"time"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
	Username   string `json:"username" binding:"required,min=2,max=50"`
	StudentID  string `json:"student_id"`
	EmployeeID string `json:"employee_id"`
	Phone      string `json:"phone" binding:"required,len=11"`
	College    string `json:"college" binding:"required"`
	Password   string `json:"password" binding:"required,min=6"`
	Role       string `json:"role" binding:"required,oneof=student admin"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	RememberMe bool `json:"remember_me"`
}

func Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(400, "参数错误: "+err.Error()))
		return
	}

	if req.Role == "student" && req.StudentID == "" {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(400, "学生必须提供学号"))
		return
	}
	if req.Role == "admin" && req.EmployeeID == "" {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(400, "管理员必须提供工号"))
		return
	}

	var existingUser models.User
	if req.StudentID != "" {
		config.DB.Where("student_id = ?", req.StudentID).First(&existingUser)
		if existingUser.ID != 0 {
			c.JSON(http.StatusBadRequest, utils.ErrorResponse(400, "学号已存在"))
			return
		}
	}
	if req.EmployeeID != "" {
		config.DB.Where("employee_id = ?", req.EmployeeID).First(&existingUser)
		if existingUser.ID != 0 {
			c.JSON(http.StatusBadRequest, utils.ErrorResponse(400, "工号已存在"))
			return
		}
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(500, "密码加密失败"))
		return
	}

	user := models.User{
		Username:   req.Username,
		StudentID:  req.StudentID,
		EmployeeID: req.EmployeeID,
		Phone:      req.Phone,
		College:    req.College,
		Password:   string(hashedPassword),
		Role:       req.Role,
		Avatar:     "",
	}

	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(500, "注册失败"))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse(user))
}

func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(400, "参数错误"))
		return
	}

	var user models.User
	config.DB.Where("username = ? OR student_id = ? OR employee_id = ? OR phone = ?", 
		req.Username, req.Username, req.Username, req.Username).First(&user)

	if user.ID == 0 {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse(401, "用户不存在"))
		return
	}

	if time.Now().Before(user.LockedUntil) {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse(401, "账户已锁定，请5分钟后再试"))
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		user.LoginAttempts++
		if user.LoginAttempts >= 3 {
			user.LockedUntil = time.Now().Add(5 * time.Minute)
			user.LoginAttempts = 0
			config.DB.Save(&user)
			c.JSON(http.StatusUnauthorized, utils.ErrorResponse(401, "密码错误3次，账户已锁定5分钟"))
			return
		}
		config.DB.Save(&user)
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse(401, "密码错误"))
		return
	}

	user.LoginAttempts = 0
	user.LockedUntil = time.Time{}
	config.DB.Save(&user)

	expireHours := 24
	if req.RememberMe {
		expireHours = 24 * 7
	}

	token, err := utils.GenerateToken(user.ID, user.Role, expireHours)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(500, "生成token失败"))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse(gin.H{
		"token": token,
		"user":  user,
	}))
}

func GetCurrentUser(c *gin.Context) {
	userID := c.GetUint("user_id")
	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, utils.ErrorResponse(404, "用户不存在"))
		return
	}
	c.JSON(http.StatusOK, utils.SuccessResponse(user))
}

func ChangePassword(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required,min=6"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(400, "参数错误"))
		return
	}

	var user models.User
	config.DB.First(&user, userID)

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword)); err != nil {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse(401, "原密码错误"))
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	user.Password = string(hashedPassword)
	config.DB.Save(&user)

	c.JSON(http.StatusOK, utils.SuccessResponse(nil))
}
