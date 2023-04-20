package messagesPool

import (
	"chatGptTeleBot/ChatGpt"
	"time"
)

var refreshTtl = time.Minute * 20
var MessagesPool = Pool{app: make(map[int64]*Chat)}

type Pool struct {
	app map[int64]*Chat
}
type Chat struct {
	pool    *Pool
	chatId  int64
	Expired int64 // time as unix
	first   *Message
	last    *Message
}
type Message struct {
	before *Message
	after  *Message

	isBot   bool
	content string
}

func (m *Message) Delete() {
	if m.before != nil && m.after != nil {
		m.before.after = m.after
		m.after.before = m.before
		m.before = nil
		m.after = nil
	} else if m.before != nil {
		m.before.after = nil
		m.before = nil
	} else if m.after != nil {
		m.after.before = nil
		m.after = nil
	}
}

func (c *Chat) Id() int64 {
	return c.chatId
}
func (c *Chat) SendMsg(content string, isBot bool) *Message {
	if c.last != nil {
		msg := &Message{content: content, isBot: isBot, before: c.last}
		c.last.after = msg
		c.last = msg
		return msg
	} else {
		msg := &Message{content: content, isBot: isBot}
		c.last = msg
		c.first = msg
		return msg
	}
}

func (c *Chat) Delete() {
	last := c.last
	deleting := c.last
	c.first = nil
	c.last = nil
	delete(c.pool.app, c.chatId)

	for last.before != nil {
		deleting = last
		last = last.before
		deleting.Delete()
	}
	last.Delete()
}
func (c *Chat) Refresh() {
	c.Expired = time.Now().Add(refreshTtl).Unix()
}

func (m *Pool) DeleteChat(chatId int64) bool {

	chat, ok := m.app[chatId]
	if !ok {
		return false
	}
	chat.Delete()
	delete(m.app, chatId)
	return true
}

func (m *Pool) SendToChat(chatId int64, content string, isBot bool) {
	chat, exist := m.app[chatId]
	if !exist {
		chat = &Chat{pool: m, chatId: chatId}
		m.app[chatId] = chat
	}

	chat.Refresh()
	chat.SendMsg(content, isBot)
}
func (m *Pool) GetAllFromChat(chatId int64) []ChatGpt.OpenAIChatMessage {
	var messages = make([]ChatGpt.OpenAIChatMessage, 0)
	chat, exist := m.app[chatId]
	if !exist {
		return messages
	}

	msg := chat.first
	from := "user"
	if msg.isBot {
		from = "assistant"
	}
	messages = append(messages, ChatGpt.OpenAIChatMessage{Role: from, Content: msg.content})

	for msg.after != nil {
		msg = msg.after
		from = "user"
		if msg.isBot {
			from = "assistant"
		}
		messages = append(messages, ChatGpt.OpenAIChatMessage{Role: from, Content: msg.content})
	}
	return messages
}

func (m *Pool) LoopChats(loop func(chat *Chat)) {
	for _, chat := range m.app {
		loop(chat)
	}
}
