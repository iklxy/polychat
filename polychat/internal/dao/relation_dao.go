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

// GetRelation 获取好友列表
func GetRelation(ownerID uint) ([]model.Relation, error) {
	var relations []model.Relation
	err := database.DB.Where("owner_id = ?", ownerID).Find(&relations).Error
	if err != nil {
		return nil, err
	}
	return relations, nil
}

// UpdateRelationNote 更改好友关系备注
func UpdateRelationNote(ownerID, targetID uint, note string) error {
	return database.DB.Model(&model.Relation{}).Where("owner_id = ? and target_id = ? ",
		ownerID, targetID).Update("note", note).Error
}
