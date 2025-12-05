package models

import (
	"core/config"
	"github.com/google/uuid"
	"time"
)

type FormToken struct {
	Uuid      uuid.UUID `gorm:"type:uuid;not null;primaryKey;unique"`
	Name      string    `gorm:"type:varchar(50)"`
	UserID    uint64
	ChatID    int64
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}

func (formToken *FormToken) Save() error {
	return config.GetDB().Save(&formToken).Error
}

func (formToken *FormToken) RevokeFormToken() error {
	return config.GetDB().Delete(&formToken).Error
}

func GetFormTokenByUuid(Uuid uuid.UUID) (*FormToken, error) {
	var formToken FormToken
	result := config.GetDB().Where("uuid = ?", Uuid).First(&formToken)
	if result.Error != nil {
		return nil, result.Error
	}
	return &formToken, nil
}
