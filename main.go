package main

import (
	"chatGptTeleBot/ChatGpt"
	"chatGptTeleBot/messagesPool"
	"chatGptTeleBot/utilities"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"strings"
	"time"
)

var loopTtlChat = time.Minute * 10
var ChatGptApp *ChatGpt.AppChatGpt

var whiteList = make(map[int]bool)

func main() {
	utilities.CheckEnvFile()
	ChatGptApp = ChatGpt.New(utilities.ChatGptToken)
	bot, err := tgbotapi.NewBotAPI(utilities.TelegramToken)
	if err != nil {
		log.Panic(err)
	}

	go LoopTtl(bot)
	bot.Debug = false

	log.Printf("Authorized on account %s", bot.Self.UserName)

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60
	updates, err := bot.GetUpdatesChan(updateConfig)
	if err != nil {
		log.Panic(err)
	}

	// Цикл получения новых сообщений
	for update := range updates {
		if update.Message == nil { // игнорируем обновления, которые не содержат сообщение
			continue
		}

		UpdateHandlers(bot, update)
	}
}

func UpdateHandlers(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	text := update.Message.Text
	log.Printf("[%s] %s", update.Message.From.UserName, text)

	if len(text) < 1 {
		return
	}

	args := strings.Split(text, " ")
	// Handlers open list
	switch args[0] {
	case "/start":
		HandlerOnCommandStart(bot, update)
		return
	case "/token":
		if len(args) > 1 {
			if _, ok := whiteList[update.Message.From.ID]; ok {
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Вы уже ввели правильный токен."))
				return
			} else if args[1] == utilities.AccessToken {
				whiteList[update.Message.From.ID] = true
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Доступ разрешен."))
				return
			} else {
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Доступ запрещён, неверный token."))
				return
			}

		}
	}

	if _, ok := whiteList[update.Message.From.ID]; !ok {
		HandlerAccessDenied(bot, update)
		return
	}
	// Handlers private list
	switch text {
	case "Закрыть чат.":
		HandlerOnCommandCloseChat(bot, update)
		return
	default:
		HandlerDialogOpenAI(bot, update)
		return
	}
}

func LoopTtl(bot *tgbotapi.BotAPI) {
	for {
		nowUnix := time.Now().Unix()

		messagesPool.MessagesPool.LoopChats(
			func(chat *messagesPool.Chat) {
				if chat.Expired < nowUnix {
					fmt.Printf("\ncloseChat %d", chat.Id())
					msg := tgbotapi.NewMessage(chat.Id(), "Чат был закрыт.")
					bot.Send(msg)
					chat.Delete()
				}
			},
		)
		time.Sleep(loopTtlChat)
	}
}
