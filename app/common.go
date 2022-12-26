package app

import (
	"context"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/noamcattan/geni/ent"
	"github.com/noamcattan/geni/ent/user"
	"time"
)

func GetFirstAndLastOfMonth(t time.Time) (time.Time, time.Time) {
	currentYear, currentMonth, _ := t.Date()
	currentLocation := t.Location()

	firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1)

	return firstOfMonth, lastOfMonth
}

func GetCategoryKeyboard(ctx context.Context, client *ent.Client, chatId int64) (tg.ReplyKeyboardMarkup, error) {
	kb := tg.NewReplyKeyboard()

	categories, err := client.User.Query().
		Where(user.TelegramIDEQ(chatId)).
		QueryAccount().
		QueryCategories().
		All(ctx)

	if err != nil {
		return kb, err
	}

	for i := 0; i < len(categories); i++ {
		kb.Keyboard = append(kb.Keyboard, tg.NewKeyboardButtonRow(tg.NewKeyboardButton(categories[i].Name)))
	}

	return kb, nil
}
