package controllers

import (
	"smart_transportation/database"
	"smart_transportation/models"
	"smart_transportation/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type AdminController struct{}

type AdminLoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (ac *AdminController) Login(c *gin.Context) {
	var req AdminLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	var admin models.Admin
	database.DB.Where("username = ?", req.Username).First(&admin)
	if admin.ID == 0 {
		utils.Unauthorized(c, "用户名或密码错误")
		return
	}

	if admin.Status != 1 {
		utils.Forbidden(c, "账户已被禁用")
		return
	}

	if !admin.CheckPassword(req.Password) {
		utils.Unauthorized(c, "用户名或密码错误")
		return
	}

	token, err := utils.GenerateToken(admin.ID, admin.Username, admin.Role)
	if err != nil {
		utils.InternalServerError(c, "生成token失败")
		return
	}

	utils.Success(c, LoginResponse{
		Token: token,
		UserInfo: gin.H{
			"id":       admin.ID,
			"username": admin.Username,
			"role":     admin.Role,
		},
	})
}

func (ac *AdminController) GetProfile(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var admin models.Admin
	database.DB.First(&admin, userID)
	if admin.ID == 0 {
		utils.NotFound(c, "管理员不存在")
		return
	}

	utils.Success(c, gin.H{
		"id":       admin.ID,
		"username": admin.Username,
		"role":     admin.Role,
		"status":   admin.Status,
	})
}

func (ac *AdminController) ListUsers(c *gin.Context) {
	var users []models.User
	var total int64

	page := 1
	pageSize := 10
	if p := c.Query("page"); p != "" {
		page = parseInt(p, 1)
	}
	if ps := c.Query("page_size"); ps != "" {
		pageSize = parseInt(ps, 10)
	}

	query := database.DB.Model(&models.User{})

	keyword := c.Query("keyword")
	if keyword != "" {
		query = query.Where("username LIKE ? OR email LIKE ? OR phone LIKE ?", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	status := c.Query("status")
	if status != "" {
		query = query.Where("status = ?", status)
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&users)

	utils.Success(c, gin.H{
		"list":      users,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func (ac *AdminController) UpdateUserStatus(c *gin.Context) {
	userID := c.Param("id")

	var req struct {
		Status int `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	if req.Status != 0 && req.Status != 1 {
		utils.BadRequest(c, "状态值无效")
		return
	}

	result := database.DB.Model(&models.User{}).Where("id = ?", userID).Update("status", req.Status)
	if result.RowsAffected == 0 {
		utils.NotFound(c, "用户不存在")
		return
	}

	utils.SuccessWithMessage(c, "状态更新成功", nil)
}

func (ac *AdminController) ListAdmins(c *gin.Context) {
	var admins []models.Admin
	var total int64

	page := 1
	pageSize := 10
	if p := c.Query("page"); p != "" {
		page = parseInt(p, 1)
	}
	if ps := c.Query("page_size"); ps != "" {
		pageSize = parseInt(ps, 10)
	}

	query := database.DB.Model(&models.Admin{})
	query.Count(&total)

	offset := (page - 1) * pageSize
	query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&admins)

	utils.Success(c, gin.H{
		"list":      admins,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func (ac *AdminController) CreateAdmin(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
		Role     string `json:"role"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	if len(req.Username) < 3 {
		utils.BadRequest(c, "用户名长度不能少于3位")
		return
	}

	if len(req.Password) < 6 {
		utils.BadRequest(c, "密码长度不能少于6位")
		return
	}

	var existingAdmin models.Admin
	database.DB.Where("username = ?", req.Username).First(&existingAdmin)
	if existingAdmin.ID != 0 {
		utils.BadRequest(c, "用户名已存在")
		return
	}

	if req.Role == "" {
		req.Role = "admin"
	}

	admin := models.Admin{
		Username: req.Username,
		Password: req.Password,
		Role:     req.Role,
		Status:   1,
	}

	if err := database.DB.Create(&admin).Error; err != nil {
		utils.InternalServerError(c, "创建失败: "+err.Error())
		return
	}

	utils.SuccessWithMessage(c, "创建成功", gin.H{
		"id":       admin.ID,
		"username": admin.Username,
	})
}

func (ac *AdminController) ChangePassword(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	if len(req.NewPassword) < 6 {
		utils.BadRequest(c, "新密码长度不能少于6位")
		return
	}

	var admin models.Admin
	database.DB.First(&admin, userID)
	if admin.ID == 0 {
		utils.NotFound(c, "管理员不存在")
		return
	}

	if !admin.CheckPassword(req.OldPassword) {
		utils.BadRequest(c, "原密码错误")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		utils.InternalServerError(c, "密码加密失败")
		return
	}

	database.DB.Model(&admin).Update("password", string(hashedPassword))
	utils.SuccessWithMessage(c, "密码修改成功", nil)
}

func parseInt(s string, defaultVal int) int {
	var result int
	for _, c := range s {
		if c >= '0' && c <= '9' {
			result = result*10 + int(c-'0')
		}
	}
	if result == 0 {
		return defaultVal
	}
	return result
}
