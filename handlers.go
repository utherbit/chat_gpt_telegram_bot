package main

import (
	"chatGptTeleBot/messagesPool"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var buttonMenu = tgbotapi.NewReplyKeyboard(tgbotapi.NewKeyboardButtonRow(
	tgbotapi.NewKeyboardButton("Закрыть чат.")))

func HandlerOnCommandCloseChat(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	ok := messagesPool.MessagesPool.DeleteChat(update.Message.Chat.ID)
	if ok {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Чат был закрыт.")
		bot.Send(msg)
	} else {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Чат уже был закрыт.")
		bot.Send(msg)
	}

}

func HandlerOnCommandStart(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf(
		`Привет, %s, я бот для общения с ChatGpt
Я выдал тебе меню управления, 
если она пропадет, просто напиши /start 

В меню ты можешь закрыть текущий диалог, 
что бы сменить тему, если ты не будешь долго писать, 
то диалог закроется автоматический.`,
		update.Message.From.UserName))
	msg.ReplyMarkup = buttonMenu
	bot.Send(msg)
}
func HandlerAccessDenied(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID,
		`Это приватный бот, доступ к функционалу ограничен. 
Если у вас есть токен доступа, введите команду
/token {ваш токен} 
для получения доступа.

Tелеграм разработчика: @utherbit
`)
	msg.ReplyMarkup = buttonMenu
	bot.Send(msg)
}
func HandlerDialogOpenAI(bot *tgbotapi.BotAPI, update tgbotapi.Update) {

	chatId := update.Message.Chat.ID
	messagesPool.MessagesPool.SendToChat(chatId, update.Message.Text, update.Message.From.IsBot)
	messages := messagesPool.MessagesPool.GetAllFromChat(chatId)

	text := "..."
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
	send, _ := bot.Send(msg)
	f := func(inp string) error {
		e := tgbotapi.NewEditMessageText(update.Message.Chat.ID, send.MessageID, inp)
		_, err := bot.Send(e)
		if err != nil {
			fmt.Printf("\nErr %v", err)
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, text)
			_, err = bot.Send(e)
			fmt.Printf("\nErr %v", err)
		}
		return nil
	}

	respMessage, err := ChatGptApp.StreamChatOpenAI(messages, f)
	if err != nil {
		panic(err)
	}

	messagesPool.MessagesPool.SendToChat(chatId, respMessage, true)

	//msg := tgbotapi.NewMessage(chatId, respMessage) // Создаем новое сообщение
	//msg.ReplyToMessageID = update.Message.MessageID // Устанавливаем идентификатор сообщения для ответа на оригинальное сообщение
	//_, err = bot.Send(msg)
}
