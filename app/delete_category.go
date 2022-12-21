package app

import (
	"context"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/noamcattan/geni/ent"
	"github.com/noamcattan/geni/ent/category"
	"github.com/noamcattan/geni/ent/user"
	"log"
)

type DeleteCategoryConversation struct {
	client *ent.Client
	bot    *tg.BotAPI
	node   ConversationNode
}

func NewDeleteCategoryConversation(client *ent.Client, bot *tg.BotAPI) *DeleteCategoryConversation {
	c := &DeleteCategoryConversation{
		client: client,
		bot:    bot,
	}

	c.node = c.getCategoryName
	return c
}

func (c *DeleteCategoryConversation) Next(update tg.Update) {
	c.node = c.node(update)
}

func (c *DeleteCategoryConversation) Node() ConversationNode {
	return c.node
}

func (c *DeleteCategoryConversation) getCategoryName(update tg.Update) ConversationNode {
	kb, err := GetCategoryKeyboard(context.Background(), c.client, update.Message.Chat.ID)
	if err != nil {
		log.Println(err.Error())
		return nil
	}

	if len(kb.Keyboard) == 0 {
		msg := tg.NewMessage(update.Message.Chat.ID, "no categories")
		msg.ReplyMarkup = tg.NewRemoveKeyboard(true)
		_, _ = c.bot.Send(msg)
		return nil
	}

	msg := tg.NewMessage(update.Message.Chat.ID, "choose category")
	msg.ReplyMarkup = kb
	_, _ = c.bot.Send(msg)

	return c.end
}

func (c *DeleteCategoryConversation) end(update tg.Update) ConversationNode {
	name := update.Message.Text
	ctx := context.Background()

	tx, err := c.client.Tx(ctx)
	if err != nil {
		log.Println(err.Error())
		return nil
	}

	ca, err := tx.User.Query().
		Where(user.TelegramIDEQ(update.Message.Chat.ID)).
		QueryAccount().
		QueryCategories().
		Where(category.NameEQ(name)).First(ctx)

	if err != nil {
		log.Println(err.Error())
		return nil
	}

	if err := tx.Category.DeleteOne(ca).Exec(ctx); err != nil {
		log.Println(err.Error())
		return nil
	}

	if err = tx.Commit(); err != nil {
		log.Println(err.Error())
	}

	msg := tg.NewMessage(update.Message.Chat.ID, "deleted")
	msg.ReplyMarkup = tg.NewRemoveKeyboard(true)
	_, _ = c.bot.Send(msg)

	return nil
}
