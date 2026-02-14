package service

import (
	"errors"
	"polychat/internal/dao"
	"polychat/internal/model"

	"gorm.io/gorm"
)

type RelationService struct{}

// AddFriend 发送好友请求（创建 relation_type=0 的待处理记录）
func (s *RelationService) AddFriend(ownerID, targetID uint, note string) error {
	// 不能加自己为好友
	if ownerID == targetID {
		return errors.New("不能添加自己为好友")
	}

	// 检查目标用户是否存在
	_, err := dao.GetUserByID(targetID)
	if err != nil {
		return errors.New("目标用户不存在")
	}

	// 检查是否已经是好友（任一方向）
	existing, err := dao.GetRelationByPair(ownerID, targetID)
	if err == nil && existing != nil {
		if existing.RelationType == 1 {
			return errors.New("已经是好友了")
		}
		if existing.RelationType == 0 {
			return errors.New("好友请求已发送，请等待对方处理")
		}
	}

	// 检查对方是否已经向我发送过请求
	reverse, err := dao.GetRelationByPair(targetID, ownerID)
	if err == nil && reverse != nil {
		if reverse.RelationType == 0 {
			return errors.New("对方已向你发送好友请求，请在消息箱中处理")
		}
		if reverse.RelationType == 1 {
			return errors.New("已经是好友了")
		}
	}

	// 创建待处理的好友请求（relation_type = 0）
	relation := &model.Relation{
		OwnerID:      ownerID,
		TargetID:     targetID,
		RelationType: 0, // 待处理
		Note:         note,
	}
	return dao.CreateRelation(relation)
}

// AcceptFriendRequest 接受好友请求
// requesterID 是发起请求的人（relation表中的owner_id）
// currentUserID 是当前用户（relation表中的target_id）
func (s *RelationService) AcceptFriendRequest(currentUserID, requesterID uint) error {
	// 验证确实存在一条待处理的好友请求
	relation, err := dao.GetRelationByPair(requesterID, currentUserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("好友请求不存在")
		}
		return err
	}
	if relation.RelationType != 0 {
		return errors.New("该请求已处理")
	}

	// 将原请求更新为已确认（relation_type = 1）
	if err := dao.UpdateRelationType(requesterID, currentUserID, 1); err != nil {
		return err
	}

	// 创建反向好友关系记录，让双方都能在好友列表中看到对方
	reverseRelation := &model.Relation{
		OwnerID:      currentUserID,
		TargetID:     requesterID,
		RelationType: 1,
		Note:         "", // 被接受方可以后续修改备注
	}
	return dao.CreateRelation(reverseRelation)
}

// RejectFriendRequest 拒绝好友请求
func (s *RelationService) RejectFriendRequest(currentUserID, requesterID uint) error {
	// 验证确实存在一条待处理的好友请求
	relation, err := dao.GetRelationByPair(requesterID, currentUserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("好友请求不存在")
		}
		return err
	}
	if relation.RelationType != 0 {
		return errors.New("该请求已处理")
	}

	// 删除待处理的请求记录
	return dao.DeleteRelation(requesterID, currentUserID)
}

// GetPendingRequests 获取当前用户收到的待处理好友请求
func (s *RelationService) GetPendingRequests(userID uint) ([]model.Relation, error) {
	return dao.GetPendingRequests(userID)
}

// DeleteFriend 删除好友（双向删除）
func (s *RelationService) DeleteFriend(ownerID, targetID uint) error {
	// 删除自己的记录
	if err := dao.DeleteRelation(ownerID, targetID); err != nil {
		return err
	}
	// 删除对方的记录
	_ = dao.DeleteRelation(targetID, ownerID)
	return nil
}

// GetFriend 获取已确认的好友列表
func (s *RelationService) GetFriend(ownerID uint) ([]model.Relation, error) {
	return dao.GetRelation(ownerID)
}

// UpdateFriendNote 更新好友备注
func (s *RelationService) UpdateFriendNote(ownerID, targetID uint, note string) error {
	return dao.UpdateRelationNote(ownerID, targetID, note)
}
