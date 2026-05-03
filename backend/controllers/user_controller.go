package controllers

import (
	"regexp"
	"smart_transportation/database"
	"smart_transportation/models"
	"smart_transportation/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type UserController struct{}

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token    string      `json:"token"`
	UserInfo interface{} `json:"user_info"`
}

func (uc *UserController) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	if len(req.Username) < 3 || len(req.Username) > 20 {
		utils.BadRequest(c, "用户名长度需要在3-20个字符之间")
		return
	}

	if len(req.Password) < 6 {
		utils.BadRequest(c, "密码长度不能少于6位")
		return
	}

	if req.Email != "" {
		emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
		if !emailRegex.MatchString(req.Email) {
			utils.BadRequest(c, "邮箱格式不正确")
			return
		}
	}

	var existingUser models.User
	database.DB.Where("username = ?", req.Username).First(&existingUser)
	if existingUser.ID != 0 {
		utils.BadRequest(c, "用户名已存在")
		return
	}

	user := models.User{
		Username: req.Username,
		Password: req.Password,
		Email:    req.Email,
		Phone:    req.Phone,
		Status:   1,
	}

	if err := database.DB.Create(&user).Error; err != nil {
		utils.InternalServerError(c, "注册失败: "+err.Error())
		return
	}

	utils.SuccessWithMessage(c, "注册成功", gin.H{
		"user_id":  user.ID,
		"username": user.Username,
	})
}

func (uc *UserController) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	var user models.User
	database.DB.Where("username = ?", req.Username).First(&user)
	if user.ID == 0 {
		utils.Unauthorized(c, "用户名或密码错误")
		return
	}

	if user.Status != 1 {
		utils.Forbidden(c, "账户已被禁用，请联系管理员")
		return
	}

	if !user.CheckPassword(req.Password) {
		utils.Unauthorized(c, "用户名或密码错误")
		return
	}

	token, err := utils.GenerateToken(user.ID, user.Username, "user")
	if err != nil {
		utils.InternalServerError(c, "生成token失败")
		return
	}

	utils.Success(c, LoginResponse{
		Token: token,
		UserInfo: gin.H{
			"id":        user.ID,
			"username":  user.Username,
			"email":     user.Email,
			"phone":     user.Phone,
			"real_name": user.RealName,
			"avatar":    user.Avatar,
			"role":      "user",
		},
	})
}

func (uc *UserController) GetProfile(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var user models.User
	database.DB.First(&user, userID)
	if user.ID == 0 {
		utils.NotFound(c, "用户不存在")
		return
	}

	utils.Success(c, gin.H{
		"id":        user.ID,
		"username":  user.Username,
		"email":     user.Email,
		"phone":     user.Phone,
		"real_name": user.RealName,
		"id_card":   user.IDCard,
		"avatar":    user.Avatar,
		"status":    user.Status,
		"created_at": user.CreatedAt,
	})
}

func (uc *UserController) UpdateProfile(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var updateData map[string]interface{}
	if err := c.ShouldBindJSON(&updateData); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	allowedFields := []string{"email", "phone", "real_name", "id_card", "avatar"}
	updates := make(map[string]interface{})
	for _, field := range allowedFields {
		if val, ok := updateData[field]; ok {
			updates[field] = val
		}
	}

	if err := database.DB.Model(&models.User{}).Where("id = ?", userID).Updates(updates).Error; err != nil {
		utils.InternalServerError(c, "更新失败: "+err.Error())
		return
	}

	utils.SuccessWithMessage(c, "更新成功", nil)
}

func (uc *UserController) ChangePassword(c *gin.Context) {
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

	var user models.User
	database.DB.First(&user, userID)
	if user.ID == 0 {
		utils.NotFound(c, "用户不存在")
		return
	}

	if !user.CheckPassword(req.OldPassword) {
		utils.BadRequest(c, "原密码错误")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		utils.InternalServerError(c, "密码加密失败")
		return
	}

	database.DB.Model(&user).Update("password", string(hashedPassword))
	utils.SuccessWithMessage(c, "密码修改成功", nil)
}
