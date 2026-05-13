package controllers

import (
	"campus-events/config"
	"campus-events/models"
	"campus-events/utils"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"
	"github.com/gin-gonic/gin"
)

type CreateEventRequest struct {
	Title           string   `json:"title" binding:"required,min=2,max=50"`
	Description     string   `json:"description"`
	CoverImage      string   `json:"cover_image"`
	StartTime       string   `json:"start_time" binding:"required"`
	EndTime         string   `json:"end_time" binding:"required"`
	RegistrationEnd string   `json:"registration_end" binding:"required"`
	MaxParticipants int      `json:"max_participants" binding:"required,min=1"`
	Tags            []string `json:"tags"`
}

func CreateEvent(c *gin.Context) {
	var req CreateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(400, "参数错误: "+err.Error()))
		return
	}

	startTime, _ := time.Parse(time.RFC3339, req.StartTime)
	endTime, _ := time.Parse(time.RFC3339, req.EndTime)
	registrationEnd, _ := time.Parse(time.RFC3339, req.RegistrationEnd)

	if !endTime.After(startTime) {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(400, "结束时间必须晚于开始时间"))
		return
	}
	if !registrationEnd.Before(startTime) {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(400, "报名截止时间必须早于活动开始时间"))
		return
	}

	signCode := fmt.Sprintf("%04d", rand.Intn(10000))

	userID := c.GetUint("user_id")
	event := models.Event{
		Title:           req.Title,
		Description:     req.Description,
		CoverImage:      req.CoverImage,
		StartTime:       startTime,
		EndTime:         endTime,
		RegistrationEnd: registrationEnd,
		MaxParticipants: req.MaxParticipants,
		Status:          models.EventStatusOpen,
		SignCode:        signCode,
		CreatorID:       userID,
	}

	tx := config.DB.Begin()
	if err := tx.Create(&event).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(500, "创建活动失败"))
		return
	}

	for _, tagName := range req.Tags {
		tag := models.EventTag{
			EventID: event.ID,
			Name:    tagName,
		}
		tx.Create(&tag)
	}
	tx.Commit()

	c.JSON(http.StatusOK, utils.SuccessResponse(event))
}

func UpdateEvent(c *gin.Context) {
	eventID, _ := strconv.Atoi(c.Param("id"))
	
	var event models.Event
	if err := config.DB.First(&event, eventID).Error; err != nil {
		c.JSON(http.StatusNotFound, utils.ErrorResponse(404, "活动不存在"))
		return
	}

	var req CreateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(400, "参数错误"))
		return
	}

	startTime, _ := time.Parse(time.RFC3339, req.StartTime)
	endTime, _ := time.Parse(time.RFC3339, req.EndTime)
	registrationEnd, _ := time.Parse(time.RFC3339, req.RegistrationEnd)

	if !endTime.After(startTime) {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(400, "结束时间必须晚于开始时间"))
		return
	}

	if req.MaxParticipants < event.CurrentCount {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(400, "最大人数不能小于当前报名人数"))
		return
	}

	event.Title = req.Title
	event.Description = req.Description
	event.CoverImage = req.CoverImage
	event.StartTime = startTime
	event.EndTime = endTime
	event.RegistrationEnd = registrationEnd
	event.MaxParticipants = req.MaxParticipants

	tx := config.DB.Begin()
	tx.Where("event_id = ?", event.ID).Delete(&models.EventTag{})
	for _, tagName := range req.Tags {
		tag := models.EventTag{
			EventID: event.ID,
			Name:    tagName,
		}
		tx.Create(&tag)
	}
	tx.Save(&event)
	tx.Commit()

	c.JSON(http.StatusOK, utils.SuccessResponse(event))
}

func GetEvent(c *gin.Context) {
	eventID, _ := strconv.Atoi(c.Param("id"))
	
	var event models.Event
	if err := config.DB.Preload("Tags").First(&event, eventID).Error; err != nil {
		c.JSON(http.StatusNotFound, utils.ErrorResponse(404, "活动不存在"))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse(event))
}

func GetEventList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	status := c.Query("status")
	keyword := c.Query("keyword")
	tag := c.Query("tag")

	query := config.DB.Model(&models.Event{}).Preload("Tags")

	if status != "" {
		query = query.Where("status = ?", status)
	}
	if keyword != "" {
		query = query.Where("title LIKE ?", "%"+keyword+"%")
	}
	if tag != "" {
		subQuery := config.DB.Model(&models.EventTag{}).Select("event_id").Where("name = ?", tag)
		query = query.Where("id IN (?)", subQuery)
	}

	var total int64
	query.Count(&total)

	var events []models.Event
	offset := (page - 1) * pageSize
	query.Order("created_at DESC").Limit(pageSize).Offset(offset).Find(&events)

	c.JSON(http.StatusOK, utils.SuccessResponse(gin.H{
		"list":  events,
		"total": total,
		"page":  page,
	}))
}

func DeleteEvent(c *gin.Context) {
	eventID, _ := strconv.Atoi(c.Param("id"))
	
	tx := config.DB.Begin()
	tx.Where("event_id = ?", eventID).Delete(&models.EventTag{})
	tx.Where("event_id = ?", eventID).Delete(&models.Registration{})
	tx.Delete(&models.Event{}, eventID)
	tx.Commit()

	c.JSON(http.StatusOK, utils.SuccessResponse(nil))
}

func EndEvent(c *gin.Context) {
	eventID, _ := strconv.Atoi(c.Param("id"))
	
	var event models.Event
	if err := config.DB.First(&event, eventID).Error; err != nil {
		c.JSON(http.StatusNotFound, utils.ErrorResponse(404, "活动不存在"))
		return
	}

	event.Status = models.EventStatusEnded
	config.DB.Save(&event)

	c.JSON(http.StatusOK, utils.SuccessResponse(event))
}
