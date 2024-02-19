package model

import "time"

type ShopType struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	Name       string    `gorm:"type:varchar(32)" json:"name"`
	Icon       string    `gorm:"type:varchar(255)" json:"icon"`
	Sort       uint      `json:"sort"`
	CreateTime time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
	UpdateTime time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
}

func (ShopType) TableName() string {
	return "tb_shop_type"
}
