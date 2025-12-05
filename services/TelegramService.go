package services

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"os"
	"strconv"
)

var (
	Bot   *tgbotapi.BotAPI
	Token string
)

func InitTelegram() {
	Token = os.Getenv("TELEGRAM_BOT_TOKEN")
	var err error
	Bot, err = tgbotapi.NewBotAPIWithAPIEndpoint(Token, os.Getenv("TELEGRAM_PROXY_URL")+"/bot%s/%s")
	if err != nil {
		log.Fatal(err)
		return
	}
	telegramDebug, _ := strconv.ParseBool(os.Getenv("TELEGRAM_DEBUG"))
	Bot.Debug = telegramDebug
	log.Printf("Authorized on account %s", Bot.Self.UserName)

	url := os.Getenv("BASE_URL") + "/" + Bot.Token
	wh, _ := tgbotapi.NewWebhook(url)
	_, err = Bot.Request(wh)
	if err != nil {
		log.Fatal(err)
	}
}

func SendTelegramMessage(to int64, body string) {
	msg := tgbotapi.NewMessage(to, body)
	msg.ParseMode = "html"
	Bot.Send(msg)
}
