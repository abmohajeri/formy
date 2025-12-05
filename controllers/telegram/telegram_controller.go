package telegram

import (
	"core/models"
	"core/services"
	"core/utils"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func TelegramWebhookHandler(c *gin.Context) {
	defer c.Request.Body.Close()

	var update tgbotapi.Update
	if err := json.NewDecoder(c.Request.Body).Decode(&update); err != nil {
		log.Println("Failed to decode request body:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	if update.Message != nil {
		handleCommand(update)
	} else if update.CallbackQuery != nil {
		handleCallbackQuery(update)
	}
}

func handleCommand(update tgbotapi.Update) {
	switch update.Message.Command() {
	case "start":
		handleStartCommand(update)
	case "get_token":
		handleGetTokenCommand(update)
	case "tokens_list":
		handleTokensListCommand(update)
	case "add_domain":
		handleAddDomainCommand(update)
	case "domains_list":
		handleDomainsListCommand(update)
	default:
		handleUnknownCommand(update)
	}
}

func handleUnknownCommand(update tgbotapi.Update) {
	chatID := update.Message.Chat.ID
	msg := tgbotapi.NewMessage(chatID, "")
	msg.ParseMode = "MarkdownV2"
	msg.Text = `I don't know that command\.`
	services.Bot.Send(msg)
}

func handleStartCommand(update tgbotapi.Update) {
	chatID := update.Message.Chat.ID
	userID := update.Message.From.ID
	msg := tgbotapi.NewMessage(chatID, "")
	msg.ParseMode = "MarkdownV2"
	telegramUserId := uint64(userID)
	_, err := models.GetByTelegramUserId(telegramUserId)
	if err != nil {
		newUser := models.User{
			TelegramUserID:   telegramUserId,
			TelegramUserName: update.Message.Chat.UserName,
			VerifiedAt:       time.Now(),
		}
		err := newUser.Save()
		if err != nil {
			msg.Text = `Error occurred\! Please try again\.`
			services.Bot.Send(msg)
			return
		}
	}

	msg.Text = fmt.Sprintf("Hello, *%s* 👋\n"+
		"Thank you for choosing Formy\\! 🎉\n\n"+
		"Below are commands you can do:\n"+
		"To get a new form token, type: \\/get\\_token FORM\\_NAME\n"+
		"To view all your form tokens, type: \\/tokens\\_list\n"+
		"To add a new domain, type: \\/add\\_domain DOMAIN\n"+
		"To view all your allowed domains, type: \\/domains\\_list\n",
		update.Message.From.FirstName)

	services.Bot.Send(msg)
}

func handleGetTokenCommand(update tgbotapi.Update) {
	chatID := update.Message.Chat.ID
	userID := update.Message.From.ID
	msg := tgbotapi.NewMessage(chatID, "")
	msg.ParseMode = "MarkdownV2"
	telegramUserId := uint64(userID)
	user, err := models.GetByTelegramUserId(telegramUserId)
	if err != nil {
		msg.Text = `User not found\! Start the bot first\.`
		services.Bot.Send(msg)
		return
	}
	if user.VerifiedAt.IsZero() {
		msg.Text = `User not validated\!`
		services.Bot.Send(msg)
		return
	}

	command := update.Message.Text
	commandRegex := regexp.MustCompile(`^/get_token\s+(\S+)$`)
	matches := commandRegex.FindStringSubmatch(command)
	if len(matches) == 0 {
		msg.Text = fmt.Sprintf("Invalid command format\\.\n\nTo get form token run: \\/get\\_token FORM\\_NAME")
		services.Bot.Send(msg)
		return
	}
	formName := matches[1]

	if ok, msgText := services.CreateUserFormToken(update, user, formName); !ok {
		msg.Text = msgText
	} else {
		msg.Text = fmt.Sprintf("Use this token for your form:\n`%s`\n\nSend \\/tokens\\_list to see tokens\\.", msgText)
	}
	services.Bot.Send(msg)
}

func handleTokensListCommand(update tgbotapi.Update) {
	chatID := update.Message.Chat.ID
	userID := update.Message.From.ID
	msg := tgbotapi.NewMessage(chatID, "")
	msg.ParseMode = "MarkdownV2"
	telegramUserId := uint64(userID)
	user, err := models.GetByTelegramUserId(telegramUserId)
	if err != nil {
		msg.Text = `User not found\! Please start the bot first\.`
		services.Bot.Send(msg)
		return
	}
	msg.Text = "*Your form tokens:*\nSelect a token below to see details\\."
	tokens := services.GetFormTokens(user)
	if len(tokens) == 0 {
		msg.Text = "You don't have any form tokens yet\\.\n\nTo get form token run: \\/get\\_token FORM\\_NAME"
		services.Bot.Send(msg)
		return
	}
	var rows [][]tgbotapi.InlineKeyboardButton
	for _, token := range tokens {
		button := tgbotapi.NewInlineKeyboardButtonData(token.Name, fmt.Sprintf("token_%s", token.Uuid))
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(button))
	}
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(rows...)
	services.Bot.Send(msg)
}

func handleAddDomainCommand(update tgbotapi.Update) {
	chatID := update.Message.Chat.ID
	userID := update.Message.From.ID
	msg := tgbotapi.NewMessage(chatID, "")
	msg.ParseMode = "MarkdownV2"
	telegramUserId := uint64(userID)
	user, err := models.GetByTelegramUserId(telegramUserId)
	if err != nil {
		msg.Text = `User not found\! Start the bot first\.`
		services.Bot.Send(msg)
		return
	}
	if user.VerifiedAt.IsZero() {
		msg.Text = `User not validated\!`
		services.Bot.Send(msg)
		return
	}

	command := update.Message.Text
	commandRegex := regexp.MustCompile(`^/add_domain\s+([a-zA-Z0-9-]+(\.[a-zA-Z0-9-]+)+|localhost|127\.0\.0\.1)$`)
	matches := commandRegex.FindStringSubmatch(command)
	if len(matches) == 0 {
		msg.Text = fmt.Sprintf("Invalid command format\\.\n\nTo add allowed domain run: \\/add\\_domain DOMAIN")
		services.Bot.Send(msg)
		return
	}
	domain := matches[1]

	if ok, msgText := services.CreateUserAllowedDomain(user, domain); !ok {
		msg.Text = msgText
	} else {
		msg.Text = fmt.Sprintf("✅ Domain created successfully\\.\n\nSend \\/domains\\_list to see domains\\.")
	}
	services.Bot.Send(msg)
}

func handleDomainsListCommand(update tgbotapi.Update) {
	chatID := update.Message.Chat.ID
	userID := update.Message.From.ID
	msg := tgbotapi.NewMessage(chatID, "")
	msg.ParseMode = "MarkdownV2"
	telegramUserId := uint64(userID)
	user, err := models.GetByTelegramUserId(telegramUserId)
	if err != nil {
		msg.Text = `User not found\! Please start the bot first\.`
		services.Bot.Send(msg)
		return
	}
	msg.Text = "*Your domains:*\nSelect a domain below to see details\\."
	domains := services.GetDomains(user.ID)
	if len(domains) == 0 {
		msg.Text = "You don't have any allowed domains yet\\.\n\nTo add allowed domain run: \\/add\\_domain DOMAIN"
		services.Bot.Send(msg)
		return
	}
	var rows [][]tgbotapi.InlineKeyboardButton
	for _, domain := range domains {
		button := tgbotapi.NewInlineKeyboardButtonData(domain.Name, fmt.Sprintf("domain_%d", domain.ID))
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(button))
	}
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(rows...)
	services.Bot.Send(msg)
}

