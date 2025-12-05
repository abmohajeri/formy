package utils

import (
	"github.com/asaskevich/govalidator"
	"github.com/google/uuid"
)

func GetUUIDFromString(uuidString string) uuid.UUID {
	if !govalidator.IsUUID(uuidString) {
		return uuid.Nil
	}
	var err error
	validUuid, err := uuid.Parse(uuidString)
	if err != nil {
		return uuid.Nil
	}
	return validUuid
}
