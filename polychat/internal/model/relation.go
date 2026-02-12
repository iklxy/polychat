package model

// Relation 好友关系表
type Relation struct {
	OwnerID      uint   `gorm:"primaryKey"`
	TargetID     uint   `gorm:"primaryKey"`
	RelationType uint   `gorm:"type:int(1);not null"`
	Note         string `gorm:"type:varchar(20);not null"`
}
