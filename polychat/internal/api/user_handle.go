package api

import (
	"net/http"
	"polychat/internal/service"

	"github.com/gin-gonic/gin"
)

type UserHandle struct {
	userService service.UserService
}

// RegisterRequest 注册请求参数
type RegisterRequest struct {
	Username string `json:"username" binding:"required"` //用户名不能为空
	Password string `json:"password" binding:"required"` //密码不能为空
}

// LoginRequest 登录请求参数
type LoginRequest struct {
	Username string `json:"username" binding:"required"` //用户名不能为空
	Password string `json:"password" binding:"required"` //密码不能为空
}

// Register 注册用户
func (h *UserHandle) Register(c *gin.Context) {
	var req RegisterRequest
	//绑定JSON参数到req结构体
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数错误 : " + err.Error()})
		return
	}

	//调用user_service的Register方法
	err := h.userService.Register(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "注册失败 : " + err.Error()})
		return
	}

	//注册成功
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "注册成功"})
}

func (h *UserHandle) Login(c *gin.Context) {
	var req LoginRequest
	//绑定JSON参数到req结构体
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数错误 : " + err.Error()})
		return
	}

	//调用user_service的Login方法
	token, userID, err := h.userService.Login(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "登录失败 : " + err.Error()})
		return
	}

	//登录成功
	c.JSON(http.StatusOK, gin.H{
		"code":     200,
		"msg":      "登录成功",
		"token":    token,
		"user_id":  userID,
		"username": req.Username,
	})
}
