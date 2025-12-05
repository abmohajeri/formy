package services

import (
	"core/config"
	"core/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
	"log"
)

func GetFormTokens(user *models.User) []models.FormToken {
	var formTokens []models.FormToken
	err := config.GetDB().Where("user_id = ?", user.ID).Find(&formTokens).Error
	if err != nil {
		log.Println("Error fetching form tokens:", err)
		return nil
	}
	return formTokens
}

func CreateUserFormToken(update tgbotapi.Update, user *models.User, formName string) (bool, string) {
	var formToken models.FormToken
	errForm := config.GetDB().Where("user_id = ? and name = ?", user.ID, formName).First(&formToken)
	if errForm.RowsAffected == 0 {
		formToken = models.FormToken{
			Uuid:   uuid.New(),
			Name:   formName,
			UserID: user.ID,
			ChatID: update.Message.Chat.ID,
		}
		err := formToken.Save()
		if err != nil {
			return false, `Error occurred\! Please try again\.`
		}
	} else {
		return false, `Form name is exist for your user\! Please try another name\.`
	}
	return true, formToken.Uuid.String()
}
