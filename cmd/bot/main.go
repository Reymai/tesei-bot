package main

import (
	"log"
	"os"
	"tesei-bot/pkg/telegram"
	"tesei-bot/pkg/tsi"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TESEI_KEY"))
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	t := tsi.NewTSI()
	telegramBot := telegram.NewBot(bot, t)
	err = telegramBot.Start()
	if err != nil {
		log.Fatal(err)
	}
}
