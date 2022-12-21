package app

import (
	"context"
	"fmt"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/noamcattan/geni/ent"
	"github.com/noamcattan/geni/ent/category"
	"github.com/noamcattan/geni/ent/user"
	"log"
	"strconv"
)

type UpadteCategoryData struct {
	Name  string
	Quota uint32
}

type UpadteCategoryConversation struct {
	client *ent.Client
	bot    *tg.BotAPI
	node   ConversationNode

	data UpadteCategoryData
}

func NewUpadteCategoryConversation(client *ent.Client, bot *tg.BotAPI) *UpadteCategoryConversation {
	c := &UpadteCategoryConversation{
		client: client,
		bot:    bot,
		data:   UpadteCategoryData{},
	}

	c.node = c.getCategoryName
	return c
}

func (c *UpadteCategoryConversation) Next(update tg.Update) {
	c.node = c.node(update)
}

func (c *UpadteCategoryConversation) Node() ConversationNode {
	return c.node
}

func (c *UpadteCategoryConversation) getCategoryName(update tg.Update) ConversationNode {
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

	return c.getCategoryQuota
}

func (c *UpadteCategoryConversation) getCategoryQuota(update tg.Update) ConversationNode {
	c.data.Name = update.Message.Text
	msg := tg.NewMessage(update.Message.Chat.ID, "enter quota")
	_, _ = c.bot.Send(msg)

	return c.end
}

func (c *UpadteCategoryConversation) end(update tg.Update) ConversationNode {
	val, err := strconv.Atoi(update.Message.Text)
	if err != nil {
		log.Println(err.Error())
		return nil
	}
	c.data.Quota = uint32(val)

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
		Where(category.NameEQ(c.data.Name)).First(ctx)

	if err != nil {
		log.Println(err.Error())
		return nil
	}

	ca, err = ca.Update().
		SetQuota(c.data.Quota).
		Save(ctx)

	if err != nil {
		log.Println(err.Error())
		return nil
	}

	if err = tx.Commit(); err != nil {
		log.Println(err.Error())
		return nil
	}

	_, _ = c.bot.Send(tg.NewMessage(
		update.Message.Chat.ID,
		fmt.Sprintf("cateory updated: %s (%d)", c.data.Name, c.data.Quota),
	))

	return nil
}
