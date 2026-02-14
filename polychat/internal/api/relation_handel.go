package api

import (
	"net/http"
	"polychat/internal/dao"
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

// PendingRequestDTO 待处理好友请求响应对象
type PendingRequestDTO struct {
	OwnerID      uint   `json:"owner_id"`
	OwnerName    string `json:"owner_name"`
	TargetID     uint   `json:"target_id"`
	RelationType uint   `json:"relation_type"`
	Note         string `json:"note"`
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

// AcceptFriendReq 接受好友请求参数
type AcceptFriendReq struct {
	RequesterID uint `json:"requester_id" binding:"required"`
}

// RejectFriendReq 拒绝好友请求参数
type RejectFriendReq struct {
	RequesterID uint `json:"requester_id" binding:"required"`
}

// AddFriend 发送好友请求
func (h *RelationHandler) AddFriend(ctx *gin.Context) {
	var req AddFriendReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "用户未登录"})
		return
	}
	ownerID := userID.(uint)

	if err := h.relationService.AddFriend(ownerID, req.TargetID, req.Note); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	// 通过 WebSocket 通知对方有新的好友请求
	if ws.ClientMgr.IsUserOnline(req.TargetID) {
		notification := ws.Message{
			Type:       ws.TypeFriendRequest,
			SenderID:   ownerID,
			ReceiverID: req.TargetID,
			Content:    req.Note,
		}
		ws.ClientMgr.SendMessage(notification)
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 200, "message": "好友请求已发送"})
}

// AcceptFriend 接受好友请求
func (h *RelationHandler) AcceptFriend(ctx *gin.Context) {
	var req AcceptFriendReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "用户未登录"})
		return
	}
	currentUserID := userID.(uint)

	if err := h.relationService.AcceptFriendRequest(currentUserID, req.RequesterID); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	// 通过 WebSocket 通知请求方好友请求已被接受
	if ws.ClientMgr.IsUserOnline(req.RequesterID) {
		notification := ws.Message{
			Type:       ws.TypeFriendAccept,
			SenderID:   currentUserID,
			ReceiverID: req.RequesterID,
			Content:    "对方已接受你的好友请求",
		}
		ws.ClientMgr.SendMessage(notification)
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 200, "message": "已接受好友请求"})
}

// RejectFriend 拒绝好友请求
func (h *RelationHandler) RejectFriend(ctx *gin.Context) {
	var req RejectFriendReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "用户未登录"})
		return
	}
	currentUserID := userID.(uint)

	if err := h.relationService.RejectFriendRequest(currentUserID, req.RequesterID); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 200, "message": "已拒绝好友请求"})
}

// GetPendingRequests 获取待处理的好友请求
func (h *RelationHandler) GetPendingRequests(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "用户未登录"})
		return
	}
	currentUserID := userID.(uint)

	relations, err := h.relationService.GetPendingRequests(currentUserID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 转换为DTO，附带请求方的用户名
	var dtos []PendingRequestDTO
	for _, r := range relations {
		ownerName := ""
		user, err := dao.GetUserByID(r.OwnerID)
		if err == nil {
			ownerName = user.Username
		}
		dtos = append(dtos, PendingRequestDTO{
			OwnerID:      r.OwnerID,
			OwnerName:    ownerName,
			TargetID:     r.TargetID,
			RelationType: r.RelationType,
			Note:         r.Note,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 200, "data": dtos})
}

// DeleteFriend 删除好友
func (h *RelationHandler) DeleteFriend(ctx *gin.Context) {
	var req DeleteFriendReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

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
