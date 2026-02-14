package api

import (
	"net/http"
	"polychat/internal/service"
	"polychat/internal/ws"

	"github.com/gin-gonic/gin"
)

type RelationHandler struct {
	relationService service.RelationService
}

// FriendDTO 好友信息响应对象
type FriendDTO struct {
	OwnerID      uint   `json:"owner_id"`
	TargetID     uint   `json:"target_id"`
	RelationType uint   `json:"relation_type"`
	Note         string `json:"note"`
	IsOnline     bool   `json:"is_online"`
}

// AddFriendReq 添加好友请求参数
type AddFriendReq struct {
	TargetID uint   `json:"target_id" binding:"required"`
	Note     string `json:"Desc"`
}

// DeleteFriendReq 删除好友请求参数
type DeleteFriendReq struct {
	TargetID uint `json:"target_id" binding:"required"`
}

// UpdateFriendNoteReq 更新好友备注请求参数
type UpdateFriendNoteReq struct {
	TargetID uint   `json:"target_id" binding:"required"`
	Note     string `json:"note" binding:"required"`
}

// AddFriend 添加好友
func (h *RelationHandler) AddFriend(ctx *gin.Context) {
	var req AddFriendReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 从上下文中获取当前登录用户ID
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "用户未登录"})
		return
	}
	ownerID := userID.(uint)

	if err := h.relationService.AddFriend(ownerID, req.TargetID, req.Note); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 200, "message": "添加好友成功"})
}

// DeleteFriend 删除好友
func (h *RelationHandler) DeleteFriend(ctx *gin.Context) {
	var req DeleteFriendReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 从上下文中获取当前登录用户ID
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "用户未登录"})
		return
	}
	ownerID := userID.(uint)

	if err := h.relationService.DeleteFriend(ownerID, req.TargetID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 200, "message": "删除好友成功"})
}

// GetFriend 获取好友列表
func (h *RelationHandler) GetFriend(ctx *gin.Context) {
	// 从上下文中获取当前登录用户ID
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "用户未登录"})
		return
	}
	ownerID := userID.(uint)

	relations, err := h.relationService.GetFriend(ownerID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 转换为DTO并填充在线状态
	var friendDTOs []FriendDTO
	for _, r := range relations {
		friendDTOs = append(friendDTOs, FriendDTO{
			OwnerID:      r.OwnerID,
			TargetID:     r.TargetID,
			RelationType: r.RelationType,
			Note:         r.Note,
			IsOnline:     ws.ClientMgr.IsUserOnline(r.TargetID),
		})
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 200, "data": friendDTOs})
}

// UpdateFriendNote 更新好友备注
func (h *RelationHandler) UpdateFriendNote(ctx *gin.Context) {
	var req UpdateFriendNoteReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 从上下文中获取当前登录用户ID
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "用户未登录"})
		return
	}
	ownerID := userID.(uint)

	if err := h.relationService.UpdateFriendNote(ownerID, req.TargetID, req.Note); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 200, "message": "更新好友备注成功"})
}