func handleCallbackQuery(update tgbotapi.Update) {
	callbackData := update.CallbackQuery.Data
	if strings.HasPrefix(callbackData, "token_") {
		handleTokenCallbackQuery(update, callbackData)
	} else if strings.HasPrefix(callbackData, "revoke_token_") {
		handleRevokeTokenCallbackQuery(update, callbackData)
	} else if strings.HasPrefix(callbackData, "domain_") {
		handleDomainCallbackQuery(update, callbackData)
	} else if strings.HasPrefix(callbackData, "delete_domain_") {
		handleDeleteDomainCallbackQuery(update, callbackData)
	}
}

func handleTokenCallbackQuery(update tgbotapi.Update, data string) {
	chatID := update.CallbackQuery.Message.Chat.ID
	messageID := update.CallbackQuery.Message.MessageID
	tokenUUID := utils.GetUUIDFromString(strings.TrimPrefix(data, "token_"))
	formToken, err := models.GetFormTokenByUuid(tokenUUID)
	if err != nil {
		return
	}
	revokeButton := tgbotapi.NewInlineKeyboardButtonData("Revoke", "revoke_token_"+tokenUUID.String())
	inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(revokeButton),
	)
	msgText := fmt.Sprintf("*%s*: `%s`\n", formToken.Name, formToken.Uuid)
	editedMsg := tgbotapi.NewEditMessageText(chatID, messageID, msgText)
	editedMsg.ParseMode = "MarkdownV2"
	editedMsg.ReplyMarkup = &inlineKeyboard
	services.Bot.Send(editedMsg)
}

