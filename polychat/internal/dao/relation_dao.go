package dao

import (
	"polychat/internal/model"
	"polychat/pkg/database"
)

// CreateRelation 创建好友关系
func CreateRelation(relation *model.Relation) error {
	return database.DB.Create(relation).Error
}

// DeleteRelation 删除好友关系
func DeleteRelation(ownerID, targetID uint) error {
	return database.DB.Delete(&model.Relation{}, "owner_id = ? AND target_id = ?", ownerID, targetID).Error
}

// GetRelation 获取已确认的好友列表（relation_type = 1）
func GetRelation(ownerID uint) ([]model.Relation, error) {
	var relations []model.Relation
	err := database.DB.Where("owner_id = ? AND relation_type = 1", ownerID).Find(&relations).Error
	if err != nil {
		return nil, err
	}
	return relations, nil
}

// UpdateRelationNote 更改好友关系备注
func UpdateRelationNote(ownerID, targetID uint, note string) error {
	return database.DB.Model(&model.Relation{}).Where("owner_id = ? AND target_id = ?",
		ownerID, targetID).Update("note", note).Error
}

// GetPendingRequests 获取待处理的好友请求（当前用户是被请求方）
func GetPendingRequests(targetID uint) ([]model.Relation, error) {
	var relations []model.Relation
	err := database.DB.Where("target_id = ? AND relation_type = 0", targetID).Find(&relations).Error
	if err != nil {
		return nil, err
	}
	return relations, nil
}

// GetRelationByPair 查询两个用户之间的关系记录
func GetRelationByPair(ownerID, targetID uint) (*model.Relation, error) {
	var relation model.Relation
	err := database.DB.Where("owner_id = ? AND target_id = ?", ownerID, targetID).First(&relation).Error
	if err != nil {
		return nil, err
	}
	return &relation, nil
}

// UpdateRelationType 更新关系类型（0=待处理, 1=已确认）
func UpdateRelationType(ownerID, targetID uint, relationType uint) error {
	return database.DB.Model(&model.Relation{}).Where("owner_id = ? AND target_id = ?",
		ownerID, targetID).Update("relation_type", relationType).Error
}
