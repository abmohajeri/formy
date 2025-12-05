package models

import (
	"core/config"
	"time"
)

type User struct {
	ID               uint64    `gorm:"autoIncrement;not null;primaryKey;unique"`
	TelegramUserID   uint64    `gorm:"not null"`
	TelegramUserName string    `gorm:"type:varchar(50);unique"`
	CreatedAt        time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	VerifiedAt       time.Time `gorm:"default:null"`
}

func (f *User) Save() error {
	return config.GetDB().Save(&f).Error
}

func GetByTelegramUserId(telegramUserId uint64) (*User, error) {
	var user User
	result := config.GetDB().Where("telegram_user_id = ?", telegramUserId).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}