func handleRevokeTokenCallbackQuery(update tgbotapi.Update, data string) {
	chatID := update.CallbackQuery.Message.Chat.ID
	messageID := update.CallbackQuery.Message.MessageID
	tokenUUID := utils.GetUUIDFromString(strings.TrimPrefix(data, "revoke_token_"))
	formToken, err := models.GetFormTokenByUuid(tokenUUID)
	if err != nil {
		return
	}
	err = formToken.RevokeFormToken()
	if err != nil {
		return
	}
	editedMsg := tgbotapi.NewEditMessageText(chatID, messageID, "✅ Token revoked successfully!\n\nSend /tokens_list to see tokens.")
	services.Bot.Send(editedMsg)
}

func handleDomainCallbackQuery(update tgbotapi.Update, data string) {
	chatID := update.CallbackQuery.Message.Chat.ID
	messageID := update.CallbackQuery.Message.MessageID
	domainId, _ := strconv.ParseUint(strings.TrimPrefix(data, "domain_"), 10, 64)
	domain, err := models.GetDomainById(domainId)
	if err != nil {
		return
	}
	revokeButton := tgbotapi.NewInlineKeyboardButtonData("Delete", "delete_domain_"+fmt.Sprintf("%d", domain.ID))
	inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(revokeButton),
	)
	msgText := fmt.Sprintf("*%s*", strings.ReplaceAll(domain.Name, ".", "\\."))
	editedMsg := tgbotapi.NewEditMessageText(chatID, messageID, msgText)
	editedMsg.ParseMode = "MarkdownV2"
	editedMsg.ReplyMarkup = &inlineKeyboard
	services.Bot.Send(editedMsg)
}

func handleDeleteDomainCallbackQuery(update tgbotapi.Update, data string) {
	chatID := update.CallbackQuery.Message.Chat.ID
	messageID := update.CallbackQuery.Message.MessageID
	domainId, _ := strconv.ParseUint(strings.TrimPrefix(data, "delete_domain_"), 10, 64)
	domain, err := models.GetDomainById(domainId)
	if err != nil {
		return
	}
	err = domain.DeleteDomain()
	if err != nil {
		return
	}
	editedMsg := tgbotapi.NewEditMessageText(chatID, messageID, "✅ Domain deleted successfully!\n\nSend /domains_list to see domains.")
	services.Bot.Send(editedMsg)
}
