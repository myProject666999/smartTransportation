package controllers

import (
	"smart_transportation/database"
	"smart_transportation/models"
	"smart_transportation/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TicketController struct{}

func (tc *TicketController) ListRoutes(c *gin.Context) {
	var routes []models.Route
	var total int64

	page := 1
	pageSize := 10
	if p := c.Query("page"); p != "" {
		page = parseInt(p, 1)
	}
	if ps := c.Query("page_size"); ps != "" {
		pageSize = parseInt(ps, 10)
	}

	query := database.DB.Model(&models.Route{}).Where("status = ?", 1)

	keyword := c.Query("keyword")
	if keyword != "" {
		query = query.Where("route_name LIKE ? OR start_station LIKE ? OR end_station LIKE ?", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&routes)

	utils.Success(c, gin.H{
		"list":      routes,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func (tc *TicketController) GetRouteDetail(c *gin.Context) {
	routeID := c.Param("id")

	var route models.Route
	result := database.DB.First(&route, routeID)
	if result.Error != nil || route.ID == 0 {
		utils.NotFound(c, "路线不存在")
		return
	}

	utils.Success(c, route)
}

func (tc *TicketController) ListTickets(c *gin.Context) {
	var tickets []models.Ticket
	var total int64

	page := 1
	pageSize := 10
	if p := c.Query("page"); p != "" {
		page = parseInt(p, 1)
	}
	if ps := c.Query("page_size"); ps != "" {
		pageSize = parseInt(ps, 10)
	}

	query := database.DB.Model(&models.Ticket{}).Where("status = ?", 1)

	routeName := c.Query("route_name")
	if routeName != "" {
		query = query.Where("route_name LIKE ?", "%"+routeName+"%")
	}

	startStation := c.Query("start_station")
	if startStation != "" {
		query = query.Where("start_station = ?", startStation)
	}

	endStation := c.Query("end_station")
	if endStation != "" {
		query = query.Where("end_station = ?", endStation)
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	query.Order("departure_time ASC").Offset(offset).Limit(pageSize).Find(&tickets)

	utils.Success(c, gin.H{
		"list":      tickets,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func (tc *TicketController) GetTicketDetail(c *gin.Context) {
	ticketID := c.Param("id")

	var ticket models.Ticket
	result := database.DB.First(&ticket, ticketID)
	if result.Error != nil || ticket.ID == 0 {
		utils.NotFound(c, "车票不存在")
		return
	}

	utils.Success(c, ticket)
}

func (tc *TicketController) BookTicket(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req struct {
		TicketID       uint   `json:"ticket_id" binding:"required"`
		PassengerName  string `json:"passenger_name" binding:"required"`
		PassengerIDCard string `json:"passenger_id_card" binding:"required"`
		PassengerPhone string `json:"passenger_phone" binding:"required"`
		Quantity       int    `json:"quantity" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	if req.Quantity < 1 {
		utils.BadRequest(c, "购票数量至少为1")
		return
	}

	var ticket models.Ticket
	result := database.DB.First(&ticket, req.TicketID)
	if result.Error != nil || ticket.ID == 0 {
		utils.NotFound(c, "车票不存在")
		return
	}

	if ticket.Status != 1 {
		utils.BadRequest(c, "该车次已停运")
		return
	}

	if ticket.AvailableSeats < req.Quantity {
		utils.BadRequest(c, "余票不足")
		return
	}

	orderNo := generateOrderNo()

	booking := models.Booking{
		OrderNo:        orderNo,
		UserID:         userID.(uint),
		TicketID:       req.TicketID,
		PassengerName:  req.PassengerName,
		PassengerIDCard: req.PassengerIDCard,
		PassengerPhone: req.PassengerPhone,
		Quantity:       req.Quantity,
		TotalAmount:    ticket.Price * float64(req.Quantity),
		Status:         0,
	}

	tx := database.DB.Begin()

	if err := tx.Create(&booking).Error; err != nil {
		tx.Rollback()
		utils.InternalServerError(c, "购票失败: "+err.Error())
		return
	}

	if err := tx.Model(&ticket).Update("available_seats", ticket.AvailableSeats-req.Quantity).Error; err != nil {
		tx.Rollback()
		utils.InternalServerError(c, "更新余票失败")
		return
	}

	tx.Commit()

	utils.SuccessWithMessage(c, "购票成功，请支付", gin.H{
		"order_no":    orderNo,
		"total_amount": booking.TotalAmount,
		"booking_id":  booking.ID,
	})
}

func (tc *TicketController) PayBooking(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req struct {
		BookingID uint `json:"booking_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	var booking models.Booking
	result := database.DB.Where("id = ? AND user_id = ?", req.BookingID, userID).First(&booking)
	if result.Error != nil || booking.ID == 0 {
		utils.NotFound(c, "订单不存在")
		return
	}

	if booking.Status != 0 {
		utils.BadRequest(c, "订单状态不正确")
		return
	}

	booking.Status = 1
	booking.PaymentTime = time.Now()

	if err := database.DB.Save(&booking).Error; err != nil {
		utils.InternalServerError(c, "支付失败")
		return
	}

	utils.SuccessWithMessage(c, "支付成功", gin.H{
		"order_no": booking.OrderNo,
		"status":   booking.Status,
	})
}

func (tc *TicketController) ListBookings(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var bookings []models.Booking
	var total int64

	page := 1
	pageSize := 10
	if p := c.Query("page"); p != "" {
		page = parseInt(p, 1)
	}
	if ps := c.Query("page_size"); ps != "" {
		pageSize = parseInt(ps, 10)
	}

	status := c.Query("status")
	query := database.DB.Model(&models.Booking{}).Where("user_id = ?", userID)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&bookings)

	utils.Success(c, gin.H{
		"list":      bookings,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func (tc *TicketController) CancelBooking(c *gin.Context) {
	userID, _ := c.Get("user_id")
	bookingID := c.Param("id")

	var booking models.Booking
	result := database.DB.Where("id = ? AND user_id = ?", bookingID, userID).First(&booking)
	if result.Error != nil || booking.ID == 0 {
		utils.NotFound(c, "订单不存在")
		return
	}

	if booking.Status == 2 {
		utils.BadRequest(c, "订单已取消")
		return
	}

	if booking.Status == 1 {
		var ticket models.Ticket
		database.DB.First(&ticket, booking.TicketID)
		if ticket.ID != 0 {
			database.DB.Model(&ticket).Update("available_seats", ticket.AvailableSeats+booking.Quantity)
		}
	}

	booking.Status = 2
	database.DB.Save(&booking)

	utils.SuccessWithMessage(c, "订单已取消", nil)
}

func generateOrderNo() string {
	return "ORD" + time.Now().Format("20060102150405") + uuid.New().String()[:8]
}
