package routes

import (
	"campus-events/controllers"
	"campus-events/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	r.Static("/uploads", "./uploads")

	api := r.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", controllers.Register)
			auth.POST("/login", controllers.Login)
		}

		events := api.Group("/events")
		{
			events.GET("", controllers.GetEventList)
			events.GET("/:id", controllers.GetEvent)
			events.GET("/:id/stats", controllers.GetEventStats)
			
			events.Use(middleware.Auth())
			{
				events.POST("", middleware.AdminAuth(), controllers.CreateEvent)
				events.PUT("/:id", middleware.AdminAuth(), controllers.UpdateEvent)
				events.DELETE("/:id", middleware.AdminAuth(), controllers.DeleteEvent)
				events.POST("/:id/end", middleware.AdminAuth(), controllers.EndEvent)
				events.GET("/:id/registrations", middleware.AdminAuth(), controllers.GetEventRegistrations)
				
				events.POST("/:id/register", controllers.RegisterForEvent)
				events.POST("/:id/cancel", controllers.CancelRegistration)
				events.POST("/:id/signin", controllers.SignIn)
			}
		}

		registrations := api.Group("/registrations")
		registrations.Use(middleware.Auth())
		{
			registrations.GET("/my", controllers.GetMyRegistrations)
		}

		upload := api.Group("/upload")
		upload.Use(middleware.Auth())
		{
			upload.POST("/avatar", controllers.UploadAvatar)
			upload.POST("/cover", middleware.AdminAuth(), controllers.UploadCover)
		}

		stats := api.Group("/stats")
		stats.Use(middleware.Auth(), middleware.AdminAuth())
		{
			stats.GET("/dashboard", controllers.GetDashboardStats)
			stats.GET("/export", controllers.ExportRegistrations)
		}

		user := api.Group("/user")
		user.Use(middleware.Auth())
		{
			user.GET("/me", controllers.GetCurrentUser)
			user.POST("/change-password", controllers.ChangePassword)
		}
	}
}
