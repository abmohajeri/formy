package models

import (
	"core/config"
	"time"
)

type AllowedDomain struct {
	ID        uint64 `gorm:"autoIncrement;not null;primaryKey;unique"`
	Name      string `gorm:"type:varchar(50)"`
	UserID    uint64
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}

func (domain *AllowedDomain) Save() error {
	return config.GetDB().Save(&domain).Error
}

func (domain *AllowedDomain) DeleteDomain() error {
	return config.GetDB().Delete(&domain).Error
}

func GetDomainById(id uint64) (*AllowedDomain, error) {
	var domain AllowedDomain
	result := config.GetDB().Where("id = ?", id).First(&domain)
	if result.Error != nil {
		return nil, result.Error
	}
	return &domain, nil
}
