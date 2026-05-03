package controllers

import (
	"smart_transportation/database"
	"smart_transportation/models"
	"smart_transportation/utils"
	"time"

	"github.com/gin-gonic/gin"
)

type ViolationController struct{}

func (vc *ViolationController) ListUserViolations(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var violations []models.Violation
	var total int64

	page := 1
	pageSize := 10
	if p := c.Query("page"); p != "" {
		page = parseInt(p, 1)
	}
	if ps := c.Query("page_size"); ps != "" {
		pageSize = parseInt(ps, 10)
	}

	query := database.DB.Model(&models.Violation{}).Where("user_id = ?", userID)

	status := c.Query("status")
	if status != "" {
		query = query.Where("status = ?", status)
	}

	violationType := c.Query("type")
	if violationType != "" {
		query = query.Where("violation_type LIKE ?", "%"+violationType+"%")
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	query.Order("violation_time DESC").Offset(offset).Limit(pageSize).Find(&violations)

	utils.Success(c, gin.H{
		"list":      violations,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func (vc *ViolationController) ListAllViolations(c *gin.Context) {
	var violations []models.Violation
	var total int64

	page := 1
	pageSize := 10
	if p := c.Query("page"); p != "" {
		page = parseInt(p, 1)
	}
	if ps := c.Query("page_size"); ps != "" {
		pageSize = parseInt(ps, 10)
	}

	query := database.DB.Model(&models.Violation{})

	keyword := c.Query("keyword")
	if keyword != "" {
		query = query.Where("plate_number LIKE ? OR violation_type LIKE ? OR location LIKE ?", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	status := c.Query("status")
	if status != "" {
		query = query.Where("status = ?", status)
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	query.Order("violation_time DESC").Offset(offset).Limit(pageSize).Find(&violations)

	utils.Success(c, gin.H{
		"list":      violations,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func (vc *ViolationController) GetDetail(c *gin.Context) {
	violationID := c.Param("id")

	var violation models.Violation
	result := database.DB.First(&violation, violationID)
	if result.Error != nil || violation.ID == 0 {
		utils.NotFound(c, "违规记录不存在")
		return
	}

	utils.Success(c, violation)
}

func (vc *ViolationController) Create(c *gin.Context) {
	var violation models.Violation
	if err := c.ShouldBindJSON(&violation); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	violation.Status = 0
	if violation.ViolationTime.IsZero() {
		violation.ViolationTime = time.Now()
	}

	if err := database.DB.Create(&violation).Error; err != nil {
		utils.InternalServerError(c, "创建失败: "+err.Error())
		return
	}

	utils.SuccessWithMessage(c, "创建成功", gin.H{"id": violation.ID})
}

func (vc *ViolationController) Update(c *gin.Context) {
	violationID := c.Param("id")

	var violation models.Violation
	result := database.DB.First(&violation, violationID)
	if result.Error != nil || violation.ID == 0 {
		utils.NotFound(c, "违规记录不存在")
		return
	}

	var updateData map[string]interface{}
	if err := c.ShouldBindJSON(&updateData); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	allowedFields := []string{"user_id", "plate_number", "violation_type", "location", "violation_time", "fine_amount", "points", "status", "payment_time", "description"}
	updates := make(map[string]interface{})
	for _, field := range allowedFields {
		if val, ok := updateData[field]; ok {
			updates[field] = val
		}
	}

	if err := database.DB.Model(&violation).Updates(updates).Error; err != nil {
		utils.InternalServerError(c, "更新失败: "+err.Error())
		return
	}

	utils.SuccessWithMessage(c, "更新成功", nil)
}

func (vc *ViolationController) PayViolation(c *gin.Context) {
	userID, _ := c.Get("user_id")
	violationID := c.Param("id")

	var violation models.Violation
	result := database.DB.Where("id = ? AND user_id = ?", violationID, userID).First(&violation)
	if result.Error != nil || violation.ID == 0 {
		utils.NotFound(c, "违规记录不存在")
		return
	}

	if violation.Status == 1 {
		utils.BadRequest(c, "该违规已处理")
		return
	}

	violation.Status = 1
	violation.PaymentTime = time.Now()

	if err := database.DB.Save(&violation).Error; err != nil {
		utils.InternalServerError(c, "处理失败")
		return
	}

	utils.SuccessWithMessage(c, "违规处理成功", nil)
}

func (vc *ViolationController) Delete(c *gin.Context) {
	violationID := c.Param("id")

	result := database.DB.Delete(&models.Violation{}, violationID)
	if result.RowsAffected == 0 {
		utils.NotFound(c, "违规记录不存在")
		return
	}

	utils.SuccessWithMessage(c, "删除成功", nil)
}

type FeedbackController struct{}

func (fc *FeedbackController) ListUserFeedbacks(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var feedbacks []models.Feedback
	var total int64

	page := 1
	pageSize := 10
	if p := c.Query("page"); p != "" {
		page = parseInt(p, 1)
	}
	if ps := c.Query("page_size"); ps != "" {
		pageSize = parseInt(ps, 10)
	}

	query := database.DB.Model(&models.Feedback{}).Where("user_id = ?", userID)

	status := c.Query("status")
	if status != "" {
		query = query.Where("status = ?", status)
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&feedbacks)

	utils.Success(c, gin.H{
		"list":      feedbacks,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func (fc *FeedbackController) ListAllFeedbacks(c *gin.Context) {
	var feedbacks []models.Feedback
	var total int64

	page := 1
	pageSize := 10
	if p := c.Query("page"); p != "" {
		page = parseInt(p, 1)
	}
	if ps := c.Query("page_size"); ps != "" {
		pageSize = parseInt(ps, 10)
	}

	query := database.DB.Model(&models.Feedback{})

	keyword := c.Query("keyword")
	if keyword != "" {
		query = query.Where("title LIKE ? OR content LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	status := c.Query("status")
	if status != "" {
		query = query.Where("status = ?", status)
	}

	feedbackType := c.Query("type")
	if feedbackType != "" {
		query = query.Where("feedback_type = ?", feedbackType)
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&feedbacks)

	utils.Success(c, gin.H{
		"list":      feedbacks,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func (fc *FeedbackController) GetDetail(c *gin.Context) {
	feedbackID := c.Param("id")

	var feedback models.Feedback
	result := database.DB.First(&feedback, feedbackID)
	if result.Error != nil || feedback.ID == 0 {
		utils.NotFound(c, "反馈记录不存在")
		return
	}

	utils.Success(c, feedback)
}

func (fc *FeedbackController) Create(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var feedback models.Feedback
	if err := c.ShouldBindJSON(&feedback); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	feedback.UserID = userID.(uint)
	feedback.Status = 0

	if err := database.DB.Create(&feedback).Error; err != nil {
		utils.InternalServerError(c, "提交失败: "+err.Error())
		return
	}

	utils.SuccessWithMessage(c, "提交成功", gin.H{"id": feedback.ID})
}

func (fc *FeedbackController) Reply(c *gin.Context) {
	feedbackID := c.Param("id")

	var feedback models.Feedback
	result := database.DB.First(&feedback, feedbackID)
	if result.Error != nil || feedback.ID == 0 {
		utils.NotFound(c, "反馈记录不存在")
		return
	}

	var req struct {
		Reply string `json:"reply" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	feedback.Reply = req.Reply
	feedback.Status = 1
	feedback.ReplyTime = time.Now()

	if err := database.DB.Save(&feedback).Error; err != nil {
		utils.InternalServerError(c, "回复失败")
		return
	}

	utils.SuccessWithMessage(c, "回复成功", nil)
}

func (fc *FeedbackController) Update(c *gin.Context) {
	feedbackID := c.Param("id")

	var feedback models.Feedback
	result := database.DB.First(&feedback, feedbackID)
	if result.Error != nil || feedback.ID == 0 {
		utils.NotFound(c, "反馈记录不存在")
		return
	}

	var updateData map[string]interface{}
	if err := c.ShouldBindJSON(&updateData); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	allowedFields := []string{"title", "content", "feedback_type", "status", "reply", "reply_time"}
	updates := make(map[string]interface{})
	for _, field := range allowedFields {
		if val, ok := updateData[field]; ok {
			updates[field] = val
		}
	}

	if err := database.DB.Model(&feedback).Updates(updates).Error; err != nil {
		utils.InternalServerError(c, "更新失败: "+err.Error())
		return
	}

	utils.SuccessWithMessage(c, "更新成功", nil)
}

func (fc *FeedbackController) Delete(c *gin.Context) {
	feedbackID := c.Param("id")

	result := database.DB.Delete(&models.Feedback{}, feedbackID)
	if result.RowsAffected == 0 {
		utils.NotFound(c, "反馈记录不存在")
		return
	}

	utils.SuccessWithMessage(c, "删除成功", nil)
}
