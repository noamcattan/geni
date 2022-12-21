package app

import (
	"context"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/noamcattan/geni/ent"
	"log"
)

type ConversationNode func(update tg.Update) ConversationNode

type Conversation interface {
	Next(update tg.Update)
	Node() ConversationNode
}

type App struct {
	bot    *tg.BotAPI
	client *ent.Client
	out    chan tg.MessageConfig

	conversations map[int64]Conversation
}

func NewApp(bot *tg.BotAPI, client *ent.Client) *App {
	return &App{
		bot:           bot,
		client:        client,
		out:           make(chan tg.MessageConfig),
		conversations: make(map[int64]Conversation),
	}
}

var commandStartKB = tg.NewReplyKeyboard(
	tg.NewKeyboardButtonRow(tg.NewKeyboardButton("report")),
	tg.NewKeyboardButtonRow(tg.NewKeyboardButton("delete")),
	tg.NewKeyboardButtonRow(tg.NewKeyboardButton("summary")),
	tg.NewKeyboardButtonRow(tg.NewKeyboardButton("add category")),
	tg.NewKeyboardButtonRow(tg.NewKeyboardButton("update category")),
	tg.NewKeyboardButtonRow(tg.NewKeyboardButton("delete category")),
)

func (a *App) chooseTopic(update tg.Update) {
	switch update.Message.Text {
	case "report":
		conv := NewCreateExpenseConversation(a.client, a.bot)
		a.conversations[update.Message.Chat.ID] = conv
		conv.Next(update)
		break
	case "delete":
		break
	case "summary":
		break
	case "add category":
		conv := NewAddCategoryConversation(a.client, a.bot)
		a.conversations[update.Message.Chat.ID] = conv
		conv.Next(update)
		break
	case "update category":
		conv := NewUpadteCategoryConversation(a.client, a.bot)
		a.conversations[update.Message.Chat.ID] = conv
		conv.Next(update)
		break
	case "delete category":
		conv := NewDeleteCategoryConversation(a.client, a.bot)
		a.conversations[update.Message.Chat.ID] = conv
		conv.Next(update)
		break
	default:
		msg := tg.NewMessage(update.Message.Chat.ID, "choose option")
		msg.ReplyMarkup = commandStartKB
		_, _ = a.bot.Send(msg)
	}
}

func (a *App) Run(ctx context.Context) {
	u := tg.NewUpdate(0)
	u.Timeout = 60
	updates := a.bot.GetUpdatesChan(u)

	for {
		select {
		case update := <-updates:
			chatId := update.Message.Chat.ID
			log.Printf("new update from: %s %s. user id=%d. message=%s",
				update.Message.Chat.FirstName,
				update.Message.Chat.LastName,
				chatId,
				update.Message.Text,
			)
			conv, ok := a.conversations[chatId]
			if !ok {
				a.chooseTopic(update)
				continue
			}

			if conv.Node() == nil {
				delete(a.conversations, chatId)
				continue
			}

			conv.Next(update)

			if conv.Node() == nil {
				delete(a.conversations, chatId)
				continue
			}

			a.conversations[chatId] = conv

		case <-ctx.Done():
			return
		}
	}
}
