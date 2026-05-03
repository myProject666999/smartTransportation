package controllers

import (
	"smart_transportation/database"
	"smart_transportation/models"
	"smart_transportation/utils"
	"time"

	"github.com/gin-gonic/gin"
)

type NewsController struct{}

func (nc *NewsController) List(c *gin.Context) {
	var news []models.News
	var total int64

	page := 1
	pageSize := 10
	if p := c.Query("page"); p != "" {
		page = parseInt(p, 1)
	}
	if ps := c.Query("page_size"); ps != "" {
		pageSize = parseInt(ps, 10)
	}

	query := database.DB.Model(&models.News{}).Where("status = ?", 1)

	category := c.Query("category")
	if category != "" {
		query = query.Where("category = ?", category)
	}

	keyword := c.Query("keyword")
	if keyword != "" {
		query = query.Where("title LIKE ? OR summary LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&news)

	utils.Success(c, gin.H{
		"list":      news,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func (nc *NewsController) GetDetail(c *gin.Context) {
	newsID := c.Param("id")

	var news models.News
	result := database.DB.First(&news, newsID)
	if result.Error != nil || news.ID == 0 {
		utils.NotFound(c, "资讯不存在")
		return
	}

	database.DB.Model(&news).Update("view_count", news.ViewCount+1)

	utils.Success(c, news)
}

func (nc *NewsController) Create(c *gin.Context) {
	var news models.News
	if err := c.ShouldBindJSON(&news); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	news.ViewCount = 0
	news.Status = 1

	if err := database.DB.Create(&news).Error; err != nil {
		utils.InternalServerError(c, "创建失败: "+err.Error())
		return
	}

	utils.SuccessWithMessage(c, "创建成功", gin.H{"id": news.ID})
}

func (nc *NewsController) Update(c *gin.Context) {
	newsID := c.Param("id")

	var news models.News
	result := database.DB.First(&news, newsID)
	if result.Error != nil || news.ID == 0 {
		utils.NotFound(c, "资讯不存在")
		return
	}

	var updateData map[string]interface{}
	if err := c.ShouldBindJSON(&updateData); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	allowedFields := []string{"title", "content", "summary", "image", "category", "author", "status"}
	updates := make(map[string]interface{})
	for _, field := range allowedFields {
		if val, ok := updateData[field]; ok {
			updates[field] = val
		}
	}

	if err := database.DB.Model(&news).Updates(updates).Error; err != nil {
		utils.InternalServerError(c, "更新失败: "+err.Error())
		return
	}

	utils.SuccessWithMessage(c, "更新成功", nil)
}

func (nc *NewsController) Delete(c *gin.Context) {
	newsID := c.Param("id")

	result := database.DB.Delete(&models.News{}, newsID)
	if result.RowsAffected == 0 {
		utils.NotFound(c, "资讯不存在")
		return
	}

	utils.SuccessWithMessage(c, "删除成功", nil)
}

type AnnouncementController struct{}

func (ac *AnnouncementController) List(c *gin.Context) {
	var announcements []models.Announcement
	var total int64

	page := 1
	pageSize := 10
	if p := c.Query("page"); p != "" {
		page = parseInt(p, 1)
	}
	if ps := c.Query("page_size"); ps != "" {
		pageSize = parseInt(ps, 10)
	}

	query := database.DB.Model(&models.Announcement{}).Where("status = ?", 1)
	now := time.Now()
	query = query.Where("(expire_time IS NULL OR expire_time > ?)", now)

	keyword := c.Query("keyword")
	if keyword != "" {
		query = query.Where("title LIKE ?", "%"+keyword+"%")
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	query.Order("priority DESC, created_at DESC").Offset(offset).Limit(pageSize).Find(&announcements)

	utils.Success(c, gin.H{
		"list":      announcements,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func (ac *AnnouncementController) GetDetail(c *gin.Context) {
	announcementID := c.Param("id")

	var announcement models.Announcement
	result := database.DB.First(&announcement, announcementID)
	if result.Error != nil || announcement.ID == 0 {
		utils.NotFound(c, "公告不存在")
		return
	}

	utils.Success(c, announcement)
}

func (ac *AnnouncementController) Create(c *gin.Context) {
	var announcement models.Announcement
	if err := c.ShouldBindJSON(&announcement); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	announcement.Status = 1
	if announcement.PublishTime.IsZero() {
		announcement.PublishTime = time.Now()
	}

	if err := database.DB.Create(&announcement).Error; err != nil {
		utils.InternalServerError(c, "创建失败: "+err.Error())
		return
	}

	utils.SuccessWithMessage(c, "创建成功", gin.H{"id": announcement.ID})
}

func (ac *AnnouncementController) Update(c *gin.Context) {
	announcementID := c.Param("id")

	var announcement models.Announcement
	result := database.DB.First(&announcement, announcementID)
	if result.Error != nil || announcement.ID == 0 {
		utils.NotFound(c, "公告不存在")
		return
	}

	var updateData map[string]interface{}
	if err := c.ShouldBindJSON(&updateData); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	allowedFields := []string{"title", "content", "publish_time", "expire_time", "priority", "status"}
	updates := make(map[string]interface{})
	for _, field := range allowedFields {
		if val, ok := updateData[field]; ok {
			updates[field] = val
		}
	}

	if err := database.DB.Model(&announcement).Updates(updates).Error; err != nil {
		utils.InternalServerError(c, "更新失败: "+err.Error())
		return
	}

	utils.SuccessWithMessage(c, "更新成功", nil)
}

func (ac *AnnouncementController) Delete(c *gin.Context) {
	announcementID := c.Param("id")

	result := database.DB.Delete(&models.Announcement{}, announcementID)
	if result.RowsAffected == 0 {
		utils.NotFound(c, "公告不存在")
		return
	}

	utils.SuccessWithMessage(c, "删除成功", nil)
}

type CarouselController struct{}

func (cc *CarouselController) List(c *gin.Context) {
	var carousels []models.Carousel
	var total int64

	page := 1
	pageSize := 10
	if p := c.Query("page"); p != "" {
		page = parseInt(p, 1)
	}
	if ps := c.Query("page_size"); ps != "" {
		pageSize = parseInt(ps, 10)
	}

	query := database.DB.Model(&models.Carousel{}).Where("status = ?", 1)
	query.Count(&total)

	offset := (page - 1) * pageSize
	query.Order("sort_order ASC, created_at DESC").Offset(offset).Limit(pageSize).Find(&carousels)

	utils.Success(c, gin.H{
		"list":      carousels,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func (cc *CarouselController) GetAll(c *gin.Context) {
	var carousels []models.Carousel
	database.DB.Where("status = ?", 1).Order("sort_order ASC, created_at DESC").Find(&carousels)
	utils.Success(c, carousels)
}

func (cc *CarouselController) Create(c *gin.Context) {
	var carousel models.Carousel
	if err := c.ShouldBindJSON(&carousel); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	carousel.Status = 1

	if err := database.DB.Create(&carousel).Error; err != nil {
		utils.InternalServerError(c, "创建失败: "+err.Error())
		return
	}

	utils.SuccessWithMessage(c, "创建成功", gin.H{"id": carousel.ID})
}

func (cc *CarouselController) Update(c *gin.Context) {
	carouselID := c.Param("id")

	var carousel models.Carousel
	result := database.DB.First(&carousel, carouselID)
	if result.Error != nil || carousel.ID == 0 {
		utils.NotFound(c, "轮播图不存在")
		return
	}

	var updateData map[string]interface{}
	if err := c.ShouldBindJSON(&updateData); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	allowedFields := []string{"title", "image", "link", "sort_order", "status"}
	updates := make(map[string]interface{})
	for _, field := range allowedFields {
		if val, ok := updateData[field]; ok {
			updates[field] = val
		}
	}

	if err := database.DB.Model(&carousel).Updates(updates).Error; err != nil {
		utils.InternalServerError(c, "更新失败: "+err.Error())
		return
	}

	utils.SuccessWithMessage(c, "更新成功", nil)
}

func (cc *CarouselController) Delete(c *gin.Context) {
	carouselID := c.Param("id")

	result := database.DB.Delete(&models.Carousel{}, carouselID)
	if result.RowsAffected == 0 {
		utils.NotFound(c, "轮播图不存在")
		return
	}

	utils.SuccessWithMessage(c, "删除成功", nil)
}
