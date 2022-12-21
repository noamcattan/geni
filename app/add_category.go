package app

import (
	"context"
	"fmt"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/noamcattan/geni/ent"
	"github.com/noamcattan/geni/ent/account"
	"github.com/noamcattan/geni/ent/user"
	"log"
	"strconv"
)

type AddCategoryData struct {
	Name  string
	Quota uint32
}

type AddCategoryConversation struct {
	client *ent.Client
	bot    *tg.BotAPI
	node   ConversationNode

	data AddCategoryData
}

func NewAddCategoryConversation(client *ent.Client, bot *tg.BotAPI) *AddCategoryConversation {
	c := &AddCategoryConversation{
		client: client,
		bot:    bot,
		data:   AddCategoryData{},
	}

	c.node = c.getCategoryName
	return c
}

func (c *AddCategoryConversation) Next(update tg.Update) {
	c.node = c.node(update)
}

func (c *AddCategoryConversation) Node() ConversationNode {
	return c.node
}

func (c *AddCategoryConversation) getCategoryName(update tg.Update) ConversationNode {
	msg := tg.NewMessage(update.Message.Chat.ID, "enter name")
	msg.ReplyMarkup = tg.NewRemoveKeyboard(true)
	_, _ = c.bot.Send(msg)

	return c.getCategoryQuota
}

func (c *AddCategoryConversation) getCategoryQuota(update tg.Update) ConversationNode {
	c.data.Name = update.Message.Text
	msg := tg.NewMessage(update.Message.Chat.ID, "enter quota")
	_, _ = c.bot.Send(msg)

	return c.end
}

func (c *AddCategoryConversation) end(update tg.Update) ConversationNode {
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

	category, err := tx.Category.Create().
		SetName(c.data.Name).
		SetQuota(c.data.Quota).Save(ctx)
	if err != nil {
		log.Println(err.Error())
		return nil
	}

	_, err = tx.Account.Update().
		Where(account.HasMemberWith(
			user.TelegramIDEQ(update.Message.Chat.ID)),
		).
		AddCategories(category).
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
		fmt.Sprintf("new cateory: %s (%d)", c.data.Name, c.data.Quota),
	))

	return nil
}
