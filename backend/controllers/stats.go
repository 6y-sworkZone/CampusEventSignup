package controllers

import (
	"campus-events/config"
	"campus-events/models"
	"campus-events/utils"
	"fmt"
	"net/http"
	"strconv"
	"time"
	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

func GetDashboardStats(c *gin.Context) {
	var totalEvents, totalUsers, totalRegistrations int64
	
	config.DB.Model(&models.Event{}).Count(&totalEvents)
	config.DB.Model(&models.User{}).Count(&totalUsers)
	config.DB.Model(&models.Registration{}).Where("status != ?", models.RegStatusCancelled).Count(&totalRegistrations)

	type CollegeStats struct {
		College string `json:"college"`
		Count   int64  `json:"count"`
	}
	var collegeStats []CollegeStats
	config.DB.Model(&models.User{}).
		Select("college, count(*) as count").
		Group("college").
		Scan(&collegeStats)

	type DailyStats struct {
		Date  string `json:"date"`
		Count int64  `json:"count"`
	}
	var dailyStats []DailyStats
	
	now := time.Now()
	for i := 6; i >= 0; i-- {
		date := now.AddDate(0, 0, -i)
		dateStr := date.Format("2006-01-02")
		
		var count int64
		config.DB.Model(&models.Registration{}).
			Where("DATE(created_at) = ? AND status != ?", dateStr, models.RegStatusCancelled).
			Count(&count)
		
		dailyStats = append(dailyStats, DailyStats{
			Date:  dateStr,
			Count: count,
		})
	}

	c.JSON(http.StatusOK, utils.SuccessResponse(gin.H{
		"total_events":       totalEvents,
		"total_users":        totalUsers,
		"total_registrations": totalRegistrations,
		"college_stats":      collegeStats,
		"daily_stats":        dailyStats,
	}))
}

func ExportRegistrations(c *gin.Context) {
	eventID := c.Query("event_id")
	
	var regs []models.Registration
	query := config.DB.Model(&models.Registration{}).Preload("User").Preload("Event")
	
	if eventID != "" {
		query = query.Where("event_id = ?", eventID)
	}
	
	query.Where("status != ?", models.RegStatusCancelled).Find(&regs)

	f := excelize.NewFile()
	sheetName := "报名名单"
	f.SetSheetName("Sheet1", sheetName)

	headers := []string{"序号", "活动名称", "姓名", "学号/工号", "学院", "手机号", "报名时间", "是否签到", "签到时间"}
	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheetName, cell, header)
	}

	for i, reg := range regs {
		row := i + 2
		studentID := reg.User.StudentID
		if studentID == "" {
			studentID = reg.User.EmployeeID
		}
		
		signed := "否"
		signedAt := ""
		if reg.Signed {
			signed = "是"
			signedAt = reg.SignedAt.Format("2006-01-02 15:04:05")
		}

		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), i+1)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), reg.Event.Title)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), reg.User.Username)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), studentID)
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", row), reg.User.College)
		f.SetCellValue(sheetName, fmt.Sprintf("F%d", row), reg.User.Phone)
		f.SetCellValue(sheetName, fmt.Sprintf("G%d", row), reg.CreatedAt.Format("2006-01-02 15:04:05"))
		f.SetCellValue(sheetName, fmt.Sprintf("H%d", row), signed)
		f.SetCellValue(sheetName, fmt.Sprintf("I%d", row), signedAt)
	}

	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename=registrations.xlsx")
	f.Write(c.Writer)
}

func GetEventStats(c *gin.Context) {
	eventID := c.Param("id")
	
	var event models.Event
	config.DB.First(&event, eventID)

	var registeredCount, signedCount, waitlistCount int64
	config.DB.Model(&models.Registration{}).
		Where("event_id = ? AND status = ? AND is_waitlist = ?", eventID, models.RegStatusRegistered, false).
		Count(&registeredCount)
	
	config.DB.Model(&models.Registration{}).
		Where("event_id = ? AND signed = ?", eventID, true).
		Count(&signedCount)
	
	config.DB.Model(&models.Registration{}).
		Where("event_id = ? AND is_waitlist = ? AND status = ?", eventID, true, models.RegStatusRegistered).
		Count(&waitlistCount)

	c.JSON(http.StatusOK, utils.SuccessResponse(gin.H{
		"event":            event,
		"max_participants": event.MaxParticipants,
		"registered_count": registeredCount,
		"signed_count":     signedCount,
		"waitlist_count":   waitlistCount,
	}))
}
