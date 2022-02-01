package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"tesei-bot/pkg/tsi"
)

type Bot struct {
	bot *tgbotapi.BotAPI
	tsi *tsi.TSI
}

func NewBot(bot *tgbotapi.BotAPI, t *tsi.TSI) *Bot {
	return &Bot{bot: bot, tsi: t}
}

func (b *Bot) Start() error {
	log.Printf("Authorized on account %s", b.bot.Self.UserName)

	updates := b.initUpdatesChannel()

	b.handleUpdates(updates)

	return nil
}

func (b *Bot) initUpdatesChannel() tgbotapi.UpdatesChannel {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.bot.GetUpdatesChan(u)
	return updates
}

func (b *Bot) handleUpdates(updates tgbotapi.UpdatesChannel) {
	for update := range updates {
		if update.Message == nil { // Ignore any non-Message Updates
			continue
		}

		if update.Message.IsCommand() {
			b.handleCommand(update.Message)
			continue
		}
		b.handleMessage(update.Message)
	}
}

func (b *Bot) handleMessage(message *tgbotapi.Message) {
	log.Printf("[%s] %s", message.From.UserName, message.Text)

	b.sendMessage(message)
}

func (b *Bot) handleCommand(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, "")

	switch message.Command() {
	case "start":
		msg.Text = b.handleStartCommand()
	case "news":
		msg.Text = b.handleNewsCommand()
	default:
		log.Printf("Unknown command: %s", message.Command())
		msg.Text = "Я не знаю такой команды"
	}

	_, err := b.bot.Send(msg)
	if err != nil {
		return err
	}

	return nil
}

func (b *Bot) handleNewsCommand() string {
	news, err := b.tsi.GetNews()
	if err != nil {
		return "Something went wrong. Can't get news right now"
	}
	return "Новости из TSI:\n" + news
}

func (b *Bot) handleStartCommand() string {
	return "Привет! Я бот для получения информации и новостей из TSI. \n" +
		"Пока что я маленький и слабенький, но обещаю, что скоро стану большим и сильным!"
}

func (b *Bot) sendMessage(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, message.Text)
	msg.DisableWebPagePreview = true
	_, err := b.bot.Send(msg)
	if err != nil {
		return err
	}
	return nil
}
