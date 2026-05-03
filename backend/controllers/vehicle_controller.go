package controllers

import (
	"smart_transportation/database"
	"smart_transportation/models"
	"smart_transportation/utils"
	"time"

	"github.com/gin-gonic/gin"
)

type VehicleController struct{}

func (vc *VehicleController) ListVehicles(c *gin.Context) {
	var vehicles []models.Vehicle
	var total int64

	page := 1
	pageSize := 10
	if p := c.Query("page"); p != "" {
		page = parseInt(p, 1)
	}
	if ps := c.Query("page_size"); ps != "" {
		pageSize = parseInt(ps, 10)
	}

	query := database.DB.Model(&models.Vehicle{})

	keyword := c.Query("keyword")
	if keyword != "" {
		query = query.Where("plate_number LIKE ? OR owner_name LIKE ? OR vehicle_type LIKE ?", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	status := c.Query("status")
	if status != "" {
		query = query.Where("status = ?", status)
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&vehicles)

	utils.Success(c, gin.H{
		"list":      vehicles,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func (vc *VehicleController) CreateVehicle(c *gin.Context) {
	var vehicle models.Vehicle
	if err := c.ShouldBindJSON(&vehicle); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	vehicle.Status = 1

	if err := database.DB.Create(&vehicle).Error; err != nil {
		utils.InternalServerError(c, "创建失败: "+err.Error())
		return
	}

	utils.SuccessWithMessage(c, "创建成功", gin.H{"id": vehicle.ID})
}

func (vc *VehicleController) UpdateVehicle(c *gin.Context) {
	vehicleID := c.Param("id")

	var vehicle models.Vehicle
	result := database.DB.First(&vehicle, vehicleID)
	if result.Error != nil || vehicle.ID == 0 {
		utils.NotFound(c, "车辆不存在")
		return
	}

	var updateData map[string]interface{}
	if err := c.ShouldBindJSON(&updateData); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	allowedFields := []string{"plate_number", "vehicle_type", "owner_name", "owner_phone", "vehicle_model", "status"}
	updates := make(map[string]interface{})
	for _, field := range allowedFields {
		if val, ok := updateData[field]; ok {
			updates[field] = val
		}
	}

	if err := database.DB.Model(&vehicle).Updates(updates).Error; err != nil {
		utils.InternalServerError(c, "更新失败: "+err.Error())
		return
	}

	utils.SuccessWithMessage(c, "更新成功", nil)
}

func (vc *VehicleController) DeleteVehicle(c *gin.Context) {
	vehicleID := c.Param("id")

	result := database.DB.Delete(&models.Vehicle{}, vehicleID)
	if result.RowsAffected == 0 {
		utils.NotFound(c, "车辆不存在")
		return
	}

	utils.SuccessWithMessage(c, "删除成功", nil)
}

type VehicleMonitorController struct{}

func (vmc *VehicleMonitorController) List(c *gin.Context) {
	var monitors []models.VehicleMonitor
	var total int64

	page := 1
	pageSize := 10
	if p := c.Query("page"); p != "" {
		page = parseInt(p, 1)
	}
	if ps := c.Query("page_size"); ps != "" {
		pageSize = parseInt(ps, 10)
	}

	query := database.DB.Model(&models.VehicleMonitor{})

	plateNumber := c.Query("plate_number")
	if plateNumber != "" {
		query = query.Where("plate_number LIKE ?", "%"+plateNumber+"%")
	}

	status := c.Query("status")
	if status != "" {
		query = query.Where("status = ?", status)
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	query.Order("monitor_time DESC").Offset(offset).Limit(pageSize).Find(&monitors)

	utils.Success(c, gin.H{
		"list":      monitors,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func (vmc *VehicleMonitorController) GetByPlateNumber(c *gin.Context) {
	plateNumber := c.Param("plate_number")

	var monitor models.VehicleMonitor
	result := database.DB.Where("plate_number = ?", plateNumber).Order("monitor_time DESC").First(&monitor)
	if result.Error != nil || monitor.ID == 0 {
		utils.NotFound(c, "未找到该车辆的监控信息")
		return
	}

	utils.Success(c, monitor)
}

func (vmc *VehicleMonitorController) Create(c *gin.Context) {
	var monitor models.VehicleMonitor
	if err := c.ShouldBindJSON(&monitor); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	monitor.Status = 1
	if monitor.MonitorTime.IsZero() {
		monitor.MonitorTime = time.Now()
	}

	if err := database.DB.Create(&monitor).Error; err != nil {
		utils.InternalServerError(c, "创建失败: "+err.Error())
		return
	}

	utils.SuccessWithMessage(c, "创建成功", gin.H{"id": monitor.ID})
}

func (vmc *VehicleMonitorController) Update(c *gin.Context) {
	monitorID := c.Param("id")

	var monitor models.VehicleMonitor
	result := database.DB.First(&monitor, monitorID)
	if result.Error != nil || monitor.ID == 0 {
		utils.NotFound(c, "监控记录不存在")
		return
	}

	var updateData map[string]interface{}
	if err := c.ShouldBindJSON(&updateData); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	allowedFields := []string{"vehicle_id", "plate_number", "location", "longitude", "latitude", "speed", "direction", "status", "monitor_time"}
	updates := make(map[string]interface{})
	for _, field := range allowedFields {
		if val, ok := updateData[field]; ok {
			updates[field] = val
		}
	}

	if err := database.DB.Model(&monitor).Updates(updates).Error; err != nil {
		utils.InternalServerError(c, "更新失败: "+err.Error())
		return
	}

	utils.SuccessWithMessage(c, "更新成功", nil)
}

func (vmc *VehicleMonitorController) Delete(c *gin.Context) {
	monitorID := c.Param("id")

	result := database.DB.Delete(&models.VehicleMonitor{}, monitorID)
	if result.RowsAffected == 0 {
		utils.NotFound(c, "监控记录不存在")
		return
	}

	utils.SuccessWithMessage(c, "删除成功", nil)
}

type RoadConditionController struct{}

func (rcc *RoadConditionController) List(c *gin.Context) {
	var roadConditions []models.RoadCondition
	var total int64

	page := 1
	pageSize := 10
	if p := c.Query("page"); p != "" {
		page = parseInt(p, 1)
	}
	if ps := c.Query("page_size"); ps != "" {
		pageSize = parseInt(ps, 10)
	}

	query := database.DB.Model(&models.RoadCondition{})

	keyword := c.Query("keyword")
	if keyword != "" {
		query = query.Where("road_name LIKE ? OR location LIKE ? OR description LIKE ?", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	condition := c.Query("condition")
	if condition != "" {
		query = query.Where("condition = ?", condition)
	}

	severity := c.Query("severity")
	if severity != "" {
		query = query.Where("severity = ?", severity)
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	query.Order("start_time DESC").Offset(offset).Limit(pageSize).Find(&roadConditions)

	utils.Success(c, gin.H{
		"list":      roadConditions,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func (rcc *RoadConditionController) GetDetail(c *gin.Context) {
	roadConditionID := c.Param("id")

	var roadCondition models.RoadCondition
	result := database.DB.First(&roadCondition, roadConditionID)
	if result.Error != nil || roadCondition.ID == 0 {
		utils.NotFound(c, "路况信息不存在")
		return
	}

	utils.Success(c, roadCondition)
}

func (rcc *RoadConditionController) Create(c *gin.Context) {
	var roadCondition models.RoadCondition
	if err := c.ShouldBindJSON(&roadCondition); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	roadCondition.Status = 1
	if roadCondition.StartTime.IsZero() {
		roadCondition.StartTime = time.Now()
	}

	if err := database.DB.Create(&roadCondition).Error; err != nil {
		utils.InternalServerError(c, "创建失败: "+err.Error())
		return
	}

	utils.SuccessWithMessage(c, "创建成功", gin.H{"id": roadCondition.ID})
}

func (rcc *RoadConditionController) Update(c *gin.Context) {
	roadConditionID := c.Param("id")

	var roadCondition models.RoadCondition
	result := database.DB.First(&roadCondition, roadConditionID)
	if result.Error != nil || roadCondition.ID == 0 {
		utils.NotFound(c, "路况信息不存在")
		return
	}

	var updateData map[string]interface{}
	if err := c.ShouldBindJSON(&updateData); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	allowedFields := []string{"road_name", "location", "condition", "severity", "description", "start_time", "end_time", "status"}
	updates := make(map[string]interface{})
	for _, field := range allowedFields {
		if val, ok := updateData[field]; ok {
			updates[field] = val
		}
	}

	if err := database.DB.Model(&roadCondition).Updates(updates).Error; err != nil {
		utils.InternalServerError(c, "更新失败: "+err.Error())
		return
	}

	utils.SuccessWithMessage(c, "更新成功", nil)
}

func (rcc *RoadConditionController) Delete(c *gin.Context) {
	roadConditionID := c.Param("id")

	result := database.DB.Delete(&models.RoadCondition{}, roadConditionID)
	if result.RowsAffected == 0 {
		utils.NotFound(c, "路况信息不存在")
		return
	}

	utils.SuccessWithMessage(c, "删除成功", nil)
}
