package controllers

import (
	"campus-events/config"
	"campus-events/models"
	"campus-events/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetEventComments(c *gin.Context) {
	eventID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(400, "无效的活动ID"))
		return
	}

	var comments []models.Comment
	config.DB.Where("event_id = ? AND parent_id IS NULL", eventID).
		Preload("User").
		Preload("Replies.User").
		Order("created_at DESC").
		Find(&comments)

	c.JSON(http.StatusOK, utils.SuccessResponse(comments))
}

func CreateComment(c *gin.Context) {
	eventID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(400, "无效的活动ID"))
		return
	}

	userID := c.GetUint("user_id")

	var req struct {
		Content  string `json:"content" binding:"required"`
		ParentID *uint  `json:"parent_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(400, "参数错误: "+err.Error()))
		return
	}

	comment := models.Comment{
		EventID:  uint(eventID),
		UserID:   userID,
		Content:  req.Content,
		ParentID: req.ParentID,
	}

	if err := config.DB.Create(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(500, "评论失败"))
		return
	}

	config.DB.Preload("User").First(&comment, comment.ID)

	c.JSON(http.StatusOK, utils.SuccessResponse(comment))
}

func DeleteComment(c *gin.Context) {
	commentID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(400, "无效的评论ID"))
		return
	}

	userID := c.GetUint("user_id")
	userRole := c.GetString("role")

	var comment models.Comment
	if err := config.DB.First(&comment, commentID).Error; err != nil {
		c.JSON(http.StatusNotFound, utils.ErrorResponse(404, "评论不存在"))
		return
	}

	if comment.UserID != userID && userRole != models.RoleAdmin {
		c.JSON(http.StatusForbidden, utils.ErrorResponse(403, "无权删除此评论"))
		return
	}

	config.DB.Delete(&comment)

	c.JSON(http.StatusOK, utils.SuccessResponse(nil))
}

func GetHotTags(c *gin.Context) {
	type TagCount struct {
		Name  string `json:"name"`
		Count int    `json:"count"`
	}

	var tags []TagCount
	config.DB.Model(&models.EventTag{}).
		Select("name, COUNT(*) as count").
		Group("name").
		Order("count DESC").
		Limit(10).
		Scan(&tags)

	c.JSON(http.StatusOK, utils.SuccessResponse(tags))
}
