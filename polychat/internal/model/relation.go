package model

// Relation 好友关系表
type Relation struct {
	OwnerID      uint   `gorm:"primaryKey" json:"owner_id"`
	TargetID     uint   `gorm:"primaryKey" json:"target_id"`
	RelationType uint   `gorm:"type:int(1);not null" json:"relation_type"`
	Note         string `gorm:"type:varchar(20);not null" json:"note"`
}
