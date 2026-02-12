package service

import (
	"polychat/internal/dao"
	"polychat/internal/model"
)

type RelationService struct{}

// AddFriend 添加好友
func (s *RelationService) AddFriend(ownerID, targetID uint, note string) error {
	relation := &model.Relation{
		OwnerID:      ownerID,
		TargetID:     targetID,
		RelationType: 1,
		Note:         note,
	}
	return dao.CreateRelation(relation)
}

// DeleteFriend 删除好友
func (s *RelationService) DeleteFriend(ownerID, targetID uint) error {
	return dao.DeleteRelation(ownerID, targetID)
}

// GetFriend 获取好友列表
func (s *RelationService) GetFriend(ownerID uint) ([]model.Relation, error) {
	return dao.GetRelation(ownerID)
}
