package controllers

import (
	"smart_transportation/database"
	"smart_transportation/models"
	"smart_transportation/utils"
	"time"

	"github.com/gin-gonic/gin"
)

type TrafficFlowController struct{}

func (tfc *TrafficFlowController) List(c *gin.Context) {
	var trafficFlows []models.TrafficFlow
	var total int64

	page := 1
	pageSize := 10
	if p := c.Query("page"); p != "" {
		page = parseInt(p, 1)
	}
	if ps := c.Query("page_size"); ps != "" {
		pageSize = parseInt(ps, 10)
	}

	query := database.DB.Model(&models.TrafficFlow{})

	keyword := c.Query("keyword")
	if keyword != "" {
		query = query.Where("road_name LIKE ? OR location LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	flowLevel := c.Query("flow_level")
	if flowLevel != "" {
		query = query.Where("flow_level = ?", flowLevel)
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	query.Order("record_time DESC").Offset(offset).Limit(pageSize).Find(&trafficFlows)

	utils.Success(c, gin.H{
		"list":      trafficFlows,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func (tfc *TrafficFlowController) GetDetail(c *gin.Context) {
	trafficFlowID := c.Param("id")

	var trafficFlow models.TrafficFlow
	result := database.DB.First(&trafficFlow, trafficFlowID)
	if result.Error != nil || trafficFlow.ID == 0 {
		utils.NotFound(c, "交通流量记录不存在")
		return
	}

	utils.Success(c, trafficFlow)
}

func (tfc *TrafficFlowController) Create(c *gin.Context) {
	var trafficFlow models.TrafficFlow
	if err := c.ShouldBindJSON(&trafficFlow); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	if trafficFlow.RecordTime.IsZero() {
		trafficFlow.RecordTime = time.Now()
	}

	if err := database.DB.Create(&trafficFlow).Error; err != nil {
		utils.InternalServerError(c, "创建失败: "+err.Error())
		return
	}

	utils.SuccessWithMessage(c, "创建成功", gin.H{"id": trafficFlow.ID})
}

func (tfc *TrafficFlowController) Update(c *gin.Context) {
	trafficFlowID := c.Param("id")

	var trafficFlow models.TrafficFlow
	result := database.DB.First(&trafficFlow, trafficFlowID)
	if result.Error != nil || trafficFlow.ID == 0 {
		utils.NotFound(c, "交通流量记录不存在")
		return
	}

	var updateData map[string]interface{}
	if err := c.ShouldBindJSON(&updateData); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	allowedFields := []string{"road_name", "location", "vehicle_count", "average_speed", "flow_level", "record_time"}
	updates := make(map[string]interface{})
	for _, field := range allowedFields {
		if val, ok := updateData[field]; ok {
			updates[field] = val
		}
	}

	if err := database.DB.Model(&trafficFlow).Updates(updates).Error; err != nil {
		utils.InternalServerError(c, "更新失败: "+err.Error())
		return
	}

	utils.SuccessWithMessage(c, "更新成功", nil)
}

func (tfc *TrafficFlowController) Delete(c *gin.Context) {
	trafficFlowID := c.Param("id")

	result := database.DB.Delete(&models.TrafficFlow{}, trafficFlowID)
	if result.RowsAffected == 0 {
		utils.NotFound(c, "交通流量记录不存在")
		return
	}

	utils.SuccessWithMessage(c, "删除成功", nil)
}

type SignalLightController struct{}

func (slc *SignalLightController) List(c *gin.Context) {
	var signalLights []models.SignalLight
	var total int64

	page := 1
	pageSize := 10
	if p := c.Query("page"); p != "" {
		page = parseInt(p, 1)
	}
	if ps := c.Query("page_size"); ps != "" {
		pageSize = parseInt(ps, 10)
	}

	query := database.DB.Model(&models.SignalLight{})

	keyword := c.Query("keyword")
	if keyword != "" {
		query = query.Where("location LIKE ? OR intersection LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	status := c.Query("status")
	if status != "" {
		query = query.Where("status = ?", status)
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&signalLights)

	utils.Success(c, gin.H{
		"list":      signalLights,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func (slc *SignalLightController) GetDetail(c *gin.Context) {
	signalLightID := c.Param("id")

	var signalLight models.SignalLight
	result := database.DB.First(&signalLight, signalLightID)
	if result.Error != nil || signalLight.ID == 0 {
		utils.NotFound(c, "信号灯记录不存在")
		return
	}

	utils.Success(c, signalLight)
}

func (slc *SignalLightController) Create(c *gin.Context) {
	var signalLight models.SignalLight
	if err := c.ShouldBindJSON(&signalLight); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	if signalLight.LastUpdated.IsZero() {
		signalLight.LastUpdated = time.Now()
	}

	if err := database.DB.Create(&signalLight).Error; err != nil {
		utils.InternalServerError(c, "创建失败: "+err.Error())
		return
	}

	utils.SuccessWithMessage(c, "创建成功", gin.H{"id": signalLight.ID})
}

func (slc *SignalLightController) Update(c *gin.Context) {
	signalLightID := c.Param("id")

	var signalLight models.SignalLight
	result := database.DB.First(&signalLight, signalLightID)
	if result.Error != nil || signalLight.ID == 0 {
		utils.NotFound(c, "信号灯记录不存在")
		return
	}

	var updateData map[string]interface{}
	if err := c.ShouldBindJSON(&updateData); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	allowedFields := []string{"location", "intersection", "status", "green_duration", "red_duration", "yellow_duration", "last_updated"}
	updates := make(map[string]interface{})
	for _, field := range allowedFields {
		if val, ok := updateData[field]; ok {
			updates[field] = val
		}
	}

	if _, ok := updates["status"]; ok {
		updates["last_updated"] = time.Now()
	}

	if err := database.DB.Model(&signalLight).Updates(updates).Error; err != nil {
		utils.InternalServerError(c, "更新失败: "+err.Error())
		return
	}

	utils.SuccessWithMessage(c, "更新成功", nil)
}

func (slc *SignalLightController) Control(c *gin.Context) {
	signalLightID := c.Param("id")

	var signalLight models.SignalLight
	result := database.DB.First(&signalLight, signalLightID)
	if result.Error != nil || signalLight.ID == 0 {
		utils.NotFound(c, "信号灯记录不存在")
		return
	}

	var req struct {
		Status         string `json:"status" binding:"required"`
		GreenDuration  int    `json:"green_duration"`
		RedDuration    int    `json:"red_duration"`
		YellowDuration int    `json:"yellow_duration"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	updates := map[string]interface{}{
		"status":       req.Status,
		"last_updated": time.Now(),
	}

	if req.GreenDuration > 0 {
		updates["green_duration"] = req.GreenDuration
	}
	if req.RedDuration > 0 {
		updates["red_duration"] = req.RedDuration
	}
	if req.YellowDuration > 0 {
		updates["yellow_duration"] = req.YellowDuration
	}

	if err := database.DB.Model(&signalLight).Updates(updates).Error; err != nil {
		utils.InternalServerError(c, "控制失败: "+err.Error())
		return
	}

	utils.SuccessWithMessage(c, "信号灯控制成功", nil)
}

func (slc *SignalLightController) Delete(c *gin.Context) {
	signalLightID := c.Param("id")

	result := database.DB.Delete(&models.SignalLight{}, signalLightID)
	if result.RowsAffected == 0 {
		utils.NotFound(c, "信号灯记录不存在")
		return
	}

	utils.SuccessWithMessage(c, "删除成功", nil)
}

type BookingAdminController struct{}

func (bac *BookingAdminController) List(c *gin.Context) {
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

	query := database.DB.Model(&models.Booking{})

	keyword := c.Query("keyword")
	if keyword != "" {
		query = query.Where("order_no LIKE ? OR passenger_name LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	status := c.Query("status")
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

func (bac *BookingAdminController) GetDetail(c *gin.Context) {
	bookingID := c.Param("id")

	var booking models.Booking
	result := database.DB.First(&booking, bookingID)
	if result.Error != nil || booking.ID == 0 {
		utils.NotFound(c, "订单不存在")
		return
	}

	utils.Success(c, booking)
}

func (bac *BookingAdminController) Cancel(c *gin.Context) {
	bookingID := c.Param("id")

	var booking models.Booking
	result := database.DB.First(&booking, bookingID)
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

type TicketAdminController struct{}

func (tac *TicketAdminController) List(c *gin.Context) {
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

	query := database.DB.Model(&models.Ticket{})

	keyword := c.Query("keyword")
	if keyword != "" {
		query = query.Where("route_name LIKE ? OR start_station LIKE ? OR end_station LIKE ?", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	status := c.Query("status")
	if status != "" {
		query = query.Where("status = ?", status)
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	query.Order("departure_time DESC").Offset(offset).Limit(pageSize).Find(&tickets)

	utils.Success(c, gin.H{
		"list":      tickets,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func (tac *TicketAdminController) GetDetail(c *gin.Context) {
	ticketID := c.Param("id")

	var ticket models.Ticket
	result := database.DB.First(&ticket, ticketID)
	if result.Error != nil || ticket.ID == 0 {
		utils.NotFound(c, "车票不存在")
		return
	}

	utils.Success(c, ticket)
}

func (tac *TicketAdminController) Create(c *gin.Context) {
	var ticket models.Ticket
	if err := c.ShouldBindJSON(&ticket); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	ticket.Status = 1

	if err := database.DB.Create(&ticket).Error; err != nil {
		utils.InternalServerError(c, "创建失败: "+err.Error())
		return
	}

	utils.SuccessWithMessage(c, "创建成功", gin.H{"id": ticket.ID})
}

func (tac *TicketAdminController) Update(c *gin.Context) {
	ticketID := c.Param("id")

	var ticket models.Ticket
	result := database.DB.First(&ticket, ticketID)
	if result.Error != nil || ticket.ID == 0 {
		utils.NotFound(c, "车票不存在")
		return
	}

	var updateData map[string]interface{}
	if err := c.ShouldBindJSON(&updateData); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	allowedFields := []string{"route_name", "start_station", "end_station", "departure_time", "arrival_time", "price", "total_seats", "available_seats", "status"}
	updates := make(map[string]interface{})
	for _, field := range allowedFields {
		if val, ok := updateData[field]; ok {
			updates[field] = val
		}
	}

	if err := database.DB.Model(&ticket).Updates(updates).Error; err != nil {
		utils.InternalServerError(c, "更新失败: "+err.Error())
		return
	}

	utils.SuccessWithMessage(c, "更新成功", nil)
}

func (tac *TicketAdminController) Delete(c *gin.Context) {
	ticketID := c.Param("id")

	result := database.DB.Delete(&models.Ticket{}, ticketID)
	if result.RowsAffected == 0 {
		utils.NotFound(c, "车票不存在")
		return
	}

	utils.SuccessWithMessage(c, "删除成功", nil)
}

type RouteAdminController struct{}

func (rac *RouteAdminController) List(c *gin.Context) {
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

	query := database.DB.Model(&models.Route{})

	keyword := c.Query("keyword")
	if keyword != "" {
		query = query.Where("route_name LIKE ? OR route_no LIKE ? OR start_station LIKE ? OR end_station LIKE ?", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	status := c.Query("status")
	if status != "" {
		query = query.Where("status = ?", status)
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

func (rac *RouteAdminController) GetDetail(c *gin.Context) {
	routeID := c.Param("id")

	var route models.Route
	result := database.DB.First(&route, routeID)
	if result.Error != nil || route.ID == 0 {
		utils.NotFound(c, "路线不存在")
		return
	}

	utils.Success(c, route)
}

func (rac *RouteAdminController) Create(c *gin.Context) {
	var route models.Route
	if err := c.ShouldBindJSON(&route); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	route.Status = 1

	if err := database.DB.Create(&route).Error; err != nil {
		utils.InternalServerError(c, "创建失败: "+err.Error())
		return
	}

	utils.SuccessWithMessage(c, "创建成功", gin.H{"id": route.ID})
}

func (rac *RouteAdminController) Update(c *gin.Context) {
	routeID := c.Param("id")

	var route models.Route
	result := database.DB.First(&route, routeID)
	if result.Error != nil || route.ID == 0 {
		utils.NotFound(c, "路线不存在")
		return
	}

	var updateData map[string]interface{}
	if err := c.ShouldBindJSON(&updateData); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	allowedFields := []string{"route_name", "route_no", "start_station", "end_station", "stations", "distance", "duration", "price", "status"}
	updates := make(map[string]interface{})
	for _, field := range allowedFields {
		if val, ok := updateData[field]; ok {
			updates[field] = val
		}
	}

	if err := database.DB.Model(&route).Updates(updates).Error; err != nil {
		utils.InternalServerError(c, "更新失败: "+err.Error())
		return
	}

	utils.SuccessWithMessage(c, "更新成功", nil)
}

func (rac *RouteAdminController) Delete(c *gin.Context) {
	routeID := c.Param("id")

	result := database.DB.Delete(&models.Route{}, routeID)
	if result.RowsAffected == 0 {
		utils.NotFound(c, "路线不存在")
		return
	}

	utils.SuccessWithMessage(c, "删除成功", nil)
}
