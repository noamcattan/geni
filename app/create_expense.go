package app

import (
	"context"
	"fmt"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/noamcattan/geni/ent"
	"github.com/noamcattan/geni/ent/account"
	"github.com/noamcattan/geni/ent/category"
	"github.com/noamcattan/geni/ent/expense"
	"github.com/noamcattan/geni/ent/user"
	"log"
	"strconv"
	"time"
)

type CreateExpenseData struct {
	Category    string
	Amount      float64
	Description string
}

type CreateExpenseConversation struct {
	client *ent.Client
	bot    *tg.BotAPI
	node   ConversationNode

	data CreateExpenseData
}

func NewCreateExpenseConversation(client *ent.Client, bot *tg.BotAPI) *CreateExpenseConversation {
	c := &CreateExpenseConversation{
		client: client,
		bot:    bot,
		data:   CreateExpenseData{},
	}

	c.node = c.getCategoryName
	return c
}

func (c *CreateExpenseConversation) Next(update tg.Update) {
	c.node = c.node(update)
}

func (c *CreateExpenseConversation) Node() ConversationNode {
	return c.node
}

func (c *CreateExpenseConversation) getCategoryName(update tg.Update) ConversationNode {
	kb, err := GetCategoryKeyboard(context.Background(), c.client, update.Message.Chat.ID)
	if err != nil {
		log.Println(err)
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

	return c.getAmount
}

func (c *CreateExpenseConversation) getAmount(update tg.Update) ConversationNode {
	c.data.Category = update.Message.Text

	msg := tg.NewMessage(update.Message.Chat.ID, "enter amount")
	msg.ReplyMarkup = tg.NewRemoveKeyboard(true)

	_, _ = c.bot.Send(msg)

	return c.getDescription
}

func (c *CreateExpenseConversation) getDescription(update tg.Update) ConversationNode {
	val, err := strconv.ParseFloat(update.Message.Text, 16)
	if err != nil {
		log.Println(err.Error())
		return nil
	}
	c.data.Amount = val

	msg := tg.NewMessage(update.Message.Chat.ID, "description:")

	_, _ = c.bot.Send(msg)

	return c.end
}

func (c *CreateExpenseConversation) end(update tg.Update) ConversationNode {
	c.data.Description = update.Message.Text

	ctx := context.Background()

	tx, err := c.client.Tx(ctx)
	if err != nil {
		log.Println(err.Error())
		return nil
	}

	// get the category
	ca, err := tx.Account.Query().
		Where(account.HasMemberWith(
			user.TelegramIDEQ(update.Message.Chat.ID)),
		).QueryCategories().
		Where(category.NameEQ(c.data.Category)).
		First(ctx)

	if err != nil {
		log.Println(err.Error())
		return nil
	}

	_, err = tx.Expense.Create().
		SetAmount(c.data.Amount).
		SetDescription(c.data.Description).
		AddCategory(ca).
		Save(ctx)
	if err != nil {
		log.Println(err.Error())
		return nil
	}

	firstOfMonth, lastOfMonth := GetFirstAndLastOfMonth(time.Now())

	total, err := tx.Account.Query().
		Where(account.HasMemberWith(user.TelegramIDEQ(update.Message.Chat.ID))).
		QueryCategories().
		Where(category.NameEQ(ca.Name)).
		QueryExpenses().
		Where(expense.And(
			expense.CreateTimeGTE(firstOfMonth),
			expense.CreateTimeLT(lastOfMonth),
		)).
		Aggregate(
			ent.Sum(expense.FieldAmount),
		).Float64(ctx)
	if err != nil {
		log.Println(err.Error())
		return nil
	}

	if err = tx.Commit(); err != nil {
		log.Println(err.Error())
		return nil
	}

	msg := tg.NewMessage(
		update.Message.Chat.ID,
		fmt.Sprintf("category: %s, amount: %.2f, description: %s\n left %.2f",
			c.data.Category,
			c.data.Amount,
			c.data.Description,
			float64(ca.Quota)-total),
	)
	tg.NewMessage(
		update.Message.Chat.ID,
		fmt.Sprintf("category: %s, amount: %.2f, description: %s\n left %.2f",
			c.data.Category,
			c.data.Amount,
			c.data.Description,
			float64(ca.Quota)-total),
	)
	_, _ = c.bot.Send(msg)

	return nil
}
