package model

import "time"

type Blog struct {
	ID         uint `gorm:"primaryKey" json:"id"`
	ShopID     uint
	UserID     uint
	Title      string `gorm:"type:varchar(255)" json:"title"`
	Images     string `gorm:"type:varchar(2048)" json:"images"`
	Content    string `gorm:"type:varchar(2048)" json:"content"`
	Liked      uint8  `json:"liked"`
	Commnets   uint8
	Icon       string    `gorm:"-" json:"icon"`
	Name       string    `gorm:"-" json:"name"`
	CreateTime time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
	UpdateTime time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
}

func (Blog) TableName() string {
	return "tb_blog"
}
