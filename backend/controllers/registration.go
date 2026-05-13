package controllers

import (
	"campus-events/config"
	"campus-events/models"
	"campus-events/utils"
	"net/http"
	"strconv"
	"time"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterForEvent(c *gin.Context) {
	userID := c.GetUint("user_id")
	eventID, _ := strconv.Atoi(c.Param("id"))

	tx := config.DB.Begin()

	var event models.Event
	if err := tx.Set("gorm:query_option", "FOR UPDATE").First(&event, eventID).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusNotFound, utils.ErrorResponse(404, "活动不存在"))
		return
	}

	if event.Status != models.EventStatusOpen {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(400, "活动不可报名"))
		return
	}

	if time.Now().After(event.RegistrationEnd) {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(400, "报名已截止"))
		return
	}

	var existingReg models.Registration
	tx.Where("user_id = ? AND event_id = ?", userID, eventID).First(&existingReg)
	if existingReg.ID != 0 && existingReg.Status != models.RegStatusCancelled {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(400, "您已报名该活动"))
		return
	}

	maxWaitlist := int(float64(event.MaxParticipants) * 0.2)
	isWaitlist := false
	waitlistNum := 0

	if event.CurrentCount >= event.MaxParticipants {
		if event.WaitlistCount >= maxWaitlist {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, utils.ErrorResponse(400, "报名已满，候补也已满"))
			return
		}
		isWaitlist = true
		waitlistNum = event.WaitlistCount + 1
		event.WaitlistCount++
	} else {
		event.CurrentCount++
	}

	if existingReg.ID != 0 {
		existingReg.Status = models.RegStatusRegistered
		existingReg.IsWaitlist = isWaitlist
		existingReg.WaitlistNum = waitlistNum
		tx.Save(&existingReg)
	} else {
		reg := models.Registration{
			UserID:      userID,
			EventID:     uint(eventID),
			Status:      models.RegStatusRegistered,
			IsWaitlist:  isWaitlist,
			WaitlistNum: waitlistNum,
		}
		tx.Create(&reg)
	}
	tx.Save(&event)
	tx.Commit()

	message := "报名成功"
	if isWaitlist {
		message = "已加入候补名单，顺位: " + strconv.Itoa(waitlistNum)
	}
	c.JSON(http.StatusOK, utils.SuccessResponse(gin.H{"message": message, "is_waitlist": isWaitlist}))
}

func CancelRegistration(c *gin.Context) {
	userID := c.GetUint("user_id")
	eventID, _ := strconv.Atoi(c.Param("id"))

	tx := config.DB.Begin()

	var reg models.Registration
	if err := tx.Where("user_id = ? AND event_id = ?", userID, eventID).First(&reg).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusNotFound, utils.ErrorResponse(404, "未找到报名记录"))
		return
	}

	if reg.Status == models.RegStatusCancelled {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(400, "报名已取消"))
		return
	}

	var event models.Event
	tx.First(&event, eventID)

	if time.Until(event.StartTime) < 2*time.Hour {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(400, "活动开始前2小时内不可取消报名"))
		return
	}

	reg.Status = models.RegStatusCancelled
	tx.Save(&reg)

	if !reg.IsWaitlist {
		event.CurrentCount--

		var firstWaitlist models.Registration
		tx.Where("event_id = ? AND is_waitlist = ? AND status = ?", 
			eventID, true, models.RegStatusRegistered).
			Order("waitlist_num ASC").First(&firstWaitlist)

		if firstWaitlist.ID != 0 {
			firstWaitlist.IsWaitlist = false
			firstWaitlist.WaitlistNum = 0
			tx.Save(&firstWaitlist)
			event.CurrentCount++
			event.WaitlistCount--

			tx.Model(&models.Registration{}).
				Where("event_id = ? AND is_waitlist = ? AND waitlist_num > ?", 
					eventID, true, firstWaitlist.WaitlistNum).
				UpdateColumn("waitlist_num", gorm.Expr("waitlist_num - 1"))
		}
	} else {
		event.WaitlistCount--
		tx.Model(&models.Registration{}).
			Where("event_id = ? AND is_waitlist = ? AND waitlist_num > ?", 
				eventID, true, reg.WaitlistNum).
			UpdateColumn("waitlist_num", gorm.Expr("waitlist_num - 1"))
	}

	tx.Save(&event)
	tx.Commit()

	c.JSON(http.StatusOK, utils.SuccessResponse(nil))
}

func SignIn(c *gin.Context) {
	userID := c.GetUint("user_id")
	eventID, _ := strconv.Atoi(c.Param("id"))
	var req struct {
		SignCode string `json:"sign_code" binding:"required,len=4"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(400, "参数错误"))
		return
	}

	var event models.Event
	if err := config.DB.First(&event, eventID).Error; err != nil {
		c.JSON(http.StatusNotFound, utils.ErrorResponse(404, "活动不存在"))
		return
	}

	if time.Until(event.StartTime) > 30*time.Minute {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(400, "签到未开始（活动开始前30分钟开放）"))
		return
	}

	if event.SignCode != req.SignCode {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(400, "签到码错误"))
		return
	}

	var reg models.Registration
	if err := config.DB.Where("user_id = ? AND event_id = ?", userID, eventID).First(&reg).Error; err != nil {
		c.JSON(http.StatusNotFound, utils.ErrorResponse(404, "未找到报名记录"))
		return
	}

	if reg.IsWaitlist {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(400, "候补用户不能签到"))
		return
	}

	if reg.Status == models.RegStatusCancelled {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(400, "报名已取消"))
		return
	}

	if reg.Signed {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(400, "已签到"))
		return
	}

	reg.Signed = true
	reg.SignedAt = time.Now()
	reg.Status = models.RegStatusSigned
	config.DB.Save(&reg)

	c.JSON(http.StatusOK, utils.SuccessResponse(nil))
}

func GetMyRegistrations(c *gin.Context) {
	userID := c.GetUint("user_id")
	status := c.Query("status")

	query := config.DB.Model(&models.Registration{}).
		Preload("Event").
		Preload("Event.Tags").
		Where("user_id = ?", userID)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	var regs []models.Registration
	query.Order("created_at DESC").Find(&regs)

	c.JSON(http.StatusOK, utils.SuccessResponse(regs))
}

func GetEventRegistrations(c *gin.Context) {
	eventID, _ := strconv.Atoi(c.Param("id"))

	var regs []models.Registration
	config.DB.Preload("User").
		Where("event_id = ?", eventID).
		Order("created_at DESC").
		Find(&regs)

	c.JSON(http.StatusOK, utils.SuccessResponse(regs))
}
