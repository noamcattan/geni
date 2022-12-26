package main

import (
	"context"
	"fmt"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	_ "github.com/mattn/go-sqlite3"
	"github.com/noamcattan/geni/app"
	"github.com/noamcattan/geni/ent"
	_ "github.com/noamcattan/geni/ent/runtime"
	"os"

	"log"
)

func main() {
	client, err := ent.Open("sqlite3", "file:/tmp/geni.sql?cache=shared&_fk=1")
	if err != nil {
		log.Panicf("failed opening connection to sqlite: %v", err)
	}
	defer client.Close()
	if err := client.Schema.Create(context.Background()); err != nil {
		log.Panicf("failed creating schema resources: %v", err)
	}
	updates := make(chan *tg.Update)

	server := newServer(client, updates)

	botToken, ok := os.LookupEnv("BOT_TOKEN")
	if !ok {
		log.Fatal("BOT_TOKEN env is missing")
	}

	bot, err := tg.NewBotAPI(botToken)
	if err != nil {
		log.Panic(err)
	}

	wh, _ := tg.NewWebhook("https: //meni-334119.ew.r.appspot.com/updates")
	_, err = bot.Request(wh)
	if err != nil {
		log.Fatal(err)
	}
	info, err := bot.GetWebhookInfo()
	if err != nil {
		log.Fatal(err)
	}

	if info.LastErrorDate != 0 {
		log.Printf("Telegram callback failed: %s", info.LastErrorMessage)
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	router := app.NewApp(bot, client, updates)
	ctx, cancel := context.WithCancel(context.Background())
	go router.Run(ctx)
	defer cancel()

	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "8080"
	}

	if err := server.Run(fmt.Sprintf(":%s", port)); err != nil {
		log.Printf(err.Error())
	}
}

