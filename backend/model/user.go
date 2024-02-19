package model

import "time"

type User struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	Phone      string    `gorm:"type:varchar(11)"`
	Password   string    `gorm:"type:varchar(128)"`
	NickName   string    `gorm:"type:varchar(32)"`
	Icon       string    `gorm:"type:varchar(255)" json:"icon"`
	CreateTime time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
	UpdateTime time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
}

func (User) TableName() string {
	return "tb_user"
}
