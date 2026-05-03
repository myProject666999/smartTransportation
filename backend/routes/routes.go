package routes

import (
	"smart_transportation/controllers"
	"smart_transportation/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	router.Use(middleware.Cors())
	router.Use(middleware.NoCache())
	router.Use(middleware.Secure())

	api := router.Group("/api")
	{
		userController := &controllers.UserController{}
		adminController := &controllers.AdminController{}
		ticketController := &controllers.TicketController{}
		newsController := &controllers.NewsController{}
		announcementController := &controllers.AnnouncementController{}
		carouselController := &controllers.CarouselController{}
		vehicleController := &controllers.VehicleController{}
		vehicleMonitorController := &controllers.VehicleMonitorController{}
		roadConditionController := &controllers.RoadConditionController{}
		violationController := &controllers.ViolationController{}
		feedbackController := &controllers.FeedbackController{}
		trafficFlowController := &controllers.TrafficFlowController{}
		signalLightController := &controllers.SignalLightController{}
		bookingAdminController := &controllers.BookingAdminController{}
		ticketAdminController := &controllers.TicketAdminController{}
		routeAdminController := &controllers.RouteAdminController{}

		api.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "pong",
			})
		})

		user := api.Group("/user")
		{
			user.POST("/register", userController.Register)
			user.POST("/login", userController.Login)

			auth := user.Group("")
			auth.Use(middleware.JWTAuth(), middleware.UserAuth())
			{
				auth.GET("/profile", userController.GetProfile)
				auth.PUT("/profile", userController.UpdateProfile)
				auth.PUT("/password", userController.ChangePassword)

				auth.GET("/routes", ticketController.ListRoutes)
				auth.GET("/routes/:id", ticketController.GetRouteDetail)

				auth.GET("/tickets", ticketController.ListTickets)
				auth.GET("/tickets/:id", ticketController.GetTicketDetail)
				auth.POST("/tickets/book", ticketController.BookTicket)
				auth.POST("/tickets/pay", ticketController.PayBooking)

				auth.GET("/bookings", ticketController.ListBookings)
				auth.POST("/bookings/:id/cancel", ticketController.CancelBooking)

				auth.GET("/vehicle-monitor", vehicleMonitorController.List)
				auth.GET("/vehicle-monitor/plate/:plate_number", vehicleMonitorController.GetByPlateNumber)

				auth.GET("/road-conditions", roadConditionController.List)
				auth.GET("/road-conditions/:id", roadConditionController.GetDetail)

				auth.GET("/violations", violationController.ListUserViolations)
				auth.GET("/violations/:id", violationController.GetDetail)
				auth.POST("/violations/:id/pay", violationController.PayViolation)

				auth.GET("/feedbacks", feedbackController.ListUserFeedbacks)
				auth.GET("/feedbacks/:id", feedbackController.GetDetail)
				auth.POST("/feedbacks", feedbackController.Create)
			}
		}

		admin := api.Group("/admin")
		{
			admin.POST("/login", adminController.Login)

			auth := admin.Group("")
			auth.Use(middleware.JWTAuth(), middleware.AdminAuth())
			{
				auth.GET("/profile", adminController.GetProfile)
				auth.PUT("/password", adminController.ChangePassword)

				auth.GET("/users", adminController.ListUsers)
				auth.PUT("/users/:id/status", adminController.UpdateUserStatus)

				auth.GET("/admins", adminController.ListAdmins)
				auth.POST("/admins", adminController.CreateAdmin)

				auth.GET("/routes", routeAdminController.List)
				auth.GET("/routes/:id", routeAdminController.GetDetail)
				auth.POST("/routes", routeAdminController.Create)
				auth.PUT("/routes/:id", routeAdminController.Update)
				auth.DELETE("/routes/:id", routeAdminController.Delete)

				auth.GET("/tickets", ticketAdminController.List)
				auth.GET("/tickets/:id", ticketAdminController.GetDetail)
				auth.POST("/tickets", ticketAdminController.Create)
				auth.PUT("/tickets/:id", ticketAdminController.Update)
				auth.DELETE("/tickets/:id", ticketAdminController.Delete)

				auth.GET("/bookings", bookingAdminController.List)
				auth.GET("/bookings/:id", bookingAdminController.GetDetail)
				auth.POST("/bookings/:id/cancel", bookingAdminController.Cancel)

				auth.GET("/vehicles", vehicleController.ListVehicles)
				auth.POST("/vehicles", vehicleController.CreateVehicle)
				auth.PUT("/vehicles/:id", vehicleController.UpdateVehicle)
				auth.DELETE("/vehicles/:id", vehicleController.DeleteVehicle)

				auth.GET("/vehicle-monitors", vehicleMonitorController.List)
				auth.GET("/vehicle-monitors/:id", vehicleMonitorController.GetByPlateNumber)
				auth.POST("/vehicle-monitors", vehicleMonitorController.Create)
				auth.PUT("/vehicle-monitors/:id", vehicleMonitorController.Update)
				auth.DELETE("/vehicle-monitors/:id", vehicleMonitorController.Delete)

				auth.GET("/road-conditions", roadConditionController.List)
				auth.GET("/road-conditions/:id", roadConditionController.GetDetail)
				auth.POST("/road-conditions", roadConditionController.Create)
				auth.PUT("/road-conditions/:id", roadConditionController.Update)
				auth.DELETE("/road-conditions/:id", roadConditionController.Delete)

				auth.GET("/violations", violationController.ListAllViolations)
				auth.GET("/violations/:id", violationController.GetDetail)
				auth.POST("/violations", violationController.Create)
				auth.PUT("/violations/:id", violationController.Update)
				auth.DELETE("/violations/:id", violationController.Delete)

				auth.GET("/feedbacks", feedbackController.ListAllFeedbacks)
				auth.GET("/feedbacks/:id", feedbackController.GetDetail)
				auth.PUT("/feedbacks/:id/reply", feedbackController.Reply)
				auth.PUT("/feedbacks/:id", feedbackController.Update)
				auth.DELETE("/feedbacks/:id", feedbackController.Delete)

				auth.GET("/traffic-flows", trafficFlowController.List)
				auth.GET("/traffic-flows/:id", trafficFlowController.GetDetail)
				auth.POST("/traffic-flows", trafficFlowController.Create)
				auth.PUT("/traffic-flows/:id", trafficFlowController.Update)
				auth.DELETE("/traffic-flows/:id", trafficFlowController.Delete)

				auth.GET("/signal-lights", signalLightController.List)
				auth.GET("/signal-lights/:id", signalLightController.GetDetail)
				auth.POST("/signal-lights", signalLightController.Create)
				auth.PUT("/signal-lights/:id", signalLightController.Update)
				auth.PUT("/signal-lights/:id/control", signalLightController.Control)
				auth.DELETE("/signal-lights/:id", signalLightController.Delete)

				auth.GET("/news", newsController.List)
				auth.GET("/news/:id", newsController.GetDetail)
				auth.POST("/news", newsController.Create)
				auth.PUT("/news/:id", newsController.Update)
				auth.DELETE("/news/:id", newsController.Delete)

				auth.GET("/announcements", announcementController.List)
				auth.GET("/announcements/:id", announcementController.GetDetail)
				auth.POST("/announcements", announcementController.Create)
				auth.PUT("/announcements/:id", announcementController.Update)
				auth.DELETE("/announcements/:id", announcementController.Delete)

				auth.GET("/carousels", carouselController.List)
				auth.POST("/carousels", carouselController.Create)
				auth.PUT("/carousels/:id", carouselController.Update)
				auth.DELETE("/carousels/:id", carouselController.Delete)
			}
		}

		public := api.Group("/public")
		{
			public.GET("/news", newsController.List)
			public.GET("/news/:id", newsController.GetDetail)
			public.GET("/announcements", announcementController.List)
			public.GET("/announcements/:id", announcementController.GetDetail)
			public.GET("/carousels", carouselController.GetAll)
			public.GET("/routes", ticketController.ListRoutes)
			public.GET("/routes/:id", ticketController.GetRouteDetail)
			public.GET("/road-conditions", roadConditionController.List)
		}
	}
}
